package auditscripts

import (
	"sentinelsight/support"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53domains"
	log "github.com/sirupsen/logrus"
)

func StartRoute53PrivacyProtectionCheck(awsKeys *support.SentinelConfig, sess *session.Session, newLog *log.Entry) {

	newLog.Info("Starting Privacy protection lookup for rout53 entries")

	svc := route53domains.New(sess, &aws.Config{
		Region: aws.String("us-east-1")})

	res, listErr := svc.ListDomains(&route53domains.ListDomainsInput{})
	if listErr != nil {
		newLog.Errorf("Error occurred while fetching domains list -> %s", listErr.Error())
		return
	}

	newLog.Info("Rout53 domains listed successfully")

	var getRes *route53domains.GetDomainDetailOutput
	var getErr error
	for i := 0; i < len(res.Domains); i++ {
		newLog.Infof("Domain found -> %s", *res.Domains[i].DomainName)

		getRes, getErr = svc.GetDomainDetail(&route53domains.GetDomainDetailInput{DomainName: res.Domains[i].DomainName})
		if getErr != nil {
			newLog.Errorf("Error occurred while finding data for -> %s", getErr.Error())
		}

		if *getRes.AdminPrivacy {
			newLog.Infof("Privacy Protection: Enabled")
		} else {
			newLog.Error("Privacy Protection: Disabled")
		}

	}
}
