package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type UserConfig struct {
	ID      int64          `json:"id"`
	Type    string         `json:"type"`
	Data    UserConfigJSON `json:"data"`
	Created time.Time      `json:"created"`
}

type UserConfigJSON struct {
	Config Config `json:"config"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *UserConfigJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := UserConfigJSON{}
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (j UserConfigJSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Config struct {
	TriggersConfigExpirationTime            string   `json:"triggersConfig_expirationTime"`
	TriggersConfigAutomation                string   `json:"triggersConfig_automation"`
	LocaleTimezone                          string   `json:"locale_timezone"`
	CashInCashboxReset                      string   `json:"cashIn_cashboxReset"`
	NotificationsEmailSecurity              bool     `json:"notifications_email_security"`
	NotificationsSmsSecurity                bool     `json:"notifications_sms_security"`
	NotificationsNotificationCenterSecurity bool     `json:"notifications_notificationCenter_security"`
	WalletsAdvancedFeeMultiplier            string   `json:"wallets_advanced_feeMultiplier"`
	WalletsAdvancedCryptoUnits              string   `json:"wallets_advanced_cryptoUnits"`
	WalletsAdvancedAllowTransactionBatching bool     `json:"wallets_advanced_allowTransactionBatching"`
	WalletsAdvancedId                       string   `json:"wallets_advanced_id"`
	TriggersConfigCustomerAuthentication    string   `json:"triggersConfig_customerAuthentication"`
	WalletsBTCZeroConfLimit                 int      `json:"wallets_BTC_zeroConfLimit"`
	WalletsBTCCoin                          string   `json:"wallets_BTC_coin"`
	WalletsBTCWallet                        string   `json:"wallets_BTC_wallet"`
	WalletsBTCTicker                        string   `json:"wallets_BTC_ticker"`
	WalletsBTCExchange                      string   `json:"wallets_BTC_exchange"`
	WalletsBTCZeroConf                      string   `json:"wallets_BTC_zeroConf"`
	LocaleId                                string   `json:"locale_id"`
	LocaleCountry                           string   `json:"locale_country"`
	LocaleFiatCurrency                      string   `json:"locale_fiatCurrency"`
	LocaleLanguages                         []string `json:"locale_languages"`
	LocaleCryptoCurrencies                  []string `json:"locale_cryptoCurrencies"`
	CommissionsMinimumTx                    int      `json:"commissions_minimumTx"`
	CommissionsFixedFee                     int      `json:"commissions_fixedFee"`
	CommissionsCashOut                      int      `json:"commissions_cashOut"`
	CommissionsCashIn                       int      `json:"commissions_cashIn"`
	CommissionsId                           string   `json:"commissions_id"`
	WalletsETHActive                        bool     `json:"wallets_ETH_active"`
	WalletsETHTicker                        string   `json:"wallets_ETH_ticker"`
	WalletsETHWallet                        string   `json:"wallets_ETH_wallet"`
	WalletsETHExchange                      string   `json:"wallets_ETH_exchange"`
	WalletsETHZeroConf                      string   `json:"wallets_ETH_zeroConf"`
	Triggers                                []struct {
		Requirement         string      `json:"requirement"`
		SuspensionDays      interface{} `json:"suspensionDays"`
		Threshold           int         `json:"threshold"`
		ThresholdDays       interface{} `json:"thresholdDays"`
		CustomInfoRequestId string      `json:"customInfoRequestId"`
		ExternalService     string      `json:"externalService"`
		Id                  string      `json:"id"`
		Direction           string      `json:"direction"`
		TriggerType         string      `json:"triggerType"`
	} `json:"triggers"`
	WalletsUSDTTRONActive   bool   `json:"wallets_USDT_TRON_active"`
	WalletsUSDTTRONTicker   string `json:"wallets_USDT_TRON_ticker"`
	WalletsUSDTTRONWallet   string `json:"wallets_USDT_TRON_wallet"`
	WalletsUSDTTRONExchange string `json:"wallets_USDT_TRON_exchange"`
	WalletsUSDTTRONZeroConf string `json:"wallets_USDT_TRON_zeroConf"`
}
