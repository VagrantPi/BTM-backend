package model

type BTM_CIB struct {
	DataType    string `json:"datatype"`
	Pid         string `json:"pid" gorm:"unique;not null"`
	WarningDate int64  `json:"warningdate"`
	ExpireDate  int64  `json:"expiredate"`
	Issuer      string `json:"issuer"`
	Blank1      string `json:"blank_1"`
	Blank2      string `json:"blank_2"`
}
