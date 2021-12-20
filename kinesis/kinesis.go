package kinesis

import (
	"archive/zip"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"

	// "github.comcast.com/mesa/sky-soc-e2e-testing/constants"
	"io/ioutil"
	"log"
	"sky-soc-e2e-testing/constants"
	"strings"
	"time"
)

type Client struct {
	kinesisClient *kinesis.Kinesis
}

func New(sess *session.Session) *Client {
	return &Client{kinesisClient: kinesis.New(sess)}
}

func (c *Client) Listen(resultsChannel chan string) {
	iteratorOutput, err := c.kinesisClient.GetShardIterator(&kinesis.GetShardIteratorInput{
		// Shard Id is provided when making put record(s) request.
		ShardId:           aws.String("shardId-000000000000"),
		ShardIteratorType: aws.String("LATEST"),
		StreamName:        aws.String(constants.Constants.KinesisStreamName),
	})
	if err != nil {
		close(resultsChannel)

		log.Fatalln(err)
		return
	}
	var shardIterator *string
	shardIterator = iteratorOutput.ShardIterator

	for shardIterator != nil {
		getRecordsInput := &kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		}

		records, err := c.kinesisClient.GetRecords(getRecordsInput)
		if err != nil {
			close(resultsChannel)
			log.Fatalln(err)
			return
		}

		shardIterator = records.NextShardIterator

		for _, record := range records.Records {
			//log.Println(record)
			reader, err := zip.NewReader(strings.NewReader(string(record.Data)), int64(len(record.Data)))
			if err != nil {
				close(resultsChannel)
				log.Fatalln(err)
				return
			}

			// Read all the files from zip archive
			for _, zipFile := range reader.File {
				unzippedFileBytes, err := readZipFile(zipFile)
				if err != nil {
					log.Println(err)
					continue
				}

				resultsChannel <- string(unzippedFileBytes)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
