package config

import (
	"io/ioutil"
	"strings"
	"strconv"
	"common"
	"github.com/golang/glog"
	"os"
)

type ConfigManager struct {
	config map[string]string
}

func NewConfig(file string) (ret *ConfigManager) {
	ret = new(ConfigManager)
	ret.config = make(map[string]string)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) >= 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			ret.config[key] = value
		}
	}
	return
}

func (cf *ConfigManager) ReadString(key string) (string, bool) {
	data, ok := cf.config[key]
	return data, ok
}

func (cf *ConfigManager) ReadInt(key string) (int, bool) {
	data, ok := cf.config[key]
	if !ok {
		return 0, false
	}
	val, err := strconv.Atoi(data)
	if err != nil {
		return 0, false
	}
	return val, ok
}

func ReadDotConf(file string) (ret []common.InnerOuterMsg) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		msgs := strings.SplitN(line, "->", 2)
		if len(msgs) < 2 {
			glog.Warningf("outer_dot %s have not ->.\n", line)
			continue
		}

		msgs[0] = strings.TrimSpace(msgs[0])
		msgs[1] = strings.TrimSpace(msgs[1])

		tmp := strings.SplitN(msgs[0], ":", 2)
		if len(tmp) < 2 {
			glog.Warningf("outer_dot OuterIp:OuterPort %s have not :.\n", msgs[0])
			continue
		}
		ip := strings.TrimSpace(tmp[0])
		tmp[1] = strings.TrimSpace(tmp[1])
		port, err := strconv.Atoi(tmp[1])
		if err != nil {
			glog.Warningf("outer_dot OuterPort %s can't convert to int.\n", tmp[1])
			continue
		}
		omsg := common.OuterMsg{ip, port}

		tmp = strings.SplitN(msgs[1], ":", 2)
		if len(tmp) < 2 {
			glog.Warningf("outer_dot InnerIp:InnerPort %s have not :.\n", msgs[0])
			continue
		}
		ip = strings.TrimSpace(tmp[0])
		tmp[1] = strings.TrimSpace(tmp[1])
		port, err = strconv.Atoi(tmp[1])
		if err != nil {
			glog.Warningf("outer_dot InnerPort %s can't convert to int.\n", tmp[1])
			continue
		}
		imsg := common.InnerMsg{ip, port}

		ret = append(ret, common.InnerOuterMsg{
			InnerMsg: imsg,
			OuterMsg: omsg,
		})
	}
	return
}

func ReadArgs(id int, dv string) string {
	if len(os.Args) > id {
		return os.Args[id]
	}
	return dv
}