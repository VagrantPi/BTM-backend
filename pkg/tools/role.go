package tools

type Role int64

const (
	RoleAdmin = 1 << iota
	RoleEditor
)

var AllRoles = []Role{
	RoleAdmin,
	RoleEditor,
}

func (t *Role) has(input Role) bool {
	return (*t)&input != 0
}

var StringToRole = map[string]int64{
	"admin":  int64(RoleAdmin),
	"editor": int64(RoleEditor),
}

var RoleToString = map[Role]string{
	RoleAdmin:  "admin",
	RoleEditor: "editor",
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
