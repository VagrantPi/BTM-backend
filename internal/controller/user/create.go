package user

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type CreateOneReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int64  `json:"role" binding:"required"`
}

func CreateOne(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "CreateOne")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req CreateOneReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	if !tools.CheckPasswordRule(req.Password) {
		log.Error("password not match regex", zap.Any("password", req.Password))
		api.ErrResponse(c, "password not match regex", errors.BadRequest(error_code.ErrInvalidRequest, "密碼需要包含大小寫字母、數字和特殊字符"))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	tx, err := repo.NewTransactionBegin(c)
	if err != nil {
		log.Error("NewTransactionBegin error", zap.Any("err", err))
		api.ErrResponse(c, "NewTransactionBegin error", errors.InternalServer(error_code.ErrInternalError, "NewTransactionBegin error").WithCause(err))
		return
	}
	defer repo.TransactionCommit(tx)

	// check role exist
	roleAny, ok := c.Get("role")
	if !ok {
		roleAny, err = repo.GetRawRoleById(tx, req.Role)
		if err != nil {
			log.Error("GetRawRoleById error", zap.Any("err", err))
			api.ErrResponse(c, "GetRawRoleById error", errors.InternalServer(error_code.ErrInternalError, "GetRawRoleById error").WithCause(err))
			return
		}
	}
	role, ok := roleAny.(domain.BTMRole)
	if !ok || (ok && role.ID == 0) {
		log.Error("role not found")
		api.ErrResponse(c, "role not found", errors.NotFound(error_code.ErrDBError, "role not found"))
		return
	}
	if role.RoleName == "admin" {
		log.Error("role not allowed")
		api.ErrResponse(c, "role not allowed", errors.Forbidden(error_code.ErrForbidden, "不能新增 admin 權限"))
		return
	}

	// create user data
	hash, err := tools.GeneratePasswordHash(req.Password)
	if err != nil {
		log.Error("GeneratePasswordHash error", zap.Any("err", err))
		api.ErrResponse(c, "GeneratePasswordHash error", errors.InternalServer(error_code.ErrInternalError, "GeneratePasswordHash error").WithCause(err))
		return
	}
	user := domain.BTMUser{
		Account:  req.Account,
		Password: hash,
		Roles:    req.Role,
	}
	err = repo.CreateBTMUser(tx, user)
	if err != nil {
		log.Error("CreateUser error", zap.Any("err", err))
		api.ErrResponse(c, "CreateUser error", errors.InternalServer(error_code.ErrInternalError, "CreateUser error").WithCause(err))
		return
	}

	// create change log
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("FetchTokenInfo error", zap.Any("err", err))
		api.ErrResponse(c, "FetchTokenInfo error", errors.InternalServer(error_code.ErrInternalError, "FetchTokenInfo error").WithCause(err))
		return
	}

	// 資安防呆
	user.Password = ""
	createJsonData, err := json.Marshal(user)
	if err != nil {
		log.Error("json.Marshal(user)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(user)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(user)").WithCause(err))
		return
	}

	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMUsers,
		OperationType:   domain.BTMChangeLogOperationTypeCreate,
		CustomerId:      nil,
		BeforeValue:     nil,
		AfterValue:      createJsonData,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	c.JSON(200, api.DefaultRep{
		Code: 20000,
		Data: nil,
	})
	c.Done()

}
