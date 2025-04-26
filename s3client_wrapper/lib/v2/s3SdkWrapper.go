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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
)

//参考： https://docs.aws.amazon.com/zh_cn/sdk-for-go/v2/developer-guide/go_s3_code_examples.html
//参考： https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2

// 源码参考：github.com/awsdocs/aws-doc-sdk-examples/gov2/s3

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
	ObjectPrivilege int
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

	privilege, ok := privilegeInS3Inner[sc.CliCfg.ObjectPrivilege]
	if !ok {
		privilege = types.ObjectCannedACLPrivate
	}
	_, err := sc.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(sc.CliCfg.Bucket),
		Key:    aws.String(remotePath),
		Body:   data,
		ACL: privilege,
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
		log.Printf("get object fail, err: %v, obj key: %v", err, remotePath)
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

// 预签名，获取一个 url; 使用该 url 来 上传 obj.
// 使用 http put => PresignPutObject or post => PresignPostObject 方式上传 obj.
func (sc *S3ClientBase) GetUploadUrl(remotePath string, expireMilliSecond int) (string, error) {
	preSignCli := s3.NewPresignClient(sc.S3Client)
	if preSignCli == nil {
		log.Printf("create presign client fail, is nil")
		return "", fmt.Errorf("create presign client is nil")
	}

	// retPSC, err := preSignCli.PresignPostObject(context.Background(), &s3.PutObjectInput{
	// 	Bucket: aws.String(sc.CliCfg.Bucket),
	// 	Key:    aws.String(remotePath),
	// }, func(opts *s3.PresignPostOptions) {
	// 	opts.Expires = time.Duration(expireMilliSecond) * time.Millisecond
	// })

	retPSC, err := preSignCli.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(sc.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireMilliSecond) * time.Millisecond
	})

	if err != nil {
		log.Printf("get presignputobj fail, err: %v", err)
		return "", err
	}

	if retPSC == nil {
		log.Printf("get presign put obj is nil")
		return "", fmt.Errorf("get persing obj is nil")
	}

	return retPSC.URL, nil
}

// 预签名，获取 url 来下载 s3中的文件
func (sc *S3ClientBase) GetDownloadUrl(remotePath string, expMilliSecond int) (string, error) {
	preSignCli := s3.NewPresignClient(sc.S3Client)
	if preSignCli == nil {
		log.Printf("create presign client fail, is nil")
		return "", fmt.Errorf("create presign client is nil")
	}

	retPSC, err := preSignCli.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(sc.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expMilliSecond) * time.Millisecond
	})

	if err != nil {
		log.Printf("pre sign get obj fail, err: %v", err)
		return "", err
	}
	if retPSC == nil {
		log.Printf("pre sign get obj is nil")
		return "", fmt.Errorf("pre sign get obj is nil")
	}

	return retPSC.URL, nil
}


const (
	DownLoadPrivilegePrivate = 1
	DownLoadPrivilegeRead = 2 
	DownLoadPrivilegeRW = 3
	DownLoadPrivilegeAuthRead = 4
	DownLoadPrivilegeExecRead = 5 
	DownLoadPrivilegeOwnerRead = 6 
	DownLoadPrivilegeOwnerFullCtrl = 7
)


      
var privilegeInS3Inner = map[int]types.ObjectCannedACL {
	DownLoadPrivilegePrivate: types.ObjectCannedACLPrivate,
	DownLoadPrivilegeRead: types.ObjectCannedACLPublicRead,
	DownLoadPrivilegeRW: types.ObjectCannedACLPublicReadWrite,
	DownLoadPrivilegeAuthRead:types.ObjectCannedACLAuthenticatedRead ,
	DownLoadPrivilegeExecRead:types.ObjectCannedACLAwsExecRead , 
	DownLoadPrivilegeOwnerRead:types.ObjectCannedACLBucketOwnerRead , 
	DownLoadPrivilegeOwnerFullCtrl: types.ObjectCannedACLBucketOwnerFullControl,
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
