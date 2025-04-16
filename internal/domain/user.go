package domain

type UserJwt struct {
	Account string `json:"account"`
	Role    int64  `json:"role"`
	Id      int64  `json:"id"`
}
