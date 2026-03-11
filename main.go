package main

import (
	"fmt"
	"go_finsl_project/Go/pkg/api"
	"go_finsl_project/Go/pkg/db"
	"go_finsl_project/Go/pkg/server"
	"log"
)

func main() {
	err := db.Init("scheduler.db")
	defer db.DB.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	api.Init()
	err = server.ServerStart()
	if err != nil {
		log.Fatal(err)
		return
	}
}
