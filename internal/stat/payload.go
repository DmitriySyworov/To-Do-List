package stat

type ResponseMyStat struct {
	QuantityActiveTask string `json:"quantity_active_task"`
	QuantityDoneTask   string `json:"quantity_done_task"`
	QuantityDeleteTask string `json:"quantity_delete_task"`
	Error              string `json:"error"`
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
