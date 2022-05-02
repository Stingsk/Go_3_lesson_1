// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// migrations/000001_create_gauge_table.down.sql
// migrations/000001_create_gauge_table.up.sql
// migrations/000002_create_counter_table.down.sql
// migrations/000002_create_counter_table.up.sql
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _migrations000001_create_gauge_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4f\x2c\x4d\x4f\xb5\x06\x04\x00\x00\xff\xff\x69\xd9\x1c\x57\x1b\x00\x00\x00")

func migrations000001_create_gauge_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations000001_create_gauge_tableDownSql,
		"migrations/000001_create_gauge_table.down.sql",
	)
}

func migrations000001_create_gauge_tableDownSql() (*asset, error) {
	bytes, err := migrations000001_create_gauge_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/000001_create_gauge_table.down.sql", size: 27, mode: os.FileMode(420), modTime: time.Unix(1651085405, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations000001_create_gauge_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4f\x2c\x4d\x4f\xd5\xe0\x52\x50\x50\x50\xc8\x4d\x2d\x29\xca\x4c\x8e\xcf\x4c\x51\x08\x73\x0c\x72\xf6\x70\x0c\x52\xd0\x30\x35\xd0\x54\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x41\x56\x57\x96\x98\x53\x9a\xaa\xe0\xe2\x1f\x0a\x32\x34\x20\xc8\xd5\xd9\x33\xd8\xd3\xdf\x0f\xac\x42\xd3\x1a\x10\x00\x00\xff\xff\x89\x49\x3d\x5c\x72\x00\x00\x00")

func migrations000001_create_gauge_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations000001_create_gauge_tableUpSql,
		"migrations/000001_create_gauge_table.up.sql",
	)
}

func migrations000001_create_gauge_tableUpSql() (*asset, error) {
	bytes, err := migrations000001_create_gauge_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/000001_create_gauge_table.up.sql", size: 114, mode: os.FileMode(420), modTime: time.Unix(1651085428, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations000002_create_counter_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xce\x2f\xcd\x2b\x49\x2d\xb2\x06\x04\x00\x00\xff\xff\xbb\x29\x09\xf1\x1d\x00\x00\x00")

func migrations000002_create_counter_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations000002_create_counter_tableDownSql,
		"migrations/000002_create_counter_table.down.sql",
	)
}

func migrations000002_create_counter_tableDownSql() (*asset, error) {
	bytes, err := migrations000002_create_counter_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/000002_create_counter_table.down.sql", size: 29, mode: os.FileMode(420), modTime: time.Unix(1651085615, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations000002_create_counter_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xce\x2f\xcd\x2b\x49\x2d\xd2\xe0\x52\x50\x50\x50\xc8\x4d\x2d\x29\xca\x4c\x8e\xcf\x4c\x51\x08\x73\x0c\x72\xf6\x70\x0c\x52\xd0\x30\x35\xd0\x54\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x41\x56\x97\x92\x9a\x53\x92\xa8\xe0\xe4\xe9\xee\xe9\x17\x02\x16\xd7\xb4\x06\x04\x00\x00\xff\xff\xa6\xff\xd3\x31\x6a\x00\x00\x00")

func migrations000002_create_counter_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations000002_create_counter_tableUpSql,
		"migrations/000002_create_counter_table.up.sql",
	)
}

func migrations000002_create_counter_tableUpSql() (*asset, error) {
	bytes, err := migrations000002_create_counter_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/000002_create_counter_table.up.sql", size: 106, mode: os.FileMode(420), modTime: time.Unix(1651085615, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"migrations/000001_create_gauge_table.down.sql":   migrations000001_create_gauge_tableDownSql,
	"migrations/000001_create_gauge_table.up.sql":     migrations000001_create_gauge_tableUpSql,
	"migrations/000002_create_counter_table.down.sql": migrations000002_create_counter_tableDownSql,
	"migrations/000002_create_counter_table.up.sql":   migrations000002_create_counter_tableUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"migrations": &bintree{nil, map[string]*bintree{
		"000001_create_gauge_table.down.sql":   &bintree{migrations000001_create_gauge_tableDownSql, map[string]*bintree{}},
		"000001_create_gauge_table.up.sql":     &bintree{migrations000001_create_gauge_tableUpSql, map[string]*bintree{}},
		"000002_create_counter_table.down.sql": &bintree{migrations000002_create_counter_tableDownSql, map[string]*bintree{}},
		"000002_create_counter_table.up.sql":   &bintree{migrations000002_create_counter_tableUpSql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
