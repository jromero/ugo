package ugo

import (
	"fmt"
	"os"
	"strings"

	"github.com/jromero/ugo/pkg/ugo/types"
)

var _ types.Logger = (*Logger)(nil)

type Level int

const (
	DEBUG = iota
	INFO
	ERROR
)

type Logger struct {
	Level       Level
	breadcrumbs []string
}

func (l *Logger) AddBreadcrumb(breadcrumb string) {
	l.breadcrumbs = append(l.breadcrumbs, breadcrumb)
}

func (l *Logger) PopBreadcrumb(breadcrumb string) {
	if len(l.breadcrumbs) > 0 {
		last := l.breadcrumbs[len(l.breadcrumbs)-1]
		if last == breadcrumb {
			l.breadcrumbs = l.breadcrumbs[:len(l.breadcrumbs)-1]
		}
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level <= DEBUG {
		fmt.Fprintf(os.Stdout, "[debug]"+l.formatBreadcrumbs()+format+"\n", v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.Level <= INFO {
		fmt.Fprintf(os.Stdout, "[info ]"+l.formatBreadcrumbs()+format+"\n", v...)
	}
}

func (l *Logger) Error(err error) {
	if l.Level <= ERROR {
		fmt.Fprintf(os.Stdout, "[ERROR]"+l.formatBreadcrumbs()+err.Error()+"\n")
	}
}

func (l *Logger) formatBreadcrumbs() string {
	if len(l.breadcrumbs) == 0 {
		return "[*] "
	}

	return "[*][" + strings.Join(l.breadcrumbs, "][") + "] "
}
