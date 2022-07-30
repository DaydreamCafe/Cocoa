// Package control 插件控制模块
package control

// EchoLevel 回显等级
type EchoLevel int

const (
	// NeverEcho 不回显任何信息
	NeverEcho EchoLevel = iota
	// OnlyEchoError 回显错误信息
	OnlyEchoError
	// OnlyEchoNoPremission 回显权限不足信息
	OnlyEchoNoPremission
	// EchoAny 回显任何信息
	EchoAny
)
