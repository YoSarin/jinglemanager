package manager

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"bytes"
	"errors"
	"github.com/martin-reznik/logger"
	"io"
	"io/ioutil"
	"os"
)

// Song - queue item containing info about playing sound
type Song struct {
	ID           int
	File         string
	logger       *logger.Log
	playing      bool
	bytesPlayed  int64
	bytesTotal   int64
	stopPlayback chan bool
}

var songList = make(map[int]*Song)
var lastID = 0

// NewSong - creates new song
func NewSong(filename string, log *logger.Log) (*Song, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("File does not exist")
	}
	lastID++
	s := &Song{
		ID:           lastID,
		File:         filename,
		logger:       log,
		playing:      false,
		bytesPlayed:  0,
		bytesTotal:   0,
		stopPlayback: make(chan bool, 1),
	}
	songList[lastID] = s
	return s, nil
}

// FindPlayingSong - finds and returns song by play id
func FindPlayingSong(id int) (*Song, error) {
	s, ok := songList[id]
	if !ok || !s.IsPlaying() {
		return nil, errors.New("fujky")
	}
	return s, nil
}

// GetAllPlaying - lists all songs which are in progress
func GetAllPlaying() []*Song {
	out := []*Song{}
	for _, s := range songList {
		if s.IsPlaying() {
			out = append(out, s)
		}
	}
	return out
}

// IsPlaying - will return true if song is playing right now
func (s *Song) IsPlaying() bool {
	return s.playing
}

// Play - plays song
func (s *Song) Play() {
	s.bytesPlayed = 0
	go func() {
		s.playing = true
		defer s.playbackDone()
		mpg123.Initialize()
		defer mpg123.Exit()

		handle, err := mpg123.Open(s.File)
		if err != nil {
			print(err.Error())
		}
		defer handle.Close()

		ao.Initialize()
		defer ao.Shutdown()
		dev := ao.NewLiveDevice(s.aoPrepare(handle))
		defer dev.Close()

		rw := new(bytes.Buffer)
		if s.bytesTotal, err = io.Copy(rw, handle); err != nil {
			s.logger.Error(err.Error())
			return
		}
		buffer, err := ioutil.ReadAll(rw)

		bufSize := int64(1024)
		for step := int64(0); step < s.bytesTotal/bufSize; step++ {
			select {
			case <-s.stopPlayback:
				s.logger.Info("Playback stopped")
				return
			default:
				size, err := dev.Write(buffer[step*bufSize : (step+1)*bufSize])
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
	s.stopPlayback <- true
}

// Position - percentage of file played (based on size)
func (s *Song) Position() float64 {
	return float64(s.bytesPlayed) / float64(s.bytesTotal)
}

// Get the ao.SampleFormat from the mpg123.Handle
func (s *Song) aoPrepare(handle *mpg123.Handle) *ao.SampleFormat {
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
}
