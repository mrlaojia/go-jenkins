package jenkinsdk

import "bytes"

/**
* @Author: jack.walker
* @File: utils.go
* @CreateDate: 2025/3/31 14:32
* @Version: 1.0.0
* @Description:
 */

// RemoveXMLHeader
// 移除 XML 版本声明 (<?xml version='1.1' encoding='UTF-8'?>) 的方法
func RemoveXMLHeader(xmlData []byte) []byte {
	// 查找 XML 声明结束位置
	endDecl := bytes.Index(xmlData, []byte("?>"))
	if endDecl == -1 {
		return xmlData
	}

	// 跳过声明部分并去除开头空白
	return bytes.TrimSpace(xmlData[endDecl+2:])
}
