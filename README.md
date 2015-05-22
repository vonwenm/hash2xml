# hash2xml
###### Simple XML generation from a Golang map[String]interface{}

Nobody uses XML anymore, but sometimes you just have to.

I was surprised to find out it is not possible to generate XML from a Go map using the encoding/xml
package. It is however possible to marshal a map to JSON using encoding/json. Untill the encoding/xml
package decides to support marshalling of a map, I will be using this.

Although there are other tools available, they have lots of dependencies and other functionalities
that I just not need.

As the say,
> necessity is the month of invention


###### TODO
* Formatting for dates
* Custom tag names for scalar types
* Attributes?
