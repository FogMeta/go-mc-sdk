package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	apiUrl := "http://127.0.0.1:5001"
	inputPath := "./testdata"
	dataCid, err := metaClient.UploadFile(apiUrl, inputPath)
	if err != nil {
		logs.GetLogger().Error("upload failed:", err)
	}
	logs.GetLogger().Infoln("upload success, and data cid is: ", dataCid)

	inputPath = "./testdata"
	dataCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	ipfsGateway := "http://127.0.0.1:8080"
	err = metaClient.ReportMetaClientServer(inputPath, dataCid, ipfsGateway)
	if err != nil {
		logs.GetLogger().Error("report meta client server  failed:", err)
	}
	logs.GetLogger().Infoln("report meta client server success")

	dataCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	outPath := "./output"
	downloadUrl := "http://127.0.0.1:8080/ipfs/QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	host := "127.0.0.1"
	port := 6800
	secret := "my_aria2_secret"
	conf := &metacli.Aria2Conf{Host: host, Port: port, Secret: secret}
	err = metaClient.DownloadFile(dataCid, outPath, downloadUrl, conf)
	if err != nil {
		logs.GetLogger().Error("download failed:", err)
	}
	logs.GetLogger().Infoln("download success")

	fileName := "testdata"
	dataCids, err := metaClient.GetDataCIDByName(fileName)
	if err != nil {
		logs.GetLogger().Error("get data cid failed:", err)
	}
	logs.GetLogger().Infof("get data cid success: %+v", dataCids)

	pageNum := 0
	limit := 10
	showStorage := true
	sourceFileList, err := metaClient.GetFileLists(pageNum, limit, showStorage)
	if err != nil {
		logs.GetLogger().Error("get file list failed:", err)
	}
	logs.GetLogger().Infof("get file list success: %+v", sourceFileList)

	dataCid = "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	sourceFileInfo, err := metaClient.GetFileInfoByDataCid(dataCid)
	if err != nil {
		logs.GetLogger().Error("get source file info failed:", err)
	}
	logs.GetLogger().Infof("get source file info success: %+v", sourceFileInfo)

	return
}
