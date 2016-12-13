package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//"github.com/gin-gonic/gin"
//Db is the SQLITE databa se object
var Db *sql.DB
var port = ":8080"

var rootPath = "hpl12"
var dataPath = rootPath        //+ string(os.PathSeparator)        //+ string(os.PathSeparator) //+ "services"
var replicaPath = rootPath     //+ string(os.PathSeparator)     //+ "replicas"
var attachmentsPath = rootPath //+ string(os.PathSeparator) //+ "attachments"
var logPath = rootPath + string(os.PathSeparator) + "logfile.txt"
var configFile = rootPath + string(os.PathSeparator) + "config.json"

//var config map[string]interface{}
//var defaultService = ""
var uploadPath = ""
var server = ""
var refreshToken = "51vzPXXNl7scWXsw7YXvhMp_eyw_iQzifDIN23jNSsQuejcrDtLmf3IN5_bK0P5Z9K9J5dNb2yBbhXqjm9KlGtv5uDjr98fsUAAmNxGqnz3x0tvl355ZiuUUqqArXkBY-o6KaDtlDEncusGVM8wClk0bRr1-HeZJcR7ph9KU9khoX6H-DcFEZ4sRdl9c16exIX5lGIitw_vTmuomlivsGIQDq9thskbuaaTHMtP1m3VVnhuRQbyiZTLySjHDR8OVllSPc2Fpt0M-F5cPl_3nQg.."
var accessToken = "XMdOaajM4srQWx8nQ77KuOYGO8GupnCoYALvXEnTj0V_ZXmEzhrcboHLb7hGtGxZCYUGFt07HKOTnkNLah8LflMDoWmKGr4No2LBSpoNkhJqc9zPa2gR3vfZp5L3yXigqxYOBVjveiuarUo2z_nqQ401_JL-mCRsXq9NO1DYrLw."

/*
type Configuration struct {
	Services struct {
		Service struct {
			Layers struct {
				Layer struct {
					ItemID        string `json:"itemId"`
					Data          string `json:"data"`
					Name          string `json:"name"`
					Oidname       string `json:"oidname"`
					Globaloidname string `json:"globaloidname"`
				} `json:"0"`
			} `json:"layers"`
			Relationships struct {
				Relationship struct {
					OID      int    `json:"oId"`
					DID      int    `json:"dId"`
					OTable   string `json:"oTable"`
					OJoinKey string `json:"oJoinKey"`
					DJoinKey string `json:"dJoinKey"`
					DTable   string `json:"dTable"`
				} `json:"0"`
			} `json:"relationships"`
		} `json:"accommodationagreementrentals"`
	} `json:"services"`
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}
*/
type Config struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	//Services
	Services map[string]map[string]map[string]map[string]interface{} `json:"services"`
	//Services map[string]map[string]Service
	//map[string]Service
}

/*
type Service struct {
	//Names map[string]interface{}
	Names map[string]Name
}
type Name struct {
	//Layers map[string]interface{}
	Layers map[string]Layer
}
type Layer struct {
	Items         map[string]Item
	Relationships map[string]Relationship
}
type Item struct {
	ItemID        string `json:"itemId"`
	Data          string `json:"data"`
	Name          string `json:"name"`
	Oidname       string `json:"oidname"`
	Globaloidname string `json:"globaloidname"`
}
type Relationship struct {
	Oid    int    `json:"oId"`
	DId    int    `json:"dId"`
	OTable string `json:"oTable"`

	OJoinKey string `json:"oJoinKey"`
	DJoinKey string `json:"dJoinKey"`
	DTable   string `json:"dTable"`
}
*/

var config Config

//var configuration Configuration

//var resultscache={}
func main() {

	file, err1 := ioutil.ReadFile(configFile)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", configFile)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal(file, &config)
	if err2 != nil {
		log.Println("Error reading configuration file: " + configFile)
		log.Println(err2.Error())
	}

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

	logParam := flag.Bool("log", false, "a bool")
	//sqliteParam := flag.Bool("sqlite", false, "a bool")
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get current directory")
	}
	rootPath = pwd + string(os.PathSeparator) + rootPath //+ string(os.PathSeparator)
	dataPath = rootPath                                  //+ string(os.PathSeparator)        //+ defaultService + string(os.PathSeparator) + "services" + string(os.PathSeparator)
	replicaPath = rootPath                               //+ string(os.PathSeparator)     //+ defaultService + string(os.PathSeparator) + "replicas" + string(os.PathSeparator)
	attachmentsPath = rootPath                           //+ string(os.PathSeparator) //+ defaultService + string(os.PathSeparator) + "attachments" + string(os.PathSeparator)
	log.Println("Data path: " + dataPath)
	log.Println("Replica path: " + replicaPath)
	log.Println("Attachments path: " + attachmentsPath)

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
	ConfigRuntime()
	StartGorillaMux()

}
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
func InitDb() {
	var err error
	Db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

}

//ConfigRuntime print out configuration details
func ConfigRuntime() {
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
	server = ip
	fmt.Println("Public IP: " + ip)
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
func StartGorillaMux() {

	r := mux.NewRouter()

	/*
	   Download certs
	*/
	r.HandleFunc("/cert", func(w http.ResponseWriter, r *http.Request) {
		//res.sendFile("certs/server.crt", { root : __dirname})
		log.Println("Sending: " + rootPath + "certs/server.crt")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"certs"+string(os.PathSeparator)+"server.crt")
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome"))
	}).Methods("GET")
	/*
	   Root responses
	*/
	r.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": "3.8"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		//w.Write(response)
		//setHeaders(c)
		//fmt.Println(response)
		//w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest (post)")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": "3.8"})
		//w.Write(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": "3.8"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		//w.Write(response)
	}).Methods("GET")

	/*
	   authentication.  Uses phoney tokens
	*/
	/*
		type esritoken struct {
			Token   string `json:"token:`
			Expires int64  `json:"expires"`
			SSL     bool   `json:"ssl"`
		}
	*/
	r.HandleFunc("/sharing/{rest}/generateToken", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rest := vars["rest"]

		log.Println("/sharing/" + rest + "/generateToken")
		/*
			tok := esritoken{
				Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
				Expires: 1940173783033,
				SSL:     false,
			}

			response, _ := json.Marshal(tok)
		*/
		var expires int64 = 1440173783033
		response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST")

	r.HandleFunc("/sharing/generateToken", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/generateToken (post)")
		//response, _ := json.Marshal([]string{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": 1940173783033, "ssl": false}
		/*
			tok := esritoken{
				Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
				Expires: 1940173783033,
				SSL:     false,
			}

			response, _ := json.Marshal(tok)
		*/
		var expires int64 = 1440173783033
		response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	}).Methods("GET", "POST")

	/*
	   r.Methods("POST").HandleFunc("/sharing/{rest}/generateToken", func(w http.ResponseWriter, r *http.Request){
	   	//log.Println("Logging in post");
	   	//log.Println(req.query);
	   	//log.Println(req.body);
	   	log.Println("Post rest/generateToken");
	   	response, _ := json.Marshal([]string{"token":"hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp","expires":1440173783033,"ssl":false}
	   	//response, _ := json.Marshal([]string{"token":"NrCcZaQedpZJHxaqSvtwBS1ycDOd3XiDL46C-UsRzZummvdCNQrFzDh1roNmZLToDL27gEu8-E1Mx2p4_GG5qSJ4ISyL06Npizxtv0bzkkfGEwrGBQJ4q1W8kybo3H1_","expires":1940038262530,"ssl":false}
	   	w.Write(response)

	   })
	*/

	r.HandleFunc("//sharing/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Println("//sharing/oauth2/authorize")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "oauth2.html")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"oauth2.html")
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/authorize")
		http.Redirect(w, r, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U", http.StatusMovedPermanently)
		//302
		//c.Redirect(http.StatusMovedPermanently, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U")
		//http.ServeFile(w, r, rootPath + "/oauth2.html");
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/approval", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/approval")
		/*
			tok := esritoken{
				Token:   "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp",
				Expires: 1940173783033,
				SSL:     false,
			}
			response, _ := json.Marshal(tok)
		*/

		var expires int64 = 1440173783033
		response, _ := json.Marshal(map[string]interface{}{"token": "hbKgXcKhu_t6oTuiMOxycWn_ELCZ5G5OEwMPkBzbiCrrQdClpi531MbGo0P_HsukvhoIP8uzecIwpD8zoCaZoy1POpEUDwtNXLf-K6n913cayKDVD6wsePmgzYSNoogp", "expires": expires, "ssl": false})

		w.Header().Set("Content-Type", "application/json")

		w.Write(response)
		//w.Write( response)
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/signin", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/signin")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "search.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"search.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/token")

		var expires int64 = 99800
		response, _ := json.Marshal(map[string]interface{}{"access_token": accessToken, "expires_in": expires, "username": "gisuser", "refresh_token": refreshToken})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST")

	r.HandleFunc("/sharing/rest/tokens", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/tokens")
		response, _ := json.Marshal(map[string]interface{}{"token": "1.0"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST")
	/*
	   openssl req -x509 -nodes -days 365 -newkey rsa:1024 \
	       -keyout /etc/ssl/private/reais.key \
	       -out /etc/ssl/certs/reais.crt
	*/

	/*
	   End authentication
	*/

	r.HandleFunc("/sharing/{rest}/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, rootPath + "/search.json")
		log.Println("/sharing/{rest}/accounts/self")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "portals.self.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"portals.self.json")

	}).Methods("GET")

	r.HandleFunc("/sharing//accounts/self", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "account.self.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"account.self.json")
	}).Methods("GET")

	//no customization necesssary except for username
	r.HandleFunc("/sharing/rest/portals/self", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/portals/self")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "portals.self.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"portals.self.json")
		//http.ServeFile(w, r, rootPath + "/portals_self.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/users/{user}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		user := vars["user"]

		log.Println("/sharing/rest/content/users/" + user)
		//response, _ := json.Marshal([]string{ "username"{user}"),"total":0,"start":1,"num":0,"nextStart":-1,"currentFolder":nil,"items":[],"folders":[] }
		//folders := make([]int64], 0)
		//folders := make([]string], 0)
		response, _ := json.Marshal(map[string]interface{}{"folders": []string{}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("/sharing/rest/content/items/" + id)
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")

		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+id+"services"+string(os.PathSeparator)+"content.items.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/data", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("/sharing/rest/content/items/" + id + "/data")
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.data.json")

		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+id+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"content.items.data.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/search")
		//vars := mux.Vars(r)

		//q := vars["q"]
		//q := r.Queries("q")
		q := r.FormValue("q")
		if strings.Index(q, "typekeywords") == -1 {
			log.Println("Sending: " + rootPath + string(os.PathSeparator) + "community.groups.json")
			http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"community.groups.json")
		} else {
			log.Println("Sending: " + rootPath + string(os.PathSeparator) + "search.json")
			http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"search.json")
		}
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/users/{user}/notifications", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing/rest/community/users/" + user + "/notifications")
		response, _ := json.Marshal(map[string]interface{}{"notifications": []string{}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/groups", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/community/groups")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "community.groups.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"community.groups.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing/rest/community/users/" + user)
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "community.users.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"community.users.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/users", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/community/users/")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "community.users.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"community.users.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/users/{user}/info/{img}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		img := vars["img"]
		log.Println("/sharing/rest/community/users/" + user + "/info/" + img)

		var path = "photos/cat.jpg"
		log.Println("Sending: " + path)
		http.ServeFile(w, r, path)
		//var fs = require("fs')
		//var file = fs.readFileSync(path, "utf8")
		//res.end(file)

	}).Methods("GET")
	r.HandleFunc("/sharing/rest/community/groups", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/community/groups")
		log.Println("Sending: " + rootPath + string(os.PathSeparator) + "community.groups.json")
		http.ServeFile(w, r, rootPath+string(os.PathSeparator)+"community.groups.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/{img}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		img := vars["img"]
		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/" + img)
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "thumbnails" + string(os.PathSeparator) + id + ".png")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+id+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"thumbnails"+string(os.PathSeparator)+id+".png")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/ago_downloaded.png", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/ago_downloaded.png")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": "3.8"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/info", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/info")
		response, _ := json.Marshal(map[string]interface{}{"owningSystemUrl": "http://" + server,
			"authInfo": map[string]interface{}{"tokenServicesUrl": "https://" + server + "/sharing/rest/generateToken", "isTokenBasedSecurity": true}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/info", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/info")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": "10.3", "fullVersion": "10.3", "authInfo": map[string]interface{}{"isTokenBasedSecurity": false}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	/*
	   Database functions
	*/
	/*
	   POST http://services5.arcgis.com/xxxx/ArcGIS/rest/services/xxxxx/FeatureServer/unRegisterReplica HTTP/1.1
	   Content-Length: 59
	   Content-Type: application/x-www-form-urlencoded
	   Host: services5.arcgis.com
	   Connection: Keep-Alive
	   User-Agent: Collector-Android-10.3.3/ArcGIS.Android-10.2.5/4.4.4/BARNES-&-NOBLE-BN-NOOKHD+
	   Accept-Encoding: gzip

	   f=json&replicaID=
	*/
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs/replicas", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/job/replica")
		var submissionTime int64 = 1441201696150
		var lastUpdatedTime int64 = 1441201705967
		response, _ := json.Marshal(map[string]interface{}{
			"replicaName": "MyReplica", "replicaID": "58808194-921a-4f9f-ac97-5ffd403368a9", "submissionTime": submissionTime, "lastUpdatedTime": lastUpdatedTime,
			"status": "Completed", "resultUrl": "http://" + server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/unRegisterReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/unRegisterReplica")
		response, _ := json.Marshal(map[string]interface{}{"success": true})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/replicas", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/replicas")
		var fileName = replicaPath + "/" + name + ".geodatabase"
		log.Println("Sending: " + fileName)
		http.ServeFile(w, r, fileName) //, { root : __dirname})
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/createReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/createReplica (post)")
		response, _ := json.Marshal(map[string]interface{}{"statusUrl": "http://" + server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/synchronizeReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/synchronizeReplica")
		response, _ := json.Marshal(map[string]interface{}{"status": "Completed", "transportType": "esriTransportTypeUrl"})
		/*
			  "responseType": <esriReplicaResponseTypeEdits | esriReplicaResponseTypeEditsAndData| esriReplicaResponseTypeNoEdits>,
			  "resultUrl": "<url>", //path to JSON (dataFormat=JSON) or a SQLite geodatabase (dataFormat=sqlite)
			  "submissionTime": "<T1>",  //Time since epoch in milliseconds
			  "lastUpdatedTime": "<T2>", //Time since epoch in milliseconds
			  "status": "<Pending | InProgress | Completed | Failed | ImportChanges | ExportChanges | ExportingData | ExportingSnapshot
				       | ExportAttachments | ImportAttachments | ProvisioningReplica | UnRegisteringReplica | CompletedWithErrors>"
		*/
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	}).Methods("POST")
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/jobs")
		var submissionTime int64 = 1441201696150
		var lastUpdatedTime int64 = 1441201705967
		response, _ := json.Marshal(map[string]interface{}{"replicaName": "MyReplica", "replicaID": "58808194-921a-4f9f-ac97-5ffd403368a9", "submissionTime": submissionTime,
			"lastUpdatedTime": lastUpdatedTime, "status": "Completed", "resultUrl": "http://" + server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/services")
		log.Println("Sending: " + dataPath + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+"FeatureServer.json")
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/services (post)")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Object moved</title></head><body>" +
			"<h2>Object moved to <a href=\"/arcgis/rest/services\">here</a>.</h2>" +
			"</body></html>"))
	}).Methods("POST")

	r.HandleFunc("/arcgis/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/services")
		log.Println("Sending: " + dataPath + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+"FeatureServer.json")
	}).Methods("GET")

	r.HandleFunc("/llarcgis/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/services (post)")
		log.Println("Sending: " + dataPath + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+"FeatureServer.json")
	}).Methods("POST")

	r.HandleFunc("/arcgis/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/services (post)")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Object moved</title></head><body>" +
			"<h2>Object moved to <a href=\"/arcgis/rest/services\">here</a>.</h2>" +
			"</body></html>"))
	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services//services/{name}/FeatureServer/info/metadata", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/info/metadata")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metadata stuff"))
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services//{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name)
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + name + ".json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+name+".json")
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name)
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")

		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
	}).Methods("GET")

	r.HandleFunc("/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/rest/services/" + name + "/FeatureServer")
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
	}).Methods("GET", "POST")

	r.HandleFunc("/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/rest/services/" + name + "/FeatureServer")
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer")
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
	}).Methods("GET", "POST")
	/*
		r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			name := vars["name"]

			log.Println("/arcgis/rest/services/" + name + "/FeatureServer  (post)")
			log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
		}).Methods("POST")
	*/

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id)
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "  (post)")
		log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
		http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", func(w http.ResponseWriter, r *http.Request) {
		//if(req.query.outFields=='OBJECTID'){
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]

		if r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectid")
			log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectid.json")
			http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectid.json")
		} else {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")
			log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
			http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
		}
		//http.ServeFile(w, r, dataPath + "/" + id  + "query.json")

	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]

		if r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectid (post)")
			log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectid.json")
			http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectid.json")
		} else {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query (post)")
			log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
			http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/queryRelatedRecords", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		//id := vars["id"]

		var relationshipId = r.FormValue("relationshipId")
		var objectIds, _ = strconv.Atoi(r.FormValue("objectIds"))
		var outFields = r.FormValue("outFields")

		var replicaDb = rootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		var tableName = config.Services[name]["relationships"][relationshipId]["dTable"].(string)
		log.Println(tableName)
		var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))

		var jsonFile = rootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
			strconv.Itoa(layerId) + ".json"
		file, err1 := ioutil.ReadFile(jsonFile)
		if err1 != nil {
			fmt.Printf("// error while reading file %s\n", jsonFile)
			fmt.Printf("File error: %v\n", err1)
			os.Exit(1)
		}
		var fields TableField
		err2 := json.Unmarshal(file, &fields)
		if err2 != nil {
			log.Println("Error reading configuration file: " + jsonFile)
			log.Println(err2.Error())
		}

		//var jsonFields=JSON.parse(file)
		log.Println("sqlite: " + replicaDb)
		//var db = new sqlite3.Database(replicaDb)
		var sqlstr = "select " + outFields + " from " +
			config.Services[name]["relationships"][relationshipId]["dTable"].(string) +
			" where " +
			config.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
			config.Services[name]["relationships"][relationshipId]["oTable"].(string) +
			" where OBJECTID=?)"

		db, err := sql.Open("sqlite3", replicaDb)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		stmt, err := db.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		//outArr := []interface{}{}
		rows, err := stmt.Query(relationshipId)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		//final_result := map[int]map[string]string{}
		//works final_result := map[int]map[string]interface{}{}
		final_result := make([]interface{}, 0)
		result_id := 0
		for rows.Next() {
			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)

			//tmp_struct := map[string]string{}
			tmp_struct := map[string]interface{}{}

			for i, col := range columns {
				var v interface{}
				val := values[i]
				switch t := val.(type) {
				case int:
					fmt.Printf("Integer: %v\n", t)
					tmp_struct["col"] = val
				case float64:
					//fmt.Printf("Float64: %v\n", t)
					tmp_struct[col] = val
				case []uint8:
					//fmt.Printf("String: %v\n", t)
					b, _ := val.([]byte)
					tmp_struct[col] = fmt.Sprintf("%s", b)
				case int64:
					//fmt.Printf("Integer 64: %v\n", t)
					tmp_struct[col] = val
				case string:
					fmt.Printf("String: %v\n", t)
					tmp_struct[col] = fmt.Sprintf("%s", v)
				case bool:
					fmt.Printf("Bool: %v\n", t)
					tmp_struct[col] = val
				case []interface{}:
					for i, n := range t {
						fmt.Printf("Item: %v= %v\n", i, n)
					}
				default:
					var r = reflect.TypeOf(t)
					fmt.Printf("Other:%v\n", r)
				}
				/*
					b, ok := val.([]byte)
					if ok {
						v = string(b)
					} else {
						v = val
					}
				*/
				//tmp_struct[col] = fmt.Sprintf("%s", v)
			}
			record := map[string]interface{}{"attributes": tmp_struct}
			//final_result["attributes"] = tmp_struct
			//final_result[result_id] = record
			final_result = append(final_result, record)
			result_id++
		}

		//stmt.Close()
		//rows.Close() //good habit to close
		//db.Close()

		//fmt.Println(final_result)
		/*
			var results=[]

			err = stmt.QueryRow(relationshipId).Scan(&)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(name)
		*/
		//var fields = jsonFields.fields
		//console.log("Sql: " + sql + " [" + objectIds + "]")
		//db.run(sql, [objectIds]);
		/*/


		  db.each(sql,[objectIds], function(err, row) {
		      var obj=map[string]interface{}
		      for( var i in row){
		        //console.log(i + ": " + row[i]);
		        obj[i]=row[i];
		      }
		      results.push({"attributes":obj})
		      //results.push({"objectId":objectIds,"relatedRecords":[{"attributes":obj}]})

		  },function(err,rows){
		  	  console.log("Number of rows: " + rows)

		      var result={"fields":fields.Fields,"relatedRecordGroups":[{"objectId":parseInt(objectIds),"relatedRecords":results}]}
		      //console.log(JSON.stringify(result))
		      res.json(result)
		  	});
		*/
		//var results []string
		var result = map[string]interface{}{}
		result["objectId"] = objectIds //strconv.Atoi(objectIds)
		result["relatedRecords"] = final_result
		//{"objectId":strconv.Atoi(objectIds),"relatedRecords":results}
		//resultArr := make([]interface{}, 1)
		//resultArr[0] = result
		//resultArr :=

		response, _ := json.Marshal(map[string]interface{}{"fields": fields.Fields, "relatedRecordGroups": []interface{}{result}})
		//var response []byte
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		/*
			if r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
				log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectid (post)")
				log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectid.json")
				http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectid.json")
			} else {
				log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query (post)")
				log.Println("Sending: " + dataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, dataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
			}
		*/
	}).Methods("GET", "POST")

	/*
	   Attachments
	*/
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/attachments", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		row := vars["row"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/attachments")
		//{"attachmentInfos":[{"id":5,"globalId":"xxxx","parentID":"47","name":"cat.jpg","contentType":"image/jpeg","size":5091}]}
		var attachmentPath = attachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		//attachments:=[]interface{}
		attachments := make([]interface{}, 0)
		//[]interface{}
		//fields.Fields, "relatedRecordGroups": []interface{}{result}}
		files, _ := ioutil.ReadDir(attachmentPath)
		i := 0
		for _, f := range files {
			attachfile := map[string]interface{}{"id": i, "contentType": "image/jpeg", "name": f.Name()}
			attachments = append(attachments, attachfile)
			i++
		}
		response, _ := json.Marshal(map[string]interface{}{"attachmentInfos": attachments})
		//var response []byte
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		//var response={"attachmentInfos":infos}

		//w.Write([]byte(attachmentPath))
		/*
				if(fs.existsSync(attachmentPath)){
				   var files = fs.readdirSync(attachmentPath)
				   var infos=[]
				   for(var i in files)
				     infos.push({"id":i,"contentType":"image/jpeg","name":files[i]})
				     //{"id"{row}"),"contentType":"image/jpeg","name"{row}")+".jpg"}
				   response, _ := json.Marshal([]string{"attachmentInfos":infos}
				}
				else
				  response, _ := json.Marshal([]string{"attachmentInfos":[]}
			  w.Write(response)
		*/
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/attachments/{img}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		img := vars["img"]
		row := vars["row"]
		log.Println("/arcgis/rest/services/FeatureServer/attachments/img")
		var attachment = attachmentsPath + string(os.PathSeparator) + name + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator) + img + ".jpg"

		/*
			files, _ := ioutil.ReadDir("./")
			for _, f := range files {
				fmt.Println(f.Name())
			}
		*/
		if _, err := os.Stat(attachment); err == nil {
			http.ServeFile(w, r, attachment)
		} else {
			response, _ := json.Marshal(map[string]interface{}{"status": "Completed", "transportType": "esriTransportTypeUrl"})
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}

		/*
				if(fs.existsSync(attachment))
			    res.sendFile(attachment)
			  else
			  	res.sendJSON({"Error":"File not found"})
		*/
		/*
				var path="photos/cat.jpg"
				var fs = require("fs')
			  var file = fs.readFileSync(path, "utf8")
			  res.end(file)
		*/
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/addAttachment", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		row := vars["row"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/addAttachment")
		// TODO: move and rename the file using req.files.path & .name)
		//res.send(console.dir(req.files))  // DEBUG: display available fields
		var uploadPath = attachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		//w.Write([]byte(uploadPath))
		if r.Method == "GET" {
			crutime := time.Now().Unix()
			h := md5.New()
			io.WriteString(h, strconv.FormatInt(crutime, 10))
			token := fmt.Sprintf("%x", h.Sum(nil))

			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, token)
		} else {
			r.ParseMultipartForm(32 << 20)
			file, handler, err := r.FormFile("uploadfile")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Fprintf(w, "%v", handler.Header)
			f, err := os.OpenFile(uploadPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}
		response, _ := json.Marshal(map[string]interface{}{"addAttachmentResult": map[string]interface{}{"objectId": id, "globalId": nil, "success": true}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		/*
			  var mkdirp = require("mkdirp')
			  if(!fs.existsSync(uploadPath)){
			      //fs.mkdir(uploadPath,function(e){
			      mkdirp.sync(uploadPath,function(e){
			          if(!e || (e && e.code === 'EEXIST')){
			              //do something with contents

			          } else {
			              //debug
			              log.Println(e)
			          }
			      })
			  }
			  var files = fs.readdirSync(uploadPath)
			  var id=files.length
			  var fstream
			  req.pipe(req.busboy)
			  req.busboy.on("file", function (fieldname, file, filename) {
			        log.Println("Uploading: " + filename)
			        var attachment = uploadPath + "/" + id + ".jpg"
			        fstream = fs.createWriteStream(attachmentPath)
			        file.pipe(fstream)
			        fstream.on("close", function () {
			            //res.redirect("back')
				          response, _ := json.Marshal([]string{"addAttachmentResult":{"objectId"{id},"globalId":null,"success":true}}
			            w.Write(response)
			        })
			  })
		*/
		/*
			  fs.readFile(req.files.attachment.path, function (err, data) {
			    // ...
			    var newPath = uploadPath + "/" + row") + ".jpg";
			    fs.writeFile(newPath, data, function (err) {
			      //res.redirect("back");
				    response, _ := json.Marshal([]string{"addAttachmentResult":{"objectId"{row}"),"globalId":null,"success":true}}
			      w.Write(response)
			    })
			  })
		*/
	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/updateAttachment", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		row := vars["row"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/updateAttachment")
		var aid = vars["attachmentIds"]
		results := []string{aid}

		//results[0] = gin.H{"objectId": id, "globalId": nil, "success": "true"}

		response, _ := json.Marshal(map[string]interface{}{"deleteAttachmentResults": results})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/{row}/deleteAttachments", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		row := vars["row"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/deleteAttachments")
		var aid = r.FormValue("attachmentIds")

		//results := []string{"objectId": id, "globalId": nil, "success": true}
		//results := []string{aid}
		response, _ := json.Marshal(map[string]interface{}{"deleteAttachmentResults": aid})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/applyEdits", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/applyEdits")

		var replicaDb = rootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		var tableName = config.Services[name]["relationships"][relationshipId]["dTable"].(string)
		log.Println(tableName)
		var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))

		var jsonFile = rootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
			strconv.Itoa(layerId) + ".json"
		file, err1 := ioutil.ReadFile(jsonFile)
		if err1 != nil {
			fmt.Printf("// error while reading file %s\n", jsonFile)
			fmt.Printf("File error: %v\n", err1)
			os.Exit(1)
		}
		var fields TableField
		err2 := json.Unmarshal(file, &fields)
		if err2 != nil {
			log.Println("Error reading configuration file: " + jsonFile)
			log.Println(err2.Error())
		}
		//now read posted JSON
		var updates := map[]interface{}{}
		err2 := json.Unmarshal(r.FormValue("updates"), &updates)
		if err2 != nil {
			log.Println("Error reading configuration file: " + r.FormValue("updates"))
			log.Println(err2.Error())
		}
/*
decoder := json.NewDecoder(r.Body)
    var t test_struct   
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()
*/

		//var jsonFields=JSON.parse(file)
		log.Println("sqlite: " + replicaDb)
		//var db = new sqlite3.Database(replicaDb)
		var sqlstr = "update " + outFields + " from " +
			config.Services[name]["relationships"][relationshipId]["dTable"].(string) +
			" where " +
			config.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
			config.Services[name]["relationships"][relationshipId]["oTable"].(string) +
			" where OBJECTID=?)"

		db, err := sql.Open("sqlite3", replicaDb)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		stmt, err := db.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		//outArr := []interface{}{}
		rows, err := stmt.Query(relationshipId)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		//final_result := map[int]map[string]string{}
		//works final_result := map[int]map[string]interface{}{}
		final_result := make([]interface{}, 0)
		result_id := 0


		//var updates = JSON.parse(req.body.updates)//JSON.parse(req.query.updates)
		/*
				var fs = require("fs')
				var path=dataPath+"/"+name") +"/FeatureServer."+id") + ".query.json"
			  var file = fs.readFileSync(path, "utf8")
			  var json=JSON.parse(file)
			  var results=[]
			  var fields=[]
			  var values=[]

			  for(var u=0;u<updates.length;u++)
			  {
				  for(var i=0;i<json.features.length;i++)
				  {
				  	//log.Println(json.features[i]['attributes']['OBJECTID'] + ":  " + updates[u].attributes['OBJECTID'])
				  	if(json.features[i]['attributes']['OBJECTID']==updates[u].attributes['OBJECTID'])
				  	{
				  		//json.features.[i]['attributes']=updates
				  		for(var j in updates[u].attributes)
				  		{
				  			for(var k in json.features[i]['attributes'])
				  			{
				  				if(j==k)
				  				{
				  					if(json.features[i]['attributes'][k] != updates[u].attributes[j])
				  					{
				  					    log.Println("Updating record: " + updates[u].attributes['OBJECTID'] + " " + k + "   values: " + json.features[i]['attributes'][k]+ " to " + updates[u].attributes[j] )
				  					    json.features[i]['attributes'][k]=updates[u].attributes[j]
			  	              fields.push(k+"=?")
			  	              values.push(updates[0].attributes[j])
				  					    break
				  				  }
				  				}
				  			}
				  		}
				  		results.push({"objectId":updates[u].attributes['OBJECTID'],"globalId":null,"success":true})
				  		break
				  	}
				  }
			  }
			  if(fields.length>0){
				  //search for id and update all fields
				  fs.writeFileSync(path, JSON.stringify(json), "utf8")

				  //now update the replica database

				  values.push(parseInt(id")))

				  var replicaDb = replicaPath + "/"+name")+".geodatabase"
				  log.Println("sqlite: " + replicaDb)
				  var db = new sqlite3.Database(replicaDb)
				  //create update statement from json
				  log.Println("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?")
				  log.Println( values )

				  db.run("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?", values)
			  }else{
			 	  results={"objectId":updates.length>0?updates[0].attributes['OBJECTID']:0,"globalId":null,"success":true}
			 	}
		*/
		//update json file with updates
		results := []string{}
		response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	var pem = "ssl/2_gis.biz.tm.key"
	var cert = "ssl/2_gis.biz.tm.crt"
	cert = "ssl/agent2-cert.cert"
	pem = "ssl/agent2-key.pem"
	//http.Handle("/", r)
	//http.HandleFunc("/hello", HelloServer)
	//  Start HTTP
	go func() {

		// Apply the CORS middleware to our top-level router, with the defaults.
		err1 := http.ListenAndServe(":8080", handlers.CORS()(r))
		if err1 != nil {
			log.Fatal("HTTP server: ", err1)
		}
	}()

	err := http.ListenAndServeTLS(":446", cert, pem, handlers.CORS()(r))
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

type FieldsStr struct {
	Fields json.RawMessage `json:"fields"`
	//Fields []Field `json:"fields"`
}
type TableField struct {
	//Fields json.RawMessage `json:"fields"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Domain       *Domain     `json:"domain"`
	Name         string      `json:"name"`
	Nullable     bool        `json:"nullable"`
	DefaultValue interface{} `json:"defaultValue"`
	Editable     bool        `json:"editable"`
	Alias        string      `json:"alias"`
	SqlType      string      `json:"sqlType"`
	Type         string      `json:"type"`
	Length       int         `json:"length,omitempty"`
}

type Domain struct {
	CodedValues []struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"codedValues,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
