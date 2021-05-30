package redisstore

import (
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store_redis"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client               *redis.Client
	inspectionRepository *InspectionRepository
}

func New(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) Inspection() store_redis.InspectionRepository {
	if r.inspectionRepository != nil {
		return r.inspectionRepository
	}

	r.inspectionRepository = &InspectionRepository{

		redis: r,
	}

	return r.inspectionRepository
}
