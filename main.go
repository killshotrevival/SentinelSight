package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	awssupport "sentinelsight/aws_support"
	"sentinelsight/support"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	t := time.Now().UTC()

	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.json", "Path to the config file")

	flag.Parse()

	file, _ := os.Open(configFilePath)
	decoder := json.NewDecoder(file)
	configuration := support.SentinelConfig{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panicf("Error occurred while reading config -> %s", err.Error())
		log.Panicf("Please create `config.json` file in proper format")
		log.Exit(1)
	}
	file.Close()

	log.Info("Validating configuration loaded")
	err = support.ValidateSentinelConfig(&configuration)
	if err != nil {
		log.Panicf("Invalid configuration found -> %s", err.Error())
		log.Exit(1)
	}
	log.Info("Configuration loaded successfully")

	log.Info("Checking output directory")
	timeString := t.Format("2006-01-02-15-04")
	log.Info("Starting the adit")

	if _, err := os.Stat(configuration.OutputDir); os.IsNotExist(err) {
		err = os.MkdirAll(configuration.OutputDir, 0777)
		if err != nil {
			log.Panicf("Error occurred while creating output directory -> %s", err.Error())
			log.Exit(1)
		}
	} else {
		log.Info("Output directory found")
	}

	configuration.OutputDir = fmt.Sprintf("%s/%s", configuration.OutputDir, timeString)

	log.Infof("Creating timestamp directory -> %s", configuration.OutputDir)
	err = os.MkdirAll(configuration.OutputDir, 0777)
	if err != nil {
		log.Panicf("Error occurred while creating timestamp directory -> %s", err.Error())
		log.Exit(1)
	}

	err = awssupport.StartAwsAudit(&configuration)
	if err != nil {
		log.Panicf("Error occurred while aws audit -> %s", err.Error())
		log.Exit(1)
	}

	log.Info("Audits executed successfully")
	log.Infof("Please check %s directory for further investigation", configuration.OutputDir)

}
