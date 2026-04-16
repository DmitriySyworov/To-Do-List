package stat

type ResponseMyStat struct {
	ActiveTask string `json:"active_task"`
	DoneTask   string `json:"done_task"`
	DeleteTask string `json:"delete_task"`
	Error      string `json:"error"`
}
type ResponseLeaderboard struct {
	User  []UserStat `json:"users"`
	Error string     `json:"error"`
}
type UserStat struct {
	Name     string `json:"name"`
	DoneTask string `json:"done_task"`
	Place    uint   `json:"place"`
}
