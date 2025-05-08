package sumsub

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func GettingDocumentImages(inspectionId, imageId string) (output io.ReadCloser, err error) {
	path := fmt.Sprintf("/resources/inspections/%s/resources/%s", inspectionId, imageId)

	url := fmt.Sprintf("%s%s", configs.C.Sumsub.ApiUrl, path)

	header, err := buildHeader(http.MethodGet, path, nil)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for k, v := range header {
		request.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	return res.Body, nil
}

var fetchSumsubLock *sync.Mutex
var fetchSumsubLockInitialFlag sync.Once

func FetchDataAdapter(ctx context.Context, log *logger.Logger, repo domain.Repository, externalUserId string) (idNumber string, err error) {
	fetchSumsubLockInitialFlag.Do(func() {
		fetchSumsubLock = new(sync.Mutex)
	})

	// cache lock
	fetchSumsubLock.Lock()
	defer fetchSumsubLock.Unlock()

	customerID, err := uuid.Parse(externalUserId)
	if err != nil {
		log.Error("uuid.Parse(externalUserId)", zap.Any("customerID", externalUserId), zap.Any("err", err))
		return "", errors.BadRequest(error_code.ErrInvalidRequest, "Invalid UUID format").WithCause(err)
	}

	// fetch sumsub application info
	data, err := GetApplicantInfo(customerID.String())
	if err != nil {
		log.Error("sumsub.GetApplicantInfo", zap.Any("customerID", customerID), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.GetApplicantInfo").WithCause(err)
	}

	// 抓取 sumsub id number
	if len(data.Info.IdDocs) > 0 {
		for _, v := range data.Info.IdDocs {
			if v.IdDocType == "ID_CARD" {
				idNumber = v.Number
				break
			}
		}
	}

	// 當狀態為 init 或是根本沒有身分證時，log 下來 回傳 sumsub success
	if data.Review.ReviewStatus == "init" || idNumber == "" {
		log.Error("data.Review.ReviewStatus == \"init\" || idNumber == \"\"", zap.Any("customerID", customerID), zap.Any("idNumber", idNumber))
		return "", nil
	}

	// fetch id docs
	idDocs, err := GetApplicantRequiredIdDocs(data.Id, data.InspectionId)
	if err != nil {
		log.Error("sumsub.GetApplicantRequiredIdDocs", zap.Any("customerID", customerID), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.GetApplicantRequiredIdDocs").WithCause(err)
	}

	// 機敏資訊加密
	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(data)", zap.Any("data", data), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "json.Marshal").WithCause(err)
	}
	encryptedInfo, err := tools.EncryptAES256(configs.C.SensitiveDataEncryptKey, string(dataByte))
	if err != nil {
		log.Error("tools.EncryptAES256", zap.Any("data", data), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "tools.EncryptAES256").WithCause(err)
	}
	hashEmail, err := tools.HashSensitiveData(configs.C.SensitiveDataEncryptKey, data.Email)
	if err != nil {
		log.Error("tools.HashSensitiveData", zap.Any("data", data), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "tools.HashSensitiveData").WithCause(err)
	}

	insertData := domain.BTMSumsub{
		ApplicantId:      data.Id,
		CustomerId:       customerID,
		InfoHash:         encryptedInfo,
		IdNumber:         idNumber,
		BanExpireDate:    sql.NullInt64{},
		EmailHash:        hashEmail,
		Phone:            data.Phone,
		InspectionId:     data.InspectionId,
		IdCardFrontImgId: idDocs.IdCardFrontImgId,
		IdCardBackImgId:  idDocs.IdCardBackImgId,
		SelfieImgId:      idDocs.SelfieImgId,
		Name:             fmt.Sprintf("%v%v", data.Info.LastName, data.Info.FirstName),
		Status:           data.Review.ReviewResult.ReviewAnswer,
	}

	// 檢查告誡名單
	exist, fetchExpireDate, err := repo.IsBTMCIBExist(repo.GetDb(ctx), idNumber)
	if err != nil {
		log.Error("repo.IsBTMCIBExist", zap.Any("customerID", customerID), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.IsBTMCIBExist").WithCause(err)
	}
	if exist {
		insertData.BanExpireDate = sql.NullInt64{
			Int64: fetchExpireDate,
			Valid: true,
		}
	}

	// store db
	err = repo.UpsertBTMSumsub(repo.GetDb(ctx), insertData)
	if err != nil {
		log.Error("repo.UpsertBTMSumsub", zap.Any("customerID", customerID), zap.Any("err", err))
		return "", errors.InternalServer(error_code.ErrBTMSumsubCreateItem, "repo.UpsertBTMSumsub").WithCause(err)
	}

	return idNumber, nil
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

// AddAndOverwriteApplicantTags 添加或覆寫申請人標籤
func AddAndOverwriteApplicantTags(applicantId string, tags []string) (err error) {
	path := fmt.Sprintf("/resources/applicants/%s/tags", applicantId)

	_, err = sumsubSend[domain.SumsubDataResApplicant](http.MethodPost, path, tags)
	if err != nil {
		return
	}

	return nil
}

// GetApplicantRequiredIdDocs 取得申請人需要提交的身份證件
func GetApplicantRequiredIdDocs(applicantId string, inspectionId string) (respApplicant domain.SumsubImageDos, err error) {
	path := fmt.Sprintf("/resources/applicants/%s/requiredIdDocsStatus", applicantId)

	resData, err := sumsubSend[domain.SumsubDataApplicantRequiredIdDocs](http.MethodGet, path, nil)
	if err != nil {
		return
	}

	if len(resData.IDENTITY.ImageIds) == 2 {
		respApplicant.IdCardFrontImgId = fmt.Sprintf("%v", resData.IDENTITY.ImageIds[0])
		respApplicant.IdCardBackImgId = fmt.Sprintf("%v", resData.IDENTITY.ImageIds[1])
	}

	if len(resData.SELFIE.ImageIds) > 0 {
		var selfieImages []string
		for _, v := range resData.SELFIE.ImageIds {
			selfieImages = append(selfieImages, fmt.Sprintf("%v", v))
		}
		selfieImageToString, err := json.Marshal(selfieImages)
		if err != nil {
			err = errors.InternalServer(error_code.ErrSumsubApiUnmarshal, "json.Marshal").WithCause(err)
			return respApplicant, err
		}
		respApplicant.SelfieImgId = string(selfieImageToString)
	}

	return respApplicant, nil
}

// GetApplicantReviewHistory 取得申請人審核歷史
func GetApplicantReviewHistory(applicantId string) (history domain.SumsubHistoryReviewData, err error) {
	path := fmt.Sprintf("/resources/applicantTimeline/%s", applicantId)

	history, err = sumsubSend[domain.SumsubHistoryReviewData](http.MethodGet, path, nil)
	if err != nil {
		return
	}

	return
}
