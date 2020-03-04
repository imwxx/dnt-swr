package lib

import (
    "gopkg.in/redis.v4"
    //"fmt"
    //"reflect"
)

var RedisClient *redis.Client

func InitRedis(redisInfo string) (*redis.Client, error) {
    var client *redis.Client
    client = redis.NewClient(&redis.Options{
        Addr:     redisInfo,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    _, err := client.Ping().Result()
    if err != nil {
        return client, err
    } else {
        return client, nil
    }
}
