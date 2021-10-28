package querystring

import (
	"encoding/json"
	"github.com/google/martian"
	"github.com/google/martian/parse"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	parse.Register("querystring.JSONMapper", MapperFromJSON)
}

// MappingConfigJSON to Unmarshal the JSON configuration
type MappingConfigJSON struct {
	Fields map[string]string    `json:"fields"`
	Scope  []parse.ModifierType `json:"scope"`
}

// Mapping contains the private and public Marvel API key
type Mapping struct {
	fields map[string]string
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *Mapping) ModifyRequest(req *http.Request) error {
	log.Println("Request Modifier ----------------------")

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	log.Println("Body original:" + string(bodyBytes))
	log.Println("Query original:" + string(req.URL.RawQuery))

	log.Println("Start Modifier ------------------------")
	query.Set("chorizo", "yes")


	var bodyjson = make(map[string]string)
	err = json.Unmarshal(bodyBytes, &bodyjson)
	if err != nil {
		panic(err)
	}
	for actualKey, newKey := range m.fields {
		log.Println("Key from: " + actualKey + ". Key to: " + newKey)

		if query.Get(actualKey) != "" {
			query.Set(newKey, query.Get(actualKey))
			query.Del(actualKey)
		}

		bodyjson[newKey] = bodyjson[actualKey]
		delete(bodyjson, actualKey)
	}

	new_body_content, _ := json.Marshal(bodyjson)
	req.Body = ioutil.NopCloser(strings.NewReader(string(new_body_content)))
	log.Println("Body result: " + string(new_body_content))

	// Recibido por referencia (puntero), lo altera directamente
	req.URL.RawQuery = query.Encode()
	log.Println("Query result: " + req.URL.RawQuery)
	return nil
}

// MapperNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func MapperNewModifier(mappingFields map[string]string) martian.RequestModifier {
	return &Mapping{
		fields: mappingFields,
	}
}

// MapperFromJSON takes a JSON message as a byte slice and returns
// a querystring.modifier and an error.
// a body.modifier
//
// Example JSON:
// {
//  "public": "apikey",
//  "private": "apikey",
//  "scope": ["request", "response"]
// }
func MapperFromJSON(b []byte) (*parse.Result, error) {
	configByteSlice := &MappingConfigJSON{}

	if err := json.Unmarshal(b, configByteSlice); err != nil {
		return nil, err
	}

	return parse.NewResult(MapperNewModifier(configByteSlice.Fields), configByteSlice.Scope)
}
