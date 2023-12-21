package dto

type User struct {
	Id        uint   `json:"id"`
	FullName  string `json:"fullName" validate:"min=1,max=50"`
	UserRole  string `json:"userRole" validate:"max=20"`
	ClassName string `json:"className"`
	Age       uint8  `json:"age" validate:"min=1,max=250"`
}

func (u *User) GetIsStudent() bool {
	return u.UserRole == "student"
}