package server

import "log"

type logger struct {
	v bool
}

func newLogger(v bool) *logger {
	return &logger{v}
}

func (l *logger) Fatal(msg ...interface{}) {
	log.Fatal(append([]interface{}{"[ERROR]"}, msg...)...)
}

func (l *logger) Error(msg ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, msg...)...)
}

func (l *logger) Log(msg ...interface{}) {
	log.Println(msg...)
}
