package constants

func LeaderboardKey(sessionID string) string {
	return "session:" + sessionID + ":leaderboard"
}
