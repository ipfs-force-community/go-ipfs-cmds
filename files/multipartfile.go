package files

import (
	"mime"
	"mime/multipart"
	"net/http"
)

const (
	multipartFormdataType = "multipart/form-data"
	multipartMixedType    = "multipart/mixed"

	contentTypeHeader = "Content-Type"
)

// MultipartFile implements File, and is created from a `multipart.Part`.
// It can be either a directory or file (checked by calling `IsDirectory()`).
type MultipartFile struct {
	File

	Part      *multipart.Part
	Reader    *multipart.Reader
	Mediatype string
}

func NewFileFromPart(part *multipart.Part) (File, error) {
	f := &MultipartFile{
		Part: part,
	}

	contentType := part.Header.Get(contentTypeHeader)

	var params map[string]string
	var err error
	f.Mediatype, params, err = mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if f.IsDirectory() {
		boundary, found := params["boundary"]
		if !found {
			return nil, http.ErrMissingBoundary
		}

		f.Reader = multipart.NewReader(part, boundary)
	}

	return f, nil
}

func (f *MultipartFile) IsDirectory() bool {
	return f.Mediatype == multipartFormdataType || f.Mediatype == multipartMixedType
}

func (f *MultipartFile) NextFile() (File, error) {
	if !f.IsDirectory() {
		return nil, ErrNotDirectory
	}

	part, err := f.Reader.NextPart()
	if err != nil {
		return nil, err
	}

	return NewFileFromPart(part)
}

func (f *MultipartFile) FileName() string {
	return f.Part.FileName()
}

func (f *MultipartFile) Read(p []byte) (int, error) {
	if f.IsDirectory() {
		return 0, ErrNotReader
	}
	return f.Part.Read(p)
}

func (f *MultipartFile) Close() error {
	if f.IsDirectory() {
		return ErrNotReader
	}
	return f.Part.Close()
}
