package dao

import (
	"context"
	"time"

	"github.com/NTHU-LSALAB/NTHU-Distributed-System/pkg/rediskit"
	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
)

type redisCommentDAO struct {
	cache   *cache.Cache
	baseDAO CommentDAO
}

var _ CommentDAO = (*redisCommentDAO)(nil)

const (
	commentDAOLocalCacheSize     = 1024
	commentDAOLocalCacheDuration = 1 * time.Minute
	commentDAORedisCacheDuration = 3 * time.Minute
)

// create redisCommentDAO by cache and baseDAO
func NewRedisCommentDAO(client *rediskit.RedisClient, baseDAO CommentDAO) *redisCommentDAO {
	// Redis !TODO
	// also ref from comment_pg.go && video_redis.go
	return &redisCommentDAO{
		// ref : https://pkg.go.dev/github.com/go-redis/cache/v8@v8.4.3#section-readme
		cache: cache.New(&cache.Options{
			// ref : https://docs.google.com/presentation/d/1lg_vK-Man8WjaZrCDOR6759vTyJYZaYaH0VqVIbCODU/edit#slide=id.g116535da8c4_1_35
			// page 21
			// for the localcache not sure second parameter should be commentDAOLocalCacheDuration or commentDAORedisCacheDuration
			Redis:      client,
			LocalCache: cache.NewTinyLFU(commentDAOLocalCacheSize, commentDAOLocalCacheDuration),
		}),
		baseDAO: baseDAO,
	}
}

// List all the comments related to the videoID
// Notice that all the comments will be stored as a single value in the cache
// The key is generated by listCommentKey function in comment.go
// The implementation should handle both cache miss and cache hit scenarios
func (dao *redisCommentDAO) ListByVideoID(ctx context.Context, videoID string, limit, offset int) ([]*Comment, error) {
	// Redis !TODO
	// ref from comment_pg.go && video_redis.go
	var comments []*Comment
	// ref : https://docs.google.com/presentation/d/1lg_vK-Man8WjaZrCDOR6759vTyJYZaYaH0VqVIbCODU/edit#slide=id.g1193fc43cd2_1_23
	// page 23
	if err := dao.cache.Once(&cache.Item{
		Key:   listCommentKey(videoID, limit, offset),
		Value: &comments,
		TTL:   commentDAORedisCacheDuration,
		Do: func(*cache.Item) (interface{}, error) {
			return dao.baseDAO.ListByVideoID(ctx, videoID, limit, offset)
		},
	}); err != nil {
		return nil, err
	}

	return comments, nil
}

// The operation are not cacheable, just pass down to baseDAO
func (dao *redisCommentDAO) Create(ctx context.Context, comment *Comment) (uuid.UUID, error) {
	// Redis !TODO
	return dao.baseDAO.Create(ctx, comment)
}

// The operation are not cacheable, just pass down to baseDAO
func (dao *redisCommentDAO) Update(ctx context.Context, comment *Comment) error {
	// Redis !TODO
	return dao.baseDAO.Update(ctx, comment)
}

// The operation are not cacheable, just pass down to baseDAO
func (dao *redisCommentDAO) Delete(ctx context.Context, id uuid.UUID) error {
	// Redis !TODO
	return dao.baseDAO.Delete(ctx, id)
}

// The operation are not cacheable, just pass down to baseDAO
func (dao *redisCommentDAO) DeleteByVideoID(ctx context.Context, videoID string) error {
	// Redis !TODO
	return dao.baseDAO.DeleteByVideoID(ctx, videoID)
}
