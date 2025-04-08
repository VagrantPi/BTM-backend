package user

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

func GetBTMUserRoleRoutes(c *gin.Context) {
	api.OKResponse(c, domain.DefaultRoleRaw)
}

type BTMRoleOnlyNameResp struct {
	ID       uint   `json:"id"`
	RoleName string `json:"role_name"`
}

func GetBTMUserRoles(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetBTMUserRoleRoutes")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	onlyResponseName := c.Query("only_name") == "true"

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	list, err := repo.GetRawRoles(repo.GetDb(c))
	if err != nil {
		log.Error("GetRawRoles", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoles", errors.InternalServer(error_code.ErrDiError, "GetRawRoles").WithCause(err))
		return
	}

	if onlyResponseName {
		var onlyNameList []BTMRoleOnlyNameResp
		for _, v := range list {
			onlyNameList = append(onlyNameList, BTMRoleOnlyNameResp{
				ID:       v.ID,
				RoleName: v.RoleName,
			})
		}
		api.OKResponse(c, onlyNameList)
		return
	}

	api.OKResponse(c, list)
}

type CreateRoleReq struct {
	Name     string `json:"role_name" binding:"required"`
	RoleDesc string `json:"role_desc" binding:"required"`
	RoleRaw  string `json:"role_raw" binding:"required"`
}

func CreateRole(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "CreateRole")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req CreateRoleReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
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
		log.Error("repo.NewTransactionBegin()", zap.Any("err", err))
		api.ErrResponse(c, "repo.NewTransactionBegin()", errors.InternalServer(error_code.ErrDiError, "repo.NewTransactionBegin()").WithCause(err))
		return
	}
	defer repo.TransactionCommit(tx)

	createData := domain.BTMRole{
		RoleName: req.Name,
		RoleDesc: req.RoleDesc,
		RoleRaw:  req.RoleRaw,
	}
	err = repo.CreateRole(tx, createData)
	if err != nil {
		log.Error("CreateRole error", zap.Any("err", err))
		api.ErrResponse(c, "CreateRole error", errors.InternalServer(error_code.ErrDiError, "CreateRole").WithCause(err))
		return
	}

	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	createJsonData, err := json.Marshal(createData)
	if err != nil {
		log.Error("json.Marshal(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(whitelist)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(whitelist)").WithCause(err))
		return
	}

	// add change log
	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMRoles,
		OperationType:   domain.BTMChangeLogOperationTypeCreate,
		CustomerId:      nil,
		BeforeValue:     nil,
		AfterValue:      createJsonData,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog error", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog error", errors.InternalServer(error_code.ErrDiError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}

type UpdateRoleReq struct {
	Name     string `json:"name"`
	RoleDesc string `json:"role_desc"`
	RoleRaw  string `json:"role_raw"`
}

func UpdateRole(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "UpdateRole")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	var req CreateRoleReq
	err := c.BindJSON(&req)
	if err != nil {
		log.Error("BindJSON error", zap.Any("err", err))
		api.ErrResponse(c, "BindJSON error", errors.BadRequest(error_code.ErrInvalidRequest, "BindJSON error").WithCause(err))
		return
	}

	repo, err := di.NewRepo()
	if err != nil {
		log.Error("di.NewRepo()", zap.Any("err", err))
		api.ErrResponse(c, "di.NewRepo()", errors.InternalServer(error_code.ErrDiError, "di.NewRepo()").WithCause(err))
		return
	}

	operationUserInfo, err := tools.FetchTokenInfo(c)
	if err != nil {
		log.Error("repo.CreateWhitelist(whitelist)", zap.Any("err", err))
		api.ErrResponse(c, "repo.CreateWhitelist(whitelist)", errors.InternalServer(error_code.ErrDBError, "repo.CreateWhitelist(whitelist)").WithCause(err))
		return
	}

	tx, err := repo.NewTransactionBegin(c)
	if err != nil {
		log.Error("repo.NewTransactionBegin()", zap.Any("err", err))
		api.ErrResponse(c, "repo.NewTransactionBegin()", errors.InternalServer(error_code.ErrDiError, "repo.NewTransactionBegin()").WithCause(err))
		return
	}
	defer repo.TransactionCommit(tx)

	// get origin role_raw
	beforeRoleValue, err := repo.GetRawRoleByRoleName(tx, req.Name)
	if err != nil {
		log.Error("GetRawRoleByRoleName", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoleByRoleName", errors.InternalServer(error_code.ErrDBError, "GetRawRoleByRoleName").WithCause(err))
		return
	}
	if beforeRoleValue.ID == 0 {
		log.Error("GetRawRoleByRoleName", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoleByRoleName", errors.BadRequest(error_code.ErrDBError, "Not found").WithCause(err))
		return
	}
	if beforeRoleValue.RoleName == "admin" || beforeRoleValue.RoleName == "no_role" {
		log.Error("GetRawRoleByRoleName", zap.Any("err", err))
		api.ErrResponse(c, "GetRawRoleByRoleName", errors.Forbidden(error_code.ErrForbidden, "不能修改 admin 和 no_role 權限"))
		return
	}

	beforeDataJson, err := json.Marshal(beforeRoleValue)
	if err != nil {
		log.Error("json.Marshal(beforeRoleValue)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(beforeRoleValue)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeRoleValue)").WithCause(err))
		return
	}

	// update role
	updateData := domain.BTMRole{
		RoleName: req.Name,
		RoleDesc: req.RoleDesc,
		RoleRaw:  req.RoleRaw,
	}
	err = repo.UpdateRole(tx, updateData)
	if err != nil {
		log.Error("UpdateRole error", zap.Any("err", err))
		api.ErrResponse(c, "UpdateRole error", errors.InternalServer(error_code.ErrDiError, "UpdateRole").WithCause(err))
		return
	}

	// delete cache
	domain.CleanTTLMap[domain.RoleWithTTL](&domain.TTLRoleMap, fmt.Sprintf("%d", beforeRoleValue.ID))

	afterDataJson, err := json.Marshal(updateData)
	if err != nil {
		log.Error("json.Marshal(updateData)", zap.Any("err", err))
		api.ErrResponse(c, "json.Marshal(updateData)", errors.InternalServer(error_code.ErrDBError, "json.Marshal(updateData)").WithCause(err))
		return
	}

	// add change log
	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserInfo.Id,
		TableName:       domain.BTMChangeLogTableNameBTMRoles,
		OperationType:   domain.BTMChangeLogOperationTypeUpdate,
		CustomerId:      nil,
		BeforeValue:     beforeDataJson,
		AfterValue:      afterDataJson,
	})
	if err != nil {
		log.Error("CreateBTMChangeLog error", zap.Any("err", err))
		api.ErrResponse(c, "CreateBTMChangeLog error", errors.InternalServer(error_code.ErrDiError, "CreateBTMChangeLog").WithCause(err))
		return
	}

	api.OKResponse(c, nil)
}
