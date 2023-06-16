package main

import (
	"encoding/json"
	"os"
	awssupport "sentinelsight/aws_support"
	"sentinelsight/support"

	log "github.com/sirupsen/logrus"
)

func main() {

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

	log.Info("Starting the adit")

	err = awssupport.StartAwsAudit(&configuration)
	if err != nil {
		log.Errorf("Error occurred while aws audit -> %s", err.Error())
		log.Exit(1)
	}

	log.Info("Audit executed successfully")
}
