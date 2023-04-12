package main

import (
	"fmt"
	sdk "github.com/FogMeta/meta-client-sdk/client"
	"github.com/FogMeta/meta-client-sdk/demo/config"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	Conf           *sdk.ClientConf
	KeyFlag        cli.StringFlag
	TokenFlag      cli.StringFlag
	ApiUrlFlag     cli.StringFlag
	GatewayUrlFlag cli.StringFlag
	MetaUrlFlag    cli.StringFlag
)

func init() {

	Conf = config.GetConfig("./client.toml")

	KeyFlag = cli.StringFlag{
		Name:    "key",
		Aliases: []string{"k"},
		Usage:   "key from meta swan",
		Value:   Conf.Key,
	}

	TokenFlag = cli.StringFlag{
		Name:    "token",
		Aliases: []string{"t"},
		Usage:   "token from meta swan",
		Value:   Conf.Token,
	}

	ApiUrlFlag = cli.StringFlag{
		Name:  "api-url",
		Usage: "url of IPFS api server",
		Value: Conf.IpfsApiUrl,
	}

	GatewayUrlFlag = cli.StringFlag{
		Name:  "gateway-url",
		Usage: "url of IPFS gateway",
		Value: Conf.IpfsGatewayUrl,
	}

	MetaUrlFlag = cli.StringFlag{
		Name:  "meta-url",
		Usage: "url of meta server",
		Value: Conf.MetaServerUrl,
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
						Value: Conf.Aria2.Host,
					},
					&cli.IntFlag{
						Name:  "port",
						Usage: "aria2 server port.",
						Value: Conf.Aria2.Port,
					},
					&cli.StringFlag{
						Name:  "secret",
						Usage: "directory where files will be downloaded to.",
						Value: Conf.Aria2.Secret,
					},
				},
			},
			{
				Name:   "report",
				Usage:  "report to meta server",
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
					&cli.BoolFlag{
						Name:  "show-car",
						Usage: "",
						Value: true,
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
				Action: GetIpfsCidByNameDemo,
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
				Action: GetInfoByIpfsCidDemo,
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
	metaUrl := c.String("meta-url")

	metaClient := sdk.NewAPIClient(key, token, metaUrl)

	return metaClient
}

func UploadDemo(c *cli.Context) error {
	//
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	input := c.String("input")
	apiUrl := c.String("api-url")
	ipfsCid, err := metaClient.UploadFile(apiUrl, input)
	if err != nil {
		logs.GetLogger().Error("upload error:", err)
		return err
	}
	logs.GetLogger().Infoln("upload success, and data cid: ", ipfsCid)

	return nil
}

func DownloadDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	ipfsCid := c.String("data-cid")
	outPath := c.String("out-path")

	var conf *sdk.Aria2Conf

	if c.Bool("aria2") {
		host := c.String("host")
		port := c.Int("port")
		secret := c.String("secret")
		conf = &sdk.Aria2Conf{Host: host, Port: port, Secret: secret}
	}
	err := metaClient.DownloadFile(ipfsCid, outPath, "", conf)
	if err != nil {
		logs.GetLogger().Error("download error:", err)
		return err
	}
	logs.GetLogger().Infoln("download success")

	return nil
}

func Notify2MetaDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	input := c.String("input")
	ipfsCid := c.String("data-cid")
	gatewayUrl := c.String("gateway-url")

	sourceFile := sdk.SourceFileReq{}
	sourceFile.SourceName = ""

	info, err := os.Stat(input)
	if err != nil {
		logs.GetLogger().Error("get input stat information error:", err)
		return err
	}
	isDir := info.IsDir()
	dataSize := info.Size()
	downloadUrl := sdk.PathJoin(gatewayUrl, "ipfs/", ipfsCid)
	sourceFile.DataList = append(sourceFile.DataList, sdk.DataItem{IsDirector: isDir, DataSize: dataSize, IpfsCid: ipfsCid, DownloadUrl: downloadUrl})

	err = metaClient.ReportMetaClientServer(sourceFile)
	if err != nil {
		logs.GetLogger().Error("report data cid to meta client server error:", err)
		return err
	}
	logs.GetLogger().Infoln("report data cid to meta client server success")

	return nil
}

func GetFilesListDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	pageNum := c.Int("page-num")
	limit := c.Int("limit")
	showCar := c.Bool("show-car")
	_, err := metaClient.GetFileLists(pageNum, limit, sdk.WithShowCar(showCar))
	if err != nil {
		logs.GetLogger().Error("get files list from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get files list from meta server success")

	return nil
}

func GetIpfsCidByNameDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	name := c.String("name")
	_, err := metaClient.GetIpfsCidByName(name)
	if err != nil {
		logs.GetLogger().Error("get data cid from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get data cid from meta server success")

	return nil
}

func GetInfoByIpfsCidDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
	}

	ipfsCid := c.String("data-cid")
	_, err := metaClient.GetFileInfoByIpfsCid(ipfsCid)
	if err != nil {
		logs.GetLogger().Error("get detail info from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get detail info from meta server success")

	return nil
}
