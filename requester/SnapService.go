package requester

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	fdir       = "fdir/"
	snapdir    = "snap/"
	urls       = "urls/"
	filejson   = "file.json"
	taskid     = "task.id"
	fileid     = "file.id"
	urlsid     = "urls.id"
	TYPE_FIXED = "fixed"
	TYPE_RAND  = "rand"
	TYPE_FILE  = "file"
)

var fileLock *sync.Mutex = new(sync.Mutex)

//缓存文件
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

	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

}

//客户端上传到服务端的文件
type FileInfo struct {
	//文件名称
	Name string `json:"Name"`
	//文件大小
	Size int64 `json:"Size"`
	//文件 id 系统内部使用
	Fid int `json:"Fid"`
	//备注信息
	Info string `json:"Info"`
	//文件类型
	Ext string `json:"Ext"`
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
	//测试结果状态 0 正在运行
	Status int
	//Err提示信息
	Err string
	//请求头信息
	QH map[string]string

	//查询参数
	P map[string]interface{}

	//Payload body 请求体
	Payload string

	//开始时间
	StartTime int64

	//结束时间
	EndTime int64

	//持续时间 EndTime-StartTime
	Duration int64

	//ReqNums
	ReqNums int

	//multipart/data file name
	InputFileName string
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
	Status   int
	Hint     string
	Datas    []interface{}
}

func getId(filePath string) int {

	idfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("read %s occur err %s \n", filePath, err.Error())
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

func getAllIds(filePath string, curPage, pageSize int) (tids []int, total int) {
	idfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	tids = make([]int, 0)
	total = 0
	if err != nil {
		log.Printf("read %s ocuur err %s \n", filePath, err.Error())
		return tids, total
	}
	defer idfile.Close()
	bs, err := ioutil.ReadAll(idfile)
	if err != nil {
		log.Println(err.Error())
		return tids, total
	}
	json.Unmarshal(bs, &tids)
	total = len(tids)
	if curPage <= 0 {
		curPage = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPage := (total + pageSize - 1) / pageSize
	if curPage > totalPage {
		return []int{}, total
	}
	start := (curPage - 1) * pageSize
	end := curPage * pageSize
	if end > total {
		end = total
	}

	sort.Sort(sort.Reverse(sort.IntSlice(tids)))

	return tids[start:end], total
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
	datas := make([]string, 10)

	log.Printf("first read urlfileid %s \n", urlfileId)
	fileBs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error())
		return defaulturl
	}
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

func SaveTestParam(testParam *TestParam) {
	bs, err := json.Marshal(testParam)
	if err != nil {
		log.Println(err.Error())
		return
	}
	SaveFileString(snapdir, strconv.Itoa(testParam.TaskId)+".p", bs)
}

func getDataInfo(dir, fileSuffix string, id int) *map[string]interface{} {
	data := make(map[string]interface{})
	bs, err := ioutil.ReadFile(dir + strconv.Itoa(id) + fileSuffix)
	if err != nil {
		log.Println(err.Error())
		return &data
	}
	json.Unmarshal(bs, &data)
	return &data
}

func GetFileInfo(fid string) *FileInfo {
	fileInfo := &FileInfo{}
	bs, err := ioutil.ReadFile(fdir + fid + ".info_")
	if err != nil {
		log.Println(err.Error())
		return fileInfo
	}
	json.Unmarshal(bs, fileInfo)
	return fileInfo
}

//更新测试结果状态
func UpdateFinshDataInfo(taskId int, reqNums int64) {
	result := *getDataInfo(snapdir, ".p", taskId)
	endTime := time.Now()
	result["EndTime"] = endTime.UnixNano()
	result["Duration"] = endTime.UnixNano() - int64((result["StartTime"]).(float64))
	result["ReqNums"] = reqNums
	result["Status"] = 1
	result["TaskId"] = taskId
	UpdateDataInfo(&result)
}
func UpdateDataInfo(dataInfo *map[string]interface{}) {

	taskId := (*dataInfo)["TaskId"].(int)
	filePath := snapdir + strconv.Itoa(taskId) + ".p"
	idfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("update %s occur err %s \n", filePath, err.Error())
		return
	}
	defer idfile.Close()
	idfile.Truncate(0)
	idfile.Seek(0, 0)
	jsonBs, jsonErr := json.Marshal(dataInfo)
	if jsonErr != nil {
		log.Printf("json.Marshal is err %s \n", jsonErr.Error())
		return
	}
	idfile.WriteString(string(string(jsonBs)))
}

func (pageInfo *PageInfo) getPageSize() int {
	if pageInfo.PageSize <= 0 {
		return 10
	}
	return pageInfo.PageSize
}

func (pageInfo *PageInfo) getCurPage() int {
	if pageInfo.CurPage <= 0 {
		return 1
	}
	return pageInfo.CurPage
}
func (pageInfo *PageInfo) QueryPageInfo(dir, idfile, fileSuffix string) {
	curPage := pageInfo.getCurPage()
	pageSize := pageInfo.getPageSize()

	ids, total := getAllIds(dir+idfile, curPage, pageSize)
	pageInfo.Count = total
	pageInfo.Total = (total + pageSize - 1) / pageSize
	Datas := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		Datas = append(Datas, getDataInfo(dir, fileSuffix, id))
	}
	log.Printf("Datas len=%d cap=%d \n", len(Datas), cap(Datas))
	pageInfo.Datas = Datas
}

func GetTaskPage(curPage, pageSize int) (pageInfo *PageInfo) {
	pageInfo = &PageInfo{
		CurPage:  curPage,
		PageSize: pageSize,
	}
	pageInfo.QueryPageInfo(snapdir, taskid, ".p")
	return pageInfo
}

//获取测试结果
func GetTaskSnap(taskId string) *map[string]interface{} {
	snap := make(map[string]interface{})
	bs, err := ioutil.ReadFile(snapdir + taskId + ".snap")
	snapPointer := &snap
	if err != nil {
		log.Println(err.Error())
		return snapPointer
	}
	log.Printf("bs len = %d \n", len(bs))
	json.Unmarshal(bs, snapPointer)
	return snapPointer
}
