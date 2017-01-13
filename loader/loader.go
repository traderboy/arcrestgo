package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	structs "github.com/traderboy/arcrestgo/structs"
)

//data sources
const (
	PGSQL   = 1 //Postgresql
	SQLITE3 = 2 //Sqlite 3
)

//DbSource is the current data source
var DbSource = PGSQL

//Db is the database object
var Db *sql.DB

//_ "github.com/mattn/go-sqlite3"
func main() {
	/*
		pwd, err := os.Getwd()
		if err != nil {
			log.Println("Unable to get current directory")
		}

		cmd := filepath.Clean(pwd + "\\..\\gdal\\ogr2ogr.exe")
		fmt.Println(cmd)
		//gdal_data := filepath.Clean(pwd + "\\..\\gdal\\gdal-data\\")
		params := strings.Split("--formats -version", " ")
		// -lco LAUNDER=NO -lco FID=OBJECTID -preserve_fid --config OGR_SQLITE_CACHE 1024 --config OGR_SQLITE_SYNCHRONOUS OFF -gt 65536 --config GDAL_DATA " + gdal_data + " -f \"SQLITE\" -overwrite"
		//fmt.Println(cmd)

		c, err := exec.Command(cmd, params...).Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(string(c))

		return
	*/
	if len(os.Args) > 1 {
		loadServices()
		return
	}
	fmt.Println("Usage:")
	fmt.Println("go run loader.go -file <path to single file> -dir <path to services> [-sqlite <filename> | -pgsql <dbname>]")
}

func loadServices() {
	var CreateTables = true
	var inputSource string
	configFile := ""
	//isSingleFile := false
	var DbName string
	var err error
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-sqlite" {
				DbSource = SQLITE3
				if len(os.Args) > i {
					DbName = os.Args[i+1]
				} else {
					DbName = "../arcrest.sqlite"
				}
				//if sqlite db exists, don't create
				if _, err := os.Stat(DbName); err == nil {
					CreateTables = false
					// path/to/whatever exists
				}

			} else if os.Args[i] == "-pgsql" && len(os.Args) > i {
				DbSource = PGSQL
				if len(os.Args) > i {
					DbName = os.Args[i+1]
				} else {
					DbName = "user=postgres dbname=gis host=192.168.99.100"
				}
			} else if os.Args[i] == "-dir" && len(os.Args) > i {
				inputSource = os.Args[i+1]
				configFile = inputSource + string(os.PathSeparator) + "config.json"

			} else if os.Args[i] == "-file" && len(os.Args) > i {
				inputSource = os.Args[i+1]
				//isSingleFile = true
				if filepath.Ext(inputSource) != ".json" {
					fmt.Println("Invalid input file.  Must be a .json file")
					return
				}
			}
		}
	}

	/*
		for _, f := range services {
			log.Println(f.Name())
		}
	*/
	/*
		for _, f := range services {
			log.Println(f)
		}
	*/
	//log.Println(files)
	/*
		inputSource += string(os.PathSeparator) + "*"
		services, err := filepath.Glob(inputSource)
		if err != nil {
			log.Println(err.Error())
			return
		}
	*/
	if DbSource == PGSQL {
		Db, err = sql.Open("postgres", DbName)
		if err != nil {
			log.Fatal(err)
		}
	} else if DbSource == SQLITE3 {
		Db, err = sql.Open("sqlite3", DbName)
		if err != nil {
			log.Fatal(err)
		}
		defer Db.Close()
	}
	/*
		Db, err := sql.Open("postgres", "user=postgres dbname=gis host=192.168.99.100")
		if err != nil {
			log.Fatal(err)
		}
	*/

	if CreateTables {
		sql := "DROP TABLE IF EXISTS catalog"
		_, err = Db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}
		sql = "DROP TABLE IF EXISTS services"
		_, err = Db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}
		if DbSource == PGSQL {
			sql = "CREATE TABLE IF NOT EXISTS catalog (id serial, name text, type text, json jsonb)"
		} else if DbSource == SQLITE3 {
			sql = "CREATE TABLE IF NOT EXISTS catalog (id INTEGER PRIMARY KEY AUTOINCREMENT, name text, type text, json text)"
		}
		_, err = Db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}
		if DbSource == PGSQL {
			sql = "CREATE TABLE IF NOT EXISTS services (id serial, service text,name text, layerid int, type text,json jsonb)"
		} else if DbSource == SQLITE3 {
			sql = "CREATE TABLE IF NOT EXISTS services (id INTEGER PRIMARY KEY AUTOINCREMENT, service text,name text, layerid int, type text,json text)"
		}
		_, err = Db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)
		}
	}

	var fi os.FileInfo
	fi, err = os.Stat(inputSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	var services []os.FileInfo
	if fi.IsDir() {
		services, err = ioutil.ReadDir(inputSource)
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		//is a file
		//isSingleFile = true
		//make exception if single file is a catalog item, not a service item
		//services = []os.FileInfo{fi}
		//need to set the inputSource to the Directory
		//inputSource = filepath.Dir(inputSource)
		inputSource = strings.Replace(inputSource, "\\", "/", -1)
		idx := strings.Index(inputSource, "/services/")

		file, err1 := ioutil.ReadFile(inputSource)
		if err1 != nil {
			fmt.Printf("// error while reading file %s\n", fi)
			fmt.Printf("File error: %v\n", err1)
		}

		//fmt.Println(file)
		if idx > -1 {
			//is a service
			idx1 := strings.LastIndex(inputSource[0:idx], "/") + 1
			schema := inputSource[idx1:idx]
			log.Println("Schema name: " + schema)
			LoadService(fi.Name(), schema, file)
		} else {
			//is a catalog item
			LoadCatalog(fi.Name(), file)
		}
		return

	}

	/*
		sql = "TRUNCATE catalog"
		_, err = Db.Exec(sql)
		if err != nil {
			log.Println(err.Error())
			log.Println(sql)

		}

		sql = "TRUNCATE services"
		_, err = Db.Exec(sql)
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
			jsonFilesPath := inputSource + string(os.PathSeparator) + schema + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "*.json"
			//inputSource = "*"
			files, err := filepath.Glob(jsonFilesPath)
			if err != nil {
				log.Println(err.Error())
				return
			}
			//log.Println(files)
			for _, jsonFile := range files {
				stat, _ := os.Stat(jsonFile)

				//jsonFile := inputSource + string(os.PathSeparator) + file
				if filepath.Ext(jsonFile) == ".json" {
					file, err := ioutil.ReadFile(jsonFile)
					if err != nil {
						fmt.Printf("// error while reading file %s\n", f)
						fmt.Printf("File error: %v\n", err)
					}
					//fmt.Println(file)
					LoadService(stat.Name(), schema, file)

					//log.Println(filepath.Base(f))
					//fileName := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
					//tableName := strings.Replace(fileName, ".", "_", -1)
					//log.Println("Loading table: " + schema + tableName)

				}

				//loadJSON2Postgresql(Db, schema, f)
			}
		} else {
			log.Println("Is file: " + f.Name())
			//loadJSON2Postgresql(Db, "", inputSource+string(os.PathSeparator)+f.Name())

			if filepath.Ext(f.Name()) == ".json" {
				jsonFile := inputSource + string(os.PathSeparator) + f.Name()
				file, err := ioutil.ReadFile(jsonFile)
				if err != nil {
					fmt.Printf("// error while reading file %s\n", f)
					fmt.Printf("File error: %v\n", err)
				}
				LoadCatalog(f.Name(), file)
			}
		}
		// filepath.Dir(f)
	}
	if configFile != "" {
		LoadGISDocker(configFile, DbName)
	}

}

func LoadGIS(configFile string, DbName string) {
	project := LoadConfigurationFromFile(configFile)
	fmt.Println("Source FGDB: " + project.FGDB)
	dbPath, _ := filepath.Abs(DbName)
	for name := range project.Services {
		fmt.Println("Service name: " + name)
		for _, layer := range project.Services[name]["layers"] {
			if layer["type"] == "layer" {
				fmt.Println("ogr2ogr -lco LAUNDER=NO -lco FID=OBJECTID -preserve_fid --config OGR_SQLITE_CACHE 1024 --config OGR_SQLITE_SYNCHRONOUS OFF -gt 65536 --config GDAL_DATA \"../gdal/gdal-data\" -f \"SQLITE\" \"" + dbPath + "\" \"" + project.FGDB + "\" " + layer["data"].(string))
			} else {
				fmt.Println("ogr2ogr -lco LAUNDER=NO -lco FID=OBJECTID -preserve_fid --config OGR_SQLITE_CACHE 1024 --config OGR_SQLITE_SYNCHRONOUS OFF -gt 65536 --config GDAL_DATA \"../gdal/gdal-data\" -f \"SQLITE\" \"" + dbPath + "\" \"" + project.FGDB + "\" " + layer["data"].(string) + " -nlt None -overwrite")
			}
		}
	}
}

func LoadGISDocker(configFile string, DbName string) {
	project := LoadConfigurationFromFile(configFile)
	fFullPath, _ := filepath.Abs(DbName)
	dbFullPath, _ := filepath.Abs(project.FGDB)

	fmt.Println("docker cp " + fFullPath + " determined_pare:/data")
	fmt.Println("docker cp " + dbFullPath + " determined_pare:/data")

	fmt.Println("Source FGDB: " + project.FGDB)
	//dbPath, _ := filepath.Abs(DbName)

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get current directory")
	}

	cmd := filepath.Clean(pwd + "\\..\\gdal\\ogr2ogr.exe")
	fmt.Println(cmd)
	gdal_data := filepath.Clean(pwd + "\\..\\gdal\\gdal-data\\")
	//fgdb := filepath.Base(project.FGDB)
	fgdb, _ := filepath.Abs(project.FGDB)

	for name := range project.Services {
		fmt.Println("Service name: " + name)
		for _, layer := range project.Services[name]["layers"] {
			if layer["type"] == "layer" {
				params := " -lco LAUNDER=NO -lco FID=OBJECTID -preserve_fid --config OGR_SQLITE_CACHE 1024 --config OGR_SQLITE_SYNCHRONOUS OFF -gt 65536 --config GDAL_DATA \"" + gdal_data + "\" -f \"Postgresql\" PG:\"" + DbName + "\" \"" + fgdb + "\" \"" + layer["data"].(string) + "\" -overwrite"
				fmt.Println(cmd + params)

				c, err := exec.Command("cmd.exe", "/c", cmd+params).Output()
				if err != nil {
					fmt.Println("Error: ", err)
				}
				fmt.Println(string(c))

			} else {
				params := " -lco LAUNDER=NO -lco FID=OBJECTID -preserve_fid --config OGR_SQLITE_CACHE 1024 --config OGR_SQLITE_SYNCHRONOUS OFF -gt 65536 --config GDAL_DATA \"" + gdal_data + "\" -f \"Postgresql\" PG:\"" + DbName + "\" \"" + fgdb + "\" \"" + layer["data"].(string) + "\" -nlt None -overwrite"
				fmt.Println(cmd + params)

				c, err := exec.Command("cmd.exe", "/c", cmd+params).Output()
				if err != nil {
					fmt.Println("Error: ", err)
				}
				fmt.Println(string(c))
			}
		}
	}
}

func LoadConfigurationFromFile(configFile string) structs.JSONConfig {
	//configFile = RootPath + string(os.PathSeparator) + "config.json"
	//var json []byte
	var project structs.JSONConfig
	file, err1 := ioutil.ReadFile(configFile)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", configFile)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal(file, &project)
	if err2 != nil {
		log.Println("Error reading configuration file: " + configFile)
		log.Println(err2.Error())
	}
	return project
}

func LoadCatalog(fileName string, file []byte) {
	names := strings.Split(fileName, ".")
	name := names[0]
	dtype := ""
	if len(names) > 1 && names[1] != "json" {
		dtype = names[1]
	}

	log.Printf("Loading service: %v/%v", name, dtype)
	json := strings.Replace(string(file), "'", "''", -1)
	json = strings.Replace(json, "\xa0", "", -1)
	json = strings.Replace(json, "\n", "", -1)

	sql := "INSERT INTO catalog(name,type,json) VALUES($1,$2,$3)"
	_, err := Db.Exec(sql, name, dtype, json)
	if err != nil {
		log.Println(err.Error())
		log.Println(sql)
	}
}

func LoadService(fileName string, schema string, file []byte) {
	sql := "INSERT INTO services(service,name,layerid,type,json) VALUES($1,$2,$3,$4,$5)"
	//" + fileName + "','" + strings.Replace(string(file), "'", "''", -1) + "')"
	//log.Println(sql)
	names := strings.Split(fileName, ".")
	layerid := -1
	dtype := ""
	name := names[0]
	var err error

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
	log.Printf("Loading service item: %v/%v/%v", name, layerid, dtype)
	json := strings.Replace(string(file), "'", "''", -1)
	json = strings.Replace(json, "\xa0", "", -1)
	json = strings.Replace(json, "\n", "", -1)
	_, err = Db.Exec(sql, schema, name, layerid, dtype, json)

	if err != nil {
		log.Println(err.Error())
		log.Println(sql)
		//return
	}
}

/*
func loadJSON2Postgresql(Db *sql.Db, schema string, name string, layerid string, dtype string, f string) {
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
		_, err = Db.Exec(sql, schema, name, layerid, dtype, file)

		if err != nil {
			log.Println(err.Error())
			//log.Println(sql)
			return
		}
	}

	//i++

}
func _loadJSON2Postgresql(Db *sql.Db, schema string, files []string) {

	if len(schema) > 0 {
		sql := "CREATE schema IF NOT EXISTS " + schema
		//log.Println(sql)
		_, err4 := Db.Exec(sql)
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
			_, err := Db.Exec(sql)
			if err != nil {
				log.Println(err.Error())
				log.Println(sql)
				continue
			}

			sql = "CREATE TABLE " + schema + tableName + " (id serial, name text, json jsonb)"
			//log.Println(sql)
			_, err3 := Db.Exec(sql)
			if err3 != nil {
				log.Println(err3.Error())
				log.Println(sql)
				continue
			}

			sql = "INSERT INTO " + schema + tableName + "(name, json) VALUES('" + fileName + "','" + strings.Replace(string(file), "'", "''", -1) + "')"
			//log.Println(sql)
			_, err2 := Db.Exec(sql)

			if err2 != nil {
				log.Println(err2.Error())
				//log.Println(sql)
				continue
			}
		}

		i++
	}



	//   	age := 21
	//  	rows, err := Db.Query("SELECT name FROM users WHERE age = $1", age)
	//  	var userid int
	//  err := Db.QueryRow(`INSERT INTO users(name, favorite_fruit, age)
	//  	VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)

}
*/
