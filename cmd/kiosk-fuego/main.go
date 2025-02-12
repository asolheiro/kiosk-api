package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	

    fmt.Println("Goodbye...")
}

var ddl string
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

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
	
    s := fuego.NewServer(
		fuego.WithAddr(os.Getenv("URL")),
		fuego.WithLogHandler(slogHandler.Handler()),
		fuego.WithGlobalMiddlewares(
			cors.New(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			}).Handler,
		),
	)
	infoAPI(s)

	fuego.Get(s, "/healthcheck", controller.HealthCheck)
	
	
	
    
	
    routers.NewRouter(s, db, logger)
    
    fmt.Println("Starting server at port 9999...")
	
	    if err := s.Run(); err != nil {
        log.Fatal("Server error:", err)
    }
   return nil 
}

func infoAPI(s *fuego.Server) {
	url := os.Getenv(os.Getenv("URL"))

	s.OpenAPI.Description().Servers = append(
		s.OpenAPI.Description().Servers, 
		&openapi3.Server{
			URL: fmt.Sprintf("http://%v", url),
			Description: "Test server",
	})
	s.OpenAPI.Description().Info.Title = "Kiosk API"
	s.OpenAPI.Description().Info.Contact = &openapi3.Contact{
		Name: "Armando Solheiro",
		Email: "avgsolheiro@gmail.com",
		URL: "https//github.com/asolheiro",
	}
}