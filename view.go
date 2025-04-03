package jenkinsdk

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/**
* @Author: jack.walker
* @File: view.go
* @CreateDate: 2025/3/31 09:16
* @Version: 1.0.0
* @Description:
 */

type JenkinsView struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      []string
}

// GetParentPath 构建 ParentPath 路径
// []string{"a", "b"} → /job/a/job/b
func (v *JenkinsView) GetParentPath() string {
	var builder strings.Builder
	for _, elem := range v.Parent {
		builder.WriteString("/job/")
		builder.WriteString(url.PathEscape(elem))
	}
	return builder.String()
}

func (f *JenkinsView) GetFullPath() string {
	return f.GetParentPath() + "/view/" + f.Name
}

type JenkinsListViewTemplate struct {
	XMLName     xml.Name `xml:"hudson.model.ListView"`
	Name        string   `xml:"name,omitempty"`                           // 描述内容
	Description string   `xml:"description,omitempty" json:"description"` // 描述内容
}

// makeFolderConfig 创建 config.xml
// config.xml模板参考: http://172.19.89.76:48080/job/folder/config.xml
func (j *JenkinsSdk) makeListViewConfig(v *JenkinsView) (string, error) {

	// 解析 XML 到结构体
	var ViewTemplate JenkinsListViewTemplate
	// 修改 Plugin 和 Description
	ViewTemplate.Description = v.Description
	ViewTemplate.Name = v.Name

	// 生成 XML
	newXML, err := xml.MarshalIndent(ViewTemplate, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(newXML), nil
}

// CreateListView 创建 list view
// api 参考: http://172.19.89.76:48080/api/
func (j *JenkinsSdk) CreateListView(v *JenkinsView) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/createView?name=%s", j.Url, v.GetParentPath(), url.QueryEscape(v.Name))
	//fmt.Println(api)

	// getConfig
	configXml, err := j.makeListViewConfig(v)
	if err != nil {
		return err
	}
	//fmt.Println(configXml)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(configXml))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml")

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("create view error: " + err.Error())
	}

	return nil
}

// DeleteView 删除 view
func (j *JenkinsSdk) DeleteView(v *JenkinsView) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/doDelete", j.Url, v.GetFullPath())
	//fmt.Println(api)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("delete view error: " + err.Error())
	}

	return nil
}

// GetView 获取 view config.xml
// api 参考: http://172.19.89.76:48080/view/view123/api/
func (j *JenkinsSdk) GetView(v *JenkinsView) ([]byte, error) {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, v.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	data, err := j.sendHttp(req)
	if err != nil {
		return nil, errors.New("get view error: " + err.Error())
	}

	return data, nil
}

// UpdateViewDescription 通过form方式更改 description
// 没有对应 api, 从页面获取请求地址
// http://172.19.89.76:48080/view/view123/ 右上边有单独的 Edit description
func (j *JenkinsSdk) UpdateViewDescription(v *JenkinsView, description string) error {
	// 创建请求 URL
	api := fmt.Sprintf("%s%v/submitDescription", j.Url, v.GetFullPath())

	// 创建表单数据
	data := url.Values{}
	data.Set("description", description) // 假设 "description" 是表单字段名

	// 创建一个 POST 请求，设置表单数据
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头，内容类型为 x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update view error: " + err.Error())
	}

	return nil
}

// UpdateView 修改 view
// api 参考(只有全量更新): http://172.19.89.76:48080/view/view123/api/
// steps:
//  1. 获取 config.xml
//  2. 修改 config.xml
//  3. 通过 UpdateView 修改
func (j *JenkinsSdk) UpdateView(v *JenkinsView, configXml []byte) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, v.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(configXml))
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update view error: " + err.Error())
	}

	return nil
}

// AddJobToView 将 job 添加到 view
// 可以重复添加
// api: http://172.19.89.76:48080/view/view123/api/
func (j *JenkinsSdk) AddJobToView(v *JenkinsView, job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/addJobToView?name=%s", j.Url, v.GetFullPath(), url.QueryEscape(job.Name))

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("add job to view error: " + err.Error())
	}

	return nil
}

// RemoveJobFromView 将 job 从 view 移除
// api: http://172.19.89.76:48080/view/view123/api/
func (j *JenkinsSdk) RemoveJobFromView(v *JenkinsView, job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/removeJobFromView?name=%s", j.Url, v.GetFullPath(), url.QueryEscape(job.Name))

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("remove job from view error: " + err.Error())
	}

	return nil
}
