package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"bytes"

	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/silentred/glog"
	"github.com/silentred/gateway/guard"
	"github.com/silentred/gateway/reactor"
	"github.com/silentred/gateway/route"
)

var (
	routeRegex = regexp.MustCompile(`\s*Route:(.+) Prefix:(.+) Service:(.+) Strip:(.+) Targets\((.+),\s*(\d+)\)\s*`)
)

type routeArg struct {
	host   string
	prefix string
}

type svcArg struct {
	name         string
	strip        string
	targetHost   string
	targetWeight int
}

type FileBackend struct {
	file  string
	table *route.Table
	data  []byte
}

func NewFileBackend(file string, table *route.Table) *FileBackend {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	return &FileBackend{
		file:  file,
		table: table,
		data:  b,
	}
}

// Parse the whole file data
func (fb *FileBackend) Parse() error {
	var buf *bytes.Buffer
	var line string
	var err, parseErr error
	var rArg routeArg
	var sArg svcArg
	var end bool

	err = fb.table.Reset()
	if err != nil {
		return err
	}

	buf = bytes.NewBuffer(fb.data)
	for !end {
		line, err = buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		rArg, sArg, parseErr = parse(line)
		glog.Debugf("[Parse] line=%s route=%v svc=%v", line, rArg, sArg)

		if parseErr != nil {
			glog.Errorf("[Parse] line=%s route=%v svc=%v err=%v", line, rArg, sArg, parseErr)
		} else {
			r := route.NewRoute(rArg.host, rArg.prefix)
			t := route.NewTarget(sArg.targetHost, sArg.targetHost, sArg.targetWeight)
			s := route.NewService(sArg.name, "", sArg.strip, t, guard.DefaultGroup, reactor.DefaultGroup)
			glog.Debugf("[Parse] line=%s route=%v svc=%v guard.DefaultGroup=%s reactor.DefaultGroup=%s", line, r, s, guard.DefaultGroup, reactor.DefaultGroup)
			fb.table.Add(r, s)
		}

		if err == io.EOF || len(strings.Trim(line, " \n\t\r")) == 0 {
			end = true
		}
	}

	return nil
}

// Watch file change
func (fb *FileBackend) Watch() {
	var w *fsnotify.Watcher
	var err error
	var event fsnotify.Event
	w, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer w.Close()

	err = w.Add(fb.file)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case event = <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				glog.Infof("[Watch] event=%s", event.String())
				err = fb.LoadData()
				if err != nil {
					glog.Errorf("[Config] err=%v", err)
				} else {
					fb.Parse()
				}
			}
		case err = <-w.Errors:
			glog.Errorf("[Watch] err=%v", err)
		}
	}
}

func (fb *FileBackend) LoadData() error {
	f, err := os.Open(fb.file)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fb.data, err = ioutil.ReadAll(f)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

func parse(line string) (routeArg, svcArg, error) {
	var r routeArg
	var s svcArg
	var err error
	matches := routeRegex.FindAllStringSubmatch(line, -1)
	if len(matches) > 0 && len(matches[0]) == 7 {
		r.host = matches[0][1]
		r.prefix = matches[0][2]
		s.name = matches[0][3]
		s.strip = matches[0][4]
		s.targetHost = matches[0][5]
		s.targetWeight, err = strconv.Atoi(matches[0][6])
		if err != nil {
			return r, s, err
		}
		return r, s, nil
	}
	return r, s, fmt.Errorf("route format error: text=%s", line)
}
