package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"

	structs "github.com/traderboy/arcrestgo/structs"
)

//_ "github.com/mattn/go-sqlite3"
//_ "github.com/traderboy/arcrestgo/controllers"
//_ "github.com/traderboy/arcrestgo/models"

var Project structs.JSONConfig
var RootPath = "leasecompliance2016"

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

	Db, err = sql.Open("postgres", "user=postgres dbname=gis host=192.168.99.100")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Pinging Postgresql: ")
	log.Println(Db.Ping)
	GetConfiguration()

	RootPath = pwd + string(os.PathSeparator) + RootPath //+ string(os.PathSeparator)
	DataPath = RootPath                                  //+ string(os.PathSeparator)        //+ defaultService + string(os.PathSeparator) + "services" + string(os.PathSeparator)
	ReplicaPath = RootPath                               //+ string(os.PathSeparator)     //+ defaultService + string(os.PathSeparator) + "replicas" + string(os.PathSeparator)
	AttachmentsPath = RootPath                           //+ string(os.PathSeparator) //+ defaultService + string(os.PathSeparator) + "attachments" + string(os.PathSeparator)
	log.Println("Data path: " + DataPath)
	log.Println("Replica path: " + ReplicaPath)
	log.Println("Attachments path: " + AttachmentsPath)

}

func GetConfiguration() {
	sql := "select json from catalog where name=$1"
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
		GetConfigurationFromFile()
		return
	}
	err = json.Unmarshal(str, &Project)
	if err != nil {
		log.Println("Error parsing configuration table")
		log.Println(err.Error())
		GetConfigurationFromFile()
		return
	}
}
func GetConfigurationFromFile() {
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
	sql := "select json from services where service like $1 and name=$2 and layerid=$3 and type=$4"
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
	sql := "select json from catalog where name=$1 and type=$2"
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
	sql := "update services set json=$5 where service like $1 and name=$2 and layerid=$3 and type=$4"
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
	sql := "update catalog set json=$3 where name=$1 and type=$2"
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
