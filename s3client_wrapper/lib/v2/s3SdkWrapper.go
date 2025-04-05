package s3SdkV2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
)

//参考： https://docs.aws.amazon.com/zh_cn/sdk-for-go/v2/developer-guide/go_s3_code_examples.html
//参考： https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2

type S3ClientConfig struct {
	Bucket     string // bucket name
	Region     string
	Ak         string
	Sk         string
	StoreType  int64
	Endpoint   string
	FilePrefix string
	//
	CliTimeoutSecondCfg int32
}
type S3ClientBase struct {
	S3Client *s3.Client
	CliCfg   *S3ClientConfig
}

func (sc *S3ClientBase) Upload(data io.ReadSeeker, remotePath string) error {
	ctx := context.Background()
	var cancelFn context.CancelFunc = nil

	if sc.CliCfg.CliTimeoutSecondCfg > 0 {
		ctx, cancelFn = context.WithTimeout(ctx,
			time.Duration(sc.CliCfg.CliTimeoutSecondCfg)*time.Millisecond)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	_, err := sc.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(sc.CliCfg.Bucket),
		Key:    aws.String(remotePath),
		Body:   data,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", sc.CliCfg.Bucket)
		} else {
			log.Printf("Couldn't upload data to %v:%v. Here's why: %v\n",
				sc.CliCfg.Bucket, remotePath, err)
		}
		//
		return err
	}

	// err = s3.NewObjectExistsWaiter(sc.S3Client).Wait(ctx,
	// 	&s3.HeadObjectInput{Bucket: aws.String(sc.CliCfg.Bucket), Key: aws.String(remotePath)}, time.Minute)
	// if err != nil {
	// 	log.Printf("Failed attempt to wait for object %s to exist.\n", remotePath)
	// 	return err
	// }

	return nil
}

func (sc *S3ClientBase) DownloadFile(remotePath string) ([]byte, error) {
	ctx := context.Background()

	retGetObj, err := sc.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(sc.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	})

	if err != nil {
		log.Printf("get object fail, err: %v, obj key: %", err, remotePath)
		return nil, err
	}
	if retGetObj == nil {
		log.Printf("get obj is nil, obj key: %v", remotePath)
		return nil, fmt.Errorf("obj ret is nil")
	}
	defer retGetObj.Body.Close()
	//
	retBuf, err := io.ReadAll(retGetObj.Body)
	if err != nil {
		log.Printf("read data from get obj response fail, err: %v", err)
		return nil, err
	}

	return retBuf, nil
}

type S3ClientV2 struct {
	S3ClientBase
}

func NewS3ClientV2(cfg *S3ClientConfig) *S3ClientV2 {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(
		cfg.Region,
	), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
		cfg.Ak, cfg.Sk, "",
	)))
	if err != nil {
		log.Printf("load default config fail, err: %v", err)
		return nil
	}
	cli := s3.NewFromConfig(sdkConfig)

	retCli := &S3ClientV2{
		S3ClientBase: S3ClientBase{
			S3Client: cli,
			CliCfg:   cfg,
		},
	}
	return retCli
}
