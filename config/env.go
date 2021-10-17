package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/joho/godotenv"
)

// EnvType is variable in .env file
type EnvType struct {
	DbConn           string
	DbConnMock       string
	HTTPPort         string
	ArtemisConn      string
	ArtemisLogin     string
	ArtemisPass      string
	ArtemisSkip      map[string]struct{}
	ArtemisOptions   []func(*stomp.Conn) error
	ArtemisQueueName string
}

// Env is global var for EnvType
var Env = EnvType{}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	Env.DbConn = getEnv("DB_CONN", "default conn")
	Env.DbConnMock = getEnv("DB_CONN_MOCK", "default conn")
	Env.HTTPPort = getEnv("HTTP_PORT", "3000")
	Env.ArtemisConn = getEnv("ARTEMIS_CONN", "localhost:61613")
	Env.ArtemisLogin = getEnv("ARTEMIS_LOGIN", "admin")
	Env.ArtemisPass = getEnv("ARTEMIS_PASS", "admin")
	artemisSkipString := getEnv("ARTEMIS_SKIP", "")
	if artemisSkipString != "" {
		artemisSkipArray := strings.Split(artemisSkipString, ";")
		Env.ArtemisSkip = make(map[string]struct{})
		for _, skipTable := range artemisSkipArray {
			if skipTable != "" {
				Env.ArtemisSkip[skipTable] = struct{}{}
			}
		}
	}

	Env.ArtemisOptions = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(Env.ArtemisLogin, Env.ArtemisPass),
		stomp.ConnOpt.Host("/"),
		stomp.ConnOpt.HeartBeat(time.Second*120, time.Second*120),
	}
	Env.ArtemisQueueName = getEnv("ARTEMIS_QUEUE_NAME", "/queue/audit")
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valueStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valueStr); err == nil {
		return val
	}

	return defaultValue
}

func getEnvAsSlice(name string, defaultValue []string, separator string) []string {
	valueStr := getEnv(name, "")

	if valueStr == "" {
		return defaultValue
	}

	val := strings.Split(valueStr, separator)

	return val
}
