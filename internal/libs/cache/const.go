package cache

import "time"

const (
	UserInformationKey = "user:%s"
	UserInformationTTL = time.Hour
	RefreshTokenKey    = "refresh:%s"
)
