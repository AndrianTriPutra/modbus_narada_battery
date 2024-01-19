package domain

type Battery struct {
	Pack_ID uint8
	Data    Narada
}

type Narada struct {
	Current            float32
	Voltage            float32
	SoC                float32
	SoH                float32
	Rem_Charge_Time    uint16
	Rem_Discharge_Time uint16
}
