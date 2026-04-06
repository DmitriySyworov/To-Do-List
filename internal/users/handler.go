package users

import "net/http"

type HandlerUsers struct {
	*HandlerUsersDep
}
type HandlerUsersDep struct {
	*ServiceUsers
}

func NewHandlerUser(router *http.ServeMux, dep *HandlerUsersDep) {
	user := &HandlerUsers{
		HandlerUsersDep: dep,
	}
	router.Handle("GET /users/my", user.GetUser())
	router.Handle("PATCH /users/my", user.UpdateUser())
	router.Handle("DELETE /users/my", user.DeleteUser())
}
func (hl *HandlerUsers) GetUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
func (hl *HandlerUsers) UpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
func (hl *HandlerUsers) DeleteUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
