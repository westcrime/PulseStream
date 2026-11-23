package logger

import "log"

type StdLogger struct{ *log.Logger }

func (l StdLogger) Errorf(format string, v ...any) { l.Printf("ERROR: "+format, v...) }
