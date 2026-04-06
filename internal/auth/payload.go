package auth

type RequestRegister struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type RequestLoginAndRestore struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type RequestConfirm struct {
	TempCode uint `json:"temp_code" validate:"required"`
}
type ResponseAuth struct {
	Message   string `json:"message"`
	SessionId string `json:"session_id"`
	JWT       string `json:"jwt"`
	Error     string `json:"error"`
}
type ResponseConfirm struct {
	JWT   string `json:"jwt"`
	Error string `json:"error"`
}
