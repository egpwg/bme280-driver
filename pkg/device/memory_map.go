package device

const (
	DEVICE_ADDR = 0x77
)

type UserMode int

const (
	Weather UserMode = iota + 1
	HumiSensing
	Indoor
	Gaming
)

type Oversampling uint8

// register ctrl_hum(0xF2) and ctrl_meas(0xF4)
// set osrs_h(bit0~bit2) / osrs_p(bit2~bit4) / osrs_t(bit5~bit7)[2:0] and contains 3 bit
const (
	Skipped        Oversampling = 0x00 // skipped
	Oversampling1  Oversampling = 0x01 // oversampling * 1
	Oversampling2  Oversampling = 0x02 // oversampling * 2
	Oversampling4  Oversampling = 0x03 // oversampling * 4
	Oversampling8  Oversampling = 0x04 // oversampling * 8
	Oversampling16 Oversampling = 0x05 // oversampling * 16
)

func (m UserMode) GetOversampling() (os map[string]Oversampling) {
	switch m {
	case Weather:
		os["Temperature"] = Oversampling1
		os["Pressure"] = Oversampling1
		os["Humidity"] = Oversampling1
	case HumiSensing:
	case Indoor:
	case Gaming:
	}

	return
}

type SensorMode uint8

// register ctrl_meas(0xF4)
// set mode(bit0~bit1)[1:0] and contains 2 bit
const (
	Sleep  SensorMode = 0x00 // 00
	Forced SensorMode = 0x01 // 01 and 10
	Normal SensorMode = 0x03 // 11
)

type TimeStandby uint8

// register config(0xF5)
// set t_sb(bit5~bit7)[2:0] and contains 3 bit  / spi3w_en(bit0)[0]
const (
	TSb0point5  TimeStandby = 0x00 // 000: 0.5ms
	TSb62point5 TimeStandby = 0x01 // 001: 62.5ms
	TSb125      TimeStandby = 0x02 // 010: 125ms
	TSb250      TimeStandby = 0x03 // 011: 250ms
	TSb500      TimeStandby = 0x04 // 100: 500ms
	TSb1000     TimeStandby = 0x05 // 101: 1000ms
	TSb10       TimeStandby = 0x06 // 110: 10ms
	TSb20       TimeStandby = 0x07 // 111: 20ms
)

// FilterCoef filter coefficient
type FilterCoef uint8

// register config(0xF5)
// set filter(bit2~bit4)[2:0] and contains 3 bit
const (
	FilterCoefOff FilterCoef = 0x00 // 000: off
	FilterCoef2   FilterCoef = 0x01 // 001: 2
	FilterCoef4   FilterCoef = 0x02 // 010: 4
	FilterCoef8   FilterCoef = 0x03 // 011: 8
	FilterCoef16  FilterCoef = 0x04 // 100 and others: 16
)

// Calibration compensation parameter storage
type Calibration struct {
	DigT1                                                  uint16
	DigT2, DigT3                                           int16
	DigP1                                                  uint16
	DigP2, DigP3, DigP4, DigP5, DigP6, DigP7, DigP8, DigP9 int16
	DigH1, DigH3                                           uint8
	DigH2, DigH4, DigH5                                    int16
	DigH6                                                  int8
}

func newCalibration(tph, h []byte) (c Calibration) {
	c.DigT1 = uint16(tph[0]) | uint16(tph[1])<<8
	c.DigT2 = int16(tph[2]) | int16(tph[3])<<8
	c.DigT3 = int16(tph[4]) | int16(tph[5])<<8
	c.DigP1 = uint16(tph[6]) | uint16(tph[7])<<8
	c.DigP2 = int16(tph[8]) | int16(tph[9])<<8
	c.DigP3 = int16(tph[10]) | int16(tph[11])<<8
	c.DigP4 = int16(tph[12]) | int16(tph[13])<<8
	c.DigP5 = int16(tph[14]) | int16(tph[15])<<8
	c.DigP6 = int16(tph[16]) | int16(tph[17])<<8
	c.DigP7 = int16(tph[18]) | int16(tph[19])<<8
	c.DigP8 = int16(tph[20]) | int16(tph[21])<<8
	c.DigP9 = int16(tph[22]) | int16(tph[23])<<8
	c.DigH1 = uint8(tph[25])

	c.DigH2 = int16(h[0]) | int16(h[1])<<8
	c.DigH3 = uint8(h[2])
	c.DigH4 = int16(h[3])<<4 | int16(h[4])&0xF
	c.DigH5 = int16(h[4])>>4 | int16(h[5])<<4
	c.DigH6 = int8(h[6])

	return c
}

func (c *Calibration) CompensateTemperatureInt32(adcT int32) (tFine, T int32) {
	var t1, t2 int32
	t1 = (adcT>>3 - int32(c.DigT1)<<1) * int32(c.DigT2) >> 11
	t2 = (((((adcT >> 4) - (int32(c.DigT1))) * ((adcT >> 4) - (int32(c.DigT1)))) >> 12) *
		(int32(c.DigT3))) >> 14
	tFine = t1 + t2
	T = (tFine*5 + 128) >> 8
	return tFine, T
}

func (c *Calibration) CompensatePressureInt64(tFine, adcP int32) (P uint32) {
	var p, p1, p2 int64
	p1 = int64(tFine) - 128000
	p2 = p1 * p1 * int64(c.DigP6)
	p2 = p2 + p1*int64(c.DigP5)<<17
	p2 = p2 + int64(c.DigP4)<<35
	p1 = p1*p1*int64(c.DigP3)>>8 + p1*int64(c.DigP2)<<12
	p1 = (int64(1)<<47 + p1) * int64(c.DigP1) >> 33
	if p1 == 0 {
		return 0
	}

	p = ((1048576-int64(adcP))<<31 - p2) * 3125 / p1
	p1 = int64(c.DigP9) * (p >> 13) * (p >> 13) >> 25
	p2 = int64(c.DigP8) * p >> 19

	return uint32((p+p1+p2)>>8 + int64(c.DigP7)<<4)
}

func (c *Calibration) CompensateHumidityInt32(tFine, adcH int32) (H uint32) {
	var h int32
	h = tFine - int32(76800)
	h1 := (adcH<<14 - int32(c.DigH4)<<20 - int32(c.DigH5)*h + int32(16384)) >> 15
	h2 := h * int32(c.DigH6) >> 10
	h3 := h*int32(c.DigH3)>>11 + 32768
	h4 := h2*h3>>10 + 2097152
	h5 := (h4*int32(c.DigH2) + 8192) >> 14
	h = h1 * h5
	h = h - h>>15*(h>>15)>>7*int32(c.DigH1)>>4
	if h < 0 {
		h = 0
	}
	if h > 419430400 {
		h = 419430400
	}

	return uint32(h) >> 12
}

// register address
const (
	regAddrCalib1   = 0x88
	regAddrChipID   = 0xD0
	regAddrReset    = 0xE0
	regAddrCalib2   = 0xE1
	regAddrCtrlHum  = 0xF2
	regAddrStatus   = 0xF3
	regAddrCtrlMeas = 0xF4
	regAddrConfig   = 0xF5
	regAddrPress    = 0xF7
	regAddrTemp     = 0xFA
	regAddrHumi     = 0xFD
)
