package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~rxw/msgapi/db"
	"git.sr.ht/~rxw/msgapi/handler"
	"github.com/BurntSushi/toml"
	//"strconv"
)

type tomlConfig struct {
	API api `toml:"api"`
	DB database `toml:"database"`
}

type api struct {
	Port string
}

type database struct {
	Host string
	Name string
	Pass string
	Port string
	User string
}


func main() {

	var conf tomlConfig

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("API port: %d\n", conf.API.Port)
	//fmt.Printf("DB port: %d\n", conf.DB.Port)

	//addr, err := strconv.Atoi(conf.API.Port); err != nil {
	//	fmt.Println("Bad port")
	//	return
	//}

	log.Printf("%s", string(":" + string(conf.API.Port)))


	listener, err := net.Listen("tcp", ":" + string(conf.API.Port))
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}

	/* Database */

	database, err := db.Initialize(conf.DB.User, conf.DB.Pass, conf.DB.Name)

	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	defer database.Conn.Close()

	httpHandler := handler.NewHandler(database)
	server := &http.Server{
		Handler: httpHandler,
	}

	go func() {
		server.Serve(listener)
	}()

	defer Stop(server)

	log.Printf("Started server on %s", string(conf.API.Port))

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
