package main

import (
	"encoding/json"
	"fmt"
	"os"
	awssupport "sentinelsight/aws_support"
	"sentinelsight/support"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	t := time.Now().UTC()
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := support.SentinelConfig{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Errorf("Error occurred while reading config -> %s", err.Error())
		log.Exit(1)
	}
	file.Close()

	log.Info("Configuration loaded successfully")

	log.Info("Checking output directory")
	timeString := t.Format("2006-01-02-15-04")

	log.Info("Starting the adit")

	if _, err := os.Stat(configuration.OutputDir); os.IsNotExist(err) {
		err = os.MkdirAll(configuration.OutputDir, 0777)
		if err != nil {
			log.Errorf("Error occurred while creating output directory -> %s", err.Error())
			log.Exit(1)
		}
	} else {
		log.Info("Output directory found")
	}

	configuration.OutputDir = fmt.Sprintf("%s/%s", configuration.OutputDir, timeString)

	log.Infof("Creating timestamp directory -> %s", configuration.OutputDir)
	err = os.MkdirAll(configuration.OutputDir, 0777)
	if err != nil {
		log.Errorf("Error occurred while creating timestamp directory -> %s", err.Error())
		log.Exit(1)
	}

	err = awssupport.StartAwsAudit(&configuration)
	if err != nil {
		log.Errorf("Error occurred while aws audit -> %s", err.Error())
		log.Exit(1)
	}

	log.Info("Audit executed successfully")
}
