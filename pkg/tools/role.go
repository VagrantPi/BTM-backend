package tools

type Role int64

const (
	RoleAdmin = 1 << iota
)

var AllRoles = []Role{
	RoleAdmin,
}

func (t *Role) has(input Role) bool {
	return (*t)&input != 0
}

var StringToRole = map[string]int64{
	"admin": int64(RoleAdmin),
}

var RoleToString = map[Role]string{
	RoleAdmin: "admin",
}

func (t *Role) ToStrings() []string {
	var roleList []string
	if t.has(RoleAdmin) {
		roleList = append(roleList, "admin")
	}
	return roleList
}
