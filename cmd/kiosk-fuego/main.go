package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/asolheiro/kiosk-api/internal-v2/routers"
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite"
)

func main() {
    fmt.Println("Starting kiosk-api...")
    
    dbPath := os.Getenv("SQLITE_DB_PATH")
    if dbPath == "" {
        log.Fatal(fmt.Println("SQLITE_DB_PATH environment variable is not set"))
    }
    
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        log.Fatal("error connecting to database.\nerr: ", err)
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
	url := os.Getenv("URL")
    s := fuego.NewServer(
		fuego.WithAddr(url),
		fuego.WithGlobalMiddlewares(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		}).Handler),
	)
	
    
    fuego.Get(s, "/healthcheck", controller.HealthCheck)

	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalln(err)
	}

	logger = logger.Named("kiosk-api")
	defer func() { _ = logger.Sync() }()
	
    routers.NewRouter(s, db, logger)
    
    fmt.Println("Starting server at port 9999...")
	
	    if err := s.Run(); err != nil {
        log.Fatal("Server error:", err)
    }
    
    fmt.Println("Goodbye...")
}