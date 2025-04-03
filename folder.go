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
* @File: folder.go
* @CreateDate: 2025/3/29 15:33
* @ChangeDate：2025/3/29 15:33
* @Version：1.0.0
* @Description:
 */

type JenkinsFolder struct {
	Name        string
	Description string
	DisplayName string
	Parent      []string
}

// GetParentPath 构建 Jenkins job 路径
// []string{"a", "b"} → /job/a/job/b
func (f *JenkinsFolder) GetParentPath() string {
	var builder strings.Builder
	for _, elem := range f.Parent {
		builder.WriteString("/job/")
		builder.WriteString(url.PathEscape(elem))
	}
	return builder.String()
}

func (f *JenkinsFolder) GetFullPath() string {
	return f.GetParentPath() + "/job/" + f.Name
}

// JenkinsFolderTemplate 生成的 template
type JenkinsFolderTemplate struct {
	XMLName     xml.Name `xml:"com.cloudbees.hudson.plugins.folder.Folder"`
	Plugin      string   `xml:"plugin,attr"`           // Plugin 属性
	Description string   `xml:"description,omitempty"` // 描述内容
	DisplayName string   `xml:"displayName,omitempty"` // 显示名称
	Properties  struct{} `xml:"properties"`            // 空属性
}

// makeFolderConfig 创建 config.xml
// config.xml模板参考: http://172.19.89.76:48080/job/folder/config.xml
func (j *JenkinsSdk) makeFolderConfig(f *JenkinsFolder) (string, error) {

	// 解析 XML 到结构体
	var folderTemplate JenkinsFolderTemplate
	// 修改 Plugin 和 Description
	folderTemplate.Description = f.Description
	folderTemplate.DisplayName = f.DisplayName
	// 获取安装的version
	plug, err := j.GetJenkinsPlugin("cloudbees-folder")
	if err != nil {
		return "", err
	}
	folderTemplate.Plugin = fmt.Sprintf("cloudbees-folder@%v", plug.Version) // 新插件版本

	// 生成 XML
	newXML, err := xml.MarshalIndent(folderTemplate, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(newXML), nil
}

// CreateFolder 创建 folder
// api 参考: http://172.19.89.76:48080/job/folder/api/
func (j *JenkinsSdk) CreateFolder(f *JenkinsFolder) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/createItem?name=%s", j.Url, f.GetParentPath(), url.QueryEscape(f.Name))
	//fmt.Println(api)

	// getConfig
	configXml, err := j.makeFolderConfig(f)
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
		return errors.New("create folder error: " + err.Error())
	}

	return nil
}

// DeleteFolder 删除 folder
// 自动(级联)删除文件夹下所有
// api 参考: http://172.19.89.76:48080/job/folder/api/
func (j *JenkinsSdk) DeleteFolder(f *JenkinsFolder) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/", j.Url, f.GetFullPath())

	// getConfig
	//configXml, err := j.makeFolderConfig(f)
	//if err != nil {
	//	return err
	//}

	// 创建 HTTP 请求
	req, err := http.NewRequest("DELETE", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("delete folder error: " + err.Error())
	}

	return nil
}

// GetFolder 获取 folder config.xml
// api 参考: http://172.19.89.76:48080/job/folder/api/
func (j *JenkinsSdk) GetFolder(f *JenkinsFolder) ([]byte, error) {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, f.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	data, err := j.sendHttp(req)
	if err != nil {
		return nil, errors.New("get folder error: " + err.Error())
	}

	return data, nil
}

// UpdateFolderDescription 通过form方式更改 description
// 没有对应 api, 从[配置]页面获取请求地址 http://172.19.89.76:48080/job/qa/configure
// form 字段中 json 的值会被应用
func (j *JenkinsSdk) UpdateFolderDescription(f *JenkinsFolder, description string) error {
	// 创建请求 URL
	api := fmt.Sprintf("%s%v/configSubmit", j.Url, f.GetFullPath())

	// 创建表单数据
	data := url.Values{}
	// 保留 description 和 图标 和 displayNameOrNull 字段
	// 如果folder还设置了其他字段值，可能会丢
	data.Set("json", fmt.Sprintf(`{"displayNameOrNull":"%v","description":"%v","icon":{"stapler-class":"com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon","$class":"com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"}}`, f.DisplayName, description))

	// 创建一个 POST 请求，设置表单数据
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头，内容类型为 x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update folder error: " + err.Error())
	}

	return nil
}

// UpdateFolder 修改 folder
// api 参考(只有全量更新): http://172.19.89.76:48080/job/folder/api/
// steps:
//  1. 获取 config.xml
//  2. 修改 config.xml
//  3. 通过 UpdateFolder 修改
func (j *JenkinsSdk) UpdateFolder(f *JenkinsFolder, configXml []byte) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, f.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(configXml))
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update folder error: " + err.Error())
	}

	return nil
}
