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

	"H4sIAAAAAAAC/+xUTW/bOBD9K8LsHiXrywtkdUscINlLNkiaU5EDTY1tphLJkiOjbqD/XpCyHctWnKBN",
	"UfQDMGCJ1Hy9N+89Ale1VhIlWSgewfIF1sw/3qCliUFGODFYCrpmqxoluSttlEZDAv2HvLGkajT/le4N",
	"P7FaVwhFmuUhzJSpGUEBQlKeQQi00ti94hwNtCEoU6I5Z4S9aMiSLI0T9wvSpEjTIt0Jt2SEnG+jr1Q/",
	"Nk3TNMqyLIvyPM+j8Xg8HoolRay6NoL3S6f/niSv6LwNweDHRhgsoXi/bWR3oHAXm169+20+NX1ATq6d",
	"NeDHoGaco6ZhsKLE/YIsLxxeg2B14V+J1u9L8yBTd7p8WRrH+HrF1H/4+mZZbiE8VJz7VMiZb5IE+RYl",
	"a0hMK4yY1lFtI90xG3HPM4SwRGOFkm6aUTJKotur0+vby//f+aE1SqYFFJCPkpFrTjNa+EWI1wnchijr",
	"/0u03AhNXbbOaYPJpo7bIubuHH3rW+hGQ0tnqlx561WS1ovHtK4E9yHxg3U5N27unv42OIMC/oqf7D5e",
	"e338nNG3fSzJNOgPrFbSdvudJclbt7HfQB8m23CO1s6aKtgi5KAfJ+NDTKWiYKYaWfqtsE1dM7MawJrY",
	"3LpdeergFs3S23Qbgm4G6OrU/xxd3e13pGvIfH5ZuvaxPkZXG26kFj9upN92pSrsvKxf9NyfByzgw1R2",
	"917KhtVIaFzt/Sxbl3GOAoXXPYQgWY392z4/4Q7WezbY3h9wNwDZm+C7hqADIJiugp1+j0tjjgPSuEFq",
	"jLQBC6yQ8wqfQ/YC6Wx1+lTqxwP8M4rjAqnPnMaXmXMqWWBFi8+u0EssEqPGBlyVeMDhZZdkD8l/kvww",
	"4bp6cCfZkomKTSs8tGU3xuZDnzuYLJB/2BnHn/rDnXFcGjTLzeo0poICFkS6iONKcVYtlKXiJDlJYmjv",
	"2y8BAAD//72RBvX1DAAA",
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
