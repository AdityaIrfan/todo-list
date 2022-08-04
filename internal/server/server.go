package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	viperPkg "github.com/spf13/viper"
	baseApp "github.com/todo-list/internal/app"
	"github.com/todo-list/pkg/logger"
	"github.com/todo-list/pkg/postgres"
	"github.com/todo-list/pkg/viper"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func Run() {

	///load config
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../..")
	config := &viper.EnvConfig{
		FileName: "config",
		FileType: "yaml",
		Path:     basepath,
	}
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}

	//load connection postgre
	pg, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := pg.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	zap, err := logger.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	//load fiber
	app := fiber.New(fiber.Config{
		IdleTimeout: 5,
	})
	app.Use(
		recover.New(),
		compress.New(),
		etag.New(),
		cors.New(),
		fiberlog.New(),
	)

	rh := &baseApp.Handlers{
		Postgres: pg,
		R:        app,
		Logger:   zap,
	}
	rh.SetupRouter()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":" + viperPkg.GetString("server.port")); err != nil {
			log.Panicf("failed listen into port %v", err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	var _ = <-c // This blocks the main thread until an interrupt is received
	log.Println("gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	sqlDB.Close()
	fmt.Println("services was successful shutdown.")
}
