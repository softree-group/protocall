package storage

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"protocall/domain/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cheggaaa/pb"
)

const maxRetries = 12

type Config struct {
	AccessKey           string
	SecretKey           string
	SessionToken        string
	RegionName          string
	Endpoint            string
	DisableSSL          bool
	SkipSSLVerification bool
}

type progressReader struct {
	reader io.Reader
	pb     *pb.ProgressBar
}

func (r progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		return n, err
	}

	r.pb.Add(n)

	return n, nil
}

type UploadFileOptions struct {
	Acl                  string
	ServerSideEncryption string
	KmsKeyId             string
	ContentType          string
	DisableMultipart     bool
}

type s3client struct {
	client         *s3.S3
	session        *session.Session
	options        *UploadFileOptions
	progressOutput io.Writer
}

func NewS3Client(progressOutput io.Writer, awsConfig *aws.Config, useV2Signing bool) repository.VoiceStorage {
	sess := session.New(awsConfig)
	client := s3.New(sess, awsConfig)

	return &s3client{
		client:         client,
		session:        sess,
		progressOutput: progressOutput,
	}
}

func NewAwsConfig(c *Config) *aws.Config {
	var creds *credentials.Credentials

	if c.AccessKey == "" && c.SecretKey == "" {
		creds = credentials.AnonymousCredentials
	} else {
		creds = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)
	}

	if len(c.RegionName) == 0 {
		c.RegionName = "us-east-1"
	}

	var httpClient *http.Client
	if c.SkipSSLVerification {
		httpClient = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	} else {
		httpClient = http.DefaultClient
	}

	awsConfig := &aws.Config{
		Region:           aws.String(c.regionName),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(maxRetries),
		DisableSSL:       aws.Bool(disableSSL),
		HTTPClient:       httpClient,
	}

	if len(endpoint) != 0 {
		endpoint := fmt.Sprintf("%s", endpoint)
		awsConfig.Endpoint = &endpoint
	}

	return awsConfig
}

func (client *s3client) UploadFile(bucketName string, remotePath string, localPath string, options storage.UploadFileOptions) (string, error) {
	uploader := s3manager.NewUploaderWithClient(client.client)

	stat, err := os.Stat(localPath)
	if err != nil {
		return "", err
	}

	localFile, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	// Automatically adjust partsize for larger files.
	fSize := stat.Size()
	if !options.DisableMultipart {
		if fSize > int64(uploader.MaxUploadParts)*uploader.PartSize {
			partSize := fSize / int64(uploader.MaxUploadParts)
			if fSize%int64(uploader.MaxUploadParts) != 0 {
				partSize++
			}
			uploader.PartSize = partSize
		}
	} else {
		uploader.MaxUploadParts = 1
		uploader.Concurrency = 1
		uploader.PartSize = fSize + 1
		if fSize <= s3manager.MinUploadPartSize {
			uploader.PartSize = s3manager.MinUploadPartSize
		}
	}

	progress := client.newProgressBar(fSize)

	progress.Start()
	defer progress.Finish()

	uploadInput := s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remotePath),
		Body:   progressReader{localFile, progress},
		ACL:    aws.String(options.Acl),
	}
	if options.ServerSideEncryption != "" {
		uploadInput.ServerSideEncryption = aws.String(options.ServerSideEncryption)
	}
	if options.KmsKeyId != "" {
		uploadInput.SSEKMSKeyId = aws.String(options.KmsKeyId)
	}
	if options.ContentType != "" {
		uploadInput.ContentType = aws.String(options.ContentType)
	}

	uploadOutput, err := uploader.Upload(&uploadInput)
	if err != nil {
		return "", err
	}

	if uploadOutput.VersionID != nil {
		return *uploadOutput.VersionID, nil
	}

	return "", nil
}

func (client *s3client) DownloadFile(bucketName string, remotePath string, versionID string, localPath string) error {
	headObject := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remotePath),
	}

	if versionID != "" {
		headObject.VersionId = aws.String(versionID)
	}

	object, err := client.client.HeadObject(headObject)
	if err != nil {
		return err
	}

	progress := client.newProgressBar(*object.ContentLength)

	downloader := s3manager.NewDownloaderWithClient(client.client)

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	getObject := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remotePath),
	}

	if versionID != "" {
		getObject.VersionId = aws.String(versionID)
	}

	progress.Start()
	defer progress.Finish()

	_, err = downloader.Download(progressWriterAt{localFile, progress}, getObject)
	if err != nil {
		return err
	}

	return nil
}
