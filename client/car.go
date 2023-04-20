package client

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/codingsince1985/checksum"
	"github.com/filedrive-team/go-graphsplit"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	DIR_NAME_INPUT  = "input"
	DIR_NAME_OUTPUT = "output"

	JSON_FILE_NAME_CAR_UPLOAD = "car.json"
	CSV_FILE_NAME_CAR_UPLOAD  = "car.csv"
)

type CmdGoCar struct {
	OutputDir          string   //required
	Inputs             []string //required
	ParentPath         string   //required
	GraphName          string   //required
	GenerateMd5        bool     //required
	GocarFileSizeLimit int64    //required
	GocarFolderBased   bool     //required always true
	Parallel           int
}

type FileDesc struct {
	Uuid           string
	SourceFileName string
	SourceFilePath string
	SourceFileMd5  string
	SourceFileSize int64
	CarFileName    string
	CarFilePath    string
	CarFileMd5     string
	CarFileUrl     string
	CarFileSize    int64
	PayloadCid     string
	PieceCid       string
	StartEpoch     *int64
	SourceId       *int
}

func GetCmdGoCar(dataSet Group, outputDir *string, parallel int, carFileSizeLimit int64) *CmdGoCar {

	var inputs []string
	for _, fileInfo := range dataSet.Items {
		inputs = append(inputs, fileInfo.Name)
	}

	cmdGoCar := &CmdGoCar{
		Inputs:             inputs,
		ParentPath:         dataSet.Path,
		GraphName:          filepath.Base(dataSet.Path) + "-" + strconv.FormatInt(int64(dataSet.Index), 10),
		GocarFileSizeLimit: carFileSizeLimit,
		GenerateMd5:        false,
		GocarFolderBased:   true,
		Parallel:           parallel,
	}

	if !utils.IsStrEmpty(outputDir) {
		cmdGoCar.OutputDir = *outputDir
	} else {
		cmdGoCar.OutputDir = filepath.Join(*outputDir, time.Now().Format("2006-01-02_15:04:05")) + "_" + uuid.NewString()
	}

	return cmdGoCar
}

func CreateGoCarFilesByConfig(dataSet Group, outputDir *string, parallel int, carFileSizeLimit int64) ([]*FileDesc, error) {

	cmdGoCar := GetCmdGoCar(dataSet, outputDir, parallel, carFileSizeLimit)
	fileDescs, err := cmdGoCar.CreateGoCarFiles()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return fileDescs, nil
}

func (cmdGoCar *CmdGoCar) CreateGoCarFiles() ([]*FileDesc, error) {

	err := utils.CreateDirIfNotExists(cmdGoCar.OutputDir, DIR_NAME_OUTPUT)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	sliceSize := cmdGoCar.GocarFileSizeLimit
	if sliceSize <= 0 {
		err := fmt.Errorf("gocar file size limit is too smal")
		logs.GetLogger().Error(err)
		return nil, err
	}

	carDir := cmdGoCar.OutputDir
	Emptyctx := context.Background()
	cb := graphsplit.CommPCallback(carDir, false, false)

	parentPath := cmdGoCar.ParentPath
	targetPaths := cmdGoCar.Inputs
	graphName := cmdGoCar.GraphName
	parallel := cmdGoCar.Parallel
	logs.GetLogger().Info("Creating car file for ", parentPath, " ", graphName)
	err = graphsplit.ChunkMulti(Emptyctx, sliceSize, parentPath, targetPaths, carDir, graphName, parallel, cb)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	logs.GetLogger().Info("Car file for ", parentPath, " created")

	fileDescs, err := cmdGoCar.createFilesDescFromManifest()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	logs.GetLogger().Info(len(fileDescs), " car files have been created to directory:", carDir)
	logs.GetLogger().Info("Please upload car files to web server or ipfs server.")

	return fileDescs, nil
}

type ManifestDetail struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size int    `json:"Size"`
	Link []struct {
		Name string `json:"Name"`
		Hash string `json:"Hash"`
		Size int64  `json:"Size"`
	} `json:"Link"`
}

const (
	ColumnPayloadCID  = "payload_cid"
	ColumnFilename    = "filename"
	ColumnPieceCID    = "piece_cid"
	ColumnPieceSize   = "piece_size"
	ColumnDetail      = "detail"
	ColumnPayloadSize = "payload_size"
)

var colums = []string{ColumnPayloadCID, ColumnFilename, ColumnPieceCID, ColumnPieceSize, ColumnDetail}

func (cmdGoCar *CmdGoCar) createFilesDescFromManifest() ([]*FileDesc, error) {
	manifestFilename := "manifest.csv"

	fileFullPath := filepath.Join(cmdGoCar.OutputDir, manifestFilename)
	file, err := os.Open(fileFullPath)
	if err != nil {
		logs.GetLogger().Error("failed opening file: ", fileFullPath)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	//
	//lines, err := utils.ReadAllLines(cmdGoCar.OutputDir, manifestFilename)
	//if err != nil {
	//	logs.GetLogger().Error(err)
	//	return nil, err
	//}

	fileDescs := []*FileDesc{}
	colMap := make(map[string]int)
	for i, fields := range records {
		if i == 0 {
			for pos, col := range fields {
				colMap[col] = pos
			}

			for _, col := range colums {
				if _, ok := colMap[col]; !ok {
					return nil, fmt.Errorf("column %s not found", col)
				}
			}
			continue
		}
		fileDesc := FileDesc{}
		fileDesc.PayloadCid = fields[colMap[ColumnPayloadCID]]
		fileDesc.CarFileName = fileDesc.PayloadCid + ".car"
		fileDesc.CarFileUrl = fileDesc.CarFileName
		fileDesc.CarFilePath = filepath.Join(cmdGoCar.OutputDir, fileDesc.CarFileName)
		fileDesc.PieceCid = fields[colMap[ColumnPieceCID]]
		fileDesc.CarFileSize = utils.GetInt64FromStr(fields[colMap[ColumnPieceCID]])

		carFileDetail := fields[colMap[ColumnDetail]]

		manifestDetail := ManifestDetail{}
		err = json.Unmarshal([]byte(carFileDetail), &manifestDetail)
		if err != nil {
			logs.GetLogger().Error("Failed to parse: ", carFileDetail)
			return nil, err
		}

		fileDesc.SourceFileName = cmdGoCar.GraphName
		fileDesc.SourceFilePath = cmdGoCar.ParentPath
		for _, link := range manifestDetail.Link {
			fileDesc.SourceFileSize = fileDesc.SourceFileSize + link.Size
		}

		if cmdGoCar.GenerateMd5 {
			if utils.IsFileExistsFullPath(fileDesc.SourceFilePath) {
				srcFileMd5, err := checksum.MD5sum(fileDesc.SourceFilePath)
				if err != nil {
					logs.GetLogger().Error(err)
					return nil, err
				}
				fileDesc.SourceFileMd5 = srcFileMd5
			}

			carFileMd5, err := checksum.MD5sum(fileDesc.CarFilePath)
			if err != nil {
				logs.GetLogger().Error(err)
				return nil, err
			}
			fileDesc.CarFileMd5 = carFileMd5
		}

		fileDescs = append(fileDescs, &fileDesc)
	}

	_, err = WriteCarFilesToFiles(fileDescs, cmdGoCar.OutputDir, JSON_FILE_NAME_CAR_UPLOAD, CSV_FILE_NAME_CAR_UPLOAD)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return fileDescs, nil
}

func WriteCarFilesToFiles(carFiles []*FileDesc, outputDir, jsonFilename, csvFileName string) (*string, error) {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	jsonFilePath, err := WriteFileDescsToJsonFile(carFiles, outputDir, jsonFilename)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = WriteCarFilesToCsvFile(carFiles, outputDir, csvFileName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return jsonFilePath, nil
}

func WriteFileDescsToJsonFile(fileDescs []*FileDesc, outputDir, jsonFileName string) (*string, error) {
	jsonFilePath := filepath.Join(outputDir, jsonFileName)
	content, err := json.MarshalIndent(fileDescs, "", " ")
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = ioutil.WriteFile(jsonFilePath, content, 0644)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	logs.GetLogger().Info("Metadata json file generated: ", jsonFilePath)
	return &jsonFilePath, nil
}

func WriteCarFilesToCsvFile(carFiles []*FileDesc, outDir, csvFileName string) error {
	csvFilePath := filepath.Join(outDir, csvFileName)
	var headers []string
	headers = append(headers, "uuid")
	headers = append(headers, "source_file_name")
	headers = append(headers, "source_file_path")
	headers = append(headers, "source_file_md5")
	headers = append(headers, "source_file_size")
	headers = append(headers, "car_file_name")
	headers = append(headers, "car_file_path")
	headers = append(headers, "car_file_md5")
	headers = append(headers, "car_file_url")
	headers = append(headers, "car_file_size")
	headers = append(headers, "pay_load_cid")
	headers = append(headers, "piece_cid")
	headers = append(headers, "start_epoch")
	headers = append(headers, "source_id")
	headers = append(headers, "deals")

	file, err := os.Create(csvFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(headers)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	for _, carFile := range carFiles {
		var columns []string
		columns = append(columns, carFile.Uuid)
		columns = append(columns, carFile.SourceFileName)
		columns = append(columns, carFile.SourceFilePath)
		columns = append(columns, carFile.SourceFileMd5)
		columns = append(columns, strconv.FormatInt(carFile.SourceFileSize, 10))
		columns = append(columns, carFile.CarFileName)
		columns = append(columns, carFile.CarFilePath)
		columns = append(columns, carFile.CarFileMd5)
		columns = append(columns, carFile.CarFileUrl)
		columns = append(columns, strconv.FormatInt(carFile.CarFileSize, 10))
		columns = append(columns, carFile.PayloadCid)
		columns = append(columns, carFile.PieceCid)

		if carFile.StartEpoch != nil {
			columns = append(columns, strconv.FormatInt(*carFile.StartEpoch, 10))
		} else {
			columns = append(columns, "")
		}

		if carFile.SourceId != nil {
			columns = append(columns, strconv.Itoa(*carFile.SourceId))
		} else {
			columns = append(columns, "")
		}

		// no deals
		columns = append(columns, "")

		err = writer.Write(columns)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
	}

	logs.GetLogger().Info("Metadata csv generated: ", csvFilePath)

	return nil
}
