/*
@Title        plugin_init.go
@Description  插件相关基础操作
@Author       WhitePaper233 2020.7.13
@Update       WhitePaper233 2020.7.13
*/
package plugin

import (
	logger "github.com/sirupsen/logrus"

	_ "github.com/DaydreamCafe/Cocoa/V2/src/logger"
)

type Metadata struct {
	Name        string
	Version     string
	Description string
	Author      string
}

var PluginMetadata Metadata

func Initialization(metadata Metadata) {
	PluginMetadata = metadata
	logger.Infof("%s 插件加载成功", metadata.Name)
}