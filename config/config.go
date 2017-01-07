package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	structs "github.com/traderboy/arcrestgo/structs"
)

//_ "github.com/mattn/go-sqlite3"
//_ "github.com/traderboy/arcrestgo/controllers"
//_ "github.com/traderboy/arcrestgo/models"
//var DbSource = "sqlite3"
const (
	PGSQL   = 1
	SQLITE3 = 2
)

var Schema = "postgres."
var DbSource = PGSQL
var Project structs.JSONConfig
var RootPath = "leasecompliance2016"
var HTTPPort = ":80"
var HTTPSPort = ":443"

//"github.com/gin-gonic/gin"
//Db is the SQLITE databa se object

var configFile = RootPath + string(os.PathSeparator) + "config.json"
var ArcGisVersion = "3.8"

var Db *sql.DB
var port = ":8080"

var DataPath = RootPath        //+ string(os.PathSeparator)        //+ string(os.PathSeparator) //+ "services"
var ReplicaPath = RootPath     //+ string(os.PathSeparator)     //+ "replicas"
var AttachmentsPath = RootPath //+ string(os.PathSeparator) //+ "attachments"

var CertificatePath = "ssl" + string(os.PathSeparator) + "agent2-cert.cert"

//var config map[string]interface{}
//var defaultService = ""
var uploadPath = ""
var Server = ""
var RefreshToken = "51vzPXXNl7scWXsw7YXvhMp_eyw_iQzifDIN23jNSsQuejcrDtLmf3IN5_bK0P5Z9K9J5dNb2yBbhXqjm9KlGtv5uDjr98fsUAAmNxGqnz3x0tvl355ZiuUUqqArXkBY-o6KaDtlDEncusGVM8wClk0bRr1-HeZJcR7ph9KU9khoX6H-DcFEZ4sRdl9c16exIX5lGIitw_vTmuomlivsGIQDq9thskbuaaTHMtP1m3VVnhuRQbyiZTLySjHDR8OVllSPc2Fpt0M-F5cPl_3nQg.."
var AccessToken = "XMdOaajM4srQWx8nQ77KuOYGO8GupnCoYALvXEnTj0V_ZXmEzhrcboHLb7hGtGxZCYUGFt07HKOTnkNLah8LflMDoWmKGr4No2LBSpoNkhJqc9zPa2gR3vfZp5L3yXigqxYOBVjveiuarUo2z_nqQ401_JL-mCRsXq9NO1DYrLw."

func Init() {
	//var err error
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get current directory")
	}
	RootPath = pwd + string(os.PathSeparator) + RootPath //+ string(os.PathSeparator)
	var DbName string

	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-sqlite" {
				DbSource = SQLITE3
				if len(os.Args) > i {
					DbName = os.Args[i+1]
				} else {
					DbName = pwd + string(os.PathSeparator) + "arcrest.sqlite"
				}
				Schema = ""
			} else if os.Args[i] == "-pgsql" && len(os.Args) > i {
				DbSource = PGSQL
				if len(os.Args) > i {
					DbName = os.Args[i+1]
				} else {
					DbName = "user=postgres dbname=gis host=192.168.99.100"
				}
				Schema = "postgres."
			} else if os.Args[i] == "-root" && len(os.Args) > i {
				RootPath, _ = filepath.Abs(os.Args[i+1])
			} else if os.Args[i] == "-p" && len(os.Args) > i {
				HTTPPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-https" && len(os.Args) > i {
				HTTPSPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-h" {
				fmt.Println("Usage:")
				fmt.Println("go run server.go -p HTTP Port -https HTTPS Port -root <path to service folder> -sqlite <path to service .sqlite> -pgsql <connection string for Postgresql> -h [show help]")
				os.Exit(0)
			}
		}
	}

	if DbSource == PGSQL {
		Db, err = sql.Open("postgres", DbName)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Postgresql database: " + DbName)
		log.Print("Pinging Postgresql: ")
		log.Println(Db.Ping)
	} else if DbSource == SQLITE3 {
		Db, err = sql.Open("sqlite3", DbName)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Sqlite database: " + DbName)
		//defer db.Close()
	}
	/*
		Db, err = sql.Open("postgres", "user=postgres DbSource=gis host=192.168.99.100")
		if err != nil {
			log.Fatal(err)
		}
	*/

	DataPath = RootPath        //+ string(os.PathSeparator)        //+ defaultService + string(os.PathSeparator) + "services" + string(os.PathSeparator)
	ReplicaPath = RootPath     //+ string(os.PathSeparator)     //+ defaultService + string(os.PathSeparator) + "replicas" + string(os.PathSeparator)
	AttachmentsPath = RootPath //+ string(os.PathSeparator) //+ defaultService + string(os.PathSeparator) + "attachments" + string(os.PathSeparator)

	log.Println("Root path: " + RootPath)
	log.Println("Data path: " + DataPath)
	log.Println("Replica path: " + ReplicaPath)
	log.Println("Attachments path: " + AttachmentsPath)
	LoadConfiguration()

}

func GetParam(i int) string {
	if DbSource == SQLITE3 {
		return "?"
	}
	return "$" + strconv.Itoa(i)
}

func LoadConfiguration() {
	sql := "select json from catalog where name=" + GetParam(1)
	log.Printf("Query: select json from catalog where name='config'")
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	var str []byte
	err = stmt.QueryRow("config").Scan(&str)
	if err != nil {
		log.Println("Error reading configuration table")
		log.Println(err.Error())
		LoadConfigurationFromFile()
		return
	}
	err = json.Unmarshal(str, &Project)
	if err != nil {
		log.Println("Error parsing configuration table")
		log.Println(err.Error())
		LoadConfigurationFromFile()
		return
	}
}

func LoadConfigurationFromFile() {
	configFile = RootPath + string(os.PathSeparator) + "config.json"
	//var json []byte
	file, err1 := ioutil.ReadFile(configFile)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", configFile)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal(file, &Project)
	if err2 != nil {
		log.Println("Error reading configuration file: " + configFile)
		log.Println(err2.Error())
	}
}

//GetArcService queries the database for service layer entries
func GetArcService(catalog string, service string, layerid int, dtype string) []byte {
	sql := "select json from services where service like " + GetParam(1) + " and name=" + GetParam(2) + " and layerid=" + GetParam(3) + " and type=" + GetParam(4)
	log.Printf("Query: select json from services where service like '%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	var json []byte
	err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&json)
	if err != nil {
		log.Println(err.Error())
		//log.Println(sql)
	}
	return json
}

//GetArcCatalog queries the database for top level catalog entries
func GetArcCatalog(service string, dtype string) []byte {
	sql := "select json from catalog where name=" + GetParam(1) + " and type=" + GetParam(2)
	log.Printf("Query: select json from catalog where name='%v' and type='%v'", service, dtype)

	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}

	var json []byte
	err = stmt.QueryRow(service, dtype).Scan(&json)
	if err != nil {
		log.Println(err.Error())
		//log.Println(sql)
	}

	return json
}

func SetArcService(json string, catalog string, service string, layerid int, dtype string) bool {
	sql := "update services set json=" + GetParam(5) + " where service like " + GetParam(1) + "and name=" + GetParam(2) + " and layerid=" + GetParam(3) + " and type=" + GetParam(4)
	log.Printf("Query: update services set json=<json> where service like '%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
	_, err = stmt.Exec(catalog, service, layerid, dtype, json)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

//GetArcCatalog queries the database for top level catalog entries
func SetArcCatalog(json string, service string, dtype string) bool {
	sql := "update catalog set json=$3 where name=" + GetParam(1) + " and type=" + GetParam(2)
	log.Printf("Query: update catalog set json=<json> where name='%v' and type='%v'", service, dtype)

	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(service, dtype, json)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}
