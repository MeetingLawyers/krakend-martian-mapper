package mapper

import (
	"encoding/json"
	"github.com/google/martian"
	"github.com/google/martian/parse"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	parse.Register("mapper.JSONMapper", MapperFromJSON)
}

// MappingConfigJSON to Unmarshal the JSON configuration
type MappingConfigJSON struct {
	CopyFields map[string]string   `json:"copy_fields"`
	MapFields map[string]string    `json:"map_fields"`
	Scope  []parse.ModifierType    `json:"scope"`
}

// Mapping contains the private and public Marvel API key
type Mapping struct {
	copyFields map[string]string
	mapFields map[string]string
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *Mapping) ModifyRequest(req *http.Request) error {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	var bodyJson = make(map[string]string)
	err = json.Unmarshal(bodyBytes, &bodyJson)
	var bodyEmpty = false
	if err != nil {
		bodyEmpty = true
	}
	for actualKey, newKey := range m.copyFields {
		if query.Get(actualKey) != "" {
			query.Set(newKey, query.Get(actualKey))
		}

		if bodyEmpty {
			continue
		}
		bodyJson[newKey] = bodyJson[actualKey]
	}
	for actualKey, newKey := range m.mapFields {
		if query.Get(actualKey) != "" {
			query.Set(newKey, query.Get(actualKey))
			query.Del(actualKey)
		}

		if bodyEmpty {
			continue
		}
		if val, ok := bodyJson[actualKey]; ok {
			bodyJson[newKey] = val
			delete(bodyJson, actualKey)
		}
	}

	if !bodyEmpty {
		new_body_content, _ := json.Marshal(bodyJson)
		req.Body = ioutil.NopCloser(strings.NewReader(string(new_body_content)))
	}

	req.URL.RawQuery = query.Encode()
	return nil
}

// MapperNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func MapperNewModifier(copyFields map[string]string, mapFields map[string]string) martian.RequestModifier {
	return &Mapping{
		copyFields: copyFields,
		mapFields: mapFields,
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

	return parse.NewResult(MapperNewModifier(configByteSlice.CopyFields, configByteSlice.MapFields), configByteSlice.Scope)
}
