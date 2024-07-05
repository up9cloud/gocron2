package setting

import (
	_ "embed"
	"errors"

	"github.com/up9cloud/gocron2/internal/modules/logger"
	"github.com/up9cloud/gocron2/internal/modules/utils"
	"gopkg.in/ini.v1"
)

const DefaultSection = "default"

type Setting struct {
	Db struct {
		Engine       string `ini:"db.engine"`
		Host         string `ini:"db.host"`
		Port         int `ini:"db.port"`
		User         string `ini:"db.user"`
		Password     string `ini:"db.password"`
		Database     string `ini:"db.database"`
		Prefix       string `ini:"db.prefix"`
		Charset      string `ini:"db.charset"`

		MaxIdleConns int `ini:"db.max.idle.conns"`
		MaxOpenConns int `ini:"db.max.open.conns"`

		Sslmode       string `ini:"db.sslmode"`
		SslCaFile     string `ini:"db.ssl_ca_file"`
		SslCertFile   string `ini:"db.ssl_cert_file"`
		SslKeyFile    string `ini:"db.ssl_key_file"`
		SslServerName string `ini:"db.ssl_server_name"`
	}
	AllowIps             string `ini:"allow_ips"`
	AllowLocalhostBypass bool `ini:"allow_localhost_bypass"`
	AppName              string `ini:"app.name"`
	ApiKey               string `ini:"api.key"`
	ApiSecret            string `ini:"api.secret"`
	ApiSignEnable        bool `ini:"api.sign.enable"`

	EnableTLS bool `ini:"enable_tls"`
	CAFile    string `ini:"ca_file"`
	CertFile  string `ini:"cert_file"`
	KeyFile   string `ini:"key_file"`

	ConcurrencyQueue int `ini:"concurrency.queue"`
	AuthSecret       string `ini:"auth_secret"`
}

// 读取配置
func Read(filename string) (*Setting, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	s := new(Setting)
	err = config.Section(DefaultSection).MapTo(&s.Db)
	if err != nil {
		return nil, err
	}
	err = config.Section(DefaultSection).MapTo(s)
	if err != nil {
		return nil, err
	}
	if s.EnableTLS {
		if !utils.FileExist(s.CAFile) {
			logger.Fatalf("failed to read ca cert file: %s", s.CAFile)
		}

		if !utils.FileExist(s.CertFile) {
			logger.Fatalf("failed to read client cert file: %s", s.CertFile)
		}

		if !utils.FileExist(s.KeyFile) {
			logger.Fatalf("failed to read client key file: %s", s.KeyFile)
		}
	}
	if s.AuthSecret == "" {
		s.AuthSecret = utils.RandAuthToken()
	}
	logger.Infof("config load: %#v\n", s)
	return s, nil
}

// 写入配置
//go:embed app.ini
var defaultConfig []byte

func Write(rawConfig []string, filename string) error {
	if len(rawConfig) == 0 {
		return errors.New("参数不能为空")
	}
	if len(rawConfig)%2 != 0 {
		return errors.New("参数不匹配")
	}
	// config := ini.Empty()
	config, _ := ini.Load(defaultConfig)

	ini.PrettyFormat = false
	ini.PrettyEqual = true

	section, err := config.NewSection(DefaultSection)
	if err != nil {
		return err
	}
	for i := 0; i < len(rawConfig); {
		_, err = section.NewKey(rawConfig[i], rawConfig[i+1])
		if err != nil {
			return err
		}
		i += 2
	}
	err = config.SaveTo(filename)

	return err
}
