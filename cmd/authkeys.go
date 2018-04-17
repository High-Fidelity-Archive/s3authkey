// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"fmt"

	"github.com/highfidelity/s3authkey/storage"
	"github.com/spf13/cobra"
)

// authkeysCmd represents the authkeys command
var authkeysCmd = &cobra.Command{
	Use:   "authkeys",
	Short: "List available auth keys",
	Long: `List available auth keys in a format compatible with the
AuthorizedKeysCommand sshd_config settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		s3bucket := storage.S3Bucket{Region: regionName, Name: bucketName}
		// todo: handle this error
		pubkeys, _ := s3bucket.List()
		for pubkey := range pubkeys {
			fmt.Printf(pubkey)
		}
	},
}

func init() {
	rootCmd.AddCommand(authkeysCmd)
}
