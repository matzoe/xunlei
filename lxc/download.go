package main

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/golang/glog"
	. "github.com/matzoe/xunlei/api"
	. "github.com/matzoe/xunlei/fetch"
)

var worker Fetcher = DefaultFetcher

func dl(uri, filename string, echo bool) error { //TODO: check file existence
	if len(M.Gid) == 0 {
		return errors.New("gdriveid missing.")
	}
	return worker.Fetch(uri, M.Gid, filename, echo)
}

func download(t *Task, filter string, echo, verify bool) error {
	if t.IsBt() {
		m, err := t.FillBtList()
		if err != nil {
			return err
		}
		for j, _ := range m.Record {
			if m.Record[j].Status == "2" {
				if ok, _ := regexp.MatchString(`(?i)`+filter, m.Record[j].FileName); ok {
					glog.V(2).Infoln("Downloading", m.Record[j].FileName, "...")
					if err = dl(m.Record[j].DownURL, filepath.Join(t.TaskName, m.Record[j].FileName), echo); err != nil {
						return err
					}
				} else {
					glog.V(3).Infof("Skip unselected task %s", m.Record[j].FileName)
				}
			} else {
				glog.V(2).Infof("Skip incompleted task %s", m.Record[j].FileName)
			}
		}
	} else {
		if len(t.LixianURL) == 0 {
			return errors.New("Target file not ready for downloading.")
		}
		glog.V(2).Infoln("Downloading", t.TaskName, "...")
		if err := dl(t.LixianURL, t.TaskName, echo); err != nil {
			return err
		}
	}
	if verify && !t.Verify(t.TaskName) {
		return errors.New("Verification failed.")
	}
	return nil
}
