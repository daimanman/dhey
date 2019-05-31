package main

import (
	"dhey/requester"
	"encoding/json"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"regexp"
	"runtime"
	"unicode/utf8"
)

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
	authRegexp   = `^(.+):([^\s].+)`
	heyUA        = "MaoTaoT/0.0.1"
)

var (
	m           = flag.String("m", "GET", "")
	headers     = flag.String("h", "", "")
	body        = flag.String("d", "", "")
	bodyFile    = flag.String("D", "", "")
	accept      = flag.String("A", "", "")
	contentType = flag.String("T", "text/html", "")
	authHeader  = flag.String("a", "", "")
	hostHeader  = flag.String("host", "", "")

	output = flag.String("o", "", "")

	c = flag.Int("c", 50, "")
	n = flag.Int("n", 200, "")
	q = flag.Float64("q", 0, "")
	t = flag.Int("t", 20, "")
	z = flag.Duration("z", 0, "")

	h2                 = flag.Bool("h2", false, "")
	cpus               = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
	disableCompression = flag.Bool("disable-compression", false, "")
	disableKeepAlives  = flag.Bool("disable-keepalive", false, "")
	disableRedirects   = flag.Bool("disable-redirects", false, "")
	proxyAddr          = flag.String("x", "", "")
)
var usage = `Usage: hey [options...] <url>

Options:
  -n  Number of requests to run. Default is 200.
  -c  Number of requests to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -q  Rate limit, in queries per second (QPS). Default is no rate limit.
  -z  Duration of application to send requests. When duration is reached,
      application stops and exits. If duration is specified, n is ignored.
      Examples: -z 10s -z 3m.
  -o  Output type. If none provided, a summary is printed.
      "csv" is the only supported alternative. Dumps the response
      metrics in comma-separated values format.

  -m  HTTP method, one of GET, POST, PUT, DELETE, HEAD, OPTIONS.
  -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
      For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
  -t  Timeout for each request in seconds. Default is 20, use 0 for infinite.
  -A  HTTP Accept header.
  -d  HTTP request body.
  -D  HTTP request body from file. For example, /home/user/file.txt or ./file.txt.
  -T  Content-type, defaults to "text/html".
  -a  Basic authentication, username:password.
  -x  HTTP Proxy address as host:port.
  -h2 Enable HTTP/2.

  -host	HTTP Host header.

  -disable-compression  Disable compression.
  -disable-keepalive    Disable keep-alive, prevents re-use of TCP
                        connections between different HTTP requests.
  -disable-redirects    Disable following of HTTP redirects
  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
`

type headerSlice []string

func errAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func parseInputWithRegexp(input, regx string) ([]string, error) {
	re := regexp.MustCompile(regx)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not parse the provided input ; input = %v ", input)
	}
	return matches, nil
}

func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

//新增修改标会元素信息
func testGetSaveFiregroundpm() {
	lonhData := &requester.LonhData{}
	jsonBytes, _ := json.Marshal(lonhData.GetSaveFiregroundpmParam())
	urlParam := map[string]string{"data": string(jsonBytes)}
	requester.LhSendPost("http://task.lonhcloud.net/webmvc/v1/fireground/saveFiregroundpm", urlParam)
}

//查询标会元素
func testGetFindFiregroundpm() {
	lonhData := &requester.LonhData{}
	urlParam := lonhData.GetFindFiregroundpmParams()
	requester.LhSendPost("http://task.lonhcloud.net/webmvc/v1/fireground/findFiregroundpm", urlParam)
}

//上线签到接口
func testSignOnline() {
	lonhData := &requester.LonhData{}
	urlParam := lonhData.GetSignOnlineParam()
	requester.LhSendPost("http://task.lonhcloud.net/webmvc/v1/elementquery/signOnline", urlParam)
}

//获取在线人员列表，返回政区/单位本级和下一级的在线人员列表
func testFindGpsOnlinelist() {
	lonhData := &requester.LonhData{}
	urlParam := lonhData.GetFindGpsOnlinelistParams()
	urlParam["all"] = "1"
	requester.LhSendPost("http://task.lonhcloud.net/webmvc/v1/elementquery/findGpsOnlinelist", urlParam)
}

func testLonhApi(method string) {
	lonhData := &requester.LonhData{}
	lonhUrlMap := map[string]string{
		"saveFiregroundpm":  "http://task.lonhcloud.net/webmvc/v1/fireground/saveFiregroundpm",
		"findFiregroundpm":  "http://task.lonhcloud.net/webmvc/v1/fireground/findFiregroundpm",
		"signOnline":        "http://task.lonhcloud.net/webmvc/v1/elementquery/signOnline",
		"findGpsOnlinelist": "http://task.lonhcloud.net/webmvc/v1/elementquery/findGpsOnlinelist",
	}
	jsonBytes, _ := json.Marshal(lonhData.GetSaveFiregroundpmParam())
	saveFiregroundpmParam := map[string]string{"data": string(jsonBytes)}

	lonhParamMap := map[string](map[string]string){
		"saveFiregroundpm":  saveFiregroundpmParam,
		"findFiregroundpm":  lonhData.GetFindFiregroundpmParams(),
		"signOnline":        lonhData.GetSignOnlineParam(),
		"findGpsOnlinelist": lonhData.GetFindGpsOnlinelistParams(),
	}
	targetUrl := lonhUrlMap[method]
	targetParam := lonhParamMap[method]
	requester.LhSendPost(targetUrl, targetParam)

}

func testStr() {
	var name string
	name = "abc我是中国人"
	fmt.Printf("%d\n", len(name))
	fmt.Printf("%d\n", utf8.RuneCountInString(name))
	for i, c := range []rune(name) {
		fmt.Printf("(%d - %c ) \n", i, c)
	}
}
func main() {
	//requester.StartServer()
	requester.StartStaticServer()
	//var idWorker requester.IdWorker
	// idWorker = requester.IdWorker{}
	// fmt.Println(idWorker.NextId())

	// fmt.Println(lonhData.GetRndGpsId())
	// fmt.Println(lonhData.GetRndAbcd())
	// fmt.Println(lonhData.GetGpsPoints(10))
	// fmt.Println(lonhData.GetSignOnlineParam())
	//testFindGpsOnlinelist()
	//testGetFindFiregroundpm()

	//testLonhApi("saveFiregroundpm")
	//testLonhApi("findFiregroundpm")
	//testLonhApi("signOnline")
	//testLonhApi("findGpsOnlinelist")

	//testStr()

}
