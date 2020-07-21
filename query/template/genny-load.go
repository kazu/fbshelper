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

var _template_genny_basic_nogo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\x4f\x4f\x84\x30\x10\xc5\xcf\x3b\x9f\x62\xc2\xc1\x80\x21\x60\x3c\x9a\x70\x59\x13\x4f\x66\x3d\xac\x37\xe3\xa1\xad\x03\x34\x58\x8a\xfd\xa3\x61\x37\xfb\xdd\x4d\xc1\xc2\xe2\xa5\xe9\xeb\xfb\xbd\x37\xcd\x0c\x4c\x74\xac\x21\xfc\xf2\x64\x46\x00\xa9\x06\x6d\x1c\xa6\xb0\x4b\x1a\xe9\x5a\xcf\x0b\xa1\x55\x29\x5a\xa2\x6e\xe4\xd2\xd9\xb2\xa1\xbe\x1f\xc3\x49\x46\x8a\x64\x8b\x75\xec\xe4\xcb\x9a\xdb\x96\x3e\x07\x32\xe5\x54\x59\x72\x66\x29\x81\x0c\xc0\x8d\x03\xe1\x9e\x59\x29\x5e\xc3\xed\xaf\xa2\x38\x78\xc5\xc9\x00\xd4\xbe\x17\xf8\x64\xb4\x5a\x90\xf4\x7b\xc5\x33\xbc\x7d\xd4\x4a\xe9\xfe\xa0\x3f\x08\xcf\x00\x3b\x31\x49\x7c\xa8\xf0\x26\x8c\x28\x56\xfb\x7c\x89\x6e\x11\xe4\xb3\xb4\x0e\x23\x15\x1f\xae\x19\xa6\x08\x2b\x4c\x96\x59\xc9\x26\x8e\x15\xce\x49\xfa\x09\xf2\x3e\x8d\x6a\xcf\x2c\xa5\x8a\x75\x94\xbe\xbd\xf3\xd1\x51\x3e\x73\x47\x79\xa2\x97\x7a\xfd\x78\x96\xe3\x5d\x8e\xce\x78\xca\x36\xbd\x13\x18\xcb\xff\x85\x16\xf0\x48\xee\x6a\x1d\x19\xec\x0c\x39\x6f\x7a\x9c\x6d\xb8\xc0\x6f\x00\x00\x00\xff\xff\xfa\xfb\xaf\x05\xbf\x01\x00\x00")

func template_genny_basic_nogo() ([]byte, error) {
	return bindata_read(
		_template_genny_basic_nogo,
		"../template/genny/basic.nogo",
	)
}

var _template_genny_dummy_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8e\xcb\x6a\xf3\x30\x10\x85\xd7\x9e\xa7\x18\xb2\x08\x49\x16\xd1\xbf\xf8\x57\x85\x2e\xda\x42\x97\x2e\xb4\x7d\x01\x59\x1e\xdb\x22\xba\x55\x1a\x51\x1c\xe3\x77\x2f\xca\x05\x07\x0a\x25\x3b\x71\xe6\xe8\x3b\x5f\x90\xea\x20\x7b\xc2\xaf\x4c\x71\x04\x10\x02\x79\xd0\x09\x75\xc2\x36\x5b\x3b\x62\x4f\xce\x8d\xc8\x64\x83\x91\x4c\xd8\x69\x43\xa0\x6d\xf0\x91\x71\x03\xd5\xaa\xd7\x3c\xe4\x66\xaf\xbc\x15\x07\x79\xcc\xa2\x6b\xd2\x40\x26\x50\x14\x27\xa0\x68\x64\xa2\x15\x54\x0d\xde\xd9\xdc\x02\xf0\x18\x08\xdf\x29\x65\xc3\x9f\xe5\x59\xfb\x96\x6a\x69\xe9\x7a\xf1\xfe\x9c\x27\x8e\x59\x31\x4e\x50\xed\x9a\xfd\x8b\xb7\xd6\xbb\x52\x85\xf9\x52\x7c\x32\x5a\xa6\xf2\xf1\xaf\x66\x97\x9d\x42\xb2\x81\xc7\x65\x71\xb3\xc5\xdd\xcd\xfe\x04\x55\x24\xce\xd1\xe1\x7a\x49\xa7\x05\xf3\x80\xeb\xe2\x7e\x03\x9e\xe6\xb9\xb0\x95\x77\x89\xf1\x55\x93\x69\xeb\x6c\x51\x3b\xc6\x47\xfc\x77\xcd\x3f\xf4\x91\xde\xba\x67\x99\xb4\x3a\xcd\x9c\xcf\xff\x2f\x4a\x35\x7d\xdf\x21\xf4\x4b\x1c\x66\xf8\x09\x00\x00\xff\xff\x03\x7e\x6b\x0a\xd1\x01\x00\x00")

func template_genny_dummy_go() ([]byte, error) {
	return bindata_read(
		_template_genny_dummy_go,
		"../template/genny/dummy.go",
	)
}

var _template_genny_field_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x53\x4d\x8f\xda\x30\x10\x3d\x33\xbf\x62\x94\x93\x8d\x68\x72\xaf\x94\x4a\xa8\x5f\xaa\xd4\xe5\x50\xe8\xa1\xa7\x55\x3e\x26\x60\x91\xd8\xa9\x63\xd3\x4d\x57\xfc\xf7\xca\x0e\xb1\xd9\x96\xa2\x96\x03\x78\x3c\x6f\xde\x7b\x33\x63\xfa\xa2\x3a\x16\x7b\xc2\xef\x96\xf4\x08\x20\xba\x5e\x69\x83\x0c\x16\xc9\x5e\x98\x83\x2d\xd3\x4a\x75\x59\x75\x20\x3a\x8e\xa5\x30\x43\xb6\x27\x29\x47\xf7\x4d\x5a\x54\xc9\x4b\xd8\xb1\xf8\x69\xb3\xa6\x1c\x0e\xd4\xf6\xa4\x33\x4f\x99\x95\xc5\x40\x09\x70\x80\x6c\x09\xbe\x18\x3b\x3b\x18\x2c\x09\xab\xa2\x6d\xa9\xc6\x9e\x34\x7e\x10\xd4\xd6\xb0\xcc\x00\xcc\xd8\xd3\x14\x6e\x8a\x8e\xf0\x22\x94\xee\xc6\x9e\xae\x93\x2e\xfe\x2d\x79\x2a\xb4\xf3\xbd\x51\x35\xb9\xd2\xc7\x40\x72\x39\xd9\x0e\x85\x34\x98\xa3\x73\x94\xae\x8d\x12\x1b\xf5\x5e\x6b\xe6\x4e\x2c\x99\x31\x09\xe7\xb7\x38\x30\x7c\x26\x8e\x3b\x2a\xbe\xd7\x0c\x59\xe3\x62\x14\x92\xea\xa7\x15\x4e\x81\xf3\xcf\xf1\xd5\x1b\x8c\xf5\x9f\xea\xa7\x9d\xf2\x0d\x38\xff\xef\xbe\x3e\x3c\x7c\x7b\xbc\x21\x5f\x2a\xd5\x62\x8e\x5b\x32\x73\xd2\xe7\x06\x96\xcc\x71\xb2\xc2\x24\xe0\x43\xe0\x98\x93\xd5\x3d\xbb\x1c\xa0\xb1\xb2\x42\x26\x55\x4d\x01\xc8\xe3\x0e\x18\x47\xa6\x69\xb0\xad\xc1\xe5\x17\xff\xbb\xf3\x6d\x3c\xc3\xe2\x72\x9d\x23\x75\xbd\x19\x63\x92\x71\x58\x54\xaa\xeb\x94\xc4\xd7\x39\x3a\xe2\xd4\xd3\xad\x0d\xbb\xeb\x64\x21\x1a\x9c\xea\x52\x87\xc3\x3c\x47\x29\x5a\xa7\x14\xa5\x36\xf4\xe3\xa5\xd0\xc2\xf3\x6f\xc9\xfc\x83\xc4\x0a\x27\x9e\x74\x4b\x6d\xb3\x1e\xde\x7a\x2d\x87\x67\x6e\xed\xb3\xe7\xff\xb2\x7c\x86\x79\x0e\xa9\xdf\x54\x1e\x3a\x28\x3a\x8a\x29\x55\xd3\x67\x31\x98\xab\xf4\xe5\x26\x40\xc2\x4b\x88\x98\xf8\x38\xfe\x00\x7d\xd4\xca\xf6\x37\x90\xfe\xde\x5b\x32\x56\x4b\x38\xff\x65\xbb\xf3\xb8\xfc\x82\x4f\xb8\xf4\xff\x89\x38\x0e\x8e\xa4\xb5\xd2\xf8\x1c\x98\xa6\x91\x44\xc4\xf5\xc0\xe3\x74\x4f\x1c\xce\xf0\x2b\x00\x00\xff\xff\xb8\x34\x9c\xdb\x53\x04\x00\x00")

func template_genny_field_go() ([]byte, error) {
	return bindata_read(
		_template_genny_field_go,
		"../template/genny/field.go",
	)
}

var _template_genny_list_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcf\x8a\xdb\x30\x10\xc6\xcf\xd1\x53\x0c\x39\x04\x3b\xb8\xf6\xf6\x1a\xf0\x21\x2c\x14\x16\x42\x2e\x5b\xf6\xb2\x2c\x45\x71\xc6\x59\x11\x4b\x72\x65\x69\x4b\x6a\xfc\xee\x65\xfc\x57\x4d\xeb\x4d\x29\x7b\x30\x98\x99\xd1\xf7\xfd\x66\x34\x2a\x79\x76\xe6\x27\x84\xef\x0e\xcd\x85\x31\x21\x4b\x6d\x2c\x2c\x4f\xc2\xbe\xba\x43\x9c\x69\x99\x9c\xf9\x4f\x97\xe4\x87\xea\x15\x8b\x12\x4d\xd2\x16\x26\x07\x5e\xe1\x92\x31\x7b\x29\x11\x76\xa2\xb2\x5f\xe9\xa7\xb2\xc6\x65\x16\x6a\x48\x12\x38\xa1\x52\x17\xb6\x58\xdf\x6b\x29\xb5\xda\xeb\x23\xb2\x86\xb1\x24\x01\xfa\xdd\x73\x89\x7d\x45\xee\x54\x06\x7b\xfc\x31\x88\x04\x21\xac\x47\xc1\x9a\xb1\x45\x21\x2a\x0b\x9b\x14\x50\x96\xf6\x32\x55\x75\x89\x98\xd4\x28\x08\x29\xac\x88\x69\x0c\xd4\x4d\x5f\x31\x01\xc4\xad\x6d\x0a\xcb\xe7\x97\x01\x62\xd9\x1b\xc4\x0f\x4a\x58\x3a\x46\xc2\x06\xad\x33\x0a\x28\x4e\xcc\x2d\xe1\x95\xfb\x6f\x8c\xc3\x81\xd5\x10\xab\x27\xcf\x4d\x8f\x35\x45\xea\xa6\x19\x55\x03\xa5\x8f\xd3\xfc\x42\xd8\xda\x40\x80\x50\x36\x84\xc0\x60\xe5\x0a\x0b\xeb\x81\x34\x02\x04\x34\x46\x9b\xb0\x73\x6c\xb3\x29\xac\x86\x3c\xf5\xdb\x45\x3d\x2f\x3a\x94\x02\x99\xf8\x63\x20\x97\xb1\xcd\x59\x96\x47\xb4\x03\x4e\x04\x6f\x13\x48\xd8\x61\x78\x7d\x5f\xeb\xf7\x27\x23\x78\xf3\xa2\xe1\xac\xd1\x17\x61\x68\xee\xb7\x5b\x9e\xcc\xb6\x36\xb8\x9b\x17\xdc\xf1\xff\xd0\x13\xca\xb6\x3a\xe3\x02\xc5\x4f\xbc\x70\xf8\xa0\x72\x1d\x3f\xed\x50\x85\xf0\x09\x3e\xcf\x7b\x3e\x62\x81\x99\x0d\x72\x05\x94\x0e\xbc\x69\x1d\xb4\x2e\x26\x98\xe7\x17\x2f\xe5\x5f\xa4\xe4\x67\x0c\xbc\x6c\x04\x77\x11\xdc\x84\x0a\xd9\x22\x6b\x47\x5c\xd1\x13\xf9\xf3\x22\x3a\x28\x22\xca\x24\x78\x6f\xb1\xc3\x22\x80\x61\x0e\xb9\x0a\xa6\x5d\xf2\xf7\x37\x93\x4d\xc8\x16\xf4\xe5\xda\xc0\xb7\x08\x32\x49\x5e\x86\xab\x13\xc2\x60\xde\x09\xf5\xad\xf0\xb2\x44\x75\xec\x3b\x8e\xe0\x5d\xd9\xf1\x1e\xba\xea\x77\x96\x44\x1d\xe7\xa7\x3b\x06\xbc\x91\x6e\x6e\x3d\x8e\xbf\xbc\x8c\xce\xe5\x83\xc6\xf5\x6f\x8d\x6d\x8b\x22\x08\xfd\xb5\xb8\xde\x4d\xff\x12\xaf\x1a\x87\x1a\xfa\x4a\x6b\x1c\x42\x33\xbf\x9e\xf7\xda\x29\x7a\x13\x42\x59\x4f\xff\xe6\x7a\xb1\x86\xfd\x0a\x00\x00\xff\xff\xaa\x3a\xef\x5e\x20\x06\x00\x00")

func template_genny_list_go() ([]byte, error) {
	return bindata_read(
		_template_genny_list_go,
		"../template/genny/list.go",
	)
}

var _template_genny_node_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x56\xdf\x4f\xe3\x38\x10\x7e\x8e\xff\x8a\xb9\x3e\xac\x6c\x54\xd1\xbb\x57\xa4\xea\xc4\xb2\x14\x55\x82\x72\x2a\xdd\xdb\x93\x10\xaa\x42\xea\x04\x1f\x89\x9d\x75\x1c\x0e\x2e\x9b\xff\xfd\x34\xb6\xf3\xa3\x69\xca\x72\x2f\x84\x8e\xbf\xf9\xe6\xf3\xcc\x78\xec\x3c\x8c\x9e\xc3\x84\xc3\xf7\x92\xeb\x37\x42\x44\x96\x2b\x6d\x80\x92\x60\x92\x08\xf3\x54\x3e\x9e\x46\x2a\x9b\x3d\x87\xff\x96\xb3\xf8\xb1\x78\xe2\x69\xce\xf5\xcc\x62\x67\x8f\x61\xc1\x27\x1f\xc0\xa5\x2a\x99\x10\x46\xc8\xec\x84\x00\x40\x56\x16\x06\xa2\x30\x4d\xe1\x37\x30\x22\xe3\x05\xe4\x5c\xc3\x26\x7c\x4c\xf9\xac\x30\xba\x8c\x0c\x5d\xa9\x1d\x5f\x85\x19\x67\xe4\x64\x46\x88\x79\xcb\x39\x34\x26\x70\x10\xa8\x48\x70\x82\xf1\x4f\x2f\x54\x96\x29\x89\xcb\xa4\x26\x24\x2e\x65\x04\x3c\xcb\xcd\x5b\xe3\x40\x19\x9c\xb4\xce\x15\x09\x34\x37\xa5\x96\xf0\xa9\xb1\x55\x1d\xc1\x19\x7c\x1a\x50\x56\x75\x8d\xac\x2f\xa1\x6e\x05\x6c\x97\xbb\xd7\x8d\xda\xa0\xa6\x2c\xcc\xef\x85\x34\x0f\x42\x1a\x98\xf7\x7f\x55\xf5\x11\x97\x2b\xad\xca\xfc\x7f\xfa\x59\xe1\x0d\xa8\x30\x5a\xc8\xa4\xe7\xe5\x0c\x43\x47\xfc\xb3\x51\xcb\xdd\xab\xc5\x39\x4c\x2f\x5c\x67\xa8\xfc\xee\xbe\x7c\xbd\xb9\xf9\x63\xdb\xf8\x2f\x8b\x3b\x97\xe5\x47\xa5\x52\x98\x83\xcd\xca\x1d\x37\xdd\x9a\x8c\x0c\x9d\x34\xf0\xc9\xd4\x21\x36\xea\xb3\x52\x29\x9d\x34\xee\x13\xc6\x7c\x45\xd0\xd7\x83\x17\x82\xa7\xbb\x82\x4a\xfc\x7f\x0a\xb1\xff\xd8\x7c\x3a\x59\x68\x2c\x33\x10\xd2\x30\x17\xbf\x22\x24\xb0\xfc\x6b\xfe\xbd\xe4\x85\xb9\xe3\xc6\x08\x99\xbc\x4f\xe6\x58\x18\x21\x01\x97\x65\xb6\x30\xd6\xa4\x9e\xe1\xcc\xef\xc6\x65\x08\x91\xf7\x16\xff\x40\x02\x11\x23\xa0\x22\x41\xd0\x68\xbd\xe3\xa6\x2d\x1d\x45\xbe\x29\xb4\x6c\x8c\x04\x35\x09\x66\xb3\xc5\xf2\xaf\x9b\xcb\x33\x24\x15\x11\xd8\x56\x55\x32\x7d\x83\xc2\x28\xcd\x7f\x27\xe4\x80\xcb\xf6\xa4\xe3\xb2\x81\x51\x62\xa2\x73\x14\xd6\x40\xaf\xb8\x69\x9b\x85\x7a\x50\x9f\xa7\xb7\x68\x79\x12\x9d\xb3\x5e\xa4\xae\x6d\xee\x71\xfd\x01\xe6\x2e\x52\x1f\xd2\x36\xc8\xbd\x4d\x9a\xc5\xac\xca\xcc\x27\xda\x97\xda\x02\xf6\xca\x7c\xe8\x8f\x81\xfd\x91\x32\xba\xe4\x84\xd4\xae\xe2\xa3\xdb\x16\x58\xd5\x29\x14\xbe\xd2\x0c\x93\x3d\xa2\x5a\xa0\x9c\x82\x74\x62\x3a\x86\x51\x31\xed\x32\x6b\x47\xc0\x68\x05\x9f\xa7\xf0\xb2\xd7\x57\xc1\xe1\x09\xbd\x7f\xc6\xe0\x2f\x83\xd8\xd6\xfd\x78\x6c\x5f\xa2\x7e\x1e\x46\x84\x74\x65\xfb\x88\x10\x8b\x3c\x50\xb3\xb7\xfa\x13\x45\x16\xf3\x13\x59\x7b\xad\xd6\xd5\x85\x6a\x5e\x94\xa9\x71\x1a\xbb\xa1\x69\x55\xec\xbb\x74\x29\xa7\x52\xed\xba\x39\xcd\x20\x6a\xa7\x28\x8e\xe0\xc1\x60\x45\x52\x11\x03\xba\xf4\xad\xf3\x39\x48\x61\x13\x12\xa4\x2a\x39\xbd\x56\x09\xb5\xdf\xdb\xab\xed\xb7\xf3\xf5\x6a\x0a\x18\x89\x32\xf0\x8b\xe7\x3a\x29\x2c\xb8\x11\x88\xf6\x05\x9d\xf4\x18\xa5\x32\x10\xab\x52\xee\x5a\x65\x13\x46\x82\xa0\xc6\xe3\x0b\x3c\x2d\x38\x88\x18\x52\x2e\xe9\x40\xca\xa9\xdb\xc4\x7c\x0e\xbf\xc2\x8f\x1f\xa3\x88\xae\xf4\x0e\x86\x42\xc6\x58\x60\x0e\x5d\xa1\x46\x30\xdd\x85\x32\x1f\xa9\xe2\x7b\x0e\xee\x3a\x19\xf3\xb2\x2b\x76\x42\xf9\xcc\x0c\x38\x9a\x63\x3a\x2c\xda\x1d\x0f\x75\xf4\xb4\x94\xb1\xa2\xb9\x2a\xdc\x71\x8d\x7d\xe1\xd7\x3c\x5a\xc8\x29\x44\x4a\xee\x16\xde\x74\x61\xff\x67\x76\x44\xdb\x08\xfd\xa2\x9f\xee\x93\x21\x51\xe3\xcd\xc8\xd1\xb6\xb1\x70\x06\x54\xc8\x58\xb9\x20\x68\x71\x21\xfa\x7b\xd9\x8b\xe4\x9c\xde\x21\x2d\xae\x79\x18\x9f\x1b\xfa\xf7\xe0\x5e\x39\xce\xd8\x7a\xbc\x43\x7b\xa1\x4a\x69\x6e\x63\x7b\x0b\x51\x86\xd4\xbd\xc3\x82\x2d\x33\x36\x26\x8e\x91\xfd\x19\xa6\x25\xb7\x1b\x11\x5e\x24\x6e\xbe\xb5\xf6\x98\x0f\xd5\xf6\x7c\x8f\x07\xb0\x32\xcf\x0d\x15\xbb\x57\x17\x60\xec\x50\x1e\x0d\xd1\xf3\xb6\x21\xf6\x5e\x65\xdf\x84\x79\xba\xd4\xba\xff\x38\x6b\x96\x48\x80\x0b\x5c\x6b\xa5\x0f\x87\xa2\x90\x49\xca\x9d\xce\xd6\x61\x0a\xdc\xc1\xd9\x01\x7d\xa7\x6f\xb0\x52\x35\xbf\xcf\xac\xf2\x29\x5c\x6a\x7d\x06\xbc\xee\x22\xf2\x7f\x8e\x3d\x08\x6d\x12\xce\xf3\x3c\x7d\xf3\xef\x8b\xde\xc3\x82\xb9\xb6\xc6\xcb\x79\xf0\xa8\xf4\xfd\x8e\x86\x6b\x51\xe0\xc3\xca\xbd\x1d\x1b\x43\x55\x93\x8f\x0c\x03\x0b\x59\x4a\x61\x28\xdb\xef\xc6\xf7\xab\x68\x8f\x77\x41\xd9\xde\x53\xf2\x30\x3b\xc3\x79\x70\x8c\x73\xad\x94\xc1\x33\x87\x5f\xf7\x78\xf2\x05\x70\x53\xfa\x17\xaf\xd2\xc3\x70\xce\xf9\x40\x8d\x43\x55\x4f\xed\xec\xbd\x5c\xaf\xb7\xab\xdb\xed\x72\x75\x71\xfd\xf5\xcb\xe5\x76\x7d\x7b\xbb\x71\x43\x48\x29\x83\x49\x34\xca\x72\x58\xbe\xcf\x61\xd1\xbb\x34\x11\x31\xc5\xe1\x4f\x6a\xf2\x5f\x00\x00\x00\xff\xff\xba\x67\xd3\xd4\x92\x0c\x00\x00")

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

var _template_genny_root_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x98\x5f\x6f\xa3\x38\x10\xc0\x9f\xf1\xa7\xf0\xf2\x50\x99\xdb\x88\x6c\xfe\x08\x55\xab\xcd\xc3\xf5\xa4\xde\x56\x5a\xa5\xa7\xee\xae\xf6\xe1\x74\x0f\x90\x9a\x04\x35\xc1\x1c\xd8\x6d\x69\x95\xef\x7e\x1a\xdb\x18\x43\xd2\x64\x39\x71\xdc\x4b\x85\xc7\xe3\xf9\xf3\x9b\xf1\x40\x93\x85\xab\x87\x70\x4d\xf1\xdf\x82\xe6\x25\x42\xc9\x2e\x63\x39\xc7\x04\x39\x6e\xc2\x5c\xe4\xb8\x05\xcf\x57\x2c\x7d\x74\x11\x72\xe2\x6d\xc8\x23\x11\xc7\x34\x2f\xb0\xbb\x4e\xf8\x46\x44\xfe\x8a\xed\xc6\x6b\xc6\xd6\x5b\x3a\xb6\xb6\xc7\x6b\x79\xd6\xd2\x79\x08\x5f\xc4\x38\x8e\x8a\x0d\xdd\x66\x34\x1f\x4b\x6f\xe3\x28\x2c\xa8\x8b\x3c\x84\x78\x99\x51\x7c\xc5\xd8\x16\x2f\x30\x08\x7d\x78\xd6\xd2\x92\x53\x23\x2d\x39\x55\xd2\x9b\x94\x5f\x56\x52\x78\x36\xd2\x49\x60\x89\x27\x81\x91\xcf\xa6\x96\x7c\x36\x35\xf2\x60\x6e\xc9\x83\xb9\x92\x7f\x4f\x2c\xf3\x72\x51\xcb\x6b\x07\x6a\x55\xef\xd4\x2e\xd4\xaa\xde\xa9\x9d\xa8\x95\xda\xb9\xde\xb2\xd0\x3a\xa4\x97\xd6\x5e\x7d\x4c\x2f\x35\xa8\xdf\xd8\x6e\xc7\xd2\x25\xbb\x37\x60\x6a\xc9\x81\xca\x97\xa4\xe0\xc7\xd4\x62\x91\xae\x30\xdd\x65\xbc\xac\xa5\xc4\xc3\xbf\xb4\x14\xf1\x2b\x72\x72\xca\x45\x9e\xe2\x8b\xd6\xd6\xeb\x1e\xed\xb5\x9d\x25\x7d\xfa\x59\x2b\x07\x2e\x8d\x91\xd6\x0e\x04\x0e\xa6\x5a\xa9\x58\xf1\x34\x77\xde\x0a\xe7\xac\x9d\xa3\x7e\x8d\xad\xdb\x8c\xa6\x24\xc7\x09\xf3\xef\x68\x78\x4f\xf3\x11\x5e\x85\x19\x4e\x52\xee\xe1\x3b\xc6\xf8\x37\xa0\x2d\x6d\x15\x62\xcb\xf1\xc7\x85\x91\xbe\xd6\x26\x3f\x2a\xfc\xca\x94\x34\xe0\xed\xab\x23\x16\x25\x7f\x19\xee\xa0\xa4\x6e\x65\xc2\x45\x8e\x3c\xf8\x95\x72\x10\xc1\x36\x39\x7e\xca\x33\xe6\xae\x29\x5f\x6d\x6e\xd2\x7b\xfa\x4c\x3c\x93\xa2\xda\x43\xce\x78\xac\x05\x3f\x1b\xa4\x4d\xe1\xaa\xbc\x12\x31\x89\x44\x8c\xff\xfc\x2b\x2a\x39\xfd\x77\x00\x8c\x95\xff\x07\x41\x95\x10\x67\x60\x8e\x44\xba\x55\xaf\xc2\xa2\x95\xcf\x4a\x5a\x87\x7c\x2e\x1a\x2d\xaf\x37\x7c\xeb\x6e\xa9\x8b\x61\xf5\xa1\xad\x53\xdd\xbd\x25\x7d\x92\xfd\x1e\x8d\xa0\x79\x88\x35\x2c\xfd\xdf\x29\xff\x7e\x1b\xc7\x05\xe5\xdf\x48\xe4\xdf\x91\x0f\x9e\xe7\x41\xe0\x8c\xbd\xcd\x53\x79\xd8\x2b\xad\x33\x00\xa5\xca\x71\x28\x8c\xd5\x48\xe0\xc4\x75\xce\x76\xca\x18\x59\x1d\xdc\xe2\x76\xc1\x9b\x9d\x84\x1c\xa7\x11\xa0\x0f\x3b\xda\x96\x37\x42\x4e\xdd\x4c\x24\x05\x2e\xd5\x41\x0f\x2f\xe9\x33\xdc\x50\xdb\x78\xc1\xc3\x5c\x26\x0f\xaa\xfe\x17\x9a\x12\x0f\x21\x27\x89\xcd\x1a\x9a\xc8\xc3\x9f\xb0\x52\x04\xe7\x3a\x1e\x50\x00\x67\xc8\x49\xe9\x13\x94\xd5\x58\xb9\x52\x65\x78\xe6\xf0\x40\xe4\x41\x30\x7a\x80\xf9\x10\x2a\x5e\x1c\x19\x5b\x4a\xe9\x64\x1b\x18\x8d\x76\x13\xe8\xd0\x4e\xb7\x82\x56\xb2\x1a\xe2\x48\xd1\xda\x2c\x6f\xd2\x82\xe6\x1c\xe8\x64\xac\x18\xe1\x22\x79\xa1\x6a\x58\xbd\x22\x47\x52\xb0\x5a\xe5\x88\xae\xf7\xa6\xe1\xcf\x61\xa1\xeb\x14\xc1\x7b\xfa\x15\x21\x9b\xb8\x29\xc9\xfb\x39\xfe\x64\x17\xed\x2d\x73\x3f\x12\xbe\xf9\x2c\x07\x6a\xab\xf2\x49\x8c\xdf\xa5\x2a\x3a\x79\x43\x65\xe4\xce\x46\xaa\x42\x99\x76\xe1\x03\x25\x6a\x00\x8d\xf0\xa5\x87\x1c\x27\x63\x85\xa9\xb1\xcc\xeb\x0f\x56\xe0\xf7\xf8\x12\x39\xf6\x17\x8b\xff\x23\x4f\x38\x55\x6f\x66\xb2\xd1\xa3\x5c\xa8\x65\xc6\x0a\xe0\x5b\x01\xca\x4a\x52\x55\x4b\xb6\x8a\xd2\xf6\x46\xf8\xc3\x08\x5f\xaa\xbf\xd0\x38\x4e\xcb\xe3\x02\x5c\xee\x1b\x58\x4c\xfa\xbf\x72\x96\x90\x02\x17\x3c\x4f\xd2\xb5\x87\x49\x92\xf2\x11\xa6\x79\xce\x72\x4f\xa2\x4c\xe5\x0a\xd2\xd0\xdf\x5b\xbe\x3a\xe1\x49\x20\xb0\xf5\x6e\x81\xd3\x64\x8b\x2f\x2e\x70\x81\x17\x0b\xec\x5e\x27\x74\x7b\xbf\x14\x3b\x57\xf2\x01\x0d\xa9\xd0\x88\x40\x1a\x35\x31\xc0\xdd\x86\xef\x2a\xf2\x28\x4b\xd8\x78\x25\xca\x20\xac\x91\x77\xf8\xaa\xef\x34\xf7\xf4\x0c\x02\x6f\xee\xa9\x69\x38\x6d\x82\x6e\xd4\x56\x8d\xfd\xe4\x85\xde\xc6\x32\x5c\xc5\x9f\xe7\x02\x26\xbd\x65\x52\xea\x54\x76\x6b\x7d\xa3\xf3\x95\x72\x95\x75\x7d\x83\xd4\x4e\x13\x4c\xc9\x29\x80\x91\xef\xb5\xff\x1e\x4c\xc9\x69\x2f\x60\x20\xdc\x2e\x60\xe0\x03\xda\x06\x23\xb3\x3e\x09\x06\xbe\xae\xc9\x23\xcc\x90\xcb\x01\xc0\x80\xb7\x3e\xc0\xc8\x70\x3b\x80\x91\xdf\xf8\x16\x18\x95\xf5\x39\x30\x93\x40\x91\x99\x04\xc3\xa0\x99\x04\x3d\xb1\x99\x04\x1d\xe1\x4c\x82\x16\x1d\x48\xfd\x1c\x9e\xd9\x54\xe1\x99\x4d\x87\xc1\x33\x9b\xf6\x84\x67\x36\xed\x88\x67\x36\x6d\xe1\x81\xd4\xcf\xe1\x09\xe6\x0a\x4f\x30\x1f\x06\x4f\x30\xef\x09\x4f\x30\xef\x88\x27\x98\xb7\xf0\x40\xea\x27\xf1\xc8\x7f\xba\xc9\xa3\x7c\x41\x0f\x31\x77\xa4\xbf\x3e\xf0\x88\xae\x93\x47\xb4\x47\x8f\x4e\xfd\x2c\x1e\x39\x7c\xc4\x50\xd3\x47\x79\xec\x8b\x50\xb7\xf9\x23\x0e\x06\x50\x95\xff\x59\x48\x72\x04\x89\xa1\x66\x90\xf2\xd8\x17\xa4\x6e\x53\x48\x1c\x8c\xa1\x2a\xff\xb3\x90\xe4\x20\x12\x43\x4d\x22\xe5\xb1\x2f\x48\xdd\x66\x91\x38\x18\x46\x55\xfe\x27\x21\xe9\x9f\xe6\xc8\x23\x8e\xd5\xd3\x00\x98\xb4\xcf\x3e\x38\x55\x41\x77\x00\xa5\x8f\xd8\xa4\x0c\x84\xf3\xa8\x64\x43\xc5\xea\x69\x28\x54\xfd\xb4\x54\x15\x74\x57\x54\xcd\xa6\x32\x10\x8e\xa0\xfa\x27\x00\x00\xff\xff\x75\x4c\x34\xbd\x70\x17\x00\x00")

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
	"../template/genny/basic.nogo": template_genny_basic_nogo,
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
				"basic.nogo": &_bintree_t{template_genny_basic_nogo, map[string]*_bintree_t{
				}},
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
