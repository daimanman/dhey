package requester

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
)

const (
	fdir     = "fdir/"
	snapdir  = "snap/"
	urls     = "urls/"
	filejson = "file.json"
	taskid   = "task.id"
	fileid   = "file.id"
	urlsid   = "urls.id"
)

var URLFILE_MAP_LIST map[string][]string

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func createDirAndIdFile(dir string, idfile string) {
	if !PathExists(dir) {
		err := os.Mkdir(dir, 0777)
		if err != nil {
			log.Printf("dir=%s idfile=%s %s \n", dir, idfile, err.Error())
		}
	}

	if !PathExists(dir + idfile) {
		file, err := os.Create(dir + idfile)
		if err != nil {
			log.Println(dir, " ", idfile, err.Error())
		} else {
			defer file.Close()
		}
	}
}

func init() {
	log.Printf(" init dir and idfile ")
	createDirAndIdFile(fdir, fileid)
	createDirAndIdFile(snapdir, taskid)
	createDirAndIdFile(urls, urlsid)

	//url文件id对应的请求列表
	URLFILE_MAP_LIST = make(map[string][]string)
}

//客户端上传到服务端的文件
type FileInfo struct {
	//文件名称
	Name string
	//文件大小
	Size int64
	//文件 id 系统内部使用
	Fid int
	//备注信息
	Info string
	//文件类型
	Ext string
}

// 压力测试参数
type TestParam struct {
	//并发
	C int
	//持续时间
	Z string
	//请求总量 Z > 0 式 此参数默认为最大整数
	N int
	//测试备注
	Remark string
	//压测的目标URL
	Url string
	//访问方法
	Method string
	//任务id
	TaskId int
	//测试类型TYPE RAND:从文件中随机获取url FIXED:固定请求url FILE:提交文件
	Type string
	//File 上传的文件id
	FileId string
	//测试结果状态
	Status int
	//Err提示信息
	Err string
	//请求头信息
	QH map[string]string
}

type UrlParam struct {
	Remark string
	Datas  []string
}

type PageInfo struct {
	CurPage  int
	PageSize int
	Total    int
	Count    int
	Datas    []interface{}
}

func getId(filePath string) int {
	idfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("read %s ocuur err %s \n", filePath, err.Error())
		return -1
	}
	defer idfile.Close()
	var tid int
	tids := make([]int, 0)
	bs, _ := ioutil.ReadAll(idfile)
	json.Unmarshal(bs, &tids)
	if len(tids) == 0 {
		tid = 1
	} else {
		tid = tids[len(tids)-1] + 1
	}
	tids = append(tids, tid)
	idfile.Truncate(0)
	idfile.Seek(0, 0)
	newJsonBs, _ := json.Marshal(tids)
	idfile.WriteString(string(newJsonBs))
	log.Printf("%s id is %d \n", filePath, tid)
	return tid
}
func getFilePath(dir string, filename string) string {
	return dir + filename
}

//生成任务id
func GetTaskId() int {
	return getId(snapdir + taskid)
}

//上产文件的id
func GetFileId() int {
	return getId(fdir + fileid)
}

//随机请求列表文件的id
func GetUrlId() int {
	return getId(urls + urlsid)
}

func SaveFileInfo(fileinfo *FileInfo) {
	finfo, err := os.Create(getFilePath(fdir, strconv.Itoa(fileinfo.Fid)+".info_"))
	if err == nil {
		defer finfo.Close()
		bs, _ := json.Marshal(fileinfo)
		finfo.WriteString(string(bs))
	}
}

//保存上传的文件 multipart源文件
func SaveFile(file io.Reader, filename string, filesize int64, info string) *FileInfo {
	fileid := GetFileId()
	fileinfo := &FileInfo{
		Name: filename,
		Size: filesize,
		Fid:  fileid,
		Info: info,
		Ext:  path.Ext(filename),
	}
	newFileName := getFilePath(fdir, strconv.Itoa(fileinfo.Fid)+fileinfo.Ext)
	nf, err := os.Create(newFileName)
	if err != nil {
		log.Println(err.Error())
	} else {
		defer nf.Close()
	}
	//拷贝文件
	io.Copy(nf, file)

	//存储文件信息
	SaveFileInfo(fileinfo)

	return fileinfo
}

func SaveFileString(dir string, filename string, b []byte) error {
	file, err := os.Create(dir + filename)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer file.Close()
	file.WriteString(string(b))
	return nil

}

//保存请求列表 json参数格式
func SaveUrlListInfo(r io.Reader) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
		return err
	}

	urlParam := &UrlParam{}
	json.Unmarshal(bs, urlParam)
	remark := urlParam.Remark
	datas := urlParam.Datas
	size := len(datas)
	urlid := GetUrlId()
	info := make(map[string]interface{})
	info["remark"] = remark
	info["size"] = size
	info["urlid"] = urlid

	infoBs, jsonInfoErr := json.Marshal(info)
	if jsonInfoErr != nil {
		log.Println(jsonInfoErr.Error())
		return jsonInfoErr
	}

	dataBs, jsonDataErr := json.Marshal(datas)
	if jsonDataErr != nil {
		log.Println(jsonDataErr.Error())
		return jsonDataErr
	}
	urlidStr := strconv.Itoa(urlid)
	SaveFileString(urls, urlidStr+".info", infoBs)
	SaveFileString(urls, urlidStr+".data", dataBs)
	return nil
}

func SaveSnapInfo(report Report) {
	reportbs, err := json.Marshal(report)
	if err != nil {
		log.Println(err.Error())
		return
	}
	SaveFileString(snapdir, strconv.Itoa(report.TaskId)+".snap", reportbs)
}

//获取随机的url 如果解析随机文件失败则返回defaulturl
func GetRandUrlFormUrlFile(urlfileId string, defaulturl string) string {
	filename := urls + urlfileId + ".data"
	if !PathExists(filename) {
		log.Println(filename, " is not found")
		return defaulturl
	}
	urlDatas := URLFILE_MAP_LIST[urlfileId]
	urlLength := len(urlDatas)
	if urlLength > 0 {
		return urlDatas[rand.Intn(urlLength-1)]
	}

	log.Printf("first read urlfileid %s \n", urlfileId)
	fileBs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error())
		return defaulturl
	}

	datas := make([]string, 10)
	jsonErr := json.Unmarshal(fileBs, &datas)

	if jsonErr != nil {
		log.Println(jsonErr.Error())
		return defaulturl
	}
	if len(datas) == 0 {
		return defaulturl
	}
	URLFILE_MAP_LIST[urlfileId] = datas
	return datas[0]
}
