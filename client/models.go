package client

type ClientConf struct {
	Key            string    `toml:"key"`
	Token          string    `toml:"token"`
	IpfsApiUrl     string    `toml:"ipfs_api_url"`
	IpfsGatewayUrl string    `toml:"ipfs_gateway_url"`
	MetaServerUrl  string    `toml:"meta_server_url"`
	Aria2          Aria2Conf `toml:"aria2"`
}

type Aria2Conf struct {
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	Secret string `toml:"secret"`
}

type JsonRpcParams struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type IpfsCidInfo struct {
	IpfsCid     string `json:"ipfs_cid"`
	DataSize    int64  `json:"data_size"`
	IsDirectory bool   `json:"is_directory"`
}

// StoreSourceFile
// StoreSourceFile(ctx context.Context, datasetName string, ipfsData []IpfsData) APIResp

type IpfsData struct {
	IpfsCid     string `json:"ipfs_cid"`
	SourceName  string `json:"source_name"`
	DataSize    int64  `json:"data_size"`
	IsDirectory bool   `json:"is_directory"`
	DownloadUrl string `json:"download_url"`
}

type StoreSourceFileResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
		Data    string `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

// GetDatasetList
// GetDatasetList(ctx context.Context, req GetDatasetListReq) APIResp

type GetDatasetListReq struct {
	DatasetName string `json:"dataset_name"`
	PageNum     int    `json:"page_num"`
	Size        int    `json:"size"`
}

type GetDatasetListResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string              `json:"code"`
		Message string              `json:"message,omitempty"`
		Data    GetDatasetListPager `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

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
	SourceName  string `json:"source_name"`
}

// GetSourceFileInfo
// func (api *ApiImpl) GetSourceFileInfo(ctx context.Context, ipfsCid string) APIResp

type GetSourceFileInfoResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string           `json:"code"`
		Message string           `json:"message,omitempty"`
		Data    []IpfsDataDetail `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

// GetSourceFileStatus
// func (api *ApiImpl) GetSourceFileStatus(ctx context.Context, req GetSourceFileStatusReq) APIResp

type GetSourceFileStatusReq struct {
	DatasetName string `json:"dataset_name"`
	IpfsCid     string `json:"ipfs_cid"`
	PageNum     int    `json:"page_num"`
	Size        int    `json:"size"`
}

type GetSourceFileStatusResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string                   `json:"code"`
		Message string                   `json:"message,omitempty"`
		Data    GetSourceFileStatusPager `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

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
	DealCid           string `json:"deal_cid"` // proposal cid or uuid
}

// GetDownloadFileInfoByIpfsCid
// func (api *ApiImpl) GetDownloadFileInfoByIpfsCid(ctx context.Context, ipfsCid string) APIResp

type DownloadFileInfoResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string             `json:"code"`
		Message string             `json:"message,omitempty"`
		Data    []DownloadFileInfo `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

type DownloadFileInfo struct {
	SourceName  string `json:"source_name"`
	DownloadUrl string `json:"download_url"`
	IsDirectory bool   `json:"is_directory"`
}

// StoreSourceFileByGroup
// func (api *ApiImpl) StoreSourceFileByGroup(ctx context.Context, groupName string, dataList [][]IpfsData) APIResp
//type IpfsData struct {
//	IpfsCid     string `json:"ipfs_cid"`
//	SourceName  string `json:"source_name"`
//	DataSize    int64  `json:"data_size"`
//	IsDirectory bool   `json:"is_directory"`
//	DownloadUrl string `json:"download_url"`
//}

type StoreSourceFileByGroupResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code string `json:"code"` //success
	} `json:"result"`
	Id int `json:"id"`
}

// GetDatasetsByGroupName
// func (api *ApiImpl) GetDatasetsByGroupName(ctx context.Context, groupName string) APIResp

type GetDatasetsByGroupNameResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code string                       `json:"code"`
		Data []GetDatasetsByGroupNameResp `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}
type GetDatasetsByGroupNameResp struct {
	DatasetId     int64  `json:"dataset_id"`
	DatasetName   string `json:"dataset_name"`
	DatasetStatus string `json:"dataset_status"`
}

// StoreCarFiles
// func (api *ApiImpl) StoreCarFiles(ctx context.Context, req StoreCarFilesReq) APIResp
type StoreCarFilesReq struct {
	DatasetId int64      `json:"dataset_id"`
	CarList   []*CarInfo `json:"car_list"`
}

type CarInfo struct {
	FileName    string `json:"file_name"`
	DataCid     string `json:"data_cid"`
	SourceSize  int64  `json:"source_size"`
	CarSize     int64  `json:"car_size"`
	PieceCid    string `json:"piece_cid"`
	DownloadUrl string `json:"download_url"`
}

type StoreCarFilesResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code string `json:"code"` //success
	} `json:"result"`
	Id int `json:"id"`
}

//////////////////////////////////////////

type DagLink struct {
	Hash struct {
		Target string `json:"/"`
	} `json:"Hash"`
	Name  string `json:"Name"`
	Tsize int64  `json:"Tsize"`
}
type DagGetResponse struct {
	Data struct {
		Target struct {
			Bytes string `json:"bytes"`
		} `json:"/"`
	} `json:"Data"`
	Links []DagLink `json:"Links,omitempty"`
}

type TreeNode struct {
	Path  string
	Name  string
	Hash  string
	Size  uint64
	Dir   bool
	IsTop bool
	Group int
	Child []*TreeNode
}

// list option
type listOption struct {
	ShowCar bool
}

type ListOption interface {
	apply(*listOption)
}

type funcOption struct {
	f func(*listOption)
}

func (fdo *funcOption) apply(do *listOption) {
	fdo.f(do)
}

func showStorageOption(f func(*listOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithShowCar(show bool) ListOption {
	return showStorageOption(func(o *listOption) {
		o.ShowCar = show
	})
}
func defaultOptions() listOption {
	return listOption{
		ShowCar: false,
	}
}
