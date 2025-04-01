package jenkinsdk

import "testing"

/**
* @Author: jack.walker
* @File: csrf_test.go
* @CreateDate: 2025/3/30 15:02
* @ChangeDate：2025/3/30 15:02
* @Version：1.0.0
* @Description:
 */

func TestJenkinsSdk_GetJenkinsCrumb(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "abc123")

	data, err := j.GetJenkinsCrumb()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("get csrf success: %v", data)
	}

}
