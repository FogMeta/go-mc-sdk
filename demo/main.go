package main

import (
	"errors"
	"fmt"
	sdk "github.com/FogMeta/go-mc-sdk/client"
	"github.com/FogMeta/go-mc-sdk/demo/config"
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
		Name:  "key",
		Usage: "key from meta swan",
		Value: Conf.Key,
	}

	TokenFlag = cli.StringFlag{
		Name:  "token",
		Usage: "token from meta swan",
		Value: Conf.Token,
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
						Name:     "source",
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
				Action: Report2MetaServerDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dataset",
						Usage:    "Dataset name.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "source",
						Usage:    "Source name.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ipfs-cid",
						Usage:    "IPFS cid which will be reported",
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
				Usage:  "get dataset list from meta server",
				Action: GetDatasetListDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dataset",
						Usage:    "",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "page-num",
						Usage:    "",
						Value:    0,
						Required: true,
					},
					&cli.IntFlag{
						Name:     "size",
						Usage:    "",
						Value:    10,
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
				Usage:  "get source file info from meta server",
				Action: GetSourceFileInfoDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ipfs-cid",
						Usage:    "IPFS cid which will be query from meta server.",
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
				Name:   "status",
				Usage:  "get source file status from  meta server",
				Action: GetSourceFileStatusDemo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dataset",
						Usage:    "",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ipfs-cid",
						Usage:    "IPFS cid which will be query from meta server.",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "page-num",
						Usage:    "",
						Value:    0,
						Required: true,
					},
					&cli.IntFlag{
						Name:     "size",
						Usage:    "",
						Value:    10,
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
				Name:  "backup",
				Usage: "Generate CAR files from IPFS data to backup",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "task-name",
						Usage: "name of a task.",
					},
					&cli.StringFlag{
						Name:  "ipfs_repo",
						Usage: "directory which IPFS repo is.",
					},
					&cli.StringFlag{
						Name:  "out-dir",
						Usage: "directory where CAR file(s) will be generated.",
						Value: "/tmp/tasks",
					},
					&cli.IntFlag{
						Name:  "parallel",
						Usage: "number goroutines run when building ipld nodes",
						Value: 5,
					},
					&cli.Int64Flag{
						Name:  "slice-size",
						Usage: "bytes of each piece",
						Value: 17179869184,
					},
					&cli.Int64Flag{
						Name:  "dataset-size",
						Usage: "bytes of each dataset",
						Value: 171798691840,
					},
					&KeyFlag,
					&TokenFlag,
					&ApiUrlFlag,
					&GatewayUrlFlag,
					&MetaUrlFlag,
				},
				Action: BackupIpfsDataDemo,
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

	// logs.GetLogger().Debugf("buildClient MetaUrl:%s ApiKey:%s ApiToken:%s", key, token, metaUrl)
	metaClient := sdk.NewAPIClient(key, token, metaUrl)

	return metaClient
}

func UploadDemo(c *cli.Context) error {
	//
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}

	input := c.String("source")
	apiUrl := c.String("api-url")
	ipfsCid, err := metaClient.UploadFile(apiUrl, input)
	if err != nil {
		logs.GetLogger().Error("upload error:", err)
		return err
	}
	logs.GetLogger().Infoln("upload success, and ipfs cid: ", ipfsCid)

	return nil
}

func DownloadDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
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

func Report2MetaServerDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}

	dataset := c.String("dataset")
	sourceName := c.String("source")
	ipfsCid := c.String("ipfs-cid")
	gatewayUrl := c.String("gateway-url")
	apiUrl := c.String("api-url")

	info, err := sdk.GetIpfsCidInfo(apiUrl, ipfsCid)
	if err != nil {
		logs.GetLogger().Error("get ipfs cid stat information error:", err)
		return err
	}
	oneItem := sdk.IpfsData{}
	oneItem.IpfsCid = ipfsCid
	oneItem.SourceName = sourceName
	oneItem.DataSize = info.DataSize
	oneItem.IsDirectory = info.IsDirectory
	oneItem.DownloadUrl = sdk.PathJoin(gatewayUrl, "ipfs/", ipfsCid)
	ipfsData := []sdk.IpfsData{oneItem}

	err = metaClient.ReportMetaClientServer(dataset, ipfsData)
	if err != nil {
		logs.GetLogger().Error("report ipfs cid to meta client server error:", err)
		return err
	}

	logs.GetLogger().Infoln("report ipfs cid to meta client server success")

	return nil
}

func GetDatasetListDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}

	datasetName := c.String("dataset")
	pageNum := c.Int("page-num")
	size := c.Int("size")
	_, err := metaClient.GetDatasetList(datasetName, pageNum, size)
	if err != nil {
		logs.GetLogger().Error("get dataset list from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get dataset list from meta server success")

	return nil
}

func GetSourceFileInfoDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}

	ipfsCid := c.String("ipfs-cid")
	_, err := metaClient.GetSourceFileInfo(ipfsCid)
	if err != nil {
		logs.GetLogger().Error("get data cid from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get data cid from meta server success")

	return nil
}

func GetSourceFileStatusDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}
	datasetName := c.String("dataset")
	ipfsCid := c.String("ipfs-cid")
	pageNum := c.Int("page-num")
	size := c.Int("size")
	_, err := metaClient.GetSourceFileStatus(datasetName, ipfsCid, pageNum, size)
	if err != nil {
		logs.GetLogger().Error("get source file status from meta server error:", err)
		return err
	}
	logs.GetLogger().Infoln("get source file status from meta server success")

	return nil
}

func BackupIpfsDataDemo(c *cli.Context) error {
	metaClient := buildClient(c)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return errors.New("create meta client failed")
	}

	taskName := c.String("task-name")
	repoDir := c.String("ipfs_repo")
	outputDir := c.String("out-dir")
	carLimit := c.Int64("slice-size")
	dataSetLimit := c.Int64("dataset-size")
	parallel := c.Int("parallel")

	apiUrl := c.String("api-url")
	gatewayUrl := c.String("gateway-url")

	//
	err := metaClient.BackupIpfsData(taskName, repoDir, outputDir, apiUrl, gatewayUrl, dataSetLimit, carLimit, parallel)
	if err != nil {
		logs.GetLogger().Error("failed to  generate CAR for datastore:", err)
		return err
	}
	logs.GetLogger().Infoln("generate CAR for datastore successfully")

	return nil
}
