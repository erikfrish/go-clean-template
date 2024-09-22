package common

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

func GetFuncName() string {
	var buffer bytes.Buffer
	const pcSize = 10
	pc := make([]uintptr, pcSize)
	const skip = 4
	runtime.Callers(skip, pc)
	frame, _ := runtime.CallersFrames(pc).Next()
	function := frame.Function
	line := frame.Line
	buffer.WriteString(function)
	buffer.WriteString(fmt.Sprintf(":%d", line))

	return filepath.Base(buffer.String())
}

type GeneralOpts struct {
	AppVersion string
	InstanceID uuid.UUID
	Env        string
	AppName    string
}
