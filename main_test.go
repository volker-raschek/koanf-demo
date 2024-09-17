package main

import (
	"fmt"
	"os"
	"testing"

	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/stretchr/testify/require"
)

func TestKoanf(t *testing.T) {
	require := require.New(t)

	testCases := []struct {
		configFile     string
		args           []string
		envs           map[string]string
		expectedValues map[string]interface{}
	}{
		// Test: env
		{
			configFile:     "",
			envs:           map[string]string{"DEMO_LOG_LEVEL": "fatal"},
			args:           []string{},
			expectedValues: map[string]interface{}{"log.level": "fatal"},
		},
		// Test: flag
		{
			configFile:     "",
			envs:           map[string]string{},
			args:           []string{"--log-level=fatal"},
			expectedValues: map[string]interface{}{"log.level": "fatal"},
		},
		// Test: env + flag
		{
			configFile:     "",
			envs:           map[string]string{"DEMO_LOG_LEVEL": "fatal"},
			args:           []string{"--log-level=error"},
			expectedValues: map[string]interface{}{"log.level": "error"},
		},
		// Test: hierarchy/order
		{
			configFile:     "config.yaml",
			envs:           map[string]string{"DEMO_LOG_LEVEL": "fatal"},
			args:           []string{"--log-level=debug"},
			expectedValues: map[string]interface{}{"log.level": "debug"},
		},
	}

	for _, testCase := range testCases {
		// set env
		for key, value := range testCase.envs {
			t.Setenv(key, value)
		}

		pflagSet := getRootPflagSet()
		pflagSet.Parse(testCase.args)

		kConfig, err := prepareKConfig(testCase.configFile, pflagSet)
		require.NoError(err)

		b, err := kyaml.Parser().Marshal(kConfig.All())
		require.NoError(err)
		t.Logf(fmt.Sprintf("%s", string(b)))

		require.Equal(testCase.expectedValues, kConfig.All())

		// unset env
		for key := range testCase.envs {
			os.Unsetenv(key)
		}
	}
}
