package core

import (
	"github.com/wonderivan/logger"
)

const (
	tryReadKeyTimes = 5
)

// Module interface
type Module interface {
	Name() string
	Group() string
	Start()
	Enable() bool
	Cleanup()
}

var (
	// Modules map
	modules         map[string]Module
	disabledModules map[string]Module
)

func init() {
	modules = make(map[string]Module)
	disabledModules = make(map[string]Module)
}

// Register register module
func Register(m Module) {
	if m.Enable() {
		modules[m.Name()] = m
		logger.Info("成功注册模块：", m.Name())
	} else {
		disabledModules[m.Name()] = m
		logger.Warn("模块未配置启动项目，当前不启动：", m.Name())
	}
}

// GetModules gets modules map
func GetModules() map[string]Module {
	return modules
}
