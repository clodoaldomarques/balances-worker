package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	AppPort          int
	MySqlDBUser      string
	MySqlDBPass      string
	MySqlDBHost      string
	MySqlDBPort      string
	MysqlDBName      string
	AwsAddress       string
	AwsRegion        string
	AccessKeyID      string
	SecretAccessKey  string
	BalancesSNSTopic string
	BalancesSQSQueue string
}

type Option func(*Config)

var (
	singleton sync.Once
	instance  *Config
)

func New(options ...Option) *Config {
	singleton.Do(func() {
		instance = &Config{
			AppPort:          GetInt("APP_PORT", 5000),
			MySqlDBUser:      GetString("MYSQL_USER", ""),
			MySqlDBPass:      GetString("MYSQL_PASSWORD", ""),
			MySqlDBHost:      GetString("MYSQL_HOST", "192.168.49.2"),
			MySqlDBPort:      GetString("MYSQL_PORT", "30001"),
			MysqlDBName:      GetString("MYSQL_DATABASE", "balances"),
			AwsAddress:       GetString("AWS_ADDRESS", ""),
			AwsRegion:        GetString("AWS_REGION", ""),
			AccessKeyID:      GetString("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:  GetString("AWS_SECRET_ACCESS_KEY", ""),
			BalancesSNSTopic: GetString("BALANCES_SNS_TOPIC", ""),
			BalancesSQSQueue: GetString("BALANCES_SQS_QUEUE", ""),
		}
	})

	for _, optFunc := range options {
		optFunc(instance)
	}

	return instance
}

func WithAppPort(appPort int) Option {
	return func(c *Config) {
		c.AppPort = appPort
	}
}

func WithMySqlDBUser(mySqlDBUser string) Option {
	return func(c *Config) {
		c.MySqlDBUser = mySqlDBUser
	}
}

func WithMySqlDBPass(MySqlDBPass string) Option {
	return func(c *Config) {
		c.MySqlDBPass = MySqlDBPass
	}
}

func WithMySqlDBHost(MySqlDBHost string) Option {
	return func(c *Config) {
		c.MySqlDBHost = MySqlDBHost
	}
}
func WithMySqlDBPort(MySqlDBPort string) Option {
	return func(c *Config) {
		c.MySqlDBPort = MySqlDBPort
	}
}
func WithMysqlDBName(MysqlDBName string) Option {
	return func(c *Config) {
		c.MysqlDBName = MysqlDBName
	}
}
func WithAwsAddress(AwsAddress string) Option {
	return func(c *Config) {
		c.AwsAddress = AwsAddress
	}
}
func WithAwsRegion(AwsRegion string) Option {
	return func(c *Config) {
		c.AwsRegion = AwsRegion
	}
}

func WithSnsAccountTopic(SnsAccountTopic string) Option {
	return func(c *Config) {
		c.BalancesSNSTopic = SnsAccountTopic
	}
}

func (c Config) GetMySQLConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySqlDBUser,
		c.MySqlDBPass,
		c.MySqlDBHost,
		c.MySqlDBPort,
		c.MysqlDBName,
	)
}

func GetString(env string, def string) string {
	if e := os.Getenv(env); e != "" {
		return e
	}
	return def
}

func GetInt(env string, def int) int {
	i, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return def
	}
	return i
}
