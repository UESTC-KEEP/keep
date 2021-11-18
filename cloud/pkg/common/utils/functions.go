package utils

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/wonderivan/logger"
	naive_engine "keep/cloud/pkg/k8sclient/naive-engine"
	"keep/constants"
	"sync"
)

func GracefulExit() {
	logger.Info("准备退出...")
	deleteRedis()
}

// FindProcess 根据进程名找当前是不是有进程在执行
func FindProcess(name string) (bool, error) {
	processes, err := ps.Processes()
	if err != nil {
		return false, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return true, nil
		}
	}
	return false, nil
}

// deleteRedis 删除redis各组件
func deleteRedis() {
	ns := constants.DefaultNameSpace
	var compoenents = []string{constants.DefaultRedisConfigMap, constants.DefaultRedisSVC, constants.DefaultRedisStatefulSet}
	var wg sync.WaitGroup
	for i := 0; i < len(compoenents); i++ {
		wg.Add(1)
		go func(i int) {
			logger.Debug("开始删除redis组件 " + compoenents[i] + " ...")
			err := naive_engine.DeleteResourceByYAML(compoenents[i], ns)
			if err != nil {
				wg.Done()
				logger.Error(err)
				return
			} else {
				wg.Done()
			}
		}(i)
	}
	wg.Wait()
}

// EnvironmentCheck  check the environment before keep start
// if Check failed,  return errors
func EnvironmentCheck() error {
	// if kubelet is running, return error
	if find, err := FindProcess("cloudcore"); err != nil {
		return err
	} else if !find {
		return errors.New("kubeedge edgecore未在运行,请检查")
	}
	logger.Debug("环境检测通过...")
	return nil
}

func PrintKEEPLogo() {
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, "  __       ___    __________     ___________   ____________                                  __     __      \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |  |     |  |   |  ________|   |  ________|  |   ______   |                                |  |   |  |   \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |  |   |  |     |  |           |  |          |  |      |  | 	 _      _    _______       |  |   |  |     \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |   | |  |      |  |_______    |  |_______   |  |______|  |    | |    | |  |  ___  |     |  |   |  |   \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |   ||  |       |   _______|   |  _______|   |   _________|    | |    | |  | |___| |    |  |   |  |       \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |  |   |  |     |  |           |  |          |  |              | |    | |  | ______|   |__|   |__|         \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |  |    |  |    |  |_______    |  |_______   |  |              | |____| |  | |         __     __     \n", 0x1B)
	fmt.Printf("%c[1;5;34m%s%c[0m", 0x1B, " |__|     |__|   |__________|   |__________|  |__|    V 1.0.0   |________|  |_|        |__|   |__|          \n", 0x1B)
}
