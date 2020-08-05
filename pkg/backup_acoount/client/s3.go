package client

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type s3Client struct {
	Vars map[string]string
	Sess session.Session
}

func NewS3Client(vars map[string]string) (*s3Client, error) {

	var accessKey string
	var secretKey string
	var endpoint string
	var region string
	if _, ok := vars["accessKey"]; ok {
		accessKey = vars["accessKey"]
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["secretKey"]; ok {
		secretKey = vars["secretKey"]
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["endpoint"]; ok {
		endpoint = vars["endpoint"]
	} else {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["region"]; ok {
		region = vars["region"]
	} else {
		return nil, errors.New(ParamEmpty)
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	return &s3Client{
		Vars: vars,
		Sess: *sess,
	}, nil
}

func (s3C s3Client) ListBuckets() ([]interface{}, error) {
	var result []interface{}
	svc := s3.New(&s3C.Sess)
	res, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}
	for _, b := range res.Buckets {
		result = append(result, b.Name)
	}
	return result, nil
}

func (s3C s3Client) Exist(path string) (bool, error) {
	bucket, err := s3C.getBucket()
	if err != nil {
		return false, err
	}
	svc := s3.New(&s3C.Sess)
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})
	if err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok {
			if aerr.StatusCode() == 404 {
				return false, nil
			}
		} else {
			return false, aerr
		}
	}
	return true, nil
}

func (s3C s3Client) Delete(path string) (bool, error) {
	bucket, err := s3C.getBucket()
	if err != nil {
		return false, err
	}
	svc := s3.New(&s3C.Sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(path)})
	if err != nil {
		return false, err
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s3C s3Client) Upload(src, target string) (bool, error) {
	bucket, err := s3C.getBucket()
	if err != nil {
		return false, err
	}
	file, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(&s3C.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(target),
		Body:   file,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s3C s3Client) Download(src, target string) (bool, error) {
	bucket, err := s3C.getBucket()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			os.Remove(target)
		} else {
			return false, err
		}
	}
	file, err := os.Create(target)
	if err != nil {
		return false, err
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(&s3C.Sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(src),
		})
	if err != nil {
		os.Remove(target)
		return false, err
	}
	return true, nil
}

func (s3C *s3Client) getBucket() (string, error) {
	if _, ok := s3C.Vars["bucket"]; ok {
		return s3C.Vars["bucket"], nil
	} else {
		return "", errors.New(ParamEmpty)
	}
}
