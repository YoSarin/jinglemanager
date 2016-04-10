package manager

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// LogI - intrface for logger
type LogI interface {
	Info(string)
	Warning(string)
	Error(string)
	Debug(string)
	Notice(string)
}

// Song - queue item containing info about playing sound
type Song struct {
	id           string
	File         string
	logger       LogI
	stream       []byte
	ao           *ao.SampleFormat
	playing      bool
	bytesPlayed  int64
	bytesTotal   int64
	stopPlayback chan bool
	done         chan bool
}

// NewSong - creates new song
func NewSong(filename string, log LogI) (*Song, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("File %v does not exist", filename)
	}
	stream, ao := getMusicStream(filename)
	s := &Song{
		id:           fmt.Sprintf("%x", md5.Sum(stream)),
		File:         filename,
		logger:       log,
		stream:       stream,
		ao:           ao,
		playing:      false,
		bytesPlayed:  0,
		bytesTotal:   int64(len(stream)),
		stopPlayback: make(chan bool, 1),
		done:         make(chan bool, 1),
	}
	return s, nil
}

// MarshalJSON - will convert song to JSON
func (s *Song) MarshalJSON() ([]byte, error) {
	data := struct {
		ID        string
		File      string
		IsPlaying bool
		Position  float64
	}{
		s.id, s.File, s.IsPlaying(), s.Position(),
	}
	return json.Marshal(data)
}

// ########### SongI implementation ##############

// IsPlaying - will return true if song is playing right now
func (s *Song) IsPlaying() bool {
	return s.playing
}

// Play - plays song
func (s *Song) Play() {
	go func() {
		s.playing = true
		defer s.playbackDone()

		ao.Initialize()
		defer ao.Shutdown()
		dev := ao.NewLiveDevice(s.ao)
		defer dev.Close()

		bufSize := int64(1024)
		for step := int64(s.bytesPlayed / bufSize); step < s.bytesTotal/bufSize; step++ {
			select {
			case <-s.stopPlayback:
				return
			default:
				size, err := dev.Write(s.stream[step*bufSize : (step+1)*bufSize])
				if err != nil {
					s.logger.Error(err.Error())
					return
				}
				s.bytesPlayed += int64(size)
			}
		}
	}()
}

// Stop - stops playing
func (s *Song) Stop() {
	s.Pause()
	// reset where we were in song to 0
	s.bytesPlayed = 0
}

// Pause - pauses playing, so it can be resumed from the point where it was stopped
func (s *Song) Pause() {
	s.stopPlayback <- true
	// wait for confirmation that playback has stopped
	_ = <-s.done
}

// Position - percentage of file played (based on size)
func (s *Song) Position() float64 {
	return float64(s.bytesPlayed) / float64(s.bytesTotal)
}

// ########### FileListItem implementation ##############

// ID - will return songs ID
func (s Song) ID() string {
	return s.id
}

// FileName - will return filename of song
func (s Song) FileName() string {
	return s.File
}

// ########### private methods ##############
// Get the ao.SampleFormat from the mpg123.Handle
func aoPrepare(handle *mpg123.Handle) *ao.SampleFormat {
	const bitsPerByte = 8

	rate, channels, encoding := handle.Format()

	return &ao.SampleFormat{
		BitsPerSample: handle.EncodingSize(encoding) * bitsPerByte,
		Rate:          int(rate),
		Channels:      channels,
		ByteFormat:    ao.FormatNative,
		Matrix:        nil,
	}
}

func (s *Song) playbackDone() {
	s.playing = false
	s.done <- true
}

func getMusicStream(filename string) ([]byte, *ao.SampleFormat) {
	mpg123.Initialize()
	defer mpg123.Exit()

	handle, err := mpg123.Open(filename)
	if err != nil {
		print(err.Error())
	}
	defer handle.Close()

	ao.Initialize()
	defer ao.Shutdown()

	dev := ao.NewLiveDevice(aoPrepare(handle))
	defer dev.Close()

	rw := new(bytes.Buffer)
	_, err = io.Copy(rw, handle)
	if err != nil {
		panic("File read error " + err.Error())
	}
	buffer, err := ioutil.ReadAll(rw)
	aoData := aoPrepare(handle)

	return buffer, aoData
}
