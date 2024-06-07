package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"html/template"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path != "" {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	} else {
		tmpl := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>File Upload</title>
        </head>
        <body>
            <h1>Upload a File</h1>
            <form method="post" enctype="multipart/form-data">
                <input type="file" name="file">
                <input type="submit" value="Upload">
            </form>
        </body>
        </html>
        `
        t, _ := template.New("upload").Parse(tmpl)
        t.Execute(w, nil)		
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create a new file in the specified directory
		uploadPath := "/usr/share/external_persistent"
		filename := "uploaded_file.txt"
		dst, err := os.Create(uploadPath + filename)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "File uploaded successfully!")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
