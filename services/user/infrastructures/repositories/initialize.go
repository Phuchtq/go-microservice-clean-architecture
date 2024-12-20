package repositories

import (
	external_services "architecture_template/services/user/external_services/role"
	"architecture_template/services/user/interfaces"
	external_iservices "architecture_template/services/user/interfaces/external_services"
	"database/sql"
	"log"

	"github.com/redis/go-redis/v9"
)

type repo struct {
	db                  *sql.DB
	logger              *log.Logger
	redisCache          *redis.Client
	externalRoleService external_iservices.IRole
}

func InitializeRepository(db *sql.DB, logger *log.Logger, redisClient *redis.Client) interfaces.IRepository {
	return &repo{
		db:                  db,
		logger:              logger,
		redisCache:          redisClient,
		externalRoleService: external_services.InitializeExternalRoleService(redisClient, logger),
	}
}
