package web

import (
	"net/http"
	"sort"
)

func ApplyMethods(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		methods := sort.StringSlice(methods)
		sort.Sort(methods)
		i := sort.SearchStrings(methods, r.Method)
		if i < len(methods) && methods[i] == r.Method {
			next.ServeHTTP(rw, r)
			return
		}
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
