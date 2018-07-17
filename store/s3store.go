package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const Suffix string = ".json"

type S3Store struct {
	Bucket string
	Prefix string
	svc    s3.S3
}

// ensure S3Store conforms to Store interface
var _ Store = &S3Store{}

// Establish an AWS Session and deduce the region:
// - from $CHAMBER_AWS_REGION
// - from EC2 metadata
func awsSessionAndRegion() (*session.Session, *string) {
	var region *string

	if regionOverride, ok := os.LookupEnv("CHAMBER_AWS_REGION"); ok {
		region = aws.String(regionOverride)
	}
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: region,
			},
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	// If region is still not set, attempt to determine it via ec2 metadata API
	region = nil
	if aws.StringValue(sess.Config.Region) == "" {
		sessEc2 := session.New()
		ec2metadataSvc := ec2metadata.New(sessEc2)
		if regionOverride, err := ec2metadataSvc.Region(); err == nil {
			region = aws.String(regionOverride)
		}
	}
	return sess, region
}

func NewS3Store(numRetries int, bucket string, prefix string) *S3Store {
	session, region := awsSessionAndRegion()

	svc := s3.New(session, &aws.Config{
		MaxRetries: aws.Int(numRetries),
		Region:     region,
	})
	s := S3Store{
		svc:    *svc,
		Bucket: bucket,
		Prefix: prefix,
	}

	return &s
}

func (s *S3Store) surround(id string) string {
	return s.Prefix + id + Suffix
}

// Return's Secrets.Meta.LastModifed will be time.Time 0-value
func (s *S3Store) WriteAll(id string, secrets RawSecrets) (string, error) {
	err := secrets.Validate()
	if err != nil {
		return "", err
	}
	key := s.surround(id)

	buf, err := json.Marshal(secrets)
	if err != nil {
		return "", err
	}

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(buf)),
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	result, err := s.svc.PutObject(input)
	if err != nil {
		return "", err
	}
	return *(result.VersionId), nil
}

// Updates Secrets at s3://<bucket>/<prefix><id><suffix> with key=value.
// If object doesn't exist yet, a new empty one is created.
// As with WriteAll, return's Secrets.Meta.LastModifed will be
// time.Time 0-value
func (s *S3Store) Write(id, key, value string) (string, error) {
	secs, err := s.ReadAll(id, "")
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != s3.ErrCodeNoSuchKey {
			return "", err
		}
		secs = &Secrets{}
	}

	(*secs).Secrets[key] = value
	return s.WriteAll(id, (*secs).Secrets)
}

// version == "" => latest
func (s *S3Store) ReadAll(id, version string) (*Secrets, error) {
	var versionId *string
	key := s.surround(id)
	if version == "" {
		versionId = nil
	} else {
		versionId = &version
	}
	result, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket:    &s.Bucket,
		Key:       &key,
		VersionId: versionId,
	})

	if err != nil {
		return nil, err
	}

	rawSecrets := make(RawSecrets)
	dec := json.NewDecoder(result.Body)
	if err = dec.Decode(&rawSecrets); err != nil {
		return nil, err
	}

	return &Secrets{
		Secrets: rawSecrets,
		Meta: &SecretsMetadata{
			Version:      *result.VersionId,
			LastModified: *result.LastModified,
		},
	}, nil
}

func (s *S3Store) Read(id, key, version string) (string, *SecretsMetadata, error) {
	secrets, err := s.ReadAll(id, version)
	if err != nil {
		return "", nil, err
	}
	val, ok := (*secrets).Secrets[key]
	if !ok {
		return "", nil, errors.New("key not found")
	}
	return val, secrets.Meta, nil
}

func (s *S3Store) DeleteAll(id string) error {
	panic("not implemented")
}

func (s *S3Store) Delete(id, key string) error {
	panic("not implemented")
}
