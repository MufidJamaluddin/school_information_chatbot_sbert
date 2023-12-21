package shared

type UserData struct {
	UserName     string `json:"username"`
	RoleGroupId  uint64 `json:"roleGroupId"`
	FullName     string `json:"fullName"`
	Position     string `json:"position"`
	PasswordHash string `json:"-"`
}
