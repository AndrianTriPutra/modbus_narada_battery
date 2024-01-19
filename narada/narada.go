package narada

import (
	"atp/modbus/domain"
	"context"
	"errors"
	"time"

	"github.com/simonvetter/modbus"
)

type strucT struct {
	setting Setting
}

type Setting struct {
	Port     string
	Baudrate uint
	Timeout  time.Duration
}

func NewRepository(setting Setting) RepositoryI {
	return &strucT{
		setting: setting,
	}
}

type RepositoryI interface {
	Modbus(ctx context.Context, id uint8) (battery domain.Battery, err error)
}

func (m strucT) Modbus(ctx context.Context, id uint8) (battery domain.Battery, err error) {
	battery.Pack_ID = id

	url := "rtu://" + m.setting.Port
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      url,
		Speed:    m.setting.Baudrate,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  m.setting.Timeout,
	})
	if err != nil {
		errN := errors.New("E0->" + err.Error())
		return battery, errN
	}

	err = client.Open()
	if err != nil {
		errN := errors.New("E1->" + err.Error())
		return battery, errN
	}
	defer client.Close()

	client.SetUnitId(38 + id)

	//datas, err := client.ReadRegisters(0x0FFF, 2, modbus.INPUT_REGISTER)
	datas, err := client.ReadRegisters(0x0FFF, 2, modbus.INPUT_REGISTER)
	for i, data := range datas {
		switch i {
		case 0:
			battery.Data.Voltage = float32(uint16(data)) * 0.01
		case 1:
			battery.Data.Current = (float32(int16(data)) - 10000) * 0.1
		}
	}
	if err != nil {
		errN := errors.New("TO1->" + err.Error())
		return battery, errN
	}

	//datas, err = client.ReadRegisters(0x1007, 3, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(0x1007, 3, modbus.INPUT_REGISTER)
	for i, data := range datas {
		switch i {
		case 0:
			battery.Data.SoC = float32(uint16(data)) * 0.01
		case 2:
			battery.Data.SoH = float32(uint16(data)) * 0.01
		}
	}
	if err != nil {
		errN := errors.New("TO2->" + err.Error())
		return battery, errN
	}

	//datas, err = client.ReadRegisters(0x102F, 2, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(0x102F, 2, modbus.INPUT_REGISTER)
	if err != nil {
		errN := errors.New("TO3->" + err.Error())
		return battery, errN
	}
	for i, data := range datas {
		switch i {
		case 0:
			battery.Data.Rem_Charge_Time = data
		case 1:
			battery.Data.Rem_Discharge_Time = data
		}
	}

	return battery, nil
}
