package pages

import (
	"net/http"
	"strings"
)

func ServePages(root string) {
	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		result, err := FormatPageList(root)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(result))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		paths, err := ScanAllFilepaths(root)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		var filepath string

		for _, v := range paths {
			if FormatFilepathToRoute(root, v) == r.URL.Path {
				filepath = v
				break
			}
		}

		// Determine status code for special error pages
		statusCode := 200
		if filepath != "" {
			// extract filename from full path
			parts := strings.Split(filepath, "/")
			filename := parts[len(parts)-1]
			if filename == "403.md" {
				statusCode = 403
			} else if filename == "404.md" {
				statusCode = 404
			} else if filename == "500.md" {
				statusCode = 500
			}
		}

		if filepath == "" {
			// try to serve 404.md
			notFoundPath := root + "/404.md"
			page, err := ReadPageFile(root, notFoundPath)
			if err == nil {
				tmp, _ := ReadTemplateFile(root)
				result := FormatTemplate(tmp, page)
				w.WriteHeader(404)
				w.Header().Add("Content-Type", "text/html")
				w.Write([]byte(result))
				return
			}
			w.WriteHeader(404)
			w.Write([]byte("404 Not Found"))
			return
		}

		page, err := ReadPageFile(root, filepath)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		tmp, err := ReadTemplateFile(root)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		result := FormatTemplate(tmp, page)

		w.WriteHeader(statusCode)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(result))
	})
}
