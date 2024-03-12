package logger

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type mockIoWriter struct {
	called     int
	data       [][]byte
	startTimes []time.Time
	endTimes   []time.Time
	writeDelay time.Duration
}

func (m *mockIoWriter) Write(p []byte) (n int, err error) {
	m.called++
	m.data = append(m.data, p)
	m.startTimes = append(m.startTimes, time.Now())

	if m.writeDelay > 0 {
		time.Sleep(m.writeDelay)
	}

	m.endTimes = append(m.endTimes, time.Now())

	return len(p), nil
}

func TestSafeConsoleWriter_CallsDefaultWriter(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConsoleWriter)

	m := mockIoWriter{}
	consoleWriter = &m

	w := newSafeConsoleWriter()
	expected := []byte{1, 2}
	w.Write(expected)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))
	assert.Equal(expected, m.data[0])
}

func TestSafeConsoleWriter_DoesNotWriteConcurrently(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetConsoleWriter)

	m := mockIoWriter{
		writeDelay: 2 * time.Second,
	}
	consoleWriter = &m

	w := newSafeConsoleWriter()

	data1 := []byte{1, 2}
	data2 := []byte{3, 4}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		w.Write(data1)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(1 * time.Second)
		w.Write(data2)
	}()

	wg.Wait()

	assert.Equal(2, m.called)
	assert.Equal(2, len(m.data))
	assert.Equal(data2, m.data[0])
	assert.Equal(data1, m.data[1])

	assert.True(m.endTimes[0].Before(m.startTimes[1]))
}

func resetConsoleWriter() {
	consoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
}
