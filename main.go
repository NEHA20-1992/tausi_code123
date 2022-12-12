package main

import (
	l "github.com/NEHA20-1992/tausi_code/pkg/logger"
	"github.com/NEHA20-1992/tausi_code/pkg/server"
)

func main() {
	l.BootstrapLogger.Infoln("Tausi Server starting")
	server.ApiServer.Run()
}
