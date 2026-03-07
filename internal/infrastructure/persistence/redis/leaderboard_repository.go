package redis

import (
	"context"

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
	key := "session:" + sessionID + ":leaderboard"

	return r.Client.ZAdd(
		context.Background(),
		key,
		redis.Z{
			Score:  float64(score),
			Member: userID,
		},
	).Err()
}

func (r *LeaderboardRepository) GetLeaderboard(sessionID string) ([]leaderboard.Entry, error) {
	key := "session:" + sessionID + ":leaderboard"

	res, err := r.Client.ZRevRangeWithScores(
		context.Background(),
		key,
		0,
		10,
	).Result()
	if err != nil {
		return nil, err
	}

	var list []leaderboard.Entry

	for _, item := range res {
		member, ok := item.Member.(string)
		if !ok {
			continue
		}

		list = append(list, leaderboard.Entry{
			UserID: member,
			Score:  int(item.Score),
		})
	}

	return list, nil
}
