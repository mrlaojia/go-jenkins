package jenkinsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/**
* @Author: jack.walker
* @File: job.go
* @CreateDate: 2025/4/1 13:44
* @Version: 1.0.0
* @Description:
 */

type JenkinsJob struct {
	Name      string
	ConfigXml string
	Parent    []string
}

// GetParentPath 构建 Jenkins job 路径
// []string{"a", "b"} → /job/a/job/b
func (f *JenkinsJob) GetParentPath() string {
	var builder strings.Builder
	for _, elem := range f.Parent {
		builder.WriteString("/job/")
		builder.WriteString(url.PathEscape(elem))
	}
	return builder.String()
}

func (f *JenkinsJob) GetFullPath() string {
	return f.GetParentPath() + "/job/" + f.Name
}

// CreateJob 创建 job
// api 参考: http://172.19.89.76:48080/job/job1/api/
func (j *JenkinsSdk) CreateJob(job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/createItem?name=%s", j.Url, job.GetParentPath(), url.QueryEscape(job.Name))
	//fmt.Println(api)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(job.ConfigXml))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/xml")

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("create job error: " + err.Error())
	}

	return nil
}

// CopyJob 复制 job
// api 参考: http://172.19.89.76:48080/job/job1/api/
// 使用 Jenkins 远程访问 API 从现有作业中复制作业。但是，除非我转到 UI 并按下配置，然后单击保存（即使没有进行任何更改），否则无法构建此作业（没有“立即构建”按钮来）。
// 默认情况下，您复制的任何作业都是禁用的（或者说，不可构建）
// 故 copy 后, 需要 先禁用job，再启用job 才可以直接build
// 参考网文 https://stackoverflow.com/questions/66620833/jenkins-api-created-jobs-dont-have-build-now-button
func (j *JenkinsSdk) CopyJob(job, from *JenkinsJob) error {

	fromList := append(from.Parent, from.Name)
	fromPath := strings.Join(fromList, "/")

	fmt.Println(from.Parent, from.Name, fromList)

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/createItem?name=%s&mode=copy&from=%v", j.Url, job.GetParentPath(), url.QueryEscape(job.Name), url.QueryEscape(fromPath))
	//fmt.Println(api)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("copy job error: " + err.Error())
	}

	return nil
}

// EnableJob 启用job
// 没有找到 api
// 参考网文 https://stackoverflow.com/questions/66620833/jenkins-api-created-jobs-dont-have-build-now-button
func (j *JenkinsSdk) EnableJob(job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/enable", j.Url, job.GetFullPath())
	//fmt.Println(api)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("enable job error: " + err.Error())
	}

	return nil
}

// DisableJob 禁用job
// 没有找到 api，参考网文
// https://stackoverflow.com/questions/66620833/jenkins-api-created-jobs-dont-have-build-now-button
func (j *JenkinsSdk) DisableJob(job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/disable", j.Url, job.GetFullPath())
	//fmt.Println(api)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("disable job error: " + err.Error())
	}

	return nil
}

// DeleteJob 删除 job
// api 参考: http://172.19.89.76:48080/job/job1/api/
func (j *JenkinsSdk) DeleteJob(job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/", j.Url, job.GetFullPath())

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

// GetJob 获取 job config.xml
// api 参考: http://172.19.89.76:48080/job/job1/api/
func (j *JenkinsSdk) GetJob(v *JenkinsJob) ([]byte, error) {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, v.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	data, err := j.sendHttp(req)
	if err != nil {
		return nil, errors.New("get job error: " + err.Error())
	}

	return data, nil
}

// UpdateJobDescription 修改描述
// api 参考: http://172.19.89.76:48080/job/job1/api/
// 仅支持 流水线 和 自由风格；像：文件夹，多分支流水线 不支持
func (j *JenkinsSdk) UpdateJobDescription(job *JenkinsJob, desc string) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/description", j.Url, job.GetFullPath())

	formData := url.Values{}
	formData.Set("description", desc)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update job description error: " + err.Error())
	}

	return nil
}

// UpdateJob 修改 job
// api 参考(只有全量更新):  http://172.19.89.76:48080/job/job1/api/
// steps:
//  1. 获取 config.xml
//  2. 修改 config.xml
//  3. 通过 UpdateJob 修改
func (j *JenkinsSdk) UpdateJob(job *JenkinsJob) error {

	// 创建请求 URL
	api := fmt.Sprintf("%s%v/config.xml", j.Url, job.GetFullPath())

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", api, bytes.NewBufferString(job.ConfigXml))
	if err != nil {
		return err
	}

	_, err = j.sendHttp(req)
	if err != nil {
		return errors.New("update job error: " + err.Error())
	}

	return nil
}
