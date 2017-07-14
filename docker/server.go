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

	"github.com/gorilla/handlers"
	//_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"

	config "github.com/traderboy/arcrestgo/config"
	routes "github.com/traderboy/arcrestgo/routes"
)

//_ "github.com/mattn/go-sqlite3"
var logPath = "logfile.txt"

func main() {

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//load directory of json files into postgresql

	logParam := flag.Bool("log", false, "a bool")
	//sqliteParam := flag.Bool("sqlite", false, "a bool")

	if *logParam {
		InitLog()
		log.Println("Writing log file to : logfile.txt")
	} else {
		log.SetOutput(os.Stdout)
		log.Println("Writing log file to stdOut")
	}

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
