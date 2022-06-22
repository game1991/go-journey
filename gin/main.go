package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"pkg.deepin.com/golang/lib/crypto/sha"
	libGin "pkg.deepin.com/service/lib/gin"
	"pkg.deepin.com/service/lib/response"
)

var (
	port      = "8888"
	accessID  = "242f8ed1f2fb4233aefb68fe4e86c3b3"
	accessKey = "ce6976bee7014999b4852ac47d7bb128"
	// accessID  = "2a306b19d3874321b123e246897d7b90"
	// accessKey = "a81282c01afc4bd0816482c195a3fa38"
	address = ":"
)

func init() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(log.Ldate | log.Lshortfile | log.LstdFlags | log.Ltime)
	flag.StringVar(&port, "p", "8888", "app端口号")
	flag.StringVar(&accessID, "id", "242f8ed1f2fb4233aefb68fe4e86c3b3", "输入应用id进行构建对应app")
	flag.StringVar(&accessKey, "key", "ce6976bee7014999b4852ac47d7bb128", "应用的key")
}

// APIReq ...
type APIReq struct {
	OpenID    string `json:"OpenID"`    // 用户openid
	TaskID    string `json:"TaskID"`    // 任务id
	UID       int    `json:"UID"`       // 用户uid
	TimeStamp int64  `json:"TimeStamp"` // 请求时间的时间戳
	UUID      string `json:"UUID"`      // 请求id
}

func main() {
	flag.Parse()
	appName := os.Getenv("APP_NAME")

	fmt.Println("获取环境变量", appName)

	fmt.Println("--------", accessID, accessKey, port, "-------")

	if accessID == "" || accessKey == "" || port == "" {
		panic(fmt.Sprintln("命令行参数不能为空", accessID, accessKey, port))
	}

	engine := gin.Default()

	engine.Use(
		func() gin.HandlerFunc {
			store, err := redis.NewStore(50, "tcp", "127.0.0.1:6379", "", []byte("我的密钥"))
			if err != nil {
				panic(err)
			}
			store.Options(sessions.Options{
				Path:   "/",
				Domain: "",
				MaxAge: 1200,
			})

			return sessions.Sessions("mysession", &libGin.Store{
				Store:       store,
				MaxAgeZero:  true,
				ConstExp:    false,
				AutoRefresh: false,
			})
		}())

	route := engine.Group("")

	route.Any("", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 500)
		c.Set("sum", 10)

		type request struct {
			ClientID string `form:"client_id"`
			Name     string `form:"name"`
			Log      *Log
		}

		req := &request{}

		if err := c.ShouldBindQuery(req); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  err.Error(),
			})
			return
		}

		ageStr := c.Request.FormValue("age")
		age, _ := strconv.Atoi(ageStr)
		ti := c.Request.FormValue("time")
		warnStr := c.Request.FormValue("warn")

		var warn bool
		if warnStr == "true" {
			warn = true
		}

		req.Log = &Log{
			Time: ti,
			Age:  age,
			Line: "",
			Warn: warn,
		}
		bts, err := json.Marshal(req.Log)
		if err != nil {
			fmt.Println(err)
			return
		}

		if req.Log != nil {
			c.Set("req", req.Log)

			param := url.Values{}
			param.Add("time", time.Now().Format(time.RFC3339))
			param.Add("age", "18")
			param.Add("warn", "true")
			param.Add("log", string(bts))

			baseURL, _ := url.Parse("http://www.baidu.com")
			baseURL.RawQuery = param.Encode()
			urlPath := baseURL.String()
			fmt.Println(urlPath)

		}

		/*
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"sum":  10,
			})*/

		resp, err := http.Get("https://www.baidu.com")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bts, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatal("访问失败")
		}

		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)
		c.Writer.WriteString(string(bts))
		c.Abort()
		return
	})

	route.POST("timeout", func(c *gin.Context) {
		now := time.Now()
		for i := 0; ; i++ {
			select {
			case <-time.After(time.Second * 5):
				fmt.Printf("我没有完成，当前循环第%v次\n", i+1)
				continue
			case <-c.Done():
				fmt.Printf("耗时：%v秒", now.Sub(time.Now())/time.Second)
				break
			}
		}
	})

	route.POST("callback", Auth, func(c *gin.Context) {
		// 取数据，进行处理哦
		req := &APIReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// 校验类型
		tp := c.GetHeader("X-BusinessType")
		sucessCode := http.StatusOK

		msg := fmt.Sprintf("任务ID[%s]已完成用户openid:%s的类型%s请求", req.TaskID, req.OpenID, tp)

		// if isPrime3(rand.Intn(99)) {
		// 	sucessCode = http.StatusCreated
		// 	msg = fmt.Sprintf("任务ID[%s]拒绝用户openid:%s的类型%s请求", req.TaskID, req.OpenID, tp)
		// }

		// if rand.Intn(99)%2 == 0 {
		// 	sucessCode = http.StatusInternalServerError
		// 	msg = fmt.Sprintf("任务ID[%s]失败用户openid:%s的类型%s请求", req.TaskID, req.OpenID, tp)
		// }

		// if tp == "3" {
		// 	sucessCode = http.StatusInternalServerError
		// }

		c.JSON(http.StatusOK, response.Response{
			Result:  true,
			Code:    sucessCode,
			Message: msg,
		})
	})

	engine.Run(address + port)
}

// Log ...
type Log struct {
	Time string `form:"time"`
	Age  int    `form:"age"`
	Line string `form:"line"`
	Warn bool   `form:"warn"`
}

func doSomething(ctx *gin.Context) {
	var incr int
	if v, ok := ctx.Get("sum"); ok && v != nil {
		log.Println(v.(int))
		incr = v.(int)
	}

	ch := time.After(time.Millisecond * 1500)

LOOP:
	for {
		incr++
		select {
		case t := <-ch:
			fmt.Println("任务成功", t.Format(time.RFC3339))
			break LOOP
		default:
			fmt.Println(time.Now().Format("15:04:05.000"), incr)
			time.Sleep(time.Millisecond * 500)
		}
	}

	if v, ok := ctx.Get("req"); ok && v != nil {
		fmt.Println("req", v.(*Log))
	}

}

// Auth 接口校验
func Auth(c *gin.Context) {
	// 读取http请求的body数据，然后缓存起来
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println(string(data))

	sign := accessID + ":" + sha.SHA256BASE64([]byte(string(data)+"||"+accessKey))
	// 校验签名
	if sign != c.GetHeader("Authorization") {
		c.AbortWithStatusJSON(http.StatusOK, response.Response{
			Result:  false,
			Code:    http.StatusUnauthorized,
			Message: fmt.Sprintf("签名校验失败，应该是%s,传入的是%s", sign, c.GetHeader("Authorization")),
		})
		return
	}
	// 使用一个新的io.Reader存储
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

}

/*
首先看一个关于质数分布的规律：大于等于5的质数一定和6的倍数相邻。例如5和7，11和13,17和19等等；
证明：令x≥1，将大于等于5的自然数表示如下：
······ 6x-1，6x，6x+1，6x+2，6x+3，6x+4，6x+5，6(x+1），6(x+1)+1 ······
可以看到，不在6的倍数两侧，即6x两侧的数为6x+2，6x+3，6x+4，由于2(3x+1)，3(2x+1)，2(3x+2)，
所以它们一定不是素数，再除去6x本身，显然，素数要出现只可能出现在6x的相邻两侧。
这里要注意的一点是，在6的倍数相邻两侧并不是一定就是质数。
根据以上规律，判断质数可以6个为单元快进，即将方法循环中i++步长加大为6，加快判断速度
孪生素数自行研究
*/
func isPrime3(num int) bool {
	//两个较小数另外处理
	if 2 == num || 3 == num {
		return true
	}

	//不在6的倍数两侧的一定不是质数
	if num%6 != 1 && num%6 != 5 {
		return false
	}
	mid := int(math.Sqrt(float64(num)))
	// 在6的倍数两侧的也可能不是质数
	for i := 5; i <= mid; i += 6 {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}

	//排除所有，剩余的是质数
	return true
}
