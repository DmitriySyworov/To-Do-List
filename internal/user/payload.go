package user

type RequestUpdateUser struct {
	OriginalPassword string `json:"original_password" validate:"required_without_all=Name"`
	Name             string `json:"name" validate:"required_without_all=Email NewPassword"`
	Email            string `json:"email" validate:"required_without_all=Name NewPassword,email"`
	NewPassword      string `json:"new_password" validate:"required_without_all=Name Email"`
}
type RequestDeleteUser struct {
	Password string `json:"password" validate:"required"`
}
