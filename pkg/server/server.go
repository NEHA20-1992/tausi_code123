package server

import (
	"fmt"
	"net/http"

	"github.com/NEHA20-1992/tausi_code/api"
	"github.com/NEHA20-1992/tausi_code/api/middleware"
	"github.com/NEHA20-1992/tausi_code/api/seeder"
	"github.com/NEHA20-1992/tausi_code/pkg/config"
	"github.com/NEHA20-1992/tausi_code/pkg/logger"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
	//_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type Server struct {
	DB       *gorm.DB
	Router   *mux.Router
	prepared bool
}

var ApiServer Server

func (server *Server) initialize() {
	var err error

	configuration := config.ServerConfiguration.Database["tausi"]

	if configuration.Driver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", configuration.UserName, configuration.Password, configuration.HostName, configuration.PortNumber, configuration.DatabaseName)
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			logger.BootstrapLogger.Fatalln("Cannot connect to %s database", configuration.Driver)
			logger.BootstrapLogger.Fatal("This is the error:", err)
		} else {
			logger.BootstrapLogger.Infoln(fmt.Sprintf("We are connected to the %s database", configuration.Driver))
		}
	}
	// if configuration.Driver == "postgres" {
	// 	DBURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", configuration.HostName, configuration.PortNumber, configuration.UserName, configuration.DatabaseName, configuration.Password)
	// 	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	// 	if err != nil {
	// 		logger.BootstrapLogger.Fatalln("Cannot connect to %s database", configuration.Driver)
	// 		logger.BootstrapLogger.Fatal("This is the error:", err)
	// 	} else {
	// 		logger.BootstrapLogger.Infoln(fmt.Sprintf("We are connected to the %s database", configuration.Driver))
	// 	}
	// }

	//server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	server.Router = mux.NewRouter()
}
func (server *Server) Prepare() {
	logger.BootstrapLogger.Infoln("Connecting to the database")
	server.initialize()

	logger.BootstrapLogger.Infoln("Seeding the database")
	seeder.Initialize(server.DB)

	logger.BootstrapLogger.Infoln("Initializing the routers")
	api.Initialize(server.DB, server.Router)

	server.Router.Use(middleware.Cors)

	server.prepared = true
}

func (server *Server) Run() {
	if !server.prepared {
		server.Prepare()
	}

	logger.BootstrapLogger.Infoln("Server running", config.ServerConfiguration.Application.Name)
	listenAddress := fmt.Sprintf(":%d", config.ServerConfiguration.Http.PortNumber)
	logger.BootstrapLogger.Infoln(http.ListenAndServe(listenAddress, server.Router))
}
