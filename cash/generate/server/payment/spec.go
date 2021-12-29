// Package generate provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package payment

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXTW/jNhD9KwLbW+3owz6kumUT7LaXrRGnpyIHihrbXEgklxwGFQL/94KUbEmWbCtZ",
	"dIEFFhBgeciZ4bz3PEO/EiZLJQUINCR9JYbtoKT+9YEqvbaZYZpn4AxKSwUaOfhlZTNjM0FLv5aD26eQ",
	"S0HSZo3MCPxLS1VA14SVct8Nai62ZD8jWlociVGbuyFCpoF62yAGSsXZMEZt7sY4WE4i7I8WmX0Bhi7m",
	"Ixhc0aoEgcPqKWOg8LN07230OI7jeZIkyXyxWCzmy+VyOXZaZg3KEvSfec89ThYzspG6pEhSwgUuktab",
	"C4QtaOcudQ76gdagtcmTKInDyD1BHKVxnMbJWHLv/c6De981UrSm7x/FY9tVjd47k2lgwBWOFzqP3BMk",
	"i9TVmoxrAmmx0pz13ePfb6MJMJ9TxN8qpwg/dTGJ6At7n7z9uoa+mUUNXy3XkJP0n84523q7uPUo6OXu",
	"n/t5IA6XiYuNBwA5+oIEtcizAuZUqXnjTmbkBbSp21N8E91E8/Xnu9X6j7+ePAkKBFWcpGRxE90kPi3u",
	"vLLCnCodmm5D3gIOe17ea9s+pKZuzanqpKk7dIySwtTiTaLIfTApsBE3VargzLuHX4yLfxgR7o0jlN7x",
	"Vw0bkpJfwnaYhM0kCfsZ258V1ZpWNXD9AoxlDIzZ2CI4nt0zaWxZUl2N1Ih0axy/nZa9Bv3iiHt2nmGD",
	"fngWtUdAq4UJWpr6uH3kIvdsaFoCgnbpTmP0tMOd5asFXZEZqadkf0ML4wQVn6bykg1cJwo2WpZn8h2F",
	"/bHe06YcjL8LCVBeC/8kLwZ//h4y607rd4tsRpbRcigOITHYSCvyExk6URwUY67JcEaUHRFePU+C1Rnd",
	"1cukbmNg8IPMqzeBdw2z/jzb9zsmagv7b2RvMmn/E0kDgCc3i+a66Wa8NCPU3fv1Q+RgI3Wwvvt0N+Dw",
	"/nBtncYhK6TN4cUB9dsQzmOzyKqxq3BD4JCvN/fY89VNBlCDG473V2B89LuCqWg+doP+YJheLXUytK/H",
	"q8y+PkkBY3+iHrw9oGfnWr3h2mTr3pt4bcBdOwi6y/3m8baxsJwK6htbQINCc8wgq4Luia907YvXBRoY",
	"LrYFBH4WDtD9BPihWnVyXUR59X1Q/rHat5+x9W0kq4LVVN58ENAvB6StLkhKdogqDcNCMlrspMH0NrqN",
	"QrJ/3v8XAAD///L1IF8IEQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
