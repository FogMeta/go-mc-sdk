package client

type ClientConf struct {
	Key            string `json:"key"`
	Token          string `json:"token"`
	IpfsApiUrl     string `json:"ipfs_api_url"`
	IpfsGatewayUrl string `json:"ipfs_gateway_url"`
	MetaServerUrl  string `json:"meta_server_url"`
}

type Aria2Conf struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Secret string `json:"secret"`
}

type JsonRpcParams struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

// StoreSourceFile
type StoreSourceFileReq struct {
	SourceName  string `json:"source_name"`
	IsDirector  bool   `json:"is_director"`
	SourceSize  int64  `json:"source_size"`
	DataCid     string `json:"data_cid"`
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
	PageNum int `json:"page_num"`
	Size    int `json:"size"`
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
	SourceName  string       `json:"source_name"`
	DataCid     string       `json:"data_cid"`
	DownloadUrl string       `json:"download_url"`
	StorageList []*SplitFile `json:"storage_list"`
	SourceSize  int64        `json:"source_size"`
	IsDirector  bool         `json:"is_director"`
}

//GetDataCidByName

type DataCidResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Code    string   `json:"code"`
		Message string   `json:"message,omitempty"`
		Data    []string `json:"data,omitempty"`
	} `json:"result"`
	Id int `json:"id"`
}

// GetSourceFileByDataCid

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
