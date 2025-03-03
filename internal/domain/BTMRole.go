package domain

type BTMRole struct {
	ID       uint   `json:"id"`
	RoleName string `json:"role_name"`
	RoleDesc string `json:"role_desc"`
	Role     int64  `json:"role"`
	RoleRaw  string `json:"role_raw"`
}

const DefaultRoleRaw = `
[
   {
      "path":"/permission",
      "alwaysShow":true,
      "name":"Permission",
      "meta":{
         "title":"Permission",
         "roles":[
            "admin"
         ]
      },
      "children":[
         {
            "path":"role",
            "name":"RolePermission",
            "meta":{
               "title":"Role Permission",
               "roles":[
                  "admin"
               ]
            }
         }
      ]
   },
   {
      "path":"/whitelist",
      "redirect":"/whitelist",
      "alwaysShow":true,
      "name":"Whitelist",
      "meta":{
         "title":"Whitelist"
      },
      "children":[
         {
            "path":"index",
            "name":"Whitelist List",
            "meta":{
               "title":"Whitelist List",
               "roles":[
                  "admin"
               ]
            }
         },
         {
            "path":"view",
            "name":"Whitelist View",
            "meta":{
               "title":"Whitelist View"
            }
         }
      ]
   },
   {
      "path":"*",
      "redirect":"/404",
      "hidden":true
   }
]
`
