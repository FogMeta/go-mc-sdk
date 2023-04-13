# APIs

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
Creates a Meta Client instance.

```shell
func NewAPIClient(key, token, metaUrl string) *MetaClient
```

Inputs:

```shell
key                    # Swan API key. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
token                  # Swan API access token. Acquire from [Swan Platform](https://console.filswan.com/#/dashboard) -> "My Profile"->"Developer Settings". 
metaUrl                # Address of Meta Server.
```

Outputs:

```shell
*MetaClient            # Created Meta Client instance.
```

## UploadFile

Definition:
Upload a file or directory to the IPFS service.

```shell
func (m *MetaClient) UploadFile(ipfsApiUrl, inputPath string) (ipfsCid string, err error) 
```

Inputs:

```shell
ipfsApiUrl             # API address of IPFS service.
inputPath              # File or directory path to be uploaded to the IPFS server
```

Outputs:

```shell
ipfsCid                # IPFS cid returned by IPFS server
error                  # error or nil
```

## ReportMetaClientServer

Definition:
Report the file or directory information to the Meta Client Service.

```shell
func (m *MetaClient) ReportMetaClientServer(datasetName string, ipfsData []IpfsData)   error 
```

Inputs:

```shell
datasetName             # The dataset name.
ipfsData                # A list of IPFS data that report to Meta Client Server, refer to the `IpfsData` struct
```

Outputs:

```shell
error                  # error or nil
```

**About `IpfsData`:**

```go
type IpfsData struct {
    IpfsCid     string `json:"ipfs_cid"`
    SourceName  string `json:"source_name"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}
```

| Field Name | Data Type | Explanation |
| --- | --- | --- |
| IpfsCid | string | The CID (Content Identifier) of the data in IPFS, which is used to uniquely identify the data |
| SourceName | string | The name of the data source |
| DataSize | int64 | The size of the data in bytes |
| IsDirectory | bool | The type of data, used to differentiate whether it is a directory or not |
| DownloadUrl | string | The download link for the data, used to download the data file from IPFS |


## DownloadFile

Definition:
Downloads the file or directory corresponding to the specified IPFS cid from IPFS. If the Aria2 option is configured, Aria2 tool will be used for downloading. Otherwise, IPFS API will be used.

```shell
func (m *MetaClient) DownloadFile(ipfsCid string, outPath string, downloadUrl string, conf *Aria2Conf) error
```

Inputs:

```shell
ipfsCid                # IPFS cid to be downloaded
outPath                # Output path for the downloaded file
downloadUrl            # Download url address of the IPFS service, if the option is not provided, it will automatically query the meta client server.
conf                   # Aria2 related options, including:host, aria2 server address; port,aria2 server port; secret, must be the same value as rpc-secure in aria2 conf
```

Outputs:

```shell
error                  # error or nil
```


## GetDatasetList

Definition:
Get the dataset list from the Meta Client Server based on the dataset name, page number, and size of records per page.

```shell
func (m *MetaClient) GetDatasetList(datasetName string, pageNum, size int) (*GetDatasetListPager, error) 
```

Inputs:

```shell
datasetName            # The dataset name.
pageNum                # Page number which to be queried
size                   # Size of records per page
```

Outputs:

```shell
*GetDatasetListPager   # The pointer to the dataset descriptions, refer to the `Structs`.
error                  # error or nil
```

**About `GetDatasetListPager`:**
```go
type GetDatasetListPager struct {
    Total       int64            `json:"total"`
    PageCount   int64            `json:"pageCount"`
    DatasetList []*DatasetDetail `json:"dataset_list"`
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| Total | int64 | The total number of datasets |
| PageCount | int64 | The total number of pages |
| DatasetList | []*DatasetDetail | A list of dataset details, which contains the details of each dataset |

```
type DatasetDetail struct {
    DataSetName   string            `json:"source_name"`
    DealFile      string            `json:"deal_file"`
    TaskName      string            `json:"task_name"`
    DatasetStatus string            `json:"dataset_status"`
    IpfsList      []*IpfsDataDetail `json:"ipfs_list"`
}

```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| DataSetName | string | The name of the dataset |
| DealFile | string | The file name of the dataset |
| TaskName | string | The name of the task associated with the dataset |
| DatasetStatus | string | The status of the dataset, indicating downloading, downloaded, or others |
| IpfsList | []*IpfsDataDetail | A list of data info associated with the dataset |


```
type IpfsDataDetail struct {
    DatasetName string `json:"dataset_name"`
    IpfsCid     string `json:"ipfs_cid"`
    DataSize    int64  `json:"data_size"`
    IsDirectory bool   `json:"is_directory"`
    DownloadUrl string `json:"download_url"`
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| DatasetName | string | The name of the dataset |
| IpfsCid | string | The CID (Content Identifier) of the IPFS data, which is used to uniquely identify the data |
| DataSize | int64 | The size of the IPFS data in bytes |
| IsDirectory | bool | The type of data, used to differentiate whether it is a directory or not |
| DownloadUrl | string | The download url for the IPFS data, used to download the data file from the IPFS gateway |


## GetSourceFileInfo

Definition:
Get the source file information by the IPFS CID from the Meta Server.

```shell
func (m *MetaClient) GetSourceFileInfo(ipfsCid string) ([]IpfsDataDetail, error)
```

Inputs:

```shell
ipfsCid               # IPFS cid to be queried
```

Outputs:

```shell
[]IpfsDataDetail       # List of IPFS data details corresponding to the queried IPFS CID
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
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| DatasetName | string | The name of the dataset |
| IpfsCid | string | The CID (Content Identifier) of the IPFS data, which is used to uniquely identify the data |
| DataSize | int64 | The size of the IPFS data in bytes |
| IsDirectory | bool | The type of data, used to differentiate whether it is a directory or not |
| DownloadUrl | string | The download link for the IPFS data, used to download the data file from the IPFS gateway |


## GetSourceFileStatus

Definition:
Get the information of the file or directory corresponding to the IPFS cid from the Meta Client Server.

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
*GetSourceFileStatusPager  # The pointer to source file status corresponding to the queried IPFS CID
error                      # error or nil
```


**About `GetSourceFileStatusPager`:**
```go
type GetSourceFileStatusPager struct {
    Total     int64              `json:"total"`
    PageCount int64              `json:"pageCount"`
    CarList   []*SplitFileDetail `json:"car_list"`
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| Total | int64 | The total number of split files |
| PageCount | int64 | The total number of pages in the split file list |
| CarList | []*SplitFileDetail | A list of split file details, which contains the details of each split file |


```
type SplitFileDetail struct {
    FileName         string            `json:"file_name"`
    DataCid          string            `json:"data_cid"`
    FileSize         int64             `json:"file_size"`
    PieceCid         string            `json:"piece_cid"`
    DownloadUrl      string            `json:"download_url"`
    StorageProviders []StorageProvider `json:"storage_providers"`
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| FileName | string | The name of the split file |
| DataCid | string | The CID (Content Identifier) of the data associated with the split file |
| FileSize | int64 | The size of the split file in bytes |
| PieceCid | string | The CID (Content Identifier) of the piece associated with the split file |
| DownloadUrl | string | The download url for the CAR of the split file, used to download the CAR file |
| StorageProviders | []StorageProvider | A list of storage providers that have stored the split file |


```
type StorageProvider struct {
    StorageProviderId string `json:"storage_provider_id"`
    StorageStatus     string `json:"storage_status"`
    DealId            int64  `json:"deal_id"`
    DealCid           string `json:"deal_cid"`               
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| StorageProviderId | string | The ID of the storage provider |
| StorageStatus | string | The status of the deal in the Filecoin network |
| DealId | int64 | The dealID in the filecoin network |
| DealCid | string | The proposal CID or deal UUID in the filecoin network |


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
IpfsCidInfo            # Status information corresponds to the queried IPFS CID
error                  # error or nil
```

**About `IpfsCidInfo`:**
```go
type IpfsCidInfo struct {
    IpfsCid     string `json:"ipfs_cid"`     // The IPFS CID
    DataSize    int64  `json:"data_size"`    // File or directory size corresponding to the IPFS cid
    IsDirectory bool   `json:"is_directory"` // Is a directory corresponding to the IPFS CID
}
```
| Field Name | Data Type | Explanation |
| --- | --- | --- |
| IpfsCid | string | The CID (Content Identifier) of the IPFS data |
| DataSize | int64 | The size of the file or directory corresponding to the IPFS CID |
| IsDirectory | bool | Whether the data corresponding to the IPFS CID is a directory or not |
