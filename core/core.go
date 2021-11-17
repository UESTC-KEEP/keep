package core

import (
	"github.com/wonderivan/logger"
	beehiveContext "keep/core/context"
	"os"
	"os/signal"
	"syscall"
)

// StartModules starts modules that are registered
func StartModules() {
	beehiveContext.InitContext(beehiveContext.MsgCtxTypeChannel)

	modules := GetModules()
	for name, module := range modules {
		//Init the module
		beehiveContext.AddModule(name)
		//Assemble typeChannels for sendToGroup
		beehiveContext.AddModuleGroup(name, module.Group())
		go module.Start()
		logger.Info("启动模块： %v", name)
	}
}

// GracefulShutdown 截取到特殊信号时处理退出流程
func GracefulShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		logger.Info("截获到： %v信号", s.String())
		//Cleanup each modules
		beehiveContext.Cancel()
		modules := GetModules()
		for name, _ := range modules {
			logger.Info("开始清理模块 %v", name)
			beehiveContext.Cleanup(name)
		}
	}
}

// Run 启动注册了的module同时准备有序退出
func Run() {
	// Address the module registration and start the core
	StartModules()
	// monitor system signal and shutdown gracefully
	GracefulShutdown()
}
