package core

import (
	"github.com/wonderivan/logger"
)

type Module interface {
	Name() string
	Group() string
	Start()
	Enable() bool
}

var (
	modules         map[string]Module
	disabledModules map[string]Module
)

func init() {
	modules = make(map[string]Module)
	disabledModules = make(map[string]Module)
}

// Register 注册系统中需要的模块
func Register(m Module) {
	if m.Enable() {
		modules[m.Name()] = m
		logger.Info("模块 %v 注册成功", m.Name())
	} else {
		disabledModules[m.Name()] = m
		logger.Warn("Module %v is disabled, do not register", m.Name())
	}
}

// GetModules 获取module的map
func GetModules() map[string]Module {
	return modules
}
