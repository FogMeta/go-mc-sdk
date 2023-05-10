package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	aria2AddURI = "aria2.addUri"
	aria2Status = "aria2.tellStatus"
)

type Aria2Client struct {
	token     string
	serverUrl string
}

type Aria2Payload struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Aria2DownloadOption struct {
	Out string `json:"out"`
	Dir string `json:"dir"`
}

type Aria2Download struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Error   *Aria2Error `json:"error"`
	Gid     string      `json:"result"`
}

type Aria2Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Aria2StatusResult struct {
	Bitfield        string                  `json:"bitfield"`
	CompletedLength string                  `json:"completedLength"`
	Connections     string                  `json:"connections"`
	Dir             string                  `json:"dir"`
	DownloadSpeed   string                  `json:"downloadSpeed"`
	ErrorCode       string                  `json:"errorCode"`
	ErrorMessage    string                  `json:"errorMessage"`
	Gid             string                  `json:"gid"`
	NumPieces       string                  `json:"numPieces"`
	PieceLength     string                  `json:"pieceLength"`
	Status          string                  `json:"status"`
	TotalLength     string                  `json:"totalLength"`
	UploadLength    string                  `json:"uploadLength"`
	UploadSpeed     string                  `json:"uploadSpeed"`
	Files           []Aria2StatusResultFile `json:"files"`
}

type Aria2StatusResultFile struct {
	CompletedLength string                     `json:"completedLength"`
	Index           string                     `json:"index"`
	Length          string                     `json:"length"`
	Path            string                     `json:"path"`
	Selected        string                     `json:"selected"`
	Uris            []Aria2StatusResultFileUri `json:"uris"`
}

type Aria2StatusResultFileUri struct {
	Status string `json:"status"`
	Uri    string `json:"uri"`
}

func NewAria2Client(aria2Host, aria2Secret string, aria2Port int) *Aria2Client {
	return &Aria2Client{
		token:     aria2Secret,
		serverUrl: fmt.Sprintf("http://%s:%d/jsonrpc", aria2Host, aria2Port),
	}
}

func (aria2Client *Aria2Client) DownloadFile(uri string, outDir, outFilename string) *Aria2Download {
	payload := aria2Client.GenPayload4Download(aria2AddURI, uri, outDir, outFilename)
	response, err := httpRequest(http.MethodPost, aria2Client.serverUrl, "", payload, nil)
	if err != nil {
		return nil
	}

	var aria2Download Aria2Download
	err = json.Unmarshal(response, &aria2Download)
	if err != nil {
		return nil
	}
	return &aria2Download
}

func (aria2Client *Aria2Client) GenPayload4Download(method string, uri string, outDir, outFilename string) *Aria2Payload {
	options := Aria2DownloadOption{
		Out: outFilename,
		Dir: outDir,
	}
	return &Aria2Payload{
		JsonRpc: "2.0",
		Id:      uri,
		Method:  method,
		Params:  []interface{}{"token:" + aria2Client.token, []string{uri}, options},
	}
}

func httpRequest(httpMethod, uri, tokenString string, params interface{}, timeoutSecond *int) (body []byte, err error) {
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

	if len(strings.Trim(tokenString, " ")) > 0 {
		request.Header.Set("Authorization", "Bearer "+tokenString)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Transport: customTransport}
	if timeoutSecond != nil {
		client.Timeout = time.Duration(*timeoutSecond) * time.Second
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %s, code:%d, url:%s", response.Status, response.StatusCode, uri)
	}

	return io.ReadAll(response.Body)
}
