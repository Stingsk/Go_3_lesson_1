package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __000001_create_gauge_table_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4f\x2c\x4d\x4f\xb5\x06\x04\x00\x00\xff\xff\x69\xd9\x1c\x57\x1b\x00\x00\x00")

func _000001_create_gauge_table_down_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_gauge_table_down_sql,
		"000001_create_gauge_table.down.sql",
	)
}

var __000001_create_gauge_table_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4f\x2c\x4d\x4f\xd5\xe0\x52\x50\x50\x50\xc8\x4d\x2d\x29\xca\x4c\x8e\xcf\x4c\x51\x08\x73\x0c\x72\xf6\x70\x0c\x52\xd0\x30\x35\xd0\x54\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x41\x56\x57\x96\x98\x53\x9a\xaa\xe0\xe2\x1f\x0a\x32\x34\x20\xc8\xd5\xd9\x33\xd8\xd3\xdf\x8f\x4b\xd3\x9a\x0b\x10\x00\x00\xff\xff\x18\xfb\x83\x91\x6f\x00\x00\x00")

func _000001_create_gauge_table_up_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_gauge_table_up_sql,
		"000001_create_gauge_table.up.sql",
	)
}

var __000002_create_counter_table_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xce\x2f\xcd\x2b\x49\x2d\xb2\x06\x04\x00\x00\xff\xff\xbb\x29\x09\xf1\x1d\x00\x00\x00")

func _000002_create_counter_table_down_sql() ([]byte, error) {
	return bindata_read(
		__000002_create_counter_table_down_sql,
		"000002_create_counter_table.down.sql",
	)
}

var __000002_create_counter_table_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x48\xce\x2f\xcd\x2b\x49\x2d\xd2\xe0\x52\x50\x50\x50\xc8\x4d\x2d\x29\xca\x4c\x8e\xcf\x4c\x51\x08\x73\x0c\x72\xf6\x70\x0c\x52\xd0\x30\x35\xd0\x54\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x41\x56\x97\x92\x9a\x53\x92\xa8\xe0\xe4\xe9\xee\xe9\x17\xc2\xa5\x69\xcd\x05\x08\x00\x00\xff\xff\x29\xac\x73\x96\x67\x00\x00\x00")

func _000002_create_counter_table_up_sql() ([]byte, error) {
	return bindata_read(
		__000002_create_counter_table_up_sql,
		"000002_create_counter_table.up.sql",
	)
}

var _bindata_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func bindata_go() ([]byte, error) {
	return bindata_read(
		_bindata_go,
		"bindata.go",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"000001_create_gauge_table.down.sql":   _000001_create_gauge_table_down_sql,
	"000001_create_gauge_table.up.sql":     _000001_create_gauge_table_up_sql,
	"000002_create_counter_table.down.sql": _000002_create_counter_table_down_sql,
	"000002_create_counter_table.up.sql":   _000002_create_counter_table_up_sql,
	"bindata.go":                           bindata_go,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"000001_create_gauge_table.down.sql":   {_000001_create_gauge_table_down_sql, map[string]*_bintree_t{}},
	"000001_create_gauge_table.up.sql":     {_000001_create_gauge_table_up_sql, map[string]*_bintree_t{}},
	"000002_create_counter_table.down.sql": {_000002_create_counter_table_down_sql, map[string]*_bintree_t{}},
	"000002_create_counter_table.up.sql":   {_000002_create_counter_table_up_sql, map[string]*_bintree_t{}},
	"bindata.go":                           {bindata_go, map[string]*_bintree_t{}},
}}
