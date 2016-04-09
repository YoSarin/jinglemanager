package manager

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"bytes"
	"errors"
	"fmt"
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
	stream       []byte
	ao           *ao.SampleFormat
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
		return nil, fmt.Errorf("File %v does not exist", filename)
	}
	lastID++
	stream, ao := getMusicStream(filename)
	s := &Song{
		ID:           lastID,
		File:         filename,
		logger:       log,
		stream:       stream,
		ao:           ao,
		playing:      false,
		bytesPlayed:  0,
		bytesTotal:   int64(len(stream)),
		stopPlayback: make(chan bool, 1),
	}
	songList[lastID] = s
	return s, nil
}

// FindSong - finds and returns song by play id
func FindSong(id int) (*Song, error) {
	s, ok := songList[id]
	if !ok {
		return nil, errors.New("Song not found")
	}
	return s, nil
}

// FindSongByFile - finds if we already have this song prepared
func FindSongByFile(filename string) *Song {
	for _, s := range songList {
		if s.File == filename {
			return s
		}
	}
	return nil
}

// GetAllPlaying - lists all songs which are in progress
func GetAllPlaying() []*Song {
	out := []*Song{}
	for _, s := range songList {
		out = append(out, s)
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

		ao.Initialize()
		defer ao.Shutdown()
		dev := ao.NewLiveDevice(s.ao)
		defer dev.Close()

		bufSize := int64(1024)
		for step := int64(0); step < s.bytesTotal/bufSize; step++ {
			select {
			case <-s.stopPlayback:
				s.logger.Info("Playback stopped")
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
	s.stopPlayback <- true
}

// Position - percentage of file played (based on size)
func (s *Song) Position() float64 {
	return float64(s.bytesPlayed) / float64(s.bytesTotal)
}

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
	fmt.Printf("%v", handle)
	_, err = io.Copy(rw, handle)
	if err != nil {
		panic("File read error " + err.Error())
	}
	buffer, err := ioutil.ReadAll(rw)
	aoData := aoPrepare(handle)

	return buffer, aoData
}
