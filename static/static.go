package static

import (
	"net/http"
	"os"
)

func ServeStatic(root string) {
	fsAssets := http.FileServer(http.Dir(root + "/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fsAssets))

	publicAssets := http.FileServer(http.Dir(root + "/public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicAssets))

	// serve public files at root level so /style.css, /favicon.ico etc. resolve
	// correctly (more specific patterns take priority over pages catch-all /)
	entries, err := os.ReadDir(root + "/public")
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		filePath := root + "/public/" + name
		http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filePath)
		})
	}
}
