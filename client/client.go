package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

type MetaClient struct {
	key   string
	token string
	conf  *MetaConf
}

func NewClient(key, token string, conf ...*MetaConf) *MetaClient {
	var cnf *MetaConf
	if len(conf) > 0 {
		cnf = conf[0]
	}
	c := &MetaClient{
		key:   key,
		token: token,
		conf:  cnf,
	}
	return c
}

func (c *MetaClient) WithMetaServer(url string) *MetaClient {
	if c.conf == nil {
		c.conf = &MetaConf{}
	}
	c.conf.MetaServer = url
	return c
}

func (c *MetaClient) WithIpfs(api, gateway string) *MetaClient {
	if c.conf == nil {
		c.conf = &MetaConf{}
	}
	c.conf.IpfsApi = api
	c.conf.IpfsGateway = gateway
	return c
}

func (c *MetaClient) WithAria2Conf(conf *Aria2Conf) *MetaClient {
	if c.conf == nil {
		c.conf = &MetaConf{}
	}
	c.conf.Aria2Conf = conf
	return c
}

// Upload uploads file or directory to ipfs
func (m *MetaClient) Upload(inputPath string) (ipfsData *IpfsData, err error) {
	if m.conf == nil || m.conf.IpfsApi == "" || m.conf.IpfsGateway == "" {
		return nil, errors.New("ipfs api or gateway is required")
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return
	}

	// create an IPFS Shell client.
	sh := shell.NewShell(m.conf.IpfsApi)
	var ipfsCid string
	if !info.IsDir() {
		ipfsCid, err = uploadFileToIpfs(sh, inputPath)
	} else {
		ipfsCid, err = uploadDirToIpfs(sh, inputPath)
	}
	if err != nil {
		return
	}
	return &IpfsData{
		IpfsCid:     ipfsCid,
		SourceName:  inputPath,
		DataSize:    info.Size(),
		IsDirectory: info.IsDir(),
		DownloadUrl: PathJoin(m.conf.IpfsGateway, "ipfs", ipfsCid),
	}, nil
}

// Download downloads all the files related with the specified ipfsCid default,
// and downloads specific files with the specified downloadUrl
func (m *MetaClient) Download(ipfsCid, outPath string, downloadUrl ...string) error {
	if m.conf == nil || m.conf.Aria2Conf == nil {
		return errors.New("aria2 config is required")
	}

	// check cid from meta server
	downInfo, err := m.DownloadFileInfo(ipfsCid)
	if err != nil {
		return err
	}
	if len(downInfo) == 0 {
		return errors.New("there are no available download links")
	}

	if len(downloadUrl) > 0 && downloadUrl[0] != "" {
		download := downloadUrl[0]
		if !strings.Contains(download, ipfsCid) {
			log.Printf("datacid: %s should be included in the url %s, but it is not.", ipfsCid, download)
		}

		downloadFile := PathJoin(outPath, filepath.Base(downInfo[0].SourceName))
		if downInfo[0].IsDirectory {
			downloadFile = downloadFile + ".tar"
		}

		if err := downloadFileByAria2(m.conf.Aria2Conf, download, downloadFile); err != nil {
			return err
		}

	} else {
		// aria2 download file
		for _, info := range downInfo {
			realUrl := info.DownloadUrl
			if !strings.Contains(realUrl, ipfsCid) {
				log.Printf("datacid: %s should be included in the url %s, but it is not.", ipfsCid, realUrl)
				continue
			}

			downloadFile := PathJoin(outPath, filepath.Base(info.SourceName))
			if info.IsDirectory {
				realUrl = realUrl + "?format=tar"
				downloadFile = downloadFile + ".tar"
			}

			err := downloadFileByAria2(m.conf.Aria2Conf, realUrl, downloadFile)
			if err == nil {
				return nil
			}
		}
	}

	return nil
}

// Backup backups the uploaded files with the datasetName,
// support multiple IpfsData
func (m *MetaClient) Backup(datasetName string, ipfsDataList ...*IpfsData) error {
	if len(ipfsDataList) == 0 {
		return errors.New("ipfsData is required")
	}

	response, err := m.httpPost(JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.StoreSourceFile",
		Params:  []interface{}{datasetName, ipfsDataList},
		Id:      1,
	})
	if err != nil {
		return err
	}

	var res StoreSourceFileResponse
	if err = json.Unmarshal(response, &res); err != nil {
		return err
	}

	if res.Result.Code != "success" {
		return errors.New("failed message from meta server")
	}

	return nil
}

// List lists the backup files with the given datasetName
func (m *MetaClient) List(datasetName string, pageNum, size int) (*DatasetListPager, error) {
	response, err := m.httpPost(JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDatasetList",
		Params:  []interface{}{DatasetListReq{datasetName, pageNum, size}},
		Id:      1,
	})
	if err != nil {
		return nil, err
	}

	var res DatasetListResponse
	if err = json.Unmarshal(response, &res); err != nil {
		return nil, err
	}

	return &res.Result.Data, nil
}

// ListStatus lists the status of backup files
func (m *MetaClient) ListStatus(datasetName, ipfsCid string, pageNum, size int) (*SourceFileStatusPager, error) {
	response, err := m.httpPost(JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileStatus",
		Params:  []interface{}{SourceFileStatusReq{datasetName, ipfsCid, pageNum, size}},
		Id:      1,
	})
	if err != nil {
		return nil, err
	}

	var res SourceFileStatusResponse
	if err = json.Unmarshal(response, &res); err != nil {
		return nil, err
	}

	return &res.Result.Data, nil
}

func (m *MetaClient) SourceFileInfo(ipfsCid string) ([]*IpfsDataDetail, error) {
	response, err := m.httpPost(JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileInfo",
		Params:  []interface{}{ipfsCid},
		Id:      1,
	})
	if err != nil {
		return nil, err
	}
	var res SourceFileInfoResponse
	if err = json.Unmarshal(response, &res); err != nil {
		return nil, err
	}

	return res.Result.Data, nil
}

func (m *MetaClient) DownloadFileInfo(ipfsCid string) ([]*DownloadFileInfo, error) {
	response, err := m.httpPost(JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDownloadFileInfoByIpfsCid",
		Params:  []interface{}{ipfsCid},
		Id:      1,
	})
	if err != nil {
		return nil, err
	}

	var res DownloadFileInfoResponse
	if err = json.Unmarshal(response, &res); err != nil {
		return nil, err
	}

	return res.Result.Data, nil
}

func (m *MetaClient) httpPost(params interface{}) ([]byte, error) {
	if m.key == "" || m.token == "" {
		return nil, errors.New("key or token is required")
	}
	if m.conf == nil {
		return nil, errors.New("meta server is required")
	}
	return httpRequestWithKey(http.MethodPost, m.conf.MetaServer, m.key, m.token, params)

}
