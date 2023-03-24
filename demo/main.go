package main

import (
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/meta-client-sdk/client"
)

func main() {
	key := ""
	token := ""
	ipfsApiUrl := "http://127.0.0.1:5001"
	gatewatUrl := "http://127.0.0.1:8080"
	metaUrl := ""

	metaClient := sdk.NewAPIClient(key, token, ipfsApiUrl, gatewatUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return
	}

	// testUpDownDir(metaClient)

	// testUpDownFile(metaClient)

	testAria2File(metaClient)

}

func testUpDownDir(c *sdk.MetaClient) {
	dataCid, err := c.UploadFile("./testdata/about")
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

func testUpDownFile(c *sdk.MetaClient) {
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

func testAria2File(c *sdk.MetaClient) {
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
