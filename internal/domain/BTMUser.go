package domain

type BTMUser struct {
	Id       uint
	Account  string
	Password string
	Roles    int64
}

type BTMUserWithRoles struct {
	Id       uint   `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
	RoleId   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
}
