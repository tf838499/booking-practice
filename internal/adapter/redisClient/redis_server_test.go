package RedisCache

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var testClient *redis.ClusterClient

const testRedisName = "redis_test"

func getTestRedis() *redis.ClusterClient {
	return testClient
}
func initRepository(t *testing.T, db *redis.ClusterClient, files ...string) (repo *RedisRepository) {

	// Setup DB again

	return NewRedisRepository(db)
}

func TestMain(m *testing.M) {
	_, closeDB, err := buildTestRedis()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeDB()
	// db := getTestPostgresDB()
	// repo := initRepository(t, db, testdata.Path(testdata.TestDataGood))
	// fmt.Println(repo)
	// testClient = db
	m.Run()
}
func buildTestRedis() (*redis.Client, func(), error) {
	db, cb, err := startRedisContainer()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to start redis container")
	}
	return db, cb, nil
}

func startRedisContainer() (*redis.Client, func(), error) {
	// new a docker pool
	var db *redis.Client
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}
	// start a PG container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.0.5-alpine",
		Env: []string{
			"ENV=develop",
			"requirepass",
		},
		Cmd: []string{"redis-server", "--requirepass", testRedisName},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start redis: %s", err)
	}

	// Get host and port(random) info from the postgres container
	// hostAndPort := resource.GetHostPort("6379/tcp")

	// build a call back function to destroy the docker pool
	cb := func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		db = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
			Password: testRedisName,
			DB:       1,  // use default DB
			PoolSize: 10, // 連接詞數量
		})
		err = db.Ping(context.Background()).Err()
		if err != nil {
			return err
		}
		return err
	}); err != nil {
		cb()
		return nil, nil, fmt.Errorf("could not connect to redis: %s", err)
	}

	return db, cb, nil
}
