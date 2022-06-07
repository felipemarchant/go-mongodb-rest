package utils

import "log"

func LogFatal(value ...any) {
	log.Fatal(value)
}

func LogPrintln(value ...any) {
	log.Println(value)
}
