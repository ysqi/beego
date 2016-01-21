// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config is used to parse config
// Usage:
// import(
//   "github.com/astaxie/beego/config"
// )
//
//  cnf, err := config.NewConfig("ini", "config.conf")
//
//  cnf APIS:
//
//  cnf.Set(key, val string) error
//  cnf.String(key string) string
//  cnf.Strings(key string) []string
//  cnf.Int(key string) (int, error)
//  cnf.Int64(key string) (int64, error)
//  cnf.Bool(key string) (bool, error)
//  cnf.Float(key string) (float64, error)
//  cnf.DefaultString(key string, defaultVal string) string
//  cnf.DefaultStrings(key string, defaultVal []string) []string
//  cnf.DefaultInt(key string, defaultVal int) int
//  cnf.DefaultInt64(key string, defaultVal int64) int64
//  cnf.DefaultBool(key string, defaultVal bool) bool
//  cnf.DefaultFloat(key string, defaultVal float64) float64
//  cnf.DIY(key string) (interface{}, error)
//  cnf.GetSection(section string) (map[string]string, error)
//  cnf.SaveConfigFile(filename string) error
//
//  more docs http://beego.me/docs/module/config.md
package config

import (
	"fmt"
)

// Configer defines how to get and set value from configuration raw data.
// 配置文件操作接口，该接口定义一系列方法来操作配置文件，包括设置，配置参数等。
type Configer interface {
	//设置指定 Key 的 Value，支持多级模式，如： section::key 形式。
	Set(key, val string) error //support section::key type in given key when using ini type.
	//获取 KEY 对应的 Value，支持多级模式
	String(key string) string //support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	//获取 KEY 对应的 Value 组数据
	Strings(key string) []string //get string slice
	//获取 Int 类型的的 Value
	Int(key string) (int, error)
	//获取 Int64 类型的 Value
	Int64(key string) (int64, error)
	//获取 Bool 类型的 Value
	Bool(key string) (bool, error)
	//获取 Float64 类型的 Value
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultVal []string) []string //get string slice
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

// Config is the adapter interface for parsing config file to get raw data to Configer.
// 配置文件解析器接口。
type Config interface {
	//指定文件路径解析配置文件
	Parse(key string) (Configer, error)
	//提供配置数据进行解析
	ParseData(data []byte) (Configer, error)
}

// 已注册解析器实例池。
var adapters = make(map[string]Config)

// Register makes a config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
// 每实现一个配置文件解析器均需要注册，以便可以通过解析器名称来构建示例。
// 注意：不得重复注册，否则会panic。
func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	//判断是否已注册，如果已注册，则 Panic
	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}
	// 追加到解析器示例池中。
	adapters[name] = adapter
}

// NewConfig adapterName is ini/json/xml/yaml.
// filename is the config file path.
// 通过解析器解析配置文件，获得配置文件操作对象。当前可解析 ini/json/xml/yaml 格式文件。
func NewConfig(adapterName, filename string) (Configer, error) {

	//提取解析器，如果解析器不存在，则返回错误信息。
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	//对应解析器，解析配置文件。
	return adapter.Parse(filename)
}

// NewConfigData adapterName is ini/json/xml/yaml.
// data is the config data.
// 指定解析器直接解析配置数据，返回数配置数据操作对象。
func NewConfigData(adapterName string, data []byte) (Configer, error) {

	//提取解析器，如果解析器不存在，则返回错误信息。
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	//对应解析器，解析配置数据。
	return adapter.ParseData(data)
}
