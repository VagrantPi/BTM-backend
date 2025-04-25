package domain

import (
	"sync/atomic"
	"time"
)

type BTMRole struct {
	ID       uint   `json:"id"`
	RoleName string `json:"role_name"`
	RoleDesc string `json:"role_desc"`
	RoleRaw  string `json:"role_raw"`
}

const RoleAdminName string = "admin"

const DefaultRoleRaw = `[{"id":"viewBar","path":"/view","name":"圖表","meta":{"roles":["admin"],"title":"圖表","icon":"component"},"children":[{"id":"viewDashboard","path":"dashboard","name":"Dashboard","meta":{"title":"Dashboard","icon":"dashboard","roles":["admin"]}}]},{"id":"permissionBar","path":"/permission","name":"後台管理","meta":{"roles":["admin"],"title":"後台管理","icon":"lock"},"children":[{"id":"adminUsers","path":"user","name":"後台用戶管理","meta":{"title":"後台用戶管理","icon":"user","roles":["admin"]}},{"id":"permission","path":"index","name":"權限管理","meta":{"title":"權限管理","icon":"lock","roles":["admin"]}}]},{"id":"userBar","path":"/user","redirect":"/user/page","name":"會員總覽","meta":{"title":"會員總覽","icon":"user","roles":["admin"]},"children":[{"id":"users","path":"users","name":"會員總覽","meta":{"title":"會員總覽","icon":"user","roles":["admin"],"noCache":true}},{"id":"userInfo","path":"info","name":"會員總覽 > 基本資料","meta":{"title":"基本資料","icon":"user","roles":["admin"],"noCache":true}},{"id":"riskView","path":"view","name":"會員總覽 > 風控限額編輯","meta":{"title":"風控限額編輯","roles":["admin"]}},{"id":"transaction","path":"/transaction","name":"交易記錄","meta":{"title":"交易記錄","icon":"el-icon-s-order","roles":["admin"],"noCache":true}}]},{"id":"reviewBar","path":"/review","redirect":"/review/cibs","name":"審核作業","meta":{"title":"審核作業","icon":"list","roles":["admin"]},"children":[{"id":"cibs","path":"cibs","name":"告誡名單","meta":{"title":"告誡名單","icon":"el-icon-warning","roles":["admin"],"noCache":true}},{"id":"cibsUpload","path":"cibs/upload","name":"上傳告誡名單","meta":{"title":"上傳告誡名單","icon":"el-icon-upload","roles":["admin"],"noCache":true}},{"id":"addresslist","path":"addresslist","name":"綁定地址","meta":{"title":"綁定地址","icon":"education","roles":["admin"],"noCache":true}},{"id":"addresslistView","path":"/addresslist/view","name":"綁定地址編輯","meta":{"title":"綁定地址編輯","roles":["admin"],"noCache":true}}]},{"id":"riskBar","path":"/risk_control","redirect":"/risk_control/page","name":"風險管理","meta":{"title":"風險管理","icon":"component","roles":["admin"]},"children":[{"id":"riskMemberList","path":"risk_member_list","name":"會員風控管理","meta":{"title":"會員風控管理","icon":"el-icon-warning","roles":["admin"]}},{"id":"riskControlHistory","path":"history","name":"會員風控管理 > 修改紀錄","meta":{"title":"風控修改紀錄","roles":["admin"]}},{"id":"eddList","path":"edd_list","name":"EDD名單","meta":{"title":"EDD名單","icon":"el-icon-warning-outline","roles":["admin"]}}]},{"id":"settingBar","path":"/setting","redirect":"/setting/page","name":"系統參數設定","meta":{"title":"系統參數設定","icon":"el-icon-document","roles":["admin"]},"children":[{"id":"systemSetting","path":"page","name":"設定內容","meta":{"title":"設定內容","icon":"el-icon-edit-outline","roles":["admin"]}},{"id":"systemSettingHistory","path":"history","name":"設定紀錄","meta":{"title":"設定紀錄","icon":"el-icon-document-copy","roles":["admin"]}}]}]`

type RoleMeta struct {
	Roles []string `json:"roles"`
}

type RoleItem struct {
	ID       string     `json:"id"`
	Path     string     `json:"path"`
	Children []RoleItem `json:"children,omitempty"`
	Meta     RoleMeta   `json:"meta"`
}

// TODO: 未來改使用 redis
type RoleWithTTL struct {
	CacheRole  BTMRole
	Expiration int64
}

// TODO: 未來改使用 redis
type UserReviewHistory struct {
	CacheHistory SumsubHistoryReviewData
	Expiration   int64
}

// TODO: 未來改使用 redis
type DeviceList struct {
	DeviceList map[string]Device
	Expiration int64
}

type TTLMap[T any] struct {
	Cache  T
	Expire int64
}

var (
	TTLRoleMap        atomic.Value
	TTLUserHistoryMap atomic.Value
	TTLDeviceListMap  atomic.Value
)

func init() {
	TTLRoleMap.Store(make(map[string]TTLMap[RoleWithTTL]))
	TTLUserHistoryMap.Store(make(map[string]TTLMap[UserReviewHistory]))
	TTLDeviceListMap.Store(make(map[string]TTLMap[DeviceList]))
}

func GetTTLMap[T any](mapName *atomic.Value, key string) (*T, bool) {
	m := mapName.Load().(map[string]TTLMap[T])
	item, exists := m[key]
	if !exists {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expire {
		CleanTTLMap[T](mapName, key)
		return nil, false
	}

	return &item.Cache, true
}

func SetTTLMap[T any](mapName *atomic.Value, key string, item T, expiration int64) {
	m := mapName.Load().(map[string]TTLMap[T])
	newMap := make(map[string]TTLMap[T])
	for k, v := range m {
		newMap[k] = v
	}
	newMap[key] = TTLMap[T]{
		Cache:  item,
		Expire: expiration,
	}
	mapName.Store(newMap)
}

func CleanTTLMap[T any](mapName *atomic.Value, key string) error {
	m := mapName.Load().(map[string]TTLMap[T])
	newMap := make(map[string]TTLMap[T])
	for k, v := range m {
		if k != key {
			newMap[k] = v
		}
	}
	mapName.Store(newMap)
	return nil
}

func PageIds(routes []RoleItem) []string {
	var ids []string
	var dfs func([]RoleItem)

	dfs = func(routes []RoleItem) {
		for _, route := range routes {
			ids = append(ids, route.ID)
			if len(route.Children) > 0 {
				dfs(route.Children)
			}
		}
	}

	dfs(routes)
	return ids
}
