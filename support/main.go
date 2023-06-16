package support

import "errors"

type SentinelConfig struct {
	AccessId   string   `json:"accessId"`
	SecretKey  string   `json:"secretKey"`
	Region     []string `json:"region"`
	OutputDir  string   `json:"outputDir"`
	MaxThreads int      `json:"maxThreads"`
}

// This function can be used for validating the Sentinel Config data loaded
func ValidateSentinelConfig(sentinelConfig *SentinelConfig) error {
	if sentinelConfig.AccessId == "" {
		return errors.New("invalid Access Id")
	}

	if sentinelConfig.SecretKey == "" {
		return errors.New("invalid Secret Key")
	}

	if len(sentinelConfig.Region) < 1 {
		return errors.New("no region provided")
	}

	return nil
}
