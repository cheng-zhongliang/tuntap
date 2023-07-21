package tuntap

import "testing"

func TestNew(t *testing.T) {
	c := Config{
		Type:    TUN,
		Name:    "test",
		Persist: true,
	}

	device, err := New(c)
	if err != nil {
		t.Fatal(err)
	}

	device.Close()
}
