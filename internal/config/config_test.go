package config

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/server"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestUnit_DefaultConfig_Server(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := server.Config{
		BasePath: "/v1",
		Port:     uint16(80),
	}
	assert.Equal(expected, conf.Server)
}

func TestUnit_DefaultConfig_ApiKey(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	expected := service.ApiConfig{
		ApiKeyValidity: 3 * time.Hour,
	}
	assert.Equal(expected, conf.ApiKey)
}

func TestUnit_DefaultConfig_Database_AssumesDockerLocalhost(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("172.17.0.1", conf.Database.Host)
}

func TestUnit_DefaultConfig_Database_AssumesPort5432(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(uint16(5432), conf.Database.Port)
}

func TestUnit_DefaultConfig_Database_AssumesPoolSizeOfOne(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(uint(1), conf.Database.ConnectionsPoolSize)
}

func TestUnit_DefaultConfig_Database_AssumesTwoSecondsConnectTimeout(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal(2*time.Second, conf.Database.ConnectTimeout)
}

func TestUnit_DefaultConfig_Database_DoesNotDefineDb(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.Name)
}

func TestUnit_DefaultConfig_Database_DoesNotDefineUser(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.User)
}

func TestUnit_DefaultConfig_Database_DoesNotDefinePassword(t *testing.T) {
	assert := assert.New(t)

	conf := DefaultConf()

	assert.Equal("", conf.Database.Password)
}

func TestUnit_DefaultConfig_Database_AssumesDebugLogLevel(t *testing.T) {
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

func TestUnit_LoadConfiguration_LooksForYamlFiles(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal("yaml", m.confType)
}

func TestUnit_LoadConfiguration_LooksForFilesInExpectedPath(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal([]string{"configs"}, m.confPaths)
}

func TestUnit_LoadConfiguration_LooksForTheRightFile(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(defaultConfigName, m.confName)
}

func TestUnit_LoadConfiguration_AppliesEnvironmentOverrides(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.True(m.automaticEnv)
	assert.Equal("ENV", m.envPrefix)
	assert.Equal(strings.NewReplacer(".", "_"), m.envKeyReplacer)
}

func TestUnit_LoadConfiguration_ReadsConfiguration(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.Equal(1, m.readCalled)
}

func TestUnit_LoadConfiguration_WhenError_ExpectDefaultAndError(t *testing.T) {
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

func TestUnit_LoadConfiguration_UnmarshalsInConfiguration(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConfigurationParser)

	m := mockConfigurationParser{}
	configurator = &m

	LoadConfiguration(defaultConfigName, DefaultConf())

	assert.IsType(&Configuration{}, m.unmarshalVal)
	assert.Equal(0, len(m.unmarshalOpts))
}

func TestUnit_LoadConfiguration_WhenUnmarshalFails_ExpectDefaultAndError(t *testing.T) {
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
