package main

import (
	"ClubTennis/config"
	"ClubTennis/initializers"
)

func main() {
	err := config.LoadConfig("config/.env")
	if err != nil {
		panic("fatal: could not load config")
	}
	r := initializers.GetEngine()

	r.Run(":8080")
}
