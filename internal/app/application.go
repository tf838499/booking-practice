package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	// "github.com/golang-migrate/migrate/v4/database/postgres"
	"ticket-booking/internal/adapter/rabbitMQ"
	RedisCache "ticket-booking/internal/adapter/redisClient"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
	"ticket-booking/internal/app/service/seats"

	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"ticket-booking/internal/app/service/customers"

	_ "github.com/lib/pq"
	// "github.com/chatbotgang/go-clean-architecture-template/internal/adapter/repository/postgres"
	// "github.com/chatbotgang/go-clean-architecture-template/internal/adapter/server"
	// "github.com/chatbotgang/go-clean-architecture-template/internal/app/service/auth"
	// "github.com/chatbotgang/go-clean-architecture-template/internal/app/service/barter"
)

type Application struct {
	Params          ApplicationParams
	SeatService     *seats.SeatService
	CustomerService *customers.CustomerService
}

type ApplicationParams struct {
	Env         string
	DatabaseDSN string
	DBHost      string
	DBPort      string
	DBUser      string
	DBname      string
	DBPassword  string

	RedisHost     string
	RedisPort     []string
	Redisname     int
	RedisPassword string
	RedisPoolSize int
}

func MustNewApplication(ctx context.Context, wg *sync.WaitGroup, params ApplicationParams) *Application {
	app, err := NewApplication(ctx, wg, params)
	if err != nil {
		log.Panicf("fail to new application, err: %s", err.Error())
	}
	return app
}

func NewApplication(ctx context.Context, wg *sync.WaitGroup, params ApplicationParams) (*Application, error) {
	// // Create repositories
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		params.DBHost,
		params.DBPort,
		params.DBUser,
		params.DBname,
		params.DBPassword,
	)
	fmt.Println(dsn)
	db := sqlx.MustOpen("postgres", dsn)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	pgRepo := PostgresDB.NewPostgresRepository(db)

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: params.RedisPort,

		ReadOnly:      true,
		RouteRandomly: true,
		PoolSize:      200,
	})

	err := client.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
	redisRepo := RedisCache.NewRedisRepository(client)

	mqConn, err := amqp.Dial("amqp://ticket:tf838499@localhost:5672/")
	if err != nil {
		panic(err)
	}
	// defer mqConn.Close()
	ch, err := mqConn.Channel()
	if err != nil {
		panic(err)
	}
	// defer ch.Close()
	rabbitMQ, err := rabbitMQ.NewRabbitMQ(ch)
	if err != nil {
		panic(err)
	}
	app := &Application{
		Params: params,
		SeatService: seats.NewSeatService(ctx, seats.SeatServiceParam{
			SeatRepo:  pgRepo,
			SeatRedis: redisRepo,
			DataRepo:  pgRepo,
		}),
		CustomerService: customers.NewCustomerService(ctx, customers.CustomerServiceParam{
			OrderRepo:     pgRepo,
			CustomerRedis: redisRepo,
			OrderMQ:       rabbitMQ,
		}),
	}

	return app, nil
}
