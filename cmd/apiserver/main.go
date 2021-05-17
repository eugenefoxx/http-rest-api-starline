package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/apiserver"
	_ "github.com/lib/pq"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")

	//flag.StringVar(&configPath, "config-path", "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images"))))
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
		ErrorLogger.Printf(err.Error())
	}

	//	s := apiserver.New(config)
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
		ErrorLogger.Printf(err.Error())
	}

}

// GeneralLogger exported
var GeneralLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func initLog() {
	//	absPath, err := filepath.Abs("./outputs/log")
	//	if err != nil {
	//		fmt.Println("Error reading given path:", err)
	//	}

	//	generalLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//generalLog, err := os.OpenFile(viper.GetString("log.logfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	generalLog, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api/logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//generalLog, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/starLine/motivationUpdate/log/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

/*
type Logger struct {
	//	logg     *log.Logger
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewLogger() *Logger {
	f, err := os.OpenFile("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api/logfile.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	s := &Logger{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	return s
}
*/
