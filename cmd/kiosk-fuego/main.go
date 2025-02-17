package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/asolheiro/kiosk-api/internal-v2/logging"
	"github.com/asolheiro/kiosk-api/internal-v2/routers"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite"
)

func main() {
    fmt.Println("Starting kiosk-api...")
    
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	

    fmt.Println("Goodbye...")
}


func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	handler := logging.NewZapHandler(logger)
	slogHandler := slog.New(handler)
	
	logger = logger.Named("kiosk-api")
	defer func() { _ = logger.Sync() }()

	dbPath := os.Getenv("SQLITE_DB_PATH")
    if dbPath == "" {
        log.Fatal(fmt.Println("SQLITE_DB_PATH environment variable is not set"))
    }
    
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        log.Fatal("error connecting to database.\nerr: ", err)
    }
	defer db.Close()

	var ddl string
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
	addr := fmt.Sprintf("%s:%s", os.Getenv("URL"), os.Getenv("PORT"))
    s := fuego.NewServer(
		fuego.WithAddr(addr),
		fuego.WithLogHandler(slogHandler.Handler()),
		fuego.WithGlobalMiddlewares(
			cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			}).Handler,
		),
	)

	_ = infoAPI(s)

	fuego.Get(s, "/healthcheck", controller.HealthCheck)
	
	
	
    
	
    routers.NewRouter(s, db, logger)
    
    fmt.Printf("Starting server at port %s...", os.Getenv("PORT"))
	
	    if err := s.Run(); err != nil {
        log.Fatal("Server error:", err)
    }
   return nil 
}

func infoAPI(s *fuego.Server) error {
	url := fmt.Sprintf("http://%s:%v", os.Getenv("URL"), os.Getenv("PORT"))

	s.OpenAPI.Description().Servers = append(
		s.OpenAPI.Description().Servers, 
		&openapi3.Server{
			URL: url,
			Description: "Test server",
	})

	s.OpenAPI.Description().Info.Title = "Kiosk API"
	s.OpenAPI.Description().Info.Contact = &openapi3.Contact{
		Name: "Armando Solheiro",
		Email: "avgsolheiro@gmail.com",
		URL: "https//github.com/asolheiro",
	}
	s.OpenAPI.Description().Info.Version = "v2.0.0"
	return nil
}