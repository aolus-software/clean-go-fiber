package selectoptions

type PermissionOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type PermissionGroupOption struct {
	Group       string             `json:"group"`
	Permissions []PermissionOption `json:"permissions"`
}

type RoleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
