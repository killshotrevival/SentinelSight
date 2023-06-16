package support

type SentinelConfig struct {
	AccessId  string   `json:accessId`
	SecretKey string   `json:secretKey`
	Region    []string `json:region`
	OutputDir string   `json:outputDir`
}
