package query

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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _template_genny_dummy_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8e\xcb\x6a\xf3\x30\x10\x85\xd7\x9e\xa7\x18\xb2\x08\x49\x16\xd1\xbf\xf8\x57\x85\x2e\xda\x42\x97\x2e\xb4\x7d\x01\x59\x1e\xdb\x22\xba\x55\x1a\x51\x1c\xe3\x77\x2f\xca\x05\x07\x0a\x25\x3b\x71\xe6\xe8\x3b\x5f\x90\xea\x20\x7b\xc2\xaf\x4c\x71\x04\x10\x02\x79\xd0\x09\x75\xc2\x36\x5b\x3b\x62\x4f\xce\x8d\xc8\x64\x83\x91\x4c\xd8\x69\x43\xa0\x6d\xf0\x91\x71\x03\xd5\xaa\xd7\x3c\xe4\x66\xaf\xbc\x15\x07\x79\xcc\xa2\x6b\xd2\x40\x26\x50\x14\x27\xa0\x68\x64\xa2\x15\x54\x0d\xde\xd9\xdc\x02\xf0\x18\x08\xdf\x29\x65\xc3\x9f\xe5\x59\xfb\x96\x6a\x69\xe9\x7a\xf1\xfe\x9c\x27\x8e\x59\x31\x4e\x50\xed\x9a\xfd\x8b\xb7\xd6\xbb\x52\x85\xf9\x52\x7c\x32\x5a\xa6\xf2\xf1\xaf\x66\x97\x9d\x42\xb2\x81\xc7\x65\x71\xb3\xc5\xdd\xcd\xfe\x04\x55\x24\xce\xd1\xe1\x7a\x49\xa7\x05\xf3\x80\xeb\xe2\x7e\x03\x9e\xe6\xb9\xb0\x95\x77\x89\xf1\x55\x93\x69\xeb\x6c\x51\x3b\xc6\x47\xfc\x77\xcd\x3f\xf4\x91\xde\xba\x67\x99\xb4\x3a\xcd\x9c\xcf\xff\x2f\x4a\x35\x7d\xdf\x21\xf4\x4b\x1c\x66\xf8\x09\x00\x00\xff\xff\x03\x7e\x6b\x0a\xd1\x01\x00\x00")

func template_genny_dummy_go() ([]byte, error) {
	return bindata_read(
		_template_genny_dummy_go,
		"../template/genny/dummy.go",
	)
}

var _template_genny_field_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x53\x4d\x8f\xda\x30\x10\x3d\xc7\xbf\x62\xe4\x93\x8d\x68\x72\xaf\x94\x4a\xa8\x5f\xaa\xd4\xe5\x50\xe8\xa1\xa7\x55\x3e\x26\x60\x11\xdb\xa9\x63\xd3\x4d\x57\xfc\xf7\xca\x0e\x89\xd9\x6d\x84\x5a\x0e\xe0\xf1\xbc\x79\xf3\x3c\x6f\xe8\x8a\xea\x54\x1c\x10\x7e\x3a\x34\x03\x21\x42\x76\xda\x58\x60\x24\xa1\x07\x61\x8f\xae\x4c\x2b\x2d\xb3\xea\x88\x78\x1a\x4a\x61\xfb\xec\x80\x4a\x0d\xfe\x1b\x8d\xa8\xe8\x4b\xd8\xa9\xf8\xed\xb2\xa6\xec\x8f\xd8\x76\x68\xb2\x40\x99\x95\x45\x8f\x94\x70\x42\xb2\x15\x09\xc5\x20\x5d\x6f\xa1\x44\xa8\x8a\xb6\xc5\x1a\x3a\x34\xf0\x49\x60\x5b\x93\x55\x46\x88\x1d\x3a\x1c\xc3\x6d\x21\x11\xae\x8d\xd2\xfd\xd0\xe1\x6d\xd2\xc7\xaf\x92\xe7\xc2\x78\xdd\x5b\x5d\xa3\x2f\x7d\x9c\x49\xae\x27\x27\x41\x28\x0b\x39\x78\x45\xe9\xc6\x6a\xb1\xd5\x1f\x8d\x61\xfe\xc4\xe8\x84\xa1\x9c\x2f\x71\xc0\xfc\x19\x39\xee\x74\x09\x6f\xcd\x80\x35\x3e\x06\xa1\xb0\x7e\x5a\xc3\x18\x78\xfd\x1c\xde\xbc\x83\x58\xff\xa5\x7e\xda\xeb\xf0\x00\xaf\xff\xc3\xf7\x87\x87\x1f\x8f\x0b\xed\x4b\xad\x5b\xc8\x61\x87\x76\x4a\x86\x5c\xcf\xe8\x14\xd3\x35\xd0\x19\x3f\x07\x9e\x99\xae\xef\xc9\xe5\x84\x34\x4e\x55\xc0\x94\xae\x71\x06\xf2\xe8\x01\xe3\xc0\x0c\xf6\xae\xb5\xb0\xfa\x16\x7e\xf7\xe1\x19\xcf\x24\xb9\x5e\xe7\x80\xb2\xb3\x43\x4c\x32\x4e\x92\x4a\x4b\xa9\x15\xbc\xcd\xc1\x13\xa7\x81\x6e\x63\xd9\x5d\x25\x89\x68\x60\xac\x4b\x3d\x0e\xf2\x1c\x94\x68\x7d\xa7\xd8\x6a\x8b\xbf\x5e\x36\x4a\x02\xff\x0e\xed\x3f\xb4\x58\xc3\xc8\x93\xee\xb0\x6d\x36\xfd\xfb\xd0\xcb\xe3\x99\xb7\x7d\xd2\xfc\x5f\x92\x2f\x64\x9a\x43\x1a\x9c\xca\xe7\x17\x14\x12\x63\x4a\xd7\xf8\x55\xf4\xf6\x26\x7d\xbd\x99\x21\xf3\x26\x44\x4c\x5c\x8e\xbf\x40\x9f\x8d\x76\xdd\x02\x32\xdc\x07\x49\xd6\x19\x45\x2e\x7e\x17\x17\xfd\x9d\x06\x16\x2c\x3e\xc3\x2a\xfc\x2b\xe2\x40\x38\xa0\x31\xda\xc0\xf3\xf2\x76\xbc\xae\xbe\x5d\x8c\xa9\x70\x12\x31\x4e\x33\x52\xdf\x7a\x15\x8d\x39\x2f\x7a\x72\x21\x7f\x02\x00\x00\xff\xff\x1e\x42\x12\xfd\xa1\x04\x00\x00")

func template_genny_field_go() ([]byte, error) {
	return bindata_read(
		_template_genny_field_go,
		"../template/genny/field.go",
	)
}

var _template_genny_list_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x56\xc1\x6e\xe3\x36\x10\x3d\x9b\x5f\x31\xf5\x21\x20\x0d\x45\x4a\xae\x46\x75\x48\x83\x16\x08\x10\xb8\x40\x53\xa4\x87\x20\x08\x68\x99\xb2\xd9\x50\x94\x4a\x51\x49\x5d\x43\xff\x5e\x8c\x28\x59\x94\xb3\xb2\xbd\xbb\x39\xed\x21\x80\x42\x0e\xe7\xbd\x79\xf3\x66\x92\x82\x27\xaf\x7c\x2d\xe0\x9f\x4a\x98\x2d\x21\x32\x2b\x72\x63\x61\xba\x96\x76\x53\x2d\xc3\x24\xcf\xa2\x57\xfe\x5f\x15\xa5\xcb\x72\x23\x54\x21\x4c\xd4\x04\x46\x4b\x5e\x8a\x29\x21\x76\x5b\x08\xb8\x97\xa5\xfd\x13\x3f\x4a\x6b\xaa\xc4\xc2\x0e\xa2\x08\xd6\x42\xeb\x2d\x99\xcc\x6e\xf3\x2c\xcb\xf5\x22\x5f\x09\x52\x13\x12\x45\x80\x9f\x0b\x9e\x89\x36\x22\xad\x74\x02\x0b\xf1\xde\x25\xa1\x0c\x66\xfb\x84\x3b\x42\x26\x4a\x96\x16\xe6\x31\x88\xac\xb0\xdb\x3e\xca\x5d\x84\x98\x0d\x0f\x21\x86\x0b\xe4\xb4\x3f\xd8\xd5\x6d\x44\x4f\x20\x6c\x60\x63\x98\x3e\x3d\x77\x24\xa6\x84\x4c\xe8\xac\x79\x88\x8f\x18\x3d\x78\xc2\xc2\x3b\x2d\x2d\x5e\x21\xa2\x11\xb6\x32\x1a\x30\x06\x8b\x69\xa8\x1f\xd0\x1a\x90\xef\x1e\x5c\x74\x67\xbb\x3e\xf3\xbc\xe5\xdb\x9f\xec\xea\x7a\x9f\x95\xea\x7c\xd5\x0b\xcb\xe0\xc6\x52\x09\x52\x5b\x06\xd4\x88\xb2\x52\x16\x66\x5d\x09\x01\x08\x10\xc6\xe4\x86\x39\xc4\xe6\x36\x86\x8b\xee\x1e\x85\x70\xa7\x1e\x16\x3e\x8a\x61\x50\x39\x22\x0e\x2a\x47\xcc\x7d\xd1\x47\x98\xfd\x25\x37\x79\x65\x7f\x45\x0a\xa3\x24\x3d\x6e\x08\x3d\x8f\xa1\x81\x6b\x21\x64\x0a\x02\x7e\x8a\x41\x4b\x85\x71\x7d\x11\x5a\x2a\x32\xa9\x4f\x72\x78\x10\xb6\x13\x28\x80\x37\x1f\xb5\x11\xc6\xeb\xc4\xf1\x8a\xdb\x3c\x01\xbc\xf9\xc7\xa4\x1e\xa9\x7c\xb5\xa2\x6f\xde\xef\x5f\x87\x85\xaf\x07\x11\x03\x4c\x36\x5a\xeb\x1f\x5c\xaf\x05\x2d\x2d\x37\x36\x00\xc5\x4b\xeb\x14\x1f\xd8\x4e\xa1\xc0\xc7\xe1\x17\xe2\x9d\x36\xf7\xbf\x17\xf6\x43\x4a\xe6\x7a\xa2\x20\xf6\x7b\xd2\x14\x35\xec\xc8\x88\xb3\x5b\x68\x0f\x8f\x2a\x36\x6e\xee\xdf\xa4\xc1\xf9\x3a\x6d\x6d\xc7\xa0\x35\xce\xd5\xb8\x46\xf7\xfc\x1b\xf2\x49\x6d\x9d\x4c\xdd\x06\x09\x1f\xb9\xaa\xc4\x9d\x4e\xf3\xf0\xf1\x5e\x68\x06\x97\x70\x3d\x8e\xf9\x20\x94\x48\x2c\x4d\x35\xe0\x35\xf5\x3c\xb8\xcc\x73\xd5\x93\x79\x7a\xfe\xd2\x50\x40\x0c\x19\x7f\x15\xd4\xbb\x0d\xe0\x2a\x80\x93\xa4\x18\x99\x24\x8d\xcc\xe5\xe9\x9e\x77\x14\x91\x5f\x92\x81\xb7\x9a\x1d\x49\xbf\xcf\xa9\xa6\xfd\x06\xf1\x7b\x9b\x64\x35\x23\x13\xfc\x49\x73\x03\x2f\x01\x24\x19\x22\x1b\xf4\x10\x74\x54\x06\x43\xcc\x8b\x42\xe8\x15\xed\xa6\xff\x68\xda\x7d\x57\x5c\xf4\x11\xcb\xe8\xd5\xb8\xd6\xfb\x03\x4f\xe0\xf9\xa9\x95\x78\x72\x1f\x3a\xcc\x4f\x12\xef\xbc\x32\x6f\x94\xa2\xcc\xb7\xcc\xa1\x6f\xfd\x96\x1e\xc8\x00\x3b\x68\x23\xad\xa9\x04\xd4\xe3\xd6\xbd\xcd\x2b\x8d\xf3\x22\xb5\xf5\xf2\x9f\xb4\xde\xf8\x28\xbc\xf3\xc2\xed\xd1\xbf\xdd\x72\xfa\xb8\x19\xc7\x1c\xda\xbf\x3c\x92\x3e\x37\xf6\x97\x2d\x55\xa2\x2c\x5d\xff\x7b\x20\xd7\xfe\xf3\xe1\xfa\x4c\xac\xfd\xef\xe4\x41\x70\x93\x6c\x20\x0c\x43\x58\x4a\xcd\xcd\x16\xca\xe6\x64\x6c\xe8\xf1\xee\xbb\x8c\x48\x26\xf2\x9c\xc9\x45\x9c\x3b\xbd\x12\xff\x36\x7b\xea\x78\x34\x76\x87\x32\x16\xc0\x67\x99\x55\xa6\x20\xe1\x67\x38\x1b\xd9\x9b\xff\x00\x5e\x60\xf8\xc7\x1e\x4b\x3e\xcf\xfd\x7e\xd5\xa3\x12\x3b\xd3\xfe\xc0\x2a\xba\x61\x1c\xe8\x76\x79\x4d\x6a\xf2\x7f\x00\x00\x00\xff\xff\xb6\xce\xe0\x86\xb9\x0b\x00\x00")

func template_genny_list_go() ([]byte, error) {
	return bindata_read(
		_template_genny_list_go,
		"../template/genny/list.go",
	)
}

var _template_genny_node_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x56\xdf\x6f\xdb\x36\x10\x7e\x16\xff\x8a\xab\x1f\x0a\x31\x35\xec\xee\x35\x9d\x31\x64\xa9\x5d\x18\x48\xed\xc1\x71\xd7\x01\x41\x60\x28\x32\xe5\x70\x91\x48\x85\xa4\xb2\x64\xae\xfe\xf7\xe1\x48\xfd\xa0\x6c\x39\xcd\x5e\xa2\xf8\xf8\xdd\x77\x1f\xef\x78\x47\xe6\x51\xfc\x10\xed\x18\x3c\x16\x4c\xbd\x10\xc2\xb3\x5c\x2a\x03\x21\x09\x06\x3b\x6e\xee\x8b\xbb\x51\x2c\xb3\xf1\x43\xf4\x6f\x31\x4e\xee\xf4\x3d\x4b\x73\xa6\xc6\x16\x3b\xbe\x8b\x34\x1b\xbc\x01\x97\xca\xdd\x80\x50\x42\xc6\x67\x04\x00\xb2\x42\x1b\x88\xa3\x34\x85\x5f\xc0\xf0\x8c\x69\xc8\x99\x82\x75\x74\x97\xb2\xb1\x36\xaa\x88\x4d\xb8\x90\x5b\xb6\x88\x32\x46\xc9\xd9\x98\x10\xf3\x92\x33\xa8\x4d\xe0\x20\xb0\x27\xc1\x19\xc6\x1f\x5d\xca\x2c\x93\x02\x97\x49\x49\x48\x52\x88\x18\x58\x96\x9b\x97\xda\x21\xa4\x70\xd6\x38\xef\x49\xa0\x98\x29\x94\x80\xf7\xb5\x6d\xdf\x12\x9c\xc3\xfb\x03\xca\x7d\x59\x22\xeb\x53\xa4\x1a\x01\x9b\xf9\xf6\x79\x2d\xd7\xa8\x29\x8b\xf2\x1b\x2e\xcc\x2d\x17\x06\x26\xfe\xaf\x7d\x79\xc2\xe5\x8b\x92\x45\xfe\x3f\xfd\xac\xf0\x1a\xa4\x8d\xe2\x62\xe7\x79\x39\xc3\xa1\x23\xfe\x59\xcb\xf9\xf6\xd9\xe2\x1c\xc6\x0b\xd7\x1a\xf6\xd5\xee\x3e\x7f\xfb\xfa\xf5\x8f\x4d\xed\x3f\xd7\xd7\x2e\xcb\x77\x52\xa6\x30\x01\x9b\x95\x6b\x66\xda\x35\x11\x9b\x70\x50\xc3\x07\x43\x87\x58\xcb\xdf\xa5\x4c\xc3\x41\xed\x3e\xa0\xb4\xaa\x08\xfa\x56\xe0\x19\x67\xe9\x56\x87\x02\xff\x1f\x42\x52\x7d\x6c\x3e\x9d\x2c\x34\x16\x19\x70\x61\xa8\x8b\xbf\x27\x24\xb0\xfc\x2b\xf6\x58\x30\x6d\xae\x99\x31\x5c\xec\x5e\x27\x73\x2c\x94\x90\x80\x89\x22\x9b\x19\x6b\x92\x0f\x70\x5e\xed\xc6\x65\x08\x91\x37\x16\x7f\x4b\x02\x9e\x20\x60\x4f\x82\xa0\xd6\x7a\xcd\x4c\x53\xba\x10\xf9\x86\xd0\xb0\x51\x12\x94\x24\x18\x8f\x67\xf3\xbf\xbe\x4e\xcf\x91\x94\xc7\x60\x8f\xaa\x14\xe9\x0b\x68\x23\x15\xfb\x8d\x90\x23\x2e\x7b\x26\x1d\x97\x0d\x8c\x12\x77\x2a\x47\x61\x35\xf4\x0b\x33\xcd\x61\x09\x2b\x90\xcf\xe3\x2d\x5a\x9e\x9d\xca\xa9\x17\xa9\x3d\x36\x37\xb8\x7e\x0b\x13\x17\xc9\x87\x34\x07\xe4\xc6\x26\xcd\x62\x16\x45\x56\x25\xba\x2a\xb5\x05\x74\xca\x7c\xec\x8f\x81\xab\x96\x32\xaa\x60\x84\x94\xae\xe2\xbd\xdb\xe6\x58\xd5\x21\xe8\xaa\xd2\x14\x93\xdd\xa3\x9a\xa3\x1c\x4d\x5a\x31\x2d\x43\xaf\x98\x66\x99\x36\x23\xa0\xb7\x82\x0f\x43\x78\xea\x9c\xab\xe0\xb8\x43\x6f\x1e\x30\xf8\xd3\x41\x6c\xeb\x7e\x3a\x76\x55\x22\x3f\x0f\x3d\x42\xda\xb2\xbd\x45\x88\x45\x1e\xa9\xe9\xac\xfe\x44\x91\xc5\xfc\x44\x56\xe7\xa8\xb5\x75\x09\x15\xd3\x45\x6a\x9c\xc6\x76\x68\x5a\x15\x5d\x97\x36\xe5\xa1\x90\xdb\x76\x4e\x53\x88\x9b\x29\x8a\x23\xf8\x60\xb0\x22\x29\x4f\x00\x5d\x7c\xeb\x64\x02\x82\xdb\x84\x04\xa9\xdc\x8d\xae\xe4\x2e\xb4\xdf\xe5\x97\xcd\xf7\x8b\xd5\x62\x08\x18\x29\xa4\x50\x2d\x5e\xa8\x9d\xb6\xe0\x5a\x20\xda\x67\xe1\xc0\x63\x14\xd2\x40\x22\x0b\xb1\x6d\x94\x0d\x28\x09\x82\x12\xdb\x17\x58\xaa\x19\xf0\x04\x52\x26\xc2\x03\x29\x23\xb7\x89\xc9\x04\x3e\xc2\x8f\x1f\xbd\x88\xb6\xf4\x0e\x86\x42\xfa\x58\x60\x02\x6d\xa1\x7a\x30\xed\x85\x32\xe9\xa9\xe2\x6b\x0e\xee\x3a\xe9\xf3\xb2\x2b\x76\x42\x55\x99\x39\xe0\xa8\xdb\xf4\xb0\x68\xd7\x2c\x52\xf1\xfd\x5c\x24\x32\xcc\xa5\x76\xed\x9a\x54\x85\x5f\xb1\x78\x26\x86\x10\x4b\xb1\x9d\x55\xa6\x4b\xfb\x3f\xb5\x23\xda\x46\xf0\x8b\x3e\xea\x92\x21\x51\xed\x4d\xc9\xc9\x63\x63\xe1\x14\x42\x2e\x12\xe9\x82\xa0\xc5\x85\xf0\xf7\xd2\x89\xe4\x9c\x5e\x21\xd5\x57\x2c\x4a\x2e\x4c\xf8\xf7\xc1\xbd\x72\x9a\xb1\xf1\x78\x85\xf6\x52\x16\xc2\x2c\x13\x7b\x0b\x85\x14\xa9\xbd\x66\xc1\x23\xd3\x37\x26\x4e\x91\xfd\x19\xa5\x05\xb3\x1b\xe1\x95\x48\xdc\x7c\x63\xf5\x98\x8f\xd5\x7a\xbe\xa7\x03\x58\x99\x17\x26\xe4\xdb\x67\x17\xa0\xaf\x29\x4f\x86\xf0\xbc\xbd\x10\x5a\xc5\x5e\x84\xe9\x63\x11\xa5\xe1\x56\x1b\xcf\x56\x4f\xb8\x44\x2a\xe0\x78\xcb\x7d\xfc\x04\x1c\x7e\x05\xad\xe2\x51\x37\x7f\x9f\x80\x7f\xf8\x60\xbb\x88\x27\xf0\x0e\xd7\x9b\x98\x74\xd4\x50\x7b\x46\xda\xe9\xfd\x24\x4a\x35\x76\x4b\xe9\x1f\xfb\x7a\xe8\x75\xde\x90\xdf\xb9\xb9\x9f\x2a\xe5\x3f\x25\xeb\x25\x12\xe0\x02\x53\x4a\xaa\xe3\x11\xce\xc5\x2e\x65\x2e\xab\x8d\xc3\x10\x98\x83\xd3\x23\xfa\x36\x9b\x07\x2b\xfb\xfa\xf7\xb9\xcd\xf3\x10\xa6\x4a\x9d\x03\x2b\xdb\x88\xec\x9f\x53\xcf\x57\x5b\xb2\x8b\x3c\x4f\x5f\xaa\xd7\x90\xf7\x0c\xa2\xae\x09\x31\xc9\x07\x4f\xe0\xaa\x3b\xd1\x70\xc5\x35\x3e\x03\xdd\x4b\xb7\x36\xec\x4b\xf2\x96\xd1\x65\x21\x73\xc1\x4d\x48\xbb\xbd\xf3\xfa\x99\xb3\xc3\x48\x87\xb4\xf3\xf0\x3d\xce\xce\xe1\xf4\x3a\xc5\xb9\x92\xd2\xe0\x84\xc0\xaf\x7b\xea\x55\x05\x70\x77\xca\xbb\x4a\x65\x05\xc3\x23\x52\x05\xaa\x1d\xf6\xe5\xd0\xde\x14\xd3\xd5\x6a\xb3\x58\x6e\xe6\x8b\xcb\xab\x6f\x9f\xa7\x9b\xd5\x72\xb9\x76\x67\x47\x4a\x83\x49\x34\xd2\x72\x38\xbe\x65\x7b\x93\xe2\xfa\x10\x2f\x2a\x52\x92\xff\x02\x00\x00\xff\xff\xa0\xb8\xf0\xa1\x3e\x0d\x00\x00")

func template_genny_node_go() ([]byte, error) {
	return bindata_read(
		_template_genny_node_go,
		"../template/genny/node.go",
	)
}

var _template_genny_node_test_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\xcf\x6e\xf2\x30\x10\xc4\xcf\xde\xa7\x58\xe5\x80\xe2\xef\x8b\xe2\x9e\xab\xd2\x23\x12\x17\x84\x04\xf7\xca\x4e\x36\xc4\x22\xb1\x5d\xff\x39\x00\xe2\xdd\x2b\x9b\x50\xa4\xde\x3c\x1e\xff\x66\xd6\xeb\x64\x77\x96\x27\xc2\xef\x44\xfe\x02\xa0\x67\x67\x7d\xc4\x1a\x58\x15\x29\x44\x6d\x4e\x15\x00\x53\x58\x9d\x74\x1c\x93\x6a\x3b\x3b\x8b\xb3\xbc\x26\x31\xa8\x30\xd2\xe4\xc8\x8b\x42\x0a\x25\x03\x55\xc0\x01\x86\x64\x3a\x3c\x52\x88\x5f\x81\xa4\xef\xc6\xad\x19\x6c\x1d\xf1\xdf\x12\xd7\x1e\x39\xde\x00\x98\xc4\xf7\x35\xd2\xec\xe2\x65\x67\x7b\xda\xc9\x99\x6a\x0e\x4c\x88\xa7\xda\x68\x9a\xfa\xb0\x4b\xf3\xd6\xe8\x98\xad\xce\x9a\x3e\x33\x39\xbf\x76\x36\xa0\x36\xb1\x41\x6d\x06\x8b\xaa\xcd\x25\x1c\x95\xb5\x13\xde\x80\x31\x4f\x31\x79\x53\xcc\x76\x6f\x03\x7e\xac\x31\x13\xab\x15\xd6\xcf\xbb\xff\xe5\x70\xd0\x57\xe2\xf8\x99\x5d\x60\x4c\x88\x05\x8c\x3e\x11\xb0\x3b\x30\x4f\xdd\xc6\xfc\xb6\x06\x54\x6d\x9e\x6f\x2f\xe3\xf8\xa7\xfa\x56\xf0\xdc\x1c\xd2\x14\x71\x8d\xd2\x39\x32\x7d\xfd\xd0\x0d\x06\xbe\x3c\xc8\x54\x78\xf9\x45\x3e\xb2\x78\x69\x94\xed\xe1\xb5\xb6\xb7\x06\xcb\x04\x0d\xe6\xdf\x73\x80\x3b\xfc\x04\x00\x00\xff\xff\x93\x01\x9e\x3b\xb0\x01\x00\x00")

func template_genny_node_test_go() ([]byte, error) {
	return bindata_read(
		_template_genny_node_test_go,
		"../template/genny/node_test.go",
	)
}

var _template_genny_root_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x98\xcf\x6f\xdb\x36\x14\xc7\xcf\xe2\x5f\xc1\xea\x10\x90\xab\x21\xd7\x3f\x20\x04\x45\x7d\x58\x06\x64\x0d\x50\xd8\x43\xda\xa2\x87\x61\x07\xd9\xa1\x6c\x21\xb6\xa8\x49\x54\x52\x25\xf0\xff\x3e\x3c\x92\xa2\x28\xc9\xb1\x27\x40\x10\x7a\x29\xc4\xc7\xc7\xf7\xe3\x43\xf2\x6b\x36\x49\xb0\x79\x0c\xb6\x0c\xff\x9b\xb3\xb4\x40\x28\x3a\x24\x3c\x15\x98\x20\xc7\x8d\xb8\x8b\x1c\x37\x13\xe9\x86\xc7\x4f\x2e\x42\x4e\xb8\x0f\xc4\x3a\x0f\x43\x96\x66\xd8\xdd\x46\x62\x97\xaf\xbd\x0d\x3f\x8c\xb7\x9c\x6f\xf7\x6c\x6c\x4d\x8f\xb7\x72\xad\xe5\xf3\x18\xbc\xe4\xe3\x70\x9d\xed\xd8\x3e\x61\xe9\x58\x66\x1b\xaf\x83\x8c\xb9\x88\x22\x24\x8a\x84\xe1\x1b\xce\xf7\x78\x81\xc1\xe8\xc1\xb7\xb6\x16\x82\x19\x6b\x21\x98\xb2\xde\xc5\xe2\xba\xb4\xc2\xb7\xb1\x4e\x7c\xcb\x3c\xf1\x8d\x7d\x36\xb5\xec\xb3\xa9\xb1\xfb\x73\xcb\xee\xcf\x95\xfd\x7b\x64\x85\x97\x83\xca\x5e\x25\x50\xa3\x6a\xa6\x4a\xa1\x46\xd5\x4c\x95\x44\x8d\xd4\xcc\xed\x9e\x07\xd6\x22\x3d\xb4\xe6\xaa\x65\x7a\xa8\x41\xfd\xc1\x0f\x07\x1e\x2f\xf9\x83\x01\x53\x59\x94\xc7\x97\x28\x13\xe5\x1c\x7c\xb7\x16\xda\x0e\xd6\x62\x14\xe6\xf1\x06\xb3\x43\x22\x8a\xca\x4a\x28\xfe\xad\xe1\x88\x5f\x91\x93\x32\x91\xa7\x31\xbe\x6a\x4c\xbd\x1e\xd1\x51\xc7\x59\xb2\xe7\xff\x1b\xa5\x95\xd2\x04\x69\xcc\x40\xe1\x10\xaa\xd1\x8a\x55\x4f\x7d\xe6\xad\x72\x2e\xc6\x39\x99\xb7\x5e\x55\x19\xc3\x60\xb6\x96\x93\xca\x4a\x49\x9b\x11\xb5\xab\x2a\xe3\x9c\xaa\xa0\x91\x77\x95\xb0\x98\xa4\x38\xe2\xde\x3d\x0b\x1e\x58\x3a\xc2\x9b\x20\xc1\x51\x2c\x46\x98\x27\x22\xc3\x9e\xe7\xc9\x5c\xab\x44\x44\x3c\xa6\xf8\x9e\x73\xf1\x0d\xb6\x5e\x86\xcd\xf2\xbd\xc0\x1f\x17\xc6\xfa\x5a\x95\xf4\x11\xeb\x75\x10\x5f\x46\x55\x11\x3d\xcf\xa3\xc7\x72\xad\xd5\x82\xb7\x0c\x0e\x70\xfc\xdc\x32\x96\x8b\x1c\x19\xe1\x2b\x13\x60\x82\x69\x72\x7a\x15\x35\xe1\x6e\x99\xd8\xec\xee\xe2\x07\xf6\x93\x50\xd3\xb6\x9a\xab\x75\x7c\x53\xdc\xe4\x21\x59\xe7\x21\xfe\xfb\x9f\x75\x21\x58\x8f\xcd\x9a\xd0\xbf\x48\xbb\x82\x43\x38\xb2\xd6\x9a\xb4\xaa\x77\xb5\x91\xa1\xa1\xab\xab\xda\x8d\xd3\x13\x9e\x75\xb5\xd5\x99\xb3\xae\x81\xed\x53\x5e\xfd\x25\x7b\x96\xd7\x6d\x3d\x82\x33\x44\x2c\x05\xf7\xfe\x64\xe2\xfb\x2a\x0c\x33\x26\xbe\x91\xb5\x77\x4f\x3e\x50\x4a\xa1\x6a\xce\xdf\xa6\xaa\x32\x1c\x95\xd7\x05\x7a\xd2\xe5\x34\x11\xce\x2b\x1e\xb0\xe2\x36\xe5\x07\x15\x8c\x6c\x5a\x22\xd2\xdc\x76\x19\xc2\x54\x87\x1c\xa7\x56\xa0\x07\x33\x3a\x16\x1d\x21\xa7\x52\x07\x12\x03\x97\x72\x21\xc5\x4b\xf6\x13\x2e\xa5\x1d\x3c\x13\x41\x2a\x9b\x07\x57\xef\x0b\x8b\x09\x45\xc8\x89\x42\x33\x86\xa3\x44\xf1\x27\xac\x1c\x21\xb9\xae\x07\x1c\x20\x19\x72\x62\xf6\x7c\x13\x64\xcc\x44\xb9\x5b\x79\x32\x93\x5c\x02\xe1\x5a\x80\xdb\x38\xf1\xe2\x84\x5e\x2a\xa7\xb3\x07\xc0\x78\x34\xb7\x5f\x17\x75\xfe\x10\x68\x27\xeb\x28\x9c\xd8\xae\x26\xc5\xbb\x38\x63\xa9\x00\x2e\x09\xcf\x46\x38\x8b\x5e\x18\x24\xa1\x00\x47\xf6\x6f\x1d\x92\x13\xbe\xf4\xcd\xc0\x9f\x83\x4c\xef\xd0\x1a\x9e\x0d\xaf\x08\xd9\xac\xcd\x66\xbc\x9f\xe3\x4f\xf6\x76\xbd\x15\xee\x47\x24\x76\x9f\xa5\xa2\x36\xf6\x3c\x0a\xf1\x3b\xb5\x51\xb1\xbc\x98\xb2\x72\x67\x27\x5d\x61\x9b\x0e\xc1\x23\x23\xa5\x2a\x5d\x53\xe4\x38\x09\xcf\xcc\xee\xca\xbe\xfe\xe2\x19\x7e\x8f\xaf\x91\x63\x3f\xa0\xbc\x1f\x69\x24\x98\x7a\x28\x90\x9d\xd6\xf2\x5c\x0d\x13\x9e\x01\xdf\x12\x50\x52\x90\x72\xb7\x60\x03\xb4\x37\x1d\xe1\x0f\x23\x7c\xad\xfe\x85\x83\xe3\x34\x32\x2e\x20\xe5\xb1\x86\xc5\xb4\xff\xbb\xe0\x11\xc9\x70\x26\xd2\x28\xde\x52\x4c\xe4\x0f\x08\x4b\x53\x9e\x52\x89\x32\x96\x23\x68\x43\x3f\xff\x3c\xb5\x82\x4a\x20\x30\xf5\x6e\x81\xe3\x68\x8f\xaf\xae\x70\x86\x17\x0b\xec\xde\x46\x6c\xff\xb0\xcc\x0f\xae\xe4\x03\x1e\xd2\xa1\x56\x81\x0c\x6a\x6a\x80\x5b\x0d\xcf\x3c\xf2\x24\xb7\xb0\xf6\x5b\x2c\x8b\xb0\xc4\xae\xfd\xc6\xe8\xa4\x78\x5a\x7d\x20\x9b\x7b\x4e\x07\xa7\x75\xd0\xb5\xbd\x55\x6a\x1f\xbd\xb0\x55\x28\xcb\x55\xfc\x45\x9a\x83\xc0\x5b\x21\xa5\x4f\x19\xb7\xf2\x37\x3e\x5f\x99\x50\x5d\x57\x37\x48\xcd\xd4\xc1\x14\x82\x01\x98\x42\xb0\x21\xc0\x14\x82\xf5\x02\x06\xca\xed\x02\x06\xde\xf3\x36\x18\xd9\xf5\x59\x30\xf0\xd8\x27\x4f\xa0\x21\xd7\x03\x80\x81\x6c\x7d\x80\x91\xe5\x76\x00\x23\xff\xcb\x61\x81\x51\x5d\x5f\x02\x33\xf1\x15\x99\x89\x3f\x0c\x9a\x89\xdf\x13\x9b\x89\xdf\x11\xce\xc4\x6f\xd0\x81\xd6\x2f\xe1\x99\x4d\x15\x9e\xd9\x74\x18\x3c\xb3\x69\x4f\x78\x66\xd3\x8e\x78\x66\xd3\x06\x1e\x68\xfd\x12\x1e\x7f\xae\xf0\xf8\xf3\x61\xf0\xf8\xf3\x9e\xf0\xf8\xf3\x8e\x78\xfc\x79\x03\x0f\xb4\x7e\x16\x8f\xfc\x1b\x00\x79\x92\x3f\xd0\x43\xe8\x8e\xcc\xd7\x07\x9e\xbc\xab\xf2\xe4\x4d\xe9\xd1\xad\x5f\xc4\x23\xc5\x27\x1f\x4a\x7d\x54\xc6\xbe\x08\x75\xd3\x9f\xbc\x25\x40\x65\xff\x17\x21\x49\x09\xca\x87\xd2\x20\x95\xb1\x2f\x48\xdd\x54\x28\x6f\xc9\x50\xd9\xff\x45\x48\x52\x88\xf2\xa1\x94\x48\x65\xec\x0b\x52\x37\x2d\xca\x5b\x62\x54\xf6\x7f\x16\x92\xfe\x4b\x21\x79\xc2\xa1\xfa\x1a\x00\x93\xce\xd9\x07\xa7\xb2\xe8\x0e\xa0\xf4\x12\x9b\x94\x81\x70\x19\x95\x3c\x50\xa1\xfa\x1a\x0a\x55\x3f\x47\xaa\x2c\xba\x2b\xaa\xfa\xa1\x32\x10\x4e\xa0\xfa\x2f\x00\x00\xff\xff\xcf\x22\xe9\xc3\xff\x17\x00\x00")

func template_genny_root_go() ([]byte, error) {
	return bindata_read(
		_template_genny_root_go,
		"../template/genny/root.go",
	)
}

var _template_genny_union_alias_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x90\x4d\x6b\xf3\x30\x10\x84\xcf\xd9\x5f\xb1\xe8\x10\x6c\xf3\x62\xdd\x03\x3e\xbc\xb4\x04\x0a\x4d\x28\x94\x9c\x8b\x6c\x6f\x12\x13\x7d\xb8\xfa\x08\x4d\x43\xfe\x7b\x51\xd2\x4a\x3e\xa4\xd0\x9b\x66\xf6\xd1\xb0\x3b\xa3\xe8\x0e\x62\x47\xf8\x1e\xc8\x9e\x00\x78\x05\x3b\xd2\xfa\x84\x2a\x38\x8f\x2d\x61\x27\xa4\xa4\x1e\x47\xb2\xb8\xd1\x83\xd1\x6b\xa1\x08\x2a\x0e\x30\xa8\xd1\x58\x8f\x05\xcc\xd8\x6e\xf0\xfb\xd0\xd6\x9d\x51\xfc\x20\x3e\x03\xdf\xb6\x6e\x4f\x72\x24\xcb\xaf\xa9\xbc\x15\x8e\x18\x94\x00\x47\x61\xf1\x71\xb3\x5a\xbd\xbc\xa5\xac\xff\x72\x10\x2e\x3e\xb0\x35\x46\x62\x83\x11\xae\x5f\xc9\x5f\x07\x05\x4b\x20\xfb\x87\x2c\xc1\xac\x04\xd8\x06\xdd\x61\xa1\x4d\x4f\x79\xb3\x12\x13\x52\x4c\xde\x78\x86\x19\xe7\x96\x5c\x90\x1e\x17\x4d\x1e\x9c\x1f\x8c\x52\x46\xaf\x4d\x4f\x0b\x8c\x51\x75\x36\x2e\x30\xbb\xf7\x23\xd9\x13\x14\x1b\x9c\x67\x35\x41\xa2\x7c\x1e\x9c\xc7\xe6\x96\xfe\xa3\xef\x64\xd4\xd7\x45\x9b\xe9\x91\x89\x5a\x92\xef\xf6\x4f\xba\xa7\x8f\xa2\x8c\xa6\x0f\x56\xe3\x6d\x06\x97\xef\x26\x52\x07\x4b\x6b\x54\x6e\xe1\x88\x55\x12\x25\x56\x89\x8a\x95\xe4\xf3\xe6\xc9\xff\xed\xbe\xe3\x44\xfe\x71\xaf\xaf\x00\x00\x00\xff\xff\x77\x9a\xc7\x0e\x5b\x02\x00\x00")

func template_genny_union_alias_go() ([]byte, error) {
	return bindata_read(
		_template_genny_union_alias_go,
		"../template/genny/union-alias.go",
	)
}

var _template_genny_union_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8f\xbd\x6e\xb3\x30\x18\x85\xe7\xbc\x57\x71\xc4\x10\x01\x03\xd6\xb7\x46\xca\xf4\xcd\x65\xeb\x05\x18\xf3\x42\xac\xe0\x9f\xfa\x47\x15\x45\xdc\x7b\x05\x51\x49\x95\xa1\x8b\x07\xfb\x3c\x8f\xcf\xf1\x52\xdd\xe5\xc8\xf8\xc8\x1c\x66\x22\x51\x13\x80\x91\xad\x9d\x61\x72\x4c\xe8\x18\x4a\x4e\x13\xf7\xf8\x87\xa4\x0d\x47\x78\x0e\x68\x5d\xcf\xad\x34\x4c\xb5\x20\xd2\xc6\xbb\x90\x50\xd2\xa9\x18\x75\xba\xe5\xae\x51\xce\x88\xbb\xfc\xca\x62\xe8\xe2\x8d\x27\xcf\x41\xec\x7e\xd1\xc9\xc8\x05\x55\x44\x69\xf6\x8c\x77\xab\x9d\xdd\x34\x88\x29\x64\x95\xb0\xd0\xa9\xde\x22\xcd\x7f\x67\x8c\xb3\xdb\x2f\xb4\x12\x0d\xd9\x2a\xb0\xf1\x69\x3e\x88\xb2\x42\xfd\xc4\x17\x3a\x05\x8e\x79\x4a\xb8\x5c\x71\x3e\xee\x97\xa7\xe6\x82\xf3\x8b\x78\x59\xd7\x1f\xaa\xd9\x25\x57\x14\x07\x59\x6c\x4f\x29\x07\x8b\x47\xe2\x68\x51\x5a\xd7\xff\x2a\x5e\xe1\x8d\x4d\xc7\xa1\xd4\xd0\x36\x55\xdb\xc1\x61\x90\x8a\x97\xf5\x51\x6a\x77\x58\x3d\x1d\x82\x96\x3f\xff\x1a\xb1\xe7\x5f\xa7\xd2\x4a\xdf\x01\x00\x00\xff\xff\x98\x1d\x38\x8d\xa7\x01\x00\x00")

func template_genny_union_go() ([]byte, error) {
	return bindata_read(
		_template_genny_union_go,
		"../template/genny/union.go",
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
	"../template/genny/dummy.go": template_genny_dummy_go,
	"../template/genny/field.go": template_genny_field_go,
	"../template/genny/list.go": template_genny_list_go,
	"../template/genny/node.go": template_genny_node_go,
	"../template/genny/node_test.go": template_genny_node_test_go,
	"../template/genny/root.go": template_genny_root_go,
	"../template/genny/union-alias.go": template_genny_union_alias_go,
	"../template/genny/union.go": template_genny_union_go,
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
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"..": &_bintree_t{nil, map[string]*_bintree_t{
		"template": &_bintree_t{nil, map[string]*_bintree_t{
			"genny": &_bintree_t{nil, map[string]*_bintree_t{
				"dummy.go": &_bintree_t{template_genny_dummy_go, map[string]*_bintree_t{
				}},
				"field.go": &_bintree_t{template_genny_field_go, map[string]*_bintree_t{
				}},
				"list.go": &_bintree_t{template_genny_list_go, map[string]*_bintree_t{
				}},
				"node.go": &_bintree_t{template_genny_node_go, map[string]*_bintree_t{
				}},
				"node_test.go": &_bintree_t{template_genny_node_test_go, map[string]*_bintree_t{
				}},
				"root.go": &_bintree_t{template_genny_root_go, map[string]*_bintree_t{
				}},
				"union-alias.go": &_bintree_t{template_genny_union_alias_go, map[string]*_bintree_t{
				}},
				"union.go": &_bintree_t{template_genny_union_go, map[string]*_bintree_t{
				}},
			}},
		}},
	}},
}}
