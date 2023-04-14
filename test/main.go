package main

import (
	metacli "github.com/FogMeta/go-mc-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
	"os"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	// update file(s) in testdata to IPFS server
	apiUrl := "http://127.0.0.1:5001"
	inputPath := "./testdata"
	ipfsCid, err := metaClient.UploadFile(apiUrl, inputPath)
	if err != nil {
		logs.GetLogger().Error("upload failed:", err)
		return
	}
	logs.GetLogger().Infoln("upload success, and ipfs cid is: ", ipfsCid)

	// report ipfs cid to meta server
	datasetName := "dataset-name"
	ipfsGateway := "http://127.0.0.1:8080"
	sourceName := inputPath
	ipfsCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"

	info, err := os.Stat(sourceName)
	if err != nil {
		logs.GetLogger().Error("get ipfs cid stat information error:", err)
		return
	}
	oneItem := metacli.IpfsData{}
	oneItem.SourceName = sourceName
	oneItem.IpfsCid = ipfsCid
	oneItem.DataSize = info.Size()
	oneItem.IsDirectory = info.IsDir()
	oneItem.DownloadUrl = metacli.PathJoin(ipfsGateway, "ipfs/", ipfsCid)
	ipfsData := []metacli.IpfsData{oneItem}
	err = metaClient.ReportMetaClientServer(datasetName, ipfsData)
	if err != nil {
		logs.GetLogger().Error("report meta client server  failed:", err)
		return
	}
	logs.GetLogger().Infoln("report meta client server success")

	// download file(s) from IPFS server
	ipfsCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	outPath := "./output"
	downloadUrl := "http://127.0.0.1:8080/ipfs/QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	host := "127.0.0.1"
	port := 6800
	secret := "my_aria2_secret"
	conf := &metacli.Aria2Conf{Host: host, Port: port, Secret: secret}
	err = metaClient.DownloadFile(ipfsCid, outPath, downloadUrl, conf)
	if err != nil {
		logs.GetLogger().Error("download failed:", err)
		return
	}
	logs.GetLogger().Infoln("download success")

	// get dataset list from meta server
	datasetName = "dataset-name"
	pageNum := 0
	size := 10
	datasetListPager, err := metaClient.GetDatasetList(datasetName, pageNum, size)
	if err != nil {
		logs.GetLogger().Error("get dataset list failed:", err)
		return
	}
	logs.GetLogger().Infof("get dataset list success: %+v", datasetListPager)

	// get source file information
	ipfsCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	ipfsDataDetail, err := metaClient.GetSourceFileInfo(ipfsCid)
	if err != nil {
		logs.GetLogger().Error("get source file information failed:", err)
		return
	}
	logs.GetLogger().Infof("get source file information success: %+v", ipfsDataDetail)

	// get source file status
	datasetName = "dataset-name"
	ipfsCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	pageNum = 0
	size = 10
	sourceFileStatusPager, err := metaClient.GetSourceFileStatus(datasetName, ipfsCid, pageNum, size)
	if err != nil {
		logs.GetLogger().Error("get source file status failed:", err)
		return
	}
	logs.GetLogger().Infof("get source file status success: %+v", sourceFileStatusPager)

	return
}
