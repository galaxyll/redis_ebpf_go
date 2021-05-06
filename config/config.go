package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Addr    string `yaml:"addr"`
	Grafana string `yaml:"grafana"`
}

type Influxdb struct {
	Addr      string `yaml:"addr"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Precision string `yaml:"precision"`
	Tag       Tag    `yaml:"tag"`
}

type Tag struct {
	Host    string `yaml:"host"`
	Service string `yaml:"service"`
}

type Plug struct {
	BinaryPath string `yaml:"binarypath"`
	Get        string `yaml:"get"`
	Set        string `yaml:"set"`
	Incr       string `yaml:"incr"`
	Decr       string `yaml:"decr"`
	Lpush      string `yaml:"lpush"`
	Rpush      string `yaml:"rpush"`
	Lpop       string `yaml:"lpop"`
	Rpop       string `yaml:"rpop"`
	Sadd       string `yaml:"sadd"`
	Hset       string `yaml:"hset"`
	Spop       string `yaml:"spop"`
	Mset       string `yaml:"mset"`
}

type Config struct {
	ServerConf   Server   `yaml:"server"`
	InfluxdbConf Influxdb `yaml:"influxdb"`
	PlugConf     Plug     `yaml:"plug"`
}

var Conf Config
var COMMAND map[string]string

func init() {
	bytes, err := ioutil.ReadFile("/home/nacl/redis_ebpf_go/config/config.yaml")
	if err != nil {
		log.Fatalf("read config from config.yam faild: %s\n", err)
	}
	err = yaml.Unmarshal(bytes, &Conf)
	if err != nil {
		log.Fatalf("faild tp parse config.yaml: %s\n", err)
	}
	COMMAND = map[string]string{
		"GET":   Conf.PlugConf.Get,
		"SET":   Conf.PlugConf.Set,
		"INCR":  Conf.PlugConf.Incr,
		"DECR":  Conf.PlugConf.Decr,
		"LPUSH": Conf.PlugConf.Lpush,
		"RPUSH": Conf.PlugConf.Rpush,
		"LPOP":  Conf.PlugConf.Lpop,
		"RPOP":  Conf.PlugConf.Rpop,
		"SADD":  Conf.PlugConf.Sadd,
		"HSET":  Conf.PlugConf.Hset,
		"SPOP":  Conf.PlugConf.Spop,
		"MSET":  Conf.PlugConf.Mset,
	}
}
