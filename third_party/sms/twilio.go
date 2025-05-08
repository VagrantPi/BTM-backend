package sms

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(phone, sessionId, timeAt, fiat, crypto, cryptoSymbol, rate, txId string) error {
	url := "https://api.twilio.com/2010-04-01/Accounts/AC6507a79bd1b4a1de37475e8fb87b52b9/Messages.json"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`To=%v&From=+14067630496&Body=RECEIPT
Operator information:
CoinNow
support@coin-now.com
02-23688839
Session: %v
Time: %v
Direction: Cash-in
Fiat: %v TWD
Crypto: %v %v
Rate: 1 %v = %v TWD
TXID: %v
`, phone, sessionId, timeAt, fiat, crypto, cryptoSymbol, cryptoSymbol, rate, txId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	username := configs.C.Twilio.Username
	password := configs.C.Twilio.Password
	auth := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+encoded)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func GetSmsHistory(dateSent time.Time, limit int) ([]domain.TwilioHistorySms, error) {
	from := "+14067630496"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: configs.C.Twilio.Username,
		Password: configs.C.Twilio.Password,
	})

	params := &twilioApi.ListMessageParams{}
	params.SetFrom(from)
	params.SetDateSent(dateSent)
	params.SetLimit(limit)
	resp, err := client.Api.ListMessage(params)
	if err != nil {
		return nil, err
	}
	messages := make([]domain.TwilioHistorySms, len(resp))
	for i := len(resp) - 1; i >= 0; i-- {
		messages[len(resp)-i-1] = domain.TwilioHistorySms{
			DateSent: resp[i].DateSent,
			Status:   resp[i].Status,
			Body:     resp[i].Body,
		}
	}
	return messages, nil
}
