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

func (app *Config) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Use the template.ParseFiles() function to read the template file into
	// template set. If there's an error, log detailed error message
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	// initialize slice containing two files. It's importnant
	// the base template should be first one

	files := []string{
		"./templates/base.html",
		"./templates/partials/nav.html",
		"./templates/pages/home.html",
	}
	// template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, err)
		return
	}

	packagesWhl, err := filepath.Glob(filepath.Join(packageDir, "*.whl"))
	if err != nil {
		app.serveError(w, err)
		return
	}
	packagesTar, err := filepath.Glob(filepath.Join(packageDir, "*.tar.gz"))
	if err != nil {
		app.serveError(w, err)
		return
	}

	allPackages := append(packagesWhl, packagesTar...)

	sortedPackages := sortPackages(allPackages)
	data := struct {
		Packages []string
	}{
		Packages: sortedPackages,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *Config) aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}
	// Use the template.ParseFiles() function to read the template file into
	// template set. If there's an error, log detailed error message
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	// initialize slice containing two files. It's importnant
	// the base template should be first one

	files := []string{
		"./templates/base.html",
		"./templates/partials/nav.html",
		"./templates/pages/about.html",
	}
	// template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *Config) contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		http.NotFound(w, r)
		return
	}
	// Use the template.ParseFiles() function to read the template file into
	// template set. If there's an error, log detailed error message
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	// initialize slice containing two files. It's importnant
	// the base template should be first one

	files := []string{
		"./templates/base.html",
		"./templates/partials/nav.html",
		"./templates/pages/contact.html",
	}
	// template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serveError(w, err)
		return
	}
}

func (app *Config) uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload request received")
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		log.Println("Invalid request method", r.Method)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // > 10 MB
	if err != nil {
		app.serveError(w, err)
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
			app.serveError(w, err)
			return
		}

	}
	defer file.Close()
	log.Println("Received file:", header.Filename)

	filePath := filepath.Join(packageDir, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating the file:", err)
		app.serveError(w, err)
		return
	}
	defer dst.Close()
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Println("Error saving the file", err)
		app.serveError(w, err)
		return
	}
	log.Printf("File uploaded successfully: %s (%d bytes)\n", header.Filename, written)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}

func (app *Config) simpleHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.RawPath
	if rawPath == "" {
		rawPath = r.URL.Path
	}
	packageName := strings.TrimPrefix(rawPath, "/simple/")
	packageName = strings.TrimSuffix(packageName, "/")
	// Create both underscore and hyphen versions of the package name
	underscorePackageName := strings.ReplaceAll(packageName, "-", "_")
	// hyphenPackageName := strings.ReplaceAll(packageName, "_", "-")

	// Search for packages using both versions
	packagesUnderscore, _ := filepath.Glob(filepath.Join(packageDir, underscorePackageName+"*"))
	// packagesHyphen, _ := filepath.Glob(filepath.Join(packageDir, hyphenPackageName+"*"))

	// Combine the results
	var packages []string
	packages = append(packages, packagesUnderscore...)

	if len(packages) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(
		w,
		"<!DOCTYPE html><html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>",
		packageName,
		packageName,
	)
	for _, pkg := range packages {
		fileName := filepath.Base(pkg)
		fmt.Fprintf(w, "<a href=\"/packages/%s\">%s</a><br>", fileName, fileName)

	}
	fmt.Fprintf(w, "</body></html>")
}

func (app *Config) packageHandler(w http.ResponseWriter, r *http.Request) {
	packageName := strings.TrimPrefix(r.URL.Path, "/packages/")
	packagePath := filepath.Join(packageDir, packageName)

	app.clientLog("Attemting to serve the package %s", packagePath)

	// Get file info to set Content-Length
	Fileinfo, err := os.Stat(packagePath)
	if err != nil {
		app.errorLog("Error accessing package file: %v", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file using http.FileServer
	http.ServeFile(w, r, packagePath)

	app.clientLog("Successfully served package %s with size %d", packagePath, Fileinfo.Size())
}

func (app *Config) favIconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func (app *Config) testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Invalid request method", r.Method)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	s := "Liveness and Readiness"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(
		w,
		"<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Test for %s</h1>",
		s,
	)
	fmt.Fprintf(w, "</body></html>")
}
