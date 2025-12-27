package models

type Toj struct {
	JenisPekerjaan string `json:"jenis_pekerjaan"`
	Total          int    `json:"total"`
}

type TypeMotorStat struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}

type JenisPekerjaan struct {
	Pekerjaan string `json:"pekerjaan"`
	Total     int    `json:"total"`
}
