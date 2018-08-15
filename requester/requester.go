package requester

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const maxResult = 1000000
const maxIdleConns = 500

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
	var size int64
	var code int
	var dnsStart, connStart, resStart, reqStart, delayStart time.Duration
	var dnsDuration, connDuration, resDuration, reqDuration, delayDuration time.Duration
	//req := cloneRequest(b.Request, b.RequestBody)
	req := createRandSimpleRequest(*b.Header, b.Method, b.SrcUrl, b.UrlFileId)
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

func (b *Work) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(b.C)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         b.Request.Host,
		},
		MaxIdleConnsPerHost: min(b.C, maxIdleConns),
		DisableCompression:  b.DisableCommpression,
		DisableKeepAlives:   b.DisableKeepAlives,
		Proxy:               http.ProxyURL(b.ProxyAddr),
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

	for i := 0; i < b.C; i++ {
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

	req, err := http.NewRequest(testParam.Method, testParam.Url, nil)
	if err != nil {
		testParam.Err = err.Error()
		testParam.Status = 1
		return
	}
	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
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
