package config

import (
	"fmt"
	"github.com/spf13/viper"
	"nftshopping-store-api/pkg/flags"
)

var configInstance *Configuration

func GetConfig() (instance *Configuration, err error) {
	if configInstance == nil {
		instance, err = newConfig()
		if err != nil {
			return nil, err
		}
		configInstance = instance
	}
	return configInstance, nil
}

func newConfig() (*Configuration, error) {
	env := flags.Env
	viper.SetConfigName("config-" + env)

	// Set the path to look for the configurations files
	viper.AddConfigPath("./resources/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config files, %s", err)
		return nil, err
	}

	//// Set undefined variables
	config := &Configuration{}
	err := viper.Unmarshal(config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}
	return config, nil
}

type Configuration struct {
	Server   *Server
	Database *Database
	Logger   *Logger
	Amazon   *Amazon
	Casbin   *Casbin
	Item     *Item
}

type Server struct {
	Port int
}

type Database struct {
	Mongo *Mongo
	Redis *Redis
}

type Mongo struct {
	Uri    string
	Source string
}

type Redis struct {
	Uri      string
	Password string
}

type Logger struct {
	Level    string
	Mode     string
	Encoding string
}

type Amazon struct {
	S3 *S3
}

type S3 struct {
	Region     string
	BucketName string
	Endpoint   string
	AccessKey  string
	SecretKey  string
}

type Casbin struct {
	Model  string
	Policy string
}

type Item struct {
	Domain string
}
