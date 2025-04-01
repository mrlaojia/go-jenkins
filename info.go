package jenkinsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
* @Author: jack.walker
* @File: info.go
* @CreateDate: 2025/3/30 16:36
* @ChangeDate：2025/3/30 16:36
* @Version：1.0.0
* @Description:
 */

// JenkinsInfo 基础信息响应结构
type JenkinsInfo struct {
	//JenkinsVersion string `json:"jenkinsVersion"`
	//QuietingDown   bool   `json:"quietingDown"`
	UseCrumbs   bool `json:"useCrumbs"`   // 是否启用 CSRF
	UseSecurity bool `json:"useSecurity"` // 是否启用安全
}

func (j *JenkinsSdk) GetJenkinsInfo() (*JenkinsInfo, error) {
	// 构造 API 地址
	apiURL := fmt.Sprintf("%s/api/json", j.Url)
	// 创建请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	data, err := j.sendHttp(req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	info := &JenkinsInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	return info, nil
}
