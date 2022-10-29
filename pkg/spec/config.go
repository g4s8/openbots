package spec

type Config struct {
	Api         *ApiConfig
	Persistence *PersistenceConfig `yaml:"persistence"`
}

type ApiConfig struct {
	Address string `yaml:"address"`
}

type PersistenceType string

const (
	MemoryPersistence   PersistenceType = "memory"
	DatabasePersistence PersistenceType = "database"
)

type PersistenceConfig struct {
	Type     PersistenceType `yaml:"type"`
	DBConfig *DBConfig       `yaml:"db_config"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	NoSSL    bool   `yaml:"no_ssl"`
}
