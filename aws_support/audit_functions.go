package awssupport

import (
	auditscripts "sentinelsight/aws_support/audit_scripts"
	"sentinelsight/support"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

// This function can be used for generating a list of all audit functions to be executed
func InitializeAuditFunctions() ([]func(*support.SentinelConfig, *session.Session, *log.Entry), error) {
	var allAwsAuditFunctions []func(*support.SentinelConfig, *session.Session, *log.Entry)

	allAwsAuditFunctions = append(allAwsAuditFunctions,
		auditscripts.StartRoute53PrivacyProtectionCheck)

	return allAwsAuditFunctions, nil
}
