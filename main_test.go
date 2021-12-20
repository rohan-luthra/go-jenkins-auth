package main

import (
	"testing"
)

var test map[string]string = map[string]string{
	"video1": "QUEUED",
}

func Test(t *testing.T) {
	resultsChannel := make(chan string)

	for i := 0; i < 10; i++ {
		go pushToChannel(resultsChannel, i)
	}

	res, err := checkResp(resultsChannel)
	if err != nil {
		t.Error(err)
	}

	for id, status := range res {
		if test[id] != "" && test[id] != status {
			t.Error("Invalid status")
		}
	}
}
