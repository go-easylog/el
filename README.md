# go-easylog

## What's this ?

This library is easy logger.  
It has these functions.

* 5 levels of output level for log
    * FATAL
    * ERROR
    * WARN <- Default
    * INFO
    * TRACE
* Output file
* File rotation

## Install
```
go get github.com/go-easylog/el
```

## Exsample
```go
package main

import "github.com/go-easylog/el"

func main() {

	// Default output level is "WARN"
	el.Warn("Output warning")
	el.Trace("This comment is don't output to log")

	// Set level for log
	el.SetLogLevel(el.TRACE)

	el.Info("This will be outputed")

	if err := el.SetRotateLog("./%Y/%M/%D.log"); err != nil {
		panic(err)
	}

	el.Trace("this comment outputing to file")

	if err := el.SetRotateLog("./no-rotate.log"); err != nil {
		panic(err)
	}
	el.Error("this file is no rotate")

	el.Fatal("error exit") // <- Forced kill here

	el.Trace("this text is no outputed") // <- This is not output
}

```