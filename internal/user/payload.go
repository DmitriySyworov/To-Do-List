package user

type RequestUpdateUser struct {
	OriginalPassword string `json:"original_password" validate:"required_without=Name"`
	Name             string `json:"name" validate:"required_without=OriginalPassword"`
	Email            string `json:"email" validate:"omitempty,email,required_with=OriginalPassword"`
	NewPassword      string `json:"new_password" validate:"required_with=OriginalPassword"`
}
type RequestDeleteUser struct {
	Password string `json:"password" validate:"required"`
}
