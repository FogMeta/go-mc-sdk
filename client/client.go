package client

import (
	"encoding/json"
	"fmt"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
	shell "github.com/ipfs/go-ipfs-api"
)

type MetaClient struct {
	ApiKey         string
	ApiToken       string
	IpfsApiUrl     string
	IpfsGatewayUrl string
	MetaUrl        string

	sh    *shell.Shell
	aria2 *client.Aria2Client
}

func NewAPIClient(key, token, ipfsApiUrl, ipfsGatewayUrl, metaUrl string) *MetaClient {

	c := &MetaClient{
		ApiKey:         key,
		ApiToken:       token,
		IpfsApiUrl:     ipfsApiUrl,
		IpfsGatewayUrl: ipfsGatewayUrl,
		MetaUrl:        metaUrl,
	}
	// TODO: check key and token ,need meta server api

	return c
}

func (m *MetaClient) UploadFile(targetPath string) (dataCid string, err error) {
	// Creates an IPFS Shell client.
	sh := shell.NewShell(m.IpfsApiUrl)

	isInputFile, err := isFile(targetPath)
	if err != nil {
		return "", err
	}

	if *isInputFile {
		dataCid, err = uploadFileToIpfs(sh, targetPath)
	} else {
		dataCid, err = uploadDirToIpfs(sh, targetPath)
	}
	if err != nil {
		return "", err
	}

	// TODO: notify meta server the result
	sourceSize := walkDirSize(targetPath)
	logs.GetLogger().Infoln("upload total size is:", sourceSize)
	//err = m.notifyMetaServer(targetPath, sourceSize, dataCid)
	//if err != nil {
	//	return "", err
	//}

	return dataCid, nil
}

func (m *MetaClient) notifyMetaServer(sourceName string, sourceSize int64, dataCid string) error {

	var params []interface{}
	params = append(params, SourceFileReq{sourceName, sourceSize, dataCid, m.IpfsGatewayUrl})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "",
		Params:  params,
		Id:      1,
	}

	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return err
	}

	res := JsonRpcResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) Error: %s \n", response, err)
		return err
	}

	return nil
}

func (m *MetaClient) DownloadFile(dataCid, outPath string, conf *Aria2Conf) error {
	// Creates an IPFS Shell client.
	sh := shell.NewShell(m.IpfsApiUrl)

	isDir, err := dataCidIsDir(sh, dataCid)
	if err != nil || isDir == nil {
		return err
	}

	// aria2 download file
	if !(*isDir) && (conf != nil) {

		downUrl := utils.UrlJoin(m.IpfsGatewayUrl, "ipfs/", dataCid)
		return downloadFileByAria2(conf, downUrl, PathJoin(outPath, dataCid))
	}

	return downloadFromIpfs(sh, dataCid, outPath)
}

func (m *MetaClient) GetFileLists(page, limit uint64, showStorageInfo bool) ([]FileDetails, error) {

	var params []interface{}
	params = append(params, FileListsParams{page, limit, showStorageInfo})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := FileListsResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return res.FileLists, nil
}

func (m *MetaClient) GetFileDataCID(fileName string) ([]string, error) {
	var params []interface{}
	params = append(params, FileDataCIDParams{fileName})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}
	res := FileDataCIDResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return res.DataCids, nil
}

func (m *MetaClient) GetFileInfo(fileName string) (*FileDetails, error) {

	var params []interface{}
	params = append(params, FileInfoParams{fileName})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := FileInfoResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return &res.Info, nil
}
