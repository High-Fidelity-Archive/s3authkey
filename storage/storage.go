// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package storage

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/highfidelity/s3authkey/sshkey"
)

func s3Key(k *sshkey.SshKey) string {
	return fmt.Sprintf("%d/%d.%s",
		k.ExpireBucket().Unix(),
		k.Expiration.Unix(),
		k.URLSafeFingerprint(),
	)
}

type S3Bucket struct {
	Region string
	Name   string
}

func (b S3Bucket) Upload(sshKey *sshkey.SshKey) error {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(b.Region),
	})
	if err != nil {
		log.Fatal("Fatal error ", err.Error())
		return err
	}
	uploader := s3manager.NewUploader(session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(s3Key(sshKey)),
		Body:   strings.NewReader(sshKey.OpenSSHPubKey()),
	})
	if err != nil {
		log.Fatal("Fatal error ", err.Error())
		return err
	}
	return nil
}

func (b S3Bucket) NewSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(b.Region),
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (b S3Bucket) NewDownloader() (*s3manager.Downloader, error) {
	session, err := b.NewSession()
	if err != nil {
		return nil, err
	}
	return s3manager.NewDownloader(session), nil
}

func (b S3Bucket) downloadKeys(in <-chan s3.GetObjectInput) (<-chan string, error) {
	downloader, err := b.NewDownloader()
	if err != nil {
		return nil, err
	}
	out := make(chan string)
	go func() {
		for objInput := range in {
			buf := aws.NewWriteAtBuffer([]byte{})
			_, err := downloader.Download(buf, &objInput)
			if err == nil {
				// todo: don't throw away errors
				out <- string(buf.Bytes())
			}
		}
		close(out)
	}()
	return out, nil
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (b S3Bucket) List() (<-chan string, error) {
	in := make(chan s3.GetObjectInput)

	session, err := b.NewSession()
	if err != nil {
		return nil, err
	}

	bucket := aws.String(b.Name)
	service := s3.New(session)
	resp, err := service.ListObjects(&s3.ListObjectsInput{Bucket: bucket})
	if err != nil {
		return nil, err
	}

	go func() {
		for _, item := range resp.Contents {
			// Skip files that are too big to be a key file.
			if *item.Size > 2048 {
				continue
			}
			in <- s3.GetObjectInput{Bucket: bucket, Key: aws.String(*item.Key)}
		}
		close(in)
	}()

	// TODO: don't throw away errors
	// TODO: make this configurable
	c1, _ := b.downloadKeys(in)
	c2, _ := b.downloadKeys(in)

	out := merge(c1, c2)

	return out, nil
}
