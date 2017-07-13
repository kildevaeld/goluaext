// Code generated by go-bindata.
// sources:
// prelude.lua
// DO NOT EDIT!

package goluaext

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _preludeLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x54\xb1\x6e\xdc\x3a\x10\xec\xf5\x15\xf3\x2a\xdf\x01\x92\xf0\x52\xb8\xb3\xdc\xa4\x49\x91\x22\x85\xbb\x20\x38\x50\xd4\xca\x5a\x9b\x22\x05\x92\xf2\xf9\x62\xf8\xdf\x83\xa5\xa8\x3b\xdd\x05\x48\x25\x6a\x35\x33\x9c\xdd\x1d\xa8\xa8\x2a\x68\xa3\x42\xa8\xcd\xac\xe4\xe5\xab\x1b\x27\x15\xb9\x35\x84\x23\xc7\x01\xdf\x67\x85\xfb\xfa\x0b\x76\xd6\x45\xdc\xd7\xff\xef\xeb\xa2\x9f\xad\x8e\xec\xec\x42\xdc\xb5\x2a\x50\x09\xb6\x1c\xf7\x05\x00\xe3\xb4\x32\xd0\x68\xf0\xf1\x09\x00\x55\x05\x05\x4b\xc7\x05\x0d\xb6\x21\x2a\xab\x49\xa0\xdc\x43\x54\x85\x0a\x65\x3b\xc4\xd3\x44\x49\x6e\x8f\xa6\xc1\xdd\x7a\xcf\x1d\xe2\x40\x56\x08\xc2\x11\x70\x03\x41\xe5\x8a\x1c\xd1\xc0\xb2\x91\x02\x99\x40\xdc\xdf\x4a\x45\xd5\x1a\xda\xe8\x54\x15\xdc\xec\xb7\xb6\x02\x14\xc2\xa0\x8c\x71\x47\x68\x37\x9d\xe0\x7a\x81\x2f\xea\x09\xf3\x5f\xbe\xaf\x77\x1e\x5c\xbe\x81\x2d\x26\xc5\x3e\xe4\x6b\x3a\x97\xbf\x03\xd0\x3f\xf9\x17\x1a\xbc\xe5\x0a\xd9\x2e\x9f\x74\x7d\xc8\x76\xd7\x06\xf2\xb7\xaa\x4a\xb7\x2d\x66\x8e\x6c\x0c\x5a\x4a\x95\x91\xa2\x4a\xee\xd3\xbd\xca\x18\x70\x0c\x70\xed\x0b\xe9\x18\xca\x4c\x4d\xc3\x1b\xe8\xb4\x30\x8d\x73\xaf\x98\x27\xa9\xb0\x17\x81\xc1\x75\x32\x78\x70\xac\x8b\xc5\xc4\x81\x6d\x47\xef\x68\xa0\x8b\x2c\x41\xef\x93\x0b\x04\x05\xed\x6c\x88\x7e\xd6\xd1\x79\x1c\x07\xd6\x03\xb4\xb2\xe2\x46\x2b\x63\xa8\x43\x7b\xc2\x43\xb2\x69\xd5\x48\x8f\xbb\x07\xe5\x9f\xc3\xe3\x66\xf3\x63\x4c\xab\x97\xc2\x18\xeb\xc3\x41\x68\x68\xb0\x2e\x73\x97\xb8\x87\xd8\x9a\x12\x75\x5d\x6f\x88\xae\x7d\x39\x33\x03\xc5\x73\xe3\x3b\xd7\xbe\x94\x7a\x9f\x03\x93\xf6\x7f\x93\x87\x84\x58\xc5\x24\x00\xc8\x1f\xab\x0a\xa3\x7a\x25\x84\xd9\xcb\x34\x95\xc4\xec\x84\x10\xe7\xbe\x47\xef\xdd\x78\xb3\x61\x49\x81\xe8\xb1\x32\xfc\x9b\xba\x75\xe1\xdc\x2f\x18\x99\xb2\x1c\xea\x5b\x0f\xe7\x62\x32\x72\x6e\xeb\xb2\xdd\xfc\xf0\x14\x67\x6f\xa5\xd1\x4d\x51\xd7\x39\xd2\xf2\xc8\x85\x70\x50\xdb\x91\x05\x32\x7d\x89\x57\xb1\xb8\x0a\xe7\x59\xa3\xc1\xf3\x76\x54\x82\x5c\x21\xc7\x81\x0d\x61\x44\xe7\x70\x49\x26\xf7\x42\x6a\x16\xb1\xd4\xc3\xea\x2a\xfa\x99\x36\x49\x95\xf5\xa1\xc1\xb8\x04\xf6\xaf\x20\x67\x52\xaf\xcc\x55\x8c\xaf\xf6\xa6\x4b\x8c\xcb\x3f\x21\xa3\x75\x21\xb0\xa2\xf8\xf6\xf4\xf4\x63\xdb\x9f\x9a\x78\xbf\x78\xcc\xc0\x2b\x99\x8f\xcf\x12\x1f\x67\x57\x97\xe8\x9e\x7f\x43\x79\x3e\x3d\x93\xe9\x0e\x92\xca\xfd\xa5\x87\x8d\xd5\x33\x7c\x52\x71\x28\xe1\x26\x79\x0b\xd7\xd8\x0d\x5e\x4d\xbc\x0b\xd1\xb3\x7d\xae\xe7\x69\x22\xbf\xdb\xe8\x97\xf8\x87\xc6\x76\x86\xeb\xf9\x73\x9f\x5a\xff\x13\x00\x00\xff\xff\xb1\x7b\x4b\xa5\x6e\x05\x00\x00")

func preludeLuaBytes() ([]byte, error) {
	return bindataRead(
		_preludeLua,
		"prelude.lua",
	)
}

func preludeLua() (*asset, error) {
	bytes, err := preludeLuaBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "prelude.lua", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"prelude.lua": preludeLua,
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
	"prelude.lua": &bintree{preludeLua, map[string]*bintree{}},
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

