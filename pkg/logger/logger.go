package logger

import "github.com/labstack/echo/v4"

func New(prefix string) echo.Logger {
	l := loggerImpl{
		prefix: prefix,
		out:    newSafeConsoleWriter(),
	}

	l.log = prettyLogger.Output(l.out)

	return &l
}
