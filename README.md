# hash2xml
###### Simple XML generation from a Golang map[String]interface{}

Nobody uses XML anymore, but sometimes you just have to.

I was surprised to find out it is not possible to generate XML from a Go map using the encoding/xml
package. It is however possible to marshal a map to JSON using encoding/json. Until the encoding/xml
package decides to support marshalling of a map, I will be using this.

Although there are other tools available, they have lots of dependencies and other functionalities
that I just don't need. This is meant to be easily embeddable without any other dependencies.

As they say, necessity is the mother of invention.


#### Basic Usage
```go
hash := make(map[string]interface{})
hash["key1"] = 1
hash["key2"] = "2"
hash["key3"] = []interface{}{"Array value 1", "Array value 2"}
hash["key4"] = []interface{}{1, 2, 3, 4, 5, 6, 7}

bytes, _ := hash2xml.ToXML("docroot", hash)
log.Printf(string(bytes))
```
prints out

```xml
<docroot>
 <key2>2</key2>
 <key3>
  <string>Array value 1</string>
  <string>Array value 2</string>
 </key3>
 <key4>
  <int>1</int>
  <int>2</int>
  <int>3</int>
  <int>4</int>
  <int>5</int>
  <int>6</int>
  <int>7</int>
 </key4>
 <key1>1</key1>
</docroot>
```

#### Addvanced Usage
```go
// hash containing user defined types
hash := map[string]interface{}{
  "key1": "value of type string",
  "key2": myType{MyInt: 42, MyString: "Hallo world"},
}

// create a new byte buffer and serializer
var b bytes.Buffer
serializer := hash2xml.NewSerializer(&b, " ", true)

// add a custom converter for myType
serializer.AddConverter(func(s *hash2xml.Serializer, raw interface{}, path string, key ...string) (bool, error) {
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
    return false, nil
  }
})

serializer.Encode("docroot", hash)

log.Printf(string(b.Bytes()))
```
prints out

```xml
<key1>value of type string</key1>
 <key2>
  <myint>42</myint>
  <mystr>Hallo world</mystr>
 </key2>
</docroot>
```

#### TODO
* Formatting for dates
* Custom tag names for scalar types
* Attributes?
