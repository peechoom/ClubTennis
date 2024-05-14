package main

import "ClubTennis/initializers"

func main() {
	r := initializers.GetEngine()

	r.Run(":8080")
}
