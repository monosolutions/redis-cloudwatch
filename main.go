package main

import (
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var version = "0.0.1"

func main() {
	configuration := ConfigurationService{}
	config := configuration.Read()
	redis := NewRedisClient(config.Password, config.Address, config.Port)
	cw := NewCloudWatchClient("eu-central-1", config.AccessKey, config.SecretKey)

	app := cli.NewApp()
	app.Name = "redis-metrics"
	app.Usage = "CLI for sending redis metrics to CloudWatch"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:  "verbose, V",
			Usage: "verbose",
		},
		cli.BoolTFlag{
			Name:  "DontSend, D",
			Usage: "Don't send to CloudWatch",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Configure CLI",
			Flags:   []cli.Flag{},
			Subcommands: []cli.Command{
				{
					Name:  "set",
					Usage: "set a configuration",
					Action: func(c *cli.Context) error {
						if len(c.Args()[0]) < 2 {
							color.Red("Please include which value you want to set, '$ support config set token value'")
							return nil
						}
						return configuration.Set(c.Args()[0], c.Args()[1])

					},
				},
			},
			Action: func(c *cli.Context) error {

				return configuration.Run(c)

			},
		},
		{
			Name:    "metrics",
			Aliases: []string{"m"},
			Usage:   "getting metrics",
			Flags:   []cli.Flag{},
			Subcommands: []cli.Command{
				{
					Name:  "stat",
					Usage: "Getting General Stats",
					Action: func(c *cli.Context) error {
						stats := redis.Stats()
						for key, value := range stats.ToMap() {
							val, err := strconv.ParseFloat(value, 32)
							log.Printf("%v - %v", key, val)
							_, err = cw.SendMetric("redis-stats", key, val)
							if err != nil {
								log.Fatal(err)
							}

						}
						return nil
					},
				},
				{
					Name:  "cpu",
					Usage: "Getting Cpu Metrics",
					Action: func(c *cli.Context) error {

						cpu := redis.Cpu()
						for key, value := range cpu.ToMap() {
							val, err := strconv.ParseFloat(value, 32)
							log.Printf("%v - %v", key, val)
							_, err = cw.SendMetric("redis-cpu", key, val)
							if err != nil {
								log.Fatal(err)
							}

						}
						return nil
					},
				},
				{
					Name:  "replication",
					Usage: "Getting Replication Metrics",
					Action: func(c *cli.Context) error {

						replication := redis.Replication()
						for key, value := range replication.ToMap() {
							val, err := strconv.ParseFloat(value, 32)
							log.Printf("%v - %v", key, val)
							_, err = cw.SendMetric("redis-replication", key, val)
							if err != nil {
								log.Fatal(err)
							}

						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				stats := redis.Stats()
				for key, value := range stats.ToMap() {
					log.Printf("%v - %v \n", key, value)
					val, err := strconv.ParseFloat(value, 32)
					_, err = cw.SendMetric("redis-stats", key, val)
					if err != nil {
						log.Fatal(err)
					}

				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
