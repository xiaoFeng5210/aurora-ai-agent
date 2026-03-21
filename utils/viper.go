package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// FileType
const (
	JSON = "json"
	YAML = "yaml"
	ENV  = "env"
)

// Viper可以解析JSON、TOML、YAML、HCL、INI、ENV等格式的配置文件。甚至可以监听配置文件的变化(WatchConfig)，不需要重启程序就可以读到最新的值。
func InitViper(dir, file, FileType string) *viper.Viper {
	config := viper.New()
	for _, configPath := range buildConfigSearchPaths(dir) {
		config.AddConfigPath(configPath)
	}
	config.SetConfigName(file)     // 文件名(不带路径，不带后缀)
	config.SetConfigType(FileType) // 文件类型

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("解析配置文件%s出错:%s", filepath.Join(dir, file)+"."+FileType, err)) //系统初始化阶段发生任何错误，直接结束进程。logger还没初始化，不能用logger.Fatal()
	}

	return config
}

func buildConfigSearchPaths(dir string) []string {
	paths := []string{dir}
	seen := map[string]struct{}{dir: {}}

	workingDir, err := os.Getwd()
	if err != nil {
		return paths
	}

	for current := workingDir; ; current = filepath.Dir(current) {
		candidate := filepath.Join(current, dir)
		if _, ok := seen[candidate]; !ok {
			paths = append(paths, candidate)
			seen[candidate] = struct{}{}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}

	return paths
}
