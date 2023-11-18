package entities

import (
	"github.com/littleurl/api/internal/timestamp"
)

type EngineConfig struct {
	EngineID string `json:"engine_id" dynamodbav:"engine_id"`
	SyncMode string `json:"sync_mode" dynamodbav:"sync_mode"`

	// deployment
	EnableAutoDeploy bool                 `json:"enable_auto_deploy"         dynamodbav:"enable_auto_deploy"`
	TargetVersion    *string              `json:"target_version,omitempty"   dynamodbav:"target_version,omitempty"`
	DeployedVersion  *string              `json:"deployed_version,omitempty" dynamodbav:"deployed_version,omitempty"`
	DeployedAt       *timestamp.Timestamp `json:"deployed_at,omitempty"      dynamodbav:"deployed_at,omitempty"`
	TestedAt         *timestamp.Timestamp `json:"tested_at,omitempty"        dynamodbav:"tested_at,omitempty"`

	// redirection
	ErrorPassthrough  bool `json:"error_passthrough"  dynamodbav:"error_passthrough"`
	OpaqueRedirection bool `json:"opaque_redirection" dynamodbav:"opaque_redirection"`

	// analytics
	AnalyticsProvider string `json:"analytics_provider" dynamodbav:"analytics_provider"`
	WaitForAnalytics  bool   `json:"wait_for_analytics" dynamodbav:"wait_for_analytics"`
}

const (
	AnalyticsProviderGoogle  = "google"
	AnalyticsProviderWebhook = "webhook"
	AnalyticsProviderMatomo  = "matomo"
)

const (
	EngineProviderCloudflareWorkers = "cloudflare-workers"
	EngineProviderAmazonS3          = "aws-s3"
	EngineProviderStandalone        = "standalone"
)

type Engine struct {
	Features EngineFeatures `json:"features"  dynamodbav:"features"`
	ID       string         `json:"id"        dynamodbav:"id"`
	Name     string         `json:"name"      dynamodbav:"name"`
	SyncPush bool           `json:"sync_push" dynamodbav:"sync_push"`
	SyncPull bool           `json:"sync_pull" dynamodbav:"sync_pull"`
	SyncCron bool           `json:"sync_cron" dynamodbav:"sync_cron"`
}

type EngineFeatures struct {
	OpaqueRedirection bool `json:"opaque_redirect" dynamodbav:"opaque_redirect"`
	Analytics         bool `json:"analytics"       dynamodbav:"analytics"`
	AccessLogging     bool `json:"access_logging"  dynamodbav:"access_logging"`
	AutoDeploy        bool `json:"auto_deploy"     dynamodbav:"auto_deploy"`
}
