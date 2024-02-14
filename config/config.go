package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Service struct {
	Name string
}

type App struct {
	Name     string   `mapstructure:"name"`
	Version  string   `mapstructure:"version"`
	Log      Log      `mapstructure:"log"`
	BasePath string   `mapstructure:"basepath"`
	Port     string   `mapstructure:"port"`
	Kafka    Kafka    `mapstructure:"kafka"`
	Database Database `mapstructure:"database"`
	TZ       string   `mapstructure:"timezone"`
	Redis    Redis    `mapstructure:"redis"`
	Gateway  Gateway  `mapstructure:"gateway"`
}

type Database struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	UserName    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	SSL         bool   `mapstructure:"ssl"`
	RootCert    string `mapstructure:"rootca"`
	PrivateKey  string `mapstructure:"key"`
	Certificate string `mapstructure:"cert"`
	Timeout     int    `mapstructure:"timeout"`
	Schema      string `mapstructure:"schema"`
}

type Kafka struct {
	Producer Producer `mapstructure:"producer"`
	Consumer Consumer `mapstructure:"consumer"`
	Addr     []string `mapstructure:"addr"`
}

type Redis struct {
	Addr     string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Producer struct {
	TopicID string `mapstructure:"topicId"`
}

type Consumer struct {
	TopicIDs []string `mapstructure:"topicIds"`
	GroupID  string   `mapstructure:"groupId"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Name   string `mapstructure:"name"`
}

type Gateway struct {
	Endpoint string `mapstructure:"endpoint"`
}

func New(name string) (*App, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetConfigFile(name)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	configuration := new(App)
	if err := viper.Unmarshal(configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
