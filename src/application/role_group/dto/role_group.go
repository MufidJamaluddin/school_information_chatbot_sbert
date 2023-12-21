package dto

import "time"

type CreateRoleGroupDTO struct {
	RoleGroup string `json:"roleGroup"`
	CreatedBy string `json:"-"`
}

type UpdateRoleGroupDTO struct {
	RoleGroup string `json:"roleGroup"`
	UpdatedBy string `json:"-"`
}

type RoleItemGroupDTO struct {
	RoleGroupId uint64     `json:"id"`
	RoleGroup   string     `json:"roleGroup"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   string     `json:"updatedBy"`
}

type ResponseMessageDTO struct {
	Id      uint64 `json:"id"`
	Message string `json:"message"`
}
