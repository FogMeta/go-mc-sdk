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

type JsonRpcResponse struct {
	JsonRpc string        `json:"jsonrpc"`
	result  []interface{} `json:"result"`
}

type SourceFileReq struct {
	SourceName  string `json:"source_name"`
	SourceSize  int64  `json:"source_size"`
	DataCid     string `json:"data_cid"`
	DownloadUrl string `json:"download_url"`
}

type SourceFile struct {
	SourceName   string      `json:"source_name"`
	DataCid      string      `json:"data_cid"`
	DownloadLink string      `json:"download_link"`
	StorageList  []SplitFile `json:"storage_list"`
	SourceSize   int64       `json:"source_size"`
}

type SplitFile struct {
	DataCid           string `json:"data_cid"`
	FileSize          int64  `json:"file_size"`
	StorageProviderId string `json:"storage_provider_id"`
	StorageStatus     string `json:"storage_status"`
	DealId            int64  `json:"deal_id"`
	DealCid           string `json:"deal_cid"` // proposalcid or uuid
}

type FileDetails struct {
	FileName     string           `json:"file_name"`
	DataCID      string           `json:"data_cid"`
	DownloadLink string           `json:"download_link"`
	FileSize     int              `json:"file_size"`
	StorageInfo  []StorageDetails `json:"storage_info"`
}

type SpDetail struct {
	StorageProviderId string `json:"storage_provider_id"`
	StorageStatus     string `json:"storage_status"`
}

type StorageDetails struct {
	DataCID  string     `json:"data_cid"`
	Size     int        `json:"size"`
	SpDetail []SpDetail `json:"sp_detail"`
}

type FileListsParams struct {
	Page        uint64 `json:"page"`
	Limit       uint64 `json:"limit"`
	ShowStorage bool   `json:"show_storage"`
}

type FileListsResponse struct {
	FileLists []FileDetails `json:"file_lists"`
}

type FileDataCIDParams struct {
	FileName string `json:"file_name"`
}

type FileDataCIDResponse struct {
	DataCids []string `json:"data_cids"`
}

type FileInfoParams struct {
	FileName string `json:"file_name"`
}

type FileInfoResponse struct {
	Info FileDetails `json:"info"`
}
