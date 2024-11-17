package s3repo

import (
	"bytes"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/rahul-aut-ind/service-user/internal/awsconfig"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"golang.org/x/net/context"
)

type (
	S3Handler interface {
		Save(uID string, imageID uuid.UUID, ext string, f *[]byte) (string, error)
		Delete(uID string, imageID string) error
		DeleteAll(uID string) error
	}

	S3Repo struct {
		log       *logger.Logger
		client    *s3.Client
		bucket    string
		directory string
	}
)

// New creates a new instance of S3Repo
func New(l *logger.Logger, cfg *awsconfig.AWSConfig, env *config.Env) *S3Repo {
	return &S3Repo{
		log:       l,
		client:    initializeClient(cfg.Config, env.DynamoDBConnectionString),
		bucket:    env.S3Bucket,
		directory: env.S3Directory,
	}
}

func initializeClient(cfg *aws.Config, connString string) *s3.Client {
	return s3.NewFromConfig(*cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(connString)
	})
}

func (r *S3Repo) Save(uID string, imageID uuid.UUID, ext string, d *[]byte) (string, error) {
	f := r.getPath(uID, fmt.Sprintf("%s%s", imageID.String(), ext))

	_, err := r.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &f,
		Body:   bytes.NewReader(*d),
	})
	if err != nil {
		r.log.Fatalf("S3 PutObject err %+v", err)
		return f, err
	}

	return f, nil
}

// nolint:unused // needed for local debug sometimes, not part of functionality
func (r *S3Repo) logBucketList() {
	buckets, err := r.client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		r.log.Fatalf("failed to list S3 buckets:%v", err)
	}
	for _, bucket := range buckets.Buckets {
		r.log.Infof("bucket list :%v", *bucket.Name)
	}
}

func (r *S3Repo) getPath(uID, imageID string) string {
	return path.Join(r.directory, uID, imageID)
}

func (r *S3Repo) Delete(uID, imageID string) error {
	f := r.getPath(uID, imageID)

	_, err := r.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key:    &f,
	})

	return err
}

func (r *S3Repo) DeleteAll(uID string) error {
	prefix := r.getPath(uID, "")

	listOutput, err := r.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &r.bucket,
		Prefix: &prefix,
	})
	if err != nil {
		r.log.Errorf("error listing objects in bucket %s with prefix %s: %v", r.bucket, prefix, err)
		return err
	}

	// Prepare the list of objects to delete
	var objectsToDelete = make([]types.ObjectIdentifier, len(listOutput.Contents))
	for i, object := range listOutput.Contents {
		objectsToDelete[i] = types.ObjectIdentifier{
			Key: object.Key,
		}
	}

	if len(objectsToDelete) == 0 {
		return nil
	}

	_, err = r.client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: &r.bucket,
		Delete: &types.Delete{
			Objects: objectsToDelete,
		},
	})

	if err != nil {
		r.log.Errorf("error deleting objects in bucket %s with prefix %s: %v", r.bucket, prefix, err)
	}

	return err
}
