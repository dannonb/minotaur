package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dannonb/go-network-monitor/config"
	"github.com/dannonb/go-network-monitor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var hostCollection *mongo.Collection = config.GetCollection(config.DB, "hosts")

func GetHosts() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var hosts []string
		defer cancel()

		results, err := hostCollection.Find(ctx, bson.M{})

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Host
			if err = results.Decode(&singleUser); err != nil {
				log.Fatal(err)
				return nil, err
			}

			hosts = append(hosts, singleUser.Hostname)
		}

		return hosts, nil

}

func AddSingleHostToRedis(host string) error {
	var ctx = context.Background()
	var key string = "hosts"

	redis := config.ConnectRedis()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err := redis.RPush(ctx, key, host).Result()
	if err != nil {
		return err
	}

	return nil
}

func AddToRedis(hosts []string) error {
	var ctx = context.Background()
	var key string = "hosts"

	redis := config.ConnectRedis()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	redisHosts := make([]interface{}, len(hosts))
	for i, h := range hosts {
		redisHosts[i] = h
	}

	_, pushErr := redis.RPush(ctx, key, redisHosts...).Result()
	if pushErr != nil {
		return pushErr
	}

	return nil
}

func UpdateCache() error {
	hosts, err := GetHosts()

	if err != nil {
		log.Fatal(err)
		return err
	}

	AddToRedis(hosts)
	return nil
}

func GetHostsFromCache() ([]string, error) {
	fmt.Println("GET HOSTS FROM CACHE")
	var ctx = context.Background()
	redis := config.ConnectRedis()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	pong, err := redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis:", pong)

	hosts, err := redis.LRange(ctx, "hosts", 0, -1).Result()
	if err != nil {
		fmt.Println("err from range")
		log.Fatal(err)
		return nil, err
	}
	return hosts, nil
}

func ClearCache() {
	var ctx = context.Background()
	redis := config.ConnectRedis()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	pong, err := redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis:", pong)


	err = redis.FlushDB(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cleared Redis Cache")
}