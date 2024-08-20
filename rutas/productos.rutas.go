package rutas

import "net/http"

func GetProductosHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get productos"))
}

func GetProductoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get producto"))
}

func PostProductoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post producto"))
}

func DeleteProductoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete producto"))
}
