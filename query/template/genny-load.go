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

var _template_genny_dummy_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8e\x3d\x6e\xc3\x30\x0c\x85\x67\xf3\x14\x44\x86\x20\xc9\x10\xed\xdd\x8a\xee\x1e\x8a\x5e\x40\x72\x68\x5b\x88\x28\xa9\x12\x89\xc2\x35\x7c\xf7\xc2\xfd\x81\x3d\x15\xd9\x88\xc7\x8f\xdf\x63\xb6\xdd\xdd\x0e\x84\xef\x4a\x65\x02\x30\x06\x65\xf4\x15\x7d\xc5\x9b\x32\x4f\x38\x50\x8c\x13\x0a\x71\x0e\x56\x08\x7b\x1f\x08\x3c\xe7\x54\x04\x4f\xd0\x1c\x06\x2f\xa3\xba\x6b\x97\xd8\xdc\xed\xa7\x9a\xde\xd5\x91\x42\xa6\x62\xbe\x85\xc6\xd9\x4a\x07\x68\x1c\x3e\x48\x9e\x01\x64\xca\x84\xaf\x54\x35\xc8\xdb\x3a\xb6\xe9\x46\xad\x65\xfa\xdb\xa4\xf4\x93\x57\x29\xda\x09\xce\xd0\x5c\xdc\xf5\x25\x31\xa7\xb8\xa2\xb0\xfc\x82\xcf\xc1\xdb\xba\x1e\xfe\x47\xf6\x1a\x3b\x6c\xe9\x63\xeb\x3b\x9d\xf1\xb2\x6b\x9f\xa1\x29\x24\x5a\x22\x1e\xb7\x74\xde\x24\x4f\x78\x5c\x3f\xdf\x69\xe7\x65\x81\x05\xbe\x02\x00\x00\xff\xff\x9f\x7a\x36\xb0\x58\x01\x00\x00")

func template_genny_dummy_go() ([]byte, error) {
	return bindata_read(
		_template_genny_dummy_go,
		"../template/genny/dummy.go",
	)
}

var _template_genny_field_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x51\xcb\x6e\xdb\x30\x10\x3c\x7b\xbf\x62\xc1\x13\x69\xb8\xe2\xbd\x80\x0a\x18\xe8\x03\x05\x6a\x1d\x5a\xf7\xd0\x93\xa1\xc7\xda\x26\x2c\x91\x2a\x45\x26\x56\x82\xfc\x7b\x40\x4a\x96\x9c\xc4\xb0\x0e\x12\x67\x67\x76\x76\x96\x6a\xf3\xf2\x94\x1f\x08\xff\x7b\xb2\x3d\x80\x6a\x5a\x63\x1d\x72\x58\xb0\x83\x72\x47\x5f\x24\xa5\x69\x64\x79\x24\x3a\xf5\x85\x72\x9d\x3c\x90\xd6\x7d\x78\x93\x55\x25\x7b\x2b\x3b\xe5\x4f\x5e\xee\x8b\xee\x48\x75\x4b\x56\x46\x4b\x59\xe4\x1d\x31\x10\x00\x72\x09\xb1\x19\x1b\xdf\x39\x2c\x08\xcb\xbc\xae\xa9\xc2\x96\x2c\x7e\x57\x54\x57\xb0\x94\x00\xae\x6f\x69\x80\x59\xde\x10\x8e\x83\x92\x6d\xdf\xd2\x35\x19\xf0\x3b\xf2\x21\xb7\x21\x77\x66\x2a\x0a\xad\xbb\xc9\x64\x3c\xf9\x06\x95\x76\x98\x62\x48\x94\xac\x9d\x51\x99\xf9\x66\x2d\x0f\x27\xce\x2e\x1a\x26\xc4\x2d\x0f\x9c\x9e\xc1\xe3\xce\x94\xb8\xab\x44\xbe\x0f\x18\x95\xa6\xea\xbc\xc2\x01\x84\xfc\x02\x3f\x7d\xc1\xb9\xff\x67\x75\xde\x9a\xb8\x40\xc8\xff\xf5\xef\x66\xf3\x6f\x77\x63\x7c\x61\x4c\x8d\x29\xfe\x21\x77\x21\x23\xd7\x71\x76\xc1\x6c\x85\x6c\xd2\x4f\x20\x38\xb3\xd5\xbd\xb8\x02\x60\xef\x75\x89\x5c\x9b\x8a\x26\xa1\x98\xff\x01\x17\xc8\x2d\x75\xbe\x76\xb8\xfc\x1d\xbf\xdb\xb8\xc6\x33\x2c\xc6\x72\x8a\x19\x3d\xce\x14\x17\xb0\x28\x4d\xd3\x18\x8d\x9f\x53\x0c\xb6\x49\x34\x5b\x3b\x7e\x3f\xc7\xe8\x97\xc4\x8d\x53\x1c\x3c\x22\x9a\x29\x53\xd1\x2f\xd5\xb9\x2b\x7a\xac\x4c\x92\xe9\x46\x67\xcd\x7c\xc9\x1f\x44\x3f\xac\xf1\xed\x0d\x65\xac\xc7\x48\xce\x5b\x0d\x2f\xf0\x1a\x00\x00\xff\xff\xcd\x2a\xfa\x2f\x29\x03\x00\x00")

func template_genny_field_go() ([]byte, error) {
	return bindata_read(
		_template_genny_field_go,
		"../template/genny/field.go",
	)
}

var _template_genny_node_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x55\x4b\x6f\xe3\x36\x10\x3e\x6b\x7e\xc5\x40\x87\x05\x19\x18\x51\x7b\x0d\x20\x14\x69\x5a\x07\x06\x36\xde\x62\xe3\x3e\x00\x23\x08\x64\x99\x72\x58\x4b\xa4\x4a\x51\x69\x52\xaf\xff\x7b\x41\x52\x0f\x4a\x96\x9d\xf6\xb2\xeb\x90\xdf\x8b\x33\x1c\xaa\x4c\xd2\x7d\xb2\x63\xf8\x57\xcd\xd4\x3b\x00\x2f\x4a\xa9\x34\x12\x08\xc2\x1d\xd7\x2f\xf5\xe6\x3a\x95\x45\xb4\x4f\xfe\xa9\xa3\x6c\x53\xbd\xb0\xbc\x64\x2a\xb2\xd8\x68\x93\x54\x2c\x84\xe0\x19\x3f\x44\xe6\x72\x17\x02\x05\x80\xe8\x0a\x10\x11\x8b\xba\xd2\x98\x26\x79\x8e\xdf\xa3\xe6\x05\xab\xb0\x64\x0a\x57\xc9\x26\x67\x51\xa5\x55\x9d\x6a\xb2\x94\x5b\xb6\x4c\x0a\x46\xe1\x2a\x02\xd0\xef\x25\xc3\x76\x09\x1d\x04\x0f\x10\x5c\x99\x0c\xd7\x77\xb2\x28\xa4\x30\xdb\x70\x04\xc8\x6a\x91\xe2\x92\xfd\xdd\xc2\x09\xc5\xab\x8e\x7a\x80\x40\x31\x5d\x2b\x81\x9f\xda\xb5\x43\x4f\xbf\xc1\x4f\x23\xc1\xc3\xf1\x68\x34\x5f\x13\xd5\xd9\x3f\x2f\xb6\x6f\x2b\xb9\x32\x89\x8a\xa4\x5c\x73\xa1\x9f\xb8\xd0\x18\xfb\x7f\x1d\x8e\x67\x28\xf7\x4a\xd6\xe5\xff\xe4\xd9\xe0\x2d\xa8\xd2\x8a\x8b\x9d\xc7\x72\x0b\x63\xa2\xf9\x67\x25\x17\xdb\x37\x8b\x73\x18\xcf\xae\x5f\x38\x1c\xc1\x1d\xef\xa7\x5f\x1f\x1e\x7e\x79\x6e\x05\x16\xd5\xa3\x2b\xf2\x46\xca\x1c\x63\xb4\x65\x79\x64\xba\xdf\x13\xa9\x26\x61\x0b\x0f\x67\x0e\xb1\x92\x3f\x4a\x99\x93\xb0\xa5\x87\x94\x36\x0d\x31\xdc\x06\x3c\xe7\x2c\xdf\x56\x44\x98\xdf\x33\xcc\x9a\xff\x6c\x41\x5d\x2e\xb3\x58\x17\xc8\x85\xa6\xce\xff\x00\x01\x13\x75\x31\x37\xd7\x60\x86\x72\x8f\x37\x4d\x22\x77\x4c\x43\x5d\x5b\x81\x27\x08\x78\x66\x00\x07\x08\x82\xd6\xef\x91\xe9\xae\xfe\xc4\x28\xcf\xb0\x53\xa3\x10\x1c\x21\x88\xa2\xf9\xe2\x8f\x87\x9f\x6f\x8c\x28\x4f\xd1\xde\x36\x29\xf2\x77\xac\xb4\x54\xec\x07\x80\x13\x2d\x7b\xb1\x9c\x96\x35\xa6\x00\xc1\x4e\x95\x26\x58\x0b\xbd\x67\xba\xeb\x38\x69\x40\xbe\x8e\xb7\x69\x75\x76\xaa\xa4\x9e\x53\xdf\xfb\xb5\xd9\x7f\xc2\xd8\x39\xf9\x90\xae\xcb\x6b\x5b\x45\x8b\x59\xd6\x05\x04\x7e\xbb\x2c\x60\xd0\xaa\x53\xbe\x31\x6e\xe6\x42\xab\x9a\x01\x1c\x9b\x31\x9a\x3a\x36\x37\x9d\x99\x61\xd5\x74\x8b\x9a\x62\x4f\xa4\xe6\x26\x4e\x05\x7d\x98\x5e\x61\x32\x4c\xb7\x4d\xfb\x29\x9e\xea\xe0\x7e\x86\xaf\xc3\xbb\x71\x3a\x66\xeb\xbd\x31\x7f\x1d\x79\x5b\xfa\x79\xef\xa6\x45\x7e\x1d\x26\x82\xf4\x6d\xfb\x2f\x41\x2c\xf2\x24\xcd\x60\xf7\x83\x44\x16\xf3\x41\xac\xc1\x55\xeb\xfb\x42\x14\xab\xea\x5c\xbb\x8c\xfd\xcb\x67\x53\x0c\x29\x7d\xc9\x89\x90\xdb\xfe\xa9\xa5\x98\x76\x4f\xa1\x79\x47\x47\xaf\xa3\x11\xe5\x19\x1a\x8a\xbf\x1a\xc7\x28\xb8\x2d\x88\x3b\xf1\x67\xb9\x23\xee\xc7\x97\xfb\xe7\xdf\x6f\xbf\x2e\x67\x68\xbc\x08\xc5\x76\xfb\x56\xed\x2a\x8b\x1f\x64\x9c\x93\xd0\x53\x15\x52\x63\x26\x6b\xb1\xed\xd2\x85\x14\x82\xe0\x68\x46\x98\xe5\x15\x43\x9e\x61\xce\x04\x19\xa5\xb9\x76\xe7\x88\x63\xfc\x0e\xbf\x7d\x9b\x44\xf4\xdd\x77\x30\x13\x64\x4a\x05\x63\xec\x7b\x35\x81\xe9\x3f\x0c\xf1\x44\x23\x2f\x11\xdc\x67\x61\x8a\x65\x77\xec\x23\xd5\x54\x66\x5c\xec\x76\x54\xc7\x8d\x7b\x64\x89\x4a\x5f\x16\x22\x93\xa4\x94\x95\x1b\xd9\xac\x29\xec\x57\x96\xce\xc5\x0c\x53\x29\xb6\xf3\x66\xe9\xce\xfe\x36\xd7\x04\x5c\x4c\xbf\xf1\xd7\x43\x31\x23\xd4\xb2\x29\x04\xe6\xee\x4c\x67\xb0\x04\x8a\x84\x8b\x4c\x3a\x1b\xb3\xe2\x4c\xfc\xe3\x0c\xbc\x1c\x09\xce\xde\xc8\x45\xf5\x99\x25\xd9\xad\x26\x7f\x0e\x86\xef\x92\x62\xc7\xb8\x20\x7b\x27\x6b\xa1\xbf\x64\xf6\xe3\x44\xa8\x91\xf6\x46\xc6\xdc\x9a\xa9\xc7\xe2\x9c\xd8\x6f\x49\x5e\x33\x7b\x10\xde\x84\x34\x87\xef\x56\x3d\xe5\xd3\xb4\x1e\xd7\x1a\x9c\xb1\xb0\x41\x6f\x35\xe1\xdb\x37\x67\x41\xc6\xd3\x49\x2f\xda\x78\x7c\x6b\xf3\x6f\x00\x00\x00\xff\xff\x76\x1e\x2b\xf4\x08\x0a\x00\x00")

func template_genny_node_go() ([]byte, error) {
	return bindata_read(
		_template_genny_node_go,
		"../template/genny/node.go",
	)
}

var _template_genny_node_test_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\xcf\x6a\xf3\x30\x10\xc4\xcf\xda\xa7\x58\x7c\x08\xd6\xf7\x19\xab\xe7\xd2\xf4\x18\xc8\xc5\x04\x92\x7b\x91\xec\x75\x2c\x62\x4b\xaa\xfe\x50\x9a\x90\x77\x2f\x52\x92\x06\x7a\xd3\x68\xf4\x9b\x59\xad\x93\xfd\x49\x1e\x09\x3f\x13\xf9\x6f\x00\xbd\x38\xeb\x23\xd6\xc0\xaa\x48\x21\x6a\x73\xac\x00\x98\xc2\xea\xa8\xe3\x94\x54\xdb\xdb\x45\x9c\xe4\x39\x89\x51\x85\x89\x66\x47\x5e\x14\x52\x28\x19\xa8\x02\x0e\x30\x26\xd3\xe3\x81\x42\xfc\x08\x24\x7d\x3f\x6d\xcd\x68\xeb\x88\xff\xee\x71\xed\x81\xe3\x05\x80\x49\x7c\x5d\x63\x47\x5f\x9d\x1d\xa8\x93\x0b\xd5\x1c\x98\x10\x0f\xb5\xd1\x34\x0f\xa1\x4b\xcb\xd6\xe8\x98\xad\xde\x9a\x21\x13\x39\xbd\x76\x36\xa0\x36\xb1\x41\x6d\x46\x8b\xaa\xcd\x15\x1c\x95\xb5\x33\x5e\x80\x31\x4f\x31\x79\x53\xcc\x76\x67\x03\xbe\xad\x31\x13\xab\x15\xd6\x8f\xbb\xff\xe5\xb0\xd7\x67\xe2\xf8\x9e\x5d\x60\x4c\x88\x3b\x18\x7d\x22\x60\x57\x60\x9e\xfa\x8d\xf9\x6d\x0d\xa8\xda\x3c\xdf\x4e\xc6\xe9\x4f\xf5\xa5\xe0\xb9\x39\xa4\x39\xe2\x1a\xa5\x73\x64\x86\xfa\xa6\x1b\x0c\xfc\xfe\x20\x53\xe1\xe9\x17\x79\xcb\xe2\xa5\x51\xb6\xfb\xe7\xd2\x5e\x1a\x2c\x13\x34\x98\x7f\xcf\x01\xae\xf0\x13\x00\x00\xff\xff\x79\x9d\x3b\x9e\xae\x01\x00\x00")

func template_genny_node_test_go() ([]byte, error) {
	return bindata_read(
		_template_genny_node_test_go,
		"../template/genny/node_test.go",
	)
}

var _template_genny_root_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x94\xcf\x8a\xdb\x30\x10\xc6\xcf\x9e\xa7\x98\xfa\x10\xa4\x36\xd8\xdd\xdd\x10\xca\xd2\x1c\x9a\xc2\xb6\x85\x92\x85\xb0\x7b\x2a\x3d\xd8\xc9\xd8\x2b\xd6\x91\x5c\x59\x6e\xea\x86\xbc\x7b\x91\xfc\x47\xca\x42\x2f\x7b\xd3\x7c\x33\xf3\x1b\xcd\x67\xe4\x3a\xdb\x3d\x67\x25\xe1\xaf\x96\x74\x07\x20\x0e\xb5\xd2\x06\x19\x44\xb1\x50\x31\x44\x71\x63\xf4\x4e\xc9\xdf\x31\x40\x54\x54\x99\xc9\xdb\xa2\x20\xdd\x60\x5c\x0a\xf3\xd4\xe6\xc9\x4e\x1d\xd2\x52\xa9\xb2\xa2\x34\x48\xa7\xa5\xeb\x0d\x6a\x9e\xb3\xbf\x6d\x5a\xe4\xcd\x13\x55\x35\xe9\xd4\x4d\x4b\xf3\xac\xa1\x18\x38\x80\xe9\x6a\xc2\xb5\x52\x15\xae\xd0\x8a\x89\x3d\x0f\x6a\x67\x68\x52\x3b\x43\xbd\xfa\x4d\x9a\x0f\xa3\x6a\xcf\x93\x7a\xb5\x0c\xe4\xab\xe5\xa4\xdf\x5c\x07\xfa\xcd\xf5\xa4\x2f\x17\x81\xbe\x5c\xf4\xfa\xa3\x08\xf0\x2e\xf0\xba\x1f\xd0\x47\x3e\xe3\x47\xf4\x91\xcf\xf8\x21\x7d\xd4\x67\xee\x2a\x95\x05\x4d\x43\x18\xe4\x7c\xdb\x10\x0e\x46\x7d\x56\x87\x83\x92\x1b\xb5\x9f\x8c\xf1\x0a\x40\xd1\xca\x1d\x6e\xe8\xe8\x35\xc6\xf1\xed\x8b\x32\x3c\x41\xa4\xc9\xb4\x5a\xe2\xec\x45\xea\x74\x86\xf3\x40\xb9\xaf\x49\x32\x8d\x42\x25\x5b\xca\xf6\xa4\xe7\xb8\xcb\x6a\x14\xd2\x70\xdc\x2a\x65\x1e\xec\x65\x3c\x68\x94\x4e\x9e\x75\xdb\x5f\xaf\xe7\xb8\x6e\x7e\x49\x5f\x77\xeb\xb6\x60\x79\x5b\xe0\x8f\x9f\x79\x67\xe8\x15\xe0\x09\x11\xa0\x99\xb4\x3b\x8e\x7d\x1c\x37\xf4\xc7\xb0\x4b\x76\x63\x32\x6d\xf0\x76\x85\xb6\x34\xf9\x4e\x92\x71\x80\x48\x14\x53\x6c\xa9\xfc\xdd\x02\x3f\x62\x5f\x7a\x82\x68\xbc\x90\x2d\x81\xe8\x0c\x10\x49\x3a\xae\xb3\x86\x26\x8e\x0d\x12\x3b\xcd\x1e\x98\x6b\xb4\x58\xad\x94\x9b\x35\x6d\x72\xee\xb5\x24\xfc\x88\x1b\x3a\xba\xaf\x35\x30\xe7\xd6\x68\x16\x3c\xaa\xe4\x0b\x99\xc7\xfb\xa2\x68\xc8\x3c\x8c\x45\xc9\x96\xbd\xe7\x9c\xf3\xc9\x2b\x4b\xfd\xaf\x0d\x5f\xb3\x66\x70\x22\xb7\x6f\xed\x04\x10\x6e\x74\xb9\x74\x60\xcb\x88\xfb\x64\x94\x60\x0d\x36\x46\x0b\x59\x72\x64\x42\x9a\x39\x92\xd6\x4a\x73\xc7\x92\x2e\xb2\x7b\x0e\x3f\x8d\xa4\xef\xe0\xce\x57\x9b\x7a\xb3\x42\x29\x2a\x9c\xcd\xb0\xc1\xd5\x0a\xe3\x3b\x41\xd5\x7e\xd3\x1e\x62\xe7\xae\xad\x70\x05\xd6\xdb\xe9\x66\x0e\x0a\x67\xf8\x17\x00\x00\xff\xff\x84\x02\x6c\xb5\xa9\x04\x00\x00")

func template_genny_root_go() ([]byte, error) {
	return bindata_read(
		_template_genny_root_go,
		"../template/genny/root.go",
	)
}

var _template_genny_union_alias_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8f\xcb\x6a\xeb\x30\x10\x86\xd7\x9e\xa7\x18\xb4\x38\xd8\xe1\x60\xed\x03\x5e\x94\x96\x42\xa1\x09\x85\x92\x75\x91\xed\x89\x2d\xa2\x8b\xab\x4b\x69\x1a\xf2\xee\x45\x0e\x95\xbd\xc8\x6e\xe6\xff\x3f\x7d\x8c\x26\xd1\x9d\xc4\x40\xf8\x19\xc9\x9d\x01\xf8\x06\x06\x32\xe6\x8c\x3a\xfa\x80\x2d\x61\x27\x94\xa2\x1e\x27\x72\x78\x30\xd2\x9a\xbd\xd0\x04\x1b\x0e\x20\xf5\x64\x5d\xc0\x12\x0a\x36\xc8\x30\xc6\xb6\xee\xac\xe6\x27\xf1\x13\xf9\xb1\xf5\x23\xa9\x89\x1c\x9f\xad\xbc\x15\x9e\x18\x54\x00\x5f\xc2\xe1\xd3\x61\xb7\x7b\xfb\xc8\xae\x07\x25\x85\x4f\x03\xb6\xd6\x2a\x6c\x30\xc1\xf5\x3b\x85\xb9\x28\x59\x06\xd9\x7f\x64\x19\x66\x15\xc0\x31\x9a\x0e\x4b\x63\x7b\x5a\x2e\xab\x30\x23\xe5\x6a\xc6\x0b\x14\x9c\x3b\xf2\x51\x05\xdc\x36\x4b\x71\x79\xb4\x5a\x5b\xb3\xb7\x3d\x6d\x31\xa9\xea\x25\xb8\x42\x71\xef\x45\x8e\x57\x28\x36\xf8\x6f\xd9\x56\x48\x5a\x5f\xa5\x0f\xd8\xdc\xec\x7f\xfb\x1d\x47\x3d\x1f\xda\xac\x3f\x99\xa9\x67\x0a\xdd\xf8\x62\x7a\xfa\x2e\xab\x14\x86\xe8\x0c\xde\x3a\xb8\xc2\x6f\x00\x00\x00\xff\xff\xbc\x73\x10\xa1\xc3\x01\x00\x00")

func template_genny_union_alias_go() ([]byte, error) {
	return bindata_read(
		_template_genny_union_alias_go,
		"../template/genny/union-alias.go",
	)
}

var _template_genny_union_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8e\xbd\x6e\xac\x30\x10\x85\x6b\xe6\x29\x8e\x28\x56\x40\x81\x75\xdb\x95\xb6\xba\x75\xe8\xf2\x00\xc6\x0c\xac\xb5\xd8\x26\xfe\x51\x44\x10\xef\x1e\x99\x55\xd8\x28\xcd\x14\x33\xe7\x7c\xf3\x2d\x52\x3d\xe4\xc4\xf8\x48\xec\x57\x22\xd1\x10\x80\x89\xad\x5d\x61\x52\x88\xe8\x19\x4a\xce\x33\x0f\xf8\x87\xa8\x0d\x07\x2c\xec\xd1\xb9\x81\x3b\x69\x98\x1a\x41\xa4\xcd\xe2\x7c\x44\x45\x45\x39\xe9\x78\x4f\x7d\xab\x9c\x11\x0f\xf9\x95\xc4\xd8\x87\x3b\xcf\x0b\x7b\x71\xf0\x45\x2f\x03\x97\x54\x13\xc5\x75\x61\xbc\x5b\xed\x6c\xc6\x20\x44\x9f\x54\xc4\x46\x45\x93\x23\xed\x7f\x67\x8c\xb3\xf9\x0b\xed\x44\x63\xb2\x0a\x1d\x7f\x9e\xf9\xaa\x46\xf3\x2a\x6f\x54\x78\x0e\x69\x8e\xb8\xde\x70\x39\xf7\xdb\x0b\x72\xc5\xe5\x0f\x76\xdb\xf7\x9f\x56\x7b\x40\x6e\x28\xcf\x66\x99\x4f\x31\x79\x8b\x67\xe2\x74\xa8\xac\x1b\x7e\x69\xd7\x78\x63\xd3\xb3\xaf\x34\xb4\x8d\x75\x1e\xec\x47\xa9\x78\xdb\x9f\x52\x07\xc3\xea\x99\x76\xfa\x0e\x00\x00\xff\xff\x89\x66\x96\xee\x68\x01\x00\x00")

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
