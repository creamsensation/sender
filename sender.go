package sender

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	
	"github.com/creamsensation/auth"
	"github.com/creamsensation/util/constant/contentType"
	"github.com/creamsensation/util/constant/dataType"
)

type Send interface {
	Header() http.Header
	Status(StatusCode int) Send
	Error(e any) error
	Json(value any) error
	Html(value string) error
	Xml(value string) error
	Text(value string) error
	Bool(value bool) error
	Redirect(url string) error
	File(name string, bytes []byte) error
}

type Sender struct {
	Bytes       []byte
	DataType    string
	ContentType string
	Value       string
	StatusCode  int
	request     *http.Request
	response    http.ResponseWriter
	auth        auth.Manager
}

func New() *Sender {
	return &Sender{}
}

func (s *Sender) Header() http.Header {
	return s.response.Header()
}

func (s *Sender) Status(StatusCode int) Send {
	s.StatusCode = StatusCode
	return s
}

func (s *Sender) Error(e any) error {
	var err error
	switch v := e.(type) {
	case nil:
		return s.Bool(true)
	case string:
		err = errors.New(v)
	case error:
		err = v
	default:
		err = errors.New(fmt.Sprintf("%v", e))
	}
	bytes, err := json.Marshal(Error{Error: err.Error()})
	s.Bytes = bytes
	s.DataType = dataType.Error
	s.ContentType = contentType.Json
	if s.StatusCode == http.StatusOK {
		s.StatusCode = http.StatusBadRequest
	}
	return err
}

func (s *Sender) Json(value any) error {
	bytes, err := json.Marshal(Json{Result: value})
	s.Bytes = bytes
	s.DataType = dataType.Json
	s.ContentType = contentType.Json
	return err
}

func (s *Sender) Html(value string) error {
	s.Bytes = []byte(value)
	s.DataType = dataType.Html
	s.ContentType = contentType.Html
	return nil
}

func (s *Sender) Xml(value string) error {
	bytes, err := json.Marshal(Json{Result: value})
	s.Bytes = bytes
	s.DataType = dataType.Xml
	s.ContentType = contentType.Xml
	return err
}

func (s *Sender) Text(value string) error {
	bytes, err := json.Marshal(Json{Result: value})
	s.Bytes = bytes
	s.DataType = dataType.Text
	s.ContentType = contentType.Json
	return err
}

func (s *Sender) Bool(value bool) error {
	bytes, err := json.Marshal(Json{Result: value})
	s.Bytes = bytes
	s.DataType = dataType.Bool
	s.ContentType = contentType.Json
	return err
}

func (s *Sender) Redirect(url string) error {
	s.Value = url
	s.DataType = dataType.Redirect
	return nil
}

func (s *Sender) File(name string, bytes []byte) error {
	s.Value = name
	s.Bytes = bytes
	s.DataType = dataType.Stream
	s.ContentType = contentType.OctetStream
	return nil
}
