package main

import (
	"context"
	"net/http"
	"time"
	"to-do-list/app/configs"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/stat"
	"to-do-list/app/internal/task"
	"to-do-list/app/internal/user"
	"to-do-list/app/pkg/event_bus"
	"to-do-list/app/pkg/open_Db"
)

func main() {
	//conf
	conf := configs.NewConfigs()
	//context
	parentCtx := context.Background()
	redisCtx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	//EventBus
	eventBus := event_bus.NewEventBus()
	//open_Db
	postgresDb := open_Db.NewOpenPostgres(conf.DbConf)
	redisDb := open_Db.NewOpenRedis(conf.DbConf)
	//Repository
	repoAuth := auth.NewRepositoryAuth(redisDb, redisCtx)
	repoUser := user.NewRepositoryUsers(postgresDb)
	repoTask := task.NewRepositoryTask(postgresDb)
	repoStat := stat.NewRepositoryStat(redisDb, redisCtx)
	//Service
	serviceAuth := auth.NewServiceAuth(repoAuth, &auth.ServiceAuthDep{IUserRepo: repoUser, Configs: conf})
	serviceUsers := user.NewServiceUsers(repoUser)
	serviceTask := task.NewServiceTask(repoTask, &task.ServiceTaskDep{IUserRepo: repoUser, EventBus: eventBus})
	serviceStat := stat.NewServiceStat(repoStat, &stat.ServiceStatDep{EventBus: eventBus, IUserRepo: repoUser})
	//Stat
	go serviceStat.AddTaskInStat()
	//Handler
	router := http.NewServeMux()
	auth.NewHandlerAuth(router, &auth.HandlerAuthDep{ServiceAuth: serviceAuth, Configs: conf})
	user.NewHandlerUser(router, &user.HandlerUserDep{ServiceUser: serviceUsers, Configs: conf})
	task.NewHandlerTask(router, &task.HandlerTaskDep{ServiceTask: serviceTask, Configs: conf})
	stat.NewHandlerStat(router, &stat.HandlerStatDep{ServiceStat: serviceStat, Configs: conf})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
