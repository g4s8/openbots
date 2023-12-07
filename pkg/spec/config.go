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
	// Api configuration specify API server parameters.
	Api *ApiConfig `yaml:"api"`
	// Persistence configuration specify persistence type and its parameters.
	Persistence *PersistenceConfig `yaml:"persistence"`
	// Assets configuration specify assets providers and its parameters.
	Assets *AssetsConfig `yaml:"assets"`
	// PaymentProviders is for payment providers tokens and parameters.
	PaymentProviders []PaymentProvider `yaml:"paymentProviders"`
}

type ApiConfig struct {
	// Address is the address to listen on.
	Address string `yaml:"address"`
}

type PersistenceType string

const (
	MemoryPersistence   PersistenceType = "memory"
	DatabasePersistence PersistenceType = "database"
)

type PersistenceConfig struct {
	// Type is the type of persistence to use.
	Type PersistenceType `yaml:"type"`
	// DBConfig is the configuration for database persistence.
	DBConfig *DBConfig `yaml:"db_config"`
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
