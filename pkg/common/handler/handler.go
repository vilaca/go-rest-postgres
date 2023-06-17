package handler

import (
	"github.com/pauljamescleary/gomin/pkg/common/config"
	"github.com/pauljamescleary/gomin/pkg/common/db"
	"github.com/pauljamescleary/gomin/pkg/common/redis"
)

type Handler struct {
	UserRepo db.UserRepository
	Queue    redis.QueueRepository
}

func NewHandler(ur db.UserRepository, queue redis.QueueRepository) *Handler {
	return &Handler{UserRepo: ur, Queue: queue}
}

func LoadHandler(configPath *string) *Handler {
	cfg, _ := config.LoadConfig(configPath)
	return LoadHandlerFromConfig(cfg)
}

func LoadHandlerFromConfig(cfg config.Config) *Handler {
	rc, _ := redis.StartRedisClient(cfg)
	rp, _ := redis.NewQueue(rc, cfg.NewClusterStream)
	rp.Read()
	database := db.NewDatabase(cfg)
	userRepo, _ := db.NewUserRepository(database)
	handler := NewHandler(userRepo, rp)
	return handler
}
