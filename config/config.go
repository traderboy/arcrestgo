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
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"
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
var DbSource = 0 // SQLITE3

var Schema = "" //= "postgres."
var TableSuffix = ""
var DbTimeStamp = ""

var Project structs.JSONConfig
var RootPath = "catalogs"
var SqlFlags = "?cache=shared&mode=wrc"
var SqlWalFlags = "?PRAGMA journal_mode=WAL"

//leasecompliance2016
var RootName string
var HTTPPort string  // = ":80"
var HTTPSPort string //= ":443"
var Pem = "ssl/reais.x10host.com.key.pem"
var Cert = "ssl/2_reais.x10host.com.crt"
var UUID = ""

//"github.com/gin-gonic/gin"
//Db is the SQLITE databa se object

var configFile = RootPath + string(os.PathSeparator) + "config.json"
var ArcGisVersion = "3.8"

var Db *sql.DB
var DbQuery *sql.DB
var DbSqliteQuery *sql.DB
var DbSqliteDbName string

//var port = ":8080"

var DataPath = RootPath        //+ string(os.PathSeparator)        //+ string(os.PathSeparator) //+ "services"
var ReplicaPath = RootPath     //+ string(os.PathSeparator)     //+ "replicas"
var AttachmentsPath = RootPath //+ string(os.PathSeparator) //+ "attachments"

var CertificatePath = "ssl" + string(os.PathSeparator) + "agent2-cert.cert"

//var config map[string]interface{}
//var defaultService = ""
var UploadPath = ""
var Server = ""
var RefreshToken = "51vzPXXNl7scWXsw7YXvhMp_eyw_iQzifDIN23jNSsQuejcrDtLmf3IN5_bK0P5Z9K9J5dNb2yBbhXqjm9KlGtv5uDjr98fsUAAmNxGqnz3x0tvl355ZiuUUqqArXkBY-o6KaDtlDEncusGVM8wClk0bRr1-HeZJcR7ph9KU9khoX6H-DcFEZ4sRdl9c16exIX5lGIitw_vTmuomlivsGIQDq9thskbuaaTHMtP1m3VVnhuRQbyiZTLySjHDR8OVllSPc2Fpt0M-F5cPl_3nQg.."
var AccessToken = "XMdOaajM4srQWx8nQ77KuOYGO8GupnCoYALvXEnTj0V_ZXmEzhrcboHLb7hGtGxZCYUGFt07HKOTnkNLah8LflMDoWmKGr4No2LBSpoNkhJqc9zPa2gR3vfZp5L3yXigqxYOBVjveiuarUo2z_nqQ401_JL-mCRsXq9NO1DYrLw."

func Initialize() {
	SqlFlags = ""
	SqlWalFlags = ""

	//var err error
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get current directory")
	}
	RootPath = pwd + string(os.PathSeparator) + RootPath //+ string(os.PathSeparator)
	var DbName string

	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			//log.Println(os.Args[i][0] == 45)
			if os.Args[i] == "-sqlite" {
				DbSource = SQLITE3
				if len(os.Args) > i+1 && os.Args[i+1][0] != 45 { //&& len(os.Args[i+1]) > 0 && os.Args[i+1][0] != 45
					DbName = os.Args[i+1]
				} else {
					DbName = pwd + string(os.PathSeparator) + "arcrest.sqlite"
				}

			} else if os.Args[i] == "-pgsql" {
				DbSource = PGSQL
				if len(os.Args) > i+1 && os.Args[i+1][0] != 45 { // && len(os.Args[i+1]) > 0 && os.Args[i+1][0] != 45
					DbName = os.Args[i+1]
				} else {

					DbName = "user=postgres dbname=gis host=192.168.99.100"
				}

			} else if os.Args[i] == "-root" {
				if len(os.Args) > i+1 && os.Args[i+1][0] != 45 {
					RootPath, _ = filepath.Abs(os.Args[i+1])
				} else {
					fmt.Println("No root path entered")
					os.Exit(1)
				}
				//RootName = filepath.Base(os.Args[i+1])
			} else if os.Args[i] == "-p" && len(os.Args) > i+1 {
				HTTPPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-https" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				HTTPSPort = ":" + os.Args[i+1]
			} else if os.Args[i] == "-pem" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				Pem = os.Args[i+1]
			} else if os.Args[i] == "-cert" && len(os.Args) > i && len(os.Args[i+1]) > 0 {
				Cert = os.Args[i+1]
			} else if os.Args[i] == "-file" {
				DbSource = FILE
				LoadConfigurationFromFile()
			} else if os.Args[i] == "-h" {
				fmt.Println("Usage:")
				fmt.Println("go run server.go -p HTTP Port -https HTTPS Port -root <path to service folder> -sqlite <path to service .sqlite> -pgsql <connection string for Postgresql> -h [show help]")
				os.Exit(0)
			}
		}
	} else if len(os.Getenv("DB_SOURCE")) > 0 {
		//read in from environment variables
		RootPath, _ = filepath.Abs(os.Getenv("ROOT_PATH"))

		tmpSrc := os.Getenv("DB_SOURCE")
		if len(tmpSrc) > 0 {
			if tmpSrc == "PG" {
				DbSource = PGSQL
				DbName = os.Getenv("DB_NAME")
			} else if tmpSrc == "SQLITE" {
				DbSource = SQLITE3
				DbName = os.Getenv("DB_NAME")
			} else {
				DbSource = FILE
			}
		} else {
			DbSource = FILE
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
			//Schema = "postgres."
		} else if Project.DataSource == "sqlite" {
			DbSource = SQLITE3
			DbName = Project.SqliteDb
		} else if Project.DataSource == "file" {
			DbSource = FILE
		}
	}
	//for docker, environment variables override command line parameters?
	if len(HTTPPort) == 0 {
		HTTPPort = ":" + os.Getenv("HTTP_PORT") //80
		if len(HTTPPort) == 1 {
			HTTPPort = ":80"
		}
	}
	if len(HTTPSPort) == 0 {
		HTTPSPort = ":" + os.Getenv("HTTPS_PORT") //443
		if len(HTTPSPort) == 1 {
			HTTPSPort = ":443"
		}
	}
	if DbSource == 0 {
		tmpSrc := os.Getenv("DB_SOURCE")
		if len(tmpSrc) > 0 {
			if tmpSrc == "PG" {
				DbSource = PGSQL
			} else if tmpSrc == "SQLITE" {
				DbSource = SQLITE3
			} else {
				DbSource = FILE
			}
		} else {
			DbSource = FILE
		}
	}

	if DbSource == PGSQL {
		Db, err = sql.Open("postgres", DbName)
		if err != nil {
			log.Fatal(err)
		}
		DbQuery = Db
		Schema = "postgres."
		UUID = "('{'||md5(random()::text || clock_timestamp()::text)::uuid||'}')"
		//DbTimeStamp = "(CAST (to_char(now(), 'J') AS INT) - 2440587.5)*86400.0*1000"
		DbTimeStamp = "(now())"

		log.Print("Postgresql database: " + DbName)
		log.Print("Pinging Postgresql: ")
		log.Println(Db.Ping)
		LoadConfiguration()
	} else if DbSource == SQLITE3 {
		Schema = ""
		TableSuffix = "_evw"
		UUID = "(select '{'||upper(substr(u,1,8)||'-'||substr(u,9,4)||'-4'||substr(u,13,3)||'-'||v||substr(u,17,3)||'-'||substr(u,21,12))||'}' from ( select lower(hex(randomblob(16))) as u, substr('89ab',abs(random()) % 4 + 1, 1) as v) as foo)"
		DbTimeStamp = "(julianday('now') - 2440587.5)*86400.0*1000"

		//use 2 different sqlite files:
		//1st: contains configuration information and JSON data
		//2nd: contains actual data
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

		//Db, err = sql.Open("sqlite3", "file:"+DbName+"?PRAGMA journal_mode=WAL")
		sql.Register("sqlite3_with_extensions",
			&sqlite3.SQLiteDriver{
				Extensions: []string{
					"stgeometry_sqlite",
				},
			})

		Db, err = sql.Open("sqlite3", DbName+SqlFlags)
		if err != nil {
			log.Fatal(err)
		}
		err = Db.Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}

		//&sqlite3.SQLiteConn.LoadExtension("stgeometry_sqlite", "sqlite3_stgeometrysqlite_init")

		//sqlite3.LoadExtension("stgeometry_sqlite", "sqlite3_stgeometrysqlite_init")

		/*
		   conn := &SQLiteConn{db: Db, loc: loc, txlock: txlock}
		   conn.LoadExtensions()
		   	if len(d.Extensions) > 0 {
		   		if err := conn.loadExtensions(d.Extensions); err != nil {
		   			conn.Close()
		   			return nil, err
		   		}
		   	}
		*/

		//_, err = DbQuery.Exec("SELECT load_extension('stgeometry_sqlite')")
		//_, err = DbQuery.("stgeometry_sqlite","sqlite3_stgeometrysqlite_init")
		//sqlite3conn := []*sqlite3.SQLiteConn{}
		//c *sqlite3.SQLiteConn
		/*
		   sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{
		   		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
		   			return conn.CreateModule("github", &githubModule{})
		   		},
		   	})
		*/

		//_, err = Db.Exec("SELECT load_extension('stgeometry_sqlite','SDE_SQL_funcs_init')")
		//SELECT load_extension('stgeometry_sqlite.dll','SDE_SQL_funcs_init');
		//if err != nil {
		//	log.Fatalf("Error on loading extension stgeometry_sqlite: %s", err.Error())
		//}

		//Db.Exec(initializeStr)
		log.Println("Sqlite config database: " + DbName)
		//defer Db.Close()
		//Db.SetMaxOpenConns(1)

		LoadConfiguration()
		//get RootName
		DbQueryName := RootPath + string(os.PathSeparator) + RootName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + RootName + ".geodatabase"

		//DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+"?PRAGMA journal_mode=WAL")
		DbQuery, err = sql.Open("sqlite3", DbQueryName+SqlFlags)
		if err != nil {
			log.Fatal(err)
		}
		err = DbQuery.Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}
		log.Println("Sqlite replica database: " + DbQueryName)

		//testQuery()
		//os.Exit(0)

		/*
			sql := "select * from grazing_inspections where GlobalGUID in (select substr(GlobalID, 2, length(GlobalID)-2) from grazing_permittees where OBJECTID in(?))"
			stmt, err := Db.Prepare(sql)
			if err != nil {
				log.Fatal(err)
			}
			rows, err := stmt.Query(16) //relationshipIdInt
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			columns, _ := rows.Columns()
			count := len(columns)
			values := make([]interface{}, count)
			valuePtrs := make([]interface{}, count)
			//for i, _ := range columns {
			//	log.Println(columns[i])
			//}

			for rows.Next() {
				for i, _ := range columns {
					valuePtrs[i] = &values[i]
					log.Println(i)
				}
				rows.Scan(valuePtrs...)
				for i, col := range columns {
					log.Println(i)
					log.Println(col)
				}
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			rows.Close()

			sql = "select * from grazing_inspections where GlobalGUID in (select substr(GlobalID, 2, length(GlobalID)-2) from grazing_permittees where OBJECTID in(16))"
			sql = "select substr(GlobalID, 2, length(GlobalID)-2) as GlobalGUID from grazing_permittees where OBJECTID in(16)"
			sql = "select OBJECTID,cows,yearling_heifers,steer_calves,yearling_steers,bulls,mares,geldings,studs,fillies,colts,ewes,lambs,rams,wethers,kids,billies,nannies,Comments,GlobalGUID,created_user,created_date,last_edited_user,last_edited_date,reviewer_name,reviewer_date,reviewer_title,GlobalID from grazing_inspections"
			sql = "select * from grazing_permittees"
			sql = "select OBJECTID from grazing_inspections"

			log.Println(sql)
			rows, err = Db.Query(sql) //relationshipIdInt
			columns, _ = rows.Columns()
			count = len(columns)
			values = make([]interface{}, count)
			valuePtrs = make([]interface{}, count)

			for rows.Next() {
				for i := range columns {
					valuePtrs[i] = &values[i]
					log.Println(i)
				}
				rows.Scan(valuePtrs...)
				for i, col := range columns {
					log.Println(i)
					log.Println(col)
					//val := values[i]
					//log.Printf("%v", val.([]uint8))
				}
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			rows.Close()

			os.Exit(1)
		*/

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
	AttachmentsPath = RootPath //+ string(os.PathSeparator) + RootName + string(os.PathSeparator) + "attachments" //+ string(os.PathSeparator)
	UploadPath = RootPath + string(os.PathSeparator) + RootName + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "attachments"

	log.Println("Root catalog: " + RootName)
	log.Println("Root path: " + RootPath)
	log.Println("Data path: " + DataPath)
	log.Println("Replica path: " + ReplicaPath)
	log.Println("Attachments path: " + AttachmentsPath)
	var DbSourceName string
	switch DbSource {
	case FILE:
		DbSourceName = "Filesystem"
		break
	case PGSQL:
		DbSourceName = "Postgresql"
		break
	case SQLITE3:
		DbSourceName = "Sqlite"
		break
	default:
		DbSourceName = "Unknown"
	}
	log.Println("Data source: " + DbSourceName)
	log.Printf("HTTP Port: %v\n", HTTPPort)
	log.Printf("HTTPS Port: %v\n", HTTPSPort)

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
		Db, err = sql.Open("sqlite3", "file:"+Project.SqliteDb+SqlWalFlags)
		if err != nil {
			log.Fatal(err)
		}
		DbQueryName := RootPath + string(os.PathSeparator) + RootName + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + RootName + ".geodatabase"
		log.Println("DbQueryName: " + DbQueryName)
		DbQuery, err = sql.Open("sqlite3", "file:"+DbQueryName+SqlWalFlags)

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
func GetArcService(catalog string, service string, layerid int, dtype string, dbPath string) []byte {
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
func GetArcCatalog(service string, dtype string, dbPath string) []byte {
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

func SetArcService(json []byte, catalog string, service string, layerid int, dtype string, dbPath string) bool {
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
func SetArcCatalog(json []byte, service string, dtype string, dbPath string) bool {
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
		sql := "select json from " + Schema + "services where service=$1 and name=$2 and layerid=$3 and type=$4"
		log.Printf("select json from "+Schema+"services where service='%v' and name='%v' and layerid=%v and type='%v'", catalog, service, layerid, dtype)
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
		//var fieldsArr []structs.Field
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
		//globalIdFieldName=GlobalID
		//objectIdFieldName=OBJECTID
		featureObj.GlobalIDField = "GlobalID"
		featureObj.ObjectIDFieldName = "OBJECTID"
		/*
			GlobalIDField     string `json:"globalIdField,omitempty"`
			GlobalIDFieldName string `json:"globalIdFieldName,omitempty"`
			ObjectIDField      string `json:"objectIdField,omitempty"`
			ObjectIDFieldName string    `json:"objectIdFieldName,omitempty"`
		*/

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

func DblQuote(s string) string {
	return "\"" + s + "\""
}

func testQuery() {

	sql := "insert into grazing_inspections(yearling_heifers,studs,lambs,wethers,kids,reviewer_name,reviewer_date,reviewer_title,cows,steer_calves,mares,fillies,nannies,Comments,OBJECTID,colts,ewes,rams,billies,yearling_steers,bulls,geldings,GlobalGUID,GlobalID) values( ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	vals := []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, 13, nil, nil, nil, nil, nil, 20, nil, nil, nil, nil, nil, nil, nil, "{6FC17403-5889-4A23-AC77-3B060E4C6DC4}", "{6FC17403-5889-4A23-AC77-3B060E4C6DC4}"}
	stmt, err := DbQuery.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Println(err.Error())
	}
	stmt.Close()

}
