package stat

import (
	"log"
	"strconv"
	"to-do-list/app/pkg/di"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/event_bus"
)

type ServiceStat struct {
	*RepositoryStat
	*ServiceStatDep
}
type ServiceStatDep struct {
	di.IUserRepo
	*event_bus.EventBus
}

func NewServiceStat(repo *RepositoryStat, dep *ServiceStatDep) *ServiceStat {
	return &ServiceStat{
		RepositoryStat: repo,
		ServiceStatDep: dep,
	}
}
func (s *ServiceStat) GetLeaderBoard(userId uint, limitStr string) (*ResponseLeaderboard, error) {
	limit, errLimit := strconv.Atoi(limitStr)
	if errLimit != nil {
		return nil, ErrLimit
	}
	_, errGetUser := s.IUserRepo.GetUserById(userId)
	if errGetUser != nil {
		return nil, errors_custom.ErrNoExistUser
	}
	leaderboard, errGetLeaderboard := s.RepositoryStat.GetLeaderboard()
	if errGetLeaderboard != nil || len(leaderboard) == 0 {
		return nil, ErrLeaderboard
	}
	for i := range leaderboard {
		for j := 0; j < len(leaderboard)-i-1; j++ {
			if leaderboard[j].QuantityDoneTask < leaderboard[j+1].QuantityDoneTask {
				temp := leaderboard[j].QuantityDoneTask
				leaderboard[j].QuantityDoneTask = leaderboard[j+1].QuantityDoneTask
				leaderboard[j].QuantityDoneTask = temp
			}
		}
	}
	var place uint = 0
	var resLeaderboard []UserStat
	for i := 0; i < len(leaderboard); i++ {
		place++
		leaderboard[i].Place = place
		resLeaderboard = append(resLeaderboard, leaderboard[i])
		if place == uint(limit) {
			return &ResponseLeaderboard{
				User: resLeaderboard,
			}, nil
		}
	}
	return nil, ErrLeaderboard
}
func (s *ServiceStat) AddTaskInStat() {
	for {
		for event := range s.EventBus.Subscribe() {
			userId, ok := event.Data.(uint)
			user, errGetUser := s.IUserRepo.GetUserById(userId)
			if errGetUser != nil {
				log.Println(errGetUser)
			}
			if !ok {
				log.Println("eventbus data incorrect")
			}
			if event.Name == event_bus.EventCreateTask {
				errAddCreate := s.RepositoryStat.AddCreateTask(userId, user.Name)
				if errAddCreate != nil {
					log.Println(errAddCreate)
				}
			} else if event.Name == event_bus.EventDoneTask {
				errAddDone := s.AddDoneTask(userId)
				if errAddDone != nil {
					log.Println(errAddDone)
				}
			} else if event.Name == event_bus.EventDeleteActiveTask {
				errAddDelete := s.AddDeleteTask(userId, fieldCreate)
				if errAddDelete != nil {
					log.Println(errAddDelete)
				}
			} else if event.Name == event_bus.EventDeleteDoneTask {
				errAddDelete := s.AddDeleteTask(userId, fieldDone)
				if errAddDelete != nil {
					log.Println(errAddDelete)
				}
			}
		}
	}
}
