// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"fmt"
	"time"

	"github.com/highfidelity/s3authkey/sshkey"
	"github.com/highfidelity/s3authkey/storage"
	"github.com/spf13/cobra"
)

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate and authorize an SSH key",
	Long:  "Generate and authorize an SSH key",
	Run: func(cmd *cobra.Command, args []string) {
		s3bucket := storage.S3Bucket{Region: regionName, Name: bucketName}
		key, err := sshkey.NewSshKey(time.Duration(6) * time.Hour)
		if err != nil {
			// TODO: handle this error
			return
		}
		err = s3bucket.Upload(key)
		if err != nil {
			// TODO: handle this error
			return
		}
		fmt.Println(key.PEM())
	},
}

func init() {
	rootCmd.AddCommand(genkeyCmd)
}
