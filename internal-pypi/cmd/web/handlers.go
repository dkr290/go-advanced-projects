package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Use the template.ParseFiles() function to read the template file into
	// template set. If there's an error, log detailed error message
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	//initialize slice containing two files. It's importnant
	//the base template should be first one

	files := []string{
		"./templates/base.html",
		"./templates/partials/nav.html",
		"./templates/pages/home.html",
	}
	//template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	packages, err := filepath.Glob(filepath.Join(packageDir, "*.whl"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Packages []string
	}{
		Packages: packages,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload request received")
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		log.Println("Invalid request method", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 10 MB
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("content")
	if err != nil {
		for _, fheaders := range r.MultipartForm.File {
			if len(fheaders) > 0 {
				file, err = fheaders[0].Open()
				header = fheaders[0]
				break
			}
		}
		if err != nil {
			log.Println("Error retreiving the file", err)
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}

	}
	defer file.Close()
	log.Println("Received file:", header.Filename)

	filePath := filepath.Join(packageDir, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating the file:", err)
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Println("Error saving the file", err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	log.Printf("File uploaded successfully: %s (%d bytes)\n", header.Filename, written)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	packageName := strings.TrimPrefix(r.URL.Path, "/simple/")
	packageName = strings.TrimSuffix(packageName, "/")
	packages, err := filepath.Glob(filepath.Join(packageDir, packageName+"*.whl"))
	if err != nil {
		log.Println("Filepath error ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>", packageName, packageName)
	for _, pkg := range packages {
		fileName := filepath.Base(pkg)
		fmt.Fprintf(w, "<a href=\"/packages/%s#sha256=placeholder\">%s</a><br>", fileName, fileName)
		log.Println(filepath.Base(pkg))

	}
	fmt.Fprintf(w, "</body></html>")
}
func packageHandler(w http.ResponseWriter, r *http.Request) {
	packageName := filepath.Base(r.URL.Path)
	packagePath := filepath.Join(packageDir, packageName)

	file, err := os.Open(packagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", packageName))
	io.Copy(w, file)
}
