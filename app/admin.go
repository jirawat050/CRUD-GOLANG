package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jirawat050/GolangWeb/model"
	"github.com/jirawat050/GolangWeb/view"
)

func adminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		userID, err := model.Login(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "user",
			Value:    userID,
			MaxAge:   int(10 * time.Minute / time.Second),
			HttpOnly: true,
			Path:     "/",
		})
		fmt.Println(userID)

		http.Redirect(w, r, "/admin/list", http.StatusSeeOther)
		return
	}
	view.AdminLogin(w, nil)
}
func adminRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := model.Register(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	view.AdminRegister(w, nil)
}
func adminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
func adminList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		id := r.FormValue("id")
		if action == "delete" {
			err := model.DeleteNews(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/admin/list", http.StatusSeeOther)
		return
	}
	list, err := model.ListNews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.AdminList(w, &view.AdminListData{
		List: list,
	})
}
func adminCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		n := model.News{
			Title:  r.FormValue("title"),
			Detail: r.FormValue("detail"),
		}

		if file, handle, err := r.FormFile("image"); err == nil {
			defer file.Close()
			fileName := time.Now().Format(time.RFC3339) + "-" + handle.Filename
			fp, err := os.Create("upload/" + fileName)
			if err == nil {
				io.Copy(fp, file)
			}
			fp.Close()
			n.Image = "/upload/" + fileName

		}
		err := model.CreateNews(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/admin/list", http.StatusSeeOther)
		return
	}
	view.AdminCreate(w, nil)
}
func adminEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	n, err := model.GetNews(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodPost {
		n.Title = r.FormValue("title")
		n.Detail = r.FormValue("detail")
		if file, handle, err := r.FormFile("image"); err == nil {
			defer file.Close()
			fileName := time.Now().Format(time.RFC3339) + "-" + handle.Filename
			fp, err := os.Create("upload/" + fileName)
			if err == nil {
				io.Copy(fp, file)
			}
			fp.Close()
			n.Image = "/upload/" + fileName

		}
		err := model.UpdateNews(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/list", http.StatusFound)
		return
	}

	view.AdminEdit(w, n)
}
