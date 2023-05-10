package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func PathJoin(url string, parts ...string) string {
	for _, part := range parts {
		url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(part, "/")
	}
	return strings.TrimRight(url, "/")
}

func GetIpfsCidInfo(ipfsApiUrl string, ipfsCid string) (IpfsCidInfo, error) {
	info := IpfsCidInfo{IpfsCid: ipfsCid}
	sh := shell.NewShell(ipfsApiUrl)
	stat, err := sh.FilesStat(context.Background(), PathJoin("/ipfs/", ipfsCid))
	if err != nil {
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
	aria2 := NewAria2Client(conf.Host, conf.Secret, conf.Port)
	outDir := filepath.Dir(outPath)
	fileName := filepath.Base(outPath)
	aria2Download := aria2.DownloadFile(downUrl, outDir, fileName)
	if aria2Download == nil {
		return errors.New("no response when asking aria2 to download")
	}

	if aria2Download.Error != nil {
		return errors.New(aria2Download.Error.Message)
	}

	if aria2Download.Gid == "" {
		return errors.New("no gid returned when asking aria2 to download")
	}

	return nil
}

const (
	contentTypeForm = "application/x-www-form-urlencoded"
	contentTypeJson = "application/json; charset=UTF-8"
)

func httpRequestWithKey(httpMethod, uri, key, token string, params interface{}) (body []byte, err error) {
	var request *http.Request

	switch params := params.(type) {
	case io.Reader:
		request, err = http.NewRequest(httpMethod, uri, params)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", contentTypeForm)
	default:
		jsonReq, errJson := json.Marshal(params)
		if errJson != nil {
			return nil, errJson
		}

		request, err = http.NewRequest(httpMethod, uri, bytes.NewBuffer(jsonReq))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", contentTypeJson)
	}

	if len(strings.Trim(key, " ")) > 0 {
		request.Header.Set("api-key", key)
	}

	if len(strings.Trim(token, " ")) > 0 {
		request.Header.Set("api-token", token)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status: %s, code:%d, url:%s", response.Status, response.StatusCode, uri)
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func uploadFileToIpfs(sh *shell.Shell, fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ipfsCid, err := sh.Add(file)
	if err != nil {
		return "", err
	}

	if err = sh.FilesCp(context.Background(), PathJoin("/ipfs/", ipfsCid), "/"); err != nil {
		return "", err
	}
	return ipfsCid, nil
}

func uploadDirToIpfs(sh *shell.Shell, dirName string) (string, error) {
	ipfsCid, err := sh.AddDir(dirName)
	if err != nil {
		return "", err
	}

	if err = sh.FilesCp(context.Background(), PathJoin("/ipfs/", ipfsCid), "/"); err != nil {
		return "", err
	}
	return ipfsCid, nil
}
