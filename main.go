// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START speech_quickstart]
// Sample speech-quickstart uses the Google Cloud Speech API to transcribe
// audio.
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	// Imports the Google Cloud Speech API client package.
	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

	"github.com/gordonklaus/portaudio"
)

func main() {
	exportKey()

	ctx := context.Background()
	c, err := speech.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	stream, err := c.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		portaudio.Initialize()
		defer portaudio.Terminate()

		bufMic := make([]int32, 4)

		// (numInputChannels, numOutputChannels int, sampleRate float64, framesPerBuffer int, args ...interface{}) (*Stream, error)
		micStream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(bufMic), bufMic)
		if err != nil {
			log.Fatal(err)
		}

		defer micStream.Close()

		micStream.Start()

		bufGoogle := make([]byte, 16)

		if err := stream.Send(&speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
				StreamingConfig: &speechpb.StreamingRecognitionConfig{
					Config: &speechpb.RecognitionConfig{
						Encoding:        speechpb.RecognitionConfig_LINEAR16,
						SampleRateHertz: 44100,
					},
				},
			},
		}); err != nil {
			log.Fatal(err)
		}

		for {
			micStream.Read()

			var startSend bool
			for k, v := range bufMic {
				start := k * 4

				fmt.Println(v)
				if v != 0 {
					v = 255
				}
				binary.BigEndian.PutUint32(bufGoogle[start:start+4], uint32(v))
				if v != 0 {
					startSend = true
				}
			}

			// fmt.Println(bufGoogle[1])
			// fmt.Println("byte")
			// fmt.Println(fmt.Sprintf("%08b", bufGoogle[1]))
			// fmt.Println("bit")

			if startSend {
				// fmt.Println("Sending")
				// fmt.Println(bufGoogle)
				if err = stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: bufGoogle[:],
					},
				}); err != nil {
					log.Printf("sending audio error: %v", err)
				}
			}
		}

	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: handle error
			continue
		}
		if resp.Error != nil {
			// TODO: handle error
			continue
		}
		for _, result := range resp.Results {
			fmt.Printf("result: %+v\n", result)
		}
	}
}

func exportKey() {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var file string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".json", f.Name())
			if err == nil && r {
				file = f.Name()
			}
		}
		return nil
	})
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", pathS+"/"+file)
}
