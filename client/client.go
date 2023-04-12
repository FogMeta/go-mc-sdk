package client

import (
	"encoding/json"
	"errors"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	shell "github.com/ipfs/go-ipfs-api"
	"path/filepath"
	"strings"
)

type MetaClient struct {
	ApiKey   string
	ApiToken string
	MetaUrl  string

	sh    *shell.Shell
	aria2 *client.Aria2Client
}

func NewAPIClient(key, token, metaUrl string) *MetaClient {

	c := &MetaClient{
		ApiKey:   key,
		ApiToken: token,
		MetaUrl:  metaUrl,
	}

	return c
}

func (m *MetaClient) UploadFile(ipfsApiUrl, inputPath string) (ipfsCid string, err error) {
	// Creates an IPFS Shell client.
	sh := shell.NewShell(ipfsApiUrl)

	isInputFile, err := isFile(inputPath)
	if err != nil {
		return "", err
	}

	if *isInputFile {
		ipfsCid, err = uploadFileToIpfs(sh, inputPath)
	} else {
		ipfsCid, err = uploadDirToIpfs(sh, inputPath)
	}
	if err != nil {
		return "", err
	}

	return ipfsCid, nil
}

func (m *MetaClient) DownloadFile(ipfsCid, outPath string, downloadUrl string, conf *Aria2Conf) error {

	if conf == nil {
		return errors.New("need aria2 server config")
	}

	// check data cid from meta server
	downInfo, err := m.GetDownloadFileInfoByIpfsCid(ipfsCid)
	if err != nil || len(downInfo) == 0 {
		logs.GetLogger().Errorf("Get Download File Info Error: %s \n", err)
		return err
	}

	if downloadUrl != "" {
		if !strings.Contains(downloadUrl, ipfsCid) {
			logs.GetLogger().Warnf("datacid: %s should be included in the url %s, but it is not.", ipfsCid, downloadUrl)
		}

		downloadFile := PathJoin(outPath, filepath.Base(downInfo[0].SourceName))
		if downInfo[0].IsDirector {
			downloadFile = downloadFile + ".tar"
		}

		err := downloadFileByAria2(conf, downloadUrl, downloadFile)
		if err == nil {
			logs.GetLogger().Info("download ", ipfsCid, "by aria2 success")
			return nil
		}
		logs.GetLogger().Warn("download ", ipfsCid, " url ", downloadUrl, " by aria2 error:", err)

	} else {
		// aria2 download file
		for _, info := range downInfo {
			realUrl := info.DownloadUrl
			if !strings.Contains(realUrl, ipfsCid) {
				logs.GetLogger().Warnf("datacid: %s should be included in the url %s, but it is not.", ipfsCid, realUrl)
				continue
			}

			downloadFile := PathJoin(outPath, filepath.Base(info.SourceName))
			if info.IsDirector {
				realUrl = realUrl + "?format=tar"
				downloadFile = downloadFile + ".tar"
			}

			err := downloadFileByAria2(conf, realUrl, downloadFile)
			if err == nil {
				logs.GetLogger().Info("download ", ipfsCid, "by aria2 success")
				return nil
			}

			logs.GetLogger().Warn("download ", ipfsCid, " url ", realUrl, " by aria2 error:", err)
		}
	}

	return errors.New("there are no available download links")
}

func (m *MetaClient) ReportMetaClientServer(sourceFile SourceFileReq) error {
	var params []interface{}
	params = append(params, sourceFile)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.StoreSourceFile",
		Params:  params,
		Id:      1,
	}

	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return err
	}

	res := StoreSourceFileResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return err
	}
	logs.GetLogger().Info(res)

	return nil
}

func (m *MetaClient) GetFileLists(pageNum, limit int, opts ...ListOption) ([]SourceFile, error) {

	op := defaultOptions()
	for _, opt := range opts {
		opt.apply(&op)
	}

	logs.GetLogger().Info("with opts is:", op.ShowCar)
	var params []interface{}
	params = append(params, SourceFilePageReq{pageNum, limit, op.ShowCar})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFiles",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := SourceFilePageResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	sources := res.Result.Data.Sources
	for index, source := range sources {
		logs.GetLogger().Infof("Index: %d, Source: %+v", index, source)

		stores := source.StorageList
		for i, store := range stores {
			logs.GetLogger().Infof("Store-%d: %+v", i, store)
			providers := store.StorageProviders
			for ii, provider := range providers {
				logs.GetLogger().Infof("Provider-%d: %+v", ii, provider)
			}
		}
		dataList := source.DataList
		for i, data := range dataList {
			logs.GetLogger().Infof("Data-%d: %+v", i, data)
		}
	}

	return nil, nil
}

func (m *MetaClient) GetIpfsCidByName(fileName string) ([]string, error) {
	var params []interface{}
	params = append(params, fileName)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetIpfsCidByName",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}
	res := IpfsCidResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	return nil, nil
}

func (m *MetaClient) GetFileInfoByIpfsCid(ipfsCid string) (*SourceFile, error) {

	var params []interface{}
	params = append(params, ipfsCid)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileByIpfsCid",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}

	res := SourceFileResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	source := res.Result.Data
	logs.GetLogger().Infof("Source: %+v", source)
	stores := source.StorageList
	for _, store := range stores {
		logs.GetLogger().Infof("Store: %+v", store)
		providers := store.StorageProviders
		for ii, provider := range providers {
			logs.GetLogger().Infof("Provider-%d: %+v", ii, provider)
		}
	}

	return nil, nil
}

func (m *MetaClient) GetDownloadFileInfoByIpfsCid(ipfsCid string) ([]DownloadFileInfo, error) {

	var params []interface{}
	params = append(params, ipfsCid)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDownloadFileInfoByIpfsCid",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}

	res := DownloadFileInfoResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	return res.Result.Data, nil
}

func (m *MetaClient) RebuildIpfsCid(fileName string) error {
	// TODO
	return nil
}
