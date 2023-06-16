package awssupport

import (
	"sentinelsight/support"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"
)

// Main starter function for all aws audit functions
func StartAwsAudit(sentinelConfig *support.SentinelConfig) error {
	newLog := log.WithFields(log.Fields{
		"filName": "awssupport",
	})

	newLog.Info("Creating session object")

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(sentinelConfig.AccessId, sentinelConfig.SecretKey, ""),
	})
	if err != nil {
		return err
	}
	newLog.Infof("Session Created successfully %s", *sess.Config.Region)

	allAwsAuditFunctions, err := InitializeAuditFunctions()

	if err != nil {
		newLog.Infof("Error occurred while generating aws audit functions list -> %s", err.Error())
		return err
	}

	newLog.Infof("Audits list generated successfully, will be running %d audits on the target", len(allAwsAuditFunctions))

	var wg = sync.WaitGroup{}
	maxGoroutines := sentinelConfig.MaxThreads
	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < len(allAwsAuditFunctions); i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(n func(*support.SentinelConfig, *session.Session, *log.Entry)) {
			n(sentinelConfig, sess, newLog)
			<-guard
			wg.Done()
		}(allAwsAuditFunctions[i])
	}

	wg.Wait()

	newLog.Info("All AWS audits completed successfully :)")
	return nil
}
