package main

import (
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/meta-client-sdk/client"
)

func main() {
	key := ""
	token := ""
	ipfsUploadUrl := ""
	ipfsDownloadUrl := ""
	metaUrl := ""

	metaClient := sdk.NewAPIClient(key, token, ipfsUploadUrl, ipfsDownloadUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return
	}

	dataCid, err := metaClient.UploadFile("./testdata")
	if err != nil {
		logs.GetLogger().Error("upload file error:", err)
		return
	}
	logs.GetLogger().Infoln("upload file success, and data cid: ", dataCid)

	err = metaClient.DownloadFile(dataCid, "./output", nil)
	if err != nil {
		logs.GetLogger().Error("download file error:", err)
		return
	}
	logs.GetLogger().Infoln("download file success")

	logs.GetLogger().Infoln("Demo Test Over ...")
}
