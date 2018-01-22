package main

import (
	"fmt"
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

type RedisClient struct {
	Password string
	Address  string
	Port     string
}

func NewRedisClient(password string, address string, port string) RedisClient {
	redis := RedisClient{
		Password: password,
		Address:  address,
		Port:     port,
	}
	return redis
}

func (redisClient RedisClient) do(field string) string {
	address := fmt.Sprintf("%v:%v", redisClient.Address, redisClient.Port)
	client, err := redis.Dial("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.Cmd("AUTH", redisClient.Password).Str(); err != nil {
		log.Fatal(err)
	}

	info, er := client.Cmd("INFO", field).Str()
	if er != nil {
		log.Fatal(er)
	}
	return info
}

func (redisClient RedisClient) Cpu() *Cpu {
	info := redisClient.do("cpu")

	cpu := &Cpu{}

	return cpu.Parse(info)
}

func (redisClient RedisClient) Stats() *Stats {
	info := redisClient.do("stats")

	stats := &Stats{}

	return stats.Parse(info)
}

func (redisClient RedisClient) Replication() *Replication {
	info := redisClient.do("replication")

	replication := &Replication{}

	return replication.Parse(info)
}
