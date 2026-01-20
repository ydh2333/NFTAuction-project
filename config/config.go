package config

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server     ServerConfig
	MySQL      MySQLConfig
	Blockchain BlockchainConfig
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	RPCEndpoint     string // 区块链节点RPC地址
	ContractAddress string // 已部署的拍卖合约地址
	PrivateKey      string // 后端操作合约的私钥
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 从当前目录查找配置文件
	viper.AutomaticEnv()     // 优先读取环境变量

	// 环境变量映射
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("mysql.dsn", "MYSQL_DSN")
	viper.BindEnv("blockchain.rpcendpoint", "RPC_ENDPOINT")

	// 加载配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Err(err).Msg("未找到配置文件，使用默认配置和环境变量")
	}

	// 初始化默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.readTimeout", 10*time.Second)
	viper.SetDefault("server.writeTimeout", 10*time.Second)
	viper.SetDefault("mysql.maxOpenConns", 100)
	viper.SetDefault("mysql.maxIdleConns", 20)
	viper.SetDefault("mysql.connMaxLifetime", 30*time.Minute)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("配置解析失败")
		os.Exit(1)
	}

	return &cfg
}
