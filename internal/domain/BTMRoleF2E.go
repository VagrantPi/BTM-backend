package domain

// F2ERouterId 前端路由 id
type F2ERouterId string

const (
	// F2EPermission - 權限設定
	F2EPermission F2ERouterId = "permission"
	// F2ETx - 交易紀錄
	F2ETx F2ERouterId = "transaction"
	// F2EWhitelist - 風控白名單
	F2EWhitelist F2ERouterId = "whitelist"
	// F2EGraylist - 風控灰名單
	F2EGraylist F2ERouterId = "graylist"
	// F2EBlacklist - 風控黑名單
	F2EBlacklist F2ERouterId = "blacklist"
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
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
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
	case method == "GET" && uri == "/api/risk_control/roles":
		return []string{
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
		}
	case method == "GET" && uri == "/api/risk_control/:customer_id/role":
		return []string{
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
		}
	case method == "PATCH" && uri == "/api/risk_control/:customer_id/role":
		return []string{
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
		}
	case method == "PATCH" && uri == "/api/risk_control/:customer_id/limit":
		return []string{
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
		}
	case method == "GET" && uri == "/api/btm/logs":
		return []string{
			F2EWhitelist.String(),
			F2EGraylist.String(),
			F2EBlacklist.String(),
		}
	default:
		return []string{}
	}
}
