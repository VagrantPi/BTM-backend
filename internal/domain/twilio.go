package domain

type TwilioHistorySms struct {
	DateSent *string `json:"date_sent"`
	Status   *string `json:"status"`
	Body     *string `json:"body"`
}
