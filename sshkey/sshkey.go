// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func urlSafeb64(s string) string {
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	return s
}

type SshKey struct {
	Expiration time.Time
	privateKey *rsa.PrivateKey
	publicKey  *ssh.PublicKey
}

func NewSshKey(age time.Duration) (*SshKey, error) {
	// Generate a new private key.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("couldn't generate private key")
		return nil, err
	}

	pubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	return &SshKey{
		Expiration: time.Now().Add(age),
		privateKey: key,
		publicKey:  &pubKey,
	}, nil
}

func (k *SshKey) PEM() string {
	block := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(k.privateKey),
	}
	return string(pem.EncodeToMemory(&block))
}

func (k SshKey) OpenSSHPubKey() string {
	return string(ssh.MarshalAuthorizedKey(*k.publicKey))
}

func (k *SshKey) Fingerprint() string {
	return ssh.FingerprintSHA256(*k.publicKey)
}

func (k *SshKey) URLSafeFingerprint() string {
	return strings.Replace(urlSafeb64(k.Fingerprint()), ":", ".", 1)
}

func (k *SshKey) ExpireBucket() time.Time {
	return k.Expiration.Round(time.Hour)
}
