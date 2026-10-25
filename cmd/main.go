package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"generator"
	"generator/internal/logger/zaplog"
	"generator/pkg/fileutils"
	"generator/storage"
)

// Error is a default error type for generator cli.
var Error = errs.Class("generator cli error")

// Config contains configurable values for generator project.
type Config struct {
	Database         string `json:"database"`
	generator.Config `json:"config"`
}

// commands.
var (
	// generator root cmd.
	rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "cli for interacting with generator project",
	}

	// generator setup cmd.
	setupCmd = &cobra.Command{
		Use:         "setup",
		Short:       "setups the program config",
		RunE:        cmdSetup,
		Annotations: map[string]string{"type": "setup"},
	}
	runCmd = &cobra.Command{
		Use:         "run",
		Short:       "runs the program",
		RunE:        cmdRun,
		Annotations: map[string]string{"type": "run"},
	}
	destroyCmd = &cobra.Command{
		Use:         "destroy",
		Short:       "deletes config folder",
		RunE:        cmdDestroy,
		Annotations: map[string]string{"type": "run"},
	}
	setupCfg Config
	runCfg   Config

	configPath = fileutils.ApplicationDir("generator")
)

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func cmdSetup(cmd *cobra.Command, args []string) (err error) {
	log := zaplog.NewLog()

	setupDir, err := filepath.Abs(configPath)
	if err != nil {
		return Error.Wrap(err)
	}

	err = os.MkdirAll(setupDir, os.ModePerm)
	if err != nil {
		return Error.Wrap(err)
	}

	configFile, err := os.Create(path.Join(configPath, "/config.json"))
	if err != nil {
		log.Error("could not create config file", Error.Wrap(err))
		return Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, configFile.Close())
	}()

	jsonData, err := json.MarshalIndent(setupCfg, "", "    ")
	if err != nil {
		log.Error("could not marshal config", Error.Wrap(err))
		return Error.Wrap(err)
	}

	_, err = configFile.Write(jsonData)
	if err != nil {
		log.Error("could not write to config", Error.Wrap(err))
		return Error.Wrap(err)
	}

	return nil
}

func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	log := zaplog.NewLog()

	runCfg, err = readConfig()
	if err != nil {
		log.Error("Could not read config from default place", Error.Wrap(err))
		return Error.Wrap(err)
	}

	db, err := storage.New(runCfg.Database)
	if err != nil {
		log.Error("Error starting master database on generator bank service", Error.Wrap(err))
		return Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	// TODO: remove for production.
	err = db.CreateSchema(ctx)
	if err != nil {
		log.Error("Error creating schema", Error.Wrap(err))
	}

	generator, err := generator.New(log, runCfg.Config, db)
	if err != nil {
		log.Error("Error starting generator bank service", Error.Wrap(err))
		return Error.Wrap(err)
	}

	runError := generator.Run(ctx)
	closeError := generator.Close()

	return Error.Wrap(errs.Combine(runError, closeError))
}

func cmdDestroy(cmd *cobra.Command, args []string) (err error) {
	return os.RemoveAll(configPath)
}

// readConfig reads config from default config dir.
// ${VAR} and ${VAR:-default} in the file are replaced with environment variables.
func readConfig() (config Config, err error) {
	configBytes, err := ioutil.ReadFile(path.Join(configPath, "/config.json"))
	if err != nil {
		return Config{}, err
	}

	return config, json.Unmarshal([]byte(expandEnv(string(configBytes))), &config)
}

// expandEnv substitutes ${VAR} with the variable value and ${VAR:-default}
// with the default value when the variable is unset or empty.
func expandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		name, fallback := key, ""
		if i := strings.Index(key, ":-"); i >= 0 {
			name, fallback = key[:i], key[i+2:]
		}
		if value := os.Getenv(name); value != "" {
			return value
		}
		return fallback
	})
}
