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

	packagesWhl, err := filepath.Glob(filepath.Join(packageDir, "*.whl"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	packagesTar, err := filepath.Glob(filepath.Join(packageDir, "*.tar.gz"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
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
		"./templates/pages/about.html",
	}
	//template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
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
		"./templates/pages/contact.html",
	}
	//template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
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
	rawPath := r.URL.RawPath
	if rawPath == "" {
		rawPath = r.URL.Path
	}
	packageName := strings.TrimPrefix(rawPath, "/simple/")
	packageName = strings.TrimSuffix(packageName, "/")
	// Create both underscore and hyphen versions of the package name
	underscorePackageName := strings.ReplaceAll(packageName, "-", "_")
	hyphenPackageName := strings.ReplaceAll(packageName, "_", "-")

	// Search for packages using both versions
	packagesUnderscore, _ := filepath.Glob(filepath.Join(packageDir, underscorePackageName+"*"))
	packagesHyphen, _ := filepath.Glob(filepath.Join(packageDir, hyphenPackageName+"*"))

	// Combine the results
	packages := append(packagesUnderscore, packagesHyphen...)

	if len(packages) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>", packageName, packageName)
	for _, pkg := range packages {
		fileName := filepath.Base(pkg)
		fmt.Fprintf(w, "<a href=\"/packages/%s\">%s</a><br>", fileName, fileName)

	}
	fmt.Fprintf(w, "</body></html>")
}

func packageHandler(w http.ResponseWriter, r *http.Request) {
	packageName := strings.TrimPrefix(r.URL.Path, "/packages/")
	packagePath := filepath.Join(packageDir, packageName)

	log.Printf("Attempting to serve package: %s", packagePath)
	file, err := os.Open(packagePath)
	if err != nil {
		log.Printf("Error opening package file: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", packageName))
	written, err := io.Copy(w, file)
	if err != nil {
		log.Printf("Error serving package file: %v", err)
	} else {
		log.Printf("Successfully served package %s (%d bytes)", packageName, written)
	}
}
