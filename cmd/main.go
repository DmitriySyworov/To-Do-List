package main

import (
	"context"
	"net/http"
	"time"
	"to-do-list/app/configs"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/users"
	"to-do-list/app/pkg/openDb"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	//context
	parentCtx := context.Background()
	redisCtx, cancel := context.WithTimeout(parentCtx, time.Minute*5)
	defer cancel()
	//openDb
	postgresDb := openDb.NewOpenPostgres(conf.DbConf)
	redisDb := openDb.NewOpenRedis(conf.DbConf)
	//Repository
	repoAuth := auth.NewRepositoryAuth(redisDb, redisCtx)
	usersRepo := users.NewRepositoryUsers(postgresDb)
	//Service
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IUserRepo: usersRepo, Configs: conf})
	serviceUsers := users.NewServiceUsers(usersRepo)
	//Handler
	router := http.NewServeMux()
	auth.NewHandlerAuth(router, &auth.HandlerAuthDep{ServiceAuth: serviceAuth})
	users.NewHandlerUser(router, &users.HandlerUsersDep{ServiceUsers: serviceUsers})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
