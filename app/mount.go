package app

import (
	"net/http"

	"github.com/jirawat050/GolangWeb/model"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/", index)
	mux.Handle("/upload/", http.StripPrefix("/upload", http.FileServer(http.Dir("upload"))))

	//news
	mux.Handle("/news/", http.StripPrefix("/news", http.HandlerFunc(newsView)))
	// mux.Handle("/news/", http.StripPrefix("/news", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	id := r.URL.Path[1:]
	// 	newsView(id).ServeHTTP(w, r)
	// })))
	mux.HandleFunc("/register", adminRegister)
	mux.HandleFunc("/login", adminLogin)
	adminMux := http.NewServeMux()

	adminMux.HandleFunc("/logout", adminLogout)
	adminMux.HandleFunc("/list", adminList)
	adminMux.HandleFunc("/create", adminCreate)
	adminMux.HandleFunc("/edit", adminEdit)
	mux.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(adminMux)))

}
func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		ok, err := model.CheckUserID(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
