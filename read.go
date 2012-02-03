package odf

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

const (
	MimeTypePfx = "application/vnd.oasis.opendocument."
)

type File struct {
	*zip.ReadCloser
	MimeType string
}

// Open an Open Document File, uncompress its ZIP archive,
// and check its mime type. The returned *File provides --
// via its Open method -- access to files containted in the
// ZIP file, like content.xml.
func Open(odfName string) (f *File, err error) {
	of := new(File)
	of.ReadCloser, err = zip.OpenReader(odfName)
	if err != nil {
		return
	}
	mf, err := of.Open("mimetype")
	if err != nil {
		of.Close()
		return
	}
	defer mf.Close()

	b, err := ioutil.ReadAll(mf)
	if err != nil {
		return
	}
	of.MimeType = string(b)

	if !strings.HasPrefix(of.MimeType, MimeTypePfx) {
		err = errors.New("not an Open Document mime type")
	} else {
		f = of
	}

	return
}

func (of *File) Open(name string) (f io.ReadCloser, err error) {
	for _, zf := range of.File {
		if zf.Name == name {
			return zf.Open()
		}
	}
	err = errors.New("odf: open " + name + ": no such file")
	return
}
