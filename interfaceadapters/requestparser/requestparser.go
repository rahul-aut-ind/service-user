package requestparser

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type (
	RequestParser struct {
		Log         *logger.Logger
		Body        []byte
		ContentType string
	}

	MultiPartData struct {
		Files  map[string]*Image
		Values map[string]*string
	}

	Image struct {
		Bytes []byte
		Ext   string
	}
)

var allowedExt = []string{".jpg"}

// ParseMultipart parses the multipart form data and returns MultiPartData
func (rp *RequestParser) ParseMultipart() (*MultiPartData, error) {
	_, params, err := mime.ParseMediaType(rp.ContentType)
	if err != nil {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("could not parse media"))
	}
	b := params["boundary"]
	if b == "" {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("could not parse media"))
	}

	reader := bytes.NewReader(rp.Body)
	multipartReader := multipart.NewReader(reader, b)
	formData, err := multipartReader.ReadForm(10 << 20)
	if err != nil {
		return nil, err
	}

	if len(formData.File) == 0 {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("no files found"))
	}
	if len(formData.Value) == 0 {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("no values found"))
	}

	res := &MultiPartData{}
	res.Files = make(map[string]*Image, len(formData.File))
	res.Values = make(map[string]*string, len(formData.Value))

	for k, v := range formData.File {
		// we only handle one file per key
		f, err := rp.getFileBytes(v[0])
		if err != nil {
			return nil, err
		}

		ext, err := rp.getExt(v[0].Filename)
		if err != nil {
			return nil, err
		}

		p := Image{
			Bytes: f,
			Ext:   *ext,
		}

		res.Files[k] = &p
	}

	for k, v := range formData.Value {
		// we only handle one value per key
		vs := v[0]
		res.Values[k] = &vs
	}

	return res, nil
}

// getFileBytes reads the file from the multipart form
func (rp *RequestParser) getFileBytes(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	if err != nil {
		rp.Log.Errorf("error opening file %v", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error opening file"))
	}
	fd, err := io.ReadAll(f)
	if err != nil {
		rp.Log.Errorf("error reading file %v", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error reading file"))
	}

	return fd, nil
}

/*
getExt returns the extension of a file
it uses regex to get the extension of a file following last (.) in the filename
*/
func (rp *RequestParser) getExt(filename string) (*string, error) {
	if filename == "" {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("filename is empty"))
	}
	re := regexp.MustCompile(`\.[0-9a-z]+$`)
	ext := strings.ToLower(re.FindString(filename))

	if ext == "" {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("no extension found in filename"))
	}

	if !contains(allowedExt, ext) {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("file extension not allowed"))
	}

	return &ext, nil
}

// contains checks if a string is in a list of strings
func contains(list []string, ext string) bool {
	for _, e := range list {
		if e == ext {
			return true
		}
	}
	return false

}
