package goflying

import (
	"io"
	"log"
	"os"
)

var Debugging bool

type logger interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
}

var Logger logger = log.New(io.Discard, "", log.LstdFlags)

func init() {
	debug := os.Getenv("GOFLYING_DEBUG")

	if debug == "1" {
		Debugging = true
		Logger = log.Default()
	}
}

func Debugf(format string, v ...any) {
	if !Debugging {
		return
	}

	Logger.Printf(format, v...)
}
func Debug(v ...any) {
	if !Debugging {
		return
	}

	Logger.Print(v...)
}
func Debugln(v ...any) {
	if !Debugging {
		return
	}

	Logger.Println(v...)
}
