package user

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
	"go.uber.org/zap"
)

type UpdateOneReq struct {
	Id      uint   `json:"id" binding:"required"`
	Account string `json:"account" binding:"required"`
	Role    int64  `json:"role" binding:"required"`
}

func UpdateOne(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateOne")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req UpdateOneReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	if req.Account == "admin" {
		log.Error("user not allowed")
		api.ErrResponse(c, "user not allowed", errors.BadRequest(error_code.ErrInvalidRequest, "不能修改 admin Account"))
		return
	}

	repo, err := di.NewRepo(configs.C.Mock)
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

	role, err := repo.GetRawRoleById(tx, req.Role)
	if err != nil {
		log.Error("GetRawRoleById error", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoleById error", errors.InternalServer(error_code.ErrInternalError, "GetRawRoleById error").WithCause(err))
		return
	}
	if role.ID == 0 {
		log.Error("role not found")
		api.ErrResponse(c, "role not found", errors.NotFound(error_code.ErrDBError, "role not found"))
		return
	}
	if role.RoleName == "admin" {
		log.Error("role not allowed")
		api.ErrResponse(c, "role not allowed", errors.Forbidden(error_code.ErrForbidden, "不能修改成 admin 權限"))
		return
	}

	// check user exist
	beforeAccount, err := repo.GetBTMUserById(tx, uint(req.Id))
	if err != nil {
		log.Error("GetBTMUserById error", zap.Any("err", err))
		api.ErrResponse(c, "GetBTMUserById error", errors.InternalServer(error_code.ErrInternalError, "GetBTMUserById error").WithCause(err))
		return
	}
	if beforeAccount.Account == "admin" {
		log.Error("user not allowed")
		api.ErrResponse(c, "user not allowed", errors.BadRequest(error_code.ErrInvalidRequest, "不能修改 admin 權限"))
		return
	}
	beforeAccount.Password = ""
	beforeJsonData, err := json.Marshal(beforeAccount)
	if err != nil {
		log.Error("json.Marshal(user)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(user)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(user)").WithCause(err))
		return
	}

	// update user
	user := domain.BTMUser{
		Account: req.Account,
		Roles:   req.Role,
	}
	err = repo.UpdateUserNameRoles(tx, uint(req.Id), req.Account, uint(req.Role))
	if err != nil {
		log.Error("UpdateOne error", zap.Any("err", err))
		api.ErrResponse(c, "UpdateOne error", errors.InternalServer(error_code.ErrInternalError, "UpdateOne error").WithCause(err))
		return
	}

	// create change log
	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("FetchTokenInfo error", zap.Any("err", err))
		api.ErrResponse(c, "FetchTokenInfo error", errors.InternalServer(error_code.ErrInternalError, "FetchTokenInfo error").WithCause(err))
		return
	}

	user.Password = ""
	afterJsonData, err := json.Marshal(user)
	if err != nil {
		log.Error("json.Marshal(user)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(user)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(user)").WithCause(err))
		return
	}

	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMUsers,
		OperationType:   domain.BTMChangeLogOperationTypeUpdate,
		CustomerId:      nil,
		BeforeValue:     beforeJsonData,
		AfterValue:      afterJsonData,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog err", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog", errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	// 更新權限登出用戶
	err = repo.DeleteLastLoginToken(tx, uint(req.Id))
	if err != nil {
		log.Error("DeleteLastLoginToken error", zap.Any("err", err))
		api.ErrResponse(c, "DeleteLastLoginToken error", errors.InternalServer(error_code.ErrInternalError, "DeleteLastLoginToken error").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
