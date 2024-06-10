package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"html/template"
	"path/filepath"
	"math/rand"
	"time"
	"strings"
	
)

var port = os.Getenv("PORT")
// Function to generate a random alphabetic string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Function to generate a unique filename with a specific extension
func generateFilename(extension string, length int) string {
	return randomString(length) + "." + strings.TrimLeft(extension, ".")
}

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
			<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
            <title>File Upload</title>
        </head>
        <body>           
            
			<div class="container my-5">
			<div class="p-5 text-center bg-body-tertiary rounded-3">
				<svg class="bi mt-4 mb-3" style="color: var(--bs-indigo);" width="100" height="100"><use xlink:href="#bootstrap"></use></svg>
				<h1 class="text-body-emphasis">Upload a File</h1>
				<p class="col-lg-8 mx-auto fs-5 text-muted">
				This is an example of file upload using Openshift
				</p>
				<div class="d-inline-flex gap-2 mb-5">
					<form method="post" class="form-control needs-validation" action="/upload" enctype="multipart/form-data">
						<div class="row">
							<div class="col-sm-12">
							<label class="visually-hidden" for="file">File</label>
								<input class="form-control" type="file" id="file" name="file">
								<div class="invalid-feedback">
									Valid file required.
								</div>
							</div>
							<div class="col-auto">
								<br/>
								<button class="btn btn-primary" type="submit" value="Upload">Submit</button>
							</div>
						</div>
					</form>
				</div>
			</div>
			</div>
        </body>
        </html>
        `
        t, _ := template.New("upload").Parse(tmpl)
        t.Execute(w, nil)		
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		extension := filepath.Ext(handler.Filename)
		// Create a new file in the specified directory
		uploadPath := "/usr/share/external_persistent/"
		// Generate filename add sequence if already exists

		filename := generateFilename(extension, 12)
		if err != nil {
			// Handle the error appropriately (e.g., log it, return an error)
		}
		dst, err := os.Create(uploadPath + filename )
		if err != nil {
			http.Error(w, "Error creating file on server: " + err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error saving file " + filename+" on server: " + err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl := `
        <!DOCTYPE html>
        <html>
        <head>
			<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
            <title>File Upload</title>
        </head>
        <body>           
            
			<div class="container my-5">
			<div class="p-5 text-center bg-body-tertiary rounded-3">
				<h1 class="text-body-emphasis">Upload Successful</h1>
				<p class="col-lg-8 mx-auto fs-5 text-muted">
				File uploaded successfully. ( <a href="/">Upload Another File</a> )
				</p>
			</div>
        </body>
        </html>
        `
        t, _ := template.New("upload").Parse(tmpl)
        t.Execute(w, nil)		
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
