# BB DC Motor

Standalone service to control beaglebone's DCMotors shield through rabbit MQ

## Config

    {
      "I2CAddress": 75, <--- address of the I2C bus (0x4B)
      "I2CLane": 2,     <--- I2C Lane
      "GpioPin": 49,     <--- GPIO pin of the motorsbridge (49 on BBG green)
      "RmqServer": "amqp://guest:guest@localhost:5672/" <--- address of the AMQP server
    }

By default all motors are enabled and stopped.

## AMQP


| AMQP exchange     | IN/OUT | Content-Type                           | Data                  | Description                                              |
| ----------------- | ------ | -------------------------------------- | --------------------- | -------------------------------------------------------- |
| bbdcmotors_ctrl   | IN     | application/dcmotor_forward            | uint32 [1-4]          | Set motor to forward state (CW)                          |
| bbdcmotors_ctrl   | IN     | application/dcmotor_backward           | uint32 [1-4]          | Set motor to backward state (CCW)                        |
| bbdcmotors_ctrl   | IN     | application/dcmotor_forward_for_ticks  | uint32 [1-4] uint32 x | Set motor to forward state (CW) for x ticks              |
| bbdcmotors_ctrl   | IN     | application/dcmotor_backward_for_ticks | uint32 [1-4] uint32 x | Set motor to backward state (CCW) for x ticks            |
| bbdcmotors_ctrl   | IN     | application/dcmotor_stop               | uint32 [1-4]          | Set motor to stopped state                               |
| bbdcmotors_ctrl   | IN     | application/dcmotor_speed              | uint32 [1-100]        | Change all motors speed in % (default duty is 30).       |
| bbdcmotors_events | OUT    | application/dcmotor_ticks_per_rotation | uint32                | Number of ticks per rotation (sent at regular intervals) |
