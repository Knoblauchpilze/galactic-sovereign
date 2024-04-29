package logger

import (
	"fmt"
)

type Color int

const (
	Red     Color = 31
	Green   Color = 32
	Yellow  Color = 33
	Blue    Color = 34
	Magenta Color = 35
	Cyan    Color = 36
	Gray    Color = 90
)

func FormatWithColor(data interface{}, color Color) string {
	return fmt.Sprintf("\033[1;%dm%v\033[0m", color, data)
}
