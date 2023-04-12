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
	} `json:"result"`
	Id int `json:"id"`
}

//GetSourceFiles

type SourceFilePageReq struct {
	PageNum int  `json:"page_num"`
	Size    int  `json:"size"`
	ShowCar bool `json:"show_car"`
}

type SourceFilePageResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Total     int64         `json:"total"`
			PageCount int64         `json:"pageCount"`
			Sources   []*SourceFile `json:"files"`
		} `json:"data"`
	} `json:"result"`
	Id int `json:"id"`
}

type SourceFile struct {
	SourceName  string             `json:"source_name"`
	DealFile    string             `json:"deal_file"`
	TaskName    string             `json:"task_name"`
	StorageList []*SplitFileDetail `json:"storage_list"`
	DataList    []*IpfsDataDetail  `json:"data_list"`
}

type IpfsDataDetail struct {
	DataId       int64  `json:"data_id"`
	SourceFileId int64  `json:"source_file_id"`
	IpfsCid      string `json:"ipfs_cid"`
	DataSize     int64  `json:"data_size"`
	IsDirector   bool   `json:"is_director"`
	DownloadUrl  string `json:"download_url"`
}

type SplitFileDetail struct {
	FileName         string            `json:"file_name"`
	DataCid          string            `json:"data_cid"`
	FileSize         int64             `json:"file_size"`
	StorageProviders []StorageProvider `json:"storage_providers"`
}

//GetIpfsCidByName

type IpfsCidResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string   `json:"code"`
		Message string   `json:"message,omitempty"`
		Data    []string `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

// GetSourceFileByIpfsCid

type SourceFileResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string     `json:"code"`
		Message string     `json:"message,omitempty"`
		Data    SourceFile `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

type SplitFile struct {
	FileName         string            `json:"file_name"`
	IpfsCid          string            `json:"data_cid"`
	FileSize         int64             `json:"file_size"`
	StorageProviders []StorageProvider `json:"storage_providers"`
}

type StorageProvider struct {
	StorageProviderId string `json:"storage_provider_id"`
	StorageStatus     string `json:"storage_status"`
	DealId            int64  `json:"deal_id"`
	DealCid           string `json:"deal_cid"` // proposal cid or uuid
}

// GetDownloadFileInfoByIpfsCid

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
	IsDirector  bool   `json:"is_director"`
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
