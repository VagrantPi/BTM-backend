package customer

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GetCustomerRiskControlRoleReq struct {
	CustomerId string `uri:"customer_id"`
}

type GetBTMUserInfoDetailRes struct {
	ApplicantId      string   `json:"applicant_id"`
	Name             string   `json:"name"`
	IDNumber         string   `json:"id_number"`
	Birthday         string   `json:"birthday"`
	Address          string   `json:"address"`
	Occupation       string   `json:"occupation"`
	PhoneNumber      string   `json:"phone_number"`
	Email            string   `json:"email"`
	InspectionId     string   `json:"inspection_id"`
	IdCardFrontImgId string   `json:"id_card_front_img_id"`
	IdCardBackImgId  string   `json:"id_card_back_img_id"`
	SelfieImgIds     []string `json:"selfie_img_ids"`
	ReviewStatus     string   `json:"review_status"`
	ReviewCreateDate string   `json:"review_create_date"`
	SuspendedUntil   string   `json:"suspended_until"`
}

func GetBTMUserInfoDetail(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMUserInfoDetail")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetCustomerRiskControlRoleReq{}
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Error("c.ShouldBindUri(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.ShouldBindUri(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.ShouldBindUri(&req)").WithCause(err))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	customerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		log.Error("uuid.Parse(req.CustomerId)", zap.Any("err", err))
		api.ErrResponse(c, "uuid.Parse(req.CustomerId)", errors.BadRequest(error_code.ErrInvalidRequest, "uuid.Parse(req.CustomerId)").WithCause(err))
		return
	}

	sumsubInfo, err := repo.GetBTMSumsub(repo.GetDb(c), req.CustomerId)
	if err != nil {
		log.Error("repo.GetBTMSumsub()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetBTMSumsub()", errors.InternalServer(error_code.ErrInternalError, "repo.GetBTMSumsub()").WithCause(err))
		return
	}

	customer, err := repo.GetCustomerById(repo.GetDb(c), customerId)
	if err != nil {
		log.Error("repo.GetCustomerById()", zap.Any("err", err))
		api.ErrResponse(c, "repo.GetCustomerById()", errors.InternalServer(error_code.ErrInternalError, "repo.GetCustomerById()").WithCause(err))
		return
	}

	if sumsubInfo == nil {
		log.Error("repo.GetBTMSumsub()", zap.Any("sumsubInfo", sumsubInfo))
		api.ErrResponse(c, "repo.GetBTMSumsub()", errors.NotFound(error_code.ErrBTMSumsubGetItem, "repo.GetBTMSumsub()").WithMetadata(map[string]string{
			"customerId": req.CustomerId,
		}))
		return
	}

	var info domain.SumsubData
	if sumsubInfo.InfoHash != "" {
		// decrypt
		decryptedInfo, err := tools.DecryptAES256(configs.C.SensitiveDataEncryptKey, sumsubInfo.InfoHash)
		if err != nil {
			log.Error("tools.DecryptAES256", zap.Any("err", err))
			api.ErrResponse(c, "tools.DecryptAES256", errors.BadRequest(error_code.ErrToolsHashSensitiveData, "tools.DecryptAES256").WithCause(err))
			return
		}

		err = json.Unmarshal([]byte(decryptedInfo), &info)
		if err != nil {
			log.Error("json.Unmarshal", zap.Any("err", err))
			api.ErrResponse(c, "json.Unmarshal", errors.BadRequest(error_code.ErrToolsHashSensitiveData, "json.Unmarshal").WithCause(err))
			return
		}
	}

	// 地址取得
	address := ""
	if len(info.Info.Addresses) > 0 {
		address = info.Info.Addresses[0].FormattedAddress
	}

	occupation := ""
	if len(info.Questionnaires) > 0 {
		for _, v := range info.Questionnaires {
			if v.Id == "newQuestionnaire" {
				occupation = domain.KYCDocOccupation[v.Sections.JiBenZiXun.Items.NinDeZhiYe.Value]
				break
			}
		}
	}

	selfieImgIds := []string{}
	err = json.Unmarshal([]byte(sumsubInfo.SelfieImgId), &selfieImgIds)
	if err != nil {
		log.Error("json.Unmarshal", zap.Any("err", err))
		api.ErrResponse(c, "json.Unmarshal", errors.BadRequest(error_code.ErrToolsHashSensitiveData, "json.Unmarshal").WithCause(err))
		return
	}

	api.OKResponse(c, GetBTMUserInfoDetailRes{
		ApplicantId:      sumsubInfo.ApplicantId,
		Name:             sumsubInfo.Name,
		IDNumber:         sumsubInfo.IdNumber,
		Birthday:         info.Info.Dob,
		Address:          address,
		Occupation:       occupation,
		PhoneNumber:      info.Phone,
		Email:            info.Email,
		InspectionId:     sumsubInfo.InspectionId,
		IdCardFrontImgId: sumsubInfo.IdCardFrontImgId,
		IdCardBackImgId:  sumsubInfo.IdCardBackImgId,
		SelfieImgIds:     selfieImgIds,
		ReviewStatus:     info.Review.ReviewStatus,
		ReviewCreateDate: info.Review.CreateDate,
		SuspendedUntil:   customer.SuspendedUntil,
	})
}
