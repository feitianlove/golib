package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/go-redis/redis"
	"github.com/jinzhu/now"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

func RandomNum(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func GetLocalIP() string {
	conn, _ := net.Dial("udp", "10.1.1.1:80")
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

var InnerIPMark = []string{"10.0.0.0/8", "9.0.0.0/8", "11.0.0.0/8", "30.0.0.0/8", "100.64.0.0/16", "172.16.0.0/16", "192.168.0.0/16"} // 公司内网网段

func IsInnerIP(ip string) bool {
	ipA, _, err := net.ParseCIDR(fmt.Sprintf("%s/24", ip))
	if err != nil {
		return false
	}
	for _, mark := range InnerIPMark {
		_, inetIn, err := net.ParseCIDR(mark)
		if err != nil {
			return false
		}
		if inetIn.Contains(ipA) {
			return true
		}
	}
	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsInsideDocker(pid int64) bool {
	cgroup, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		return false
	}
	if strings.Contains(string(cgroup), "docker") {
		return true
	}
	return false
}

func InitLogger(path string, reserveDay int) error {
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(os.Args[0]), path)
	}
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(reserveDay)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		return err
	}
	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
	}, &log.TextFormatter{}))

	log.SetOutput(ioutil.Discard)
	log.SetLevel(log.DebugLevel)
	//log.SetReportCaller(true)
	return nil
}

func InitLoggerWithLevel(path string, reserveDay int, level log.Level) error {
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(os.Args[0]), path)
	}
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(reserveDay)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		return err
	}
	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
		log.TraceLevel: writer,
	}, &log.TextFormatter{}))

	log.SetOutput(ioutil.Discard)
	log.SetLevel(level)
	//log.SetReportCaller(true)

	return nil
}

func InitFileNoLoggerWithLevel(path string, reserveDay int, level log.Level) error {
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(os.Args[0]), path)
	}
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(reserveDay)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		return err
	}
	log.SetReportCaller(true)
	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
	}, &log.TextFormatter{CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		repopath := fmt.Sprintf("%sgit.code.oa.com/storage-ops", os.Getenv("GOPATH"))
		filename := strings.Replace(f.File, repopath, "", -1)
		return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
	}}))
	log.SetOutput(ioutil.Discard)
	log.SetLevel(level)
	//log.SetReportCaller(true)
	return nil
}

func ReadFileAsLines(path string) ([]string, error) {
	var lines []string
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return lines, err
	}
	lines = strings.Split(string(data), "\n")
	return lines, nil
}

func FileExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if stat.Mode().IsRegular() {
		return true
	}
	return false
}

func DirExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if stat.IsDir() {
		return true
	}
	return false
}

func Chunk(whole []time.Time, chunkSize int) [][]time.Time {
	var divided [][]time.Time
	for i := 0; i < len(whole); i += chunkSize {
		end := i + chunkSize

		if end > len(whole) {
			end = len(whole)
		}

		divided = append(divided, whole[i:end])
	}
	return divided
}

func RemoveStrInSlices(list []string, str string) []string {
	var idx int
	for i := range list {
		if list[i] == str {
			idx = i
			break
		}
	}
	return append(list[:idx], list[idx+1:]...)
}

func IfValidIP(ip string) bool {
	ipv4Regex := regexp.MustCompile(`^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if ipv4Regex.MatchString(ip) {
		return true
	}
	return false
}

func NewClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConns:        5,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 60 * time.Second,
	}
	return client
}

func HttpGet(url string, timeout int) (int, []byte, error) {
	resp, err := resty.SetTimeout(time.Duration(timeout) * time.Second).R().Get(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

func HttpPost(url string, params map[string]string, timeout int) (int, []byte, error) {
	resp, err := resty.SetTimeout(time.Duration(timeout) * time.Second).R().SetFormData(params).Post(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

func HttpPostAsJson(url string, params interface{}, timeout int) (int, []byte, error) {
	resp, err := resty.SetTimeout(time.Duration(timeout) * time.Second).R().SetBody(params).Post(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

// 与HttpGet的区别是可以支持传http client 达到复用tcp连接(长连接)的目的
func HttpGetWithClient(client *http.Client, url string, timeout int) (int, []byte, error) {
	resp, err := resty.NewWithClient(client).SetTimeout(time.Duration(timeout) * time.Second).R().Post(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

func HttpPostWithClient(client *http.Client, url string, params map[string]string, timeout int) (int, []byte, error) {
	resp, err := resty.NewWithClient(client).SetTimeout(time.Duration(timeout) * time.Second).R().SetFormData(params).Post(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

func HttpPostAsJsonWithClient(client *http.Client, url string, params interface{}, timeout int) (int, []byte, error) {
	resp, err := resty.NewWithClient(client).SetTimeout(time.Duration(timeout) * time.Second).R().SetBody(params).Post(url)
	if err != nil {
		return -1, nil, err
	}
	return resp.StatusCode(), resp.Body(), nil
}

func NewRedisClient(addr string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}

func PrettyJson(data interface{}) {
	content, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(content))
}

func BeginOfMonth(t time.Time) time.Time {
	return now.New(t).BeginningOfMonth()
}

func EndOfMonth(t time.Time) time.Time {
	return now.New(t).EndOfMonth()
}

func BeginOfDay(t time.Time) time.Time {
	return now.New(t).BeginningOfDay()
}

func EndOfDay(t time.Time) time.Time {
	return now.New(t).EndOfDay()
}

func NewDefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func IsNum(s string) bool {
	pattern := "\\d+" //反斜杠要转义
	result, _ := regexp.MatchString(pattern, s)
	return result
}

//这个方法移动到golib/common/file
//func Md5sum(filePath string) (string, error) {

// 判断obj是否在target中，target支持的类型arrary,slice,map
// 反射效率较低 慎用
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// 是否包含string
func ContainStr(obj string, target []string) bool {
	for _, s := range target {
		if s == obj {
			return true
		}
	}
	return false
}

func Sha1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

func Sha1Bytes(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum([]byte("")))
}

//乱序string数组
func ShuffleString(tasks []string) []string {
	newTasks := make([]string, len(tasks))
	idS := GenerateRandomNumber(0, len(tasks), len(tasks))
	index := 0
	for _, id := range idS {
		newTasks[index] = tasks[id]
		index += 1
	}
	return newTasks
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn(end-start) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

//sleep
func SleepSecond(n int) {
	time.Sleep(time.Second * time.Duration(n))
}

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.WithField("funcName", name).WithField("elapsed", elapsed).Info("TimeTrack")
}

func TimeTrackUseLogger(start time.Time, logger *log.Logger) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	logger.WithField("funcName", name).WithField("elapsed", elapsed).Info("TimeTrack")
}

// GBK To UTF-8
func ConvertToByte(src string, srcCode string, targetCode string) ([]byte, error) {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, err := tagCoder.Translate([]byte(srcResult), true)
	return cdata, err
}

// 判断时间是否没有初始化 默认 1990-01-01 01:01:01 为最早时间
func IsZeroTime(t time.Time) bool {
	return t.Unix() < 631126861
}

// 把一个list按照大小拆分段
func SplitSlice(s []string, number int) [][]string {
	if len(s) <= number {
		return [][]string{s}
	}
	var res [][]string
	for i := 0; i < len(s); i += number {
		res = append(res, s[i:int(math.Min(float64(i+number), float64(len(s))))])
	}
	return res
}
