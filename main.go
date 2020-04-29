package hidrelay

import (
	"fmt"
	"runtime"

	"github.com/spetr/hid"
)

type (
	Relay struct {
		info *hid.DeviceInfo
		dev  *hid.Device
	}
	IoStatus      int
	ChannelNumber int
	ChannelStatus struct {
		Channel_1 IoStatus
		Channel_2 IoStatus
		Channel_3 IoStatus
		Channel_4 IoStatus
		Channel_5 IoStatus
		Channel_6 IoStatus
		Channel_7 IoStatus
		Channel_8 IoStatus
	}
)

// List returns all installed USB relay devices
func List() (list []*Relay) {
	list = make([]*Relay, 0)
	for _, devInfo := range hid.Enumerate(0x16C0, 0x05DF) {
		relay, err := devInfo.Open()
		if err != nil {
			continue
		}
		list = append(list, &Relay{info: &devInfo})
		relay.Close()
	}
	return
}

// Open creates HID connection to relay
func (r *Relay) Open() (err error) {
	r.dev, err = r.info.Open()
	return
}

// Close destroys HID connection to relay
func (r *Relay) Close() error {
	return r.dev.Close()
}

// Set switch selected channel to IoStatus (ON/OFF)
func (r *Relay) Set(no ChannelNumber, s IoStatus) error {
	cmd := make([]byte, 9)
	cmd[0] = 0x0
	if no < 0 || no > C8 {
		return fmt.Errorf("Illegal channel number (%d)", no)
	}
	if no == ALL {
		if s == ON {
			cmd[1] = 0xFE
		} else {
			cmd[1] = 0xFC
		}
	} else {
		if s == ON {
			cmd[1] = 0xFF
		} else {
			cmd[1] = 0xFD
		}
		cmd[2] = byte(no)
	}
	_, err := r.dev.SendFeatureReport(cmd)
	return err
}

// SetOn switch selected channel to ON
func (r *Relay) SetOn(num ChannelNumber) error {
	return r.Set(num, ON)
}

// SetAllOn turns all relay to ON
func (r *Relay) SetAllOn() error {
	return r.Set(ALL, ON)
}

// SetOff turns selected channel to OFF
func (r *Relay) SetOff(num ChannelNumber) error {
	return r.Set(num, OFF)
}

// SetAllOff turns all relay to OFF
func (r *Relay) SetAllOff() error {
	return r.Set(ALL, OFF)
}

// GetAll returns relay status (ON/OFF)
func (r *Relay) GetAll() (status *ChannelStatus, err error) {
	cmd := make([]byte, 9)
	_, err = r.dev.GetFeatureReport(cmd)
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		cmd = cmd[1:]
	}
	status = &ChannelStatus{
		Channel_1: IoStatus(cmd[7] >> 0 & 0x01),
		Channel_2: IoStatus(cmd[7] >> 1 & 0x01),
		Channel_3: IoStatus(cmd[7] >> 2 & 0x01),
		Channel_4: IoStatus(cmd[7] >> 3 & 0x01),
		Channel_5: IoStatus(cmd[7] >> 4 & 0x01),
		Channel_6: IoStatus(cmd[7] >> 5 & 0x01),
		Channel_7: IoStatus(cmd[7] >> 6 & 0x01),
		Channel_8: IoStatus(cmd[7] >> 7 & 0x01),
	}
	return
}

// SetSN writes serial number (max 5 characters)
func (r *Relay) SetSN(sn string) (err error) {
	if len(sn) > 5 {
		err = fmt.Errorf("The length of '%s' is large than 5 bytes", sn)
		return
	}
	cmd := make([]byte, 9)
	cmd[0] = 0x00
	cmd[1] = 0xFA
	copy(cmd[2:], sn)
	_, err = r.dev.SendFeatureReport(cmd)
	return
}

// GetSN reads serial number (max 5 characters)
func (r *Relay) GetSN() (sn string, err error) {
	cmd := make([]byte, 9)
	_, err = r.dev.GetFeatureReport(cmd)
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		cmd = cmd[1:]
	}
	sn = string(cmd[:5])
	return
}
