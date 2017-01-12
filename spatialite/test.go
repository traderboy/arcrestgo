package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const DbName = "D:/workspace/go/src/github.com/traderboy/arcrestgo/leasecompliance2016/leasecompliance2016/replicas/leasecompliance2016.geodatabase"

var Db *sql.DB

func main() {
	var err error
	Db, err = sql.Open("sqlite3", DbName)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Sqlite database: " + DbName)
	sql := "select count(*) from grazing"
	log.Printf(sql)
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//var json []byte
	var ret int
	err = stmt.QueryRow().Scan(&ret)
	if err != nil {
		log.Println(err.Error())
		//log.Println(sql)
	}
	log.Println(ret)

}
