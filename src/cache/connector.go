package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"os"
	"time"
)

var (
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPass = os.Getenv("REDIS_PASS")
)

var conn *redis.Client

func getConnection() *redis.Client {
	fmt.Println("Connecting to redis: ", RedisHost+":"+RedisPort)
	ctx, _ := GetContext()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
		Password: RedisPass,
		DB:       0,
	})
	status := client.Ping(ctx)
	if status.Err() != nil {
		panic(fmt.Sprintf("Could not connect to redis: %v", status.Err()))
	}
	return client
}

func GetConnection() *redis.Client {
	return conn
}

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func init() {
	conn = getConnection()
}
