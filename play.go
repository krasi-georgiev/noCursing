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
	"log"

	// Imports the Google Cloud Speech API client package.

	"github.com/gordonklaus/portaudio"
)

func main() {

	portaudio.Initialize()
	defer portaudio.Terminate()

	bufMic := make([]int32, 1)

	// (numInputChannels, numOutputChannels int, sampleRate float64, framesPerBuffer int, args ...interface{}) (*Stream, error)
	micStream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(bufMic), bufMic)
	if err != nil {
		log.Fatal(err)
	}

	defer micStream.Close()

	micStream.Start()

	for {
		micStream.Read()

		bufGoogle := make([]byte, 4)
		for _, v := range bufMic {

			binary.BigEndian.PutUint32(bufGoogle[:], uint32(v))

		}

		fmt.Printf("%v ", bufGoogle[0])
		fmt.Printf("%v ", bufGoogle[1])
		fmt.Printf("%v ", bufGoogle[2])
		fmt.Printf("%v ", bufGoogle[3])

	}

}
