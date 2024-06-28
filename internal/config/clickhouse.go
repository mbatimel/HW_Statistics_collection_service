package config

type ClickHouse struct {
	Host 		string `yaml:"host"`
	Port 		string `yaml:"port"`
	DB   		string `yaml:"db"`
	Username	string `yaml:"username"`
	Password	string `yaml:"password"`
}