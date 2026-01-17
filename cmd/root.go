// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fido-device-onboard/go-fdo/tpm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFile    string
	configErr     error
	debug         bool
	blobPath      string
	tpmc          tpm.Closer
	tpmPath       string
	clientContext context.Context
	v             *viper.Viper
)

var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage: true,
	Use:          "fdo_client",
	Short:        "FIDO Device Onboard Client",
	Long:         `FIDO Device Onboard Client`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if configErr != nil {
			return configErr
		}

		// Sync global persistent flags from config first
		if err := syncPersistentFlagsFromConfig(cmd); err != nil {
			return err
		}

		// Sync command-specific flags based on which command is running
		switch cmd.Name() {
		case "device-init":
			if err := syncFlagsFromConfig(cmd, "device-init"); err != nil {
				return err
			}
			// Get server URL from CLI arg or config
			if len(args) > 0 {
				diURL = args[0]
			} else if v != nil && v.IsSet("device-init.server-url") {
				diURL = v.GetString("device-init.server-url")
			}
		case "onboard":
			if err := syncFlagsFromConfig(cmd, "onboard"); err != nil {
				return err
			}
		}

		return nil
	},
}

// Called by main to parse the command line and execute the subcommand
func Execute() error {
	// Catch interrupts
	var cancel context.CancelFunc
	clientContext, cancel = context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		defer signal.Stop(sigs)
		select {
		case <-clientContext.Done():
		case <-sigs:
			cancel()
		}
	}()

	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&configFile, "config", "", "Configuration file (YAML format)")
	pflags.StringVar(&blobPath, "blob", "", "File path of device credential blob")
	pflags.BoolVar(&debug, "debug", false, "Print HTTP contents")
	pflags.StringVar(&tpmPath, "tpm", "", "Use a TPM at path for device credential secrets")
	rootCmd.MarkFlagsOneRequired("blob", "tpm")
	rootCmd.MarkFlagsMutuallyExclusive("blob", "tpm")
}

// initConfig reads in config file
func initConfig() {
	v = viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)

		if err := v.ReadInConfig(); err != nil {
			configErr = fmt.Errorf("failed to read config file %q: %w", configFile, err)
			return
		}
	}
}

// syncPersistentFlagsFromConfig syncs global persistent flags from config
func syncPersistentFlagsFromConfig(cmd *cobra.Command) error {
	return syncFlagsFromViper(cmd.Root().PersistentFlags(), v)
}

// syncFlagsFromConfig syncs command-specific flags from a config section
func syncFlagsFromConfig(cmd *cobra.Command, section string) error {
	if v == nil {
		return nil
	}
	sub := v.Sub(section)
	if sub == nil {
		return nil
	}
	return syncFlagsFromViper(cmd.Flags(), sub)
}

// syncFlagsFromViper syncs flags from a viper instance
func syncFlagsFromViper(flags *pflag.FlagSet, vp *viper.Viper) error {
	if vp == nil {
		return nil
	}

	var syncErr error

	flags.VisitAll(func(flag *pflag.Flag) {
		// Skip if flag was explicitly set via CLI
		if flag.Changed {
			return
		}

		if !vp.IsSet(flag.Name) {
			return
		}

		configValue := vp.Get(flag.Name)
		if configValue == nil {
			return
		}

		// Set flag value from config using FlagSet.Set to mark as Changed
		var err error
		switch val := configValue.(type) {
		case []interface{}:
			for _, item := range val {
				if setErr := flags.Set(flag.Name, fmt.Sprintf("%v", item)); setErr != nil {
					err = setErr
					break
				}
			}
		case []string:
			for _, item := range val {
				if setErr := flags.Set(flag.Name, item); setErr != nil {
					err = setErr
					break
				}
			}
		default:
			err = flags.Set(flag.Name, fmt.Sprintf("%v", configValue))
		}

		if err != nil && syncErr == nil {
			syncErr = fmt.Errorf("failed to set flag %s from config: %w", flag.Name, err)
		}
	})

	return syncErr
}
