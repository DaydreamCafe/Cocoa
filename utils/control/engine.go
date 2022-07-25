// Package control 插件控制模块
package control

import zero "github.com/wdvxdr1123/ZeroBot"

func getEngine(pluginMetadata Metadata) zero.Engine {
	// 初始化engine对象
	engine := zero.Engine{}

	// 添加prehandler
	// 插件prehandler
	engine.UsePreHandler(pluginCheck(pluginMetadata))

	// 返回engine
	return engine
}