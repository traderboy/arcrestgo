package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	structs "github.com/traderboy/arcrestgo/structs"
)

//"github.com/traderboy/arcrestgo/config"
//_ "github.com/mattn/go-sqlite3"
//_ "github.com/traderboy/arcrestgo/controllers"
//_ "github.com/traderboy/arcrestgo/models"
//var DbSource = "sqlite3"
const (
	PGSQL   = 1
	SQLITE3 = 2
	FILE    = 3
)

var Catalogs map[string]structs.Catalog

//map[string]Service
var DbSource = SQLITE3

var Schema = "postgres."

var Project structs.JSONConfig
var RootPath = "catalogs"

//leasecompliance2016
var RootName string
var HTTPPort = ":80"
var HTTPSPort = ":443"
var Pem = "ssl/reais.x10host.com.key.pem"
var Cert = "ssl/2_reais.x10host.com.crt"

//"github.com/gin-gonic/gin"
//Db is the SQLITE databa se object

var configFile = RootPath + string(os.PathSeparator) + "config.json"
var ArcGisVersion = "3.8"

var Db *sql.DB
var DbQuery *sql.DB

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

func Initialize() {
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
				//RootName = filepath.Base(os.Args[i+1])
			} else if os.Args[i] == "-p" && len(os.Args) > i {
				HTTPPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-https" && len(os.Args) > i {
				HTTPSPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-pem" && len(os.Args) > i {
				Pem = ":" + os.Args[i+1]
			} else if os.Args[i] == "-cert" && len(os.Args) > i {
				Cert = ":" + os.Args[i+1]
			} else if os.Args[i] == "-file" {
				DbSource = FILE
				LoadConfigurationFromFile()
			} else if os.Args[i] == "-h" {
				fmt.Println("Usage:")
				fmt.Println("go run server.go -p HTTP Port -https HTTPS Port -root <path to service folder> -sqlite <path to service .sqlite> -pgsql <connection string for Postgresql> -h [show help]")
				os.Exit(0)
			}
		}
	} else {
		//read all folder in catalogs
		/*
			files, _ := ioutil.ReadDir(RootPath)

			for _, f := range files {
				if f.IsDir() {
			}
			}
		*/

		LoadConfigurationFromFile()
		//RootPath, _ = filepath.Abs(os.Args[i+1])
		//RootName = filepath.Base(os.Args[i+1])
		if Project.DataSource == "pg" {
			DbSource = PGSQL
			DbName = Project.PG
			Schema = "postgres."
		} else if Project.DataSource == "sqlite" {
			DbSource = SQLITE3
			DbName = Project.SqliteDb
		} else if Project.DataSource == "file" {
			DbSource = FILE
		}
	}

	if DbSource == PGSQL {
		Db, err = sql.Open("postgres", DbName)
		if err != nil {
			log.Fatal(err)
		}
		DbQuery = Db
		log.Print("Postgresql database: " + DbName)
		log.Print("Pinging Postgresql: ")
		log.Println(Db.Ping)
		LoadConfiguration()
	} else if DbSource == SQLITE3 {
		/*
					initializeStr := `PRAGMA automatic_index = ON;
			        PRAGMA cache_size = 32768;
			        PRAGMA cache_spill = OFF;
			        PRAGMA foreign_keys = ON;
			        PRAGMA journal_size_limit = 67110000;
			        PRAGMA locking_mode = NORMAL;
			        PRAGMA page_size = 4096;
			        PRAGMA recursive_triggers = ON;
			        PRAGMA secure_delete = ON;
			        PRAGMA synchronous = NORMAL;
			        PRAGMA temp_store = MEMORY;
			        PRAGMA journal_mode = WAL;
			        PRAGMA wal_autocheckpoint = 16384;
					`
		*/
		//log.Println(initializeStr)
		//initializeStr = "PRAGMA synchronous = OFF;PRAGMA cache_size=100000;PRAGMA journal_mode=WAL;"
		//log.Println(initializeStr)

		Db, err = sql.Open("sqlite3", "file:"+DbName+"?PRAGMA journal_mode=WAL")
		if err != nil {
			log.Fatal(err)
		}
		//Db.Exec(initializeStr)
		log.Print("Sqlite database: " + DbName)
		//defer Db.Close()
		//Db.SetMaxOpenConns(1)

		LoadConfiguration()
		//get RootName
		DbQueryName := RootPath + string(os.PathSeparator) + RootName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + RootName + ".geodatabase"
		log.Println("DbQueryName: " + DbQueryName)
		DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+"?PRAGMA journal_mode=WAL")
		//DbQuery, err = sql.Open("sqlite3", DbQueryName)
		if err != nil {
			log.Fatal(err)
		}
		//defer DbQuery.Close()
		//DbQuery.SetMaxOpenConns(1)
		//log.Print("Sqlite database: " + DbQueryName)
		//DbQuery.Exec(initializeStr)
		//defer db.Close()
	}
	/*
		else if DbSource == FILE {
			LoadConfigurationFromFile()
		}
	*/

	/*
		Db, err = sql.Open("postgres", "user=postgres DbSource=gis host=192.168.99.100")
		if err != nil {
			log.Fatal(err)
		}
	*/

	DataPath = RootPath        //+ string(os.PathSeparator)        //+ defaultService + string(os.PathSeparator) + "services" + string(os.PathSeparator)
	ReplicaPath = RootPath     //+ string(os.PathSeparator)     //+ defaultService + string(os.PathSeparator) + "replicas" + string(os.PathSeparator)
	AttachmentsPath = RootPath //+ string(os.PathSeparator) //+ defaultService + string(os.PathSeparator) + "attachments" + string(os.PathSeparator)

	log.Println("Root catalog: " + RootName)
	log.Println("Root path: " + RootPath)
	log.Println("Data path: " + DataPath)
	log.Println("Replica path: " + ReplicaPath)
	log.Println("Attachments path: " + AttachmentsPath)

}

func GetParam(i int) string {
	if DbSource == SQLITE3 {
		return "?"
	}
	return "$" + strconv.Itoa(i)
}

func SetDatasource(newDatasource int) {
	if DbSource == newDatasource {
		return
	}
	if newDatasource == FILE {
		DbSource = FILE
		//close db
		Db.Close()
		return
	}

	var err error
	if newDatasource == PGSQL {

		Db, err = sql.Open("postgres", Project.PG)
		if err != nil {
			log.Fatal(err)
		}
		DbQuery = Db
		log.Print("Postgresql database: " + Project.PG)
		log.Print("Pinging Postgresql: ")
		log.Println(Db.Ping)
	} else if newDatasource == SQLITE3 {
		Db, err = sql.Open("sqlite3", "file:"+Project.SqliteDb+"?PRAGMA journal_mode=WAL")
		if err != nil {
			log.Fatal(err)
		}
		DbQueryName := RootPath + string(os.PathSeparator) + RootName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + RootName + ".geodatabase"
		log.Println("DbQueryName: " + DbQueryName)
		DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+"?PRAGMA journal_mode=WAL")

		if err != nil {
			log.Fatal(err)
		}
	}
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
	for i, _ := range Project.Services {
		RootName = i
		//RootName = val[i]
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
	for i, _ := range Project.Services {
		RootName = i
		//RootName = val[i]
	}

}

//GetArcService queries the database for service layer entries
func GetArcService(catalog string, service string, layerid int, dtype string) []byte {
	if DbSource == FILE {
		if len(service) > 0 {
			service += "."
		}
		sp := ""
		if layerid > -1 {
			sp = fmt.Sprint(layerid, ".")
		}

		if len(dtype) > 0 {
			if dtype == "data" && service == "content" {
				dtype = "items." + dtype + "."
			} else {
				dtype += "."
			}

		}
		jsonFile := fmt.Sprint(DataPath, string(os.PathSeparator), catalog, string(os.PathSeparator), "services", string(os.PathSeparator), service, sp, dtype, "json")
		file, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			log.Println(err)
		}
		return file
	}
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
	if DbSource == FILE {
		if len(service) > 0 {
			service += "."
		}

		if len(dtype) > 0 {
			dtype += "."
		}

		jsonFile := fmt.Sprint(DataPath, string(os.PathSeparator), service, dtype, "json")
		file, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			log.Println(err)
		}

		return file

	}
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

func SetArcService(json []byte, catalog string, service string, layerid int, dtype string) bool {
	if DbSource == FILE {
		if len(service) > 0 {
			service += "."
		}
		sp := ""
		if layerid > -1 {
			sp = fmt.Sprint(layerid, ".")
		}

		if len(dtype) > 0 {
			if dtype == "data" && service == "content" {
				dtype = "items." + dtype + "."
			} else {
				dtype += "."
			}
		}

		jsonFile := fmt.Sprint(DataPath, string(os.PathSeparator), catalog, string(os.PathSeparator), "services", string(os.PathSeparator), service, sp, dtype, "json")
		err := ioutil.WriteFile(jsonFile, json, 0644)
		if err != nil {
			return false
		}
		return true
	}

	sql := "update services set json=" + GetParam(1) + " where service like " + GetParam(2) + " and name=" + GetParam(3) + " and layerid=" + GetParam(4) + " and type=" + GetParam(5)
	log.Printf("Query: update services set json=<json> where service like '%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
	_, err = stmt.Exec(json, catalog, service, layerid, dtype)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

//GetArcCatalog queries the database for top level catalog entries
func SetArcCatalog(json []byte, service string, dtype string) bool {
	if DbSource == FILE {
		if len(service) > 0 {
			service += "."
		}
		if len(dtype) > 0 {
			dtype += "."
		}

		jsonFile := fmt.Sprint(DataPath, string(os.PathSeparator), service, dtype, "json")
		err := ioutil.WriteFile(jsonFile, json, 0644)
		if err != nil {
			return false
		}
		return true

	}

	sql := "update catalog set json=" + GetParam(1) + " where name=" + GetParam(2) + " and type=" + GetParam(3)
	log.Printf("Query: update catalog set json=<json> where name='%v' and type='%v'", service, dtype)

	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = stmt.Exec(json, service, dtype)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func GetArcQuery(catalog string, service string, layerid int, dtype string, objectIds string, where string) []byte {
	//objectIdsInt, _ := strconv.Atoi(objectIds)
	objectIdsArr := strings.Split(objectIds, ",")
	var objectIdsFloat = []float64{}
	for _, i := range objectIdsArr {
		j, err := strconv.ParseFloat(i, 64)
		if err != nil {
			panic(err)
		}
		objectIdsFloat = append(objectIdsFloat, j)
	}

	if DbSource == FILE {
		//config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json"

		jsonFile := fmt.Sprint(DataPath, string(os.PathSeparator), catalog+string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", layerid, ".query.json")
		log.Println(jsonFile)
		file, err1 := ioutil.ReadFile(jsonFile)
		if err1 != nil {
			log.Println(err1)
		}
		var srcObj structs.FeatureTable

		//map[string]map[string]map[string]
		err := json.Unmarshal(file, &srcObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(file))
			log.Println(err.Error())
		}
		var results []structs.Feature
		for _, i := range srcObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {
			if in_float_array(i.Attributes["OBJECTID"].(float64), objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		srcObj.Features = results
		jsonstr, err := json.Marshal(srcObj)
		if err != nil {
			log.Println(err)
		}
		return jsonstr
	} else if DbSource == PGSQL {
		sql := "select json from services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
		stmt, err := Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		var fields []byte
		err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("")
		}
		var featureObj structs.FeatureTable
		err = json.Unmarshal(fields, &featureObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		var results []structs.Feature
		for _, i := range featureObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {
			if in_float_array(i.Attributes["OBJECTID"].(float64), objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		featureObj.Features = results
		fields, err = json.Marshal(featureObj)
		if err != nil {
			log.Println(err)
		}
		return fields
	} else if DbSource == SQLITE3 {
		sql := "select json from services where service=? and name=? and layerid=? and type=?"
		log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
		stmt, err := Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		var fields []byte
		err = stmt.QueryRow(catalog, service, layerid, dtype).Scan(&fields)
		if err != nil {
			log.Println(err.Error())
			//w.Header().Set("Content-Type", "application/json")
			//w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
			return []byte("")
		}
		var featureObj structs.FeatureTable
		err = json.Unmarshal(fields, &featureObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(fields))
			log.Println(err.Error())
		}
		var results []structs.Feature
		for _, i := range featureObj.Features {
			//if int(i.Attributes["OBJECTID"].(float64)) == objectIdsInt {
			if in_float_array(i.Attributes["OBJECTID"].(float64), objectIdsFloat) {
				//oJoinVal = i.Attributes[oJoinKey]
				results = append(results, i)
				//break
			}
		}
		featureObj.Features = results
		fields, err = json.Marshal(featureObj)
		if err != nil {
			log.Println(err)
		}
		return fields
	}
	return []byte("")
	/*
		sql := "select * from  set json=" + GetParam(1) + " where name=" + GetParam(2) + " and type=" + GetParam(3)
		log.Printf("Query: update catalog set json=<json> where name='%v' and type='%v'", service, dtype)

		stmt, err := Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		_, err = stmt.Exec(json, service, dtype)
		if err != nil {
			log.Println(err.Error())
			return false
		}

		return true
	*/
}
func in_string_array(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}
func in_float_array(val float64, array []float64) bool {
	for i := range array {
		if array[i] == val {
			return true
		}
	}
	return false
}

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}
