package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashInTx struct {
	ID                   uuid.UUID       `json:"id"`
	DeviceID             string          `json:"deviceId"`
	ToAddress            string          `json:"toAddress"`
	CryptoAtoms          decimal.Decimal `json:"cryptoAtoms"`
	CryptoCode           string          `json:"cryptoCode"`
	Fiat                 decimal.Decimal `json:"fiat"`
	FiatCode             string          `json:"fiatCode"`
	Fee                  int64           `json:"fee"`
	TxHash               string          `json:"txHash"`
	Phone                string          `json:"phone"`
	Error                string          `json:"error"`
	Created              time.Time       `json:"created"`
	Send                 bool            `json:"send"`
	SendConfirmed        bool            `json:"sendConfirmed"`
	Timedout             bool            `json:"timedout"`
	SendTime             time.Time       `json:"sendTime"`
	ErrorCode            string          `json:"errorCode"`
	OperatorCompleted    bool            `json:"operatorCompleted"`
	SendPending          bool            `json:"sendPending"`
	CashInFee            decimal.Decimal `json:"cashInFee"`
	MinimumTx            int             `json:"minimumTx"`
	CustomerID           uuid.UUID       `json:"customerId"`
	TxVersion            int             `json:"txVersion"`
	TermsAccepted        bool            `json:"termsAccepted"`
	CommissionPercentage decimal.Decimal `json:"commissionPercentage"`
	RawTickerPrice       decimal.Decimal `json:"rawTickerPrice"`
	IsPaperWallet        bool            `json:"isPaperWallet"`
	Discount             int16           `json:"discount"`
	BatchID              uuid.UUID       `json:"batchId"`
	Batched              bool            `json:"batched"`
	BatchTime            time.Time       `json:"batchTime"`
	DiscountSource       string          `json:"discountSource"`
	TxCustomerPhotoAt    time.Time       `json:"txCustomerPhotoAt"`
	TxCustomerPhotoPath  string          `json:"txCustomerPhotoPath"`
	WalletScore          int16           `json:"walletScore"`
	Email                string          `json:"email"`
}

type CashInTxWithInfo struct {
	CashInTx
	DeviceName string `json:"deviceName"`
	InvoiceNo  string `json:"invoiceNo"`
}
