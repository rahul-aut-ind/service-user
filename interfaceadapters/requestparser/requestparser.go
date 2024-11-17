package requestparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
)

type (
	RequestParser struct {
		Body        []byte
		ContentType string
	}

	MultiPartData struct {
		Image    *Image
		Metadata *models.Metadata
	}

	Image struct {
		Bytes []byte
		Ext   string
	}
)

const (
	QueryParamBoundary = "boundary"
	ImageKey           = "image"
	MetadataKey        = "metadata"
	JPGImageExtension  = ".jpg"
)

var (
	allowedExt    = []string{JPGImageExtension}
	mediaExtRegEx = regexp.MustCompile(`\.[0-9a-z]+$`)
)

// ParseMultipart parses the multipart form data and returns MultiPartData
func (rp *RequestParser) ParseMultipart() (*MultiPartData, error) {
	_, params, err := mime.ParseMediaType(rp.ContentType)
	if err != nil {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("could not parse media %s", rp.ContentType))
	}

	boundary := params[QueryParamBoundary]
	if boundary == "" {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("could not parse media"))
	}

	formData, err := multipart.NewReader(bytes.NewReader(rp.Body), boundary).ReadForm(10 << 20)
	if err != nil {
		return nil, err
	}

	data := &MultiPartData{}

	if err := rp.parseFiles(formData, data); err != nil {
		return nil, err
	}
	if err := rp.parseValues(formData, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (rp *RequestParser) parseFiles(formData *multipart.Form, data *MultiPartData) error {
	if fileHeaders, exists := formData.File[ImageKey]; exists && len(fileHeaders) > 0 {
		fileBytes, err := rp.getFileBytes(fileHeaders[0])
		if err != nil {
			return err
		}
		ext, err := rp.getExt(fileHeaders[0].Filename)
		if err != nil {
			return err
		}
		data.Image = &Image{Bytes: fileBytes, Ext: *ext}
	} else {
		return errors.New(errors.ErrCodeBadRequest, fmt.Errorf("image not found, please check request params"))
	}
	return nil
}

func (rp *RequestParser) parseValues(formData *multipart.Form, data *MultiPartData) error {
	if metadataValues, exists := formData.Value[MetadataKey]; exists && len(metadataValues) > 0 {
		metadata := &models.Metadata{}
		if err := json.Unmarshal([]byte(metadataValues[0]), metadata); err != nil {
			return errors.New(errors.ErrCodeGeneric, fmt.Errorf("err unmarshalling metadata"))
		}
		data.Metadata = metadata
	} else {
		return errors.New(errors.ErrCodeBadRequest, fmt.Errorf("metadata not found, please check request params"))
	}
	return nil
}

// getFileBytes reads the file from the multipart form
func (rp *RequestParser) getFileBytes(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	if err != nil {
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error opening file"))
	}
	fd, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error reading file"))
	}

	return fd, nil
}

// getExt returns the extension of a file.
func (rp *RequestParser) getExt(filename string) (*string, error) {
	ext := strings.ToLower(mediaExtRegEx.FindString(filename))
	if ext == "" || !contains(allowedExt, ext) {
		return nil, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("file extension not allowed or not found"))
	}
	return &ext, nil
}

// contains checks if a string is in a list of strings
func contains(list []string, item string) bool {
	for _, e := range list {
		if e == item {
			return true
		}
	}
	return false
}
