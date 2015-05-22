package hash2xml

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"reflect"
	"time"
)

var (
	header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// Printer is responsible for writing bytes to its writer
type Printer struct {
	*bufio.Writer
	Spacer string
	Pretty bool
	depth  int
}

// Serializer is resposible for converting a hash to XML
type Serializer struct {
	Printer
}

// ToXML is a factory/helper method that converts a map to an XML document
func ToXML(rootName string, hash map[string]interface{}) ([]byte, error) {
	var b bytes.Buffer

	serializer := Serializer{
		Printer: Printer{
			Writer: bufio.NewWriter(&b),
			Spacer: " ",
			Pretty: true,
		},
	}

	err := serializer.encode(rootName, hash)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// encode is resposible for doing the actual work of writing
// a hash as xml into a Writer
func (s *Serializer) encode(rootName string, hash map[string]interface{}) error {
	s.WriteString(header)
	s.encodeHash(hash, rootName)
	err := s.Flush()

	return err
}

// encodeHash recurses a hash and converts key/values into XML
func (s *Serializer) encodeScalar(value interface{}, key ...string) {
	var defaultKey string

	if len(key) > 0 {
		s.writeStartTag(key[0])
	} else {
		defaultKey = s.getDefaultKey(value)
		s.writeStartTag(defaultKey)
	}

	s.writeScalar(value)

	if len(key) > 0 {
		s.writeEndTag(key[0])
	} else {
		s.writeEndTag(defaultKey)
	}
}

// encodeHash recurses a hash and converts key/values into XML
func (s *Serializer) encodeHash(hash map[string]interface{}, key ...string) {
	if len(key) > 0 {
		s.writeStartTag(key[0])
		s.newline()
		s.indent()
	}

	// recursively serialize the hash
	for k, v := range hash {
		s.recurse(v, k)
	}

	if len(key) > 0 {
		s.dedent()
		s.writeIndentation()
		s.writeEndTag(key[0])
	}
}

// encodeHash recurses a hash and converts key/values into XML
func (s *Serializer) encodeArray(array []interface{}, key ...string) {
	if len(key) > 0 {
		s.writeStartTag(key[0])
		s.newline()
		s.indent()
	}

	// iterate the array and serialize all the values
	for _, value := range array {
		s.recurse(value)
	}

	if len(key) > 0 {
		s.dedent()
		s.writeIndentation()
		s.writeEndTag(key[0])
	}
}

// recurse is delegating enttries of a hash to the correct encoding method
func (s *Serializer) recurse(raw interface{}, key ...string) {
	switch v := raw.(type) {

	// map
	case map[string]interface{}:
		s.encodeHash(v, key...)

	// arrays
	case []interface{}:
		s.encodeArray(v, key...)

		// scalar
	case string, float64, bool, int, int32, int64, float32, time.Time:
		s.encodeScalar(v, key...)

	default:
		log.Printf("XML serializer not supporting type %#v", v)
	}
}

func (s *Serializer) getDefaultKey(value interface{}) string {
	return reflect.TypeOf(value).String()
}

// vvvvvvv Printer methods vvvvvvv

// indent increases current depth
func (p *Printer) indent() {
	p.depth++
}

// dedent decreases current depth
func (p *Printer) dedent() {
	p.depth--
}

// writeIndentation prints the current depth of spaces
func (p *Printer) writeIndentation() {
	if p.Pretty && p.depth > 0 {
		for i := 0; i < p.depth; i++ {
			p.WriteString(p.Spacer)
		}
	}
}

// newLine does indeed print a newline character
func (p *Printer) newline() {
	if p.Pretty {
		p.WriteByte('\n')
	}
}

// writeStartTag creates an start tag
func (p *Printer) writeStartTag(name string) {
	p.writeIndentation()
	p.WriteByte('<')
	p.WriteString(name)
	p.WriteByte('>')
}

// writeEndTag creates an start tag
func (p *Printer) writeEndTag(name string) {
	p.WriteString("</")
	p.WriteString(name)
	p.WriteByte('>')
	p.newline()
}

// writeScalar creates an start tag
func (p *Printer) writeScalar(raw interface{}) {
	p.WriteString(fmt.Sprintf("%v", raw))
}
