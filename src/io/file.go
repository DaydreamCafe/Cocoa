// Package io 常用io相关函数库
package io

import "io/ioutil"

// ReadConfig 读取配置文件并返回配置文件内容
func ReadConfig() ([]byte, error) {
	return ioutil.ReadFile("config.yaml")
}