package client

import (
	"github.com/filswan/go-swan-lib/client"
)

func DownloadByAria2(url, outDir, outFilename string) bool {
	aria2Host := 0
	aria2Secret := ""
	aria2Port := ""

	aria2Client = client.GetAria2Client(aria2Host, aria2Secret, aria2Port)

	aria2Download := aria2Client.DownloadFile(url, outDir, outFilename)

	if aria2Download == nil {
		return false
	}

	if aria2Download.Error != nil {
		return false
	}

	return true
}
