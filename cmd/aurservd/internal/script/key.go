// Copyright (C) liasica. 2024-present.
//
// Created at 2024-03-02
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func keyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "密钥工具",
	}
	cmd.AddCommand(genrsaCmd())
	return cmd
}

func genrsaCmd() *cobra.Command {
	var (
		bits int
		path string
	)

	cmd := &cobra.Command{
		Use:   "genrsa",
		Short: "生成 RSA 密钥对",
		Run: func(cmd *cobra.Command, args []string) {
			privateKey, err := rsa.GenerateKey(rand.Reader, bits)
			if err != nil {
				fmt.Printf("私钥生成失败: %s\n", err)
				os.Exit(1)
			}

			// Encode the private key to the PEM format
			privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
			privateKeyPEM := &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: privateBytes,
			}
			privateKeyFile, err := os.Create(fmt.Sprintf("%s_private_key.pem", path))
			if err != nil {
				fmt.Printf("私钥保存失败: %s\n", err)
				os.Exit(1)
			}
			_ = pem.Encode(privateKeyFile, privateKeyPEM)
			_ = privateKeyFile.Close()

			// Extract the public key from the private key
			publicKey := &privateKey.PublicKey

			// Encode the public key to the PEM format
			publicBytes := x509.MarshalPKCS1PublicKey(publicKey)
			publicKeyPEM := &pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: publicBytes,
			}
			publicKeyFile, err := os.Create(fmt.Sprintf("%s_public_key.pem", path))
			if err != nil {
				fmt.Printf("公钥保存失败: %s\n", err)
				os.Exit(1)
			}
			_ = pem.Encode(publicKeyFile, publicKeyPEM)
			_ = publicKeyFile.Close()

			fmt.Printf("RSA 密钥对生成成功\nPRIVATE: %s\nPUBLIC: %s\n", base64.StdEncoding.EncodeToString(privateBytes), base64.StdEncoding.EncodeToString(publicBytes))
		},
	}

	cmd.Flags().IntVarP(&bits, "bits", "b", 4096, "密钥长度")
	cmd.Flags().StringVarP(&path, "path", "p", "", "密钥路径名称, 例如: /tmp/rsa/[name]")

	_ = cmd.MarkFlagRequired("path")
	return cmd
}
