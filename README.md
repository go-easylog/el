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

### main.go
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
	// %Y <- year  (YYYY)
	// %M <- month (MM)
	// %D <- day   (DD)
	//
	// "./prefix-%Y-%M.log"  <---  It is possible to specify like this
	//                             This pattern in one month's rotation


	el.Trace("this comment outputing to file")

	
	// no rotate pattern
	if err := el.SetRotateLog("./no-rotate.log"); err != nil {
		panic(err)
	}

	el.Error("this file is no rotate")
	el.Fatal("error exit") // <- Forced kill here

	el.Trace("this text is no outputed") // <- This is not output
}
```

### terminal
```
$ go run main.go
2017-02-04 04:08:03 [WARN] Output warning
2017-02-04 04:08:03 [INFO] This will be outputed
2017-02-04 04:08:03 [FATAL] error exit
exit status 255
```

### 2017/02/04.log
```
2017-02-04 04:08:03 [TRACE] this comment outputing to file
```

### rotator_test.go
```
2017-02-04 04:08:03 [ERROR] this file is no rotate
```