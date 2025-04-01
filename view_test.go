package jenkinsdk

import (
	"regexp"
	"testing"
)

/**
* @Author: jack.walker
* @File: view_test.go
* @CreateDate: 2025/3/31 09:39
* @Version: 1.0.0
* @Description:
 */

func TestJenkinsSdk_CreateListView(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name:        "view123",
		Description: "view123 Description",
	}

	err := j.CreateListView(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("create view success: %v", f)
	}
}

func TestJenkinsSdk_CreateListViewInFolder(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name:        "test1",
		Description: " test1 Description",
		Parent:      []string{"folder-b"},
	}

	err := j.CreateListView(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("create view success: %v", f)
	}
}

func TestJenkinsSdk_DeleteView(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name:   "test1",
		Parent: []string{"dev"},
	}

	err := j.DeleteView(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("delete view success: %v", f)
	}
}

func TestJenkinsSdk_GetView(t *testing.T) {

	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name:   "view1",
		Parent: []string{"dev"},
	}

	config, err := j.GetView(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("get view success: %v", string(RemoveXMLHeader(config)))
	}

}

func TestJenkinsSdk_UpdateView(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name:        "view1",
		Description: "我是view1 Description",
		Parent:      []string{"dev"},
	}

	// 获取folder当前配置
	config, err := j.GetView(f)
	if err != nil {
		t.Error(err)
	}

	// 模拟 改变 config
	// 将 <description>folder-b description</description>, 会存在 <description/> 或 没有 description字段 两种情形
	// 改为 <description>folder-b xxxxx </description>
	re := regexp.MustCompile(`<description>.*?</description>`)
	newConfig := re.ReplaceAll(config, []byte("<description>"+f.Description+"</description>"))

	err = j.UpdateView(f, newConfig)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update view success: %v", f)
	}
}

func TestJenkinsSdk_UpdateViewDescription(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsView{
		Name: "view123",
		//Parent: []string{"dev"},
	}

	err := j.UpdateViewDescription(f, "我是新的描述内容")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update view description success: %v", f)
	}
}
