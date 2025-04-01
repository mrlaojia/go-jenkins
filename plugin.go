package jenkinsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
* @Author: jack.walker
* @File: plugin.go
* @CreateDate: 2025/3/30 13:10
* @ChangeDate：2025/3/30 13:10
* @Version：1.0.0
* @Description:
 */

type JenkinsPlugin struct {
	ShortName string `json:"shortName"` // 插件短名称 (如 "cloudbees-folder")
	Version   string `json:"version"`   // 插件版本 (如 "6.991.v1d5f531726d0")
	Url       string `json:"url"`       // 插件Url (如  "https://plugins.jenkins.io/apache-httpcomponents-client-4-api")
}

type JenkinsPluginResponse struct {
	Plugins []*JenkinsPlugin `json:"plugins"` // 插件列表
}

// GetJenkinsPlugin 获取某个插件的名字和版本
func (j *JenkinsSdk) GetJenkinsPlugin(name string) (*JenkinsPlugin, error) {

	allPlugins, err := j.GetAllJenkinsPlugins()
	if err != nil {
		return nil, err
	}

	// 查找目标插件
	targetPlugin := "cloudbees-folder"
	for _, plugin := range allPlugins.Plugins {
		if plugin.ShortName == targetPlugin {
			return plugin, nil
		}
	}

	return nil, fmt.Errorf("插件 %s 未安装", targetPlugin)
}

// GetAllJenkinsPlugins 获取所有插件及版本
func (j *JenkinsSdk) GetAllJenkinsPlugins() (*JenkinsPluginResponse, error) {
	// 构建 API 地址
	fullURL := j.Url + "/pluginManager/api/json?depth=1"

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	data, err := j.sendHttp(req)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 响应
	pluginResponse := &JenkinsPluginResponse{
		Plugins: make([]*JenkinsPlugin, 0),
	}

	if err := json.Unmarshal(data, &pluginResponse); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	return pluginResponse, nil

}
