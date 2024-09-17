package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	ENV_PREFIX string = "DEMO_"
)

func main() {
	var (
		kConfig *koanf.Koanf = nil
	)

	rootCmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			if kConfig != nil {
				b, err := yaml.Parser().Marshal(kConfig.All())
				if err != nil {
					return err
				}
				fmt.Fprint(os.Stdout, string(b))
			}

			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			// Especially paring flags for koanf
			err = cmd.Flags().Parse(args)
			if err != nil {
				return err
			}

			kConfig, err = prepareKConfig(configFile, cmd.Flags())
			if err != nil {
				return err
			}
			return nil
		},
		Use: "demo",
	}
	rootCmd.Flags().AddFlagSet(getRootPflagSet())

	rootCmd.Execute()
}

func getRootPflagSet() *pflag.FlagSet {
	rootPflagSet := pflag.NewFlagSet("root", pflag.ContinueOnError)
	rootPflagSet.String("config", "config.yaml", "Configuration file")
	rootPflagSet.String("log-level", "info", "Log level")
	return rootPflagSet
}

// prepareKConfig loads sources in following order
// 1. config file
// 2. env
// 3. pflagSet (no defaults)
func prepareKConfig(configFile string, parsedPflags *pflag.FlagSet) (*koanf.Koanf, error) {
	kConfig := koanf.New(".")

	// 1. config file

	if _, err := os.Stat(configFile); !errors.Is(err, os.ErrNotExist) {
		err = kConfig.Load(file.Provider(configFile), kyaml.Parser())
		if err != nil {
			return nil, err
		}
	}

	// 2. env vars
	err := kConfig.Load(env.Provider(ENV_PREFIX, "_", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, ENV_PREFIX)), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, err
	}

	// 3. pflag (hopefully only defined flags
	err = kConfig.Load(posflag.Provider(parsedPflags, "-", nil), nil)
	if err != nil {
		return nil, err
	}

	return kConfig, nil
}

// func customMergeFunc(src, dest map[string]interface{}) error {
// 	for srcKey, srcValue := range src {
// 		switch srcKey {
// 		// case "log-level":
// 		// 	dest["log"].(map[string]interface{})["level"] = srcValue
// 		default:
// 			dest[srcKey] = srcValue
// 		}
// 	}
// 	return nil
// }
