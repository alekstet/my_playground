package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alekstet/my_playground/api"
	"github.com/alekstet/my_playground/code"
	"github.com/alekstet/my_playground/config"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/julienschmidt/httprouter"
)

func Run(cnfName *string) {
	var cfg config.PlayConfig
	err := cleanenv.ReadConfig(*cnfName, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	codeStore, err := code.NewCodeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	store, err := api.NewStore(cfg, *codeStore)
	if err != nil {
		store.Log.Fatal(err)
	}

	router := httprouter.New()
	store.Register(router)

	log.Println("Start server on", cfg.Port)
	err = http.ListenAndServe(cfg.Port, router)
	store.Log.Fatal(err)
}

func main() {
	cnfName := flag.String("config", "config/config.yml", "config name")
	flag.Parse()

	Run(cnfName)
}
