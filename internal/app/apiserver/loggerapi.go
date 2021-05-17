package apiserver

import (
	"fmt"
	"log"
	"os"
)

// GeneralLogger exported
var GeneralLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func InitLog() {
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
