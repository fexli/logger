package logcolor

import (
	"errors"
	"github.com/xo/terminfo"
	"os"
	"sync"
)

type WriterConsole struct {
	syncMutex  sync.Mutex
	std        *os.File
	fd         uintptr // handle to the console
	colorLevel terminfo.ColorLevel
}

const (
	resetCtr = "\x1b[0m"
	startCtr = "\x1b["
	endCtrl  = "m"
)

var (
	InvalidConsole = errors.New("invalid console")
	lf             = []byte("\n")
)

func Colorable(file *os.File) *WriterConsole {
	w := &WriterConsole{
		syncMutex: sync.Mutex{},
		std:       file,
	}
	w.fd = file.Fd()
	w.colorLevel = terminfo.ColorLevelNone
	if IsTerminal(w.fd) && EnableColor {
		w.colorLevel = colorLevel
	}
	return w
}

func (w *WriterConsole) EnableColor() {
	w.colorLevel = colorLevel
}
func (w *WriterConsole) DisableColor() {
	w.colorLevel = terminfo.ColorLevelNone
}

const InvalidHandle = ^uintptr(0)

func (w *WriterConsole) Write(text *LogTextCtx, sync bool) (bool, error) {
	if w == nil || w.fd == InvalidHandle {
		return false, InvalidConsole
	}
	if sync {
		w.syncMutex.Lock()
		defer w.syncMutex.Unlock()
	}
	if w.colorLevel != terminfo.ColorLevelNone {
		text.WriteBytes(w.std, text.Color)
	} else {
		text.WriteRawBytes(w.std)
	}
	return true, nil
}

func (w *WriterConsole) Println(text *LogTextCtx) {
	if w == nil || w.fd == InvalidHandle {
		return
	}
	w.syncMutex.Lock()
	defer w.syncMutex.Unlock()
	if w.colorLevel != terminfo.ColorLevelNone {
		text.WriteBytes(w.std, text.Color)
	} else {
		text.WriteRawBytes(w.std)
	}
	w.std.Write(lf)
}
