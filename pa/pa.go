package main

/*
#cgo pkg-config: portaudio-2.0
#include <portaudio.h>
extern PaStreamCallback* streamCallbackWrap;
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
	_ "os"
	_ "runtime"
	"reflect"
)

var stream=&Stream{
	in:  new(reflect.SliceHeader),
	out:  new(reflect.SliceHeader),
}



// DeviceInfo contains information for an audio device.
type DeviceInfo struct {
	index                    C.PaDeviceIndex
	Name                     string
	MaxInputChannels         int
	MaxOutputChannels        int
	DefaultLowInputLatency   time.Duration
	DefaultLowOutputLatency  time.Duration
	DefaultHighInputLatency  time.Duration
	DefaultHighOutputLatency time.Duration
	DefaultSampleRate        float64
}

type Stream struct {
	paStream            unsafe.Pointer
	inParams, outParams *C.PaStreamParameters
	in, out             *reflect.SliceHeader
	args                []reflect.Value
	callback            reflect.Value
	closed              bool
}

func main() {

	if err:=C.Pa_Initialize();err != C.paNoError {
		fmt.Println(err)
	}
	defer func(){
		if err := C.Pa_Terminate();err != C.paNoError {
			fmt.Println(err )
		}

		}()

	inputParameters := new(C.PaStreamParameters)
	inputParameters.device = C.Pa_GetDefaultInputDevice(); /* default input device */
	if (inputParameters.device == C.paNoDevice) {
				fmt.Println("Error: No default input device.")
	}
	inputParameters.channelCount = 1
	inputParameters.sampleFormat = C.paFloat32
  inputParameters.suggestedLatency = C.Pa_GetDeviceInfo( inputParameters.device ).defaultLowInputLatency
  inputParameters.hostApiSpecificStreamInfo = nil

		err := C.Pa_OpenStream(
                 &stream.paStream,
                 inputParameters,
                 nil,                  /* &outputParameters, */
                 C.double(44100),
                 C.paFramesPerBufferUnspecified,
                 C.paClipOff,      /* we won't output out of range samples so don't bother clipping them */
                 C.streamCallbackWrap,
                 nil );
       if( err != C.paNoError ) {
				 fmt.Println(err)
				 }

       ;
       if err := C.Pa_StartStream( unsafe.Pointer(stream.paStream) ); err != C.paNoError {
				 fmt.Println(err)
				 }


			 for x:=0;x<20;x++{
				 time.Sleep(1*time.Second)
				//  x:=C.Pa_IsStreamActive(unsafe.Pointer(stream))

				//  fmt.Println(x)
			 }
      //  while( ( err = Pa_IsStreamActive( stream ) ) == 1 )
      //  {
      //      Pa_Sleep(1000);
      //      printf("index = %d\n", data.frameIndex ); fflush(stdout);
      //  }

}

func duration(paTime C.PaTime) time.Duration {
	return time.Duration(paTime * C.PaTime(time.Second))
}

//export streamCallback
func streamCallback(inputBuffer, outputBuffer unsafe.Pointer, frames C.ulong, timeInfo *C.PaStreamCallbackTimeInfo, statusFlags C.PaStreamCallbackFlags) {

	updateBuffer(stream.in, uintptr(inputBuffer),int(frames))
	updateBuffer(stream.out, uintptr(outputBuffer), int(frames))
}

func updateBuffer(buf *reflect.SliceHeader, p uintptr, frames int) {
	if p == 0 {
		return
	}
	buf.Data = data
	buf.Len = n
	buf.Cap = n
		setSlice(buf, p, frames)

		fmt.Println(buf)

}

// func setChannels(s *reflect.SliceHeader, p uintptr, frames int) {
// 	sp := s.Data
// 	for i := 0; i < s.Len; i++ {
// 		setSlice((*reflect.SliceHeader)(unsafe.Pointer(sp)), *(*uintptr)(unsafe.Pointer(p)), frames)
// 		sp += unsafe.Sizeof(reflect.SliceHeader{})
// 		p += unsafe.Sizeof(uintptr(0))
// 	}
// }

func setSlice(s *reflect.SliceHeader, data uintptr, n int) {


}
