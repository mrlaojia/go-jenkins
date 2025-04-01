package jenkinsdk

import (
	"testing"
)

/**
* @Author: jack.walker
* @File: info_test.go
* @CreateDate: 2025/3/30 16:37
* @ChangeDate：2025/3/30 16:37
* @Version：1.0.0
* @Description:
 */

func TestJenkinsSdk_GetJenkinsInfo(t *testing.T) {

	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	info, err := j.GetJenkinsInfo()
	if err != nil {
		t.Error(err)
	}

	// 输出诊断信息
	t.Logf("Jenkins 基本信息：\n")
	t.Logf("• 安全启用：%t\n", info.UseSecurity)
	t.Logf("• CSRF 保护：%t\n", info.UseCrumbs)
}
