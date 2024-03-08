package logger

import (
	"io"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

type loggerImpl struct {
	out io.Writer
	log zerolog.Logger
}

func New() echo.Logger {
	l := loggerImpl{
		out: newSafeConsoleWriter(),
	}

	l.log = prettyLogger.Output(l.out)

	return &l
}

func (l *loggerImpl) Output() io.Writer {
	return l.out
}

func (l *loggerImpl) Prefix() string {
	panic("unimplemented")
}

func (l *loggerImpl) SetHeader(h string) {
	panic("unimplemented")
}

func (l *loggerImpl) Level() log.Lvl {
	return fromZerologLevel(l.log.GetLevel())
}

func (l *loggerImpl) SetLevel(v log.Lvl) {
	l.log = l.log.Level(fromLogLevel(v))
}

func (l *loggerImpl) SetOutput(w io.Writer) {
	l.out = w
	l.log = l.log.Output(l.out)
}

func (l *loggerImpl) SetPrefix(p string) {
	panic("unimplemented")
}

func (l *loggerImpl) messageF(event *zerolog.Event, format string, args ...interface{}) {
	event.Timestamp().Msgf(format, args...)
}

func (l *loggerImpl) Print(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

func (l *loggerImpl) Printj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Debug(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Debugf(format string, args ...interface{}) {
	l.messageF(l.log.Debug(), format, args...)
}

func (l *loggerImpl) Debugj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Info(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Infof(format string, args ...interface{}) {
	l.messageF(l.log.Info(), format, args...)
}

func (l *loggerImpl) Infoj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Warn(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Warnf(format string, args ...interface{}) {
	l.messageF(l.log.Warn(), format, args...)
}

func (l *loggerImpl) Warnj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Error(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Errorf(format string, args ...interface{}) {
	l.messageF(l.log.Error(), format, args...)
}

func (l *loggerImpl) Errorj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Panic(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Panicf(format string, args ...interface{}) {
	l.messageF(l.log.Panic(), format, args...)
}

func (l *loggerImpl) Panicj(j log.JSON) {
	panic("unimplemented")
}

func (l *loggerImpl) Fatal(i ...interface{}) {
	panic("unimplemented")
}

func (l *loggerImpl) Fatalf(format string, args ...interface{}) {
	l.messageF(l.log.Fatal(), format, args...)
}

func (l *loggerImpl) Fatalj(j log.JSON) {
	panic("unimplemented")
}
