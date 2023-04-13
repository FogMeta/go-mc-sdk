# go-mc-sdk

* [NewAPIClient](#NewAPIClient)
* [UploadFile](#UploadFile)
* [ReportMetaClientServer](#ReportMetaClientServer)
* [DownloadFile](#DownloadFile)
* [GetDatasetList](#GetDatasetList)
* [GetSourceFileInfo](#GetSourceFileInfo)
* [GetSourceFileStatus](#GetSourceFileStatus)
* [GetIpfsCidStat](#GetIpfsCidStat)

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
ipfsCid                # IPFS cid returned by IPFS server
error                  # error or nil
```

## ReportMetaClientServer

Definition:
Report the Meta Client Server that the file or folder has been uploaded to the IPFS service.

```shell
func (m *MetaClient) ReportMetaClientServer(datasetName string, ipfsData []IpfsData)   error 
```

Inputs:

```shell
datasetName             # The dataset name.
ipfsData                # List of IPFS data which report to Meta Server, refer to the `Structs` below for details of `IpfsData`.
```

Outputs:

```shell
error                  # error or nil
```

Structs:
```go
type IpfsData struct {
    IpfsCid     string `json:"ipfs_cid"`
    SourceName  string `json:"source_name"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}
```


## DownloadFile

Definition:
Downloads the file or folder corresponding to the specified IPFS cid from IPFS. If the Aria2 option is configured, Aria2 tool will be used for downloading. Otherwise, IPFS API will be used.

```shell
func (m *MetaClient) DownloadFile(ipfsCid string, outPath string, downloadUrl string, conf *Aria2Conf) error
```

Inputs:

```shell
ipfsCid                # IPFS cid to be downloaded
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
Gets the dataset list from the Meta Client Server based on the specified datasetName, page number and size of records per page.

```shell
func (m *MetaClient) GetDatasetList(datasetName string, pageNum, size int) (*GetDatasetListPager, error) 
```

Inputs:

```shell
datasetName            # The dataset name.
pageNum                # Page number which to be queryed
size                   # Size of records per page
```

Outputs:

```shell
*GetDatasetListPager   # The pointer to the dataset descriptions, refer to the `Structs` below for details of `GetDatasetListPager`.
error                  # error or nil
```

Structs:
```go
type GetDatasetListPager struct {
    Total       int64            `json:"total"`
    PageCount   int64            `json:"pageCount"`
    DatasetList []*DatasetDetail `json:"dataset_list"`
}

type DatasetDetail struct {
    DataSetName   string            `json:"source_name"`
    DealFile      string            `json:"deal_file"`
    TaskName      string            `json:"task_name"`
    DatasetStatus string            `json:"dataset_status"`
    IpfsList      []*IpfsDataDetail `json:"ipfs_list"`
}

type IpfsDataDetail struct {
    DatasetName string `json:"dataset_name"`
    IpfsCid     string `json:"ipfs_cid"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}
```


## GetSourceFileInfo

Definition:
Gets the source file information corresponding by the specified IPFS cid from the Meta Server.

```shell
func (m *MetaClient) GetSourceFileInfo(ipfsCid string) ([]IpfsDataDetail, error)
```

Inputs:

```shell
ipfsCid               # IPFS cid to be queried
```

Outputs:

```shell
[]IpfsDataDetail       # List of IPFS data details corresponding to the queried IPFS cid, refer to the `Structs` below for details of `IpfsDataDetail`.
error                  # error or nil
```

Structs:
```go
type IpfsDataDetail struct {
    DatasetName string `json:"dataset_name"`
    IpfsCid     string `json:"ipfs_cid"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}
```


## GetSourceFileStatus

Definition:
Gets detailed information about the file or folder corresponding to the IPFS cid from the Meta Server.

```shell
func (m *MetaClient) GetSourceFileStatus(datasetName, ipfsCid string, pageNum, size int) (*GetSourceFileStatusPager, error)
```

Inputs:

```shell
datasetName            # The dataset name.
ipfsCid                # IPFS cid to be queried
pageNum                # Which page to query
size                   # Size of records per page
```

Outputs:

```shell
*GetSourceFileStatusPager  # The pointer to source file status corresponding to the queried IPFS cid, refer to the `Structs` below for details of `GetSourceFileStatusPager`.
error                      # error or nil
```

Structs:
```go
type GetSourceFileStatusPager struct {
    Total     int64              `json:"total"`
    PageCount int64              `json:"pageCount"`
    CarList   []*SplitFileDetail `json:"car_list"`
}

type SplitFileDetail struct {
    FileName         string            `json:"file_name"`
    DataCid          string            `json:"data_cid"`
    FileSize         int64             `json:"file_size"`
    PieceCid         string            `json:"piece_cid"`
    DownloadUrl      string            `json:"download_url"`
    StorageProviders []StorageProvider `json:"storage_providers"`
}

type StorageProvider struct {
    StorageProviderId string `json:"storage_provider_id"`
    StorageStatus     string `json:"storage_status"`
    DealId            int64  `json:"deal_id"`
    DealCid           string `json:"deal_cid"`               // proposal cid or uuid
}
```

## GetIpfsCidStat

Definition:

```shell
func GetIpfsCidInfo(ipfsApiUrl string, ipfsCid string) (IpfsCidInfo, error)
```

Inputs:

```shell
ipfsApiUrl             # API address of IPFS service.
ipfsCid                # IPFS cid to be queried
```

Outputs:

```shell
IpfsCidInfo            # Status information corresponding to the queried IPFS cid, refer to the `Structs` below for details of `IpfsCidInfo`.
error                  # error or nil
```

Structs:
```go
type IpfsCidInfo struct {
    IpfsCid     string `json:"ipfs_cid"`     // The IPFS cid
    DataSize    int64  `json:"data_size"`    // File or diractory size corresponding to the IPFS cid
    IsDirectory bool   `json:"is_directory"` // Is a diracctory corresponding to the IPFS cid
}
```