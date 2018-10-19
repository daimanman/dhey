package requester

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	heyServer = &http.Server{
		Addr:           ":7654",
		Handler:        &MyServer{},
		ReadTimeout:    10 * time.Microsecond,
		WriteTimeout:   10 * time.Microsecond,
		MaxHeaderBytes: 1 << 30,
	}
	heyHandlerMap = make(map[string]HandlersFunc)
)

var mimeTypeExt map[string]string = map[string]string{
	".woff2": "application/x-font-woff",
	".woff":  "application/x-font-woff",
}

const (
	STATIC_PREFIX = "static/admin/"
)

type ResultMsg struct {
	Code int
	Msg  string
	Data interface{}
}

type MyServer struct {
}
type HandlersFunc func(http.ResponseWriter, *http.Request)

func (*MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Path
	h := heyHandlerMap[urlStr]
	if h != nil {
		h(w, r)
		return
	}
	SendStaticFile(w, r)
}
func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

//读取静态文件
func SendStaticFile(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimLeft(r.URL.Path, "/")
	if urlPath == "" {
		urlPath = "index.html"
	}
	ext := strings.ToLower(path.Ext(urlPath))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = mimeTypeExt[ext]
	}
	w.Header().Set("Content-Type", mimeType)

	file, err := os.Open(STATIC_PREFIX + urlPath)
	if err != nil {
		w.WriteHeader(404)
		io.Copy(w, bytes.NewReader(make([]byte, 1)))
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

//读取静态文件
func SendStaticFileTest(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimLeft(r.URL.Path, "/")

	ext := strings.ToLower(path.Ext(urlPath))
	mimeType := mime.TypeByExtension(ext)
	res := make(map[string]interface{})
	res["url"] = urlPath
	res["mimeType"] = mimeType
	res["ext"] = ext
	SendJson(w, res)
}

func SendJson(w http.ResponseWriter, data interface{}) {
	setJsonHeader(w)
	json.NewEncoder(w).Encode(data)
}

func IndexController(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	r.ParseForm()
	res["name"] = "dxm"
	res["code"] = 1
	filename := r.FormValue("filename")
	res["ext"] = path.Ext(filename)
	SendJson(w, res)
}

//接受上传的文件 做为提交Multipartt请求的源文件
func UpFileController(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	uploadFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		res["msg"] = err.Error()
		res["code"] = -1
		SendJson(w, res)
		return
	}
	defer uploadFile.Close()
	filename := fileHeader.Filename
	filesize := fileHeader.Size
	info := r.FormValue("info")
	fileinfo := SaveFile(uploadFile, filename, filesize, info)
	res["result"] = fileinfo
	SendJson(w, res)
}

//接受请求列表
func SaveUrlsController(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	res["code"] = 0
	res["msg"] = "OK"
	err := SaveUrlListInfo(r.Body)
	if err != nil {
		res["msg"] = err.Error()
		res["code"] = 1
	}
	SendJson(w, res)
}

//开启一个测试任务
func SatrtTaskController(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//并发
	C := r.FormValue("C")
	//请求总数量
	N := r.FormValue("N")
	//持续时间 z > 0 则 n 无效
	Z := r.FormValue("Z")
	//压测的目标URL
	targetUrl := r.FormValue("targetUrl")
	result := make(map[string]interface{})
	result["C"] = C
	result["N"] = N
	result["Z"] = Z
	result["targetUrl"] = targetUrl
	SendJson(w, result)
}

func TestOneTask(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.FormValue("m")
	N, _ := strconv.Atoi(r.FormValue("N"))
	C, _ := strconv.Atoi(r.FormValue("C"))
	Z := r.FormValue("Z")
	z, _ := time.ParseDuration(Z)
	url := r.FormValue("url")
	fileId := r.FormValue("fileId")
	go StartOneTask(method, url, N, C, z, fileId)
	SendJson(w, "OK")
}

//开启一个任务 fileId 参数
func StartOneTaskController(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	res := ResultMsg{
		Code: 0,
		Msg:  "操作成功",
	}
	var err error
	C, err := strconv.Atoi(r.FormValue("C"))
	N, err := strconv.Atoi(r.FormValue("N"))
	log.Printf("C is %s \n", r.FormValue("C"))
	log.Printf("N is %s \n", r.FormValue("N"))
	Z := r.FormValue("Z")
	if err != nil {
		log.Println(err)
		res.Code = 1
		res.Msg = "并发数，请求数为整数"
		SendJson(w, res)
		return
	}

	if Z == "" && C > N {
		res.Code = 1
		res.Msg = "并发数 < 请求总数"
		SendJson(w, res)
		return
	}
	//JSON格式
	QH := r.FormValue("QH")
	qhMap := make(map[string]string)
	if QH != "" {
		err = json.Unmarshal([]byte(QH), &qhMap)
		if err != nil {
			res.Code = 1
			res.Msg = "请求头信息格式错误"
			log.Println(err.Error())
			SendJson(w, res)
			return
		}
	}

	payloadBs, _ := ioutil.ReadAll(r.Body)

	//组装测试参数
	reqParam := TestParam{
		Method:        strings.ToUpper(r.FormValue("Method")),
		C:             C,
		N:             N,
		Remark:        r.FormValue("Remark"),
		Url:           r.FormValue("Url"),
		Type:          r.FormValue("Type"),
		FileId:        r.FormValue("FileId"),
		Z:             Z,
		QH:            qhMap,
		InputFileName: r.FormValue("InputFileName"),
	}

	if len(payloadBs) > 0 {
		reqParam.Payload = string(payloadBs)
	}

	ps := r.FormValue("P")
	log.Println("ps", ps)
	if ps != "" {
		pm := make(map[string]interface{})
		jsonErr := json.Unmarshal([]byte(ps), &pm)
		if jsonErr != nil {
			log.Println("解析参数JSON错误", jsonErr.Error())
		} else {
			reqParam.P = pm
		}
	}
	res.Data = reqParam
	go StartTaskWork(reqParam)
	SendJson(w, res)

}

//测试json参数数据
func TestJsonController(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	bs, _ := ioutil.ReadAll(body)
	param := make(map[string]interface{})
	json.Unmarshal(bs, &param)
	log.Printf("name=%s age=%f \n", param["name"], param["age"])
	for i := 0; i < 100; i++ {
		log.Printf("rand %d \n", rand.Intn(100))
	}
	SendJson(w, param)
}

func GetRandUrl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileid := r.FormValue("urlid")
	url := GetRandUrlFormUrlFile(fileid, "baidu.com")
	SendJson(w, url)
}
func GetIntVal(r *http.Request, key string, errPlan int) int {
	value := r.FormValue(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Println(err.Error())
		return errPlan
	}
	return i
}

func QueryTaskPageController(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	curPage := GetIntVal(r, "curPage", 1)
	pageSize := GetIntVal(r, "pageSize", 10)
	SendJson(w, GetTaskPage(curPage, pageSize))
}

func GetTaskSnapController(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	SendJson(w, GetTaskSnap(r.FormValue("taskId")))
}

func TestPs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bs, _ := ioutil.ReadAll(r.Body)
	log.Printf("bs length %d \n", len(bs))
	log.Printf("name  = %s \n", r.FormValue("name"))
	log.Printf("content-type = %s \n", r.Header.Get("Content-Type"))
	urlValues := r.Form
	SendJson(w, urlValues.Encode())
}

func StartServer() {
	http.HandleFunc("/index", IndexController)
	http.HandleFunc("/UpFile", UpFileController)
	http.HandleFunc("/StartTask", SatrtTaskController)
	http.HandleFunc("/TestJson", TestJsonController)
	http.HandleFunc("/SaveUrls", SaveUrlsController)
	http.HandleFunc("/TestOneTask", TestOneTask)
	http.HandleFunc("/GetRandUrl", GetRandUrl)
	log.Printf("start a Server on 7654 \n")
	http.ListenAndServe(":7655", nil)
}

func StartStaticServer() {
	heyHandlerMap["/index"] = IndexController
	heyHandlerMap["/UpFile"] = UpFileController
	heyHandlerMap["/StartTask"] = SatrtTaskController
	heyHandlerMap["/TestJson"] = TestJsonController
	heyHandlerMap["/SaveUrls"] = SaveUrlsController
	heyHandlerMap["/TestOneTask"] = TestOneTask
	heyHandlerMap["/GetRandUrl"] = GetRandUrl
	heyHandlerMap["/StartMyTask"] = StartOneTaskController
	heyHandlerMap["/QueryTaskPage"] = QueryTaskPageController
	heyHandlerMap["/GetTaskSnap"] = GetTaskSnapController
	heyHandlerMap["/TestPs"] = TestPs
	log.Println("start")
	err := heyServer.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}
