package main

const READ_MODE = 0
const WRITE_MODE = 1

// TB_WORKMODE

const TB_SHORT_BREAK = 0
const TB_CW = 1
const TB_CCW = 2
const TB_STOP = 3
const TB_WORKMODE_NUM = 4

// TB_PORTMODE

const TB_DCM = 0
const TB_SPM = 1
const TB_PORTMODE_NUM = 2

// SVM_PORT

const SVM1 = 0
const SVM2 = 1
const SVM3 = 2
const SVM4 = 3
const SVM5 = 4
const SVM6 = 5
const SVM_PORT_NUM = 6

// SVM_STATE

const SVM_DISABLE = 0
const SVM_ENABLE = 1
const SVM_STATE_NUM = 2

// IO_MODE

const IO_IN = 0
const IO_OUT = 1
const IO_MODE_NUM = 2

// IO_PUPD

const IO_PU = 0
const IO_PD = 1
const IO_NP = 2
const IO_PUPD_NUM = 3

// IO_PPOD

const IO_PP = 0
const IO_OD = 1
const IO_PPOD_NUM = 2

// IO_STATE

const IO_LOW = 0
const IO_HIGH = 1
const IO_STATE_NUM = 2

// IO_PORT

const IO1 = 0
const IO2 = 1
const IO3 = 2
const IO4 = 3
const IO5 = 4
const IO6 = 5
const IO_NUM = 6

// PARAM_REG

const CONFIG_VALID = 0
const CONFIG_TB_PWM_FREQ = CONFIG_VALID + 4

const I2C_ADDRESS = CONFIG_TB_PWM_FREQ + 4

const TB_1A_MODE = I2C_ADDRESS + 1
const TB_1A_DIR = TB_1A_MODE + 1
const TB_1A_DUTY = TB_1A_DIR + 1
const TB_1A_SPM_SPEED = TB_1A_DUTY + 2
const TB_1A_SPM_STEP = TB_1A_SPM_SPEED + 4

const TB_1B_MODE = TB_1A_SPM_STEP + 4
const TB_1B_DIR = TB_1B_MODE + 1
const TB_1B_DUTY = TB_1B_DIR + 1
const TB_1B_SPM_SPEED = TB_1B_DUTY + 2
const TB_1B_SPM_STEP = TB_1B_SPM_SPEED + 4

const TB_2A_MODE = TB_1B_SPM_STEP + 4
const TB_2A_DIR = TB_2A_MODE + 1
const TB_2A_DUTY = TB_2A_DIR + 1
const TB_2A_SPM_SPEED = TB_2A_DUTY + 2
const TB_2A_SPM_STEP = TB_2A_SPM_SPEED + 4

const TB_2B_MODE = TB_2A_SPM_STEP + 4
const TB_2B_DIR = TB_2B_MODE + 1
const TB_2B_DUTY = TB_2B_DIR + 1
const TB_2B_SPM_SPEED = TB_2B_DUTY + 2
const TB_2B_SPM_STEP = TB_2B_SPM_SPEED + 4

const SVM1_STATE = TB_2B_SPM_STEP + 4
const SVM1_FREQ = SVM1_STATE + 1
const SVM1_ANGLE = SVM1_FREQ + 2

const SVM2_STATE = SVM1_ANGLE + 2
const SVM2_FREQ = SVM2_STATE + 1
const SVM2_ANGLE = SVM2_FREQ + 2

const SVM3_STATE = SVM2_ANGLE + 2
const SVM3_FREQ = SVM3_STATE + 1
const SVM3_ANGLE = SVM3_FREQ + 2

const SVM4_STATE = SVM3_ANGLE + 2
const SVM4_FREQ = SVM4_STATE + 1
const SVM4_ANGLE = SVM4_FREQ + 2

const SVM5_STATE = SVM4_ANGLE + 2
const SVM5_FREQ = SVM5_STATE + 1
const SVM5_ANGLE = SVM5_FREQ + 2

const SVM6_STATE = SVM5_ANGLE + 2
const SVM6_FREQ = SVM6_STATE + 1
const SVM6_ANGLE = SVM6_FREQ + 2

const IO1_STATE = SVM6_ANGLE + 2
const IO1_MODE = IO1_STATE + 1
const IO1_PUPD = IO1_MODE + 1
const IO1_PPOD = IO1_PUPD + 1

const IO2_STATE = IO1_PPOD + 1
const IO2_MODE = IO2_STATE + 1
const IO2_PUPD = IO2_MODE + 1
const IO2_PPOD = IO2_PUPD + 1

const IO3_STATE = IO2_PPOD + 1
const IO3_MODE = IO3_STATE + 1
const IO3_PUPD = IO3_MODE + 1
const IO3_PPOD = IO3_PUPD + 1

const IO4_STATE = IO3_PPOD + 1
const IO4_MODE = IO4_STATE + 1
const IO4_PUPD = IO4_MODE + 1
const IO4_PPOD = IO4_PUPD + 1

const IO5_STATE = IO4_PPOD + 1
const IO5_MODE = IO5_STATE + 1
const IO5_PUPD = IO5_MODE + 1
const IO5_PPOD = IO5_PUPD + 1

const IO6_STATE = IO5_PPOD + 1
const IO6_MODE = IO6_STATE + 1
const IO6_PUPD = IO6_MODE + 1
const IO6_PPOD = IO6_PUPD + 1

const PARAM_REG_NUM = IO6_PPOD + 1
