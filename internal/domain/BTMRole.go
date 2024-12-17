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
      "component":"layout/Layout",
      "redirect":"/permission/page",
      "alwaysShow":true,
      "name":"Permission",
      "meta":{
         "title":"Permission",
         "icon":"lock",
         "roles":[
            "admin"
         ]
      },
      "children":[
         {
            "path":"role",
            "component":"views/permission/role",
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
      "component":"layout/Layout",
      "meta":{
         "roles":[
            "admin",
            "editor"
         ]
      },
      "children":[
         {
            "path":"index",
            "component":"views/whitelist/index",
            "name":"WhiteList",
            "meta":{
               "title":"WhiteList",
               "icon":"education",
               "roles":[
                  "admin",
                  "editor"
               ],
               "noCache":true
            }
         },
         {
            "path":"/whitelist/view",
            "component":"views/whitelist/components/view",
            "name":"WhiteListView",
            "hidden":true,
            "meta":{
               "title":"WhiteList View",
               "roles":[
                  "admin",
                  "editor"
               ],
               "noCache":true
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
