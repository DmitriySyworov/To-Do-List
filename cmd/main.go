package main

import (
	"context"
	"net/http"
	"time"
	"to-do-list/app/configs"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/task"
	"to-do-list/app/internal/user"
	"to-do-list/app/pkg/open_Db"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	//context
	parentCtx := context.Background()
	redisCtx, cancel := context.WithTimeout(parentCtx, time.Second*5)
	defer cancel()
	//open_Db
	postgresDb := open_Db.NewOpenPostgres(conf.DbConf)
	redisDb := open_Db.NewOpenRedis(conf.DbConf)
	//Repository
	repoAuth := auth.NewRepositoryAuth(redisDb, redisCtx)
	repoUser := user.NewRepositoryUsers(postgresDb)
	repoTask := task.NewRepositoryTask(postgresDb)
	//Service
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IUserRepo: repoUser, Configs: conf})
	serviceUsers := user.NewServiceUsers(repoUser)
	serviceTask := task.NewServiceTask(repoTask, &task.ServiceTaskDep{IUserRepo: repoUser})
	//Handler
	router := http.NewServeMux()
	auth.NewHandlerAuth(router, &auth.HandlerAuthDep{ServiceAuth: serviceAuth, Configs: conf})
	user.NewHandlerUser(router, &user.HandlerUserDep{ServiceUser: serviceUsers, Configs: conf})
	task.NewHandlerTask(router, &task.HandlerTaskDep{ServiceTask: serviceTask, Configs: conf})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
