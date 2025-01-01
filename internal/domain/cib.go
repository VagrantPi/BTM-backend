package domain

type CibToken struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}
