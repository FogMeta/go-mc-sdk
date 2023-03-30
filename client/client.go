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

func (m *MetaClient) NotifyMetaServer(sourceName string, dataCid string) error {

	isFile, err := isFile(sourceName)
	if err != nil {
		return err
	}

	sourceSize := walkDirSize(sourceName)
	logs.GetLogger().Infoln("upload total size is:", sourceSize)

	var params []interface{}
	downUrl := pathJoin(m.IpfsGatewayUrl, "ipfs/", dataCid)
	params = append(params, StoreSourceFileReq{sourceName, !(*isFile), sourceSize, dataCid, downUrl})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.StoreSourceFile",
		Params:  params,
		Id:      1,
	}

	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return err
	}

	res := StoreSourceFileResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		fmt.Printf("Parse Response (%s) Error: %s \n", response, err)
		return err
	}
	logs.GetLogger().Info(res)

	return nil
}

func (m *MetaClient) GetFileLists(page, limit int, showStorageInfo bool) ([]SourceFile, error) {

	var params []interface{}
	params = append(params, SourceFilePageReq{page, limit})
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
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

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
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}
	res := DataCidResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
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
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := SourceFileResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}
	logs.GetLogger().Info(res)

	return nil, nil
}

func (m *MetaClient) RebuildDataCID(fileName string) error {
	// TODO
	return nil
}
