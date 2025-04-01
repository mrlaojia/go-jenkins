package jenkinsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/**
* @Author: jack.walker
* @File: csrf.go
* @CreateDate: 2025/3/29 15:43
* @ChangeDate：2025/3/29 15:43
* @Version：1.0.0
* @Description:
 */

// GetJenkinsCrumb  获取 Jenkins Crumb 处理 CSRF
func (j *JenkinsSdk) GetJenkinsCrumb() (string, error) {

	crumbURL := fmt.Sprintf("%s/crumbIssuer/api/json", j.Url)

	req, _ := http.NewRequest("GET", crumbURL, nil)
	j.setBasicAuth(req)

	resp, err := j.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Crumb string `json:"crumb"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.Crumb, nil
}

// MakeCsrf 处理 CSRF (如果启用)
func (j *JenkinsSdk) makeCsrf(req *http.Request) error {

	if j.csrfEnable && strings.ToLower(req.Method) != "get" {
		crumb, err := j.GetJenkinsCrumb()
		if err != nil {
			return err
		}

		req.Header.Add("Jenkins-Crumb", crumb)
	}

	return nil
}
