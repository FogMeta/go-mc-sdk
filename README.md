# meta-client-sdk

* [Install Aria2 Service](#Install-Aria2-Service)
* [Install IPFS Server](#Install-IPFS-Server)
* [Usage](#Usage)

## Install Aria2 Service
```shell
sudo apt install aria2
```

## Install IPFS Server
###
[Install IPFS binary](https://docs.ipfs.tech/install/command-line/#linux)
### 
[Start the IPFS daemon](https://docs.ipfs.tech/how-to/kubo-basic-cli/#install-kubo)

## Usage

Install

```go
go get github.com/filswan/go-swan-lib/meta-client-sdk@latest
```

Quick Start

```go
package main

import (
	"fmt"
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/FogMeta/meta-client-sdk/client"
	"os"
)

func main() {
	// Initialize parameters
	
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	key        := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
	token      := "fca72014744019a949248874610fXXXX"
	
	ipfsApiUrl := "http://127.0.0.1:5001"    // IPFS API address
	gatewayUrl := "http://127.0.0.1:8080"    // Gateway address
	metaUrl := "http://127.0.0.1:8099"       // Meta server address
	
	targetName := "./testdata"               // Target file or folder
	outPath    := "./output"                 // Directory to save downloaded files

	// Create Meta Client object based on actual parameters
	metaClient := sdk.NewAPIClient(key, token, ipfsApiUrl, gatewayUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters") // Fail to create object
		return
	}

	// Upload the target file or folder
	dataCid, err := metaClient.UploadFile(targetName)
	if err != nil {
		logs.GetLogger().Error("upload dir error:", err) // Fail to upload
		return
	}
	logs.GetLogger().Infoln("upload dir success, and data cid: ", dataCid) // Upload successful

	// Notify Meta Server that the file has been uploaded to IPFS
	err = metaClient.NotifyMetaServer(targetName, dataCid)
	if err != nil {
		logs.GetLogger().Error("notify meta server error:", err) // Fail to notify Meta Server
		return
	}
	logs.GetLogger().Infoln("notify meta server success") // Successfully notified Meta Server

	// Create Aria2 configuration options
	conf := &sdk.Aria2Conf{Host: "127.0.0.1", Port: 6800, Secret: "my_aria2_secret"}

	// Download the file or folder corresponding to the specified Data Cid to the specified directory
	err = metaClient.DownloadFile(dataCid, outPath, conf)
	if err != nil {
		logs.GetLogger().Error("download dir error:", err) // Fail to download
		return
	}
	logs.GetLogger().Infoln("download dir by aria2 success") // Download successful

}

```