package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

var (
	h bool
	v bool
	s string
)
var once sync.Once
var usagetest = `Usage: hey [options...] <url>

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

  -host HTTP Host header.

  -disable-compression  Disable compression.
  -disable-keepalive    Disable keep-alive, prevents re-use of TCP
                        connections between different HTTP requests.
  -disable-redirects    Disable following of HTTP redirects
  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
`

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.StringVar(&s, "s", "", "send `signal` to a master")

}

func main1() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	flag.Parse()

	fmt.Printf("h=%t  s=%s nags=%d nflag=%d \n", h, s, flag.NArg(), flag.NFlag())
	if h {
		fmt.Printf("useage-----------")
	}

	url := flag.Args()[0]
	var bodys []byte
	bodys = []byte(url)
	fmt.Printf("body length is %d \n", len(bodys))
	fmt.Printf("args length %d url is %s \n", len(flag.Args()), url)
}

func main2() {

	data_path := flag.String("D", "/home/manu/sample/", "DB data path")
	log_file := flag.String("l", "/home/manu/sample.log", "log file")
	nowait_flag := flag.Bool("W", false, "do not wait until operation completes")

	flag.Parse()

	var cmd string = flag.Arg(0)

	fmt.Printf("action   : %s\n", cmd)
	fmt.Printf("data path: %s\n", *data_path)
	fmt.Printf("log file : %s\n", *log_file)
	fmt.Printf("nowait     : %v\n", *nowait_flag)

	fmt.Printf("-------------------------------------------------------\n")

	fmt.Printf("there are %d non-flag input param\n", flag.NArg())
	for i, param := range flag.Args() {
		fmt.Printf("#%d    :%s\n", i, param)
	}

}
func m1() {
	fmt.Printf("this is m1 \n")
}
func m2() {
	fmt.Printf("this is m2 \n")
}

func maintestt() {
	a1 := []int{1, 2, 4}
	a2 := []int{9, 0, 8}
	a3 := append(a1, a2...)
	fmt.Printf("a3 = %v\n", a3)
	fmt.Printf("a1 = %v\n", a1)
	fmt.Printf("a2 = %v\n", a2)
}

func mainreg() {
	reg := `^([\w-]+):\s*(.+)`
	re := regexp.MustCompile(reg)
	matchs := re.FindStringSubmatch("Content-Type:text/html")
	for _, m := range matchs {
		fmt.Printf(" %s\n ", m)
	}
}

func mainread() {
	filename := "D:\\cm.txt"
	slurp, _ := ioutil.ReadFile(filename)
	fmt.Printf("file bytes size %d \n", len(slurp))
	str := string(slurp[200:300])
	fmt.Printf("%s\n", str)
}

func maintimer() {
	timer := time.NewTimer(2 * time.Second)
	fmt.Printf("now Time is %v \n", time.Now())
	expireTime := <-timer.C
	fmt.Printf("expiretime is %v \n", expireTime)
	fmt.Printf("stop timer %v \n", timer.Stop())
}

func maintick() {
	intChan := make(chan int, 1)
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			select {
			case intChan <- 1:
			case intChan <- 2:
			case intChan <- 3:
			}
		}
		fmt.Printf("End.[sender] \n")
	}()

	var sum int
	for e := range intChan {
		fmt.Printf("receive:%v \n", e)
		sum += e
		if sum > 10 {
			fmt.Printf("Got: %v \n", sum)
			break
		}
	}

	ticker.Stop()
	time.Sleep(4 * time.Second)
	fmt.Println("End receiver")
}

func maine6() {
	fmt.Printf("%f \n", 1e6)
}

func mainiu() {
	f, _ := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE, 0777)
	f.WriteString("DMXMXMM\n")
	defer f.Close()
}

func mainpost() {
	url := "http://192.168.227.129:9090/api/v1/oss"
	fb, _ := ioutil.ReadFile("test.txt")
	respone, _ := http.Post(url, "multipart/form-data", bytes.NewReader(fb))
	defer func() {
		respone.Body.Close()
		fmt.Printf("finish*******\n")
	}()

	body, _ := ioutil.ReadAll(respone.Body)
	fmt.Printf("%s\n", string(body))
}

func main222() {
	url := "http://192.168.227.129:9090/api/v1/oss"
	req, _ := http.NewRequest("POST", url, nil)
	fb, _ := ioutil.ReadFile("test.txt")
	req.ContentLength = int64(len(fb))
	header := make(http.Header)
	header.Set("Content-Type", "multipart/form-data")
	req.Body = ioutil.NopCloser(bytes.NewReader(fb))

	tr := &http.Transport{}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", string(body))

}

func mainmultii() {
	fb, _ := ioutil.ReadFile("test.txt")
	url := "http://192.168.227.129:9090/api/v1/oss"

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, _ := writer.CreateFormFile("file", "dxm.txt")
	io.Copy(formFile, bytes.NewReader(fb))
	contentType := writer.FormDataContentType()
	writer.Close()

	response, _ := http.Post(url, contentType, buf)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", string(body))

}

func main22() {
	fb, _ := ioutil.ReadFile("01.pn")
	filetype := http.DetectContentType(fb)
	fmt.Printf("%s\n", filetype)
}

func postFile(filename string, target_url string) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	fmt.Printf("boundary %s \n", boundary)
	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	return http.DefaultClient.Do(req)
}

func main765() {
	url := "http://192.168.227.129:9090/api/v1/oss"
	fh, _ := os.Open("test.txt")

	//buf := new(bytes.Buffer)
	//buf := bytes.NewBufferString("")
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		fmt.Println("err")
	}
	io.Copy(formFile, fh)
	contentType := writer.FormDataContentType()
	boundary := writer.Boundary()
	fmt.Println(boundary)
	writer.Close()
	//close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	//req.Body = io.MultiReader(buf, fh, close_buf)
	request_reader := io.MultiReader(buf)
	req, _ := http.NewRequest("POST", url, nil)

	rc, ok := request_reader.(io.ReadCloser)
	if !ok && request_reader != nil {
		rc = ioutil.NopCloser(request_reader)
	}

	fmt.Println(buf.Len())

	req.Body = rc

	header := make(http.Header)
	header.Set("Content-Type", contentType)
	req.Header = header
	fmt.Println(contentType)

	tr := &http.Transport{}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	//response, err := client.Post(url, contentType, buf)
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", string(body))

}

func main111() {
	url := "http://192.168.227.129:9090/api/v1/oss"
	response, err := postFile("test.txt", url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", string(body))
}

func testFile() {
	f, _ := os.Open("test.txt")
	bytes, _ := ioutil.ReadAll(f)
	fmt.Printf("content1 is\n %s \n", string(bytes))
	fmt.Printf("content2 is \n %s \n", string(bytes))
}

func main1111() {
	testFile()
}

func mainok() {
	fmt.Println("test FDR")
}
