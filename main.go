package main

import (
	"Scrunchy/api"
	"Scrunchy/initializers"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	server := api.NewApiServer(":3000")
	server.Run()
}
