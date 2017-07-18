#include <portaudio.h>

void streamCallback( void* inputBuffer, void *outputBuffer, unsigned long frames, const PaStreamCallbackTimeInfo *timeInfo, PaStreamCallbackFlags statusFlags);

int cb(const void *inputBuffer, void *outputBuffer, unsigned long frames, const PaStreamCallbackTimeInfo *timeInfo, PaStreamCallbackFlags statusFlags, void *userData) {
	streamCallback((void*)inputBuffer, outputBuffer, frames, (PaStreamCallbackTimeInfo*)timeInfo, statusFlags);
	return paContinue;
}
PaStreamCallback* streamCallbackWrap = cb;
