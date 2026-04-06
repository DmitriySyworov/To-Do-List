package users

type ServiceUsers struct {
	*RepositoryUsers
}

func NewServiceUsers(repo *RepositoryUsers) *ServiceUsers {
	return &ServiceUsers{
		RepositoryUsers: repo,
	}
}
