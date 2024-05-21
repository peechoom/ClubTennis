package main

import (
	"ClubTennis/config"
	"ClubTennis/initializers"
	"fmt"
)

func main() {
	err := config.LoadConfig("config/.env")
	if err != nil {
		panic(fmt.Sprintf("fatal: could not load config, %s", err.Error()))
	}
	r := initializers.GetEngine()

	r.Run(":8080")
}
