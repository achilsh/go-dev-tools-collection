package main

import (
	"bytes"
	"strings"

	s3SdkV1 "github.com/achilsh/go-dev-tools-collection/s3client_wrapper/lib/v1"
)

func main() {
	s3v1Client := s3SdkV1.NewS3ClientS3(
		&s3SdkV1.S3ClientConfig{
			Bucket:              "config.GetConfig().demoS3.Bucket",
			Region:              "config.GetConfig().demoS3.Region",
			Ak:                  "config.GetConfig().demoS3.Ak",
			Sk:                  "config.GetConfig().demoS3.Sk",
			FilePrefix:          "config.GetConfig().demoS3.FilePrefix",
			StoreType:           0,
			CliTimeoutSecondCfg: 10, //"config.GetConfig().demoS3.UploadTimeoutSecond",
		})

	demoData := "this is 000000"
	demoIo := strings.NewReader(demoData)

	s3v1Client.Upload(demoIo, s3SdkV1.BuildS3ObjectFile("cloud/track/", s3v1Client.BuildBaseObjFileName("demo-test_1.txt")))

	bufData := []byte("this is buf data")
	bytesData := bytes.NewReader(bufData)
	s3v1Client.Upload(bytesData, s3SdkV1.BuildS3ObjectFile("cloud/track/", s3v1Client.BuildBaseObjFileName("demo-test_2.txt")))

	dataBuf, _ := s3v1Client.DownloadFile("to_download_url")
	_ = dataBuf
	//
	s3v1Client.Download(s3SdkV1.BuildS3ObjectFile("cloud/track/", s3v1Client.BuildBaseObjFileName("demo-test_1.txt")))

	//
	downloadUrl, _ := s3v1Client.GetDownloadUrl("", 10)
	_ = downloadUrl

	//
	expired := s3v1Client.CheckPresignedURL("download_url")
	_ = expired
}
