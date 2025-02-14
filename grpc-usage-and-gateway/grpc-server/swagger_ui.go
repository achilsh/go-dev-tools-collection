package main

import (
	"fmt"
	swagger "grpc-usage-and-gateway/ui/data/swagger"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

// RegisterSwaggerUI
func RegisterSwaggerUI(mux *http.ServeMux) {
	// register swagger
	mux.HandleFunc("/swagger/", swaggerFile)
	swaggerUI(mux)
}

/*
*
serveSwaggerUI: 提供UI支持
*/
func swaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		// https://github.com/swagger-api/swagger-ui 工具dist目录下的文件在本地存放路径
		Prefix: "thirdparty/swagger_ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

/*
*
swaggerFile: 提供对swagger.json文件的访问支持
*/
func swaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		fmt.Printf("Not Found: %s\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	name := path.Join("gen/go/proto", p)
	fmt.Printf("Serving swagger-file: %s\n", name)
	http.ServeFile(w, r, name)
}
