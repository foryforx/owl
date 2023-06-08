package conf

import (
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
)

var conf *Configuration
var onceEnv sync.Once

type Configuration struct {
	Environment        string
	JwtSecret          string
	Port               string
	DBHost             string
	DBPort             string
	DBDatabase         string
	DBUser             string
	DBPwd              string
	DBSSLMode          string
	DBMaxConnection    int
	SessionStoreHost   string
	SessionStorePort   string
	AcquireConnTimeout int
}

func loadEnvOrExit(varName string) string {
	envVar := os.Getenv(varName)

	if envVar == "" {
		log.Fatal("Env var ", varName, " is required but not defined")
	}

	return envVar
}

func loadEnvOrDefault(varName string, defaultValue string) string {
	envVar := os.Getenv(varName)

	if envVar == "" {
		log.Info("Using default value ", defaultValue, " for ", varName)

		return defaultValue
	}

	return envVar
}

// GetConfiguration instantiates Configuration
func GetConfiguration() *Configuration {
	onceEnv.Do(func() {
		conf = &Configuration{}

		conf.Port = loadEnvOrDefault("PORT", "5000")

		conf.DBHost = loadEnvOrExit("PGHOST")
		conf.DBPort = loadEnvOrExit("PGPORT")
		conf.DBDatabase = loadEnvOrExit("PGDATABASE")
		conf.DBUser = loadEnvOrExit("PGUSER")
		conf.DBPwd = loadEnvOrDefault("PGPASSWORD", "")
		conf.DBSSLMode = loadEnvOrExit("PGSSLMODE")
		conf.DBMaxConnection, _ = strconv.Atoi(loadEnvOrDefault("PGMAXCONNECTIONS", "20"))
		conf.AcquireConnTimeout, _ = strconv.Atoi(loadEnvOrDefault("DB_CONN_ACQUIRE_TIMEOUT", "30"))
		conf.JwtSecret = loadEnvOrExit("JWT_SECRET")
	})

	return conf
}

// GetConnConfig returns a pointer to an instance of *pgx.ConnConfig
// https://godoc.org/github.com/jackc/pgx#ConnConfig
func (c *Configuration) GetDBConfig() (*pgx.ConnConfig, error) {
	config := &pgx.ConnConfig{}

	if c.DBHost == "" || c.DBUser == "" || c.DBDatabase == "" {

		return config, errors.New("ConnectionString is invalid")

	}
	portValue, err := strconv.Atoi(c.DBPort)
	if err != nil {
		return config, err
	}
	config.Host = c.DBHost
	config.Port = uint16(portValue)
	config.User = c.DBUser
	config.Password = c.DBPwd
	config.Database = c.DBDatabase

	return config, nil
}

// Print logs current configuration to stdout
func (c *Configuration) Print() {
	log.Debugln("Loaded configuration with settings")
	log.Debugln("DBHost: ", c.DBHost)
	log.Debugln("DBPort: ", c.DBPort)
	log.Debugln("DBUser: ", c.DBUser)
	log.Debugln("DBPassword: ", len(c.DBPwd), " characters")
	log.Debugln("SSLMode: ", c.DBSSLMode)
	log.Debugln("Database name: ", c.DBDatabase)
	log.Debugln("Max database connections:", strconv.Itoa(c.DBMaxConnection))
}
