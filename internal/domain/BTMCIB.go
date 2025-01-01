package domain

type BTMCIB struct {
	DataType    string `csv:"DATATYPE" json:"datatype"`
	Pid         string `csv:"PID" json:"pid" gorm:"unique;not null"`
	WarningDate int64  `csv:"WARNINGDATE" json:"warningdate"`
	ExpireDate  int64  `csv:"EXPIREDATE" json:"expiredate"`
	Issuer      string `csv:"ISSUER" json:"issuer"`
	Blank1      string `csv:"BLANK_1" json:"blank_1"`
	Blank2      string `csv:"BLANK_2" json:"blank_2"`
}
