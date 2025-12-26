package models

type TypeMotorStat struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}

type JenisPekerjaan struct {
	Pekerjaan string `json:"pekerjaan"`
	Total     int    `json:"total"`
}
