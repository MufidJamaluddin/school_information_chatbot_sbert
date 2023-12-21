package dto

type Admin struct {
	Username       string `json:"username" validate:"nonzero,max=10"`
	Password       string `json:"password,omitempty" validate:"min=3,max=60"`
	HashedPassword []byte `json:"-"`
	FullName       string `json:"fullName" validate:"nonzero,max=100"`
	PhoneNo        string `json:"phoneNo" validate:"nonzero,max=15"`
	Position       string `json:"position" validate:"nonzero,max=50"`
	RoleGroupId    uint64 `json:"roleGroupId"`
}

type ListAdminItem struct {
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	PhoneNo     string `json:"phoneNo"`
	Position    string `json:"position"`
	RoleGroupId uint64 `json:"roleGroupId"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}
