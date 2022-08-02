// Package control 插件控制模块
package control

import zero "github.com/wdvxdr1123/ZeroBot"

func getEngine(pluginMetadata Metadata, echoLevel EchoLevel) zero.Engine {
	// 初始化engine对象
	engine := zero.Engine{}

	// 添加prehandler
	// 添加忽略账号prehandler
	engine.UsePreHandler(ignoreUserChecker)

	// 插件prehandler
	engine.UsePreHandler(pluginChecker(pluginMetadata, echoLevel))

	// 返回engine
	return engine
}
