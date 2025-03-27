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

const DefaultRoleRaw = `[{"id":"permissionBar","path":"/permission","meta":{"roles":["admin"],"title":"後台管理","icon":"lock"},"children":[{"id":"adminUsers","path":"user","name":"後台用戶管理","meta":{"title":"後台用戶管理","icon":"user","roles":["admin"]}},{"id":"permission","path":"index","name":"權限管理","meta":{"title":"權限管理","icon":"lock","roles":["admin"]}}]},{"id":"userBar","path":"/user","redirect":"/user/page","name":"User","meta":{"title":"會員總覽","icon":"user","roles":["admin"]},"children":[{"id":"users","path":"users","name":"會員總覽","meta":{"title":"會員總覽","icon":"user","roles":["admin"],"noCache":true}},{"id":"userInfo","path":"info","name":"基本資料","meta":{"title":"基本資料","icon":"user","roles":["admin"],"noCache":true}},{"id":"transaction","path":"/transaction","name":"交易記錄","meta":{"title":"交易記錄","icon":"el-icon-s-order","roles":["admin"],"noCache":true}}]},{"id":"riskBar","path":"/risk_control","redirect":"/risk_control/page","name":"風控","meta":{"title":"風控","icon":"component","roles":["admin"]},"children":[{"id":"whitelist","path":"whitelist","name":"風控白名單","meta":{"title":"風控白名單","icon":"eye-open","roles":["admin"]}},{"id":"whitelistView","path":"whitelist/view","name":"風控白名單限額編輯","meta":{"title":"風控白名單限額編輯","roles":["admin"]}},{"id":"graylist","path":"graylist","name":"風控灰名單","meta":{"title":"風控灰名單","icon":"eye","roles":["admin"]}},{"id":"graylistView","path":"graylist/view","name":"風控灰名單限額編輯","meta":{"title":"風控灰名單限額編輯","roles":["admin"]}},{"id":"blacklist","path":"blacklist","name":"風控黑名單","meta":{"title":"風控黑名單","icon":"el-icon-s-release","roles":["admin"]}},{"id":"blacklistView","path":"blacklist/view","name":"風控黑名單限額編輯","meta":{"title":"風控黑名單限額編輯","roles":["admin"]}}]},{"id":"reviewBar","path":"/review","redirect":"/review/cibs","name":"審核作業","meta":{"title":"審核作業","icon":"list","roles":["admin"]},"children":[{"id":"cibs","path":"cibs","name":"告誡名單","meta":{"title":"告誡名單","icon":"el-icon-warning","roles":["admin"],"noCache":true}},{"id":"cibsUpload","path":"cibs/upload","name":"上傳告誡名單","meta":{"title":"上傳告誡名單","icon":"el-icon-upload","roles":["admin"],"noCache":true}},{"id":"addresslist","path":"addresslist","name":"綁定地址","meta":{"title":"綁定地址","icon":"education","roles":["admin"],"noCache":true}},{"id":"addresslistView","path":"/addresslist/view","name":"綁定地址編輯","meta":{"title":"綁定地址編輯","roles":["admin"],"noCache":true}}]}]`

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

var (
	TTLRoleMap atomic.Value
)

func init() {
	TTLRoleMap.Store(make(map[uint]RoleWithTTL))
}

func GetTTLRoleMap(key uint) (*BTMRole, bool) {
	roleMap := TTLRoleMap.Load().(map[uint]RoleWithTTL)
	item, exists := roleMap[key]
	if !exists {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		_ = CleanTTLRoleMap(key)
		return nil, false
	}

	return &item.CacheRole, true
}

func SetTTLRoleMap(key uint, role BTMRole, expiration int64) {
	roleMap := TTLRoleMap.Load().(map[uint]RoleWithTTL)
	newMap := make(map[uint]RoleWithTTL)
	for k, v := range roleMap {
		newMap[k] = v
	}
	newMap[key] = RoleWithTTL{
		CacheRole:  role,
		Expiration: expiration,
	}
	TTLRoleMap.Store(newMap)
}

func CleanTTLRoleMap(key uint) error {
	roleMap := TTLRoleMap.Load().(map[uint]RoleWithTTL)
	newMap := make(map[uint]RoleWithTTL)
	for k, v := range roleMap {
		if k != key {
			newMap[k] = v
		}
	}
	TTLRoleMap.Store(newMap)
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
