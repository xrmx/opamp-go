package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"google.golang.org/protobuf/proto"

	"github.com/open-telemetry/opamp-go/internal/examples/server/data"
	"github.com/open-telemetry/opamp-go/internal/examples/server/opampsrv"
	"github.com/open-telemetry/opamp-go/internal/examples/server/uisrv"
	"github.com/open-telemetry/opamp-go/protobufs"
)

var logger = log.New(log.Default().Writer(), "[MAIN] ", log.Default().Flags()|log.Lmsgprefix|log.Lmicroseconds)

var (
	configFile = flag.String("config", "", "Path to binary protobuf file containing default AgentConfigMap")
)

func main() {
	flag.Parse()

	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Load default config if provided
	var defaultConfig *protobufs.AgentConfigMap
	if *configFile != "" {
		logger.Printf("Loading default config from: %s", *configFile)
		data, err := os.ReadFile(*configFile)
		if err != nil {
			logger.Fatalf("Failed to read config file: %v", err)
		}

		defaultConfig = &protobufs.AgentConfigMap{}
		if err := proto.Unmarshal(data, defaultConfig); err != nil {
			logger.Fatalf("Failed to unmarshal config data: %v", err)
		}
		logger.Printf("Default config loaded successfully")
	}

	logger.Println("OpAMP Server starting...")

	uisrv.Start(curDir)
	opampSrv := opampsrv.NewServer(&data.AllAgents, defaultConfig)
	opampSrv.Start()

	logger.Println("OpAMP Server running...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	logger.Println("OpAMP Server shutting down...")
	uisrv.Shutdown()
	opampSrv.Stop()
}
