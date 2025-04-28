package settings

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/tests/settings/tools"
	dbTool "BTM-backend/third_party/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	db                 *gorm.DB     // 你要先初始化這個
	router             http.Handler // 你要初始化你的 router
	repo               domain.Repository
	apiToken           string
	latestResponseBody []byte
	currentCustomerID1 = uuid.New()
	currentCustomerID2 = uuid.New()
)

func init() {
	var err error
	db, _ = dbTool.ConnectToMockDatabase()
	tools.MigrateDB(db)
	apiToken = tools.InitApiToken(db)
	repo, err = di.NewRepo(true)

	if err != nil {
		panic(err)
	}
}

// 建立白名單用戶
func InitCustomer(arg1, arg2 int) error {
	err := repo.CreateCustomerLimit(db, currentCustomerID1)
	if err != nil {
		panic(err)
	}

	err = repo.CreateCustomerLimit(db, currentCustomerID2)
	if err != nil {
		panic(err)
	}

	return nil
}

// 更新限額
func updateDefaultLimit(table *godog.Table) error {
	params := map[string]interface{}{}
	for _, row := range table.Rows {
		key := row.Cells[0].Value
		val := row.Cells[1].Value
		if i, err := strconv.Atoi(val); err == nil {
			params[key] = i
		} else {
			params[key] = val
		}
	}
	body, _ := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPatch, "/api/config/limit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", apiToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	latestResponseBody = rec.Body.Bytes()
	return nil
}

// 驗證回傳 JSON
func apiResonseCheck(arg1 *godog.DocString) error {
	// 預期結果
	var expected map[string]interface{}
	if err := json.Unmarshal([]byte(arg1.Content), &expected); err != nil {
		return fmt.Errorf("解析預期結果失敗: %w", err)
	}

	// 真實 API 回傳
	var actual map[string]interface{}
	if err := json.Unmarshal(latestResponseBody, &actual); err != nil {
		return fmt.Errorf("解析實際回傳失敗: %w", err)
	}

	if actual["timestamp"] != nil {
		delete(actual, "timestamp")
	}

	// 比對兩個 map 是否相等
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("回傳不符\n預期: %v\n實際: %v", expected, actual)
	}
	return nil
}

// 取得限額
func getCustomerLimit(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11 int) error {
	url := fmt.Sprintf("/api/risk_control/%s/role", currentCustomerID1)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	req.Header.Set("token", apiToken)
	router.ServeHTTP(rec, req)
	latestResponseBody = rec.Body.Bytes()
	return nil
}

// 驗證回傳的限額設定為預設
func apiCustomerLimitCheckIsDefault(table *godog.Table) error {
	var actual map[string]interface{}
	if err := json.Unmarshal(latestResponseBody, &actual); err != nil {
		return err
	}
	data := actual["data"].(map[string]interface{})

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		expected := row.Cells[1].Value
		if fmt.Sprintf("%v", data[key]) != expected {
			return fmt.Errorf("欄位 %s 預期為 %s，實際為 %v", key, expected, data[key])
		}
	}
	return nil
}

func apiCustomerNLimitCheck(arg1 int, table *godog.Table) error {
	var actual map[string]interface{}
	if err := json.Unmarshal(latestResponseBody, &actual); err != nil {
		return err
	}
	data := actual["data"].(map[string]interface{})

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		expected := row.Cells[1].Value
		if fmt.Sprintf("%v", data[key]) != expected {
			return fmt.Errorf("欄位 %s 預期為 %s，實際為 %v", key, expected, data[key])
		}
	}
	return nil
}

func checkCustomerDBLimit(arg1 int, table *godog.Table) error {
	customer := currentCustomerID1
	if arg1 == 2 {
		customer = currentCustomerID2
	}
	var result model.BTMRiskControlCustomerLimitSetting
	err := db.Model(&model.BTMRiskControlCustomerLimitSetting{}).Where("customer_id = ?", customer).First(&result).Error
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"daily_limit":   result.DailyLimit.IntPart(),
		"monthly_limit": result.MonthlyLimit.IntPart(),
		"is_customized": result.IsCustomized,
	}

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		expected := row.Cells[1].Value
		if fmt.Sprintf("%v", data[key]) != expected {
			return fmt.Errorf("欄位 %s 預期為 %s，實際為 %v", key, expected, data[key])
		}
	}
	return nil
}

func currentCustomerID1CheckInit(daily, monthly int) error {
	var result model.BTMRiskControlCustomerLimitSetting
	err := db.Model(&model.BTMRiskControlCustomerLimitSetting{}).Where("customer_id = ?", currentCustomerID1).First(&result).Error
	if err != nil {
		return err
	}
	if result.DailyLimit.IntPart() != int64(daily) || result.MonthlyLimit.IntPart() != int64(monthly) {
		return fmt.Errorf("限額不符合，預期日 %d 月 %d，但實際為 %d / %d", daily, monthly, result.DailyLimit.IntPart(), result.MonthlyLimit.IntPart())
	}
	return nil
}

// 更新用戶限額
func updateCustomer2Limit(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11 int, table *godog.Table) error {
	params := map[string]interface{}{}
	for _, row := range table.Rows {
		key := row.Cells[0].Value
		val := row.Cells[1].Value
		if i, err := strconv.Atoi(val); err == nil {
			params[key] = i
		} else {
			params[key] = val
		}
	}
	body, _ := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/risk_control/%s/limit", currentCustomerID2), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", apiToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	latestResponseBody = rec.Body.Bytes()
	return nil
}

func removeDB() error {
	if err := db.Exec(`TRUNCATE "btm_login_logs" RESTART IDENTITY CASCADE;`).Error; err != nil {
		panic(err)
	}
	if err := db.Exec(`TRUNCATE "btm_risk_control_customer_limit_settings" RESTART IDENTITY CASCADE;`).Error; err != nil {
		panic(err)
	}
	if err := db.Exec(`TRUNCATE "btm_change_logs" RESTART IDENTITY CASCADE;`).Error; err != nil {
		panic(err)
	}
	return nil
}

// 註冊步驟
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^初始化用戶(\d+)限額、初始化用戶(\d+)限額$`, InitCustomer)
	ctx.Step(`^該用戶目前的日限額為 (\d+)，月限額為 (\d+)，且為非自訂限制$`, currentCustomerID1CheckInit)
	ctx.Step(`^回傳的結果應為預設限額且包含：$`, apiCustomerLimitCheckIsDefault)
	ctx.Step(`^用戶(\d+) 回傳的結果應包含：$`, apiCustomerNLimitCheck)
	ctx.Step(`^API 回傳的結果應為：$`, apiResonseCheck)
	ctx.Step(`^呼叫 GET \/api\/risk_control\/(\d+)e(\d+)f(\d+)c(\d+)-aa(\d+)b-(\d+)b(\d+)-(\d+)b-(\d+)ec(\d+)d(\d+)\/role 取得該用戶設定$`, getCustomerLimit)
	ctx.Step(`^管理者呼叫 PATCH \/api\/config\/limit API 並傳入以下參數：$`, updateDefaultLimit)
	ctx.Step(`^呼叫 PATCH \/api\/risk_control\/(\d+)e(\d+)f(\d+)c(\d+)-aa(\d+)b-(\d+)b(\d+)-(\d+)b-(\d+)ec(\d+)d(\d+)\/limit 修改用戶限額，並傳入以下參數：$`, updateCustomer2Limit)
	ctx.Step(`^用戶(\d+) DB 現在限額為：$`, checkCustomerDBLimit)
	ctx.Step(`^移除暫時DB$`, removeDB)
}

func TestMain(m *testing.M) {
	go func() {
		router = tools.SetupGin()
	}()

	status := godog.TestSuite{
		Name:                "settings",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",      // 顯示美化輸出
			Paths:  []string{"."}, // 當前資料夾找 .feature
			Strict: true,
		},
	}.Run()

	os.Exit(status)
}
