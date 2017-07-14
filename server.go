package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	"database/sql"

	"github.com/gorilla/handlers"
	//_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"

	config "github.com/traderboy/arcrestgo/config"
	routes "github.com/traderboy/arcrestgo/routes"
)

//_ "github.com/mattn/go-sqlite3"
var logPath = "logfile.txt"

//
//_ "github.com/traderboy/arcrestgo/models"
//"crypto/md5"
//"github.com/gorilla/handlers"
//"github.com/gorilla/mux"
//"time"
//"reflect"
//"net/http"
//"html/template"
//"io"

//var configuration Configuration
//var resultscache={}
func main() {
	//TestDb()
	//return
	//current_time := time.Now().Local()
	//fmt.Printf("%v/%v/%v", current_time.Year(), int(current_time.Month()), current_time.Day())

	if false {
		Db, err1 := sql.Open("postgres", "user=postgres dbname=gis host=192.168.99.100")
		if err1 != nil {
			log.Fatal(err1)
		}

		update := "{\"OBJECTID\":27,\"acres\":3.00,\"lease_site\":0,\"feature_type\":1,\"climatic_zone\":2,\"quad_name\":\"077-SE-196\",\"elevation\":6048,\"permittee\":\"Lorraine / Elsie Begay\",\"homesite_id\":\"H61A\"}"
		//update = "{'OBJECTID':27,'acres':3.00,'lease_site':0,'feature_type':1,'climatic_zone':2,'quad_name':'077-SE-196','elevation':6048,'permittee':'Lorraine / Elsie Begay','homesite_id':'H61A'}'::jsonb"
		//update = "{\"OBJECTID\":27,\"acres\":3.00,\"lease_site\":0,\"feature_type\":1,\"climatic_zone\":2,\"quad_name\":\"077-SE-196\",\"elevation\":6048,\"permittee\":\"Lorraine / Elsie Begay\",\"homesite_id\":\"H61A\"}::jsonb"

		sql := "update services set json=jsonb_set(json,'{features,26,attributes}',$1::jsonb,false) where type='query' and layerId=$2"

		stmt, err2 := Db.Prepare(sql)
		if err2 != nil {
			log.Fatal(err2)
		}
		log.Println(update)

		stmt.Exec(update, 0)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Println("Done")
		return
	}
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//load directory of json files into postgresql
	/*
		if len(os.Args) > 1 {
			loadServices(os.Args[1])
			return
		}
	*/

	logParam := flag.Bool("log", false, "a bool")
	//sqliteParam := flag.Bool("sqlite", false, "a bool")

	if *logParam {
		InitLog()
		log.Println("Writing log file to : logfile.txt")
	} else {
		log.SetOutput(os.Stdout)
		log.Println("Writing log file to stdOut")
	}
	/*
		if *sqliteParam {
			log.Println("SQLite database in memory only")
			InitDb()
		} else {
			log.Println("SQLite database on disk")
			OpenDb()
		}
	*/
	config.Initialize()
	config.Server = ConfigRuntime()
	r := routes.StartGorillaMux()

	//test with: curl -H "Origin: http://localhost" -H "Access-Control-Request-Method: PUT" -H "Access-Control-Request-Headers: X-Requested-With" -X OPTIONS --verbose http://reais.x10host.com/
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	//os.Getenv("ORIGIN_ALLOWED")
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	//http.Handle("/", r)
	//http.HandleFunc("/hello", HelloServer)
	//  Start HTTP
	go func() {
		// Apply the CORS middleware to our top-level router, with the defaults.
		err1 := http.ListenAndServe(config.HTTPPort, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
		if err1 != nil {
			log.Fatal("HTTP server: ", err1)
		}
	}()
	/*
		var pem = "ssl/2_gis.biz.tm.key"
		var cert = "ssl/2_gis.biz.tm.crt"
		cert = "ssl/agent2-cert.cert"
		pem = "ssl/agent2-key.pem"

		cert = "ssl/2_reais.x10host.com.crt"
		pem = "ssl/reais.x10host.com.key.pem"
	*/

	err := http.ListenAndServeTLS(config.HTTPSPort, config.Cert, config.Pem, handlers.CORS(originsOk, headersOk, methodsOk)(r)) //handlers.CORS()(r))
	if err != nil {
		log.Fatal("HTTP server: ", err)
	}

	/*
		err := http.ListenAndServeTLS(":443", "ca.pem", "key.pem", nil)
		if err != nil {
			log.Fatal("HTTP server: ", err)
		}
	*/

	//http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))

}

func TestDb() {
	Db, err := sql.Open("sqlite3", "file:"+"D:/workspace/go/src/github.com/traderboy/arcrestgo/arcrest.sqlite"+"?PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	sql := "select json from services where type='query' and layerId=0"
	stmt, err := Db.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	rows, err := Db.Query(sql)
	defer rows.Close()

	var row []byte
	for rows.Next() {
		err := rows.Scan(&row)
		if err != nil {
			log.Fatal(err)
		}
	}
	rows.Close()
	stmt.Close()

	sql = "delete from services where type='query' and layerId=10"
	_, err = Db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
	}

	sql = "insert into services(type,layerId) values('query',10)"
	_, err = Db.Exec(sql)
	if err != nil {
		log.Println(err.Error())
	}

	//'" + string(row) + "'
	sql = "update services set json=? where type='query' and layerId=10"
	_, err = Db.Exec(sql, string(row))
	if err != nil {
		log.Println(err.Error())
	}

	/*
		tx, err := Db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		_, err = tx.Exec("PRAGMA journal_mode=WAL")
		if err != nil {
			log.Println(err.Error())
		}



		stmt, err = tx.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		_, err = tx.Stmt(stmt).Exec(string(row))
		if err != nil {
			log.Println(err.Error())
		}
		tx.Commit()
	*/
	log.Println("Done")
}

/*
func OpenDb() {
	var err error
	Db, err = sql.Open("sqlite3", "./heroes.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	//defer Db.Close()
	rows, err1 := Db.Query("select count(*) from heroes")
	if err1 != nil {
		//log.Fatal(err)

		//LoadDb()
		return
	}
	defer rows.Close()
}
*/
func InitLog() {

	var err error
	var f *os.File
	f, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//fmt.fprintln("error opening file: %v", err)
		fmt.Printf("%v\n", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

//InitDb intialize databases
/*
func InitDb() {
	var err error
	Db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

}
*/

//ConfigRuntime print out configuration details
func ConfigRuntime() string {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.Printf("Running with %d CPUs\n", nuCPU)
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
		}
	}
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	//server = ip
	fmt.Println("Public IP: " + ip)
	return ip
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

/*
func sendJSON(w http.ResponseWriter, r *http.Request, response gin.H) {
	//setHeaders(c)
	c.JSON(http.StatusOK, response)
}

func sendFile(w http.ResponseWriter, r *http.Request, filename string) {
	//setHeaders(c)
	//c.JSON(http.StatusOK, response)
	//file, err = ioutil.ReadFile(filename)
	//if ext != "" {
	//   w.Header().Set("Content-Type", mime.TypeByExtension(ext))
	//}
	// setHeaders(c)
	//c.String(http.StatusOK,file)
	//setHeaders(c)
	c.File(filename)

}
*/

//conf := gjson.Get(string(file), "services.accommodationagreementrentals.layers.1")
//log.Println(conf)
//for key, val := range config.Services {
//config.Services[key] = val.(map[string]interface{})

//defaultService = key
//}

//defaultService = "break"
/*
	for key, _ := range config.Services {
		defaultService = key
	}
	log.Println("Using default service: " + defaultService)
	log.Println(config.Services["accommodationagreementrentals"])
*/
//log.Println(config.Hostname)
//log.Println(config.Username)
/*
	for key, val := range config.Services {
		log.Println(key)
		for key1, val1 := range val.(map[string]interface{}) {
			log.Println(key1)
			log.Println(val1)
		}
		//log.Println(i)
	}

	//log.Println(config.Services["accommodationagreementrentals"])
	if true {
		os.Exit(0)
	}
*/
/*
	for k, i := range config {
		log.Println(i)
		log.Println(k)
	}
*/
/*
   "services": {
      "accommodationagreementrentals": {
          "layers": {
              "1": {
*/
//log.Println(config["services"]["accommodationagreementrentals"]["layers"]["1"])
//log.Println(config["username"])
//log.Println(config["hostname"])
