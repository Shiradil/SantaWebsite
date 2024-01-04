package main

import (
	"SantaWeb/cmd"
	"SantaWeb/db"
)

func main() {
	// подключаем монгодб
	db.DbConnection()
	// раним сервак
	cmd.RunServer()
}
