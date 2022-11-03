/*
Copyright 2022 Manuel Bovo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kind

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	kindLog "sigs.k8s.io/kind/pkg/log"
)

// NewStdoutLogger returns a logr.Logger that prints to stdout.
func NewStdoutLogger() logr.Logger {
	return funcr.New(func(prefix, args string) {
		if prefix != "" {
			_ = fmt.Sprintf("%s: %s\n", prefix, args)
		} else {
			fmt.Println(args)
		}
	}, funcr.Options{})
}

func NewLoggerWrapper() LoggerWrapper {
	return LoggerWrapper{log: NewStdoutLogger()}
}

// LoggerWrapper implement a wrapper from kind logger to logr.Loggger used by controller runtime
type LoggerWrapper struct {
	log logr.Logger
}

// Implemente the Kind Logger interface calling the beneath log.Error
func (l LoggerWrapper) Error(message string) {
	l.log.Error(errors.New(strings.ToLower(message)), "")
}

func (l LoggerWrapper) Errorf(format string, args ...interface{}) {
	l.log.Error(errors.New("error"), fmt.Sprintf(format, args...))
}

func (l LoggerWrapper) Info(message string) {
	l.log.Info(message)
}

func (l LoggerWrapper) Infof(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l LoggerWrapper) Warn(message string) {
	l.log.Info(message)
}

func (l LoggerWrapper) Warnf(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l LoggerWrapper) Enabled() bool {
	return l.log.Enabled()
}

func (l LoggerWrapper) V(lvl kindLog.Level) (il kindLog.InfoLogger) {
	// newLog := LoggerWrapper{log: logr.Logger{}.V(int(lvl))}
	newLog := LoggerWrapper{log: l.log.V(int(lvl))}
	return newLog
}
