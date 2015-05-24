package hash2xml

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
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

// Converter is a function  in the xml conversion pipeline that converts an
// object into xml  by writing it to the serializer
type Converter func(*Serializer, interface{}, string, ...string) (bool, error)

// Serializer is resposible for converting a hash to XML
type Serializer struct {
	Printer
	converters []Converter
}

// ToXML is a factory/helper method that converts a map to an XML document
func ToXML(rootName string, hash map[string]interface{}) ([]byte, error) {
	var b bytes.Buffer

	serializer := NewSerializer(&b, " ", true)
	err := serializer.Encode(rootName, hash)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// NewSerializer instantiates a serializer with all fields set
func NewSerializer(b *bytes.Buffer, spacer string, pretty bool) Serializer {

	// default converters
	array := []Converter{
		hashConverter,
		arrayConverter,
		scalarConverter,
		timeConverter,
	}

	// instantiate the Serializer
	return Serializer{
		converters: array,
		Printer: Printer{
			Writer: bufio.NewWriter(b),
			Spacer: spacer,
			Pretty: pretty,
		},
	}
}

// AddConverter adds a new conversion function for the xml conversion pipeline
// This method actually prepends the converter, so that it take preference
// over the default ones
func (s *Serializer) AddConverter(c ...Converter) {
	s.converters = append(c, s.converters...)
}

// Encode is resposible for doing the actual work of writing
// a hash as xml into a Writer
func (s *Serializer) Encode(rootName string, hash map[string]interface{}) error {
	s.WriteString(header)
	err := s.Convert(hash, "", rootName)
	if err != nil {
		log.Printf("Hash2XML conversion error: %#v", err)
		return err
	}

	return s.Flush()
}

// Convert is delegating enttries of a hash to the correct encoding method
func (s *Serializer) Convert(raw interface{}, p string, key ...string) error {
	// update the path
	path := p
	if len(key) > 0 {
		path = fmt.Sprintf("%s/%s", p, key[0])
	}

	// look for a suitable converter
	for _, c := range s.converters {
		found, err := c(s, raw, path, key...)
		if err != nil {
			return err
		}
		// stop looking if a converter was found
		if found {
			return nil
		}
	}

	t := reflect.TypeOf(raw)
	log.Printf("Please add your own hash2xml.Converter that accepts type %s", t)
	return fmt.Errorf("Error: XML serializer did not find a converter for type: %v", t)
}

func (s *Serializer) getDefaultKey(value interface{}) string {
	return reflect.TypeOf(value).String()
}

// vvvvvvv Printer methods vvvvvvv

// Indent increases current depth
func (p *Printer) Indent() {
	p.depth++
}

// Dedent decreases current depth
func (p *Printer) Dedent() {
	p.depth--
}

// WriteIndentation prints the current depth of spaces
func (p *Printer) WriteIndentation() {
	if p.Pretty && p.depth > 0 {
		for i := 0; i < p.depth; i++ {
			p.WriteString(p.Spacer)
		}
	}
}

// GetIndentation returns the current indentation as a string
func (p *Printer) GetIndentation() string {
	if p.Pretty && p.depth > 0 {
		return strings.Repeat(p.Spacer, p.depth)
	}
	return ""
}

// Newline does indeed print a newline character
func (p *Printer) Newline() {
	if p.Pretty {
		p.WriteByte('\n')
	}
}

// WriteStartTag creates an start tag
func (p *Printer) WriteStartTag(name string, attributes ...string) {
	p.WriteIndentation()
	p.WriteByte('<')
	p.WriteString(name)

	// write some optional attributes
	for _, attr := range attributes {
		p.WriteByte(' ')
		p.WriteString(attr)
	}

	// close the tag
	p.WriteByte('>')
}

// WriteEndTag creates an start tag
func (p *Printer) WriteEndTag(name string) {
	p.WriteString("</")
	p.WriteString(name)
	p.WriteByte('>')
	p.Newline()
}

// WriteScalar creates an start tag
func (p *Printer) WriteScalar(raw interface{}) {
	p.WriteString(fmt.Sprintf("%v", raw))
}
