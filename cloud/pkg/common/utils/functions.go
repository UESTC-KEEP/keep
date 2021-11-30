package utils

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/wonderivan/logger"
)

func GracefulExit() {
	logger.Warn("准备退出...")
	//os.Exit(1)
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
