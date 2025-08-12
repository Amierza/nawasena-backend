package entity

import "github.com/Amierza/nawasena-backend/constants"

type (
	Role string
)

const (
	SuperAdminRole Role = constants.ENUM_ROLE_SUPER_ADMIN
	AdminRole      Role = constants.ENUM_ROLE_ADMIN
)

func IsValidRoleAdmin(r Role) bool {
	return r == SuperAdminRole || r == AdminRole
}
