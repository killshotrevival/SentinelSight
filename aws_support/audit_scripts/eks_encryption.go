package auditscripts

import (
	"fmt"
	"os"
	"sentinelsight/support"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	log "github.com/sirupsen/logrus"
)

// Main handler function for Rout53 Privacy Protection Check
func StartEKSEncryptionCheck(sentinelConfig *support.SentinelConfig, sess *session.Session, newLog *log.Entry) {

	newLog = newLog.WithFields(log.Fields{
		"audit": "eks_encryption",
	})
	newLog.Info("Starting encryption checkup for EKS cluster")
	for i := 0; i < len(sentinelConfig.Region); i++ {
		newLog.Infof("Working with region -> %s", sentinelConfig.Region[i])
		startEKSEncryptionCheckPerRegion(sentinelConfig, sess, newLog, sentinelConfig.Region[i])

	}
}

// region Based privacy protection check executor
func startEKSEncryptionCheckPerRegion(sentinelConfig *support.SentinelConfig, sess *session.Session, newLog *log.Entry, region string) {

	var listOfCluster []string
	svc := eks.New(sess, &aws.Config{
		Region: aws.String(region)})

	listResp, listErr := svc.ListClusters(&eks.ListClustersInput{})

	if listErr != nil {
		newLog.Errorf("Error occurred while listing EKS clusters -> %s", listErr.Error())
		return
	}

	if len(listResp.Clusters) < 1 {
		newLog.Info("No cluster found in the region, existing")
		return
	}

	var clusterName string
	var clusterArray []string
	var descResp *eks.DescribeClusterOutput
	var descErr error
	var encryptionStatus []*string
	for i := 0; i < len(listResp.Clusters); i++ {

		clusterName = *listResp.Clusters[i]
		newLog.Infof("Working with cluster -> %s", clusterName)
		clusterArray = strings.Split(clusterName, "/")

		clusterName = clusterArray[len(clusterArray)-1]

		descResp, descErr = svc.DescribeCluster(&eks.DescribeClusterInput{Name: &clusterName})

		if descErr != nil {
			newLog.Errorf("Error occurred while fetching cluster information -> %s", descErr.Error())
			return
		}

		encryptionStatus = descResp.Cluster.EncryptionConfig[0].Resources

		if len(encryptionStatus) > 0 {
			newLog.Info("Secrets Encryption: Enabled")
		} else {
			newLog.Error("Secrets Encryption: Disabled")
			listOfCluster = append(listOfCluster, clusterName)
		}
	}

	if len(listOfCluster) > 0 {
		newLog.Info("Some cluster found that needs investigation")
		dataForFile := "# Following cluster needs more investigation\n\n"
		for i := 0; i < len(listOfCluster); i++ {
			dataForFile = fmt.Sprintf("%s\n%s", dataForFile, listOfCluster[i])
		}

		err := os.WriteFile(fmt.Sprintf("%s/%s-%s", sentinelConfig.OutputDir, region, "eks_encryption.txt"), []byte(dataForFile), 0644)
		if err != nil {
			newLog.Panic("Error occurred while writing data to output file.")
		}

		newLog.Info("Data written to output file successfully.")

	} else {
		newLog.Info("Nothing to write in output file.")
	}

}
