# meta-client-sdk

* [NewAPIClient](#NewAPIClient)
* [UploadFile](#UploadFile)
* [NotifyMetaServer](#NotifyMetaServer)
* [DownloadFile](#DownloadFile)
* [GetFileLists](#GetFileLists)
* [GetDataCIDByName](#GetDataCIDByName)
* [GetFileInfoByDataCid](#GetFileInfoByDataCid)
* [Usage](#Usage)

## NewAPIClient

Definition:
创建 Meta Client实例，通过这个示例调用相关API

```shell
func NewAPIClient(key, token, ipfsApiUrl, ipfsGatewayUrl, metaUrl string) *MetaClient
```

Inputs:

```shell
key                    # 从Filswan获取的api key
token                  # 从Filswan获取的access token
ipfsApiUrl             # IPFS服务的API地址
ipfsGatewayUrl         # IPFS服务的网关地址
metaUrl                # Meta Server的地址
```


Outputs:

```shell
*MetaClient            # 创建的Meta Client实例
```

## UploadFile

Definition:
上传文件或文件夹到IPFS服务

```shell
func (m *MetaClient) UploadFile(inputPath string) (dataCid string, err error) 
```

Inputs:

```shell
inputPath              # 需要上传到IPFS服务器的文件路径

```

Outputs:

```shell
dataCid                # IPFS服务器返回的Data Cid
error                  # error or nil
```

## NotifyMetaServer

Definition:
通知Meta Server，文件或文件夹已上传到IPFS服务。

```shell
func (m *MetaClient) NotifyMetaServer(sourceName string, dataCid string) error 
```

Inputs:

```shell
sourceName             # 已经上传到IPFS服务器的文件名.
dataCid                # IPFS服务器返回的Data Cid
```

Outputs:

```shell
error                  # error or nil
```


## DownloadFile

Definition:
从IPFS下载指定Data Cid对应的文件或文件夹，如果配置aria2选项，则使用aria2工具下载，否则使用IPFS API下载

```shell
func (m *MetaClient) DownloadFile(dataCid string, outPath string, conf *Aria2Conf) error
```

Inputs:

```shell
dataCid                # 需要下载的Data CID
outPath                # 下载文件输出路径
conf                   # aria2 相关配置项，包括:  
                       # host   Aria2 server address
                       # port   Aria2 server port
                       # secret Must be the same value as rpc-secure in aria2 conf

```

Outputs:

```shell
error                  # error or nil
```


## GetFileLists

Definition:
根据指定页码和每页条数，从Meta Server获取文件列表

```shell
func (m *MetaClient) GetFileLists(pageNum int, limit int, showStorageInfo bool) ([]FileDetails, error)
```

Inputs:

```shell
pageNum                # 查询第几页
limit                  # 每页记录的条数
showStorageInfo        # 是否返回存储信息
```

Outputs:

```shell
[]FileDetails          # 返回文件描述列表
error                  # error or nil
```


## GetDataCIDByName

Definition:
从Meta Server获取指定文件或文件夹对应的Data Cid

```shell
func (m *MetaClient) GetDataCIDByName(fileName string) ([]string, error) 
```

Inputs:

```shell
fileName               # 查询文件名
```

Outputs:

```shell
[]string               # 查询文件名对应的Data CID列表
error                  # error or nil
```


## GetFileInfoByDataCid

Definition:
从Meta Server获取Data Cid对应的文件或文件详细信息

```shell
func (m *MetaClient) GetFileInfoByDataCid(dataCid string) (*FileDetails, error)
```

Inputs:

```shell
dataCid                # 查询Dada Cid
```

Outputs:

```shell
*FileDetails           # 返回Data CID对应的信息
error                  # error or nil
```

## Usage

Install

```go
go get github.com/filswan/go-swan-lib/meta-client-sdk@latest
```

Demo

```go
package main

import (
	"fmt"
	"github.com/filswan/go-swan-lib/logs"
	sdk "github.com/meta-client-sdk/client"
	"os"
)
func main() {
	key := "XXXXXXXX"
	token := "XXXXXXXX"
	ipfsApiUrl := "http://127.0.0.1:5001"
	gatewatUrl := "http://127.0.0.1:8080"
	metaUrl := "http://127.0.0.1:1234"

	targetName := "./testdata"
	outPath := "output"

	metaClient := sdk.NewAPIClient(key, token, ipfsApiUrl, gatewatUrl, metaUrl)
	if metaClient == nil {
		logs.GetLogger().Error("create meta client failed, please check the input parameters")
		return
	}

	dataCid, err := metaClient.UploadFile(targetName)
	if err != nil {
		logs.GetLogger().Error("upload dir error:", err)
		return
	}
	logs.GetLogger().Infoln("upload dir success, and data cid: ", dataCid)

	err = metaClient.NotifyMetaServer(targetName, dataCid)
	if err != nil {
		logs.GetLogger().Error("notify meta server error:", err)
		return
	}
	logs.GetLogger().Infoln("notify meta server success")

	//create aria2 config
	conf := &sdk.Aria2Conf{Host: "127.0.0.1", Port: 6800, Secret: "secret123"}

	err = metaClient.DownloadFile(dataCid, outPath, conf)
	if err != nil {
		logs.GetLogger().Error("download dir error:", err)
		return
	}
	logs.GetLogger().Infoln("download dir by aria2 success")

}

```