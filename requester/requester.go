package requester

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const maxResult = 1000000
const maxIdleConns = 500

var UGS = []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 OPR/26.0.1656.60", "Opera/8.0 (Windows NT 5.1; U; en)", "Mozilla/5.0 (Windows NT 5.1; U; en; rv:1.8.1) Gecko/20061208 Firefox/2.0.0 Opera 9.50", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0", "Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.71 Safari/537.1 LBBROWSER", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; LBBROWSER) ", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 SE 2.X MetaSr 1.0", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.3.4000 Chrome/30.0.1599.101 Safari/537.36", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 UBrowser/4.0.3214.0 Safari/537.36", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5", "Mozilla/5.0 (iPod; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5", "Mozilla/5.0 (iPad; U; CPU OS 4_2_1 like Mac OS X; zh-cn) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8C148 Safari/6533.18.5", "Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5", "Mozilla/5.0 (Linux; U; Android 2.2.1; zh-cn; HTC_Wildfire_A3333 Build/FRG83D) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1", "Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1", "MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1", "Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10", "Mozilla/5.0 (Linux; U; Android 3.0; en-us; Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13", "Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; en) AppleWebKit/534.1+ (KHTML, like Gecko) Version/6.0.0.337 Mobile Safari/534.1+", "Mozilla/5.0 (hp-tablet; Linux; hpwOS/3.0.0; U; en-US) AppleWebKit/534.6 (KHTML, like Gecko) wOSBrowser/233.70 Safari/534.6 TouchPad/1.0", "Mozilla/5.0 (SymbianOS/9.4; Series60/5.0 NokiaN97-1/20.0.019; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) BrowserNG/7.1.18124", "Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; HTC; Titan)", "UCWEB7.0.2.37/28/999", "NOKIA5700/ UCWEB7.0.2.37/28/999", "Openwave/ UCWEB7.0.2.37/28/999", "Mozilla/4.0 (compatible; MSIE 6.0; ) Opera/UCWEB7.0.2.37/28/999"}

const (
	USER_AGENT          = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
	CONTENT_TYPE_NORMAL = "application/x-www-form-urlencoded"
)

func getRandUg() string {
	var size int = len(UGS)
	rand.Seed(time.Now().UnixNano())
	var randIndex int = rand.Intn(size)
	return UGS[randIndex]

}

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	connDuration  time.Duration
	dnsDuration   time.Duration
	reqDuration   time.Duration
	resDuration   time.Duration
	delayDuration time.Duration
	contentLength int64
}

type Work struct {
	Request     *http.Request
	Header      *http.Header
	Method      string
	SrcUrl      string
	RequestBody []byte
	N           int

	C int

	H2 bool

	Timeout int

	QPS float64

	DisableCommpression bool

	DisableKeepAlives bool

	DisableRedirects bool

	Output string

	ProxyAddr *url.URL

	Writer io.Writer

	initOnce sync.Once

	results chan *result

	stopCh chan struct{}
	start  time.Duration

	report *report

	TaskId int

	UrlFileId string

	TestParam *TestParam
}

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func (b *Work) writer() io.Writer {
	if b.Writer == nil {
		return os.Stdout
	}
	return b.Writer
}

func (b *Work) Init() {
	b.initOnce.Do(func() {
		b.results = make(chan *result, min(b.C*1000, maxResult))
		b.stopCh = make(chan struct{}, b.C)
	})
}

func (b *Work) Run() {
	b.Init()
	b.start = now()
	b.report = newReport(b.writer(), b.results, b.Output, b.N)
	b.report.taskId = b.TaskId
	go func() {
		runReporter(b.report)
	}()
	b.runWorkers()
	b.Finish()
}

func (b *Work) Stop() {
	for i := 0; i < b.C; i++ {
		b.stopCh <- struct{}{}
	}
}

func (b *Work) Finish() {
	close(b.results)
	total := now() - b.start
	//等待统计report
	<-b.report.done
	b.report.finalize(total)
}

func (b *Work) makeRequest(c *http.Client) {
	s := now()
	testParam := b.TestParam
	var size int64
	var code int
	var req *http.Request
	var dnsStart, connStart, resStart, reqStart, delayStart time.Duration
	var dnsDuration, connDuration, resDuration, reqDuration, delayDuration time.Duration
	//req := cloneRequest(b.Request, b.RequestBody)
	if testParam.Type == "FIXED" {
		req = b.makeFormRequest()
	} else if testParam.Type == "FILE" {
		req = b.makeFileRequest()
	} else {
		req = b.makeRandRequest()
		//req = createRandSimpleRequest(*b.Header, testParam.Method, testParam.Url, testParam.FileId)
		//log.Printf("new Method=%s url=%s fileid=%s \n", testParam.Method, testParam.Url, testParam.FileId)
		//log.Printf("old Method=%s url=%s fileid=%s \n", b.Method, b.SrcUrl, b.UrlFileId)
		//req = createRandSimpleRequest(*b.Header, b.Method, b.SrcUrl, b.UrlFileId)
	}
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = now()
		},
		DNSDone: func(dnsinfo httptrace.DNSDoneInfo) {
			dnsDuration = now() - dnsStart
		},
		GetConn: func(h string) {
			connStart = now()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if !connInfo.Reused {
				connDuration = now() - connStart
			}
			reqStart = now()
		},
		WroteRequest: func(w httptrace.WroteRequestInfo) {
			reqDuration = now() - reqStart
			delayStart = now()
		},
		GotFirstResponseByte: func() {
			delayDuration = now() - delayStart
			resStart = now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := c.Do(req)
	if err == nil {
		size = resp.ContentLength
		code = resp.StatusCode
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	t := now()
	resDuration = t - resStart
	finish := t - s
	b.results <- &result{
		statusCode:    code,
		duration:      finish,
		err:           err,
		contentLength: size,
		connDuration:  connDuration,
		dnsDuration:   dnsDuration,
		reqDuration:   reqDuration,
		resDuration:   resDuration,
		delayDuration: delayDuration,
	}
}

func (b *Work) runWorker(client *http.Client, n int) {
	var throttle <-chan time.Time

	if b.QPS > 0 {
		log.Printf("test QPS ---- %f \n", b.QPS)
		throttle = time.Tick(time.Duration(1e6/(b.QPS)) * time.Microsecond)
	}

	if b.DisableRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	for i := 0; i < n; i++ {
		select {
		case <-b.stopCh:
			return
		default:
			if b.QPS > 0 {
				<-throttle
			}
			b.makeRequest(client)
		}
	}
}
func getProxy() *url.URL {
	resp, _ := http.Get("http://192.168.1.192:5010/get")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var proxyUrl string = string(body)
	fmt.Printf("%s \n", proxyUrl)
	url, _ := url.Parse("http://" + proxyUrl)
	return url
}
func (b *Work) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(b.C)

	for i := 0; i < b.C; i++ {
		fmt.Printf("this is a %d req\n", i)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         b.Request.Host,
			},
			MaxIdleConnsPerHost: min(b.C, maxIdleConns),
			DisableCompression:  b.DisableCommpression,
			DisableKeepAlives:   b.DisableKeepAlives,
			Proxy:               http.ProxyURL(getProxy()),
		}

		if b.H2 {
			http2.ConfigureTransport(tr)
		} else {
			tr.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
		}

		client := &http.Client{
			Transport: tr,
			Timeout:   time.Duration(b.Timeout) * time.Second,
		}

		go func() {
			b.runWorker(client, b.N/b.C)
			wg.Done()
		}()
	}
	wg.Wait()
}

func cloneRequest(r *http.Request, body []byte) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	if len(body) > 0 {
		r2.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	r2.Close = true
	return r2
}

//上传文件
func createMultiPartRequest(h http.Header) *http.Request {

	return nil
}

//随机请求
func createRandSimpleRequest(h http.Header, method string, url string, urlfileId string) *http.Request {
	targetUrl := GetRandUrlFormUrlFile(urlfileId, url)
	r2, _ := http.NewRequest(method, targetUrl, nil)
	r2.Header = make(http.Header, len(h))
	for k, s := range h {
		r2.Header[k] = append([]string(nil), s...)
	}

	return r2
}

//随机请求主要是从文件列表中获取随机访问的URL地址进行请求
func (b *Work) makeRandRequest() *http.Request {
	testParam := b.TestParam
	defaultUrl := testParam.Url
	header := b.Header
	targetUrl := GetRandUrlFormUrlFile(testParam.FileId, defaultUrl)
	r2, _ := http.NewRequest(testParam.Method, targetUrl, nil)
	r2.Header = make(http.Header)
	for k, v := range *header {
		r2.Header[k] = append([]string(nil), v...)
	}
	return r2
}

//处理-FIXED 请求
func (b *Work) makeFormRequest() *http.Request {
	testParam := b.TestParam
	defaultUrl := testParam.Url
	header := b.Header
	var reader *strings.Reader
	//fmt.Println("--------------", header.Get("Content-Type"))
	if header.Get("Content-Type") == CONTENT_TYPE_NORMAL {
		params := url.Values{}
		for key, value := range testParam.P {
			//		log.Println(key, value)
			params.Add(key, fmt.Sprintf("%s", value))
		}
		reader = strings.NewReader(params.Encode())
	} else {
		reader = strings.NewReader(testParam.Payload)
	}
	r2, _ := http.NewRequest(testParam.Method, defaultUrl, reader)
	r2.Header = make(http.Header)
	for k, v := range *header {
		r2.Header[k] = append([]string(nil), v...)
	}
	r2.Header.Set("User-Agent", getRandUg())
	return r2
}

//上传文件请求
func (b *Work) makeFileRequest() *http.Request {
	testParam := b.TestParam
	url := testParam.Url
	fid := testParam.FileId
	fileInfo := GetFileInfo(fid)
	fh, err := os.Open(fdir + fid + fileInfo.Ext)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer func() {
		fh.Close()
	}()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, ferr := writer.CreateFormFile(testParam.InputFileName, fileInfo.Name)
	if ferr != nil {
		log.Println(ferr.Error())
		return nil
	}
	io.Copy(formFile, fh)
	contentType := writer.FormDataContentType()
	writer.Close()
	request_reader := io.MultiReader(buf)
	req, _ := http.NewRequest("POST", url, nil)
	rc, ok := request_reader.(io.ReadCloser)
	if !ok && request_reader != nil {
		rc = ioutil.NopCloser(request_reader)
	}
	req.Body = rc

	header := b.Header
	req.Header = make(http.Header)
	for k, v := range *header {
		req.Header[k] = append([]string(nil), v...)
	}
	req.Header.Set("Content-Type", contentType)
	return req
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//开启一个任务
func StartOneTask(method string, url string, N int, C int, dur time.Duration, urlFileId string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.New("请求路劲错误")
	}
	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	taskId := GetTaskId()

	if dur > 0 {
		N = math.MaxInt32
	}
	w := &Work{
		Request:   req,
		N:         N,
		C:         C,
		QPS:       -1,
		TaskId:    taskId,
		UrlFileId: urlFileId,
		Header:    &header,
	}
	w.Init()
	if dur > 0 {
		go func() {
			time.Sleep(dur)
			w.Stop()
		}()
	}
	w.Run()

	return nil
}

func StartTaskWork(testParam TestParam) {
	taskId := GetTaskId()
	testParam.TaskId = taskId
	testParam.StartTime = time.Now().UnixNano()
	req, err := http.NewRequest(testParam.Method, testParam.Url, nil)
	if err != nil {
		testParam.Err = err.Error()
		testParam.Status = 1
		SaveTestParam(&testParam)
		return
	}
	SaveTestParam(&testParam)

	header := make(http.Header)

	//设置默认的请求头
	header.Set("Content-Type", CONTENT_TYPE_NORMAL)
	header.Set("User-Agent", getRandUg())
	for key, value := range testParam.QH {
		header.Set(key, value)
	}
	N := testParam.N

	dur, err := time.ParseDuration(testParam.Z)
	if err == nil && dur > 0 {
		log.Println("*******重置请求数N*************")
		N = math.MaxInt32
	}

	w := &Work{
		Request:   req,
		N:         N,
		C:         testParam.C,
		QPS:       -1,
		TaskId:    taskId,
		UrlFileId: testParam.FileId,
		Header:    &header,
		TestParam: &testParam,
		Method:    testParam.Method,
		SrcUrl:    testParam.Url,
	}
	w.Init()
	if dur > 0 {
		go func() {
			time.Sleep(dur)
			w.Stop()
		}()
	}
	w.Run()
}
