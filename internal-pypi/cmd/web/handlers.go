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
	// var newPackagesWhl []string
	// for _, pWl := range packagesWhl {
	// 	baseName := filepath.Base(pWl)
	// 	newPackagesWhl = append(newPackagesWhl, baseName)
	//
	// }
	packagesTar, err := filepath.Glob(filepath.Join(packageDir, "*.tar.gz"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// var newPackagesTar []string
	// for _, pTr := range packagesTar {
	// 	baseName := filepath.Base(pTr)
	// 	newPackagesTar = append(newPackagesTar, baseName)
	// }

	allPackages := append(packagesWhl, packagesTar...)
	log.Println(allPackages)
	// var onlyNames []string
	// var nameWithoutExt string
	// for _, file := range allPackages {
	// 	baseName := filepath.Base(file)
	// 	ext := filepath.Ext(baseName)
	// 	if ext == ".gz" {
	// 		nameWithoutExt = extractName(baseName)
	// 	}
	// 	if ext == ".whl" {
	// 		nameWithoutExt = extractName(baseName)
	// 	}
	//
	// 	onlyNames = append(onlyNames, nameWithoutExt)
	// }
	// var newData []string
	// for _, f := range onlyNames {
	// 	f += "/"
	// 	newData = append(newData, f)
	// }
	// newData = removeDuplicates(newData)

	// type PagesData struct {
	// 	Name   []string
	// 	WlPkg  []string
	// 	TarPkg []string
	// }
	//
	// data := struct {
	// 	Pages []PagesData
	// }{
	// 	Pages: []PagesData{
	// 		{Name: newData, WlPkg: newPackagesWhl, TarPkg: newPackagesTar},
	// 	},
	// }
	data := struct {
		Packages []string
	}{
		Packages: allPackages,
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
	rawPath := r.URL.RawPath
	if rawPath == "" {
		rawPath = r.URL.Path
	}
	log.Println(rawPath)
	packageName := strings.TrimPrefix(rawPath, "/simple/")
	packageName = strings.TrimSuffix(packageName, "/")
	packages, err := filepath.Glob(filepath.Join(packageDir, packageName+"*"))
	if err != nil {
		log.Println("Filepath error ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Listing packages for: %s", packageName)
	if len(packages) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>", packageName, packageName)
	for _, pkg := range packages {
		fileName := filepath.Base(pkg)
		fmt.Fprintf(w, "<a href=\"/packages/%s\">%s</a><br>", fileName, fileName)
		log.Println(filepath.Base(pkg))

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
