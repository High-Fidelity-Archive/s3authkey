// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package sshkey

import (
	"strings"
	"testing"
	"time"
)

func TestNewSshKey(t *testing.T) {
	variance := 3
	age := time.Duration(6) * time.Hour
	sshKey, _ := NewSshKey(age)
	// Ensure Expiration is set correctly
	expectedExpiration := time.Now().Add(age)
	deltaTime := int(sshKey.Expiration.Unix() - expectedExpiration.Unix())
	if deltaTime > variance {
		t.Errorf(
			"Expiration delta out of variance. Expected: %d or less, got: %d ",
			variance,
			deltaTime,
		)
	}
}

func TestPEM(t *testing.T) {
	sshKey, _ := NewSshKey(0)
	encodedKey := sshKey.PEM()
	// Quick and dirty check to make sure string looks like a key. x509 does the
	// heavy lifting, so I'm really counting on it to pull through here.
	l := len(encodedKey)
	if !(1500 < l && l < 1700) {
		t.Errorf("unexpected private key size %d", l)
	}
	if !strings.Contains(encodedKey, "RSA PRIVATE KEY") {
		t.Errorf("key doesn't container type marker")
	}
}

func TestOpenSSHPubKey(t *testing.T) {
	sshKey, _ := NewSshKey(0)
	openSsh := sshKey.OpenSSHPubKey()
	l := len(openSsh)
	if !(375 < l && l < 390) {
		t.Errorf("unexpected openssh pubkey size %d", l)
	}
	if !strings.HasPrefix(openSsh, "ssh-rsa") {
		t.Errorf("key doesn't start with ssh-rsa")
	}
}

func TestURLSafeFingerprint(t *testing.T) {
	sshKey, _ := NewSshKey(0)
	fp := sshKey.URLSafeFingerprint()
	if len(fp) != 50 {
		t.Errorf("fingerprint isn't 50 characters")
	}
	if !strings.HasPrefix(fp, "SHA256.") {
		t.Errorf("fingerprint doesn't start with 'SHA256.': %s", fp)
	}
}

func TestExpireBucket(t *testing.T) {
	sshKey, _ := NewSshKey(0)
	if sshKey.ExpireBucket() != sshKey.Expiration.Round(time.Hour) {
		t.Errorf("expire bucket isn't rounded to the hour")
	}
}
