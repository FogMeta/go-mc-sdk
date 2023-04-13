# go-mc-sdk

[![Made by FogMeta](https://img.shields.io/badge/made%20by-FogMeta-green.svg)](https://en.fogmeta.com/)
[![Twitter Follow](https://img.shields.io/twitter/follow/FogMeta)](https://twitter.com/FogMeta)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

A Golang SDK for the MC(Meta Client) service, providing an easy interface for developers to deal with the Meta-Client API. It streamlines the process of securely storing, retrieving, and recovering data on the IPFS and Filecoin network. 

Meta-Client is a Web3 data service that securely stores data backups and enables data recovery. It automatically records data storage information and stores data on both IPFS gateway and Filecoin network, providing fast retrieval and permanent backup.

## Features

go-mc-sdk provides the following features:

- Upload files or directory to the IPFS gateway
- Report data information to the Meta-Client server 
    - Meta-Client will automatically complete data processing(split or merge file and generate CAR files)
    - Store the CAR file in the IPFS gateway
    - Send CAR files to the storage providers in the Filecoin network
- Download files or directory to the local machine
- Query dataset list and details by the specified dataset name
- Get list of source file information by the specified IPFS cid
- Query storage information and status of a dataset

## Prerequisites

Before using go-mc-sdk, you need to install the following services:

- Aria2 service

```
sudo apt install aria2

```
- [IPFS service](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions)
- [Go](https://golang.org/dl/) (1.16 or later)

## Installation

To install go-mc-sdk, run the following command:

```
go get github.com/FogMeta/go-mc-sdk
```


## Usage

### Initialization

First, you need to create a MetaClient object, which can be initialized as follows:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
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
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings".
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings".
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

    return
}
```

### Report Data-related Information
To report data-related information to the Meta-Client server, you can use the following method:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)

    // report ipfs cid to meta server
    apiUrl := "http://127.0.0.1:5001"
    inputPath := "./testdata"
	
    datasetName := "dataset-name"
    ipfsGateway := "http://127.0.0.1:8080"
    sourceName := inputPath
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"

    info, err := metacli.GetIpfsCidInfo(apiUrl, ipfsCid)
    if err != nil {
        logs.GetLogger().Error("get ipfs cid stat information error:", err)
        return
    }
    oneItem := metacli.IpfsData{}
    oneItem.SourceName = sourceName
    oneItem.IpfsCid = ipfsCid
    oneItem.DataSize = info.DataSize
    oneItem.IsDirectory = info.IsDirectory
    oneItem.DownloadUrl = metacli.PathJoin(ipfsGateway, "ipfs/", ipfsCid)
    ipfsData := []metacli.IpfsData{oneItem}
    err = metaClient.ReportMetaClientServer(datasetName, ipfsData)
    if err != nil {
        logs.GetLogger().Error("report meta client server  failed:", err)
        return
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
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)

    // download file(s) from IPFS server
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    outPath := "./output"
    downloadUrl := "http://127.0.0.1:8080/ipfs/QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    host := "127.0.0.1"
    port := 6800
    secret := "my_aria2_secret"
    conf := &metacli.Aria2Conf{Host: host, Port: port, Secret: secret}
    err := metaClient.DownloadFile(ipfsCid, outPath, downloadUrl, conf)
    if err != nil {
        logs.GetLogger().Error("download failed:", err)
        return
    }
    logs.GetLogger().Infoln("download success")
    
    return
}
```

### Get Dataset List by Dataset Name
To get the dataset list by dataset name, you can use the following method:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)

    // get dataset list from meta server
    datasetName := "dataset-name"
    pageNum := 0
    size := 10
    datasetListPager, err := metaClient.GetDatasetList(datasetName, pageNum, size)
    if err != nil {
        logs.GetLogger().Error("get dataset list failed:", err)
        return
    }
    logs.GetLogger().Infof("get dataset list success: %+v", datasetListPager)

    return
}

```

### Get Source File Information by IPFS Cid
To get the dataset file information by IPFS cid, you can use the following method:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)
	
    // get source file information
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    ipfsDataDetail, err := metaClient.GetSourceFileInfo(ipfsCid)
    if err != nil {
        logs.GetLogger().Error("get source file information failed:", err)
        return
    }
    logs.GetLogger().Infof("get source file information success: %+v", ipfsDataDetail)

    return
}
```

### Get Source File Status by Dataset Name
To get the source file status by dataset name, you can use the following method:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    // Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
    key := "V0schjjl_bxCtSNwBYXXXX"
    // Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings".
    token := "fca72014744019a949248874610fXXXX"
    metaUrl := "http://{ip}:8099/rpc/v0"
    metaClient := metacli.NewAPIClient(key, token, metaUrl)

    // get source file status
    datasetName := "dataset-name"
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    pageNum := 0
    size := 10
    sourceFileStatusPager, err := metaClient.GetSourceFileStatus(datasetName, ipfsCid, pageNum, size)
    if err != nil {
        logs.GetLogger().Error("get source file status failed:", err)
        return
    }
    logs.GetLogger().Infof("get source file status success: %+v", sourceFileStatusPager)

    return
}
```

## API Documentation

For detailed API lists, please check out the [API Documentation](document/api.md ':include').

## Contributing

Contributions to go-mc-sdk are welcome! If you find any errors or want to add new features, please submit an [Issue](https://github.com/FogMeta/meta-client-sdk/issues), or initiate a [Pull Request](https://github.com/FogMeta/meta-client-sdk/pulls).

## License

go-mc-sdk is licensed under the MIT License.
