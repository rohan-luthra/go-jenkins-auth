package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	// "github.comcast.com/mesa/sky-soc-e2e-testing/constants"
	"log"
	"sky-soc-e2e-testing/constants"
)

type Client struct {
	s3Client *s3.S3
}

func New(sess *session.Session) *Client {
	return &Client{s3Client: s3.New(sess)}
}

func (c *Client) GetBucketAssets() ([]string, error) {
	loi := s3.ListObjectsInput{
		Bucket: aws.String(constants.Constants.S3BucketName),
		Prefix: aws.String(constants.Constants.BucketRoot),
	}
	out, err := c.s3Client.ListObjects(&loi)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	assets := make([]string, 0, len(out.Contents))

	for _, o := range out.Contents {
		assets = append(assets, *o.Key)
	}

	return assets, nil
}
