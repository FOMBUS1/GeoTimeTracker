package config

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     uint8  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port uint8  `yaml:"port"`
}
