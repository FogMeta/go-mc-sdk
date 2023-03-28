package client

import (
	"encoding/json"
	"fmt"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
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
	//err = m.NotifyMetaServer(targetPath, dataCid)
	//if err != nil {
	//	return "", err
	//}

	return dataCid, nil
}

func (m *MetaClient) NotifyMetaServer(sourceName string, dataCid string) error {

	sourceSize := walkDirSize(sourceName)
	logs.GetLogger().Infoln("upload total size is:", sourceSize)

	var params []interface{}
	downUrl := pathJoin(m.IpfsGatewayUrl, "ipfs/", dataCid)
	params = append(params, SourceFileReq{sourceName, sourceSize, dataCid, downUrl})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "StoreSourceFile",
		Params:  params,
		Id:      1,
	}

	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return err
	}

	res := APIResp{}
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
	if conf != nil {
		downUrl := pathJoin(m.IpfsGatewayUrl, "ipfs/", dataCid)
		outFile := pathJoin(outPath, dataCid)

		if *isDir {
			downUrl = downUrl + "?format=tar"
			outFile = outFile + ".tar"
		}

		err = downloadFileByAria2(conf, downUrl, outFile)
		if err == nil {
			return nil
		}
		logs.GetLogger().Warnln("download ", dataCid, "by aria2 error:", err)
	}

	return downloadFromIpfs(sh, dataCid, outPath)
}

func (m *MetaClient) GetFileLists(page, limit int, showStorageInfo bool) ([]FileDetails, error) {

	var params []interface{}
	params = append(params, SourceFilePageReq{page, limit})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "GetSourceFiles",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := APIResp{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return nil, nil
}

func (m *MetaClient) GetDataCIDByName(fileName string) ([]string, error) {
	var params []interface{}
	params = append(params, fileName)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "",
		Method:  "GetSourceFilesByName",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}
	res := APIResp{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return nil, nil
}

func (m *MetaClient) GetFileInfoByDataCid(dataCid string) (*FileDetails, error) {

	var params []interface{}
	params = append(params, dataCid)
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

	res := APIResp{}
	err = json.Unmarshal(response, res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return nil, nil
}
