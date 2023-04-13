# go-mc-sdk

* [NewAPIClient](#NewAPIClient)
* [UploadFile](#UploadFile)
* [ReportMetaClientServer](#ReportMetaClientServer)
* [DownloadFile](#DownloadFile)
* [GetDatasetList](#GetDatasetList)
* [GetSourceFileInfo](#GetSourceFileInfo)
* [GetSourceFileStatus](#GetSourceFileStatus)

## NewAPIClient

Definition:
Creates a Meta Client instance to make relevant API calls.

```shell
func NewAPIClient(key, token, metaUrl string) *MetaClient
```

Inputs:

```shell
key                    # Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
token                  # Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". It can be ignored if `[sender].offline_swan=true`.
metaUrl                # Address of Meta Server.
```


Outputs:

```shell
*MetaClient            # Created Meta Client instance.
```

## UploadFile

Definition:
Uploads a file or directory to IPFS service.

```shell
func (m *MetaClient) UploadFile(ipfsApiUrl, inputPath string) (ipfsCid string, err error) 
```

Inputs:

```shell
ipfsApiUrl             # API address of IPFS service.
inputPath              # File or directory path to be uploaded to IPFS server
```

Outputs:

```shell
ipfsCid                # IPFS Cid returned by IPFS server
error                  # error or nil
```

## ReportMetaClientServer

Definition:
Report the Meta Client Server that the file or folder has been uploaded to the IPFS service.

```shell
ReportMetaClientServer(datasetName string, ipfsData []IpfsData) 
func (m *MetaClient) ReportMetaClientServer(datasetName string, ipfsData []IpfsData)  error 
```

Inputs:

```shell
datasetName             # 
ipfsData                #  , refer to the `Structs` below for details.
```

Outputs:

```shell
error                  # error or nil
```


## DownloadFile

Definition:
Downloads the file or folder corresponding to the specified IPFS Cid from IPFS. If the Aria2 option is configured, Aria2 tool will be used for downloading. Otherwise, IPFS API will be used.

```shell
func (m *MetaClient) DownloadFile(ipfsCid string, outPath string, downloadUrl string, conf *Aria2Conf) error
```

Inputs:

```shell
ipfsCid                # IPFS Cid to be downloaded
outPath                # Output path for downloaded file
downloadUrl            # Download url address of IPFS service, if the option is not provided, it will automatically query the meta client server.
conf                   # Aria2 related options, including:host, aria2 server address; port,aria2 server port; secret, must be the same value as rpc-secure in aria2 conf
```

Outputs:

```shell
error                  # error or nil
```


## GetDatasetList

Definition:
Gets the file list from the Meta Server based on the specified page number and number of records per page.

```shell
func (m *MetaClient) GetDatasetList(pageNum int, limit int, showCar ...bool) ([]*SourceFile, error)
```

Inputs:

```shell
pageNum                # Which page to query
limit                  # Number of records per page
showCar                # Whether to return storage information,default is false
```

Outputs:

```shell
[]*SourceFile          # List of file descriptions, refer to the `Structs` below for details.
error                  # error or nil
```


## GetSourceFileInfo

Definition:
Gets the IPFS Cid corresponding to the specified file or folder from the Meta Server.

```shell
func (m *MetaClient) GetSourceFileInfo(fileName string) ([]string, error) 
```

Inputs:

```shell
fileName               # File or directory name to be queried
```

Outputs:

```shell
[]string               # List of IPFS Cids corresponding to the queried file name
error                  # error or nil
```


## GetSourceFileStatus

Definition:
Gets detailed information about the file or folder corresponding to the IPFS Cid from the Meta Server.

```shell
func (m *MetaClient) GetSourceFileStatus(ipfsCid string) (*SourceFile, error)
```

Inputs:

```shell
ipfsCid                # IPFS Cid to be queried
```

Outputs:

```shell
*SourceFile            # Information corresponding to the queried IPFS Cid, refer to the `Structs` below for details.
error                  # error or nil
```


## Structs

```go
type IpfsData struct {
    IpfsCid     string `json:"ipfs_cid"`
    SourceName  string `json:"source_name"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}

type SourceFile struct {
	SourceName  string             `json:"source_name"`
	DealFile    string             `json:"deal_file"`
	TaskName    string             `json:"task_name"`
	StorageList []*SplitFileDetail `json:"storage_list"`
	DataList    []*IpfsDataDetail  `json:"data_list"`
}

type SplitFileDetail struct {
    FileName         string            `json:"file_name"`
    DataCid          string            `json:"data_cid"`
    FileSize         int64             `json:"file_size"`
    StorageProviders []StorageProvider `json:"storage_providers"`
}

type StorageProvider struct {
    StorageProviderId string `json:"storage_provider_id"`
    StorageStatus     string `json:"storage_status"`
    DealId            int64  `json:"deal_id"`
    DealCid           string `json:"deal_cid"` // proposal cid or uuid
}

type IpfsDataDetail struct {
    DataId       int64  `json:"data_id"`
    SourceFileId int64  `json:"source_file_id"`
    IpfsCid      string `json:"ipfs_cid"`
    DataSize     int64  `json:"data_size"`
    IsDirector   bool   `json:"is_director"`
    DownloadUrl  string `json:"download_url"`
}
```