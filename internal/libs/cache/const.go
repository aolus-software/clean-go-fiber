package cache

import "time"

const (
	UserInformationKey         = "user:%s"
	UserInformationTTL         = time.Hour
	RefreshTokenKey            = "refresh:%s"
	PermissionSelectOptionsKey = "select_options:permissions"
	RoleSelectOptionsKey       = "select_options:roles"
	SelectOptionsTTL           = time.Hour
)
