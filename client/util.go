package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/logs"
	shell "github.com/ipfs/go-ipfs-api"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	response, err := web.HttpRequestWithKey(http.MethodPost, uri, key, token, params)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return response, nil
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

func calcDirSize(dirPath string) int64 {
	var size int64
	entrys, err := os.ReadDir(dirPath)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0
	}
	for _, entry := range entrys {
		if entry.IsDir() {
			size += calcDirSize(filepath.Join(dirPath, entry.Name()))
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

func NewNode(hash, path, name string, size uint64, dir bool) *TreeNode {
	return &TreeNode{
		Hash:  hash,
		Path:  path,
		Name:  name,
		Size:  size,
		Dir:   dir,
		Child: []*TreeNode{},
	}
}

func NewNodeByDataCid(sh *shell.Shell, dataCid string, nodePath, name string) *TreeNode {
	path := PathJoin("/ipfs/", dataCid)
	stat, err := sh.FilesStat(context.Background(), path)
	if err != nil {
		logs.GetLogger().Error(dataCid, " get dag directory info err:", err)
		return nil
	}

	if stat.Type == "directory" {
		return NewNode(dataCid, PathJoin(nodePath, dataCid), name, stat.CumulativeSize, true)
	} else if stat.Type == "file" {
		return NewNode(dataCid, PathJoin(nodePath, dataCid), name, stat.CumulativeSize, false)
	} else {
		logs.GetLogger().Warn("unknown type in build node: ", stat.Type)
	}

	return nil
}

func (n *TreeNode) AddChild(node *TreeNode) error {

	if n.Child != nil {
		n.Child = append(n.Child, node)
	}

	return nil
}

func (n *TreeNode) BuildChildTree(sh *shell.Shell) error {
	if !n.Dir || len(n.Child) == 0 {
		return nil
	}

	for _, child := range n.Child {
		if !child.Dir {
			continue
		}

		resp := DagGetResponse{}
		if err := sh.DagGet(child.Hash, &resp); err != nil {
			logs.GetLogger().Error(child.Hash, " get dag directory info err:", err)
			continue
		}

		// build all subChild
		for _, link := range resp.Links {
			subChild := NewNodeByDataCid(sh, link.Hash.Target, child.Path, link.Name)
			if subChild == nil {
				continue
			}

			child.AddChild(subChild)
		}

		child.BuildChildTree(sh)
	}

	return nil
}

func (n *TreeNode) ReduceChildTree() error {
	if n.Path != "/" {
		logs.GetLogger().Error("reduce child must from root")
		return errors.New("reduce child must from root")
	}

	for _, child := range n.Child {

		hash := child.Hash
		child.IsTop = false

		funded := false
		for _, reduceChild := range n.Child {

			// First exclude its own nodes.
			if hash == reduceChild.Hash {
				continue
			}

			n := reduceChild.Find(hash)
			if n != nil {
				funded = true
				//n.Show()
			}
		}

		if !funded {
			child.IsTop = true
		}

	}

	return nil
}

func (n *TreeNode) SortChild() error {
	sort.Sort(TreeNodeDecrement(n.Child))
	return nil
}

func (n *TreeNode) Insert(hash string, node *TreeNode) error {

	prev := n.Find(hash)
	if prev != nil {
		prev.Child = append(prev.Child, node)
	}

	return nil
}

func (n *TreeNode) Del(hash string) error {
	//TODO:
	return nil
}

func (n *TreeNode) Find(hash string) *TreeNode {
	if n.Hash == hash {
		return n
	}

	for _, child := range n.Child {
		if fn := child.Find(hash); fn != nil {
			return fn
		}
	}

	return nil
}

func (n *TreeNode) PrintAll() error {

	n.Print()

	for _, child := range n.Child {
		if !child.Dir {
			child.PrintAll()
		}
	}

	for _, child := range n.Child {
		if child.Dir {
			child.PrintAll()
		}
	}

	return nil
}

func (n *TreeNode) PrintAllTop() error {
	n.Print()

	for _, child := range n.Child {
		if child.IsTop {
			child.PrintAll()
		}
	}

	return nil
}

func (n *TreeNode) Print() error {
	//logs.GetLogger().Infof("TreeNode: hash=%s, path=%s, name=%s, size=%d, deep=%d, dir=%t, child-num=%d",
	//	n.Hash, n.Path, n.Name, n.Size, n.Deep, n.Dir, len(n.Child))
	fmt.Print("\n")
	if n.Path == "/" {
		fmt.Printf("/")
		return nil
	}

	count := strings.Count(n.Path, "/")
	for i := 0; i < count-1; i++ {
		fmt.Print("    ")
	}
	fmt.Printf("|---%s (Hash:%s Size:%d)", n.Name, n.Hash, n.Size)

	return nil
}

func (n *TreeNode) Show() error {
	logs.GetLogger().Infof("TreeNode: hash=%s, path=%s, name=%s, size=%d, group=%d, dir=%t, child-num=%d",
		n.Hash, n.Path, n.Name, n.Size, n.Group, n.Dir, len(n.Child))
	return nil
}

type TreeNodeDecrement []*TreeNode

func (s TreeNodeDecrement) Len() int           { return len(s) }
func (s TreeNodeDecrement) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TreeNodeDecrement) Less(i, j int) bool { return s[i].Size > s[j].Size }

type FileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

type DataSet struct {
	Index int
	Name  string
	Path  string
	Size  int64
	Items []FileInfo
}

func (m *MetaClient) buildDirectoryTree(ipfsApiUrl string, dataCid string) error {
	// Creates an IPFS Shell client.
	sh := shell.NewShell(ipfsApiUrl)

	pins, err := sh.Pins()
	if err != nil {
		logs.GetLogger().Errorf("List pins error: %s", err)
		return err
	}

	logs.GetLogger().Info(len(pins), " records to process ...")

	//build root node
	root := NewNode("/", "/", "/", 0, true)
	for hash, info := range pins {
		logs.GetLogger().Debug("Key:", hash, " Type:", info.Type)

		path := PathJoin("/ipfs/", hash)
		stat, err := sh.FilesStat(context.Background(), path)
		if err != nil {
			logs.GetLogger().Error(err)
			continue
		}
		logs.GetLogger().Debugf("FileStat:%+v", stat)

		if stat.Type == "directory" {
			resp := DagGetResponse{}
			sh.DagGet(hash, &resp)
			logs.GetLogger().Debugf("dag directory info resp:%+v", resp)

			n := NewNode(hash, PathJoin(root.Path, hash), hash, stat.CumulativeSize, true)
			root.AddChild(n)
			logs.GetLogger().Debugf("add node director of %s to root", hash)

		} else if stat.Type == "file" {
			n := NewNode(hash, PathJoin(root.Path, hash), hash, stat.CumulativeSize, false)
			root.AddChild(n)
			logs.GetLogger().Debugf("add node file of %s to root", hash)
		} else {
			logs.GetLogger().Warn("unknown type: ", stat.Type)
		}
	}

	root.BuildChildTree(sh)
	root.SortChild()
	root.PrintAll()
	fmt.Print("\n")

	root.ReduceChildTree()
	root.PrintAllTop()
	fmt.Print("\n")

	return nil
}

func GetFileInfoList(dirPath string) []FileInfo {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	infos := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			logs.GetLogger().Warnf("failed to get file %s info:%s", PathJoin(dirPath, entry.Name()), err)
			continue
		}

		dirSize := info.Size()
		if info.IsDir() {
			dirSize = calcDirSize(PathJoin(dirPath, info.Name()))
		}

		fileInfo := FileInfo{
			Name:  info.Name(),
			Size:  dirSize,
			IsDir: info.IsDir(),
		}

		infos = append(infos, fileInfo)
	}
	return infos
}

func GreedyDataSets(dirPath string, givenSize int64) []DataSet {

	entrys := GetFileInfoList(dirPath)
	if entrys == nil {
		log.Fatal("failed to get files information")
	}

	sort.Slice(entrys, func(i, j int) bool {
		return entrys[i].Size > entrys[j].Size
	})

	var groupDataSets []DataSet
	gIndex := 0
	curDataSet := DataSet{Index: gIndex, Size: 0, Path: dirPath}
	for _, entry := range entrys {
		if curDataSet.Size+entry.Size <= givenSize {
			curDataSet.Size += entry.Size
			curDataSet.Items = append(curDataSet.Items, entry)
		} else {
			if len(curDataSet.Items) != 0 {
				groupDataSets = append(groupDataSets, curDataSet)
				gIndex = gIndex + 1
				curDataSet = DataSet{Index: gIndex, Size: 0, Path: dirPath}
			}

			curDataSet.Size += entry.Size
			curDataSet.Items = append(curDataSet.Items, entry)
		}
	}

	if len(curDataSet.Items) > 0 {
		groupDataSets = append(groupDataSets, curDataSet)
	}

	return groupDataSets
}
