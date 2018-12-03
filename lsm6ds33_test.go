package lsm6ds33

import (
	"testing"

	"gobot.io/x/gobot/platforms/raspi"
)

func TestLSM6DS33(t *testing.T) {
	cn, err := raspi.NewAdaptor().GetConnection(I2cAddress, 1)

	if err != nil {
		t.Fatal(err)
	}

	d := New(cn)

	d.Start()

	defer d.Stop()

	id, err := d.ReadId()

	if err != nil {
		t.Fatal(err)
	}

	temp, err := d.ReadTemp()

	if err != nil {
		t.Fatal(err)
	}

	accel, err := d.ReadAccel()

	if err != nil {
		t.Fatal(err)
	}

	gyro, err := d.ReadGyro()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("[%d] T = %v, A = %v, G = %v", id, temp, accel, gyro)
}
