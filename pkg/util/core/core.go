package core

import (
	"keep/pkg/util/loggerv1.0.1"
	"os"
	"os/signal"
	"syscall"

	beehiveContext "keep/pkg/util/core/context"
)

// StartModules starts modules that are registered
func StartModules() {
	beehiveContext.InitContext(beehiveContext.MsgCtxTypeChannel)
	modules := GetModules()

	for name, module := range modules {
		// Init the module
		beehiveContext.AddModule(name)
		// Assemble typeChannels for sendToGroup
		beehiveContext.AddModuleGroup(name, module.Group())
		go module.Start()
		logger.Info("启动模块：", name)
	}
}

// GracefulShutdown is if it gets the special signals it does modules cleanup
func GracefulShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	s := <-c
	logger.Info("获取到 ", s.String(), " 信号")

	// Cleanup each modules
	beehiveContext.Cancel()
	modules := GetModules()
	for name, module := range modules {
		logger.Warn("准备清理模块：", name)
		module.Cleanup()
		beehiveContext.Cleanup(name)
		logger.Warn("模块：" + name + " 清理完成...")
	}
}

// Run starts the modules and in the end does module cleanup
func Run() {
	// Address the module registration and start the core
	StartModules()
	// monitor system signal and shutdown gracefully
	GracefulShutdown()
}
