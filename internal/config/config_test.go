package config

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestDefaultConfig_Server(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := rest.Config{
		BasePath:  "/v1",
		Prefix:    "",
		Port:      uint16(80),
		RateLimit: 10,
	}
	assert.Equal(expected, conf.Server)
}

func TestDefaultConfig_ApiKey(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := service.ApiConfig{
		ApiKeyValidity: 3 * time.Hour,
	}
	assert.Equal(expected, conf.ApiKey)
}

func TestDefaultConfig_Database_AssumesDockerLocalhost(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("172.17.0.1", conf.Database.Host)
}

func TestDefaultConfig_Database_AssumesPort5432(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(uint16(5432), conf.Database.Port)
}

func TestDefaultConfig_Database_AssumesPoolSizeOfOne(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(uint(1), conf.Database.ConnectionsPoolSize)
}

func TestDefaultConfig_Database_DoesNotDefineDb(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.Name)
}

func TestDefaultConfig_Database_DoesNotDefineUser(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.User)
}

func TestDefaultConfig_Database_DoesNotDefinePassword(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.Password)
}

func TestDefaultConfig_Database_AssumesDebugLogLevel(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(log.DEBUG, conf.Database.LogLevel)
}

type mockConfigurationParser struct {
	confType       string
	confPaths      []string
	confName       string
	envKeyReplacer *strings.Replacer
	envPrefix      string
	automaticEnv   bool

	readCalled int
	readErr    error

	unmarshalVal  interface{}
	unmarshalOpts []viper.DecoderConfigOption
	unmarshalErr  error
}

const defaultConfigName = "some-config"

func TestLoadConfiguration_LooksForYamlFiles(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal("yaml", m.confType)
}

func TestLoadConfiguration_LooksForFilesInExpectedPath(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal([]string{"configs"}, m.confPaths)
}

func TestLoadConfiguration_LooksForTheRightFile(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(defaultConfigName, m.confName)
}

func TestLoadConfiguration_AppliesEnvironmentOverrides(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.True(m.automaticEnv)
	assert.Equal("ENV", m.envPrefix)
	assert.Equal(strings.NewReplacer(".", "_"), m.envKeyReplacer)
}

func TestLoadConfiguration_ReadsConfiguration(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(1, m.readCalled)
}

func TestLoadConfiguration_WhenError_ExpectDefaultAndError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{
		readErr: errDefault,
	}
	configurator = &m

	conf, err := LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(errDefault, err)
	assert.Equal(DefaultConf(), conf)
}

func TestLoadConfiguration_UnmarshalsInConfiguration(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.IsType(&Configuration{}, m.unmarshalVal)
	assert.Equal(0, len(m.unmarshalOpts))
}

func TestLoadConfiguration_WhenUnmarshalFails_ExpectDefaultAndError(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{
		unmarshalErr: errDefault,
	}
	configurator = &m

	conf, err := LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(errDefault, err)
	assert.Equal(DefaultConf(), conf)
}

func resetConfigurationParser() {
	configurator = viper.New()
}

func (m *mockConfigurationParser) SetConfigType(extension string) {
	m.confType = extension
}

func (m *mockConfigurationParser) AddConfigPath(path string) {
	m.confPaths = append(m.confPaths, path)
}

func (m *mockConfigurationParser) SetConfigName(fileName string) {
	m.confName = fileName
}

func (m *mockConfigurationParser) SetEnvKeyReplacer(replacer *strings.Replacer) {
	m.envKeyReplacer = replacer
}

func (m *mockConfigurationParser) SetEnvPrefix(envPrefix string) {
	m.envPrefix = envPrefix
}

func (m *mockConfigurationParser) AutomaticEnv() {
	m.automaticEnv = true
}

func (m *mockConfigurationParser) ReadInConfig() error {
	m.readCalled++
	return m.readErr
}

func (m *mockConfigurationParser) Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	m.unmarshalVal = rawVal
	m.unmarshalOpts = append(m.unmarshalOpts, opts...)

	return m.unmarshalErr
}
