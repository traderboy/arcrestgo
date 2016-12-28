package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) > 1 {
		loadServices(os.Args[1])
		return
	}
	fmt.Println("Usage:")
	fmt.Println("go run loader.go <path to services>")

}

func loadServices(jsonPath string) {
	services, err := ioutil.ReadDir(jsonPath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	/*
		for _, f := range services {
			log.Println(f)
		}
	*/

	//log.Println(files)
	/*
		jsonPath += string(os.PathSeparator) + "*"
		services, err := filepath.Glob(jsonPath)
		if err != nil {
			log.Println(err.Error())
			return
		}
	*/

	db, err := sql.Open("postgres", "user=postgres dbname=gis host=192.168.99.100")
	if err != nil {
		log.Fatal(err)
	}
	sql := "DROP TABLE IF EXISTS catalog"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
		log.Println(sql)

	}
	sql = "DROP TABLE IF EXISTS catalog"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
		log.Println(sql)

	}

	sql = "CREATE TABLE IF NOT EXISTS catalog (id serial, name text, type text, json jsonb)"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
		log.Println(sql)

	}

	sql = "CREATE TABLE IF NOT EXISTS services (id serial, service text,name text, layerid int, type text,json jsonb)"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
		log.Println(sql)
	}
	/*
		sql = "TRUNCATE catalog"
		_, err = db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}

		sql = "TRUNCATE services"
		_, err = db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}
	*/

	for _, f := range services {
		//if os.FileInfo(s).IsDir() {
		if f.IsDir() {
			//dirName := strings.TrimSuffix(filepath.Base(f), filepath.Dir(f))
			schema := f.Name()
			log.Println("Schema name: " + schema)
			//log.Println("Directory name: " + filepath.Dir(f))
			jsonFilesPath := jsonPath + string(os.PathSeparator) + schema + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "*.json"
			//jsonPath = "*"
			files, err := filepath.Glob(jsonFilesPath)
			if err != nil {
				log.Println(err.Error())
				return
			}
			//log.Println(files)
			for _, jsonFile := range files {
				stat, _ := os.Stat(jsonFile)

				//jsonFile := jsonPath + string(os.PathSeparator) + file
				if filepath.Ext(jsonFile) == ".json" {
					file, err := ioutil.ReadFile(jsonFile)
					if err != nil {
						fmt.Printf("// error while reading file %s\n", f)
						fmt.Printf("File error: %v\n", err)
					}
					//fmt.Println(file)

					//log.Println(filepath.Base(f))
					//fileName := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
					//tableName := strings.Replace(fileName, ".", "_", -1)
					//log.Println("Loading table: " + schema + tableName)
					sql := "INSERT INTO services(service,name,layerid,type,json) VALUES($1,$2,$3,$4,$5)"
					//" + fileName + "','" + strings.Replace(string(file), "'", "''", -1) + "')"
					//log.Println(sql)
					names := strings.Split(stat.Name(), ".")
					layerid := 0
					dtype := ""
					name := names[0]
					log.Println("Loading service item: " + name)
					if len(names) > 1 && names[1] != "json" {
						if layerid, _ = strconv.Atoi(names[1]); err == nil {
							//fmt.Printf("%d looks like a number.\n", layerid)
							if len(names) > 2 && names[2] != "json" {
								dtype = names[2]
							}
						} else {
							dtype = names[1]
						}
					}
					json := strings.Replace(string(file), "'", "''", -1)
					json = strings.Replace(json, "\xa0", "", -1)
					_, err = db.Exec(sql, schema, name, layerid, dtype, json)

					if err != nil {
						log.Println(err.Error())
						//log.Println(sql)
						//return
					}
				}

				//loadJSON2Postgresql(db, schema, f)
			}
		} else {
			log.Println("Is file: " + f.Name())
			//loadJSON2Postgresql(db, "", jsonPath+string(os.PathSeparator)+f.Name())
			jsonFile := jsonPath + string(os.PathSeparator) + f.Name()

			if filepath.Ext(jsonFile) == ".json" {
				file, err := ioutil.ReadFile(jsonFile)
				if err != nil {
					fmt.Printf("// error while reading file %s\n", f)
					fmt.Printf("File error: %v\n", err)
				}
				names := strings.Split(f.Name(), ".")
				name := names[0]
				dtype := ""
				if len(names) > 1 && names[1] != "json" {
					dtype = names[1]
				}
				log.Println("Loading service: " + name)
				json := strings.Replace(string(file), "'", "''", -1)
				json = strings.Replace(json, "\xa0", "", -1)

				sql := "INSERT INTO catalog(name,type,json) VALUES($1,$2,$3)"
				_, err = db.Exec(sql, name, dtype, json)
				if err != nil {
					log.Println(err.Error())
					//log.Println(sql)
				}
			}
		}
		// filepath.Dir(f)
	}

}

/*
func loadJSON2Postgresql(db *sql.DB, schema string, name string, layerid string, dtype string, f string) {
	//i := 0

	//read json file to string
	//log.Println(f)
	//log.Println(filepath.Ext(f))
	//log.Println(filepath.Clean(f))

	//=="json")

	if filepath.Ext(f) == ".json" {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("// error while reading file %s\n", f)
			fmt.Printf("File error: %v\n", err)
			return
		}
		//fmt.Println(file)

		//log.Println(filepath.Base(f))
		//fileName := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		//tableName := strings.Replace(fileName, ".", "_", -1)
		//log.Println("Loading table: " + schema + tableName)

		sql := "INSERT INTO services(service,name,layerid,type,json) VALUES($1,$2,$3,$4,$5)"

		//" + fileName + "','" + strings.Replace(string(file), "'", "''", -1) + "')"
		//log.Println(sql)
		_, err = db.Exec(sql, schema, name, layerid, dtype, file)

		if err != nil {
			log.Println(err.Error())
			//log.Println(sql)
			return
		}
	}

	//i++

}
func _loadJSON2Postgresql(db *sql.DB, schema string, files []string) {

	if len(schema) > 0 {
		sql := "CREATE schema IF NOT EXISTS " + schema
		//log.Println(sql)
		_, err4 := db.Exec(sql)
		if err4 != nil {
			log.Println(err4.Error())
			log.Println(sql)
		}
		schema += "."
	} else {
		schema = ""
	}

	i := 0
	for _, f := range files {
		//read json file to string
		//log.Println(f)
		//log.Println(filepath.Ext(f))
		//log.Println(filepath.Clean(f))

		//=="json")

		if filepath.Ext(f) == ".json" {
			file, err1 := ioutil.ReadFile(f)
			if err1 != nil {
				fmt.Printf("// error while reading file %s\n", f)
				fmt.Printf("File error: %v\n", err1)
				continue
			}
			//fmt.Println(file)

			log.Println(filepath.Base(f))
			fileName := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
			tableName := strings.Replace(fileName, ".", "_", -1)
			log.Println("Loading table: " + schema + tableName)
			sql := "DROP TABLE IF EXISTS " + schema + tableName
			//log.Println(sql)
			_, err := db.Exec(sql)
			if err != nil {
				log.Println(err.Error())
				log.Println(sql)
				continue
			}

			sql = "CREATE TABLE " + schema + tableName + " (id serial, name text, json jsonb)"
			//log.Println(sql)
			_, err3 := db.Exec(sql)
			if err3 != nil {
				log.Println(err3.Error())
				log.Println(sql)
				continue
			}

			sql = "INSERT INTO " + schema + tableName + "(name, json) VALUES('" + fileName + "','" + strings.Replace(string(file), "'", "''", -1) + "')"
			//log.Println(sql)
			_, err2 := db.Exec(sql)

			if err2 != nil {
				log.Println(err2.Error())
				//log.Println(sql)
				continue
			}
		}

		i++
	}



	//   	age := 21
	//  	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	//  	var userid int
	//  err := db.QueryRow(`INSERT INTO users(name, favorite_fruit, age)
	//  	VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)

}
*/
