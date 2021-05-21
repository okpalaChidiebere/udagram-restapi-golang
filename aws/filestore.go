package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/udacity/udagram-restapi-golang/config"
)

var (
	c    = config.NewConfig()
	sess *session.Session
)

type S3client struct {
	client *s3.S3
}

// Creates a S3 client
func createS3Client() *s3.S3 {

	if c.Aws_profile != "DEPLOYED" {
		/*
			Initialize a session that the SDK will use to load
			credentials from the shared credentials file ~/.aws/credentials
			and region from the shared configuration file ~/.aws/config.
		*/
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           c.Aws_profile,
		}))
	} else {
		sess = session.Must(session.NewSession())
	}

	svc := s3.New(sess, aws.NewConfig().WithRegion(c.Aws_region))
	return svc
}

// NewDynamoDbRepo creates a new DynamoDb Repository
func NewS3client() *S3client {
	s3c := createS3Client()

	return &S3client{s3c}
}

/* GetGetSignedUrl generates an aws signed url to retreive an item
 * @Params
 *    key: string - the filename to be put into the s3 bucket
 * @Returns:
 *    a url as a string and error
 */
func (s *S3client) GetGetSignedUrl(key string) (string, error) {

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.Aws_media_bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(5 * time.Minute) //we want the expire time of the url to be about 5 minutes

	return urlStr, err

}

/* GetPutSignedUrl generates an aws signed url to put an item
 * @Params
 *    key: string - the filename to be retreived from s3 bucket
 * @Returns:
 *    a url as a string and error
 */
func (s *S3client) GetPutSignedUrl(key string) (string, error) {
	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(c.Aws_media_bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(5 * time.Minute) //we want the expire time of the url to be about 5 minutes

	return urlStr, err
}
