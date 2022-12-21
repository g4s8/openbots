package spec

var DefaultConfig = &Config{
	Persistence: &PersistenceConfig{
		Type: MemoryPersistence,
	},
	Assets: &AssetsConfig{
		Provider: "fs",
	},
}

type Config struct {
	Api              *ApiConfig         `yaml:"api"`
	Persistence      *PersistenceConfig `yaml:"persistence"`
	Assets           *AssetsConfig      `yaml:"assets"`
	PaymentProviders []PaymentProvider  `yaml:"paymentProviders"`
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

type AssetsConfig struct {
	Provider string            `yaml:"provider"`
	Params   map[string]string `yaml:"params"`
}

type PaymentProvider struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}
