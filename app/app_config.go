package app

import (
	"fmt"
	"github.com/jitsucom/bulker/base/uuid"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

type AppConfig struct {
	ClusterId  string `mapstructure:"CLUSTER_ID" default:"bulker_cluster"`
	InstanceId string `mapstructure:"INSTANCE_ID"`

	HTTPPort int `mapstructure:"HTTP_PORT" default:"3042"`

	AuthTokens   string `mapstructure:"AUTH_TOKENS"`
	TokenSecrets string `mapstructure:"TOKEN_SECRETS"`

	//TODO: CONFIGURATION_SOURCE
	ConfigSource string `mapstructure:"CONFIG_SOURCE"`

	KafkaBootstrapServers string `mapstructure:"KAFKA_BOOTSTRAP_SERVERS" default:"localhost:9092"`
	KafkaSecurityProtocol string `mapstructure:"KAFKA_SECURITY_PROTOCOL" default:"PLAINTEXT"`
	KafkaSaslMechanism    string `mapstructure:"KAFKA_SASL_MECHANISM" default:"GSSAPI"`
	KafkaSaslUsername     string `mapstructure:"KAFKA_SASL_USERNAME"`
	KafkaSaslPassword     string `mapstructure:"KAFKA_SASL_PASSWORD"`

	//TODO: Change to hours
	KafkaTopicRetentionMs       int `mapstructure:"KAFKA_TOPIC_RETENTION_MS" default:"604800000"`
	KafkaTopicPartitionsCount   int `mapstructure:"KAFKA_TOPIC_PARTITIONS_COUNT" default:"4"`
	KafkaTopicReplicationFactor int `mapstructure:"KAFKA_TOPIC_REPLICATION_FACTOR" default:"1"`
	KafkaAdminMetadataTimeoutMs int `mapstructure:"KAFKA_ADMIN_METADATA_TIMEOUT_MS" default:"1000"`
	//TODO: max.poll.interval.ms

	TopicManagerRefreshPeriodSec int `mapstructure:"TOPIC_MANAGER_REFRESH_PERIOD_SEC" default:"5"`

	ProducerWaitForDeliveryMs int `mapstructure:"PRODUCER_WAIT_FOR_DELIVERY_MS" default:"1000"`

	//TODO: per destination batch period
	BatchRunnerPeriodSec          int `mapstructure:"BATCH_RUNNER_PERIOD_SEC" default:"300"`
	BatchRunnerDefaultBatchSize   int `mapstructure:"BATCH_RUNNER_DEFAULT_BATCH_SIZE" default:"1000"`
	BatchRunnerWaitForMessagesSec int `mapstructure:"BATCH_RUNNER_WAIT_FOR_MESSAGES_SEC" default:"1"`
}

func init() {
	initViperVariables()
	viper.SetDefault("INSTANCE_ID", uuid.New())
}

func initViperVariables() {
	elem := reflect.TypeOf(AppConfig{})
	fieldsCount := elem.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := elem.Field(i)
		variable := field.Tag.Get("mapstructure")
		if variable != "" {
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				viper.SetDefault(variable, defaultValue)
			} else {
				_ = viper.BindEnv(variable)
			}
		}
	}
}

func InitAppConfig() (*AppConfig, error) {
	appConfig := AppConfig{}
	configPath := os.Getenv("BULKER_CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("bulker")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("BULKER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		//it is ok to not have config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("❗error reading config file: %s", err)
		}
	}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, fmt.Errorf("❗error unmarshalling config: %s", err)
	}
	return &appConfig, nil
}