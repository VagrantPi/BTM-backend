package domain

type UserJwt struct {
	Account string `json:"account"`
	Role    int64  `json:"role"`
	Id      uint   `json:"id"`
}
