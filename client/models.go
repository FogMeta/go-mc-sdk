package client

type FileDetails struct {
	FileName     string           `json:"file_name"`
	DataCID      string           `json:"data_cid"`
	DownloadLink string           `json:"download_link"`
	FileSize     int              `json:"file_size"`
	StorageInfo  []StorageDetails `json:"storage_info"`
}

type StorageDetails struct {
	BlockNumber   int    `json:"block_number"`
	DataCID       string `json:"data_cid"`
	Size          int    `json:"size"`
	MinerHash     string `json:"miner_hash"`
	StorageStatus string `json:"storage_status"`
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
	DataCid string `json:"data_cid"`
}

type FileStatusParams struct {
	FileName string `json:"file_name"`
}

type FileStatusResponse struct {
	Status FileDetails `json:"status"`
}
