// Copyright (C) Quickplay Media. All rights reserved.
//
// All information contained herein is, and remains the property
// of Quickplay Media. The intellectual and technical concepts
// contained herein are proprietary to Quickplay Media.
//
// Dissemination of this information or reproduction of this
// material is strictly forbidden unless prior written
// permission is obtained from Quickplay Media.

// Package logger provides utility methods to initialize
// log handles
package logger

import (
	"os"
	"strings"

	"github.com/NEHA20-1992/tausi_code/pkg/config"
	"github.com/sirupsen/logrus"
)

// BootstrapLogger is the logger handle for logging service
// bootstrap events.
var BootstrapLogger *logrus.Logger

// AccessLogger is the logger handle for logging service
// access events.
var AccessLogger *logrus.Logger

// ServiceLogger is the logger handle for logging all service events.
var ServiceLogger *logrus.Logger

// init initializes all the logger handlers during service startup.
// This will be executed only once at the start of service.
func init() {

	// initBootstrapLogger initializes the logger handle for logging
	// service bootstrap events.
	initBootstrapLogger()

	// initAccessLogger initializes the logger handle for logging
	// service access events.
	initAccessLogger()

	// initLogger initializes the logger handle for logging
	// all service events.
	initServiceLogger()
}

func initLogger(logger *logrus.Logger, lfc config.LogFileConfiguration) {
	logger.Level = getLogLevel(lfc.LogLevel)
	if lfc.Path == "" {
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
		logger.SetReportCaller(true)
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		file, err := os.OpenFile(
			lfc.Path+
				strings.ToLower(lfc.Name),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
		if err != nil {
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(file)
		}
	}
}

// initBootstrapLogger initializes the logger handle for logging
// service bootstrap events.
func initBootstrapLogger() {
	lfc := config.ServerConfiguration.Logging.LogFile["bootstrap"]
	BootstrapLogger = logrus.New()

	initLogger(BootstrapLogger, lfc)
}

// initAccessLogger initializes the logger handle for logging
// service access events.
func initAccessLogger() {
	lfc := config.ServerConfiguration.Logging.LogFile["access"]
	AccessLogger = logrus.New()

	initLogger(AccessLogger, lfc)
}

// initLogger initializes the logger handle for logging
// all service events.
func initServiceLogger() {
	lfc := config.ServerConfiguration.Logging.LogFile["service"]
	ServiceLogger = logrus.New()

	initLogger(ServiceLogger, lfc)
}

// getLogLevel returns the log level to initialize the logger handler.
// Supported log levels
// 1. DEBUG
// 2. INFO
// 3. WARN
// 4. ERROR
// 5. FATAL
// 6. PANIC
func getLogLevel(loglevel string) logrus.Level {
	switch loglevel {
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.ErrorLevel
	}
}
