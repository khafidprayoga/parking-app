package types

type AppStatus struct {
	Revenue            float64 `json:"revenue"`
	LotParkingCapacity int     `json:"area_capacity"`
	TxCount            int     `json:"tx_count"`
	CarList            []*Car  `json:"car_list"`
}
