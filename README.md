# Meta-Client-SDK

Meta-Client is a Web3 data service that securely stores data backups and enables data recovery. It automatically records data storage information and stores data on both IPFS gateway and Filecoin network, providing fast retrieval and permanent backup.

With Meta-Client-SDK, users can easily access and recover their data while ensuring its safety. The automated recording of data storage information also helps users optimize their storage strategies.
## Features

Meta-Client-SDK provides the following features:

- Upload files or directory to the IPFS gateway
- Report data information to the Meta-Client server 
    - Meta-Client will automatically complete data processing(split or merge file and generate CAR files)
    - Store the CAR file in the IPFS gateway
    - Send CAR files to the storage providers in the Filecoin network
- Download files or directory to the local machine
- Query DataCID for a file by the file name
- Get a list of all files of current user
- Query storage information and status of a single file or directory

## Prerequisites

Before using Meta-Client-SDK, you need to install the following services:

- Aria2 service

```
sudo apt install aria2
```
- [IPFS service](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions)
- [Go](https://golang.org/dl/) (1.16 or later)

## Installation

To install Meta-Client-SDK, run the following command:

```
go get github.com/FogMeta/meta-client-sdk
```


## Usage

### Initialization

First, you need to create a MetaClient object, which can be initialized as follows:

```
package main

import (
  metacli "github.com/FogMeta/meta-client-sdk/client"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)
}
```
### Upload Files or Directory
To upload files or directory to IPFS gateway and Filecoin network, you can use the following method:
```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings".
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)

    apiUrl := "http://127.0.0.1:5001"
    inputPath := "./testdata" //file or directory path
    dataCid, err := metaClient.UploadFile(apiUrl, inputPath)
    if err != nil {
	logs.GetLogger().Error("upload failed:", err)
    }
    logs.GetLogger().Infoln("upload success, and data cid is: ", dataCid)
    return
}
```
### Report Data-related Information
To report data-related information to the Meta-Client server, you can use the following method:

```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	inputPath := "./testdata" // file or directory path that has been uploaded to the IPFS gateway
	dataCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	ipfsGateway := "http://127.0.0.1:8080"
	err := metaClient.ReportMetaClientServer(inputPath, dataCid, ipfsGateway)
	if err != nil {
		logs.GetLogger().Error("report meta client server  failed:", err)
	}
	logs.GetLogger().Infoln("report meta client server success")

	return
}
```
### Download Files or Directory
To download files or directory from the IPFS gateway and Filecoin network, you can use the following method:

```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	dataCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	outPath := "./output"
	downUrl := "http://127.0.0.1:8080/ipfs/QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	host := "127.0.0.1"
	port := 6800
	secret := "my_aria2_secret"
	conf := &metacli.Aria2Conf{Host: host, Port: port, Secret: secret}
	err := metaClient.DownloadFile(dataCid, outPath, downUrl, conf)
	if err != nil {
		logs.GetLogger().Error("download failed:", err)
	}
	logs.GetLogger().Infoln("download success")

	return
}
```

### Get DataCID for a File by file name
To get the DataCID for a file by its filename, you can use the following method:
```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	name := "./testdata"
	dataCids, err := metaClient.GetDataCIDByName(name)
	if err != nil {
		logs.GetLogger().Error("get data cid failed:", err)
	}
	logs.GetLogger().Infof("get data cid success: %+v", dataCids)

	return
}

```
### View List of Files and Storage Status
To view a list of all files under the current user, or to query storage information and the status of a single file or folder, you can use the following method:
```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	page := 0
	limit := 10
	showStorage := true
	sourceFileList, err := metaClient.GetFileLists(page, limit, showStorage)
	if err != nil {
		logs.GetLogger().Error("get file list failed:", err)
	}
	logs.GetLogger().Infof("get file list success: %+v", sourceFileList)

	return
}
```


### Get Source File Information by its Filename
Todo
```
package main

import (
	metacli "github.com/FogMeta/meta-client-sdk/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	// Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
	key := "V0schjjl_bxCtSNwBYXXXX"
	// Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings".
	token := "fca72014744019a949248874610fXXXX"
	metaUrl := "http://{ip}:8099/rpc/v0"
	metaClient := metacli.NewAPIClient(key, token, metaUrl)

	dataCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
	sourceFileInfo, err := metaClient.GetFileInfoByDataCid(dataCid)
	if err != nil {
		logs.GetLogger().Error("get source file info failed:", err)
	}
	logs.GetLogger().Infof("get source file info success: %+v", sourceFileInfo)

	return
}
```

## API Documentation

For detailed API lists, please check out the [API Documentation](document/api.md ':include').

## Contributing

Contributions to Meta-Client-SDK are welcome! If you find any errors or want to add new features, please submit an [Issue](https://github.com/FogMeta/meta-client-sdk/issues), or initiate a [Pull Request](https://github.com/FogMeta/meta-client-sdk/pulls).

## License

Meta-Client-SDK is licensed under the MIT License.
