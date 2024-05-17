package main

import (
	"ClubTennis/config"
	"ClubTennis/initializers"
)

func main() {
	config, err := config.LoadConfig("config/config.json")
	if err != nil {
		panic("fatal: could not load config")
	}
	r := initializers.GetEngine(config)

	r.Run(":8080")
}
