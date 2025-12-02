package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config application configuration structure
type Config struct {
	// Network configuration
	Net string

	// Database configuration
	Database DatabaseConfig

	// Blockchain configuration
	Chain ChainConfig

	// Indexer configuration
	Indexer IndexerConfig

	// MetaApp configuration
	MetaApp MetaAppConfig

	// Metafs configuration
	Metafs MetafsConfig
}

// DatabaseConfig database configuration
type DatabaseConfig struct {
	IndexerType  string // Indexer database type: mysql, pebble
	Dsn          string // MySQL DSN
	MaxOpenConns int    // MySQL max open connections
	MaxIdleConns int    // MySQL max idle connections
	DataDir      string // PebbleDB data directory
}

// ChainConfig blockchain configuration
type ChainConfig struct {
	RpcUrl      string
	RpcUser     string
	RpcPass     string
	StartHeight int64
}

// StorageConfig storage configuration
type StorageConfig struct {
	Type  string
	Local LocalStorageConfig
	OSS   OSSStorageConfig
}

// LocalStorageConfig local storage configuration
type LocalStorageConfig struct {
	BasePath string
}

// OSSStorageConfig OSS storage configuration
type OSSStorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

// IndexerConfig indexer configuration
type IndexerConfig struct {
	Port               string // Indexer service port
	ScanInterval       int    // Scan interval in seconds
	BatchSize          int    // Batch size for processing
	StartHeight        int64  // Start block height
	MvcInitBlockHeight int64  // MVC chain initial block height to start scanning from
	BtcInitBlockHeight int64  // BTC chain initial block height to start scanning from
	SwaggerBaseUrl     string // Swagger API base URL
	ZmqEnabled         bool   // Enable ZMQ real-time monitoring
	ZmqAddress         string // ZMQ server address
}

// MetaAppConfig MetaApp configuration
type MetaAppConfig struct {
	DeployFilePath string // Deploy file path for MetaApp
}

// MetafsConfig Metafs service configuration
type MetafsConfig struct {
	Domain string // Metafs service domain (e.g., "https://file.metaid.io")
}

// UploaderConfig uploader configuration
type UploaderConfig struct {
	MaxFileSize    int64
	FeeRate        int64
	ChunkSize      int64
	SwaggerBaseUrl string // Swagger API base URL (e.g., "example.com:7282")
}

// RpcConfig RPC configuration
type RpcConfig struct {
	Url      string
	Username string
	Password string
}

// RpcConfigMap RPC configuration mapping (for multi-chain support)
var RpcConfigMap = map[string]RpcConfig{}

// Cfg global configuration instance
var Cfg *Config

// InitConfig initialize configuration
func InitConfig() error {
	viper.SetConfigFile(GetYaml())
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %s", err)
	}

	// Create configuration instance
	Cfg = &Config{
		Net: viper.GetString("net"),

		Database: DatabaseConfig{
			IndexerType:  viper.GetString("database.indexer_type"),
			Dsn:          viper.GetString("database.dsn"),
			MaxOpenConns: viper.GetInt("database.max_open_conns"),
			MaxIdleConns: viper.GetInt("database.max_idle_conns"),
			DataDir:      viper.GetString("database.data_dir"),
		},

		Chain: ChainConfig{
			RpcUrl:      viper.GetString("chain.rpc_url"),
			RpcUser:     viper.GetString("chain.rpc_user"),
			RpcPass:     viper.GetString("chain.rpc_pass"),
			StartHeight: viper.GetInt64("chain.start_height"),
		},

		Indexer: IndexerConfig{
			Port:               viper.GetString("indexer.port"),
			ScanInterval:       viper.GetInt("indexer.scan_interval"),
			BatchSize:          viper.GetInt("indexer.batch_size"),
			StartHeight:        viper.GetInt64("indexer.start_height"),
			MvcInitBlockHeight: viper.GetInt64("indexer.mvc_init_block_height"),
			BtcInitBlockHeight: viper.GetInt64("indexer.btc_init_block_height"),
			SwaggerBaseUrl:     viper.GetString("indexer.swagger_base_url"),
			ZmqEnabled:         viper.GetBool("indexer.zmq_enabled"),
			ZmqAddress:         viper.GetString("indexer.zmq_address"),
		},

		MetaApp: MetaAppConfig{
			DeployFilePath: viper.GetString("meta_app.deploy_file_path"),
		},

		Metafs: MetafsConfig{
			Domain: viper.GetString("metafs.domain"),
		},
	}

	// Set default values
	if Cfg.Indexer.Port == "" {
		Cfg.Indexer.Port = "7281"
	}
	if Cfg.Indexer.ScanInterval == 0 {
		Cfg.Indexer.ScanInterval = 10
	}
	if Cfg.Indexer.BatchSize == 0 {
		Cfg.Indexer.BatchSize = 100
	}
	if Cfg.Database.MaxOpenConns == 0 {
		Cfg.Database.MaxOpenConns = 100
	}
	if Cfg.Database.MaxIdleConns == 0 {
		Cfg.Database.MaxIdleConns = 10
	}
	if Cfg.Indexer.SwaggerBaseUrl == "" {
		Cfg.Indexer.SwaggerBaseUrl = "localhost:" + Cfg.Indexer.Port
	}
	if Cfg.MetaApp.DeployFilePath == "" {
		Cfg.MetaApp.DeployFilePath = "./deploy_data"
	}

	// Initialize RpcConfigMap (use currently configured chain)
	RpcConfigMap[Cfg.Net] = RpcConfig{
		Url:      Cfg.Chain.RpcUrl,
		Username: Cfg.Chain.RpcUser,
		Password: Cfg.Chain.RpcPass,
	}

	return nil
}
