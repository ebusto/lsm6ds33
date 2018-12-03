package lsm6ds33

import (
	"encoding/binary"
	"io"
	"time"
)

const (
	I2cAddress = 0x6B
)

const (
	Id = 0x69
)

var order = binary.LittleEndian

type LSM6DS33 struct {
	cn  io.ReadWriter
	err error
}

type RawSample struct {
	Accel [3]int16
	Gyro  [3]int16
	Temp  int16
}

func New(cn io.ReadWriter) *LSM6DS33 {
	return &LSM6DS33{cn, nil}
}

func (d *LSM6DS33) Start() error {
	d.write(CTRL1_XL, 0x00)
	d.write(CTRL2_G, 0x00)
	d.write(CTRL3_C, 0x00)

	time.Sleep(time.Millisecond * 50)

	// Accelerometer

	// 0x80 = 0b10000000
	// ODR = 1000 (1.66 kHz (high performance)); FS_XL = 00 (+/-2 g full scale)
	d.write(CTRL1_XL, 0x80)

	// Gyro

	// 0x80 = 0b010000000
	// ODR = 1000 (1.66 kHz (high performance)); FS_XL = 00 (245 dps)
	d.write(CTRL2_G, 0x80)

	// Common

	// 0x04 = 0b00000100
	// IF_INC = 1 (automatically increment register address)
	d.write(CTRL3_C, 0x04)

	return d.err
}

func (d *LSM6DS33) Stop() error {
	// Stop accelerometer.
	d.write(CTRL1_XL, 0x00)

	// Stop gyroscope.
	d.write(CTRL2_G, 0x00)

	return d.err
}

func (d *LSM6DS33) Reset() error {
	if err := d.Stop(); err != nil {
		return err
	}

	return d.Start()
}

func (d *LSM6DS33) ReadId() (byte, error) {
	d.write(WHO_AM_I)

	var b [1]byte

	d.read(b[:])

	return b[0], d.err
}

func (d *LSM6DS33) ReadAccel() ([]int16, error) {
	d.write(OUTX_L_XL)

	b := make([]byte, 6)

	d.read(b)

	x := int16(order.Uint16(b[0:2]))
	y := int16(order.Uint16(b[2:4]))
	z := int16(order.Uint16(b[4:6]))

	return []int16{x, y, z}, d.err
}

func (d *LSM6DS33) ReadGyro() ([]int16, error) {
	d.write(OUTX_L_G)

	b := make([]byte, 6)

	d.read(b)

	x := int16(order.Uint16(b[0:2]))
	y := int16(order.Uint16(b[2:4]))
	z := int16(order.Uint16(b[4:6]))

	return []int16{x, y, z}, d.err
}

func (d *LSM6DS33) ReadTemp() (int16, error) {
	d.write(OUT_TEMP_L)

	b := make([]byte, 2)

	d.read(b)

	t := int16(order.Uint16(b[0:2]))
	t = 25 + t/16

	return t, d.err
}

func (d *LSM6DS33) read(b []byte) {
	if d.err != nil {
		return
	}

	_, d.err = d.cn.Read(b)
}

func (d *LSM6DS33) write(b ...byte) {
	if d.err != nil {
		return
	}

	_, d.err = d.cn.Write(b)
}
