package handler

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
	"path/filepath"
	"sync"
	"github.com/micro/go-micro/config/encoder/json"
)

var (
	err error
	m                       sync.RWMutex
	inited                  bool
	baseConfig map[string]interface{}
)

// Init 初始化配置
func Init() {
	m.Lock()
	defer m.Unlock()
	if inited {
		log.Logf("[Init] 配置已经初始化过")
		return
	}
	// 创建新的配置
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("./", string(filepath.Separator))))

	e := json.NewEncoder()
	fmt.Println(appPath)
	fileSource := file.NewSource(
		file.WithPath(appPath+"/config/base.json"),
		source.WithEncoder(e),
	)
	conf := config.NewConfig()
	// 加载micro.yml文件
	if err = conf.Load(fileSource); err != nil {
		panic(err)
	}
	baseConfig = conf.Map()
	fmt.Println(baseConfig)
	// 侦听文件变动
	watcher, err := conf.Watch()
	if err != nil {
		log.Fatalf("[Init] 开始侦听应用配置文件变动 异常，%s", err)
		panic(err)
	}
	// 起一个协程监听配置变化
	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] 侦听应用配置文件变动 异常， %s", err)
				return
			}
			if err = conf.Load(fileSource); err != nil {
				panic(err)
			}
			log.Logf("[loadAndWatchConfigFile] 文件变动，%s", string(v.Bytes()))
		}
	}()
	// 标记已经初始化
	inited = true
}
