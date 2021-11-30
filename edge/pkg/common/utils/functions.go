package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"io"
	"keep/pkg/util/loggerv1.0.1"
	"os"
)

func GracefulExit() {
	logger.Warn("准备退出...")
	os.Exit(1)
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
	if find, err := FindProcess("edgecore"); err != nil {
		return err
	} else if !find {
		return errors.New("kubeedge cloudcore未在运行,请检查")
	}
	logger.Debug("环境检测通过...")
	return nil
}

func copyFile(dstFileName string, srcFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}
	defer srcFile.Close()
	//通过srcFile，获取到Reader
	reader := bufio.NewReader(srcFile)
	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}
	writer := bufio.NewWriter(dstFile)
	defer func() {
		writer.Flush() //把缓冲区的内容写入到文件
		dstFile.Close()

	}()
	return io.Copy(writer, reader)

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
