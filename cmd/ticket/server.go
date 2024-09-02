package ticket

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"sync"
	"syscall"
	"ticket-booking/internal/app"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	AppName    = "simple-store"
	AppVersion = "unknown_version"
	AppBuild   = "unknown_build"
)

func StartStoreServer() {

	cfg := initAppConfig()
	rootLogger := initRootLogger(cfg.LogLevel, cfg.Env)
	fmt.Println(cfg)

	// // Create root context
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())
	rootCtx = rootLogger.WithContext(rootCtx)

	rootLogger.Info().
		Str("version", AppVersion).
		Str("build", AppBuild).
		Msgf("Launching %s", AppName)

	wg := sync.WaitGroup{}
	// Create application
	app := app.MustNewApplication(rootCtx, &wg, app.ApplicationParams{
		Env:         cfg.Env,
		DatabaseDSN: cfg.DatabaseHost,
		DBHost:      cfg.DatabaseHost,
		DBPort:      cfg.DatabasePort,
		DBUser:      cfg.DatabaseUser,
		DBname:      cfg.DatabaseName,
		DBPassword:  cfg.DatabasePasswd,

		RedisHost:     cfg.RedisHost,
		RedisPort:     cfg.RedisPort,
		Redisname:     cfg.RedisDBname,
		RedisPassword: cfg.RedisPasswd,
		RedisPoolSize: cfg.RedisPoolSize,
	})

	// Run server
	wg.Add(1)
	runHTTPServer(rootCtx, &wg, cfg.Port, app)

	//  Listen to SIGTERM/SIGINT to close
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	<-gracefulStop
	rootCtxCancelFunc()

	// Wait for all services to close with a specific timeout
	var waitUntilDone = make(chan struct{})
	go func() {
		wg.Wait()
		close(waitUntilDone)
	}()
	select {
	case <-waitUntilDone:
		rootLogger.Info().Msg("success to close all services")
	case <-time.After(10 * time.Second):
		rootLogger.Err(context.DeadlineExceeded).Msg("fail to close all services")
	}
}
func initRootLogger(levelStr, env string) zerolog.Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set logger time format
	const rfc3339Micro = "2006-01-02T15:04:05.000000Z07:00"
	zerolog.TimeFieldFormat = rfc3339Micro

	serviceName := fmt.Sprintf("%s-%s", AppName, env)
	rootLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return rootLogger
}

type AppConfig struct {
	// General configuration
	Env      string
	LogLevel string

	// Database configuration
	DatabaseHost   string
	DatabasePort   string
	DatabaseUser   string
	DatabaseName   string
	DatabasePasswd string
	// HTTP configuration
	Port int

	// redis configuration
	RedisHost     string
	RedisPort     []string
	RedisPasswd   string
	RedisDBname   int
	RedisPoolSize int
}

func initAppConfig() *AppConfig {
	// Setup basic application information

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		print(err.Error())
	}
	var config AppConfig

	config.Port = viper.GetInt("application.port")
	// postgres
	config.DatabaseHost = viper.GetString("application.db.HOST")
	config.DatabasePort = viper.GetString("application.db.PORT")
	config.DatabaseUser = viper.GetString("application.db.USER")
	config.DatabaseName = viper.GetString("application.db.NAME")
	config.DatabasePasswd = viper.GetString("application.db.PASSWD")
	// redis
	config.RedisHost = viper.GetString("application.redis.HOST")
	config.RedisPort = viper.GetStringSlice("application.redis.PORT")
	config.RedisPasswd = viper.GetString("application.redis.PASSWD")
	config.RedisPoolSize = viper.GetInt("application.redis.POOLSIZE")
	config.RedisDBname = viper.GetInt("application.redis.DEFAULTDB")
	//server
	config.Env = "staging"
	config.LogLevel = "info"
	return &config
}
