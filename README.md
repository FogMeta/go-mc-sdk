# go-mc-sdk

[![Made by FogMeta](https://img.shields.io/badge/made%20by-FogMeta-green.svg)](https://en.fogmeta.com/)
[![Twitter Follow](https://img.shields.io/twitter/follow/FogMeta)](https://twitter.com/FogMeta)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

A Golang SDK for the MC([Meta Client](https://github.com/FogMeta/meta-client)) service, providing an easy interface for developers to deal with the Meta-Client API. It streamlines the process of securely storing, retrieving and recovering data on the IPFS and Filecoin network. 

Meta-Client is a Web3 data service that securely stores data backups and enables data recovery. It automatically records data storage information and stores data on both the IPFS gateway and Filecoin network, providing fast retrieval and permanent backup.

## Features

`go-mc-sdk` provides the following features:

- Upload files or directories to the IPFS gateway
- Report data information to the Meta-Client server 
    - Meta-Client will automatically complete data processing(split or merge file and generate CAR files)
    - Store the CAR file in the IPFS gateway
    - Send CAR files to the storage providers in the Filecoin network
- Download files or directories to the local machine
- Query the dataset list and details by the dataset name
- Get a source file information by the IPFS CID
- Acquire the storage information and status of the dataset.

## Prerequisites

Before using `go-mc-sdk`, you need to install the following services:

- Aria2 service

```
sudo apt install aria2

```
- [IPFS service](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions)
- [Go](https://golang.org/dl/) (1.16 or later)

## Installation

To install `go-mc-sdk`, run the following command:

```
go get github.com/FogMeta/go-mc-sdk
```


## Usage

### [Initialization](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#newapiclient)

First, you need to create a Meta Client object, which can be initialized as follows:

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
### [Upload files or directories](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#uploadfile) 
To upload files or directories to the IPFS gateway and Filecoin network, you can use the following method:

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

    // upload files in testdata to the IPFS server
    apiUrl := "http://127.0.0.1:5001"
    inputPath := "./testdata"
    ipfsCid, err := metaClient.UploadFile(apiUrl, inputPath)
    if err != nil {
        logs.GetLogger().Error("upload failed:", err)
        return
    }
    logs.GetLogger().Infoln("upload successful, IPFS CID: ", ipfsCid)

    return
}
```

### [Report the data information](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#reportmetaclientserver)
To report data information to the Meta Client server, you can use the following method:

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

    info, err := os.Stat(sourceName)
    if err != nil {
        logs.GetLogger().Error("Failed to get the file information, error:", err)
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
        logs.GetLogger().Error("failed to report the dataset info to the server, error:", err)
        return
    }
    logs.GetLogger().Infoln("the dataset has been successfully reported to the server., dataset name:", datasetName)

    return
}
```

### [Download Files or Directories](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#downloadfile)
To download files or directories from the IPFS gateway and Filecoin network, you can use the following method:

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

    // download the file from IPFS gateway to your local
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    outPath := "./output"
    downloadUrl := "http://127.0.0.1:8080/ipfs/QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    host := "127.0.0.1"
    port := 6800
    secret := "my_aria2_secret"
    conf := &metacli.Aria2Conf{Host: host, Port: port, Secret: secret}
    err := metaClient.DownloadFile(ipfsCid, outPath, downloadUrl, conf)
    if err != nil {
        logs.GetLogger().Error("failed to download the file", err)
        return
    }
    logs.GetLogger().Infoln("the file has been downloaded successfully")
    
    return
}
```

### [Get the dataset list by the dataset name](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#getdatasetlist)
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

    // get the dataset list from meta server
    datasetName := "dataset-name"
    pageNum := 0
    size := 10
    datasetListPager, err := metaClient.GetDatasetList(datasetName, pageNum, size)
    if err != nil {
        logs.GetLogger().Error("failed to get the dataset list:", err)
        return
    }
    logs.GetLogger().Infof("get the dataset list successfully: %+v", datasetListPager)

    return
}

```

### [Get the source file information by IPFS CID](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#getsourcefileinfo)
To get the dataset file information by IPFS CID, you can use the following method:

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
    logs.GetLogger().Infof("get source file information successfully: %+v", ipfsDataDetail)

    return
}
```

### [Get source file status](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#getsourcefilestatus)
To get the source file status by dataset name and IPFS CID, you can use the following method:

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
        logs.GetLogger().Error("failed to get the status of source file:", err)
        return
    }
    logs.GetLogger().Infof("get source file status successfully: %+v", sourceFileStatusPager)

    return
}
```


### [Get IPFS CID Information](https://github.com/FogMeta/go-mc-sdk/blob/main/document/api.md#getipfscidstat)
To get file or directory information of IPFS CID, you can use the following method:

```
package main

import (
    metacli "github.com/FogMeta/go-mc-sdk/client"
    "github.com/filswan/go-swan-lib/logs"
)

func main() {
    apiUrl := "http://127.0.0.1:5001"
    ipfsCid := "QmQgM2tGEduvYmgYy54jZaZ9D7qtsNETcog8EHR8XoeyEp"
    info, err := metacli.GetIpfsCidInfo(apiUrl, ipfsCid)
    if err != nil {
        logs.GetLogger().Error("Failed to get IPFS CID information:", err)
        return
    }
    logs.GetLogger().Infof("get information successfully, IPFS CID: %s, Data Size:%d, Is Directory:%t.", info.IpfsCid, info.DataSize, info.IsDirectory)
}
```

## API Documentation

For detailed API lists, please check out the [API Documentation](document/api.md ':include').

## Contributing

Contributions to go-mc-sdk are welcome! If you find any errors or want to add new features, please submit an [Issue](https://github.com/FogMeta/go-mc-sdk/issues), or initiate a [Pull Request](https://github.com/FogMeta/go-mc-sdk/pulls).

## License

go-mc-sdk is licensed under the MIT License.
