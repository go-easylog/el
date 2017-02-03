package el

import (
	"os"
	"testing"
)

func Test_Basic(t *testing.T) {

	if GetLogLevel() != WARN {
		t.Fatal("default of log level is not WARN")
	}
	if IsDebug() {
		t.Fatal("default of Debug is true")
	}

	SetLogLevel(TRACE)
	if GetLogLevel() != TRACE {
		t.Fatal("set log level")
	}

	SetDebug()
	if !IsDebug() {
		t.Fatal("set debug")
	}

	SetRelease()
	if IsDebug() {
		t.Fatal("set release")
	}

	if File() == nil {
		t.Fatal("file is nil")
	}

	if File() != os.Stdout {
		t.Fatal("default file is not Stdout")
	}
}
