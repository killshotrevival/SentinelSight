package awssupport

import (
	auditscripts "sentinelsight/aws_support/audit_scripts"
	"sentinelsight/support"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"
)

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

	// Route53 Privacy Protection Check
	auditscripts.StartRoute53PrivacyProtectionCheck(sentinelConfig, sess, newLog)

	return nil
}
