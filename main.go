package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Printf("инициализация базы данных: %v", err)
		return
	}
	defer db.DB.Close()

	if err := server.Start(); err != nil {
		log.Printf("остановка сервера: %v", err)
	}
}
