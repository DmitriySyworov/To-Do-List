package stat

import "errors"

var (
	ErrLeaderboardEmpty = errors.New("leaderboard is empty")
	ErrNotFoundStat     = errors.New("your statistics data was not found")
	ErrLeaderboard      = errors.New("failed to load leaderboard")
	ErrLimit            = errors.New("incorrect limit")
)
