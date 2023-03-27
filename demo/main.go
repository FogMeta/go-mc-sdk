package main

import (
	"fmt"
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/meta-client-sdk/client"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	DefaultConf    *sdk.ClientConf
	KeyFlag        cli.StringFlag
	TokenFlag      cli.StringFlag
	ApiUrlFlag     cli.StringFlag
	GatewayUrlFlag cli.StringFlag
	MetaUrlFlag    cli.StringFlag
)

func init() {
	DefaultConf = &sdk.ClientConf{
		Key:            "",
		Token:          "",
		IpfsApiUrl:     "http://127.0.0.1:5001",
		IpfsGatewayUrl: "http://127.0.0.1:8080",
		MetaServerUrl:  "http://127.0.0.1:1234",
	}

	KeyFlag = cli.StringFlag{
		Name:    "key",
		Aliases: []string{"k"},
		Usage:   "key from meta swan",
		Value:   DefaultConf.Key,
	}

	TokenFlag = cli.StringFlag{
		Name:    "token",
		Aliases: []string{"t"},
		Usage:   "token from meta swan",
		Value:   DefaultConf.Token,
	}

	ApiUrlFlag = cli.StringFlag{
		Name:  "api-url",
		Usage: "url of IPFS api server",
		Value: DefaultConf.IpfsApiUrl,
	}

	GatewayUrlFlag = cli.StringFlag{
		Name:  "gateway-url",
		Usage: "url of IPFS gateway",
		Value: DefaultConf.IpfsGatewayUrl,
	}

	MetaUrlFlag = cli.StringFlag{
		Name:  "meta-url",
		Usage: "url of meta server",
		Value: DefaultConf.MetaServerUrl,
	}
}

func main() {
	app := &cli.App{
		Name:  "client-sdk-demo",
		Usage: "Utility for working with meta client sdk",
		Commands: []*cli.Command{
			{
				Name:   "upload",
				Usage:  "upload file or dir to ipfs server",
				Action: UploadDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "file or directory which will upload to IPFS server.",
						Required: true,
					},
					&KeyFlag,
					&TokenFlag,
					&ApiUrlFlag,
					&GatewayUrlFlag,
					&MetaUrlFlag,
				},
			},
			{
				Name:   "download",
				Usage:  "download file or dir from ipfs server",
				Action: DownloadDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "data-cid",
						Usage:    "data cid which will be downloaded",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "out-path",
						Usage:    "directory where files will be downloaded to.",
						Required: true,
					},
					&KeyFlag,
					&TokenFlag,
					&ApiUrlFlag,
					&GatewayUrlFlag,
					&MetaUrlFlag,
					&cli.BoolFlag{
						Name:  "aria2",
						Usage: "directory where files will be downloaded to.",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "host",
						Usage: "aria2 server address.",
						Value: "127.0.0.1",
					},
					&cli.IntFlag{
						Name:  "port",
						Usage: "aria2 server port.",
						Value: 6800,
					},
					&cli.StringFlag{
						Name:  "secret",
						Usage: "directory where files will be downloaded to.",
						Value: "",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func UploadDemo(c *cli.Context) error {
	key := c.String("key")
	token := c.String("token")
	apiUrl := c.String("api-url")
	gatewayUrl := c.String("gateway-url")
	metaUrl := c.String("meta-url")

	metaClient := sdk.NewAPIClient(key, token, apiUrl, gatewayUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	input := c.String("input")
	dataCid, err := metaClient.UploadFile(input)
	if err != nil {
		logs.GetLogger().Error("upload dir error:", err)
		return err
	}
	logs.GetLogger().Infoln("upload dir success, and data cid: ", dataCid)

	return nil
}

func DownloadDemo(c *cli.Context) error {
	key := c.String("key")
	token := c.String("token")
	apiUrl := c.String("api-url")
	gatewayUrl := c.String("gateway-url")
	metaUrl := c.String("meta-url")

	metaClient := sdk.NewAPIClient(key, token, apiUrl, gatewayUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	dataCid := c.String("data-cid")
	outPath := c.String("out-path")

	var conf *sdk.Aria2Conf

	if c.Bool("aria2") {
		host := c.String("host")
		port := c.Int("port")
		secret := c.String("secret")
		conf = &sdk.Aria2Conf{Host: host, Port: port, Secret: secret}
	}
	err := metaClient.DownloadFile(dataCid, outPath, conf)
	if err != nil {
		logs.GetLogger().Error("download dir error:", err)
		return err
	}
	logs.GetLogger().Infoln("download dir success")

	return nil
}

func UpDownDirDemo(c *sdk.MetaClient) {
	dataCid, err := c.UploadFile("./testdata")
	if err != nil {
		logs.GetLogger().Error("upload dir error:", err)
		return
	}
	logs.GetLogger().Infoln("upload dir success, and data cid: ", dataCid)

	err = c.DownloadFile(dataCid, "./output", nil)
	if err != nil {
		logs.GetLogger().Error("download dir error:", err)
		return
	}
	logs.GetLogger().Infoln("download dir success")
}

func UpDownFileDemo(c *sdk.MetaClient) {
	dataCid, err := c.UploadFile("./testdata")
	if err != nil {
		logs.GetLogger().Error("upload file error:", err)
		return
	}
	logs.GetLogger().Infoln("upload file success, and data cid: ", dataCid)

	err = c.DownloadFile(dataCid, "./output", nil)
	if err != nil {
		logs.GetLogger().Error("download file error:", err)
		return
	}
	logs.GetLogger().Infoln("download file success")
}

func Aria2DownFileDemo(c *sdk.MetaClient) {
	dataCid, err := c.UploadFile("./testdata/help")
	if err != nil {
		logs.GetLogger().Error("upload file error:", err)
		return
	}
	logs.GetLogger().Infoln("upload file success, and data cid: ", dataCid)

	conf := &sdk.Aria2Conf{Host: "127.0.0.1", Port: 6800, Secret: "secret123"}
	err = c.DownloadFile(dataCid, "output", conf)
	if err != nil {
		logs.GetLogger().Error("download file error:", err)
		return
	}
	logs.GetLogger().Infoln("download file by aria2 success")
}

func Aria2DownDirDemo(c *sdk.MetaClient) {
	dataCid, err := c.UploadFile("./testdata")
	if err != nil {
		logs.GetLogger().Error("upload dir error:", err)
		return
	}
	logs.GetLogger().Infoln("upload dir success, and data cid: ", dataCid)

	conf := &sdk.Aria2Conf{Host: "127.0.0.1", Port: 6800, Secret: "secret123"}

	err = c.DownloadFile(dataCid, "output", conf)
	if err != nil {
		logs.GetLogger().Error("download dir error:", err)
		return
	}
	logs.GetLogger().Infoln("download dir by aria2 success")
}
