package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const DefaultPrefix string = "chamber/"
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
	if prefix == "" {
		prefix = DefaultPrefix
	}
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

func (s *S3Store) WriteAll(id string, secrets RawSecrets) error {
	// TODO: GET, PUT
	panic("not implemented")
}

func (s *S3Store) Write(id, key, value string) error {
	// TODO: GET, update, PUT
	panic("not implemented")
}

func (s *S3Store) ReadAll(id, version string) (*Secrets, error) {
	key := s.surround(id)
	result, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})

	if err != nil {
		return nil, err
	}

	rawSecrets := make(map[string]string)
	dec := json.NewDecoder(result.Body)
	if err = dec.Decode(&rawSecrets); err != nil {
		return nil, err
	}

	return &Secrets{
		Secrets: rawSecrets,
		Meta: &SecretMetadata{
			Version:      *result.VersionId,
			LastModified: *result.LastModified,
		},
	}, nil
}

func (s *S3Store) Read(id, key, version string) (string, *SecretMetadata, error) {
	secrets, err := s.ReadAll(id, version)
	if err != nil {
		return "", nil, err
	}
	val, ok := secrets.Secrets[key]
	if !ok {
		return "", nil, errors.New("key not found")
	}
	return val, secrets.Meta, nil
}

func (s *S3Store) Delete(id string) error {
	panic("not implemented")
}
