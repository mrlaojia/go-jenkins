package jenkinsdk

import "testing"

/**
* @Author: jack.walker
* @File: plugin_test.go
* @CreateDate: 2025/3/30 16:28
* @ChangeDate：2025/3/30 16:28
* @Version：1.0.0
* @Description:
 */

func TestJenkinsSdk_GetAllJenkinsPlugins(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb").EnableCsrf()

	plugins, err := j.GetAllJenkinsPlugins()
	if err != nil {
		t.Error(err)
	}

	for _, plugin := range plugins.Plugins {
		t.Log(plugin)
	}

}
