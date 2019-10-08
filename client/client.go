package client

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

var instance *rpcclient.Client

func GetInstance() *rpcclient.Client {
	if instance != nil {
		return instance
	}

	var err error
	conCfg := loadConfig
	instance, err = rpcclient.New(conCfg, nil, nil)
	if err != nil {
		log.Fatal(err)
		instance.Shutdown()
	}
	return instance
}

// load config file from config/app.yml with viper
// the config file should contain the correct RPC details
func loadConfig() *rpcclient.ConnConfig {

	connCfg := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:5222",
		User:         "via",
		Pass:         "via",
		HTTPPostMode: true, // Viacoin core only supports HTTP POST mode
		DisableTLS:   true, // Viacoin core does not provide TLS by default
	}

	return connCfg
}
