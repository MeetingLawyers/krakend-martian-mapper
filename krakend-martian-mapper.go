package querystring

import (
	"encoding/json"
	"github.com/google/martian"
	"github.com/google/martian/parse"
	"net/http"
)

func init() {
	parse.Register("querystring.JSONMapper", MapperFromJSON)
}

// MappingConfigJSON to Unmarshal the JSON configuration
type MappingConfigJSON struct {
	Scope   []parse.ModifierType `json:"scope"`
}

// Mapping contains the private and public Marvel API key
type Mapping struct {

}
// ModifyRequest modifies the query string of the request with the given key and value.
func (m *Mapping) ModifyRequest(req *http.Request) error {
	query := req.URL.Query()
	query.Set("chorizo", "yes")
	//	ts := strconv.FormatInt(time.Now().Unix(), 10)
	//	hash := GetMD5Hash(ts + m.private + m.public)
//	query.Set("ts", ts)
//	query.Set("hash", hash)

	// Recibido por referencia (puntero), lo altera directamente
	req.URL.RawQuery = query.Encode()

	return nil
}

// MapperNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func MapperNewModifier() martian.RequestModifier {
	return &Mapping{

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

	return parse.NewResult(MapperNewModifier(), configByteSlice.Scope)
}
