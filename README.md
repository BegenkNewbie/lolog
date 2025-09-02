### Log Writer With Log Level
There are two main parts in `log` which are __Logger__ and __Writer__.

- __Logger__: Main actor that will decide where, how and whether it should write the logs or not based on the log level defined in each `Writer`.

- __Writer__: Decide where the logs passed from `Logger` should be written to. Is it to terminal, file or whatever this `Writer` will decide that.
  We already provide pre-defined `Writer` implementer namely `console`, `file`, and `container`.
---

#### Dependency:
* [Slog](https://pkg.go.dev/golang.org/x/exp/slog) for provides structured logging severity level
* [Zap](https://pkg.go.dev/go.uber.org/zap) for provides structured logging json object and severity level
---

#### Output Configuration:
- __console__: will be printed on your console
  ![Console](https://i.postimg.cc/QDwV4Wzr/Screenshot-2024-10-24-at-4-19-15-PM.png)
- __file__: will be printed on your log files
  ![File](https://i.postimg.cc/qvFrW3X7/Screenshot-2024-10-24-at-4-20-32-PM.png)
- __container__: will be printed on your console with json stdOut
  ![Container](https://i.postimg.cc/YpJk3QCn/Screenshot-2024-10-24-at-4-19-46-PM.png)
---

#### Level Configuration:
- __debug__: if you fill this level with a `debug`, then `info`, `warning`, and `error` are included
- __info__: if you fill this level with an `info`, then `warning`, and `error` are included
- __warning__: if you fill this level with a `warning`, then `error` are included
- __error__: if you fill this level with an `error`, only `error` level will be printed on your output

---

#### Log file setup:

if you're using `output.file` you must have extra configuration for:
`path` (path of log file), `prefix_name` ex: (`'err'-YourhostName.log`)

#### log setup for your local environment
```yaml 
log:
  output:
    - console
    # - file
    # - container
  level: 'debug'
  #file_setup:
    # path: "./logs/error/"
    # prefix_name: "err-" //request from devops to add prefix file name
```

#### log setup for development, release, and production environment
```yaml 
log:
  output:
    # - console
    - file
    # - container
  level: 'debug'
  file_setup:
    path: "./logs/"
    prefix_name: "err-" //request from devops to add prefix file name
```
---

### How to use
```go
package main

import (
  "context"
  "errors"
  "fmt"
  "gitlab.dataon.com/gophers/sf7-lib/v2/log"
  stdLog "log"
  "os"
  "os/signal"
  "syscall"
)

var ll log.Logger
var err = errors.New("and i'm error information")

func main() {
  flushLogs := initLogger()
  ll.Dbg("Helo I'am log without context")

  // put another data to context, and register log to the context
  // use this technique to get data in your middleware
  ctx := context.Background()
  dt := ll.With(
    log.String("trace_id", "data0N1"),
    log.String("path", "/testing-aja"),
    log.String("ip", "127.0.0.1"),
    log.String("everything", "what you want"),
  )
  newCtx := log.WithCtx(ctx, dt)

  // and now, call the log from context and write something with method Dbg, Inf, Wrn, and Err
  wr := log.FromCtx(newCtx)
  wr.Dbg("Hi.. I'm debug message") // <-- debug level
  // output: {"level":"DEBUG","time":"2024-10-24T14:59:23+07:00","msg":"Hi.. I'm debug message","app_name":"wp-professional-std","app_env":"dev","app_version":"v11","app_timezone":"Asia/Jakarta","trace_id":"data0N1","path":"/testing-aja","ip":"127.0.0.1","everything":"what you want"}

  wr.Inf("Hi.. I'm info message") // <-- info level
  // output: {"level":"INFO","time":"2024-10-24T14:59:23+07:00","msg":"Hi.. I'm info message","app_name":"wp-professional-std","app_env":"dev","app_version":"v11","app_timezone":"Asia/Jakarta","trace_id":"data0N1","path":"/testing-aja","ip":"127.0.0.1","everything":"what you want"}

  wr.Wrn("Hi.. I'm warning message") // <-- warning level
  // output: {"level":"WARN","time":"2024-10-24T14:59:23+07:00","msg":"Hi.. I'm warning message","app_name":"wp-professional-std","app_env":"dev","app_version":"v11","app_timezone":"Asia/Jakarta","trace_id":"data0N1","path":"/testing-aja","ip":"127.0.0.1","everything":"what you want"}

  wr.Err("Hi.. I'm error message", log.Error(err))
  // output: {"level":"ERROR","time":"2024-10-24T14:59:23+07:00","msg":"Hi.. I'm error message","app_name":"wp-professional-std","app_env":"dev","app_version":"v11","app_timezone":"Asia/Jakarta","trace_id":"data0N1","path":"/testing-aja","ip":"127.0.0.1","everything":"what you want","error":"and i'm error information\n\n\n/Users/Oky/Documents/dataon/wp-professional-std/main.go:46 => main.main"}

  // this is wrapper u can use.
  // basically, I create from 'any' type on logger
  wr.Dbg(
    "With all wrappers",
    log.Query("select * from s where a = ?", "b"),
    log.Request(map[string]interface{}{"request": "foo"}),
    log.Response(map[string]interface{}{"response": "bar"}),
    log.Header(map[string]interface{}{"header": "foo bar"}),
    log.Error(err),
    log.StackTrace(),
    log.Any("object", "interface"),
  )
  // output: {"level":"DEBUG","time":"2024-10-24T14:59:23+07:00","msg":"With all wrappers","app_name":"wp-professional-std","app_env":"dev","app_version":"v11","app_timezone":"Asia/Jakarta","trace_id":"data0N1","path":"/testing-aja","ip":"127.0.0.1","everything":"what you want","query":{"statement":"select * from s where a = ?","values":["b"]},"request":{"request":"foo"},"response":{"response":"bar"},"header":{"header":"foo bar"},"error":"and i'm error information\n\n\n/Users/Oky/Documents/dataon/wp-professional-std/aaaa.go:49 => main.main","stack_trace":[{"file":"/Users/Oky/Documents/dataon/wp-professional-std/aaaa.go:50","func":"main.main"}]}

  // more example:
  // wr := log.FromCtx(ctx)
  // wr.Dbg("repo FindAllUser", log.Request(payload), log.Query(query, values))
  // wr.Inf("file has been removed", log.Response(resp), log.Request(payload))
  // wr.Warn("there is no data to process", log.Err(err), log.Query(query, values))
  // wr.Err("failed to execute query", log.Err(err), log.StackTrace(), log.Query(query, values))

  // Create channel that listen to termination signal
  sig := make(chan os.Signal, 1)
  signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
  // Block the program until terminate signal received
  <-sig

  // Flushing any leftover logs
  stdLog.Println("Flushing logs...")
  flushLogs()
  stdLog.Println("server Exiting...")
}

func initLogger() func() {
  // prepare setup file for log file writer
  fileSetup := &log.FileSetup{
    PrefixFile: "err-",
    Path:       "./logs/",
  }

  // New logger
  wr, err := log.SetupLogWriter(
    log.WithLevel("debug"), 
    log.WithOutput([]string{"container", "file", "console"}, fileSetup),
    // log.WithSlogLogger(), // if u want to setup a logger with Slog
  )
  if err != nil {
    panic(fmt.Sprintf("Failed to initialize logger: %v", err))
  }

  // Give contextual data should exist on all subsequent logs
  wr = wr.With(
    log.String("app_name", "wp-professional-std"),
    log.String("app_env", "dev"),
    log.String("app_version", "v11"),
    log.String("app_timezone", "Asia/Jakarta"),
  )

  // Register logger so that it can be used by router / etc
  ll = wr

  return wr.Flush

}

```

