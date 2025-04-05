package s3SdkV1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	logger "log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AwsS3 = 1
)

type IS3Client interface {
	Upload(data io.ReadSeeker, remotePath string)
}

// S3ClientConfig s3 client的配置项
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

// S3ClientBase s3 client 类型的基类
type S3ClientBase struct {
	S3Client *s3.S3
	CliCfg   *S3ClientConfig
}

func BuildS3ObjectFile(dirPath, fileName string) string {
	//return filepath.Join(dirPath, fileName)
	return dirPath + fileName
}

// BuildBaseObjFileName 使用配置前缀，构造一个完整文件名
func (sb *S3ClientBase) BuildBaseObjFileName(fileName string) string {
	if sb == nil || sb.CliCfg == nil || sb.CliCfg.FilePrefix == "" {
		return fileName
	}
	return sb.CliCfg.FilePrefix + fileName
}

// Upload 直接上传文件， 其中 remotePath 为桶里面的文件名
func (sb *S3ClientBase) Upload(data io.ReadSeeker, remotePath string) error {
	if sb.CliCfg.CliTimeoutSecondCfg > 0 {
		ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(sb.CliCfg.CliTimeoutSecondCfg)*time.Second)
		defer cancelFn()

		_, err := sb.S3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(sb.CliCfg.Bucket),
			Key:    aws.String(remotePath),
			Body:   data})

		if err != nil {
			var aerr awserr.Error
			if errors.As(err, &aerr) && aerr.Code() == request.CanceledErrorCode {
				logger.Printf("upload data fail timeout: %v second", sb.CliCfg.CliTimeoutSecondCfg)
			} else {
				logger.Printf("upload data fail, err: %v", err)
			}
			return err
		}
		logger.Printf("upload to s3 succ.")

	} else {
		_, err := sb.S3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(sb.CliCfg.Bucket),
			Key:    aws.String(remotePath),
			Body:   data})

		if err != nil {
			logger.Printf("upload data fail, err: %v", err)
			return err
		} else {
			logger.Printf("upload to s3 succ.")
			return nil
		}
	}
	return nil
}

// GetDownloadUrl 获取下载链接， remotePath 为桶中文件名， urlExpireSecond 为下载url的有效时间
func (sb *S3ClientBase) GetDownloadUrl(remotePath string, urlExpireSecond int32) (string, error) {
	expires := time.Now().Add(time.Duration(urlExpireSecond) * time.Second)
	sdkReq, _ := sb.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:          aws.String(sb.CliCfg.Bucket),
		Key:             aws.String(remotePath),
		ResponseExpires: &expires,
	})
	url, _, err := sdkReq.PresignRequest(time.Duration(urlExpireSecond) * time.Second)
	if err != nil {
		logger.Printf("call PresignRequest fail, err: %v", err)
		return "", err
	}
	return url, nil
}

// Download 直接从 s3上 下载数据
func (sb *S3ClientBase) Download(remotePath string) {
	response, err := sb.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(sb.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	})

	if err != nil {
		logger.Printf("download fail, err: %v", err)
		return
	}
	if response == nil {
		logger.Printf("download response is nil")
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Printf("read data fail, err: %v", err)
		return
	}
	logger.Printf("%v", string(data))
}

func (sb *S3ClientBase) DownloadFile(remotePath string) ([]byte, error) {
	response, err := sb.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(sb.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	})

	if err != nil {
		logger.Printf("download fail, err: %v", err)
		return nil, err
	}
	if response == nil {
		logger.Printf("download response is nil")
		return nil, errors.New("download response is nil")
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Printf("read data fail, err: %v", err)
		return nil, err
	}

	return data, err
}

func (sb *S3ClientBase) DownloadToLocal(remotePath string, localPath string) error {
	response, err := sb.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(sb.CliCfg.Bucket),
		Key:    aws.String(remotePath),
	})

	if err != nil {
		logger.Printf("download fail, err: %v", err)
		return err
	}
	if response == nil {
		logger.Printf("download response is nil")
		return fmt.Errorf("file is empty")
	}
	defer response.Body.Close()

	file, err := os.Create(localPath) // 替换为你要保存的文件名
	if err != nil {
		logger.Printf("Failed to create file: %v", err)
		return fmt.Errorf("open file fail, file: %v, err: %v", localPath, err)
	}
	defer file.Close()

	_, err = file.ReadFrom(response.Body)
	if err != nil {
		logger.Printf("Failed to write to file: %v, err: %v", localPath, err)
		return err
	}
	return nil
}

func (sb *S3ClientBase) CheckPresignedURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Error accessing the URL %s : %s", url, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return false // URL 已过期
	}
	return true // URL 仍然有效
}
func DownloadByPreUrl(url, localFile string) error {
	if url == "" || localFile == "" {
		return fmt.Errorf("input param is invalid")
	}
	// 使用预签名 URL 进行文件下载
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Failed to download file: %v, url: %v", err, url)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("Failed to download file: %v, url: %v", resp.Status, url)
		return fmt.Errorf("download status code not 200")
	}

	file, err := os.Create(localFile) // 替换为你要保存的文件名
	if err != nil {
		logger.Printf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		logger.Printf("Failed to read response: %v", err)
		return fmt.Errorf("read response body and write to file: %v fail: %v", localFile, err)
	}
	return nil
}
func DownloadToCacheByPreUrl(url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("input param is invalid")
	}
	// 使用预签名 URL 进行文件下载
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Failed to download file: %v, url: %v", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("Failed to download file: %v, url: %v", resp.Status, url)
		return nil, fmt.Errorf("download status code not 200")
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		logger.Printf("read response body to buf fail, err: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// S3ClientS3 s3 client
type S3ClientS3 struct {
	S3ClientBase
}

// NewS3ClientS3 创建 s3 S3Client 对象
func NewS3ClientS3(cfg *S3ClientConfig) *S3ClientS3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.Ak, cfg.Sk, ""),
	}))
	svc := s3.New(sess, aws.NewConfig().WithRegion(cfg.Region))
	cli := &S3ClientS3{
		S3ClientBase: S3ClientBase{
			S3Client: svc,
			CliCfg:   cfg,
		},
	}
	return cli
}
