package tools

type Role int64

const (
	RoleAdmin = 1 << iota
	RoleEditor
)

func (t *Role) has(input Role) bool {
	return (*t)&input != 0
}

var StringToRole = map[string]int64{
	"admin":  int64(RoleAdmin),
	"editor": int64(RoleEditor),
}

func (t *Role) ToStrings() []string {
	var roleList []string
	if t.has(RoleAdmin) {
		roleList = append(roleList, "admin")
	}
	if t.has(RoleEditor) {
		roleList = append(roleList, "editor")
	}
	return roleList
}
