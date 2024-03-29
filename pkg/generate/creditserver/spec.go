// Package creditserver provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package creditserver

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

	"H4sIAAAAAAAC/+xUy27rNhD9FWHapWS9XCDVLnGApJs0SJpVkQVNjW2lEsmSI6NuoH8vSMqOZStO0JuL",
	"i/sABEjicF7nzJxn4LJRUqAgA8UzGL7ChrnPOzQ008gIZxrLim7ZpkFB1qS0VKipQneRt4Zkg/q30v7h",
	"P6xRNUKRZnkIC6kbRlBAJSjPIATaKPS/uEQNXQhSl6gvGeHAG7IkS+PEPkGaFGlapHvuhnQlljvvGzn0",
	"TdM0jbIsy6I8z/NoOp1Ox3xJEqtvdcWHqdNfz5J3VN6FoPHvttJYQvHnrpD9hsJ9bAb5Hnfx5PwJOdly",
	"esBPQc04R0XjYEWJfYIsLyxeo2B59/+J1vdL8yhTD6p8ezVO8fWOrn/w9clruYPweOPs1UosXJFUkStR",
	"sJaqeY0RUypqTKQ8sxF3PEMIa9SmksJ2M0kmSXR/c357f/37H65phYKpCgrIJ8kkhxAUo5UbhLgPYCdE",
	"Gvcu0XBdKfLRvNIGs20eO0XM2ix9vRV8a2joQpYbJ71SUD94TKm64s4lfjI25lbN7dfPGhdQwE/xi9zH",
	"vdbHrwl9N8SSdIvuwCgpjJ/vLEk+uozDAoYwmZZzNGbR1sEOIQv9NJkeYyokBQvZitJNhWmbhunNCNbE",
	"lsbOyksF96jXTqa7EFQ7Qpff/tfo8tbPSNeY+HyzdB1ifYquLtyuWvy8Xf3Op6rRa9kw6aU7D1jAx6n0",
	"drfKmjVIqG3uwyg7lbGKAoXbewhBsAaH1iE/4R7WBzLYPR5xNwLZh+DbQ+ABCOabYK/e06uxxJHVuENq",
	"tTABC0wlljW+huwV0sXm/CXVlwf4a1yOK6QhcwrfZs5uyQpZTat/baa3aCRGrQm4LPGIxOs+yjiW7+v3",
	"lyQ/vtwXGzwItmZVzeY1Hqu47Xp70VcSzFbI/9pr3x+7073+bSDU6+2stbqGAlZEqojjWnJWr6Sh4iw5",
	"S2LoHrv/AgAA//8BlnMUJg0AAA==",
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
