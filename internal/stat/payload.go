package stat

type ResponseMyStat struct {
	ActiveTask string `json:"active_task"`
	DoneTask   string `json:"done_task"`
	DeleteTask string `json:"delete_task"`
	Error      string `json:"error"`
}
type ResponseLeaderboard struct {
	User  []UserStat `json:"user"`
	Error string     `json:"error"`
}
type UserStat struct {
	Name             string `json:"name"`
	QuantityDoneTask string `json:"quantity_done_task"`
	Place            uint   `json:"place"`
}
