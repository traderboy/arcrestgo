package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"
)

func main() {
	DbName := ""
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if len(os.Args[i]) > 0 {
				DbName = os.Args[i]
			}
		}
	}
	Db, err := sql.Open("sqlite3", "file:"+DbName+"?PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	//Db.Exec(initializeStr)
	log.Print("Sqlite database: " + DbName)
	sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
	sql = "SELECT \"ItemInfo\" FROM \"GDB_ServiceItems\""
	log.Printf("Query: " + sql)
	/*
		    stmt, err := Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
	*/
	//var datasetName []byte
	//var itemId int
	var itemInfo []byte
	//var advDrawingInfo []byte

	rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
	if err != nil {
		log.Println("Error reading configuration table")
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		//err = rows.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
		err = rows.Scan(&itemInfo)
		fmt.Println(string(itemInfo))
		//fmt.Println(string(advDrawingInfo))
	}
	rows.Close() //good habit to close
	if false {
		sql = "SELECT \"ObjectID\", \"UUID\", \"Type\", \"Name\", \"PhysicalName\", \"Path\", \"Url\", \"Properties\", \"Defaults\", \"DatasetSubtype1\", \"DatasetSubtype2\", \"DatasetInfo1\", \"DatasetInfo2\", \"Definition\", \"Documentation\", \"ItemInfo\", \"Shape\" FROM \"GDB_Items\""
		sql = "SELECT  \"Definition\" FROM \"GDB_Items\""
		var definition []byte
		rows, err = Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
		if err != nil {
			log.Println("Error reading configuration table")
			log.Println(err.Error())
			return
		}
		for rows.Next() {
			err = rows.Scan(&definition)
			fmt.Println(string(definition))
		}
	}

	//fmt.Printf("%v: %v", string(datasetName), itemId)

	//fmt.Println(string(itemInfo))
	//fmt.Println(string(advDrawingInfo))
	/*
		sql = "SELECT \"ObjectID\", \"UUID\", \"Type\", \"Name\", \"PhysicalName\", \"Path\", \"Url\", \"Properties\", \"Defaults\", \"DatasetSubtype1\", \"DatasetSubtype2\", \"DatasetInfo1\", \"DatasetInfo2\", \"Definition\", \"Documentation\", \"ItemInfo\", \"Shape\" FROM \"GDB_Items\""

		log.Printf("Query: " + sql)
		stmt, err = Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		var datasetName []byte
		var itemId int
		var itemInfo []byte
		var advDrawingInfo []byte

		err = stmt.QueryRow("config").Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
		if err != nil {
			log.Println("Error reading configuration table")
			log.Println(err.Error())
			return
		}
		fmt.Printf("%v: %v", string(datasetName), itemId)

		fmt.Println(string(itemInfo))
		fmt.Println(string(advDrawingInfo))
	*/

	/*
		err = json.Unmarshal(str, &Project)
		if err != nil {
			log.Println("Error parsing configuration table")
			log.Println(err.Error())
			LoadConfigurationFromFile()
			return
		}
	*/

}
