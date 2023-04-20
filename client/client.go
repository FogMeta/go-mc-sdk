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
		if downInfo[0].IsDirectory {
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
			if info.IsDirectory {
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

func (m *MetaClient) ReportMetaClientServer(datasetName string, ipfsData []IpfsData) error {
	var params []interface{}
	params = append(params, datasetName, ipfsData)
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
	logs.GetLogger().Infof("meta StoreSourceFile response: %+v", res)

	if res.Result.Code != "save_source_file_success" {
		return errors.New("failed message from meta server")
	}

	return nil
}

func (m *MetaClient) GetDatasetList(datasetName string, pageNum, size int) (*GetDatasetListPager, error) {
	var params []interface{}
	params = append(params, GetDatasetListReq{datasetName, pageNum, size})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDatasetList",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := GetDatasetListResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Infof("meta GetDatasetList response: %+v", res)

	//datasetList := res.Result.Data.DatasetList
	//for index, dataset := range datasetList {
	//	logs.GetLogger().Infof("Index: %d, Dataset: %+v", index, dataset)
	//
	//	ipfsList := dataset.IpfsList
	//	for i, ipfsDetail := range ipfsList {
	//		logs.GetLogger().Infof("IPFS Detail-%d: %+v", i, ipfsDetail)
	//	}
	//}

	return &res.Result.Data, nil
}

func (m *MetaClient) GetSourceFileInfo(ipfsCid string) ([]IpfsDataDetail, error) {
	var params []interface{}
	params = append(params, ipfsCid)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileInfo",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}
	res := GetSourceFileInfoResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Infof("meta GetSourceFileInfo response: %+v", res)

	return res.Result.Data, nil
}

func (m *MetaClient) GetSourceFileStatus(datasetName, ipfsCid string, pageNum, size int) (*GetSourceFileStatusPager, error) {

	var params []interface{}
	params = append(params, GetSourceFileStatusReq{datasetName, ipfsCid, pageNum, size})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetSourceFileStatus",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}

	res := GetSourceFileStatusResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Infof("meta GetSourceFileStatus response: %+v", res)

	sourceFileStatus := res.Result.Data
	logs.GetLogger().Infof("Source File Status: %+v", sourceFileStatus)

	carList := sourceFileStatus.CarList
	for _, carDetail := range carList {
		logs.GetLogger().Infof("CAR Detail: %+v", carDetail)
		providers := carDetail.StorageProviders
		for ii, provider := range providers {
			logs.GetLogger().Infof("Provider-%d: %+v", ii, provider)
		}
	}

	return &res.Result.Data, nil
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
	logs.GetLogger().Infof("meta GetDownloadFileInfoByIpfsCid response: %+v", res)

	return res.Result.Data, nil
}

func (m *MetaClient) GetDatasetsByGroupName(groupName string) ([]GetDatasetsByGroupNameResp, error) {

	var params []interface{}
	params = append(params, groupName)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.GetDatasetsByGroupName",
		Params:  params,
		Id:      1,
	}
	//logs.GetLogger().Infof("GetDatasetsByGroupName MetaUrl:%s ApiKey:%s ApiToken:%s Params:%+v", m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return nil, err
	}

	res := GetDatasetsByGroupNameResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return nil, err
	}
	logs.GetLogger().Infof("meta GetDatasetsByGroupNameResponse response: %+v", res)

	return res.Result.Data, nil
}

func (m *MetaClient) StoreSourceFileByGroup(groupName string, dataList [][]IpfsData) error {

	var params []interface{}
	params = append(params, groupName, dataList)
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.StoreSourceFileByGroup",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return err
	}

	res := StoreSourceFileByGroupResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return err
	}
	logs.GetLogger().Infof("meta StoreSourceFileByGroup response: %+v", res)

	if res.Result.Code != "success" {
		return errors.New("StoreSourceFileByGroup not success from meta server")
	}
	return nil
}

func (m *MetaClient) StoreCarFiles(datasetId int64, carList []*CarInfo) error {

	var params []interface{}
	params = append(params, StoreCarFilesReq{DatasetId: datasetId, CarList: carList})
	jsonRpcParams := JsonRpcParams{
		JsonRpc: "2.0",
		Method:  "meta.StoreCarFiles",
		Params:  params,
		Id:      1,
	}
	response, err := httpPost(m.MetaUrl, m.ApiKey, m.ApiToken, jsonRpcParams)
	if err != nil {
		logs.GetLogger().Errorf("Get Response Error: %s", err)
		return err
	}

	res := StoreCarFilesResponse{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		logs.GetLogger().Errorf("Parse Response (%s) Error: %s", response, err)
		return err
	}
	logs.GetLogger().Infof("meta StoreCarFiles response: %+v", res)

	if res.Result.Code != "success" {
		return errors.New("StoreCarFiles not success from meta server")
	}
	return nil
}

func (m *MetaClient) RebuildIpfsCid(fileName string) error {
	// TODO
	return nil
}

func (m *MetaClient) BackupIpfsData(taskName, repoDir, outputDir, apiUrl, gatewayUrl string, dataSetLimit, carSizeLimit int64, parallel int) error {

	var todoSets []GetDatasetsByGroupNameResp
	// query last task by task name
	todoSets, err := m.GetDatasetsByGroupName(taskName)
	if err != nil {
		logs.GetLogger().Error("failed to get remained datasets from meta server:", err)
		return err
	}

	if len(todoSets) == 0 {
		// split task to datasets
		logs.GetLogger().Infof("start to backup %s subdirectory of blocks and datastore", repoDir)

		var dateSets []DataSet
		blocksDataSets := GreedyDataSets(PathJoin(repoDir, "blocks"), dataSetLimit)
		datastoreDataSets := GreedyDataSets(PathJoin(repoDir, "datastore"), dataSetLimit)

		dateSets = append(dateSets, datastoreDataSets...)
		dateSets = append(dateSets, blocksDataSets...)

		logs.GetLogger().Infof("task datasets count is %d", len(dateSets))
		// update task include all datasets to meta server
		var dataSets [][]IpfsData
		for _, group := range dateSets {
			var dataSet []IpfsData
			for _, item := range group.Items {
				dataSet = append(dataSet, IpfsData{
					SourceName:  PathJoin(group.Path, item.Name),
					DataSize:    item.Size,
					IsDirectory: item.IsDir,
				})
			}
			dataSets = append(dataSets, dataSet)
		}

		err := m.StoreSourceFileByGroup(taskName, dataSets)
		if err != nil {
			logs.GetLogger().Error("failed to store source file by group:", err)
			return err
		}
		logs.GetLogger().Infof("has report store source file, datasets count is %d", len(dataSets))

		todoSets, err = m.GetDatasetsByGroupName(taskName)
		if err != nil {
			logs.GetLogger().Error("failed to get datasets from meta server:", err)
			return err
		}

	}
	logs.GetLogger().Info("start to do dataset count is: ", len(todoSets))
	for _, dataSet := range todoSets {
		if dataSet.DatasetStatus != "ReadyForCarUpload" {
			logs.GetLogger().Warn("skip one dataset due to status is: ", dataSet.DatasetStatus)
			continue
		}

		//
		logs.GetLogger().Infof("======== do dataset ======== %+v", dataSet)
		datasetListPager, err := m.GetDatasetList(dataSet.DatasetName, 0, 100000)
		if err != nil {
			logs.GetLogger().Error("failed to get dataset list from meta server and continue:", err)
			continue
		}

		if len(datasetListPager.DatasetList) != 1 {
			logs.GetLogger().Error("get dataset list count should only one from meta server and continue")
			continue
		}

		fileList := datasetListPager.DatasetList[0].IpfsList
		logs.GetLogger().Debugf("get files list count is %d", len(fileList))
		oneSet := DataSet{Size: 0, Name: dataSet.DatasetName}
		for _, fl := range fileList {
			oneSet.Items = append(oneSet.Items, FileInfo{
				Name:  fl.SourceName,
				Size:  fl.DataSize,
				IsDir: fl.IsDirectory,
			})
			oneSet.Size += fl.DataSize
		}

		if len(oneSet.Items) > 0 {
			oneSet.Path = filepath.Dir(oneSet.Items[0].Name)
		}

		// Create CAR
		carDir := PathJoin(outputDir, dataSet.DatasetName)
		logs.GetLogger().Debugf("to create car path:%s item count:%d, output:%s", oneSet.Path, len(oneSet.Items), carDir)
		fileDescs, err := CreateGoCarFilesByConfig(oneSet, &carDir, parallel, carSizeLimit)
		if err != nil {
			logs.GetLogger().Error("failed to creat car:", err)
			continue
		}

		logs.GetLogger().Infof("created CAR count is :%d", len(fileDescs))

		//break
		// update all CAR to ipfs
		sh := shell.NewShell(apiUrl)
		var carList []*CarInfo
		for i, desc := range fileDescs {
			logs.GetLogger().Debugf("File Desc Index %d: %+v", i, desc)
			ipfsCid, err := uploadFileToIpfs(sh, desc.CarFilePath)
			if err != nil {
				logs.GetLogger().Errorf("failed to upload CAR %s to IPFS: %s", desc.CarFilePath, err)
				continue
			}
			carInfo := &CarInfo{
				FileName:    desc.CarFileName,
				DataCid:     desc.PayloadCid,
				SourceSize:  desc.SourceFileSize,
				CarSize:     desc.CarFileSize,
				PieceCid:    desc.PieceCid,
				DownloadUrl: PathJoin(gatewayUrl, "ipfs/", ipfsCid),
			}
			carList = append(carList, carInfo)
			logs.GetLogger().Debugf("upload CAR to IPFS:%+v", carInfo)
			logs.GetLogger().Infof("has upload %s to IPFS, and IPFS CID is %s", carInfo.FileName, ipfsCid)
		}
		if len(carList) != len(fileDescs) {
			logs.GetLogger().Errorf("failed to upload dataset %s ALL CAR (%d/%d) to IPFS", dataSet.DatasetName, len(carList), len(fileDescs))
			continue
		}

		// report to meta server
		err = m.StoreCarFiles(dataSet.DatasetId, carList)
		if err != nil {
			logs.GetLogger().Error("failed to store CAR files to meta server:", err)
		}
		logs.GetLogger().Info("successfully process one dataset:", dataSet.DatasetName)

		// break

	}

	logs.GetLogger().Info("all datasets done!")

	return nil
}
