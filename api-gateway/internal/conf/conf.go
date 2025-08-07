package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"time"
)

// Bootstrap is the top-level configuration struct that maps to the YAML file.
type Bootstrap struct {
	Server *Server `yaml:"server"`
	Data   *Data   `yaml:"data"`
}

// Server holds the configuration for HTTP and gRPC servers.
type Server struct {
	HTTP *Server_HTTP `yaml:"http"`
	GRPC *Server_GRPC `yaml:"grpc"`
}

type Server_HTTP struct {
	Network string        `yaml:"network"`
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
	TLS     *Server_TLS   `yaml:"tls"`
}

type Server_TLS struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type Server_GRPC struct {
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}

// Data holds the configuration for database and redis.
type Data struct {
	Database *Data_Database `yaml:"database"`
	Redis    *Data_Redis    `yaml:"redis"`
}

type Data_Database struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type Data_Redis struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// LoadBootstrap loads the configuration from the `configs/config.yaml` file.
func LoadBootstrap() (*Bootstrap, error) {
	c := config.New(
		config.WithSource(
			file.NewSource("./configs"),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return nil, err
	}

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}
