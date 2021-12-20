package main

import (
	"encoding/json"
	"fmt"

	"log"
	"time"
)

// var inputJSON = `{
//   "source": "test",
//   "id": "%s",
//   "type": "processing-request",
//   "details": {
//     "responseEndpoint": {
//       "streamArn": "test",
//       "assumeRole": "test"
//     },
//     "program": {
//       "title": "titleid1234",
//       "subtitle": "",
//       "description": "",
//       "startTime": 1628274600000,
//       "expiryTime": 1630866600000,

//     },
//     "videoAnalysisFiles": [
//       {
//         "videoAnalysisFile": "test.mpg",
//         "captionDataFile": ""
//       }
//     ]
//   }
// }`

// func main() {
// 	cfg := aws.Config{
// 		Region:                        aws.String(constants.Constants.Region),
// 		CredentialsChainVerboseErrors: aws.Bool(true),
// 	}

// 	// sess := session.Must(session.NewSessionWithOptions(session.Options{Config: cfg}))

// 	// s3Client := s3.New(sess)
// 	// kinesisClient := kinesis.New(sess)

// 	// assets, err := s3Client.GetBucketAssets()
// 	// if err != nil {
// 	// 	os.Exit(1)
// 	// }

// 	// httpClient := &http.Client{
// 	// 	Timeout: 10 * time.Second,
// 	// }

// 	assets := []string{}

// 	log.Println("Asset Count: ", len(assets))

// 	resultsChannel := make(chan string, len(assets))
// 	sentAssets := make(map[string]bool, len(assets))
// 	for idx, asset := range assets {
// 		asset = asset
// 		idx = idx

// 		sentAssetID := uuid.New().String()
// 		sentAssets[sentAssetID] = false

// 		go func(asset, assetID string, idx int) {
// 			input := fmt.Sprintf(inputJSON, assetID, asset, asset)

// 			req, err := http.NewRequest(http.MethodPost, constants.Constants.APIGatewayURL, strings.NewReader(input))
// 			if err != nil {
// 				log.Println("unable to build request: ", err)
// 				return
// 			}

// 			req.Header.Set("x-api-key", "")

// 			// res, err := httpClient.Do(req)
// 			// if err != nil {
// 			// 	log.Println("unable to build response: ", err)
// 			// 	return
// 			// }

// 			// if res.StatusCode != http.StatusOK {
// 			// 	log.Println("failure sending asset: ", asset)
// 			// }

// 			// err = res.Body.Close()
// 			// if err != nil {
// 			// 	log.Println("unable to close response body: ", err)
// 			// 	return
// 			// }

// 			resultsChannel <- `{
//     "status": "QUEUED",
//     "jobRequest": {
//         "source": "test",
//         "id": "test_123",
//         "type": "test-request",
//         "details": {
//             "responseEndpoint": {
//                 "streamArn": "test",
//                 "assumeRole": "test"
//             },
//             "program": {
//                 "title": "titleid1234",
//                 "subtitle": "",
//                 "description": "",
//                 "startTime": 1628274600000,
//                 "expiryTime": 1630866600000,

//             },
//             "videoAnalysisFiles": [
//                 {
//                     "videoAnalysisFile": "test.mpg",
//                     "captionDataFile": ""
//                 }
//             ]
//         }
//     }
// }
// `
// 		}(asset, sentAssetID, idx)
// 	}

// 	go func(resultsChannel chan string) {
// 		kinesisClient.Listen(resultsChannel)
// 	}(resultsChannel)

// 	failureCnt := len(assets)

// 	for {
// 		select {
// 		case result := <-resultsChannel:

// 			var resultMap map[string]interface{}
// 			_ = json.Unmarshal([]byte(result), &resultMap)

// 			jobStatus := resultMap["status"].(string)
// 			fmt.Println(jobStatus)

// 			if resultMap["jobRequest"].(map[string]interface{})["id"] != nil {
// 				uniqueID := resultMap["jobRequest"].(map[string]interface{})["id"].(string)

// 				if _, ok := sentAssets[uniqueID]; ok && jobStatus == "COMPLETED" {
// 					sentAssets[uniqueID] = true
// 					failureCnt -= 1
// 				}

// 				if failureCnt <= 0 {
// 					fmt.Println("All assets are successful")
// 					return
// 				}
// 			}
// 		case <-time.After(360 * time.Second):
// 			if failureCnt > 0 {
// 				log.Printf("%d assets have failed", failureCnt)
// 			}
// 			return
// 		}
// 	}
// }

func pushToChannel(resultsChannel chan string, i int) error {

	resultsChannel <- `{
		"status": "QUEUED",
		"jobRequest": {
			"source": "test",
			"id": "test_123",
			"type": "test-request",
			"details": {
				"responseEndpoint": {
					"streamArn": "test",
					"assumeRole": "test"
				},
				"program": {
					"title": "titleid1234",
					"subtitle": "",
					"description": "",
					"startTime": 1628274600000,
					"expiryTime": 1630866600000
				},
				"videoAnalysisFiles": [{
					"videoAnalysisFile": "test.mpg",
					"captionDataFile": ""
				}]
			}
		}
	}`

	return nil

}

func checkStatusAndId(result string) (string, string, error) {

	var resultMap map[string]interface{}
	err := json.Unmarshal([]byte(result), &resultMap)
	if err != nil {
		return "", "", err
	}

	if resultMap["jobRequest"].(map[string]interface{})["id"] == nil {
		return "", "", fmt.Errorf("no id found in response")
	}

	uniqueID := resultMap["jobRequest"].(map[string]interface{})["id"].(string)

	jobStatus := resultMap["status"].(string)
	log.Println(jobStatus, uniqueID)
	return jobStatus, uniqueID, nil

}

func checkResp(resultsChannel chan string) (map[string]string, error) {

	sentAssets := make(map[string]string)
	failureCnt := 100

	for {
		select {
		case result := <-resultsChannel:
			jobStatus, uniqueID, err := checkStatusAndId(result)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			sentAssets[uniqueID] = jobStatus
			// TODO do something with uniueqID and Status
			// if _, ok := sentAssets[uniqueID]; ok && jobStatus == "QUEUED" {

			// 	failureCnt -= 1
			// }

		case <-time.After(10 * time.Second):
			if failureCnt > 0 {
				log.Printf("%d assets have failed", failureCnt)
			}
			return sentAssets, nil
		}
	}

}

func main() {

	resultsChannel := make(chan string)

	for i := 0; i < 10; i++ {
		go pushToChannel(resultsChannel, i)
	}

	checkResp(resultsChannel)

}
