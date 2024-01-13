package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Database struct {
	Host string `yaml:"host" mapstructure:"host"`
	User string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	Port uint `yaml:"port" mapstructure:"port"`
	DBName string `yaml:"db_name" mapstructure:"db_name"`

	Timeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
	ReadTimeout time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`

	ConnMaxLifetime int `yaml:"comm_max_lifetime" mapstructure:"conn_max_lifetime"`
	MaxOpenConns    int `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`

	MigrateFunc string `yaml:"migrate_func" mapstructure:"migrate_func"`
	Debug       bool   `yaml:"debug" mapstructure:"debug"`
	
}

func Load(path string, cfg interface{}) {
	viper.SetConfigType("yaml")
	if path != "" {
		plan, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		err = viper.ReadConfig(bytes.NewBuffer(plan))
		if err != nil {
			panic(err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("my config %+v\n", cfg)

}