package constants

type Constant struct {
	S3BucketName      string
	BucketRoot        string
	KinesisStreamName string
	Region            string
	APIGatewayURL     string
}

var Constants Constant = *&Constant{}
