package main

import (
	"flag"
	"log"

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
	}

	//	s := apiserver.New(config)
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}

}
