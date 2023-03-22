package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	web "github.com/filswan/go-swan-lib/client/web"
	utils "github.com/filswan/go-swan-lib/utils"
	shell "github.com/ipfs/go-ipfs-api"
	"os"
	"path/filepath"
)

type MetaClient struct {
	ApiKey string
	ApiSec string
	sh     *shell.Shell
}

func NewAPIClient() *MetaClient {
	c := &MetaClient{}

	return c
}

func (m *MetaClient) UploadFile(filePath string) (string, error) {
	// Creates an IPFS Shell client.
	sh := shell.NewShell("localhost:5001")

	// Iterates through the specified directory and uploads all files to IPFS.
	dirIter := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If the current path points to a file, upload the file to IPFS.
		if !info.IsDir() {
			fmt.Printf("Uploading %s...", path)
			cid, err := sh.Add(context.Background(), path)
			fmt.Printf("Done\n")
			if err != nil {
				return err
			}
			fmt.Printf("Added to IPFS: %s\n", cid)
		}
		return nil
	})

	// Start the iteration.
	if err := filepath.Walk(dirPath, dirIter); err != nil {
		return "", err
	}

	// Finally, upload the whole directory to IPFS.
	fmt.Printf("Uploading %s...", dirPath)
	dirCid, err := sh.AddDir(dirPath)
	fmt.Printf("Done\n")
	if err != nil {
		return "", err
	}
	fmt.Printf("Added to IPFS: %s\n", dirCid)

	// Returns the CID for the uploaded directory.
	return dirCid, nil
}

func (m *MetaClient) DownloadFile(hash, filePath string) error {
	url := ""
	outDir := ""
	outFilename := ""
	if !DownloadByAria2(url, outDir, outFilename) {
		fmt.Printf("Download Hash (%s) failed\n", hash)
		return errors.New("download files failed")
	}
	return nil
}

func (m *MetaClient) GetFileLists(page, limit uint64, showStorageInfo bool) ([]FileDetails, error) {
	params := FileListsParams{page, limit, showStorageInfo}

	// TODO:
	targetUrl := utils.UrlJoin("", "file_lists")
	response, err := web.HttpGetNoToken(targetUrl, params)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := FileListsResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return res.FileLists, nil
}

func (m *MetaClient) GetFileDataCID(fileName string) (string, error) {
	params := FileDataCIDParams{fileName}

	// TODO:
	targetUrl := utils.UrlJoin("", "file_data_cid")
	response, err := web.HttpGetNoToken(targetUrl, params)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return "", err
	}

	res := FileDataCIDResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) Error: %s \n", response, err)
		return "", err
	}

	return res.DataCid, nil
}

func (m *MetaClient) GetFileStatus(fileName string) (*FileDetails, error) {
	params := FileStatusParams{fileName}

	// TODO:
	targetUrl := utils.UrlJoin("", "file_status")
	response, err := web.HttpGetNoToken(targetUrl, params)
	if err != nil {
		fmt.Printf("Get Response Error: %s \n", err)
		return nil, err
	}

	res := FileStatusResponse{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("Parse Response (%s) Error: %s \n", response, err)
		return nil, err
	}

	return &(res.Status), nil
}
