package jenkinsdk

import (
	"fmt"
	"regexp"
	"testing"
)

/**
* @Author: jack.walker
* @File: folder_test.go
* @CreateDate: 2025/3/30 12:46
* @ChangeDate：2025/3/30 12:46
* @Version：1.0.0
* @Description:
 */

// TestJenkinsSdk_CreateFolder 在 folder 中创建 folder
func TestJenkinsSdk_CreateFolder(t *testing.T) {

	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsFolder{
		Name:        "test123",
		Description: "Description2222222222222",
		//DisplayName: "DisplayName1111111111111",
	}

	err := j.CreateFolder(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("create folder success: %v", f)
	}

}

// TestJenkinsSdk_CreateFolderInFolder 在 folder 中创建 folder
func TestJenkinsSdk_CreateFolderInFolder(t *testing.T) {

	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsFolder{
		Name:        "test",
		Description: "test",
		DisplayName: "test-dis",
		Parent:      []string{"DEV"},
	}

	err := j.CreateFolder(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("create folder success: %v", f)
	}

}

func TestJenkinsSdk_UpdateFolder(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	f := &JenkinsFolder{
		Name:        "qa",
		Description: "我是qa哦",
		DisplayName: "folder-b2",
		//Parent:     []string{"hi-22"},
	}

	// 获取folder当前配置
	config, err := j.GetFolder(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("config old:", string(config))

	// 模拟 改变 config
	// 将 <description>folder-b description</description>
	// 改为 <description>folder-b xxxxx </description>
	re := regexp.MustCompile(`<description>.*?</description>`)

	newConfig := re.ReplaceAll(config, []byte("<description>"+f.Description+"</description>"))
	fmt.Println("newConfig:", string(newConfig))

	// 更改
	err = j.UpdateFolder(f, newConfig)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update folder success: %v", f)
	}
}

func TestJenkinsSdk_DeleteFolder(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "abc123")

	f := &JenkinsFolder{
		Name: "folder-a",
		//Parent:     []string{"hi-22"},
	}

	err := j.DeleteFolder(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("delete folder success: %v", f)
	}
}

// 不建议使用 UpdateFolderDescription
func TestJenkinsSdk_UpdateFolderDescription(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "abc123")
	f := &JenkinsFolder{
		Name:        "test123",
		DisplayName: "test123-folder",
		//Parent: []string{"dev"},
	}

	err := j.UpdateFolderDescription(f, "test123111...")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update folder description success: %v", f)
	}
}

func TestJenkinsSdk_GetFolder(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "abc123")

	f := &JenkinsFolder{
		Name: "folder-a",
		//Parent:     []string{"hi-22"},
	}

	data, err := j.GetFolder(f)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("get folder success: %v", data)
	}
}
