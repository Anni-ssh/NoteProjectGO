package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")

func main() {
	// Параметры подключения к базе данных
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/dbname")
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}

	defer db.Close()

	q := `CREATE TABLE IF NOT EXISTS Users
    (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        age INTEGER NOT NULL,
        password TEXT NOT NULL,
        superUser INTEGER DEFAULT 0
    );`

	_, execErr := db.Exec(q)
	if execErr != nil {
		return fmt.Errorf("Ошибка создания таблицы Users в БД: %v", execErr)
	}
	return nil

}
