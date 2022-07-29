// Package control 插件控制模块
package control

type EchoLevel int

const (
	NeverEcho            EchoLevel = iota // 不回显任何信息
	OnlyEchoError                         // 回显错误信息
	OnlyEchoNoPremission                  // 回显权限不足信息
	EchoAny                               // 回显任何信息
)
