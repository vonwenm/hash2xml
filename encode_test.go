package hash2xml_test

import (
	"bytes"
	"encoding/xml"
	"hash2xml"
	"log"
	"testing"
	"time"
)

// string type
type stringType struct {
	Value string `xml:"string"`
}

// int type
type intType struct {
	Value int `xml:"int"`
}

// sample type
type testType struct {
	Key1 int          `xml:"key1"`
	Key2 string       `xml:"key2"`
	Key3 []stringType `xml:"key3"`
	Key4 []intType    `xml:"key4"`
}

// sample type
type myType struct {
	MyInt    int    `xml:"myint"`
	MyString string `xml:"mystr"`
}

// Example usage
func TestBasicExample(t *testing.T) {
	hash := make(map[string]interface{})
	hash["key1"] = 123
	hash["key2"] = "hallo world"
	hash["key3"] = []interface{}{1, 2, 3}
	hash["key4"] = map[string]interface{}{
		"name":    "John",
		"surname": "Doe",
		"number":  "555-FILK",
	}

	bytes, _ := hash2xml.ToXML("docroot", hash)
	log.Printf(string(bytes))
}

// Simple test for checking if the XML is well formed
func TestDeserialize(t *testing.T) {
	hash := make(map[string]interface{})
	hash["key1"] = 1
	hash["key2"] = "2"
	hash["key3"] = []interface{}{"Array value 1", "Array value 2"}
	hash["key4"] = []interface{}{1, 2, 3, 4, 5, 6, 7}
	bytes, err := hash2xml.ToXML("docroot", hash)
	if err != nil {
		t.Fatalf("XML encoding error encountered: %#v", err)
	}

	// deserialize with encoding/xml
	temp := testType{}
	err = xml.Unmarshal(bytes, &temp)
	if err != nil {
		t.Fatalf("XML unmarshal failed: %#v", err)
	}
	log.Printf("Key1: %#v", temp.Key1)
	log.Printf("Key2: %#v", temp.Key2)
	log.Printf("Key3: %#v", temp.Key3)
	log.Printf("Key4: %#v", temp.Key4)
	log.Printf(string(bytes))
}

// Test for a hash containing various types
func TestLargeMap(t *testing.T) {
	hash := make(map[string]interface{})

	// scalars
	hash["key1"] = 1
	hash["key2"] = "2"

	// arrays
	hash["key3"] = []interface{}{
		"Array value 1",
		"Array value 2",
	}

	// maps
	hash["key4"] = map[string]interface{}{
		"MapKey1": "Map value 1",
		"MapKey2": "Map value 2",
	}

	// mixed
	hash1 := map[string]interface{}{
		"a": "hallo world",
		"b": "sawubona mhlaba",
	}

	hash2 := map[string]interface{}{
		"c": "another key",
		"d": 123,
	}

	array1 := []interface{}{
		map[string]interface{}{"EmbeddedArray1": "This is a string"},
		map[string]interface{}{"EmbeddedArray2": 2.343},
		map[string]interface{}{"EmbeddedArray3": 3},
		map[string]interface{}{"EmbeddedArray4": time.Now()},
		map[string]interface{}{"EmbeddedArray5": true},
	}

	hash["key5"] = map[string]interface{}{
		"EmbeddedMap1": hash1,
		"EmbeddedMap2": hash2,
		"EmbeddedMap3": array1,
	}
	hash["key6"] = []interface{}{1, 2, 3, 4, 5, 6, 7}

	bytes, err := hash2xml.ToXML("docroot", hash)
	if err != nil {
		t.Fatalf("XML encoder error encountered: %#v", err)
	}
	log.Printf("\n%s", string(bytes))
}

func TestCustomEncoder(t *testing.T) {
	// hash containing user defined types
	hash := map[string]interface{}{
		"key1": "value of type string",
		"key2": myType{MyInt: 42, MyString: "Hallo world"},
	}

	var b bytes.Buffer
	// create a new serializer with a reference to a byte buffer
	serializer := hash2xml.NewSerializer(&b, " ", true)

	// add a custom encoder for myType
	serializer.AddEncoder(func(s *hash2xml.Serializer, raw interface{}, path string, key ...string) (bool, error) {
		switch v := raw.(type) {
		case myType:

			// change the root element name of myType
			wrapper := struct {
				myType
				XMLName xml.Name
			}{v, xml.Name{Local: key[0]}}

			// delegate to encoding/xml
			out, err := xml.MarshalIndent(wrapper, s.GetIndentation(), " ")

			// embed the xml in the rest of the document
			s.WriteString(string(out))
			s.Newline()
			return true, err
		default:
			// return false to indicate that other encoders must handle it
			return false, nil
		}
	})

	// encode the hash
	err := serializer.Encode("docroot", hash)
	if err != nil {
		t.Fatalf("XML encoder error encountered: %#v", err)
	}

	log.Printf(string(b.Bytes()))
}
