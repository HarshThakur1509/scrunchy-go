package main

import (
	"github.com/HarshThakur1509/scrunchy-go/api"
	"github.com/HarshThakur1509/scrunchy-go/initializers"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	server := api.NewApiServer(":3000")
	server.Run()
}
