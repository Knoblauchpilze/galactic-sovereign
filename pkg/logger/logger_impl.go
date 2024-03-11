package logger

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

type loggerImpl struct {
	header string
	prefix string
	out    io.Writer
	log    zerolog.Logger
}

func (l *loggerImpl) Level() log.Lvl {
	return fromZerologLevel(l.log.GetLevel())
}

func (l *loggerImpl) SetLevel(level log.Lvl) {
	l.log = l.log.Level(fromLogLevel(level))
}

func (l *loggerImpl) Prefix() string {
	return l.prefix
}

func (l *loggerImpl) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *loggerImpl) Output() io.Writer {
	return l.out
}

func (l *loggerImpl) SetOutput(w io.Writer) {
	l.out = w
	l.log = l.log.Output(l.out)
}

func (l *loggerImpl) SetHeader(header string) {
	l.header = header
}

func (l *loggerImpl) prependPrefixAndHeaderIfNeeded(in string) string {
	out := in

	if l.prefix != "" {
		out = fmt.Sprintf("[%s] %s", l.prefix, out)
	}
	if l.header != "" {
		out = fmt.Sprintf("[%s] %s", l.header, out)
	}

	return out
}

func (l *loggerImpl) addPrefixAndHeaderIfNeeded(event *zerolog.Event) {
	if l.header != "" {
		event.Str("header", l.header)
	}
	if l.prefix != "" {
		event.Str("id", l.prefix)
	}
}

func (l *loggerImpl) msgF(event *zerolog.Event, format string, args ...interface{}) {
	event.Timestamp()
	format = l.prependPrefixAndHeaderIfNeeded(format)
	event.Msgf(format, args...)
}

func (l *loggerImpl) json(event *zerolog.Event, data log.JSON) {
	event = event.Timestamp()
	l.addPrefixAndHeaderIfNeeded(event)

	for key, value := range data {
		raw, _ := json.Marshal(value)
		event.RawJSON(key, raw)
	}

	event.Send()
}

func (l *loggerImpl) fields(event *zerolog.Event, fields ...interface{}) {
	event.Timestamp()
	l.addPrefixAndHeaderIfNeeded(event)

	for _, data := range fields {
		event.Fields(data)
	}

	event.Send()
}

func (l *loggerImpl) Print(i ...interface{}) {
	l.log.Print(i...)
}

func (l *loggerImpl) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

func (l *loggerImpl) Printj(data log.JSON) {
	var fields []interface{}

	for key, value := range data {
		fields = append(fields, key)
		fields = append(fields, value)
	}

	l.Print(fields...)
}

func (l *loggerImpl) Debug(i ...interface{}) {
	l.fields(l.log.Debug(), i...)
}

func (l *loggerImpl) Debugf(format string, args ...interface{}) {
	l.msgF(l.log.Debug(), format, args...)
}

func (l *loggerImpl) Debugj(data log.JSON) {
	l.json(l.log.Debug(), data)
}

func (l *loggerImpl) Info(i ...interface{}) {
	l.fields(l.log.Info(), i...)
}

func (l *loggerImpl) Infof(format string, args ...interface{}) {
	l.msgF(l.log.Info(), format, args...)
}

func (l *loggerImpl) Infoj(data log.JSON) {
	l.json(l.log.Info(), data)
}

func (l *loggerImpl) Warn(i ...interface{}) {
	l.fields(l.log.Warn(), i...)
}

func (l *loggerImpl) Warnf(format string, args ...interface{}) {
	l.msgF(l.log.Warn(), format, args...)
}

func (l *loggerImpl) Warnj(data log.JSON) {
	l.json(l.log.Warn(), data)
}

func (l *loggerImpl) Error(i ...interface{}) {
	l.fields(l.log.Error(), i...)
}

func (l *loggerImpl) Errorf(format string, args ...interface{}) {
	l.msgF(l.log.Error(), format, args...)
}

func (l *loggerImpl) Errorj(data log.JSON) {
	l.json(l.log.Error(), data)
}

func (l *loggerImpl) Panic(i ...interface{}) {
	l.fields(l.log.Panic(), i...)
}

func (l *loggerImpl) Panicf(format string, args ...interface{}) {
	l.msgF(l.log.Panic(), format, args...)
}

func (l *loggerImpl) Panicj(data log.JSON) {
	l.json(l.log.Panic(), data)
}

func (l *loggerImpl) Fatal(i ...interface{}) {
	l.fields(l.log.Fatal(), i...)
}

func (l *loggerImpl) Fatalf(format string, args ...interface{}) {
	l.msgF(l.log.Fatal(), format, args...)
}

func (l *loggerImpl) Fatalj(data log.JSON) {
	l.json(l.log.Fatal(), data)
}
