package models

type Toj struct {
	JenisPekerjaan string `json:"jenis_pekerjaan"`
	Total          int    `json:"total"`
}

type Gender struct {
	LakiLaki           float64 `json:"laki_laki"`
	Perempuan          float64 `json:"perempuan"`
	TotalLkSebelumnya  float64 `json:"total_lk_sebelumnya"`
	TotalPrSebelumnya  float64 `json:"total_pr_sebelumnya"`
	PersenLk           float64 `json:"persen_lk"`
	PersenPr           float64 `json:"persen_pr"`
	PersenLkSebelumnya float64 `json:"persen_lk_sebelumnya"`
	PersenPrSebelumnya float64 `json:"persen_pr_sebelumnya"`
}

type TypeMotorStat struct {
	Type            string `json:"type"`
	Total           int    `json:"total"`
	TotalSebelumnya int    `json:"total_sebelumnya"`
}

type JenisPekerjaan struct {
	Pekerjaan string `json:"pekerjaan"`
	Total     int    `json:"total"`
}
