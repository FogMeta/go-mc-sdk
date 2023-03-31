package config

import (
	"github.com/BurntSushi/toml"
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/meta-client-sdk/client"
	"log"
	"os"
)

var config *sdk.ClientConf

func GetConfig(confFile string) *sdk.ClientConf {
	if config == nil {
		initConfig(confFile)
	}

	return config
}

func requiredFieldsAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"key"},
		{"token"},
		{"ipfs_api_url"},
		{"ipfs_gateway_url"},
		{"meta_server_url"},

		{"aria2"},
		{"aria2", "host"},
		{"aria2", "port"},
		{"aria2", "secret"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			log.Fatal("Required fields ", v)
		}
	}

	return true
}

func initConfig(confFile string) {

	_, err := os.Stat(confFile)
	if os.IsNotExist(err) {
		defaultConf := sdk.ClientConf{
			Key:            "V0schjjl_bxCtSNwBYXXXX",
			Token:          "fca72014744019a949248874610fXXXX",
			IpfsApiUrl:     "http://127.0.0.1:5001",
			IpfsGatewayUrl: "http://127.0.0.1:8080",
			MetaServerUrl:  "http://127.0.0.1:8099/rpc/v0",
			Aria2: sdk.Aria2Conf{
				Host:   "127.0.0.1",
				Port:   6800,
				Secret: "my_aria2_secret",
			},
		}
		f, err := os.Create(confFile)
		if err != nil {
			logs.GetLogger().Warn("create default config error:", err)
		}
		defer f.Close()

		if err := toml.NewEncoder(f).Encode(defaultConf); err != nil {
			logs.GetLogger().Warn("write default config error:", err)
		}
	}

	if metaData, err := toml.DecodeFile(confFile, &config); err != nil {
		logs.GetLogger().Fatal("read config error:", err)

	} else {
		if !requiredFieldsAreGiven(metaData) {
			logs.GetLogger().Fatal("Required fields not given")
		}
	}
}
