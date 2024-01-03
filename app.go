package main

import (
	"SantaWeb/cmd"
	"SantaWeb/db"
)

func main() {
	// раним сервак
	cmd.RunServer()
	// подключаем монгодб
	db.DbConnection()
}
