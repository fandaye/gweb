package modules

import (
	"log"
	"github.com/larspensjo/config"
)

type Config struct {
	ConfigFile string
}

func (c *Config) GlobalConfig() map[string]string {
	var r = make(map[string]string)
	cfg, err := config.ReadDefault(c.ConfigFile)
	if err != nil {
		log.Fatalf("Fail to find", c.ConfigFile, err)
	}

	if cfg.HasSection("global") {
		section, err := cfg.SectionOptions("global")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("global", v)
				if err == nil {
					r[v] = options
				}
			}
		}
	}
	return r
}

func (c *Config) MysqlConfig() map[string]string {
	var r = make(map[string]string)
	cfg, err := config.ReadDefault(c.ConfigFile)
	if err != nil {
		log.Fatalf("Fail to find", c.ConfigFile, err)
	}
	if cfg.HasSection("mysql") {
		section, err := cfg.SectionOptions("mysql")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mysql", v)
				if err == nil {
					r[v] = options
				}
			}
		}
	}
	return r
}


func (c *Config) RedisConfig() map[string]string {
	var r = make(map[string]string)
	cfg, err := config.ReadDefault(c.ConfigFile)
	if err != nil {
		log.Fatalf("Fail to find", c.ConfigFile, err)
	}

	if cfg.HasSection("redis") {
		section, err := cfg.SectionOptions("redis")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("redis", v)
				if err == nil {
					r[v] = options
				}
			}
		}
	}
	return r
}



func (c *Config) CheckMapCV(M map[string]string, V string) bool {
	for i := range M {
		if i == V {
			return true
		}
	}
	return false
}
