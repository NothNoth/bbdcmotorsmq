package main

import (
	"encoding/binary"
	"encoding/json"

	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	bbhw "github.com/btittelbach/go-bbhw"
	i2c "github.com/d2r2/go-i2c"
	"github.com/streadway/amqp"
)

const defaultWait = 10 * time.Millisecond
const defaultDCFrequency = 30
const defaultDuty = 1000
const directionForward = 1
const directionBackward = 2
const (
	exchangeCtrl   = "bbdcmotors_ctrl"
	exchangeEvents = "bbdcmotors_events"
)

type BBDCMotorsConfig struct {
	I2CAddress byte
	I2CLane    int
	GpioPin    uint
	RmqServer  string
}

type BBDCMotorsMQ struct {
	config    BBDCMotorsConfig
	ctrl      *bbhw.MMappedGPIO
	i2c       *i2c.I2C
	killed    bool
	conn      *amqp.Connection
	ch        *amqp.Channel
	ctrlQueue amqp.Queue
	speedDuty uint32
}

func InitBBDCMotorsMQ(configFile string) (*BBDCMotorsMQ, error) {
	var mmq BBDCMotorsMQ
	var err error
	//Load config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &mmq.config)
	if err != nil {
		return nil, err
	}

	//Setup GPIO for BB Motor bridge
	mmq.ctrl = bbhw.NewMMappedGPIO(mmq.config.GpioPin, bbhw.OUT)
	err = mmq.ctrl.SetState(true)
	if err != nil {
		mmq.ctrl.Close()
		return nil, err
	}
	time.Sleep(defaultWait)

	mmq.i2c, err = i2c.NewI2C(mmq.config.I2CAddress, mmq.config.I2CLane)
	if err != nil {
		return nil, nil
	}
	time.Sleep(defaultWait)

	//Setup AMQP
	mmq.conn, err = amqp.Dial(mmq.config.RmqServer)
	if err != nil {
		return nil, err
	}

	mmq.ch, err = mmq.conn.Channel()
	if err != nil {
		return nil, err
	}

	//Setup Control exchange & queue
	err = mmq.ch.ExchangeDeclare(
		exchangeCtrl, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	mmq.ctrlQueue, err = mmq.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	//Bind this queue to this exchange so that exchange will publish here
	err = mmq.ch.QueueBind(
		mmq.ctrlQueue.Name, // queue name
		"",                 // routing key
		exchangeCtrl,       // exchange
		false,
		nil)

	//Setup events exchange
	err = mmq.ch.ExchangeDeclare(
		exchangeEvents, // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}
	mmq.killed = false
	mmq.speedDuty = defaultDuty

	mmq.EnableDC(1, true)
	mmq.EnableDC(2, true)
	mmq.EnableDC(3, true)
	mmq.EnableDC(4, true)

	return &mmq, nil
}

func (mmq *BBDCMotorsMQ) Destroy() {
	mmq.ctrl.Close()
	mmq.i2c.Close()
}

func (mmq *BBDCMotorsMQ) ReceiveCommands() error {
	msgs, err := mmq.ch.Consume(
		mmq.ctrlQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			switch d.ContentType {
			case "application/dcmotor_forward":
				motorID := binary.BigEndian.Uint32(d.Body)
				err := mmq.MoveDC(motorID, directionForward, mmq.speedDuty)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("DC motor %d forward", motorID)
				}
			case "application/dcmotor_backward":
				motorID := binary.BigEndian.Uint32(d.Body)
				err := mmq.MoveDC(motorID, directionBackward, mmq.speedDuty)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("DC motor %d backward", motorID)
				}
			case "application/dcmotor_stop":
				motorID := binary.BigEndian.Uint32(d.Body)
				err := mmq.StopDC(motorID)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("DC motor %d stopped", motorID)
				}
			case "application/dcmotor_speed":
				speedDuty := binary.BigEndian.Uint32(d.Body)
				mmq.speedDuty = speedDuty
				err := mmq.ChangeSpeedDC(1, mmq.speedDuty)
				if err != nil {
					log.Println(err)
				}
				mmq.ChangeSpeedDC(2, mmq.speedDuty)
				mmq.ChangeSpeedDC(3, mmq.speedDuty)
				mmq.ChangeSpeedDC(4, mmq.speedDuty)
				log.Printf("Changed speed to %d", mmq.speedDuty)
			default:
				log.Printf("Received unexpected message: %s", d.Body)
			}
		}
	}()

	return nil
}

func (mmq *BBDCMotorsMQ) EmitEvents() error {
	for {
		if mmq.killed == true {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (mmq *BBDCMotorsMQ) writeWord(reg byte, value uint32) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE)       // Read/Write ?
	byteSeq = append(byteSeq, reg)              //Which register ?
	byteSeq = append(byteSeq, byte(value&0xFF)) //32 bits value
	byteSeq = append(byteSeq, byte((value>>8)&0xFF))
	byteSeq = append(byteSeq, byte((value>>16)&0xFF))
	byteSeq = append(byteSeq, byte((value>>24)&0xFF))

	fmt.Println("Write word:")
	for _, z := range byteSeq {
		fmt.Printf("%02x", z)
	}
	fmt.Println("")

	_, err := mmq.i2c.WriteBytes(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (mmq *BBDCMotorsMQ) writeHalfWord(reg byte, value uint16) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE)       // Read/Write ?
	byteSeq = append(byteSeq, reg)              //Which register ?
	byteSeq = append(byteSeq, byte(value&0xFF)) //16 bits value
	byteSeq = append(byteSeq, byte((value>>8)&0xFF))

	fmt.Println("Write halfword:")
	for _, z := range byteSeq {
		fmt.Printf("%02x", z)
	}
	fmt.Println("")
	_, err := mmq.i2c.WriteBytes(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (mmq *BBDCMotorsMQ) writeByte(reg byte, value byte) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE) // Read/Write ?
	byteSeq = append(byteSeq, reg)        //Which register ?
	byteSeq = append(byteSeq, value)      //8 bits value
	fmt.Println("Write byte:")
	for _, z := range byteSeq {
		fmt.Printf("%02x", z)
	}
	fmt.Println("")
	_, err := mmq.i2c.WriteBytes(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func getDCRegisters(dc uint32) (mode byte, direction byte, duty byte, err error) {
	switch dc {
	case 1:
		mode = TB_1A_MODE
		direction = TB_1A_DIR
		duty = TB_1A_DUTY
		break
	case 2:
		mode = TB_1B_MODE
		direction = TB_1B_DIR
		duty = TB_1B_DUTY
		break
	case 3:
		mode = TB_2A_MODE
		direction = TB_2A_DIR
		duty = TB_2A_DUTY
		break
	case 4:
		mode = TB_2B_MODE
		direction = TB_2B_DIR
		duty = TB_2B_DUTY
		break
	default:
		mode = 0
		direction = 0
		duty = 0
		err = errors.New("Invalid dc id (1-2)")
	}
	return
}

func (mmq *BBDCMotorsMQ) EnableDC(dc uint32, enable bool) error {
	if (dc <= 0) || (dc > 4) {
		return errors.New("Invalid motor ID")
	}

	modeReg, directionReg, _, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mmq.writeWord(CONFIG_TB_PWM_FREQ, mmq.speedDuty)
	time.Sleep(defaultWait)
	mmq.writeByte(modeReg, TB_DCM)
	time.Sleep(defaultWait)
	mmq.writeByte(directionReg, TB_STOP)
	time.Sleep(defaultWait)

	return nil
}

func (mmq *BBDCMotorsMQ) ChangeSpeedDC(dc uint32, duty uint32) error {
	if (dc <= 0) || (dc > 4) {
		return errors.New("Invalid motor ID (1-4)")
	}
	if (duty <= 0) || (duty > 100) {
		return errors.New("Invalid speed (1-100)")
	}
	_, _, dutyReg, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mmq.writeWord(dutyReg, duty*10)
	time.Sleep(defaultWait)
	return nil
}

func (mmq *BBDCMotorsMQ) MoveDC(dc uint32, direction byte, duty uint32) error {

	if (dc <= 0) || (dc > 4) {
		return errors.New("Invalid motor ID")
	}

	_, directionReg, dutyReg, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mmq.writeByte(directionReg, direction)
	time.Sleep(defaultWait)
	mmq.writeWord(dutyReg, duty*10)
	time.Sleep(defaultWait)
	return nil
}

func (mmq *BBDCMotorsMQ) StopDC(dc uint32) error {

	if (dc <= 0) || (dc > 4) {
		return errors.New("Invalid motor ID")
	}

	_, directionReg, _, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mmq.writeByte(directionReg, TB_STOP)
	time.Sleep(defaultWait)
	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <config file>")
		return
	}
	mmq, err := InitBBDCMotorsMQ(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to init BBDCMotorsMQ: %s", err)
		return
	}
	defer mmq.Destroy()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			mmq.killed = true
		}
	}()

	mmq.ReceiveCommands()
	mmq.EmitEvents()

	mmq.Destroy()
}
