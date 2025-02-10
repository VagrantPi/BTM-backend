package cib

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/jszwec/csvutil"
	"go.uber.org/zap"
)

func GetToken() (string, error) {
	log := logger.Zap().WithClassFunction("cib", "GetToken")
	defer func() {
		_ = log.Sync()
	}()

	targetUrl := "https://dawexchange.cib.gov.tw/WarningListServlet"

	config := configs.NewConfigs()

	data := url.Values{}
	data.Set("func", "getToken")
	data.Set("account", config.Cib.Account)
	data.Set("pwd", config.Cib.Pwd)

	client := &http.Client{}
	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(data.Encode()))

	if err != nil {
		log.Error("http.NewRequest", zap.Any("err", err))
		err = errors.InternalServer(error_code.ErrThirdPartyHttpCall, "http.NewRequest").
			WithCause(err).
			WithMetadata(map[string]string{
				"url": targetUrl,
			})
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Error("client.Do", zap.Any("err", err), zap.String("url", targetUrl), zap.Any("res", res))
		err = errors.InternalServer(error_code.ErrThirdPartyHttpCall, "client.Do").
			WithCause(err).
			WithMetadata(map[string]string{
				"url": targetUrl,
			})
		return "", err
	}
	defer res.Body.Close()

	resp := domain.CibToken{}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Error("json.NewDecoder.Decode", zap.Any("err", err), zap.String("url", targetUrl), zap.Any("res", res))
		err = errors.InternalServer(error_code.ErrCibTokenParse, "json.NewDecoder.Decode").
			WithCause(err).
			WithMetadata(map[string]string{
				"url": targetUrl,
			})
		return "", err
	}

	if resp.Success == false {
		log.Error("resp.Success is false", zap.Any("resp", resp))
		err = errors.InternalServer(error_code.ErrCibTokenFetch, "resp.Success is false").
			WithCause(err).
			WithMetadata(map[string]string{
				"url": targetUrl,
				"msg": resp.Message,
			})
		return "", err
	}
	return resp.Token, nil
}

// TODO:
func GetWarningZip(token string, destFile string) error {
	log := logger.Zap().WithClassFunction("cib", "GetWarningZip")
	defer func() {
		_ = log.Sync()
	}()

	targetUrl := "https://dawexchange.cib.gov.tw/WarningListServlet"

	fileDate := time.Now().Format("2006-01-02")
	data := url.Values{}
	data.Set("func", "getWarningListByToken")
	data.Set("token", token)
	data.Set("fileDate", fileDate)

	client := &http.Client{}
	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func ConvertCsvFileToBTMCIB(file io.Reader) ([]domain.BTMCIB, error) {
	// 讀取匯款資訊
	csvReader := csv.NewReader(file)
	decoder, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsCsvOpen, "open csv file err").WithCause(err)
		return nil, err
	}

	var results []domain.BTMCIB
	for {
		var item domain.BTMCIB

		if err1 := decoder.Decode(&item); err1 == io.EOF {
			break
		} else if err1 != nil {
			err = errors.InternalServer(error_code.ErrToolsCsvRead, "read csv row err").WithCause(err)
			return nil, err
		}

		results = append(results, item)
	}
	return results, nil
}
