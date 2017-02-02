package routes

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	config "github.com/traderboy/arcrestgo/config"
	structs "github.com/traderboy/arcrestgo/structs"
)

//_ "github.com/mattn/go-sqlite3"
func StartGorillaMux() *mux.Router {

	r := mux.NewRouter()

	//fs := http.FileServer(http.Dir("."))
	//http.Handle("/", fs)

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/config (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "config", "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("config", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "config.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"config.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/db (" + r.Method + ")")
		//vars := mux.Vars(r)
		var id int
		idstr := r.URL.Query().Get("id")

		if len(idstr) > 0 {
			id, _ = strconv.Atoi(idstr)
		} else {
			id = config.DbSource
		}

		str := ""
		//	PGSQL   = 1
		//	SQLITE3 = 2
		//	FILE    = 3

		if id == 3 {
			str += "<li>Static JSON files <b style='color:red'>active </b></li>"
			config.SetDatasource(config.FILE)
		} else {
			str += "<li>Static JSON files <a href='/db?id=3'>enable</a> </li>"
		}
		if id == 2 {
			str += "<li>Sqlite <b style='color:red'>active </b> </li>"
			config.SetDatasource(config.SQLITE3)
		} else {
			str += "<li>Sqlite <a href='/db?id=2'>enable</a> </li>"
		}
		if id == 1 {
			str += "<li>Postgresql <b style='color:red'>active </b> </li>"
			config.SetDatasource(config.PGSQL)
		} else {
			str += "<li>Postgresql <a href='/db?id=1'>enable</a> </li>"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Current data source</h1><ul>" + str + "</ul>"))

	}).Methods("GET")

	//r.StrictSlash(true)
	/*
	   Download certs
	*/
	r.HandleFunc("/cert/", func(w http.ResponseWriter, r *http.Request) {
		//res.sendFile("certs/server.crt", { root : __dirname})
		log.Println("Sending: " + config.RootPath + "certs/server.crt")
		http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"certs"+string(os.PathSeparator)+"server.crt")
	}).Methods("GET")

	/*
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("/")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Welcome"))
		}).Methods("GET", "OPTIONS")
	*/
	r.HandleFunc("/arcgis", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome"))
	}).Methods("GET")

	/*
	   Root responses
	*/
	r.HandleFunc("/sharing", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		//w.Write(response)
		//setHeaders(c)
		//fmt.Println(response)
		//w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		//w.Write(response)
		//setHeaders(c)
		//fmt.Println(response)
		//w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {

		log.Println("/sharing/rest (" + r.Method + ")")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		//w.Write(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST")

	r.HandleFunc("/sharing/rest/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest (" + r.Method + ")")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		//w.Write(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST", "PUT")

	/*
		r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {
			log.Println("/sharing/rest")
			response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
			//w.Write(response)
		}).Methods("GET")
	*/
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
		log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "oauth2.html")
		http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"oauth2.html")
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/authorize")
		http.Redirect(w, r, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U", http.StatusMovedPermanently)
		//302
		//c.Redirect(http.StatusMovedPermanently, "/sharing/rest?f=json&culture=en-US&code=KIV31WkDhY6XIWXmWAc6U")
		//http.ServeFile(w, r, config.RootPath + "/oauth2.html");
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
		log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "search.json")
		http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"search.json")
	}).Methods("GET")

	r.HandleFunc("/sharing/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/oauth2/token")

		var expires int64 = 99800
		response, _ := json.Marshal(map[string]interface{}{"access_token": config.AccessToken, "expires_in": expires, "username": "gisuser", "refresh_token": config.RefreshToken})
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
		//http.ServeFile(w, r, config.RootPath + "/search.json")
		log.Println("/sharing/{rest}/accounts/self (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "portals", "self")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}
		response := config.GetArcCatalog("portals", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "portals.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"portals.self.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing//accounts/self (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "account", "self")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("account", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "account.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"account.self.json")
		}
	}).Methods("GET", "PUT")

	//no customization necesssary except for username
	r.HandleFunc("/sharing/rest/portals/self", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/portals/self (" + r.Method + ")")
		if r.Method == "PUT" {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "portals", "self")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("portals", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "portals.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"portals.self.json")
		}
		//http.ServeFile(w, r, config.RootPath + "/portals_self.json")
	}).Methods("GET", "PUT")

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

	r.HandleFunc("/sharing/rest/content/items/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		//temp
		name = config.RootName
		log.Println("/sharing/rest/content/items/" + name + "(" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "content", -1, "items")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		//load from db
		response := config.GetArcService(name, "content", -1, "items")
		if len(response) > 0 {
			//log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+"services"+string(os.PathSeparator)+"content.items.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/rest/content/items/{name}/data", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		if config.DbSource != config.FILE {
			name = "%"

		}
		log.Println("Old name:  " + name)
		name = config.RootName
		log.Println("New name:  " + name)
		log.Println("/sharing/rest/content/items/" + name + "/data (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "content", -1, "data")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcService(name, "content", -1, "data")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.data.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"content.items.data.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/rest/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/search (" + r.Method + ")")
		//vars := mux.Vars(r)

		//q := vars["q"]
		//q := r.Queries("q")
		q := r.FormValue("q")
		if strings.Index(q, "typekeywords") == -1 {
			if r.Method == "PUT" {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.Write([]byte("Error"))
					return
				}
				ret := config.SetArcCatalog(body, "community", "groups")
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": ret})
				w.Write(response)
				return
			}

			response := config.GetArcCatalog("community", "groups")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)

			} else {
				log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.groups.json")
				http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.groups.json")
			}
		} else {
			if r.Method == "PUT" {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.Write([]byte("Error"))
					return
				}
				ret := config.SetArcCatalog(body, "search", "")
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": ret})
				w.Write(response)
				return
			}

			response := config.GetArcCatalog("search", "")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)

			} else {
				log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "search.json")
				http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"search.json")
			}
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/rest/community/users/{user}/notifications", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing/rest/community/users/" + user + "/notifications")
		response, _ := json.Marshal(map[string]interface{}{"notifications": []string{}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")
	/*
		r.HandleFunc("/sharing/rest/community/groups", func(w http.ResponseWriter, r *http.Request) {
			log.Println("/sharing/rest/community/groups")
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.groups.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.groups.json")
		}).Methods("GET")
	*/

	r.HandleFunc("/sharing//community/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing//community/users/" + user + "(" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "community", "users")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("community", "users")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {

			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.users.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.users.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/rest/community/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing/rest/community/users/" + user + "(" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "community", "users")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("community", "users")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {

			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.users.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.users.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/sharing/rest/community/users", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/community/users/ (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcCatalog(body, "community", "users")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("community", "users")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.users.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.users.json")
		}
	}).Methods("GET", "PUT")

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

		response := config.GetArcCatalog("community", "groups")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.groups.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.groups.json")
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/{img}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		img := vars["img"]
		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/" + img)
		log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "thumbnails" + string(os.PathSeparator) + id + ".png")
		http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+id+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"thumbnails"+string(os.PathSeparator)+id+".png")
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/info/thumbnail/ago_downloaded.png", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		log.Println("/sharing/rest/content/items/" + id + "/info/thumbnail/ago_downloaded.png")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/info", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/info")
		response, _ := json.Marshal(map[string]interface{}{"owningSystemUrl": "http://" + config.Server,
			"authInfo": map[string]interface{}{"tokenServicesUrl": "https://" + config.Project.Hostname + "/sharing/rest/generateToken", "isTokenBasedSecurity": true}})
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
	//http://reais.x10host.com/arcgis/rest/services/leasecompliance2016/FeatureServer/jobs/replicas?f=json
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/jobs/replicas", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/job/replica")
		var submissionTime int64 = 1441201696150
		var lastUpdatedTime int64 = 1441201705967
		response, _ := json.Marshal(map[string]interface{}{
			"replicaName": "MyReplica", "replicaID": "58808194-921a-4f9f-ac97-5ffd403368a9", "submissionTime": submissionTime, "lastUpdatedTime": lastUpdatedTime,
			"status": "Completed", "resultUrl": "http://" + config.Project.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/unRegisterReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/unRegisterReplica")
		response, _ := json.Marshal(map[string]interface{}{"success": true})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	//http://reais.x10host.com/arcgis/rest/services/leasecompliance2016/FeatureServer/replicas/?f=json
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/replicas/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/replicas")
		var fileName = config.ReplicaPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		log.Println("Sending: " + fileName)
		http.ServeFile(w, r, fileName) //, { root : __dirname})
	}).Methods("GET", "POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/createReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/createReplica (post)")
		response, _ := json.Marshal(map[string]interface{}{"statusUrl": "http://" + config.Project.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas"})
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
			"lastUpdatedTime": lastUpdatedTime, "status": "Completed", "resultUrl": "http://" + config.Project.Hostname + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/services (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}

			ret := config.SetArcCatalog(body, "FeatureServer", "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcCatalog("FeatureServer", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+"FeatureServer.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/arcgis/rest/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/services (post)")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Object moved</title></head><body>" +
			"<h2>Object moved to <a href=\"/arcgis/rest/services\">here</a>.</h2>" +
			"</body></html>"))
	}).Methods("POST")

	r.HandleFunc("/arcgis/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/services")
		log.Println("Sending: " + config.DataPath + "FeatureServer.json")
		response := config.GetArcCatalog("FeatureServer", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			http.ServeFile(w, r, config.DataPath+"FeatureServer.json")
		}
	}).Methods("GET")

	r.HandleFunc("/llarcgis/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/services (post)")

		response := config.GetArcCatalog("FeatureServer", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+"FeatureServer.json")
		}
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

		response := config.GetArcService(name, name, -1, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + name + ".json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+name+".json")
		}
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + " (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "", -1, "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcService(name, "FeatureServer", -1, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
		}
	}).Methods("GET", "PUT")

	r.HandleFunc("/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/rest/services/" + name + "/FeatureServer (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "FeatureServer", -1, "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcService(name, "FeatureServer", -1, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
		}
	}).Methods("GET", "POST", "PUT")

	/*
		r.HandleFunc("/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			name := vars["name"]

			log.Println("/rest/services/" + name + "/FeatureServer")

			response := config.GetArcService(name, "FeatureServer", 0, "")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
			}
		}).Methods("GET", "POST")
	*/
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer (" + r.Method + ")")
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "FeatureServer", -1, "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcService(name, "FeatureServer", -1, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
		}
	}).Methods("GET", "POST", "PUT")
	/*
		r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			name := vars["name"]

			log.Println("/arcgis/rest/services/" + name + "/FeatureServer  (post)")
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
		}).Methods("POST")
	*/

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id, _ := vars["id"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + " (" + r.Method + ")")

		idInt, _ := strconv.Atoi(id)
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte("Error"))
				return
			}
			ret := config.SetArcService(body, name, "FeatureServer", idInt, "")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": ret})
			w.Write(response)
			return
		}

		response := config.GetArcService(name, "FeatureServer", idInt, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
		}
	}).Methods("GET", "POST", "PUT")
	/*
		r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			name := vars["name"]
			id := vars["id"]

			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "(" + r.Method + ")")

			idInt, _ := strconv.Atoi(id)
			if r.Method == "PUT" {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.Write([]byte("Error"))
					return
				}
				ret := config.SetArcService(string(body), name, "FeatureServer", idInt, "")
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": ret})
				w.Write(response)
				return
			}

			response := config.GetArcService(name, "FeatureServer", idInt, "")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
			}
		}).Methods("POST", "PUT")
	*/

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/queryLocal", func(w http.ResponseWriter, r *http.Request) {
		//if(req.query.outFields=='OBJECTID'){
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)
		log.Println(r.FormValue("returnGeometry"))
		log.Println(r.FormValue("outFields"))

		if len(r.FormValue("where")) > 0 {
			w.Header().Set("Content-Type", "application/json")
			var response = []byte("{\"objectIdFieldName\":\"OBJECTID\",\"globalIdFieldName\":\"GlobalID\",\"geometryProperties\":{\"shapeAreaFieldName\":\"Shape__Area\",\"shapeLengthFieldName\":\"Shape__Length\",\"units\":\"esriMeters\"},\"features\":[]}")
			w.Write(response)

		} else if len(r.FormValue("objectIds")) > 0 {
			response := config.GetArcService(name, "FeatureServer", idInt, "query")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

			}
		} else if r.FormValue("returnGeometry") == "false" && strings.Index(r.FormValue("outFields"), "OBJECTID") > -1 { //r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectid")

			response := config.GetArcService(name, "FeatureServer", idInt, "objectid")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectid.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectid.json")
			}
		} else {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")

			response := config.GetArcService(name, "FeatureServer", idInt, "query")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

			}
		}
		//http.ServeFile(w, r, config.DataPath + "/" + id  + "query.json")

	}).Methods("GET", "POST")

	/*
		r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			name := vars["name"]
			id := vars["id"]
			idInt, _ := strconv.Atoi(id)

			if r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
				log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectid (post)")

				response := config.GetArcService(name, "FeatureServer", idInt, "objectid")
				if len(response) > 0 {
					w.Header().Set("Content-Type", "application/json")
					w.Write(response)
				} else {
					log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectid.json")
					http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectid.json")
				}
			} else {
				log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query (post)")

				response := config.GetArcService(name, "FeatureServer", idInt, "query")
				if len(response) > 0 {
					w.Header().Set("Content-Type", "application/json")
					w.Write(response)
				} else {
					log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
					http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")
				}
			}
		}).Methods("GET", "POST")
	*/
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
		var AttachmentPath = config.AttachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		//attachments:=[]interface{}
		attachments := make([]interface{}, 0)
		//[]interface{}
		//fields.Fields, "relatedRecordGroups": []interface{}{result}}
		files, _ := ioutil.ReadDir(AttachmentPath)
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

		//w.Write([]byte(AttachmentPath))
		/*
				if(fs.existsSync(AttachmentPath)){
				   var files = fs.readdirSync(AttachmentPath)
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
		var attachment = config.AttachmentsPath + string(os.PathSeparator) + name + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator) + img + ".jpg"

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
		idInt, _ := strconv.Atoi(id)
		row := vars["row"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/addAttachment")
		// TODO: move and rename the file using req.files.path & .name)
		//res.send(console.dir(req.files))  // DEBUG: display available fields
		var uploadPath = config.AttachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "attachments" + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
		os.MkdirAll(uploadPath, 0755)

		//w.Write([]byte(uploadPath))
		if r.Method == "GET" {
			crutime := time.Now().Unix()
			h := md5.New()
			io.WriteString(h, strconv.FormatInt(crutime, 10))
			token := fmt.Sprintf("%x", h.Sum(nil))

			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, token)
		} else {
			const MAX_MEMORY = 10 * 1024 * 1024
			if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusForbidden)
			}

			for key, value := range r.MultipartForm.Value {
				fmt.Fprintf(w, "%s:%s ", key, value)
				log.Printf("%s:%s", key, value)
			}

			for _, fileHeaders := range r.MultipartForm.File {
				for _, fileHeader := range fileHeaders {
					file, _ := fileHeader.Open()
					path := fmt.Sprintf("%s%s%s", uploadPath, string(os.PathSeparator), fileHeader.Filename)
					log.Println(path)
					buf, _ := ioutil.ReadAll(file)
					ioutil.WriteFile(path, buf, os.ModePerm)
				}
			}
			/*
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
			*/
		}
		response, _ := json.Marshal(map[string]interface{}{"addAttachmentResult": map[string]interface{}{"objectId": idInt, "globalId": nil, "success": true}})
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
			        fstream = fs.createWriteStream(AttachmentPath)
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

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", func(w http.ResponseWriter, r *http.Request) {
		//if(req.query.outFields=='OBJECTID'){
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)
		where := r.FormValue("where")
		outFields := r.FormValue("outFields")
		returnIdsOnly := r.FormValue("returnIdsOnly")
		//returnGeometry := r.FormValue("returnGeometry")
		objectIds := r.FormValue("objectIds")
		//returnIdsOnly = true

		log.Println(r.FormValue("returnGeometry"))
		log.Println(r.FormValue("outFields"))
		//sql := "select "+outFields + " from " +

		if len(where) > 0 {
			//response := config.GetArcQuery(name, "FeatureServer", idInt, "query",objectIds,where)
			w.Header().Set("Content-Type", "application/json")
			//var response = []byte("{\"objectIdFieldName\":\"OBJECTID\",\"globalIdFieldName\":\"GlobalID\",\"geometryProperties\":{\"shapeAreaFieldName\":\"Shape__Area\",\"shapeLengthFieldName\":\"Shape__Length\",\"units\":\"esriMeters\"},\"features\":[]}")
			var response = []byte(`{"objectIdFieldName":"OBJECTID","globalIdFieldName":"GlobalID","geometryProperties":{"shapeLengthFieldName":"","units":"esriMeters"},"features":[]}`)
			w.Write(response)

		} else if returnIdsOnly == "true" {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/objectids")

			response := config.GetArcService(name, "FeatureServer", idInt, "objectids")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".objectids.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".objectids.json")
			}
		} else if len(objectIds) > 0 {

			//only get the select objectIds
			//response := config.GetArcService(name, "FeatureServer", idInt, "query")
			response := config.GetArcQuery(name, "FeatureServer", idInt, "query", objectIds, "")

			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

			}
			//if returnGeometry == "false" &&
		} else if strings.Index(outFields, "OBJECTID") > -1 { //r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/outfields")

			response := config.GetArcService(name, "FeatureServer", idInt, "outfields")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".outfields.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".outfields.json")
			}
		} else {
			log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/query")

			response := config.GetArcService(name, "FeatureServer", idInt, "query")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)
			} else {
				log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json")
				http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".query.json")

			}
		}

		//http.ServeFile(w, r, config.DataPath + "/" + id  + "query.json")

	}).Methods("GET", "POST")

	//http://192.168.2.59:8080/arcgis/rest/services/accommodationagreementrentals/FeatureServer/1/queryRelatedRecords?objectIds=12&outFields=*&relationshipId=2&returnGeometry=true&f=json
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/queryRelatedRecords", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		/*
			if id == "3" {
				jsonstr := `{"fields":[{"name":"OBJECTID","type":"esriFieldTypeOID","alias":"OBJECTID","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"cows","type":"esriFieldTypeSmallInteger","alias":"Cows","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"yearling_heifers","type":"esriFieldTypeSmallInteger","alias":"Yearling heifers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"steer_calves","type":"esriFieldTypeSmallInteger","alias":"Steer calves","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"yearling_steers","type":"esriFieldTypeSmallInteger","alias":"Yearling steers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"bulls","type":"esriFieldTypeSmallInteger","alias":"Bulls","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"mares","type":"esriFieldTypeSmallInteger","alias":"Mares","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"geldings","type":"esriFieldTypeSmallInteger","alias":"Geldings","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"studs","type":"esriFieldTypeSmallInteger","alias":"Studs","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"fillies","type":"esriFieldTypeSmallInteger","alias":"Fillies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"colts","type":"esriFieldTypeSmallInteger","alias":"Colts","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"ewes","type":"esriFieldTypeSmallInteger","alias":"Ewes","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"lambs","type":"esriFieldTypeSmallInteger","alias":"Lambs","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"rams","type":"esriFieldTypeSmallInteger","alias":"Rams","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"wethers","type":"esriFieldTypeSmallInteger","alias":"Wethers","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"kids","type":"esriFieldTypeSmallInteger","alias":"Kids","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"billies","type":"esriFieldTypeSmallInteger","alias":"Billies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"nannies","type":"esriFieldTypeSmallInteger","alias":"Nannies","sqlType":"sqlTypeOther","domain":null,"defaultValue":null},{"name":"Comments","type":"esriFieldTypeString","alias":"Comments","sqlType":"sqlTypeOther","length":8000,"domain":null,"defaultValue":null},{"name":"GlobalGUID","type":"esriFieldTypeGUID","alias":"GlobalGUID","sqlType":"sqlTypeOther","length":38,"domain":null,"defaultValue":null},{"name":"created_user","type":"esriFieldTypeString","alias":"Created user","sqlType":"sqlTypeOther","length":255,"domain":null,"defaultValue":null},{"name":"created_date","type":"esriFieldTypeDate","alias":"Created date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"last_edited_user","type":"esriFieldTypeString","alias":"Last edited user","sqlType":"sqlTypeOther","length":255,"domain":null,"defaultValue":null},{"name":"last_edited_date","type":"esriFieldTypeDate","alias":"Last edited date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"reviewer_name","type":"esriFieldTypeString","alias":"Reviewer name","sqlType":"sqlTypeOther","length":50,"domain":null,"defaultValue":null},{"name":"reviewer_date","type":"esriFieldTypeDate","alias":"Reviewer date","sqlType":"sqlTypeOther","length":8,"domain":null,"defaultValue":null},{"name":"reviewer_title","type":"esriFieldTypeString","alias":"Reviewer title","sqlType":"sqlTypeOther","length":50,"domain":null,"defaultValue":null},{"name":"GlobalID","type":"esriFieldTypeGlobalID","alias":"GlobalID","sqlType":"sqlTypeOther","length":38,"domain":null,"defaultValue":null}],"relatedRecordGroups":[]}`
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonstr))
				return
			}
		*/

		var relationshipId = r.FormValue("relationshipId")
		var objectIds = r.FormValue("objectIds")
		var outFields = r.FormValue("outFields")
		var objectId, _ = strconv.Atoi(objectIds)
		//get fields for the related table
		dID := config.Project.Services[name]["relationships"][relationshipId]["dId"]

		//get the fields json

		if config.DbSource == config.FILE {
			//have to find the joinAttribute value for source and destination
			/*
				var sqlstr = "select " + outFields + " from " + config.Schema +
					config.Project.Services[name]["relationships"][relationshipId]["dTable"].(string) +
					" where " +
					config.Project.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
					config.Project.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
					config.Project.Services[name]["relationships"][relationshipId]["oTable"].(string) +
					" where OBJECTID in(" + config.GetParam(1) + "))"
			*/
			var dJoinKey = config.Project.Services[name]["relationships"][relationshipId]["dJoinKey"].(string)
			var oJoinKey = config.Project.Services[name]["relationships"][relationshipId]["oJoinKey"].(string)

			jsonFile := fmt.Sprint(config.DataPath, string(os.PathSeparator), name+string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", id, ".query.json")
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

			var oJoinVal interface{}
			for _, i := range srcObj.Features {
				//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
				//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))

				if int(i.Attributes["OBJECTID"].(float64)) == objectId {
					oJoinVal = i.Attributes[oJoinKey]
					//i.Attributes["OBJECTID"]
					//fieldObj.Features[k].Attributes = updates[num].Attributes
					break
					//record.RelatedRecord = append(record.RelatedRecord, fieldObj.Features[k].Attributes)
				}
			}

			//strconv.Itoa(int(dID.(float64)))
			jsonFile = fmt.Sprint(config.DataPath, string(os.PathSeparator), name, string(os.PathSeparator), "services", string(os.PathSeparator), "FeatureServer.", dID, ".query.json")
			log.Println(jsonFile)
			file, err1 = ioutil.ReadFile(jsonFile)
			if err1 != nil {
				log.Println(err1)
			}
			var fieldObj structs.FeatureTable

			//map[string]map[string]map[string]
			err = json.Unmarshal(file, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(file))
				log.Println(err.Error())
			}
			var relRecords structs.RelatedRecords
			relRecords.Fields = fieldObj.Fields

			var recordGroup structs.RelatedRecordGroup
			recordGroup.ObjectId = objectId

			//records.RelatedRecordGroups.ObjectId = objectId
			//records.ObjectId = objectId
			//records.RelatedRecord = map[string]interface{}
			//c := 0
			log.Printf("Finding: %v", oJoinVal)

			for k, i := range fieldObj.Features {
				//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
				//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))

				if i.Attributes[dJoinKey] == oJoinVal {
					log.Printf("Found: %v", i.Attributes[dJoinKey])
					var rec structs.RelatedRecord
					//i.Attributes["OBJECTID"]
					//fieldObj.Features[k].Attributes = updates[num].Attributes
					//break
					//var attributes structs.Attribute
					//attributes = fieldObj.Features[k].Attributes
					//rec.Attributes = append(rec.Attributes, fieldObj.Features[k].Attributes)
					rec.Attributes = fieldObj.Features[k].Attributes
					recordGroup.RelatedRecords = append(recordGroup.RelatedRecords, rec)
					//c++
				}

			}

			var jsonstr []byte
			//if c == 0 {
			//	records.RelatedRecordGroups = records.RelatedRecordGroups[:0]
			//}
			if len(recordGroup.RelatedRecords) > 0 {
				relRecords.RelatedRecordGroups = append(relRecords.RelatedRecordGroups, recordGroup)
			} else {
				relRecords.RelatedRecordGroups = make([]structs.RelatedRecordGroup, 0)
			}
			jsonstr, err = json.Marshal(relRecords)
			if err != nil {
				log.Println(err)
			}

			/*
				tx, err := config.Db.Begin()
				if err != nil {
					log.Fatal(err)
				}

				var response []byte
				if len(final_result) > 0 {
					var result = map[string]interface{}{}
					result["objectId"] = objectIds //strconv.Atoi(objectIds)
					result["relatedRecords"] = final_result
					response, _ = json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})
					response = response[1:]
				} else {
					response = []byte("\"relatedRecordGroups\":[]}")
				}
			*/

			//var response []byte
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonstr)
			return
			//response = "{fields:" + fields + "," + response[1]
			//w.Write([]byte("{\"fields\":"))
			//w.Write(fields)
			//w.Write([]byte(","))
			//w.Write(response)
		}
		//idInt, _ := strconv.Atoi(id)

		var sql string
		var fields []byte

		if config.DbSource == config.PGSQL {
			sql = "select json->'fields' from services where service=$1 and name=$2 and layerid=$3 and type=$4"
			log.Printf("select json->'fields' from services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", dID, "")
			stmt, err := config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			err = stmt.QueryRow(name, "FeatureServer", dID, "").Scan(&fields)
			if err != nil {
				log.Println(err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
				return
			}
		} else if config.DbSource == config.SQLITE3 {
			sql = "select json from services where service=? and name=? and layerid=? and type=?"
			log.Printf("select json from services where service='%v' and name='%v' and layerid=%v and type='%v'", name, "FeatureServer", dID, "")
			stmt, err := config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			err = stmt.QueryRow(name, "FeatureServer", dID, "").Scan(&fields)
			if err != nil {
				log.Println(err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"fields\":[],\"relatedRecordGroups\":[]}"))
				return
			}
			//fields = fields["fields"]
			var fieldObj structs.FeatureTable
			//map[string]map[string]map[string]
			err = json.Unmarshal(fields, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(fields))
				log.Println(err.Error())
			}
			fields, err = json.Marshal(fieldObj.Fields)
			if err != nil {
				log.Println(err)
			}

		}

		//_, err = w.Write(fields)
		//return
		//var replicaDb = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		var tableName = config.Project.Services[name]["relationships"][relationshipId]["dTable"].(string)
		log.Println(tableName)
		//var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))
		//var jsonFields=JSON.parse(file)
		//log.Println("sqlite: " + replicaDb)
		//var db = new sqlite3.Database(replicaDb)
		var sqlstr = "select " + outFields + " from " + config.Schema +
			config.Project.Services[name]["relationships"][relationshipId]["dTable"].(string) +
			" where " +
			config.Project.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Project.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
			config.Project.Services[name]["relationships"][relationshipId]["oTable"].(string) +
			" where OBJECTID in(" + config.GetParam(1) + "))"

		//_, err = w.Write([]byte(sqlstr))
		log.Println(strings.Replace(sqlstr, config.GetParam(1), objectIds, -1))

		stmt, err := config.DbQuery.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}

		//outArr := []interface{}{}
		//relationshipIdInt, _ := strconv.Atoi(relationshipId)
		rows, err := stmt.Query(objectIds) //relationshipIdInt
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		//var colLookup = map[string]interface{}{"objectid": "OBJECTID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor"}
		var colLookup = map[string]string{"objectid": "OBJECTID", "globalguid": "GlobalGUID", "globalid": "GlobalID", "creationdate": "CreationDate", "creator": "Creator", "editdate": "EditDate", "editor": "Editor", "comments": "Comments"}
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
				//var v interface{}
				val := values[i]
				if colLookup[col] != "" {
					col = colLookup[col]
				}

				switch t := val.(type) {
				case int:
					//fmt.Printf("Integer: %v=%v\n", col, t)
					tmp_struct[col] = val
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
					//fmt.Printf("String: %v=%v\n", col, val)
					tmp_struct[col] = fmt.Sprintf("%s", val)
				case bool:
					//fmt.Printf("Bool: %v\n", t)
					tmp_struct[col] = val
				case []interface{}:
					for i, n := range t {
						fmt.Printf("Item: %v= %v\n", i, n)
					}
				default:
					var r = reflect.TypeOf(t)
					tmp_struct[col] = r
					//fmt.Printf("Other:%v=%v\n", col, r)
				}
			}
			record := map[string]interface{}{"attributes": tmp_struct}
			final_result = append(final_result, record)
			result_id++
		}

		var response []byte
		if len(final_result) > 0 {
			var result = map[string]interface{}{}
			result["objectId"] = objectIds //strconv.Atoi(objectIds)
			result["relatedRecords"] = final_result
			response, _ = json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})
			response = response[1:]
		} else {
			response = []byte("\"relatedRecordGroups\":[]}")
		}

		//var response []byte
		w.Header().Set("Content-Type", "application/json")
		//response = "{fields:" + fields + "," + response[1]
		w.Write([]byte("{\"fields\":"))
		w.Write(fields)
		w.Write([]byte(","))
		w.Write(response)
		//w.Write([]byte("}"))
	}).Methods("GET", "POST")

	//http://192.168.2.59:8080/arcgis/rest/services/accommodationagreementrentals/FeatureServer/0/applyEdits?f=json&updates=[{%22attributes%22%3A{%22OBJECTID%22%3A2%2C%22suyl%22%3A40%2C%22ov%22%3A5%2C%22bo%22%3A10%2C%22eq%22%3A0%2C%22permittee%22%3A%22Anna+H.+Begay%22%2C%22range_unit%22%3A%22RU255%22%2C%22GlobalID%22%3A%22{425B5BE6-41BE-4A47-92EE-4C4138897DB8}%22}}]
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/applyEdits", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		//idInt, _ := strconv.Atoi(id)
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/applyEdits")
		var response []byte
		if config.DbSource == config.FILE {

			//get the fields json

			jsonFile := config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".query.json"
			log.Println(jsonFile)
			file, err1 := ioutil.ReadFile(jsonFile)
			if err1 != nil {
				log.Println(err1)
			}
			var fieldObj structs.FeatureTable

			//map[string]map[string]map[string]
			err := json.Unmarshal(file, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(file))
				log.Println(err.Error())
			}
			var objectid int
			//var globalID string
			var results []interface{}
			if len(r.FormValue("updates")) > 0 {
				var updates structs.Record
				decoder := json.NewDecoder(strings.NewReader(r.FormValue("updates"))) //r.Body
				err := decoder.Decode(&updates)
				if err != nil {
					panic(err)
				}
				//var objId int
				for k, i := range fieldObj.Features {
					//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
					//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
					if int(i.Attributes["OBJECTID"].(float64)) == objectid {
						//i.Attributes["OBJECTID"]
						fieldObj.Features[k].Attributes = updates[0].Attributes
						break
					}
				}
				var jsonstr []byte
				jsonstr, err = json.Marshal(fieldObj)
				if err != nil {
					log.Println(err)
				}
				err = ioutil.WriteFile(jsonFile, jsonstr, 0644)
				if err != nil {
					log.Println(err1)
				}
				//write json back to file
				result := map[string]interface{}{}
				result["objectId"] = objectid
				result["success"] = true
				result["globalId"] = nil
				results = append(results, result)
				response, _ = json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})

				//response = Updates(name, id, tableName, r.FormValue("updates"))
			} else if len(r.FormValue("adds")) > 0 {
				//response = Adds(name, id, tableName, r.FormValue("adds"))

				var adds []structs.Feature
				decoder := json.NewDecoder(strings.NewReader(r.FormValue("adds"))) //r.Body
				err := decoder.Decode(&adds)
				if err != nil {
					panic(err)
				}
				objectid = len(fieldObj.Features) + 1
				for _, i := range adds {
					i.Attributes["objectId"] = objectid
					fieldObj.Features = append(fieldObj.Features, i)
					objectid++
				}

				var jsonstr []byte
				jsonstr, err = json.Marshal(fieldObj)
				if err != nil {
					log.Println(err)
				}
				err = ioutil.WriteFile(jsonFile, jsonstr, 0644)
				if err != nil {
					log.Println(err1)
				}
				//write json back to file
				result := map[string]interface{}{}
				result["objectId"] = objectid
				result["success"] = true
				result["globalId"] = nil
				results = append(results, result)
				response, _ = json.Marshal(map[string]interface{}{"addResults": results, "updateResults": []string{}, "deleteResults": []string{}})
			} else if len(r.FormValue("deletes")) > 0 {
				//response = Deletes(name, id, tableName, r.FormValue("deletes"))
				objectid, _ = strconv.Atoi(r.FormValue("deletes"))
				if objectid == 0 {
					return
				}
				for k, i := range fieldObj.Features {
					if int(i.Attributes["OBJECTID"].(float64)) == objectid {
						//i.Attributes["OBJECTID"]
						fieldObj.Features = append(fieldObj.Features[:k], fieldObj.Features[k+1:]...)
						break
					}
				}
				var jsonstr []byte
				jsonstr, err = json.Marshal(fieldObj)
				if err != nil {
					log.Println(err)
				}
				err = ioutil.WriteFile(jsonFile, jsonstr, 0644)
				if err != nil {
					log.Println(err1)
				}
				//write json back to file
				result := map[string]interface{}{}
				result["objectId"] = objectid
				result["success"] = true
				result["globalId"] = nil
				results = append(results, result)
				response, _ = json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": []string{}, "deleteResults": results})
			}

		} else {
			var tableName = config.Schema + config.Project.Services[name]["layers"][id]["data"].(string)
			log.Println("Table name: " + tableName)
			//var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))

			if len(r.FormValue("updates")) > 0 {
				response = Updates(name, id, tableName, r.FormValue("updates"))
			} else if len(r.FormValue("adds")) > 0 {
				response = Adds(name, id, tableName, r.FormValue("adds"))
			} else if len(r.FormValue("deletes")) > 0 {
				response = Deletes(name, id, tableName, r.FormValue("deletes"))
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		/*
			sql := "select json->'fields' from services where service=$1 and name=$2 and layerid=$3 and type=$4"
			log.Println(sql)
			log.Println("Values: " + name + "," + "FeatureServer" + "," + id)
			stmt, err := Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			var fields []byte
			err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
			if err != nil {
				log.Println(err.Error())
			}
		*/
		/*
			var replicaDb = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
			//var tableName = config.Services[name]["relationships"][id]["dTable"].(string)
			//log.Println(tableName)
			//var layerId = int(config.Services[name]["relationships"][id]["dId"].(float64))
			//id = "1"
			var jsonFile = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
				id + ".query.json"
			file, err1 := ioutil.ReadFile(jsonFile)
			if err1 != nil {
				fmt.Printf("// error while reading file %s\n", jsonFile)
				fmt.Printf("File error: %v\n", err1)
				os.Exit(1)
			}
		*/
		//var features map[string]interface{}{}
		//var features map[string]interface{}
		//var features map[string]map[string]map[string]map[string]interface{}
		//var features TableField
		/*
			var features []Field
			//map[string]map[string]map[string]
			err = json.Unmarshal(fields, &features)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(fields))
				log.Println(err.Error())
			}
			log.Println("Features dump:")
			log.Print(features)
			b, err1 := json.Marshal(features)
			if err1 != nil {
				log.Println(err1)
			}
			log.Println(string(b))
		*/

		//var replicaDb = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"

		//var jsonFields=JSON.parse(file)
		//log.Println("sqlite: " + replicaDb)
		//var db = new sqlite3.Database(replicaDb)
		/*
			var sqlstr = "select " + outFields + " from " +
				config.Services[name]["relationships"][relationshipId]["dTable"].(string) +
				" where " +
				config.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
				config.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
				config.Services[name]["relationships"][relationshipId]["oTable"].(string) +
				" where OBJECTID=$1)"
		*/

		/*
			var jsonOutputFile = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." +
				id + ".query.exported.json"

			os.Remove(jsonOutputFile)

			b, err1 := json.Marshal(features)
			if err1 != nil {
				log.Println(err1)
			}
			log.Println(string(b))
			ioutil.WriteFile(jsonOutputFile, b, 0644)
		*/

		//now read posted JSON
		//var updates = map[string]interface{}{}

	}).Methods("GET", "POST")

	//http://reais.x10host.com/arcgis/rest/services/leasecompliance2016/FeatureServer/replicas/?f=json
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/db/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)
		fieldStr := r.URL.Query().Get("field")
		if len(fieldStr) == 0 {
			fieldStr = "\"ItemInfo\""
		}
		dbPath := r.URL.Query().Get("db")

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/db/" + id)

		var dbName = config.ReplicaPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		if len(dbPath) > 0 {
			dbName = "file:" + dbPath + "?PRAGMA journal_mode=WAL"
			if config.DbSqliteQuery != nil {
				config.DbSqliteQuery.Close()
				config.DbSqliteQuery = nil
				if dbPath == "close" {
					return
				}
			}
		}
		//err := config.DbSqliteQuery.Ping()

		var err error
		//if err != nil {
		if config.DbSqliteQuery == nil {
			//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
			config.DbSqliteQuery, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal(err)
			}
		}
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)
				return
			}
			//ret := config.SetArcService(body, name, "FeatureServer", idInt, "")
			sql := "update \"GDB_ServiceItems\" set " + fieldStr + "=? where OBJECTID=?"
			log.Println(sql)
			//log.Println(body)
			log.Println(id)
			stmt, err := config.DbSqliteQuery.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)

			}
			_, err = stmt.Exec(string(body), idInt)
			//db.Close()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)

				log.Println(err.Error())
				return
			}
			stmt.Close()
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": "ok"})
			w.Write(response)
			return
		}
		//Db.Exec(initializeStr)
		log.Print("Sqlite database: " + dbName)
		//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
		sql := "SELECT " + fieldStr + " FROM \"GDB_ServiceItems\" where OBJECTID=?"
		log.Printf("Query: "+sql+"%v", idInt)

		stmt, err := config.DbSqliteQuery.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			//w.Write([]byte(err.Error()))
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			return
		}
		//rows := stmt.QueryRow(id)
		var itemInfo []byte
		err = stmt.QueryRow(idInt).Scan(&itemInfo)
		//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(itemInfo)
		/*
			for rows.Next() {
				err = rows.Scan(&itemInfo)
				w.Header().Set("Content-Type", "application/json")

				w.Write(itemInfo)
				//fmt.Println(string(itemInfo))
			}
			rows.Close() //good habit to close
		*/
		//db.Close()

	}).Methods("GET", "POST", "PUT")
	//http://reais.x10host.com/arcgis/rest/services/leasecompliance2016/FeatureServer/xml/31?f=json&db=C:\Users\steve\AppData\Local\Packages\Esri.CollectorforArcGIS_eytg3kh68c6a8\LocalState\hpluser5_qd3vos1n.1xf\df5aa0e91991468eb0efadf475bea54e\n2tel3ls.beb.geodatabase
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/xml/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)
		dbPath := r.URL.Query().Get("db")

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/xml/" + id)
		var dbName = config.ReplicaPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		if len(dbPath) > 0 {
			dbName = "file:" + dbPath + "?PRAGMA journal_mode=WAL"
			if config.DbSqliteQuery != nil {
				config.DbSqliteQuery.Close()
				config.DbSqliteQuery = nil
				if dbPath == "close" {
					return
				}
			}
		}

		var err error
		//if err != nil {
		if config.DbSqliteQuery == nil {
			//config.DbSqliteQuery, err = sql.Open("sqlite3", "file:"+dbName+"?PRAGMA journal_mode=WAL")
			config.DbSqliteQuery, err = sql.Open("sqlite3", dbName)
			if err != nil {
				log.Fatal(err)
			}
		}
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)

				return
			}
			//ret := config.SetArcService(body, name, "FeatureServer", idInt, "")
			sql := "update \"GDB_Items\" set \"Definition\"=? where OBJECTID=?"
			stmt, err := config.DbSqliteQuery.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)

				return
			}
			_, err = stmt.Exec(body, id)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.Header().Set("Content-Type", "application/json")
				response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
				w.Write(response)

				return
			}
			//db.Close()
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": "ok"})
			w.Write(response)
			return
		}
		//Db.Exec(initializeStr)
		log.Print("Sqlite database: " + dbName)
		//sql := "SELECT \"DatasetName\",\"ItemId\",\"ItemInfo\",\"AdvancedDrawingInfo\" FROM \"GDB_ServiceItems\""
		sql := "SELECT \"Definition\" FROM \"GDB_Items\" where OBJECTID=?"
		log.Printf("Query: "+sql+"%v", idInt)

		stmt, err := config.DbSqliteQuery.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

		}
		//rows := stmt.QueryRow(id)
		var itemInfo []byte
		err = stmt.QueryRow(idInt).Scan(&itemInfo)
		//rows, err := Db.Query(sql) //.Scan(&datasetName, &itemId, &itemInfo, &advDrawingInfo)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(map[string]interface{}{"response": err.Error()})
			w.Write(response)

			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(itemInfo)
	}).Methods("GET", "POST", "PUT")

	//put this last - serve static content
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	return r

}

func Updates(name string, id string, tableName string, updateTxt string) []byte {
	//log.Println(updateTxt)

	var updates structs.Record
	decoder := json.NewDecoder(strings.NewReader(updateTxt)) //r.Body

	err := decoder.Decode(&updates)
	if err != nil {
		panic(err)
	}
	//defer r.Body.Close()
	cols := ""
	sep := ""
	c := 1
	//var vals := []interface{}
	//objectid := 1
	var objectid int
	//var globalID string
	var results []interface{}
	//var objId int

	for num, i := range updates {
		var vals []interface{}

		result := map[string]interface{}{}
		for key, j := range i.Attributes {
			//fmt.Println(key + ":  ")
			//var objectid = updates[0].Attributes["OBJECTID"]
			//var globalId = updates[0].Attributes["GlobalID"]
			if key == "OBJECTID" {
				objectid = int(j.(float64))
				result["objectId"] = objectid

				//objId = c
				//c++
				//} else if key == "GlobalID" {
				//	globalID = j.(string)
				//	result["globalId"] = globalID
			} else {
				cols += sep + key + "=" + config.GetParam(c)
				sep = ","
				vals = append(vals, j)
				//fmt.Println(j)
				c++
			}
		}
		vals = append(vals, objectid)
		log.Println("update " + tableName + " set " + cols + " where OBJECTID=" + config.GetParam(len(vals)))
		log.Print(vals)
		//log.Print(objId)
		var sql string
		if config.DbSource == config.PGSQL {
			sql = "update " + tableName + " set " + cols + " where OBJECTID=" + config.GetParam(len(vals))
		} else if config.DbSource == config.SQLITE3 {
			sql = "update " + tableName + " set " + cols + " where OBJECTID=?"
		}

		stmt, err := config.DbQuery.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Println(err.Error())
		}
		result["success"] = true
		result["globalId"] = nil
		results = append(results, result)

		/*
			select pos-1  from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=0 and elem->'attributes'->>'OBJECTID'='$1')::int

			update services set json=jsonb_set(json,
			'{features,26,attributes}',
			'{"OBJECTID":27,"acres":3.12,"lease_site":0,"feature_type":1,"climatic_zone":2,"quad_name":"077-SE-196","elevation":6048,"permittee":"Lorraine / Elsie Begay","homesite_id":"H61A"}'::jsonb,
			false) where type='query' and layerId=0;
		*/
		//sql = "update services set json=jsonb_set(json, array('features',elem_index::text, ,false) from (select pos - 1 as elem_index from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=0 and elem->'attributes'->>'OBJECTID'='$2')"

		updateTxt = updateTxt[15 : len(updateTxt)-2]
		if config.DbSource == config.PGSQL {
			sql = "select pos-1  from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=$1 and elem->'attributes'->>'OBJECTID'=$2"

			log.Println(sql)
			log.Print("Layer ID: ")
			log.Println(id)
			log.Print("Objectid: ")
			log.Println(objectid)
			rows, err := config.Db.Query(sql, id, objectid)
			defer rows.Close()
			var rowId int
			for rows.Next() {
				err := rows.Scan(&rowId)
				if err != nil {
					log.Fatal(err)
				}
			}
			sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
			stmt, err = config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			log.Println(updateTxt)
			log.Println(id)
			_, err = stmt.Exec(updateTxt, id)
			if err != nil {
				log.Println(err.Error())
			}

		} else if config.DbSource == config.SQLITE3 {
			sql = "select json from services where type='query' and layerId=?"
			stmt, err = config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			rows, err := config.Db.Query(sql, id, objectid)
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

			var fieldObj structs.FeatureTable
			//map[string]map[string]map[string]
			err = json.Unmarshal(row, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(row))
				log.Println(err.Error())
			}
			for k, i := range fieldObj.Features {
				//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
				//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
				if int(i.Attributes["OBJECTID"].(float64)) == objectid {
					//i.Attributes["OBJECTID"]
					fieldObj.Features[k].Attributes = updates[num].Attributes
					break
				}
			}
			var jsonstr []byte
			jsonstr, err = json.Marshal(fieldObj)
			if err != nil {
				log.Println(err)
			}
			tx, err := config.Db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			sql = "update services set json=? where type='query' and layerId=?"

			stmt, err = tx.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			idInt, _ := strconv.Atoi(id)
			//log.Printf("%v\n%v", string(jsonstr), idInt)
			//sql = "PRAGMA synchronous = OFF;PRAGMA cache_size=100000;PRAGMA journal_mode=WAL;"
			//tx.Exec(sql)

			_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
			if err != nil {
				log.Println(err.Error())
			}
			tx.Commit()
			//sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
		}
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})
	return response

	//curl -H "Content-Type: application/x-www-form-urlencoded" -X POST -d 'rollbackOnFailure=true&updates=[{"attributes":{"OBJECTID":3,"permittee":"Jack/Bessie Hatathlie","homesite_id":"9w3hdseq78dy","range_unit":551,"acres":3,"lease_site":0,"feature_type":0,"climatic_zone":2,"quad_name":"099-NW-004","elevation":6040,"permittee_globalid":"{D1A2F0B1-6F46-477A-80A9-CF550915B6BB}","has_permittee":1}}]&f=json' http://localhost:81/arcgis/rest/services/leasecompliance2016/FeatureServer/0/applyEdits

	//curl -H "Content-Type: application/x-www-form-urlencoded" -X POST -d 'rollbackOnFailure=true&adds=[{"geometry":"attributes":{"OBJECTID":3,"permittee":"Jack/Bessie Hatathlie","homesite_id":"9w3hdseq78dy","range_unit":551,"acres":3,"lease_site":0,"feature_type":0,"climatic_zone":2,"quad_name":"099-NW-004","elevation":6040,"permittee_globalid":"{D1A2F0B1-6F46-477A-80A9-CF550915B6BB}","has_permittee":1}}]&f=json' http://localhost:81/arcgis/rest/services/leasecompliance2016/FeatureServer/0/applyEdits

	//var jsonvals []interface{}
	//updateTxt := "[{\"attributes\":{\"OBJECTID\":27,\"acres\":3.15,\"lease_site\":0,\"feature_type\":1,\"climatic_zone\":2,\"quad_name\":\"077-SE-196\",\"elevation\":6048,\"permittee\":\"Lorraine / Elsie Begay\",\"homesite_id\":\"H61A\"}}]"
	//updateTxt = strings.Replace(updateTxt[15:len(updateTxt)-1], "\"", "\\\"", -1)

	//jsonvals = append(jsonvals, updateTxt)
	//jsonvals = append(jsonvals, id)
	//jsonvals = append(jsonvals, rowId)

	/*
		_, err = stmt.Exec(jsonvals...)
		if err != nil {
			log.Println(err.Error())
		}
	*/
	/*
		sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}','" + updateTxt + "'::jsonb,false) where type='query' and layerId=$1"
		stmt, err = config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = stmt.Exec(strconv.Atoi(id))
		if err != nil {
			log.Println(err.Error())
		}
	*/

	//log.Println(sql)
	//log.Println(jsonvals)
	/*
		_, err = stmt.Exec(sql, updateTxt, id)
		if err != nil {
			log.Println(err.Error())
		}
	*/

	/*
		var jsonvals []interface{}
		jsonvals = append(jsonvals, updateTxt)

		jsonvals = append(jsonvals, id)

	*/

	//find the matching OBJECTID in the query.json file and update fields and save back to disk
	/*
		for _, i := range updates {
			for _, j := range fields.Fields ["features"] {
				for _, k := range updates[i]["attributes"] {

				}

			}
		}
	*/

	/*
		err2 := json.Unmarshal(r.FormValue("updates"), &updates)
		if err2 != nil {
			log.Println("Error reading configuration file: " + r.FormValue("updates"))
			log.Println(err2.Error())
		}
	*/
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
	//log.Println("sqlite: " + replicaDb)
	//var db = new sqlite3.Database(replicaDb)
	/*
		var sqlstr = "update " + outFields + " from " +
			config.Services[name]["relationships"][id]["dTable"].(string) +
			" where " +
			config.Services[name]["relationships"][id]["dJoinKey"].(string) + " in (select " +
			config.Services[name]["relationships"][id]["oJoinKey"].(string) + " from " +
			config.Services[name]["relationships"][id]["oTable"].(string) +
			" where OBJECTID=?)"

		db, err := sql.Open("sqlite3", replicaDb)
		if err != nil {
			log.Fatal(err)
		}
		defer Db.Close()
		stmt, err := Db.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		//outArr := []interface{}{}
		rows, err := stmt.Query(id)
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
	*/

	//var updates = JSON.parse(req.body.updates)//JSON.parse(req.query.updates)
	/*
			var fs = require("fs')
			var path=DataPath+"/"+name") +"/FeatureServer."+id") + ".query.json"
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

			  var replicaDb = ReplicaPath + "/"+name")+".geodatabase"
			  log.Println("sqlite: " + replicaDb)
			  var db = new sqlite3.Database(replicaDb)
			  //create update statement from json
			  log.Println("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?")
			  log.Println( values )

			  Db.run("UPDATE " + name") + " SET "+fields.join(",")+" WHERE OBJECTID = ?", values)
		  }else{
		 	  results={"objectId":updates.length>0?updates[0].attributes['OBJECTID']:0,"globalId":null,"success":true}
		 	}
	*/
	//update json file with updates
}

func Deletes(name string, id string, tableName string, deletesTxt string) []byte {
	//deletesTxt should be a objectId
	var objectid, _ = strconv.Atoi(deletesTxt)
	var results []interface{}
	result := map[string]interface{}{}
	result["objectId"] = objectid
	result["success"] = true
	result["globalId"] = nil
	results = append(results, result)
	//delete from table
	log.Println("delete from " + tableName + " where OBJECTID=" + config.GetParam(0))
	var sql = "delete from " + tableName + " where OBJECTID=" + config.GetParam(0)
	stmt, err := config.DbQuery.Prepare(sql)
	if err != nil {
		log.Println(err.Error())
	}
	//err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
	_, err = stmt.Exec(objectid)
	if err != nil {
		log.Println(err.Error())
	}

	if config.DbSource == config.PGSQL {
		sql := "select pos-1  from services,jsonb_array_elements(json->'features') with ordinality arr(elem,pos) where type='query' and layerId=$1 and elem->'attributes'->>'OBJECTID'=$2"

		log.Println(sql)
		log.Printf("Layer ID: %v", id)
		log.Printf("Objectid: %v", objectid)

		rows, err := config.Db.Query(sql, id, objectid)
		defer rows.Close()
		var rowId int
		for rows.Next() {
			err := rows.Scan(&rowId)
			if err != nil {
				log.Fatal(err)
			}
		}
		//sql = "update services set json=json->'features' - " + strconv.Itoa(rowId) + " where type='query' and layerId=$1"
		sql = "update services set json=json #- '{features," + strconv.Itoa(rowId) + "}' where type='query' and layerId=$1"
		log.Println(sql)
		log.Printf("Row id: %v", rowId)
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			log.Println(err.Error())
		}

	} else if config.DbSource == config.SQLITE3 {
		sql := "select json from services where type='query' and layerId=?"
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		rows, err := config.Db.Query(sql, id, objectid)
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

		var fieldObj structs.FeatureTable
		//map[string]map[string]map[string]
		err = json.Unmarshal(row, &fieldObj)
		if err != nil {
			log.Println("Error unmarshalling fields into features object: " + string(row))
			log.Println(err.Error())
		}
		for k, i := range fieldObj.Features {
			//if fieldObj.Features[i].Attributes["OBJECTID"] == objectid {
			//log.Printf("%v:%v", i.Attributes["OBJECTID"].(float64), strconv.Itoa(objectid))
			if int(i.Attributes["OBJECTID"].(float64)) == objectid {
				//i.Attributes["OBJECTID"]
				//fieldObj.Features = fieldObj.Features[k]
				fieldObj.Features = append(fieldObj.Features[:k], fieldObj.Features[k+1:]...)
				//fieldObj.Features[k].Attributes = updates[num].Attributes
				break
			}
		}
		var jsonstr []byte
		jsonstr, err = json.Marshal(fieldObj)
		if err != nil {
			log.Println(err)
		}
		tx, err := config.Db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		sql = "update services set json=? where type='query' and layerId=?"

		stmt, err = tx.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}

		idInt, _ := strconv.Atoi(id)

		_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
		if err != nil {
			log.Println(err.Error())
		}
		tx.Commit()
		//sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": []string{}, "deleteResults": results})
	return response

}
func Adds(name string, id string, tableName string, addsTxt string) []byte {
	var results []interface{}
	var objectid int
	log.Println(addsTxt)
	var adds []structs.Feature
	decoder := json.NewDecoder(strings.NewReader(addsTxt)) //r.Body
	err := decoder.Decode(&adds)
	if err != nil {
		panic(err)
	}
	cols := ""
	p := ""
	sep := ""
	c := 1

	sql := "select max(OBJECTID)+1 from " + tableName
	log.Println(sql)
	rows, err := config.Db.Query(sql)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&objectid)
		if err != nil {
			//log.Fatal(err)
			objectid = 1
		}
	}
	for _, i := range adds {
		var vals []interface{}

		for key, j := range i.Attributes {
			if key == "OBJECTID" {
				cols += sep + key
				p += sep + config.GetParam(c)
				sep = ","
				vals = append(vals, objectid)
				c++
			} else {
				cols += sep + key
				p += sep + config.GetParam(c)
				sep = ","
				vals = append(vals, j)
				c++
			}
		}
		cols += sep + "GlobalId"
		p += sep + config.GetParam(c)
		vals = append(vals, "")

		log.Println("insert into " + tableName + "(" + cols + ") values( " + p + ")")
		log.Print(vals)
		sql := "insert into " + tableName + "(" + cols + ") values( " + p + ")"
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Println(err.Error())
		}
		result := map[string]interface{}{}
		result["objectId"] = objectid
		result["success"] = true
		result["globalId"] = nil

		results = append(results, result)

		if config.DbSource == config.PGSQL {

			//addsTxt = addsTxt[15 : len(addsTxt)-2]

			sql = "update services set json=jsonb_set(json,'{features}',$1::jsonb,true) where type='query' and layerId=$2"
			log.Println(sql)
			stmt, err = config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			log.Println(i)
			log.Println(id)
			var jsonstr []byte
			jsonstr, err = json.Marshal(i)
			if err != nil {
				log.Println(err)
			}

			_, err = stmt.Exec(jsonstr, id)
			if err != nil {
				log.Println(err.Error())
			}
		} else if config.DbSource == config.SQLITE3 {
			sql := "select json from services where type='query' and layerId=?"
			stmt, err := config.Db.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}
			rows, err := config.Db.Query(sql, id, objectid)
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

			var fieldObj structs.FeatureTable
			err = json.Unmarshal(row, &fieldObj)
			if err != nil {
				log.Println("Error unmarshalling fields into features object: " + string(row))
				log.Println(err.Error())
			}
			var jsonstr []byte
			jsonstr, err = json.Marshal(fieldObj)
			if err != nil {
				log.Println(err)
			}
			fieldObj.Features = append(fieldObj.Features, i)
			tx, err := config.Db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			sql = "update services set json=? where type='query' and layerId=?"

			stmt, err = tx.Prepare(sql)
			if err != nil {
				log.Println(err.Error())
			}

			idInt, _ := strconv.Atoi(id)

			_, err = tx.Stmt(stmt).Exec(string(jsonstr), idInt)
			if err != nil {
				log.Println(err.Error())
			}
			tx.Commit()
		}
		objectid++
	}
	response, _ := json.Marshal(map[string]interface{}{"addResults": results, "updateResults": []string{}, "deleteResults": []string{}})
	return response
}
