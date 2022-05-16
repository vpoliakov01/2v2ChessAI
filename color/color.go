package color

import "fmt"

type Color string

var (
	Reset      Color = "\033[0m"
	Red        Color = "\033[31m"
	Green      Color = "\033[32m"
	Yellow     Color = "\033[33m"
	Blue       Color = "\033[34m"
	Background Color = "\033[37m"
)

func init() {
	Reset = Background
	fmt.Print(Reset)
}
