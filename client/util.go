package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func PathJoin(root string, parts ...string) string {
	url := root

	for _, part := range parts {
		url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(part, "/")
	}
	url = strings.TrimRight(url, "/")

	return url
}

func GetIpfsCidInfo(ipfsApiUrl string, ipfsCid string) (IpfsCidInfo, error) {

	info := IpfsCidInfo{IpfsCid: ipfsCid}

	sh := shell.NewShell(ipfsApiUrl)
	stat, err := sh.FilesStat(context.Background(), PathJoin("/ipfs/", ipfsCid))
	if err != nil {
		logs.GetLogger().Error(err)
		return info, err
	}
	info.DataSize = int64(stat.CumulativeSize)
	info.IsDirectory = false
	if stat.Type == "directory" {
		info.IsDirectory = true
	}

	return info, nil
}

func downloadFileByAria2(conf *Aria2Conf, downUrl, outPath string) error {
	aria2 := client.GetAria2Client(conf.Host, conf.Secret, conf.Port)
	outDir := filepath.Dir(outPath)
	fileName := filepath.Base(outPath)
	logs.GetLogger().Infof("start download by aria2, downUrl:%s, outDir:%s, fileName:%s", downUrl, outDir, fileName)
	aria2Download := aria2.DownloadFile(downUrl, outDir, fileName)
	if aria2Download == nil {
		logs.GetLogger().Error("no response when asking aria2 to download")
		return errors.New("no response when asking aria2 to download")
	}

	if aria2Download.Error != nil {
		logs.GetLogger().Error(aria2Download.Error.Message)
		return errors.New(aria2Download.Error.Message)
	}

	if aria2Download.Gid == "" {
		logs.GetLogger().Error("no gid returned when asking aria2 to download")
		return errors.New("no gid returned when asking aria2 to download")
	}

	logs.GetLogger().Info("can check download status by gid:", aria2Download.Gid)
	return nil
}

func httpPost(uri, key, token string, params interface{}) ([]byte, error) {
	response, err := HttpRequestWithKey(http.MethodPost, uri, key, token, params)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return response, nil
}

const HTTP_CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const HTTP_CONTENT_TYPE_JSON = "application/json; charset=UTF-8"

func HttpRequestWithKey(httpMethod, uri, key, token string, params interface{}) ([]byte, error) {
	var request *http.Request
	var err error

	switch params := params.(type) {
	case io.Reader:
		request, err = http.NewRequest(httpMethod, uri, params)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		request.Header.Set("Content-Type", HTTP_CONTENT_TYPE_FORM)
	default:
		jsonReq, errJson := json.Marshal(params)
		if errJson != nil {
			logs.GetLogger().Error(errJson)
			return nil, errJson
		}

		request, err = http.NewRequest(httpMethod, uri, bytes.NewBuffer(jsonReq))
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		request.Header.Set("Content-Type", HTTP_CONTENT_TYPE_JSON)
	}

	if len(strings.Trim(key, " ")) > 0 {
		request.Header.Set("api_key", key)
	}

	if len(strings.Trim(token, " ")) > 0 {
		request.Header.Set("api_token", token)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status: %s, code:%d, url:%s", response.Status, response.StatusCode, uri)
		logs.GetLogger().Error(err)
		switch response.StatusCode {
		case http.StatusNotFound:
			logs.GetLogger().Error("please check your url:", uri)
		case http.StatusUnauthorized:
			logs.GetLogger().Error("Please check your key:", key, " and token:", token)
		}
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return responseBody, nil
}

func isFile(dirFullPath string) (*bool, error) {
	fi, err := os.Stat(dirFullPath)

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		isFile := false
		return &isFile, nil
	case mode.IsRegular():
		isFile := true
		return &isFile, nil
	default:
		err := fmt.Errorf("unknown path type")
		logs.GetLogger().Error(err)
		return nil, err
	}
}

func dirSize(path string) int64 {
	var size int64
	entrys, err := os.ReadDir(path)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0
	}
	for _, entry := range entrys {
		if entry.IsDir() {
			size += dirSize(filepath.Join(path, entry.Name()))
		} else {
			info, err := entry.Info()
			if err == nil {
				size += info.Size()
			}
		}
	}
	return size
}

func walkDirSize(path string) int64 {
	var totalSize int64
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileInfo, err := os.Stat(path)
			if err == nil {
				fileSize := fileInfo.Size()
				totalSize += fileSize
			}
		}
		return nil
	})
	return totalSize
}

func uploadFileToIpfs(sh *shell.Shell, fileName string) (string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}
	defer file.Close()

	ipfsCid, err := sh.Add(file)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	destPath := "/"
	srcPath := PathJoin("/ipfs/", ipfsCid)
	err = sh.FilesCp(context.Background(), srcPath, destPath)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	return ipfsCid, nil
}

func uploadDirToIpfs(sh *shell.Shell, dirName string) (string, error) {

	ipfsCid, err := sh.AddDir(dirName)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	destPath := "/"
	srcPath := PathJoin("/ipfs/", ipfsCid)
	err = sh.FilesCp(context.Background(), srcPath, destPath)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	return ipfsCid, nil
}

func ipfsCidIsDir(sh *shell.Shell, ipfsCid string) (*bool, error) {

	path := PathJoin("/ipfs/", ipfsCid)
	stat, err := sh.FilesStat(context.Background(), path)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	logs.GetLogger().Debug("FileStat:", stat)

	isFile := false
	if stat.Type == "directory" {
		isFile = true
	}

	return &isFile, nil
}

func downloadFromIpfs(sh *shell.Shell, ipfsCid, outDir string) error {
	return sh.Get(ipfsCid, outDir)
}
