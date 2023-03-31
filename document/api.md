# meta-client-sdk

* [NewAPIClient](#NewAPIClient)
* [UploadFile](#UploadFile)
* [NotifyMetaServer](#NotifyMetaServer)
* [DownloadFile](#DownloadFile)
* [GetFileLists](#GetFileLists)
* [GetDataCIDByName](#GetDataCIDByName)
* [GetFileInfoByDataCid](#GetFileInfoByDataCid)

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
Uploads a file or folder to IPFS service.

```shell
func (m *MetaClient) UploadFile(ipfsApiUrl, inputPath string) (dataCid string, err error) 
```

Inputs:

```shell
ipfsApiUrl             # API address of IPFS service.
inputPath              # File or director path to be uploaded to IPFS server
```

Outputs:

```shell
dataCid                # Data Cid returned by IPFS server
error                  # error or nil
```

## ReportMetaClientServer

Definition:
Report the Meta Client Server that the file or folder has been uploaded to the IPFS service.

```shell
func (m *MetaClient) ReportMetaClientServer(sourceName string, dataCid string, ipfsGateway string) error 
```

Inputs:

```shell
sourceName             # Name of the file that has been uploaded to the IPFS server.
dataCid                # Data Cid returned by IPFS server.
ipfsGateway            # Gateway address of IPFS service.
```

Outputs:

```shell
error                  # error or nil
```


## DownloadFile

Definition:
Downloads the file or folder corresponding to the specified Data Cid from IPFS. If the Aria2 option is configured, Aria2 tool will be used for downloading. Otherwise, IPFS API will be used.

```shell
func (m *MetaClient) DownloadFile(dataCid string, outPath string, downUrl string, conf *Aria2Conf) error
```

Inputs:

```shell
dataCid                # Data CID to be downloaded
outPath                # Output path for downloaded file
downUrl                # Download url address of IPFS service, if the option is not provided, it will automatically query the meta client server.
conf                   # Aria2 related options, including:host, aria2 server address; port,aria2 server port; secret, must be the same value as rpc-secure in aria2 conf
```

Outputs:

```shell
error                  # error or nil
```


## GetFileLists

Definition:
Gets the file list from the Meta Server based on the specified page number and number of records per page.

```shell
func (m *MetaClient) GetFileLists(pageNum int, limit int, showStorage bool) ([]FileDetails, error)
```

Inputs:

```shell
pageNum                # Which page to query
limit                  # Number of records per page
showStorage            # Whether to return storage information
```

Outputs:

```shell
[]FileDetails          # List of file descriptions returned
error                  # error or nil
```


## GetDataCIDByName

Definition:
Gets the Data Cid corresponding to the specified file or folder from the Meta Server.

```shell
func (m *MetaClient) GetDataCIDByName(fileName string) ([]string, error) 
```

Inputs:

```shell
fileName               # Name of the file to be queried
```

Outputs:

```shell
[]string               # List of Data Cids corresponding to the queried file name
error                  # error or nil
```


## GetFileInfoByDataCid

Definition:
Gets detailed information about the file or folder corresponding to the Data Cid from the Meta Server.

```shell
func (m *MetaClient) GetFileInfoByDataCid(dataCid string) (*FileDetails, error)
```

Inputs:

```shell
dataCid                # Data Cid to be queried
```

Outputs:

```shell
*FileDetails           # Information corresponding to the queried Data Cid
error                  # error or nil
```
