package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var idWorker *IdWorker = &IdWorker{}

type LonhData struct {
	Abcd   []string `json:"Abcd"`
	Gpsid  []string `json:"Gpsid"`
	Abcdyn []string `json:Abcdyn`
}

//龙慧数据缓存
var LonhDataParam LonhData

//初始化数据
func (lonhData *LonhData) initData() {
	dataFilePath := os.Getenv("LONH_FILE")
	strBytes, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		log.Printf("读取数据错误 %s %s \n", dataFilePath, err)
		return
	}
	LonhDataParam = LonhData{}
	json.Unmarshal(strBytes, &LonhDataParam)
}

//获取随机gps设备id
func (lonhData *LonhData) GetRndGpsId() string {
	if len(LonhDataParam.Abcd) == 0 {
		fmt.Println("初始化lonhData")
		lonhData.initData()
	}
	rand.Seed(time.Now().Unix())
	size := len(LonhDataParam.Gpsid)
	rnd := rand.Intn(size)
	return LonhDataParam.Gpsid[rnd]
}

//获取随机区域id
func (lonhData *LonhData) GetRndAbcd() string {
	if len(LonhDataParam.Abcd) == 0 {
		fmt.Println("初始化lonhData")
		lonhData.initData()
	}
	rand.Seed(time.Now().Unix())
	size := len(LonhDataParam.Abcd)
	rnd := rand.Intn(size)
	return LonhDataParam.Abcd[rnd]
}

//获取随机云南省区域id
func (lonhData *LonhData) GetRndAbcdyn() string {
	if len(LonhDataParam.Abcdyn) == 0 {
		fmt.Println("初始化lonhData")
		lonhData.initData()
	}
	rand.Seed(time.Now().Unix())
	size := len(LonhDataParam.Abcdyn)
	rnd := rand.Intn(size)
	return LonhDataParam.Abcdyn[rnd]
}

//获取点线数据 如果是线 则为多个点
func (lonhData *LonhData) GetGpsPoints(size int) string {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		buffer.WriteString("102.579081448891,25.0737392461037,0")
		if i < size-1 {
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
}

//随机获取上线签到接口数据
func (lonhData *LonhData) GetSignOnlineParam() map[string]string {
	param := make(map[string]string)
	param["projectid"] = "slfh"
	param["gpsid"] = lonhData.GetRndGpsId()
	if time.Now().Unix()%2 == 0 {
		param["type"] = "手机"
	} else {
		param["type"] = "电脑"
	}
	return param
}

//获取新增修改标会元素信息参数
func (lonhData *LonhData) GetSaveFiregroundpmParam() map[string]string {
	param := make(map[string]string)
	param["creategpsid"] = "171019"
	param["createperson"] = "汉高祖-刘邦"
	param["createunit"] = "沛丰邑中阳里防火办"
	param["createunitid"] = "133487"
	param["delflag"] = "0"
	param["description"] = "刘邦的测试数据"
	param["groundid"] = "ProjectTest_530300000000000_20190225132654349"
	param["iconhref"] = "com/lonhwin/source/staticIcon/nperson.png"
	param["iconscal"] = "1"
	param["labelcolor"] = "FFFFFFFF"
	param["labelscal"] = "1"
	param["linewidth"] = "0"
	param["linkid"] = "171009"
	param["linktype"] = "0"
	param["name"] = "汉高祖-刘邦"
	param["pmid"] = strings.Join([]string{"Test", strconv.FormatInt(idWorker.NextId(), 10)}, "_")
	param["projectid"] = "slfh"
	param["stateid"] = strings.Join([]string{"TestStates", lonhData.GetRndGpsId()}, "_")

	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(10000)
	rnd1 := rand.Intn(800) + 1
	if rnd%2 == 0 {
		//点
		param["shapetype"] = "1"
		param["coordinates"] = lonhData.GetGpsPoints(1)
	} else {
		//面
		param["shapetype"] = "2"
		param["coordinates"] = lonhData.GetGpsPoints(rnd1)
	}
	return param
}

//构造查询标会元素参数
func (lonhData *LonhData) GetFindFiregroundpmParams() map[string]string {
	param := make(map[string]string)
	param["groundid"] = "ProjectTest_530300000000000_20190225132654349"
	param["stateid"] = strings.Join([]string{"TestStates", lonhData.GetRndGpsId()}, "_")
	return param
}

//构造获取在线人员接口查询参数
func (lonhData *LonhData) GetFindGpsOnlinelistParams() map[string]string {
	param := make(map[string]string)
	param["projectid"] = "slfh"
	param["adcd"] = lonhData.GetRndAbcdyn()
	return param
}
