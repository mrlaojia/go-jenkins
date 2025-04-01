package jenkinsdk

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/**
* @Author: jack.walker
* @File: jenkins.go
* @CreateDate: 2025/3/29 15:35
* @ChangeDate：2025/3/29 15:35
* @Version：1.0.0
* @Description:
 */

type JenkinsSdk struct {
	Url      string
	Username string
	Token    string // password create 时报错

	csrfEnable bool
	http       *http.Client
}

func NewJenkinsSdk(url string, username string, token string) *JenkinsSdk {
	return &JenkinsSdk{
		Url:      strings.TrimRight(url, "/"),
		Username: username,
		Token:    token,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (j *JenkinsSdk) EnableCsrf() *JenkinsSdk {
	j.csrfEnable = true
	return j
}

func (j *JenkinsSdk) setBasicAuth(req *http.Request) {
	req.SetBasicAuth(j.Username, j.Token)
}

func (j *JenkinsSdk) sendHttp(req *http.Request) ([]byte, error) {

	// 处理 csrf
	if err := j.makeCsrf(req); err != nil {
		return nil, err
	}
	// 设置认证
	j.setBasicAuth(req)

	// 发送请求
	resp, err := j.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return io.ReadAll(resp.Body)
	} else {
		return nil, fmt.Errorf("fail to send http request, Code: %d", resp.StatusCode)
	}
}
