package tuntap

import "testing"

func TestNew(t *testing.T) {
	c := Config{
		Type:       TUN,
		Name:       "tun1",
		Persist:    false,
		MultiQueue: false,
	}

	device, err := New(c)
	if err != nil {
		t.Fatal(err)
	}

	err = device.Close()
	if err != nil {
		t.Fatal(err)
	}
}
