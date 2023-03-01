package configinit

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	//  DB
	DBName     string
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string

	// Gin
	HostIp   string
	HostPort string

	// OpenTelemetry
	ServiceName      string
	CollectorURLIP   string
	CollectorURLPort string
	Insecure         string
)

func LoadEnv() {

	var (
		_, b, _, _ = runtime.Caller(0)

		projectRootPath = filepath.Join(filepath.Dir(b), "../../")
	)
	err := godotenv.Load(os.ExpandEnv(projectRootPath + "/.env"))
	if err != nil {
		log.Printf("Error getting env %v\n", err)
	}

	DBName = os.Getenv("DB_DBNAME")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUsername = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")

	HostIp = os.Getenv("HOST_IP")
	HostPort = os.Getenv("HOST_PORT")

	ServiceName = os.Getenv("SERVICE_NAME")
	CollectorURLIP = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT_IP")
	CollectorURLPort = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT_PORT")
	Insecure = os.Getenv("INSECURE_MODE")
}
