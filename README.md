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


| AMQP channel      | IN/OUT | Content-Type                  | Data         | Description                                     |
| ----------------- | ------ | ----------------------------- | ------------ | ----------------------------------------------- |
| bbdcmotors_ctrl   | IN     | application/dcmotor_forward   | uint32 [0-3] | Set motor to forward state (CW)                 |
| bbdcmotors_ctrl   | IN     | application/dcmotor_backward  | uint32 [0-3] | Set motor to backward state (CCW)               |
| bbdcmotors_ctrl   | IN     | application/dcmotor_stop      | uint32 [0-3] | Set motor to stopped state                      |
| bbdcmotors_ctrl   | IN     | application/dcmotor_speed     | uint32       | Change all motors speed (default duty is 1000). |
| bbdcmotors_events | OUT    | --                            | --           | Unused                                          |