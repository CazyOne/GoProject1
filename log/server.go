package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

//实现io.Writer接口的结构体。
type fileLog string

var _ io.Writer = (*fileLog)(nil)

/*
实现了io.Writer接口的Write方法，用于将日志数据写入指定的文件中。
每次写入都要打开文件，写入数据，然后关闭文件。
使用os.O_APPEND确保日志不会被覆盖
文件权限0600（仅所有者可读写）
*/
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

/*
初始化函数
创建标准log.Logger实例
前缀为"[go] - "
stlog.LstdFlags启用标准日志标志（日期和时间）
*/
func Run(destination string) {
	log = stlog.New(fileLog(destination), "[go] - ", stlog.LstdFlags)
}

/*
HTTP处理函数
注册/log端点
只接受POST请求，读取请求体作为日志内容
验证请求体不为空
*/
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
