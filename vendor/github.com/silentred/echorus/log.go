package echorus

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/gommon/log"
)

var (
	lvlToLevel map[log.Lvl]logrus.Level
	levelToLvl map[logrus.Level]log.Lvl

	TextFormat = &TextFormatter{QuoteCharacter: "`"}
	JSONFormat = &logrus.JSONFormatter{}
)

func init() {
	lvlToLevel = map[log.Lvl]logrus.Level{
		log.DEBUG: logrus.DebugLevel,
		log.INFO:  logrus.InfoLevel,
		log.WARN:  logrus.WarnLevel,
		log.ERROR: logrus.ErrorLevel,
	}

	levelToLvl = map[logrus.Level]log.Lvl{
		logrus.DebugLevel: log.DEBUG,
		logrus.InfoLevel:  log.INFO,
		logrus.WarnLevel:  log.WARN,
		logrus.ErrorLevel: log.ERROR,
		logrus.FatalLevel: log.ERROR,
		logrus.PanicLevel: log.ERROR,
	}
}

type Echorus struct {
	prefix    string
	logger    *logrus.Logger
	level     log.Lvl
	output    io.Writer
	shortFile bool
}

func NewLogger() *Echorus {
	l := logrus.New()
	l.Formatter = TextFormat

	e := &Echorus{
		logger:    l,
		shortFile: true,
	}
	e.SetLevel(log.DEBUG)
	e.SetOutput(os.Stdout)
	return e
}

func (e *Echorus) SetFormat(fomat logrus.Formatter) {
	e.logger.Formatter = fomat
}

func (e *Echorus) WithFields(fields logrus.Fields) *logrus.Entry {
	return e.logger.WithFields(fields)
}

func (e *Echorus) LogrusLogger() *logrus.Logger {
	return e.logger
}

func (e *Echorus) Output() io.Writer {
	return e.output
}

func (e *Echorus) SetOutput(w io.Writer) {
	e.output = w
	e.logger.Out = e.output
}

func (e *Echorus) Prefix() string {
	return e.prefix
}

func (e *Echorus) SetPrefix(p string) {
	e.prefix = p
}

func (e *Echorus) Level() log.Lvl {
	e.level = levelToLvl[e.logger.Level]
	return e.level
}

func (e *Echorus) SetLevel(level log.Lvl) {
	e.level = level
	e.logger.Level = lvlToLevel[level]
}

func (e *Echorus) Debug(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Debug(i)
}

func (e *Echorus) Debugf(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Debugf(format, args...)
}

func (e *Echorus) Debugj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Debug()
}

func (e *Echorus) Info(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Info(i)
}

func (e *Echorus) Infof(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Infof(format, args...)
}

func (e *Echorus) Infoj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Info()
}

func (e *Echorus) Warn(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Warn(i)
}

func (e *Echorus) Warnf(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Warnf(format, args...)
}

func (e *Echorus) Warnj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Warn()
}

func (e *Echorus) Error(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Error(i)
}

func (e *Echorus) Errorf(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Errorf(format, args...)
}

func (e *Echorus) Errorj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Error()
}

func (e *Echorus) Fatal(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Fatal(i)
}

func (e *Echorus) Fatalf(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Fatalf(format, args...)
}

func (e *Echorus) Fatalj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Fatal()
}

func (e *Echorus) Panic(i ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Panic(i)
}

func (e *Echorus) Panicf(format string, args ...interface{}) {
	e.logger.WithFields(logrus.Fields(e.StaticFields(2))).Panicf(format, args...)
}

func (e *Echorus) Panicj(j log.JSON) {
	static := e.StaticFields(2)
	fields := logrus.Fields(e.MergeJSON(j, static))
	e.logger.WithFields(fields).Panic()
}

func (e *Echorus) Print(i ...interface{}) {
	e.Debug(i)
}

func (e *Echorus) Printf(format string, args ...interface{}) {
	e.Debugf(format, args...)
}

func (e *Echorus) Printj(j log.JSON) {
	e.Debugj(j)
}

func (e *Echorus) StaticFields(skip int) log.JSON {
	_, file, line, _ := runtime.Caller(skip)
	if e.shortFile {
		_, file = filepath.Split(file)
	}

	file = fmt.Sprintf("%s:%d", file, line)

	return log.JSON{
		"prefix": e.prefix,
		"file":   file,
	}
}

func (e *Echorus) MergeJSON(args ...log.JSON) log.JSON {
	res := log.JSON{}

	if len(args) == 0 {
		return res
	}

	for _, json := range args {
		for key, value := range json {
			res[key] = value
		}
	}
	return res
}
