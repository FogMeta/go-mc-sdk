package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func DownloadByAria2(url, outDir, outFilename string) bool {

	aria2Download := aria2Client.DownloadFile(url, outDir, outFilename)
	if aria2Download == nil {
		UpdateStatusAndLog(deal, DEAL_STATUS_DOWNLOAD_FAILED, "no response when asking aria2 to download")
		return
	}

	if aria2Download.Error != nil {
		UpdateStatusAndLog(deal, DEAL_STATUS_DOWNLOAD_FAILED, aria2Download.Error.Message)
		return
	}

	if aria2Download.Gid == "" {
		UpdateStatusAndLog(deal, DEAL_STATUS_DOWNLOAD_FAILED, "no gid returned when asking aria2 to download")
		return
	}

	return true
}

func DownloadByIpfs(url, outDir, outFilename string) bool {
	// Creates an IPFS Shell client.
	sh := shell.NewShell("localhost:5001")

	// Determines if the CID represents a file or a directory.
	stats, err := sh.ObjectStat(context.Background(), cid)
	if err != nil {
		return fmt.Errorf("failed to get stats for %s: %v", cid, err)
	}

	// If it is a directory, use sh.Get to download the whole directory to a temporary directory.
	if stats.Type == "directory" {
		tmpDir, err := ioutil.TempDir("", "ipfs-download")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)
		err = sh.Get(context.Background(), cid, tmpDir)
		if err != nil {
			return fmt.Errorf("failed to download directory %s: %v", cid, err)
		}
		err = os.Rename(tmpDir, localPath)
		if err != nil {
			return fmt.Errorf("failed to move downloaded directory: %v", err)
		}
	} else if stats.Type == "file" {
		// If it is a file, use sh.Cat to read the content of the file and write to the local file.
		reader, err := sh.Cat(cid)
		if err != nil {
			return fmt.Errorf("failed to read content for %s: %v", cid, err)
		}
		defer reader.Close()
		file, err := os.Create(localPath)
		if err != nil {
			return fmt.Errorf("failed to create local file: %v", err)
		}
		defer file.Close()
		_, err = io.Copy(file, reader)
		if err != nil {
			return fmt.Errorf("failed to write local file: %v", err)
		}
	} else {
		return fmt.Errorf("%s is not a valid CID for a file or a directory", cid)
	}

	return true
}
