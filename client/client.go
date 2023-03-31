package client

import (
	"encoding/json"
	"errors"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	shell "github.com/ipfs/go-ipfs-api"
	"path/filepath"
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

func (m *MetaClient) UploadFile(ipfsApiUrl, inputPath string) (dataCid string, err error) {
	// Creates an IPFS Shell client.
	sh := shell.NewShell(ipfsApiUrl)

	isInputFile, err := isFile(inputPath)
	if err != nil {
		return "", err
	}

	if *isInputFile {
		dataCid, err = uploadFileToIpfs(sh, inputPath)
	} else {
		dataCid, err = uploadDirToIpfs(sh, inputPath)
	}
	if err != nil {
		return "", err
	}

	return dataCid, nil
}

func (m *MetaClient) DownloadFile(dataCid, outPath string, downUrl string, conf *Aria2Conf) error {

	if conf == nil {
		return errors.New("need aria2 server config")
	}

	// check data cid from meta server
	downInfo, err := m.GetDownloadFileInfoByDataCid(dataCid)
	if err != nil || len(downInfo) == 0 {
		logs.GetLogger().Errorf("Get Download File Info Error: %s \n", err)
		return err
	}

	if downUrl != "" {
		downFile := pathJoin(outPath, filepath.Base(downInfo[0].SourceName))
		if downInfo[0].IsDirector {
			downFile = downFile + ".tar"
		}

		err := downloadFileByAria2(conf, downUrl, downFile)
		if err == nil {
			logs.GetLogger().Info("download ", dataCid, "by aria2 success")
			return nil
		}
	}

	// aria2 download file
	for _, down := range downInfo {
		realUrl := down.DownloadUrl
		downFile := pathJoin(outPath, filepath.Base(down.SourceName))
		if down.IsDirector {
			realUrl = realUrl + "?format=tar"
			downFile = downFile + ".tar"
		}

		err := downloadFileByAria2(conf, realUrl, downFile)
		if err == nil {
			logs.GetLogger().Info("download ", dataCid, "by aria2", err)
			break
		}
		logs.GetLogger().Warn("download ", dataCid, " url ", realUrl, " by aria2 error:", err)
		//remove failed
	}

	return nil
}

func (m *MetaClient) NotifyMetaClientServer(sourceName string, dataCid string, ipfsGateway string) error {

	isFile, err := isFile(sourceName)
	if err != nil {
		return err
	}

	sourceSize := walkDirSize(sourceName)
	logs.GetLogger().Infoln("upload total size is:", sourceSize)

	var params []interface{}
	downUrl := pathJoin(ipfsGateway, "ipfs/", dataCid)
	params = append(params, StoreSourceFileReq{sourceName, !(*isFile), sourceSize, dataCid, downUrl})
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

func (m *MetaClient) GetFileLists(page, limit int, showStorage bool) ([]SourceFile, error) {

	var params []interface{}
	params = append(params, SourceFilePageReq{page, limit, showStorage})
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
	}

	return nil, nil
}

func (m *MetaClient) GetDataCIDByName(fileName string) ([]string, error) {
	var params []interface{}
	params = append(params, fileName)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDataCidByName",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}
	res := DataCidResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	return nil, nil
}

func (m *MetaClient) GetFileInfoByDataCid(dataCid string) (*SourceFile, error) {

	var params []interface{}
	params = append(params, dataCid)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileByDataCid",
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

func (m *MetaClient) GetDownloadFileInfoByDataCid(dataCid string) ([]DownloadFileInfo, error) {

	var params []interface{}
	params = append(params, dataCid)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDownloadFileInfoByDataCid",
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

func (m *MetaClient) RebuildDataCID(fileName string) error {
	// TODO
	return nil
}
