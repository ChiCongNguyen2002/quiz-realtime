package redis

import (
	"context"

	"quiz-realtime/internal/constants"
	"quiz-realtime/internal/domain/leaderboard"

	"github.com/redis/go-redis/v9"
)

type LeaderboardRepository struct {
	Client *redis.Client
}

func NewLeaderboardRepository(client *redis.Client) *LeaderboardRepository {
	return &LeaderboardRepository{
		Client: client,
	}
}

func (r *LeaderboardRepository) UpdateScore(sessionID string, userID string, score int) error {
	key := constants.LeaderboardKey(sessionID)

	return r.Client.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(score),
		Member: userID,
	}).Err()
}

func (r *LeaderboardRepository) GetLeaderboard(sessionID string) ([]leaderboard.Entry, error) {
	key := constants.LeaderboardKey(sessionID)

	res, err := r.Client.ZRevRangeWithScores(context.Background(), key, 0, 10).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]leaderboard.Entry, 0, len(res))

	for _, item := range res {
		member, ok := item.Member.(string)
		if !ok {
			continue
		}

		entries = append(entries, leaderboard.Entry{
			UserID: member,
			Score:  int(item.Score),
		})
	}

	return entries, nil
}
