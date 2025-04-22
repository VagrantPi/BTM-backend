package domain

// F2ERouterId 前端路由 id
type F2ERouterId string

const (
	// F2EAdminUsers - 後台用戶管理
	F2EAdminUsers F2ERouterId = "adminUsers"
	// F2EPermission - 權限管理
	F2EPermission F2ERouterId = "permission"

	// F2EUser - 用戶總覽
	F2EUser F2ERouterId = "users"
	// F2EUserInfo - 用戶總覽 > 個人基本資料
	F2EUserInfo F2ERouterId = "userInfo"
	// F2ERiskView - 用戶總覽 > 風險管理
	F2ERiskView F2ERouterId = "riskView"
	// F2ETx - 交易紀錄
	F2ETx F2ERouterId = "transaction"

	// F2ERiskMemberList - 風控名單
	F2ERiskMemberList F2ERouterId = "riskMemberList"

	// F2ECIBs - 告誡名單
	F2ECIBs F2ERouterId = "cibs"
	// F2EAddressList - 綁定地址
	F2EAddressList F2ERouterId = "addresslist"
)

func (r F2ERouterId) String() string {
	return string(r)
}

func URIToF2ERouterId(method string, uri string) []string {
	switch {
	// 權限設定頁面相關
	case method == "GET" && uri == "/api/user/role/routes":
		return []string{F2EPermission.String()}
	case method == "POST" && uri == "/api/user/role":
		return []string{F2EPermission.String()}
	case method == "PUT" && uri == "/api/user/role":
		return []string{F2EPermission.String()}

	case method == "GET" && uri == "/api/customer/list":
		return []string{
			F2EAddressList.String(),
			F2EUser.String(),
			F2ERiskMemberList.String(),
		}
	// 綁定地址頁面相關
	case method == "GET" && uri == "/api/customer/whitelist":
		return []string{F2EAddressList.String()}
	case method == "GET" && uri == "/api/customer/whitelist/search":
		return []string{F2EAddressList.String()}
	case method == "POST" && uri == "/api/customer/whitelist":
		return []string{F2EAddressList.String()}
	case method == "DELETE" && uri == "/api/customer/whitelist":
		return []string{F2EAddressList.String()}

	// 交易紀錄頁面相關
	case method == "GET" && uri == "/api/tx/list":
		return []string{F2ETx.String()}

	// 告誡名單頁面相關
	case method == "GET" && uri == "/api/cib/list":
		return []string{F2ECIBs.String()}
	case method == "POST" && uri == "/api/cib/upload":
		return []string{F2ECIBs.String()}

	// 風控頁面相關
	case method == "GET" && uri == "/api/risk_control/:customer_id/role":
		return []string{
			F2ERiskView.String(),
		}
	case method == "PATCH" && uri == "/api/risk_control/:customer_id/role":
		return []string{
			F2ERiskView.String(),
		}
	case method == "PATCH" && uri == "/api/risk_control/:customer_id/limit":
		return []string{
			F2ERiskView.String(),
		}
	case method == "GET" && uri == "/api/btm/logs":
		return []string{
			F2ERiskView.String(),
			F2EUserInfo.String(),
		}

	// 用戶基本資料頁面
	case method == "GET" && uri == "/api/customer/image":
		return []string{F2EUserInfo.String()}
	default:
		return []string{}
	}
}
