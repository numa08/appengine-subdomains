package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

// DOMAIN は実行するアプリの hostname.
// 環境変数とかでセットできるようにする
const DOMAIN = "myapp.team"

// ReservedSubdomains は動的に生成されたくない subdomain のリスト.
// ルーティングの段階ではチェックはしないで、アプリケーションの方で良い感じに
// するのもあり
var ReservedSubdomains = []string{
	"admin", "docs",
}

func main() {
	router := mux.NewRouter()
	router.Host(DOMAIN).PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath.Join("static", "top", "index.html"))
	})
	router.Host(fmt.Sprintf("{subdomain}.%s", DOMAIN)).PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		subdomain := vars["subdomain"]
		for _, d := range ReservedSubdomains {
			if d == subdomain {
				http.NotFound(w, req)
				return
			}
		}
		http.ServeFile(w, req, filepath.Join("static", "team", "index.html"))
	})
	http.Handle("/", router)
	appengine.Main()
}
