package main

import (
	"net/http"
	"to-do-list/app/configs"
	"to-do-list/app/internal/auth"
	"to-do-list/app/internal/stat"
	"to-do-list/app/internal/task"
	"to-do-list/app/internal/user"
	"to-do-list/app/pkg/event_bus"
	"to-do-list/app/pkg/middleware"
	"to-do-list/app/pkg/open_Db"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: App(),
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
func App() http.Handler {
	//conf
	conf := configs.NewConfigs()
	//EventBus
	eventBus := event_bus.NewEventBus()
	//open_Db
	postgresDb := open_Db.NewOpenPostgres(conf.DSN)
	redisDb := open_Db.NewOpenRedis(conf.RedisPassword)
	//Repository
	repoAuth := auth.NewRepositoryAuth(redisDb)
	repoUser := user.NewRepositoryUsers(postgresDb)
	repoTask := task.NewRepositoryTask(postgresDb)
	repoStat := stat.NewRepositoryStat(redisDb)
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
	//chain
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.RecoveryPanic,
	)
	return stack(router)
}
