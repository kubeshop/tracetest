package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"github.com/orlangure/gnomock/preset/redis"
)

type pokeshopConfig struct {
	postgres postgresConfig
	redis    redisConfig
	rabbitmq rabbitmqConfig
	jaeger   jaegerConfig
}

type postgresConfig struct {
	host     string
	port     int
	db       string
	user     string
	password string
}

type redisConfig struct {
	endpoint string
}

type rabbitmqConfig struct {
	endpoint string
}

type jaegerConfig struct {
	host string
	port int
}

type DemoAppConfig struct {
	Endpoint string
}

func GetDemoApplicationInstance() (DemoAppConfig, error) {
	demoConfig := pokeshopConfig{}
	var wg sync.WaitGroup
	wg.Add(4)
	var err error

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		postgresConfig := postgresConfig{
			db:       "pokeshop",
			user:     "ashketchum",
			password: "squirtle123",
		}

		pgErr := getDemoPostgresInstance(postgresConfig, wg, pokeshopConfig)
		if err != nil {
			err = pgErr
		}
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		redisErr := getRedisInstance(wg, pokeshopConfig)
		if err != nil {
			err = redisErr
		}
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		rabbitErr := getRabbitMQInstance(wg, pokeshopConfig)
		if err != nil {
			err = rabbitErr
		}
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		jaegerErr := getJaegerInstance(wg, pokeshopConfig)
		if err != nil {
			err = jaegerErr
		}
	}(&wg, &demoConfig)

	wg.Wait()

	if err != nil {
		return DemoAppConfig{}, err
	}

	container, err = getPokeshopInstance(demoConfig)
	if err != nil {
		return DemoAppConfig{}, err
	}

	return DemoAppConfig{
		Endpoint: fmt.Sprintf("%s:%d", container.Host, container.DefaultPort()),
	}, nil
}

func getDemoPostgresInstance(config postgresConfig, wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) error {
	defer wg.Done()
	preset := postgres.Preset(
		postgres.WithDatabase(config.db),
		postgres.WithUser(config.user, config.password),
	)

	_, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_postgres"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return fmt.Errorf("could not start demo app's postgres: %w", err)
	}

	pokeshopConfig.postgres = postgresConfig{
		host:     "pokeshop_postgres",
		port:     5432,
		db:       config.db,
		user:     config.user,
		password: config.password,
	}

	return nil
}

func getRedisInstance(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) error {
	defer wg.Done()
	preset := redis.Preset()

	_, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_redis"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return fmt.Errorf("could not start redis container: %w", err)
	}

	pokeshopConfig.redis = redisConfig{
		endpoint: fmt.Sprintf("%s", "pokeshop_redis"),
	}

	return nil
}

func getRabbitMQInstance(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) error {
	defer wg.Done()
	preset := rabbitmq.Preset()

	_, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_rabbitmq"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return fmt.Errorf("could not start rabbitMQ container: %w", err)
	}

	pokeshopConfig.rabbitmq = rabbitmqConfig{
		endpoint: fmt.Sprintf("%s:%d", "pokeshop_rabbitmq", 5672),
	}

	return nil
}

func getJaegerInstance(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) error {
	defer wg.Done()

	_, err := gnomock.StartCustom(
		"jaegertracing/all-in-one:latest",
		gnomock.NamedPorts{
			"6832":  gnomock.Port{Protocol: "udp", Port: 6832},
			"9411":  gnomock.Port{Protocol: "TCP", Port: 9411},
			"16686": gnomock.Port{Protocol: "TCP", Port: 16686},
		},
		gnomock.WithEnv("COLLECTOR_ZIPKIN_HOST_PORT=9411"),
		gnomock.WithContainerName("pokeshop_jaeger"),
		gnomock.WithNetwork("pokeshop"),
	)

	if err != nil {
		return fmt.Errorf("could not start jaeger container: %w", err)
	}

	pokeshopConfig.jaeger = jaegerConfig{
		host: "pokeshop_jaeger",
		port: 6832,
	}

	return nil
}

func getPokeshopInstance(config pokeshopConfig) (*gnomock.Container, error) {
	databaseUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?schema=public",
		config.postgres.user,
		config.postgres.password,
		config.postgres.host,
		config.postgres.port,
		config.postgres.db,
	)

	healthCheckFunction := func(ctx context.Context, c *gnomock.Container) error {
		url := fmt.Sprintf("http://%s:%d/pokemon/healthcheck", c.Host, c.DefaultPort())
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		body, _ := ioutil.ReadAll(resp.Body)
		var x = make(map[string]interface{}, 0)
		json.Unmarshal(body, &x)

		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
		}

		return nil
	}

	container, err := gnomock.StartCustom(
		"kubeshop/demo-pokemon-api:0.0.11",
		gnomock.DefaultTCP(80),
		gnomock.WithEnv(fmt.Sprintf("DATABASE_URL=%s", databaseUrl)),
		gnomock.WithEnv(fmt.Sprintf("REDIS_URL=%s", config.redis.endpoint)),
		gnomock.WithEnv(fmt.Sprintf("RABBITMQ_HOST=%s", config.rabbitmq.endpoint)),
		gnomock.WithEnv(fmt.Sprintf("JAEGER_HOST=%s", config.jaeger.host)),
		gnomock.WithEnv(fmt.Sprintf("JAEGER_PORT=%d", config.jaeger.port)),
		gnomock.WithEnv("POKE_API_BASE_URL=https://pokeapi.co/api/v2"),
		gnomock.WithEnv("NPM_RUN_COMMAND=api"),
		gnomock.WithContainerName("pokeshop-api"),
		gnomock.WithDebugMode(),
		gnomock.WithHealthCheck(healthCheckFunction),
		gnomock.WithNetwork("pokeshop"),
	)

	if err != nil {
		return nil, fmt.Errorf("could not create demo app container: %w", err)
	}

	return container, nil
}
