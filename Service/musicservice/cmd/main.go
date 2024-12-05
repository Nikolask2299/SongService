package main

import (

	"client"
	"client/server"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"musicservice/cmd/migration"
	"musicservice/interal/app"
	"musicservice/interal/server"
	"musicservice/pkg/config"
	"musicservice/pkg/sql/postgres"

    _"github.com/swaggo/http-swagger"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)

type Server struct{}

func NewServer() Server {
 return Server{}
}


func (Server) GetInfo(w http.ResponseWriter, r *http.Request, param api.GetInfoParams) {
    w.Header().Set("Content-type", "application/json")
    
    data := &api.SongDetail{
        Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
        ReleaseDate: "16.07.2006",
        Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
    }
    
    json.NewEncoder(w).Encode(data)
}

func initializing() {
    migration.Migrations()
    
    app := NewServer()

    log.Println("Initializing server...")
    e := http.NewServeMux()
    h := api.HandlerFromMux(app, e)
    
    s := &http.Server{
        Addr:    "0.0.0.0:8070",
        Handler: h,
    }
    
    if err := s.ListenAndServe(); err != nil {
        log.Fatal(err)
    }

}
// @title           Swagger Example API
// @version         2.0
// @description     This is a sample server celler server.

// @host      localhost:8080
// @BasePath  /
func main() {
    go initializing()

	loger := setupLogger("local")
	loger = loger.With(slog.String("env", "local"))

	loger.Info("initializing server") 
    
    confPost, err := config.ReturnedDatabase()
    if err!= nil {
        loger.Error("error initializing config", slog.String("error", err.Error()))
        panic(err)
    }

    loger.Info("connecting to database " + confPost.DBName)
    postgres, err := postgres.NewPostgres(confPost.User, confPost.Password, confPost.DBName, confPost.Host, confPost.Port)
    if err != nil {
       loger.Error("error initializing postgres", slog.String("error", err.Error()))
       panic(err)
    }

    loger.Info("initializing server config")
    confServer, confAPI, err := config.RetuneServerConfig()
    if err!= nil {
        loger.Error("error initializing config", slog.String("error", err.Error()))
        panic(err)
    }

    loger.Info("initializing client config")
    clientMusic, err := client.NewClientWithResponses("http://" + confAPI.Server.Host + ":" + confAPI.Server.Port, client.WithHTTPClient(&http.Client{}))
    if err != nil {
        loger.Error("error initializing client", slog.String("error", err.Error()))
        panic(err)
    }

    loger.Info("initializing server app")  
    app := app.NewApp(loger, postgres, clientMusic)
    server := server.NewMysicServer(loger, *app)

    loger.Info("Initializing server endpoints")
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        loger.Info("Received request: " + r.URL.String())
        fmt.Fprint(w, "Server listening on " + r.URL.Host)
    })
    http.HandleFunc("/search", server.GetData)
    http.HandleFunc("/text", server.GetText)
    http.HandleFunc("/delete", server.DeleteSong)
    http.HandleFunc("/update", server.UpdateSong)
    http.HandleFunc("/create", server.CreateSong)
    
    loger.Info("Starting server..." + confServer.Host + " " + confServer.Port)
    err = http.ListenAndServe(confServer.Host + ":" + confServer.Port, nil) 
    if err != nil {
        loger.Error("error starting server", slog.String("error",err.Error()))
        panic(err)
    }
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

    switch env {
    case envLocal:
        logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    default:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return logger
}

