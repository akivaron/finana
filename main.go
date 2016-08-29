package main

import (
	"Chatapp/tracer"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

const (
	GITHUB_ID     = "78d0161c3f2c89033a6d"
	GITHUB_SECRET = "bc63b27b665587cd8437aeb2569780b44972bd31"
	GOOGLE_ID     = "911008785927-0av98e5vvdedmm95fuf5qve50584ninq.apps.googleusercontent.com"
	GOOGLE_SECRET = "YWZ8e-otEYyvRU9SLiMqt_iR"
)

var addr *string

var avatars Avatar = TryAvatars{useFileSystemAvatar, UseGravatarAvatar, UseAuthAvatar}

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	addr = flag.String("addr", "10.1.50.106:8080", "The addr of the application.")
	flag.Parse() // parse the flags

	//set up gomniauth
	gomniauth.SetSecurityKey("1234567890")
	gomniauth.WithProviders(
		facebook.New("local", "1234567890", "http://10.1.50.106:8080/auth/callback/facebook"),
		github.New("e246daf9b9bd5b1a748c", "ab524e7bac482c13d6d85e304cbd70f94634dccb", "http://10.1.50.106:8080/auth/callback/github"),
		google.New(GOOGLE_ID, GOOGLE_SECRET, "http://10.1.50.106:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = tracer.New(os.Stdout)

	fsAssets := http.FileServer(http.Dir("assets"))
	fsAvatars := http.FileServer((http.Dir("avatars")))

	http.Handle("/assets/", http.StripPrefix("/assets", fsAssets))
	http.Handle("/avatars/", http.StripPrefix("/avatars", fsAvatars))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/room", r)
	server := http.Server{
		Addr:    *addr,
		Handler: nil,
	}
	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
