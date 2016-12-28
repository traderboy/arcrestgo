package routes

import (
	"crypto/md5"
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

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	config "github.com/traderboy/arcrestgo/config"
	structs "github.com/traderboy/arcrestgo/structs"
)

func StartGorillaMux() *mux.Router {
	config.Init()
	r := mux.NewRouter()

	/*
	   Download certs
	*/
	r.HandleFunc("/cert", func(w http.ResponseWriter, r *http.Request) {
		//res.sendFile("certs/server.crt", { root : __dirname})
		log.Println("Sending: " + config.RootPath + "certs/server.crt")
		http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"certs"+string(os.PathSeparator)+"server.crt")
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
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		//w.Write(response)
		//setHeaders(c)
		//fmt.Println(response)
		//w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest (post)")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
		//w.Write(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("POST")

	r.HandleFunc("/sharing/rest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest")
		response, _ := json.Marshal(map[string]interface{}{"currentVersion": config.ArcGisVersion})
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
		log.Println("/sharing/{rest}/accounts/self")

		response := config.GetArcCatalog("portals", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "portals.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"portals.self.json")
		}
	}).Methods("GET")

	r.HandleFunc("/sharing//accounts/self", func(w http.ResponseWriter, r *http.Request) {

		response := config.GetArcCatalog("account", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "account.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"account.self.json")
		}
	}).Methods("GET")

	//no customization necesssary except for username
	r.HandleFunc("/sharing/rest/portals/self", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/portals/self")

		response := config.GetArcCatalog("portals", "self")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "portals.self.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"portals.self.json")
		}
		//http.ServeFile(w, r, config.RootPath + "/portals_self.json")
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
		//load from db
		response := config.GetArcService("%", "content", 0, "")
		if len(response) > 0 {
			//log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+id+"services"+string(os.PathSeparator)+"content.items.json")
		}
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/content/items/{id}/data", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("/sharing/rest/content/items/" + id + "/data")

		response := config.GetArcService("%", "content", 0, "data")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + id + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "content.items.data.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+id+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"content.items.data.json")
		}
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/search")
		//vars := mux.Vars(r)

		//q := vars["q"]
		//q := r.Queries("q")
		q := r.FormValue("q")
		if strings.Index(q, "typekeywords") == -1 {

			response := config.GetArcCatalog("community", "groups")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)

			} else {
				log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.groups.json")
				http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.groups.json")
			}
		} else {

			response := config.GetArcCatalog("search", "")
			if len(response) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.Write(response)

			} else {
				log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "search.json")
				http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"search.json")
			}
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
	/*
		r.HandleFunc("/sharing/rest/community/groups", func(w http.ResponseWriter, r *http.Request) {
			log.Println("/sharing/rest/community/groups")
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.groups.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.groups.json")
		}).Methods("GET")
	*/

	r.HandleFunc("/sharing/rest/community/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		log.Println("/sharing/rest/community/users/" + user)
		response := config.GetArcCatalog("community", "users")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {

			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.users.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.users.json")
		}
	}).Methods("GET")

	r.HandleFunc("/sharing/rest/community/users", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/sharing/rest/community/users/")

		response := config.GetArcCatalog("community", "users")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		} else {
			log.Println("Sending: " + config.RootPath + string(os.PathSeparator) + "community.users.json")
			http.ServeFile(w, r, config.RootPath+string(os.PathSeparator)+"community.users.json")
		}
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
			"authInfo": map[string]interface{}{"tokenServicesUrl": "https://" + config.Server + "/sharing/rest/generateToken", "isTokenBasedSecurity": true}})
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
			"status": "Completed", "resultUrl": "http://" + config.Server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
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
		var fileName = config.ReplicaPath + "/" + name + ".geodatabase"
		log.Println("Sending: " + fileName)
		http.ServeFile(w, r, fileName) //, { root : __dirname})
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/createReplica", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/createReplica (post)")
		response, _ := json.Marshal(map[string]interface{}{"statusUrl": "http://" + config.Server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas"})
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
			"lastUpdatedTime": lastUpdatedTime, "status": "Completed", "resultUrl": "http://" + config.Server + "/arcgis/rest/services/" + name + "/FeatureServer/replicas/"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/arcgis/rest/services")

		response := config.GetArcCatalog("FeatureServer", "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+"FeatureServer.json")
		}
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

		response := config.GetArcService(name, name, 0, "")
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

		log.Println("/arcgis/rest/services/" + name)

		response := config.GetArcService(name, "FeatureServer", 0, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"FeatureServer.json")
		}
	}).Methods("GET")

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

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer")

		response := config.GetArcService(name, "FeatureServer", 0, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer.json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer.json")
		}
	}).Methods("GET", "POST")
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

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id)

		idInt, _ := strconv.Atoi(id)
		response := config.GetArcService(name, "FeatureServer", idInt, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
		}
	}).Methods("GET")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]

		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "  (post)")

		idInt, _ := strconv.Atoi(id)
		response := config.GetArcService(name, "FeatureServer", idInt, "")
		if len(response) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			log.Println("Sending: " + config.DataPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "FeatureServer." + id + ".json")
			http.ServeFile(w, r, config.DataPath+string(os.PathSeparator)+name+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"FeatureServer."+id+".json")
		}
	}).Methods("POST")

	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/query", func(w http.ResponseWriter, r *http.Request) {
		//if(req.query.outFields=='OBJECTID'){
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)

		if r.FormValue("returnGeometry") == "false" && r.FormValue("outFields") == "OBJECTID" {
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
		row := vars["row"]
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/" + row + "/addAttachment")
		// TODO: move and rename the file using req.files.path & .name)
		//res.send(console.dir(req.files))  // DEBUG: display available fields
		var uploadPath = config.AttachmentsPath + string(os.PathSeparator) + name + string(os.PathSeparator) + id + string(os.PathSeparator) + row + string(os.PathSeparator)
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

	//http://192.168.2.59:8080/arcgis/rest/services/accommodationagreementrentals/FeatureServer/1/queryRelatedRecords?objectIds=12&outFields=*&relationshipId=2&returnGeometry=true&f=json
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/queryRelatedRecords", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		idInt, _ := strconv.Atoi(id)

		sql := "select json->'fields' from services where service=$1 and name=$2 and layerid=$3 and type=$4"
		stmt, err := config.Db.Prepare(sql)
		if err != nil {
			log.Println(err.Error())
		}
		var fields []byte
		err = stmt.QueryRow(name, "FeatureServer", idInt, "").Scan(&fields)
		if err != nil {
			log.Println(err.Error())
		}
		//_, err = w.Write(fields)

		//return

		var relationshipId = r.FormValue("relationshipId")
		var objectIds, _ = strconv.Atoi(r.FormValue("objectIds"))
		var outFields = r.FormValue("outFields")

		//var replicaDb = config.RootPath + string(os.PathSeparator) + name + string(os.PathSeparator) + "replicas" + string(os.PathSeparator) + name + ".geodatabase"
		var tableName = config.Project.Services[name]["relationships"][relationshipId]["dTable"].(string)
		log.Println(tableName)
		//var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))

		//var jsonFields=JSON.parse(file)
		//log.Println("sqlite: " + replicaDb)
		//var db = new sqlite3.Database(replicaDb)
		var sqlstr = "select " + outFields + " from postgres." +
			config.Project.Services[name]["relationships"][relationshipId]["dTable"].(string) +
			" where " +
			config.Project.Services[name]["relationships"][relationshipId]["dJoinKey"].(string) + " in (select " +
			config.Project.Services[name]["relationships"][relationshipId]["oJoinKey"].(string) + " from " +
			config.Project.Services[name]["relationships"][relationshipId]["oTable"].(string) +
			" where OBJECTID=$1)"

		//_, err = w.Write([]byte(sqlstr))
		log.Println(sqlstr)

		stmt, err = config.Db.Prepare(sqlstr)
		if err != nil {
			log.Fatal(err)
		}

		//outArr := []interface{}{}
		relationshipIdInt, _ := strconv.Atoi(relationshipId)
		rows, err := stmt.Query(relationshipIdInt)
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
				//var v interface{}
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
					fmt.Printf("String: %v\n", val)
					tmp_struct[col] = fmt.Sprintf("%s", val)
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
			}
			record := map[string]interface{}{"attributes": tmp_struct}
			final_result = append(final_result, record)
			result_id++
		}

		var result = map[string]interface{}{}
		result["objectId"] = objectIds //strconv.Atoi(objectIds)
		result["relatedRecords"] = final_result

		response, _ := json.Marshal(map[string]interface{}{"relatedRecordGroups": []interface{}{result}})

		//var response []byte
		w.Header().Set("Content-Type", "application/json")
		//response = "{fields:" + fields + "," + response[1]
		w.Write([]byte("{\"fields\":"))
		w.Write(fields)
		w.Write([]byte(","))
		w.Write(response[1:])
	}).Methods("GET", "POST")

	//http://192.168.2.59:8080/arcgis/rest/services/accommodationagreementrentals/FeatureServer/0/applyEdits?f=json&updates=[{%22attributes%22%3A{%22OBJECTID%22%3A2%2C%22suyl%22%3A40%2C%22ov%22%3A5%2C%22bo%22%3A10%2C%22eq%22%3A0%2C%22permittee%22%3A%22Anna+H.+Begay%22%2C%22range_unit%22%3A%22RU255%22%2C%22GlobalID%22%3A%22{425B5BE6-41BE-4A47-92EE-4C4138897DB8}%22}}]
	r.HandleFunc("/arcgis/rest/services/{name}/FeatureServer/{id}/applyEdits", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		id := vars["id"]
		//idInt, _ := strconv.Atoi(id)
		log.Println("/arcgis/rest/services/" + name + "/FeatureServer/" + id + "/applyEdits")

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
		var tableName = "postgres." + config.Project.Services[name]["layers"][id]["data"].(string)
		log.Println("Table name: " + tableName)
		//var layerId = int(config.Services[name]["relationships"][relationshipId]["dId"].(float64))

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
		var updates structs.Record
		var updateTxt = r.FormValue("updates")
		log.Println(updateTxt)
		decoder := json.NewDecoder(strings.NewReader(updateTxt)) //r.Body

		err := decoder.Decode(&updates)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		cols := ""
		sep := ""
		c := 1
		//var vals := []interface{}
		//objectid := 1
		var objectid int
		//var globalID string
		var results []interface{}
		var objId string

		for _, i := range updates {
			var vals []interface{}

			result := map[string]interface{}{}
			for key, j := range i.Attributes {
				//fmt.Println(key + ":  ")
				//var objectid = updates[0].Attributes["OBJECTID"]
				//var globalId = updates[0].Attributes["GlobalID"]
				if key == "OBJECTID" {
					objectid = int(j.(float64))
					result["objectId"] = objectid
					vals = append(vals, objectid)
					objId = strconv.Itoa(c)
					c++
					//} else if key == "GlobalID" {
					//	globalID = j.(string)
					//	result["globalId"] = globalID
				} else {
					cols += sep + key + "=$" + strconv.Itoa(c)
					sep = ","
					vals = append(vals, j)
					//fmt.Println(j)
					c++
				}
			}
			log.Println("update " + tableName + " set " + cols + " where OBJECTID=$" + objId)
			log.Print(vals)
			log.Print(objId)

			sql := "update " + tableName + " set " + cols + " where OBJECTID=$" + objId
			stmt, err := config.Db.Prepare(sql)
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

			//var jsonvals []interface{}
			//updateTxt := "[{\"attributes\":{\"OBJECTID\":27,\"acres\":3.15,\"lease_site\":0,\"feature_type\":1,\"climatic_zone\":2,\"quad_name\":\"077-SE-196\",\"elevation\":6048,\"permittee\":\"Lorraine / Elsie Begay\",\"homesite_id\":\"H61A\"}}]"
			//updateTxt = strings.Replace(updateTxt[15:len(updateTxt)-1], "\"", "\\\"", -1)
			updateTxt = updateTxt[15 : len(updateTxt)-2]

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

			sql = "update services set json=jsonb_set(json,'{features," + strconv.Itoa(rowId) + ",attributes}',$1::jsonb,false) where type='query' and layerId=$2"
			log.Println(sql)
			//log.Println(jsonvals)
			/*
				_, err = stmt.Exec(sql, updateTxt, id)
				if err != nil {
					log.Println(err.Error())
				}
			*/

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

			/*
				var jsonvals []interface{}
				jsonvals = append(jsonvals, updateTxt)

				jsonvals = append(jsonvals, id)

			*/
		}
		response, _ := json.Marshal(map[string]interface{}{"addResults": []string{}, "updateResults": results, "deleteResults": []string{}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

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

	}).Methods("GET", "POST")
	return r

}
