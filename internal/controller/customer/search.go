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
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SearchCustomersReq struct {
	Query              string              `form:"query"`
	Address            string              `form:"address"`
	Email              string              `form:"email"`
	Name               string              `form:"name"`
	WhiteListDateStart time.Time           `form:"white_list_date_start"`
	WhiteListDateEnd   time.Time           `form:"white_list_date_end"`
	CustomerDateStart  time.Time           `form:"customer_date_start"`
	CustomerDateEnd    time.Time           `form:"customer_date_end"`
	CustomerType       domain.CustomerType `form:"customer_type"`
	Active             bool                `form:"active"`
	Limit              int                 `form:"limit"`
	Page               int                 `form:"page"`
}

type SearchCustomersRep struct {
	Total int                      `json:"total"`
	Items []SearchCustomersRepItem `json:"items"`
}

type SearchCustomersRepItem struct {
	ID                    uuid.UUID `json:"id"`
	Phone                 string    `json:"phone"`
	Email                 string    `json:"email"`
	Name                  string    `json:"name"`
	Created               time.Time `json:"created_at"`
	FirstWhiteListCreated time.Time `json:"first_white_list_created"`
	IsLamassuBlock        bool      `json:"is_lamassu_block"`
	IsAdminBlock          bool      `json:"is_admin_block"`
	IsCibBlock            bool      `json:"is_cib_block"`
}

func SearchCustomers(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "SearchCustomers")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := SearchCustomersReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	cid, _ := uuid.Parse(req.Query)
	phone, _ := strconv.Atoi(req.Query)
	_customerId := ""
	if cid != uuid.Nil || phone == 0 {
		// 現在搜尋條件是customer_id
		_customerId = req.Query
		req.Query = ""
	}

	reqEmailHash := ""
	if req.Email != "" {
		reqEmailHash, err = tools.HashSensitiveData(configs.C.SensitiveDataEncryptKey, req.Email)
		if err != nil {
			log.Error("tools.HashSensitiveData", zap.Any("err", err))
			api.ErrResponse(c, "tools.HashSensitiveData", errors.BadRequest(error_code.ErrToolsHashSensitiveData, "tools.HashSensitiveData").WithCause(err))
			return
		}
	}

	customers, total, err := repo.SearchCustomers(repo.GetDb(c), req.Query, _customerId, req.Address, reqEmailHash, req.Name,
		req.WhiteListDateStart, req.WhiteListDateEnd, req.CustomerDateStart, req.CustomerDateEnd,
		req.CustomerType, req.Active,
		req.Limit, req.Page)
	if err != nil {
		log.Error("repo.SearchCustomersByCustomerId", zap.Any("err", err))
		api.ErrResponse(c, "repo.SearchCustomersByCustomerId", errors.NotFound(error_code.ErrDBError, "repo.SearchCustomersByCustomerId()").WithCause(err))
		return
	}

	resp := SearchCustomersRep{
		Total: total,
		Items: make([]SearchCustomersRepItem, len(customers)),
	}

	for i, v := range customers {
		var info domain.SumsubData
		if v.InfoHash != "" {
			// decrypt
			decryptedInfo, err := tools.DecryptAES256(configs.C.SensitiveDataEncryptKey, v.InfoHash)
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

		resp.Items[i] = SearchCustomersRepItem{
			ID:                    v.ID,
			Phone:                 v.Phone,
			Email:                 tools.MaskEmail(info.Email),
			Name:                  tools.MaskName(v.Name),
			Created:               v.Created,
			FirstWhiteListCreated: v.FirstWhiteListCreated,
			IsLamassuBlock:        v.IsLamassuBlock,
			IsAdminBlock:          v.IsAdminBlock,
			IsCibBlock:            v.IsCibBlock,
		}
	}

	api.OKResponse(c, resp)
}
