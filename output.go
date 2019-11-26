package scrape

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goextension/log"
)

func copyCache(cache *Cache, msg *Content, sample bool, output string) (e error) {
	pid := filepath.Join(output, strings.ToUpper(msg.ID))
	e = copyFile(cache, msg.Image, filepath.Join(pid, "image"))
	if e != nil {
		return e
	}
	e = copyFile(cache, msg.Thumb, filepath.Join(pid, "thumb"))
	if e != nil {
		return e
	}
	for _, act := range msg.Actors {
		e = copyFile(cache, act.Image, filepath.Join(pid, ".actor", act.Name))
		if e != nil {
			return e
		}
	}
	if sample {
		for _, s := range msg.Sample {
			e = copyFile(cache, s.Image, filepath.Join(pid, ".sample", "sample"+"@"+strconv.Itoa(s.Index)))
			if e != nil {
				return e
			}
			e = copyFile(cache, s.Thumb, filepath.Join(pid, ".thumb", "thumb"+"@"+strconv.Itoa(s.Index)))
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func copyInfo(msg *Content, path string, name string) error {
	pid := filepath.Join(path, strings.ToUpper(msg.ID))
	inf := filepath.Join(pid, name)
	_ = os.MkdirAll(filepath.Dir(inf), os.ModePerm)
	info, e := os.Stat(inf)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}
	bytes, e := json.MarshalIndent(msg, "", " ")
	if e != nil {
		return e
	}
	return ioutil.WriteFile(inf, bytes, 0755)
}

func copyFile(cache *Cache, source, path string) error {
	if source == "" {
		return nil
	}
	reader, e := cache.Reader(source)
	if e != nil {
		return e
	}
	path = TrimEnd(path)
	if debug {
		log.Infow("copy", "dir", filepath.Dir(path), "path", path)
	}
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	info, e := os.Stat(path + Ext(source))
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if e == nil && info.Size() != 0 {
		return nil
	}

	file, e := os.OpenFile(path+Ext(source), os.O_SYNC|os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		return e
	}
	defer file.Close()
	written, e := io.Copy(file, reader)
	if e != nil {
		return e
	}
	_ = written
	return nil
}

func imageCache(cache *Cache, m *Content, sample bool) (e error) {
	path := make(chan string)
	go func(path chan<- string) {
		defer close(path)
		path <- m.Image
		path <- m.Thumb
		for _, act := range m.Actors {
			path <- act.Image
		}
		if sample {
			for _, s := range m.Sample {
				path <- s.Image
				path <- s.Thumb
			}
		}
	}(path)

	for p := range path {
		if p != "" {
			_, err := cache.Get(p)
			if err != nil && !os.IsExist(err) {
				log.Error(err)
			}
		}
	}
	return nil
}

// TrimEnd ...
func TrimEnd(source string) string {
	return strings.Split(source, "?")[0]
}

// Ext ...
func Ext(source string) string {
	ext := filepath.Ext(TrimEnd(source))
	if debug {
		log.Infow("ext", "source", source, "ext", ext)
	}
	return ext
}
