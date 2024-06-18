package templatecache

import (
	"html/template"
	"log"
	"path/filepath"
	"sync"
)

/* This module provides a Get function that retrieves a template
from the cache or parses and caches it if it doesn't exist.
The Get function takes a name parameter,
which is the key used to store and retrieve the template in the cache,
and a variadic paths parameter, which is a slice of file paths to parse for the template.
The Get function uses a read-write mutex to ensure thread-safety when accessing and modifying
the templateCache map. It first attempts to retrieve the template from the cache using a read lock.
If the template is not found in the cache, it acquires a write lock and checks again in case another
goroutine cached the template while it was waiting for the lock. If the template is still not found, it parses the template files and stores
the resulting *template.Template in the cache.
*/

var (
	templateCache = make(map[string]*template.Template)
	mutex         sync.RWMutex
)

func Get(name string, verbose bool, paths ...string) (*template.Template, error) {

	mutex.RLock()
	tmpl, ok := templateCache[name]
	mutex.RUnlock()

	if ok {
		if verbose {
			log.Printf("Template %s already in the cache so not adding it there\n", name)
		}
		return tmpl, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	tmpl, ok = templateCache[name]
	if ok {
		if verbose {
			log.Printf("Template %s already in the cache so not adding it there\n", name)
		}
		return tmpl, nil
	}
	// This line is responsible for creating a new *template.Template instance and parsing the template files specified in the paths slice.
	//template.New(name): This creates a new, empty *template.Template instance with the given name.
	//The name is used for error reporting and is typically the name of the template or the base name of the first file being parsed.
	//ParseFiles(paths...): This method parses the template definition from the files specified in the paths slice. The ... syntax is
	//used to pass the slice elements as separate arguments to the ParseFiles function. If any of the files cannot be read
	//or parsed, an error is returned.

	var err error
	tmpl, err = template.New(name).Funcs(template.FuncMap{
		"pathJoin": filepath.Join,
	}).ParseFiles(paths...)
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Printf("The template %s is not in cache\n", tmpl.Name())
		log.Println("Building the template in the cache and fine")
	}

	templateCache[name] = tmpl
	return tmpl, nil
}
