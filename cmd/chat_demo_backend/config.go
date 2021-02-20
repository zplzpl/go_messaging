package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigParser struct {
	Server struct {
		Addr         string
		ReadTimeout  int64
		WriteTimeout int64
	}
	RabbitMQ RabbitMQCfg
}

type RabbitMQCfg struct {
	DialUrl     string
	MsgExchange struct {
		Name        string
		Type        string
		Durable     bool
		AutoDeleted bool
		Internal    bool
		NoWait      bool
	}
	MsgQueueDeclare struct {
		Name             string
		Durable          bool
		DeleteWhenUnused bool
		Exclusive        bool
		NoWait           bool
	}
	MsgQueueBind struct {
		RoutingKey string
		Exchange   string
		NoWait     bool
	}
	MsgConsume struct {
		Consumer  string
		AutoAck   bool
		Exclusive bool
		NoLocal   bool
		NoWait    bool
	}
	MsgRespExchange struct {
		Name        string
		Type        string
		Durable     bool
		AutoDeleted bool
		Internal    bool
		NoWait      bool
	}
	MsgRespPublish struct {
		Exchange   string
		RoutingKey string
		Mandatory  bool
		Immediate  bool
	}
}

func WatchConfig(changeConfig chan struct{}) error {

	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(AppPath, "config"))
	viper.SetConfigName("app")
	viper.SetConfigType("toml")

	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		changeConfig <- struct{}{}
	})

	return nil
}

func LoadConfig() (*ConfigParser, error) {

	var config ConfigParser
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
