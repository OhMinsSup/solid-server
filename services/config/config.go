package config

import (
	"github.com/spf13/viper"
)

const (
	DefaultServerRoot = "http://localhost:8080"
	DefaultPort       = 8080
)

// Configuration json 파일에 저장된 앱 구성
type Configuration struct {
	ServerRoot string `json:"serverRoot" mapstructure:"serverRoot"`
	Port       int    `json:"port" mapstructure:"port"`

	// Db
	DBType         string `json:"dbType" mapstructure:"dbType"`
	DBConfigString string `json:"dbconfig" mapstructure:"dbconfig"`
	DBTablePrefix  string `json:"dbtableprefix" mapstructure:"dbtableprefix"`

	// Web
	UseSSL       bool   `json:"useSSL" mapstructure:"useSSL"`
	SecureCookie bool   `json:"secureCookie" mapstructure:"secureCookie"`
	WebPath      string `json:"webpath" mapstructure:"webpath"`

	// File
	FilesDriver string `json:"filesdriver" mapstructure:"filesdriver"`
	FilesPath   string `json:"filespath" mapstructure:"filespath"`
	MaxFileSize int64  `json:"maxfilesize" mapstructure:"mafilesize"`

	Telemetry         bool     `json:"telemetry" mapstructure:"telemetry"`
	TelemetryID       string   `json:"telemetryid" mapstructure:"telemetryid"`
	PrometheusAddress string   `json:"prometheusaddress" mapstructure:"prometheusaddress"`
	WebhookUpdate     []string `json:"webhook_update" mapstructure:"webhook_update"`

	// Session
	Secret                   string            `json:"secret" mapstructure:"secret"`
	SessionExpireTime        int64             `json:"session_expire_time" mapstructure:"session_expire_time"`
	SessionRefreshTime       int64             `json:"session_refresh_time" mapstructure:"session_refresh_time"`
	SessionSecretKey         string            `json:"session_secret_key" mapstructure:"session_secret_key"`
	LocalOnly                bool              `json:"localonly" mapstructure:"localonly"`
	EnableLocalMode          bool              `json:"enableLocalMode" mapstructure:"enableLocalMode"`
	LocalModeSocketLocation  string            `json:"localModeSocketLocation" mapstructure:"localModeSocketLocation"`
	EnablePublicSharedBoards bool              `json:"enablePublicSharedBoards" mapstructure:"enablePublicSharedBoards"`
	FeatureFlags             map[string]string `json:"featureFlags" mapstructure:"featureFlags"`

	AuthMode string `json:"authMode" mapstructure:"authMode"`

	LoggingCfgFile string `json:"logging_cfg_file" mapstructure:"logging_cfg_file"`
	LoggingCfgJSON string `json:"logging_cfg_json" mapstructure:"logging_cfg_json"`

	AuditCfgFile string `json:"audit_cfg_file" mapstructure:"audit_cfg_file"`
	AuditCfgJSON string `json:"audit_cfg_json" mapstructure:"audit_cfg_json"`
}

// ReadConfigFile 은 파일 시스템에서 구성을 읽습니다.
func ReadConfigFile(configFilePath string) (*Configuration, error) {
	if configFilePath == "" {
		viper.SetConfigFile("./config.json")
	} else {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetEnvPrefix("solid")
	viper.AutomaticEnv()
	viper.SetDefault("ServerRoot", DefaultServerRoot)
	viper.SetDefault("Port", DefaultPort)
	viper.SetDefault("DBType", "sqlite3")
	viper.SetDefault("DBConfigString", "./solid.db")
	viper.SetDefault("DBTablePrefix", "")
	viper.SetDefault("SecureCookie", false)
	viper.SetDefault("WebPath", "./pack")
	viper.SetDefault("FilesPath", "./files")
	viper.SetDefault("FilesDriver", "local")
	viper.SetDefault("Telemetry", true)
	viper.SetDefault("TelemetryID", "")
	viper.SetDefault("WebhookUpdate", nil)
	viper.SetDefault("SessionExpireTime",  60*60*24*30) // 30 days session lifetime
	viper.SetDefault("SessionRefreshTime", 60*60*5)    // 5 minutes session refresh
	viper.SetDefault("SessionSecretKey", "4qUMElHAy6J7ZD9nJCQTbhxIwbAyy7vp")
	viper.SetDefault("LocalOnly", false)
	viper.SetDefault("EnableLocalMode", false)
	viper.SetDefault("LocalModeSocketLocation", "/var/tmp/solid_local.socket")
	viper.SetDefault("EnablePublicSharedBoards", false)
	viper.SetDefault("FeatureFlags", map[string]string{})
	viper.SetDefault("AuthMode", "native")
	viper.SetDefault("NotifyFreqCardSeconds", 120)    // 2 minutes after last card edit
	viper.SetDefault("NotifyFreqBoardSeconds", 86400) // 1 day after last card edit
	viper.SetDefault("PrometheusAddress", "")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}

	configuration := Configuration{}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func removeSecurityData(config Configuration) Configuration {
	clean := config
	return clean
}
