package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKoanf(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	maxIterations := 10

	testCases := []struct {
		configFile     string
		args           []string
		envs           map[string]string
		expectedValues map[string]interface{}
	}{
		// Test: Default value of pflag. Does not work or is not supported?
		// {
		// 	configFile: "",
		// 	envs:       map[string]string{},
		// 	args:       []string{},
		// 	expectedValues: map[string]interface{}{
		// 		"log.level": "info",
		// 	},
		// },
		// Test: hierarchy/order
		{
			configFile: "config.yaml",
			envs: map[string]string{
				"DEMO_LOG_LEVEL": "fatal",
			},
			args: []string{"--log-level=debug"},
			expectedValues: map[string]interface{}{
				"log.level": "debug",
			},
		},
	}

	var testResults map[int]map[string]int = make(map[int]map[string]int)

	for i, testCase := range testCases {
		// set env
		for key, value := range testCase.envs {
			t.Setenv(key, value)
		}

		if _, present := testResults[i]; !present {
			testResults[i] = make(map[string]int)
		}

		for j := 0; j < maxIterations; j++ {
			pflagSet := getRootPflagSet()
			pflagSet.Parse(testCase.args)

			kConfig, err := prepareKConfig(testCase.configFile, pflagSet)
			require.NoError(err)

			if _, present := testResults[i]; !present {
				testResults[i] = make(map[string]int)
			}

			for expectedKey, expectedValue := range testCase.expectedValues {
				if !assert.True(kConfig.Exists(expectedKey), "Key %s not found", expectedKey) {
					testResults[i]["failed"]++
					continue
				}
				if !assert.Equal(expectedValue, kConfig.All()[expectedKey]) {
					testResults[i]["failed"]++
					continue
				}
				testResults[i]["successful"]++
			}
		}

		// unset env
		for key := range testCase.envs {
			os.Unsetenv(key)
		}
	}

	// results
	for testCase, testCaseResults := range testResults {
		t.Logf("TestCase %v, failed %v of %v iterations", testCase, testCaseResults["failed"], maxIterations)
		t.Logf("TestCase %v, was successful %v of %v iterations", testCase, testCaseResults["successful"], maxIterations)

	}
}
