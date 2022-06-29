package testmock

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

type DemoApp struct {
	endpoint          string
	apiContainer      *gnomock.Container
	workerContainer   *gnomock.Container
	redisContainer    *gnomock.Container
	postgresContainer *gnomock.Container
	jaegerContainer   *gnomock.Container
	rabbitmqContainer *gnomock.Container
}

func (da *DemoApp) Endpoint() string {
	return da.endpoint
}

func (da *DemoApp) JaegerEndpoint() string {
	return fmt.Sprintf("%s:%d", da.jaegerContainer.Host, da.jaegerContainer.Port("16685"))
}

func (da *DemoApp) Stop() {
	gnomock.Stop(da.apiContainer)
	gnomock.Stop(da.workerContainer)
	gnomock.Stop(da.jaegerContainer)
	gnomock.Stop(da.postgresContainer)
	gnomock.Stop(da.redisContainer)
	gnomock.Stop(da.rabbitmqContainer)
}

func GetDemoApplicationInstance() (*DemoApp, error) {
	demoConfig := pokeshopConfig{}
	var wg sync.WaitGroup
	var jaegerContainer, postgresContainer, redisContainer, rabbitContainer *gnomock.Container
	wg.Add(4)
	var err error

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		defer wg.Done()
		postgresConfig := postgresConfig{
			db:       "pokeshop",
			user:     "ashketchum",
			password: "squirtle123",
		}

		container, pgErr := getDemoPostgresInstance(postgresConfig, pokeshopConfig)
		if pgErr != nil {
			err = pgErr
		}
		postgresContainer = container
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		defer wg.Done()
		container, redisErr := getRedisInstance(pokeshopConfig)
		if redisErr != nil {
			err = redisErr
		}
		redisContainer = container
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		defer wg.Done()
		container, rabbitErr := getRabbitMQInstance(pokeshopConfig)
		if rabbitErr != nil {
			err = rabbitErr
		}
		rabbitContainer = container
	}(&wg, &demoConfig)

	go func(wg *sync.WaitGroup, pokeshopConfig *pokeshopConfig) {
		defer wg.Done()
		container, jaegerErr := getJaegerInstance(pokeshopConfig)
		if jaegerErr != nil {
			err = jaegerErr
		}
		jaegerContainer = container
	}(&wg, &demoConfig)

	wg.Wait()

	if err != nil {
		return nil, err
	}

	apiContainer, workerContainer, err := getPokeshopInstance(demoConfig)
	if err != nil {
		return nil, err
	}

	return &DemoApp{
		endpoint:          fmt.Sprintf("%s:%d", apiContainer.Host, apiContainer.DefaultPort()),
		apiContainer:      apiContainer,
		workerContainer:   workerContainer,
		redisContainer:    redisContainer,
		postgresContainer: postgresContainer,
		jaegerContainer:   jaegerContainer,
		rabbitmqContainer: rabbitContainer,
	}, nil
}

func getDemoPostgresInstance(config postgresConfig, pokeshopConfig *pokeshopConfig) (*gnomock.Container, error) {
	preset := postgres.Preset(
		postgres.WithDatabase(config.db),
		postgres.WithUser(config.user, config.password),
	)

	container, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_postgres"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start demo app's postgres: %w", err)
	}

	pokeshopConfig.postgres = postgresConfig{
		host:     "pokeshop_postgres",
		port:     5432,
		db:       config.db,
		user:     config.user,
		password: config.password,
	}

	return container, nil
}

func getRedisInstance(pokeshopConfig *pokeshopConfig) (*gnomock.Container, error) {
	preset := redis.Preset()

	container, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_redis"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start redis container: %w", err)
	}

	pokeshopConfig.redis = redisConfig{
		endpoint: "pokeshop_redis",
	}

	return container, nil
}

func getRabbitMQInstance(pokeshopConfig *pokeshopConfig) (*gnomock.Container, error) {
	preset := rabbitmq.Preset()

	container, err := gnomock.Start(
		preset,
		gnomock.WithContainerName("pokeshop_rabbitmq"),
		gnomock.WithNetwork("pokeshop"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start rabbitMQ container: %w", err)
	}

	pokeshopConfig.rabbitmq = rabbitmqConfig{
		endpoint: fmt.Sprintf("%s:%d", "pokeshop_rabbitmq", 5672),
	}

	return container, nil
}

func getJaegerInstance(pokeshopConfig *pokeshopConfig) (*gnomock.Container, error) {
	container, err := gnomock.StartCustom(
		"jaegertracing/all-in-one:latest",
		gnomock.NamedPorts{
			"6832":  gnomock.Port{Protocol: "udp", Port: 6832},
			"9411":  gnomock.Port{Protocol: "TCP", Port: 9411},
			"16685": gnomock.Port{Protocol: "TCP", Port: 16685},
			"16686": gnomock.Port{Protocol: "TCP", Port: 16686},
		},
		gnomock.WithEnv("COLLECTOR_ZIPKIN_HOST_PORT=9411"),
		gnomock.WithContainerName("pokeshop_jaeger"),
		gnomock.WithNetwork("pokeshop"),
	)

	if err != nil {
		return nil, fmt.Errorf("could not start jaeger container: %w", err)
	}

	pokeshopConfig.jaeger = jaegerConfig{
		host: "pokeshop_jaeger",
		port: 6832,
	}

	return container, nil
}

func getPokeshopInstance(config pokeshopConfig) (*gnomock.Container, *gnomock.Container, error) {
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

	apiContainer, err := gnomock.StartCustom(
		"kubeshop/demo-pokemon-api:latest",
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
		return nil, nil, fmt.Errorf("could not create demo app api container: %w", err)
	}

	workerContainer, err := gnomock.StartCustom(
		"kubeshop/demo-pokemon-api:latest",
		gnomock.DefaultTCP(80),
		gnomock.WithEnv(fmt.Sprintf("DATABASE_URL=%s", databaseUrl)),
		gnomock.WithEnv(fmt.Sprintf("REDIS_URL=%s", config.redis.endpoint)),
		gnomock.WithEnv(fmt.Sprintf("RABBITMQ_HOST=%s", config.rabbitmq.endpoint)),
		gnomock.WithEnv(fmt.Sprintf("JAEGER_HOST=%s", config.jaeger.host)),
		gnomock.WithEnv(fmt.Sprintf("JAEGER_PORT=%d", config.jaeger.port)),
		gnomock.WithEnv("POKE_API_BASE_URL=https://pokeapi.co/api/v2"),
		gnomock.WithEnv("NPM_RUN_COMMAND=worker"),
		gnomock.WithContainerName("pokeshop-worker"),
		gnomock.WithDebugMode(),
		gnomock.WithNetwork("pokeshop"),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("could not create demo app worker container: %w", err)
	}

	return apiContainer, workerContainer, nil
}
