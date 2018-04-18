// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, bucketName, regionName string

var rootCmd = &cobra.Command{
	Use:   "s3authkey",
	Short: "Leverage AWS's IAM and S3 to authorize ephemeral SSH keys.",
	Long:  "Leverage AWS's IAM and S3 to authorize ephemeral SSH keys.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3authkey.yaml)")
	rootCmd.PersistentFlags().StringVarP(&regionName, "region", "r", "", "AWS region (required)")
	rootCmd.PersistentFlags().StringVarP(&bucketName, "bucket", "b", "", "AWS bucket (required)")
	rootCmd.MarkFlagRequired("region")
	rootCmd.MarkFlagRequired("bucket")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".s3authkey" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".s3authkey")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
