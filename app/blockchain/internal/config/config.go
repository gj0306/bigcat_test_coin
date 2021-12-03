package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	// Config ...
	Config     *YamlConfig
	configOnce sync.Once
)

type AddressConf struct {
	PrivateKey string `yaml:"private_key"`
}

type NetworkConf struct {
	Port     uint16 `yaml:"port"`
	IsCreate bool   `yaml:"is_create"`
	Addr     string `yaml:"addr"`
}

type LogConf struct {
	LogPath     string `yaml:"log_path"`
	LogFileName string `yaml:"log_file_name"`
	Debug       bool   `yaml:"debug"`
}

type DbConf struct {
	DriverName string `yaml:"driver_name"`
	Dir        string `yaml:"dir"`
}

type WebServer struct {
	HttpAddr string `yaml:"http_addr"`
	HttpTimeout int64 `yaml:"http_timeout"`
	GrpcAddr string `yaml:"grpc_addr"`
	GrpcTimeout int64 `yaml:"grpc_timeout"`
	Swagger bool `yaml:"swagger"`
	Open bool `yaml:"open"`
}

// YamlConfig 配置文件
type YamlConfig struct {
	Debug     bool `yaml:"debug"`
	IsGenesis bool `yaml:"is_genesis"`
	Address   *AddressConf
	Network   *NetworkConf
	Log       *LogConf
	Db        *DbConf
	Server    *WebServer `yaml:"server"`
}

// NewConfig 获取配置文件对象
func NewConfig(path string) *YamlConfig {
	var err error
	if Config != nil {
		return Config
	}
	if strings.TrimSpace(path) == "" {
		log.Fatal("配置文件对象为nil,需要传递配置文件路径")
	}
	configOnce.Do(func() {
		err = loadConfig(path, &Config)
		if err != nil {
			log.Fatal("加载配置文件失败:", err.Error())
		}
	})
	return Config
}

// loadConfig 加载配置文件
func loadConfig(path string, cfg **YamlConfig) error {
	var err error
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configContent, cfg)
	if err != nil {
		return err
	}
	return nil
}

// GenerateConfig 生成配置文件
func GenerateConfig(path string, cfg *YamlConfig) error {
	var err error
	_tmp, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(_tmp)
	if err != nil {
		return err
	}
	return nil
}

// LoadConfigBase 配置
func LoadConfigBase(path string, cfg interface{}) error {
	var err error
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configContent, cfg)
	if err != nil {
		return err
	}
	return nil
}

func GenerateConfigBase(path string, cfg interface{}) error {
	var err error
	_tmp, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(_tmp)
	if err != nil {
		return err
	}
	return nil
}
