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
		Key:            "V0schjjl_bxCtSNwBYXXXX",
		Token:          "fca72014744019a949248874610fXXXX",
		IpfsApiUrl:     "http://192.168.2.42:5001",
		IpfsGatewayUrl: "http://192.168.2.42:8080",
		MetaServerUrl:  "http://192.168.2.34:8099",
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
			{
				Name:   "notify",
				Usage:  "notify to meta server",
				Action: Notify2MetaDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "file or directory which will upload to IPFS server.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "data-cid",
						Usage:    "data cid which will be downloaded",
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
				Name:   "list",
				Usage:  "get files list from meta server",
				Action: GetFilesListDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "page-num",
						Usage:    "",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "limit",
						Usage:    "",
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
				Name:   "datacid",
				Usage:  "get data cid from meta server",
				Action: GetDataCidByNameDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "file or directory name which will be query from meta server.",
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
				Name:   "info",
				Usage:  "get detail info from  meta server",
				Action: GetInfoByDataCidDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "data-cid",
						Usage:    "data cid which will be query from meta server.",
						Required: true,
					},
					&KeyFlag,
					&TokenFlag,
					&ApiUrlFlag,
					&GatewayUrlFlag,
					&MetaUrlFlag,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func buildClient(c *cli.Context) *sdk.MetaClient {
	key := c.String("key")
	token := c.String("token")
	apiUrl := c.String("api-url")
	gatewayUrl := c.String("gateway-url")
	metaUrl := c.String("meta-url")

	metaClient := sdk.NewAPIClient(key, token, apiUrl, gatewayUrl, metaUrl)

	return metaClient
}

func UploadDemo(c *cli.Context) error {
	//
	metaClient := buildClient(c)
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
	metaClient := buildClient(c)
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

func Notify2MetaDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	input := c.String("input")
	dataCid := c.String("data-cid")
	err := metaClient.NotifyMetaServer(input, dataCid)
	if err != nil {
		logs.GetLogger().Error("notify data cid to meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("notify data cid to meta server success")

	return nil
}

func GetFilesListDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	pageNum := c.Int("page-num")
	limit := c.Int("limit")
	_, err := metaClient.GetFileLists(pageNum, limit, true)
	if err != nil {
		logs.GetLogger().Error("get files list from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get files list from meta server success")

	return nil
}

func GetDataCidByNameDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	name := c.String("name")
	_, err := metaClient.GetDataCIDByName(name)
	if err != nil {
		logs.GetLogger().Error("get data cid from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get data cid from meta server success")

	return nil
}

func GetInfoByDataCidDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	dataCid := c.String("data-cid")
	_, err := metaClient.GetFileInfoByDataCid(dataCid)
	if err != nil {
		logs.GetLogger().Error("get detail info from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get detail info from meta server success")

	return nil
}
