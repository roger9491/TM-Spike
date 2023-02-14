package integration__test

import (
	"TM-Spike/init/sqlinit"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
)

func database(t *testing.T) *gorm.DB {
	var (
		_, b, _, _ = runtime.Caller(0)

		projectRootPath = filepath.Join(filepath.Dir(b), "../")
	)
	fmt.Println(projectRootPath)
	err := godotenv.Load(os.ExpandEnv(projectRootPath + "/.env"))
	if err != nil {
		log.Printf("Error getting env %v\n", err)
	}

	username := os.Getenv("USERNAME_TEST")
	password := os.Getenv("PASSWORD_TEST")
	database := os.Getenv("DATABASE_TEST")
	fmt.Println("bbbbbba",username)
	ctx := context.Background()
	host, port := setupMysql(ctx, t)

	dbConn = sqlinit.TestInit(username, password, host, port, database)
	return dbConn
}

func setupMysql(ctx context.Context, t *testing.T) (host, port string) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "roger",
			"MYSQL_DATABASE":      "roger",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
			wait.ForListeningPort("3306/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// perform assertions

	host, _ = container.Host(ctx)

	p, _ := container.MappedPort(ctx, "3306/tcp")
	portInt := p.Int()
	port = strconv.Itoa(portInt)
	return
}
