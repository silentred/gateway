# Echo logrus

## Usage

```
import (
    "github.com/silentred/echorus"
    elog "github.com/labstack/gommon/log"
    "github.com/silentred/kassadin/util/rotator"
)

// new a splitter
//spliter = rotator.NewDaySpliter()
spliter = rotator.NewSizeSpliter(uint64(limitSize))

// new a writer
writer = rotator.NewFileRotator("/tmp/logdir", "app_name", "log", spliter)

defaultLogger := echorus.NewLogger()
defaultLogger.SetFormat(echorus.TextFormat)
defaultLogger.SetOutput(writer)
defaultLogger.SetLevel(elog.DEBUG)

```