package manager

import (
	"bitbucket.org/weberc2/media/ao"
	"bitbucket.org/weberc2/media/mpg123"
	"errors"
	"fmt"
	"github.com/martin-reznik/logger"
	"io"
)

// Song - queue item containing info about playing sound
type Song struct {
	ID     int
	File   string
	logger *logger.Log
	handle *mpg123.Handle
}

var songList = make(map[int]*Song)

// NewSong - creates new song
func NewSong(filename string, log *logger.Log) *Song {
	s := &Song{1, filename, log, nil}
	songList[1] = s
	return s
}

// FindPlayingSong - finds and returns song by play id
func FindPlayingSong(id int) (*Song, error) {
	s, ok := songList[id]
	if !ok {
		return nil, errors.New("fujky")
	}
	return s, nil
}

// Play - plays song
func (s *Song) Play() {
	go func() {
		defer s.close()
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
		s.handle = handle
		defer dev.Close()

		if _, err := io.Copy(dev, handle); err != nil {
			s.logger.Warning(fmt.Sprintf("Stoping playback: %v", err))
		}
	}()
}

// Stop - stops playing (bit nasty way, but, well...)
func (s *Song) Stop() {
	s.handle.Close()
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

func (s *Song) close() {
	delete(songList, s.ID)
}
