# Meta-Client-SDK

Meta-Client-SDK is a Web3 data service that helps users store backups and recover data. This SDK supports automated recording of data storage information and can store data automatically to IPFS gateway and Filecoin network, providing fast data retrieval, permanent backup, and data recovery capabilities.

## Features

Meta-Client-SDK provides the following features:

- Upload files or folders to the IPFS gateway
- Report data information to the Meta-Client server (server will automatically complete data processing and deal-sending functions)
- Store data to the Filecoin network
- Store data in the IPFS gateway
- Download files or folders to the local machine
- Query DataCID for a file by its filename
- View a list of all files under the current user
- Query storage information and status of a single file or folder

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
go get github.com/FogMeta/Meta-Client-SDK
```


## Usage

### Initialization

First, you need to create a MetaClient object, which can be initialized as follows:

```go
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()
}
```
### Upload Files or Folders
To upload files or folders to IPFS gateway and Filecoin network, you can use the following method:
```
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()

    // Upload a single file
    path := "/path/to/file"
    result, err := cli.UploadFile(path)
    if err != nil {
        // handle error
    }
    fmt.Println(result)

    // Upload a folder
    path := "/path/to/folder"
    result, err := cli.UploadFolder(path)
    if err != nil {
        // handle error
    }
    fmt.Println(result)
}

```
### Report Data-related Information
To report data-related information to the Meta-Client server, you can use the following method:

```
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()

    cid := "bafkreiaugtrtgvtpnv23ipcfivfzljavh2olndt7blbnslvnhdbt7awqve"
    result, err := cli.Report(cid)
    if err != nil {
        // handle error
    }
    fmt.Println(result)
}

```
### Download Files or Folders
To download files or folders from the IPFS gateway and Filecoin network, you can use the following method:

```
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()

    // Download a single file
    cid := "bafkreiaugtrtgvtpnv23ipcfivfzljavh2olndt7blbnslvnhdbt7awqve"
    path := "/path/to/save/file"
    err := cli.DownloadFile(cid, path)
    if err != nil {
        // handle error
    }

    // Download a folder
    cid := "bafkreiaugtrtgvtpnv23ipcfivf

```

### Get DataCID for a File by its Filename
To get the DataCID for a file by its filename, you can use the following method:
```
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()

    filename := "file.txt"
    cid, err := cli.GetDataCID(filename)
    if err != nil {
        // handle error
    }
    fmt.Println(cid)
}

```
### View List of Files and Storage Status
To view a list of all files under the current user, or to query storage information and the status of a single file or folder, you can use the following method:
```
import "github.com/FogMeta/Meta-Client-SDK/client"

func main() {
    cli := client.NewMetaClient()

    // View list of files
    list, err := cli.List()
    if err != nil {
        // handle error
    }
    fmt.Println(list)

    // Query storage information and status of a single file or folder
    cid := "bafkreiaugtrtgvtpnv23ipcfivfzljavh2olndt7blbnslvnhdbt7awqve"
    info, err := cli.Info(cid)
    if err != nil {
        // handle error
    }
    fmt.Println(info)
}

```
## API Documentation

For detailed API lists, please check out the [API Documentation](meta-client-sdk/document/api.md ':include').

## Contributing

Contributions to Meta-Client-SDK are welcome! If you find any errors or want to add new features, please submit an [Issue](https://github.com/FogMeta/meta-client-sdk/issues), or initiate a [Pull Request](https://github.com/FogMeta/meta-client-sdk/pulls).

## License

Meta-Client-SDK is licensed under the MIT License.
