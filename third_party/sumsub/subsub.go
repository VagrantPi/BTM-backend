package sumsub

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/tools"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

func sumsubSend[T any](method string, path string, body any) (output T, err error) {
	url := fmt.Sprintf("%s%s", configs.C.Sumsub.ApiUrl, path)

	header, err := buildHeader(method, path, body)
	if err != nil {
		return
	}

	statusCode, res, _, err := tools.JsonRequest(method, url, header, body)
	if err != nil {
		err = errors.InternalServer(error_code.ErrSumsubRequest, "tools.JsonRequest").WithCause(err).
			WithMetadata(map[string]string{"method": method, "url": configs.C.Sumsub.ApiUrl, "path": path})
		return
	} else if statusCode >= 400 {
		err = errors.InternalServer(error_code.ErrSumsubBadRequest, "tools.JsonRequest: statusCode >= 400").
			WithMetadata(map[string]string{
				"url":    configs.C.Sumsub.ApiUrl,
				"path":   path,
				"status": strconv.Itoa(statusCode),
				"res":    string(res),
			})
		return
	}

	if err = json.Unmarshal(res, &output); err != nil {
		err = errors.InternalServer(error_code.ErrSumsubApiUnmarshal, "json.Unmarshal").WithCause(err)
		return
	}

	return
}

func buildHeader(method string, path string, body any) (header map[string]string, err error) {
	var bodyByte []byte
	var ts = fmt.Sprintf("%d", time.Now().Unix())

	if body != nil {
		bodyByte, err = json.Marshal(body)
		if err != nil {
			err = errors.InternalServer(error_code.ErrSumsubApiUnmarshal, "json.Marshal").WithCause(err).WithMetadata(map[string]string{
				"method": method, "path": path,
			})
			return
		}
	}

	sign := sumsubSign(ts, configs.C.Sumsub.AppSecret, method, path, bodyByte)

	header = map[string]string{
		"X-App-Token":      configs.C.Sumsub.AppToken,
		"X-App-Access-Sig": sign,
		"X-App-Access-Ts":  ts,
		"Accept":           "application/json",
	}

	return
}

func sumsubSign(ts string, secret string, method string, path string, body []byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	data := []byte(ts + method + path)

	if body != nil {
		data = append(data, body...)
	}

	hash.Write(data)

	return hex.EncodeToString(hash.Sum(nil))
}

// GetApplicantInfo 取得申請人資訊
func GetApplicantInfo(externalUserId string) (applicant domain.SumsubData, err error) {
	path := fmt.Sprintf("/resources/applicants/-;externalUserId=%s/one", externalUserId)

	applicant, err = sumsubSend[domain.SumsubData](http.MethodGet, path, nil)
	if err != nil {
		return
	}

	return
}

func AddAndOverwriteApplicantTags(applicantId string, tags []string) (err error) {
	path := fmt.Sprintf("/resources/applicants/%s/tags", applicantId)

	_, err = sumsubSend[domain.SumsubDataResApplicant](http.MethodPost, path, tags)
	if err != nil {
		return
	}

	return nil
}
