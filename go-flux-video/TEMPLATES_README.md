# Template Structure

## Overview

The HTML templates have been moved out of the Go code and into separate files in the `templates/` directory. This makes it easier to:
- Edit HTML/CSS/JavaScript without recompiling Go code
- Maintain cleaner separation of concerns
- Use proper template inheritance
- Collaborate with front-end developers

## Template Files

### 1. `templates/base.html`
- Main layout template with HTML structure and CSS styles
- Defines the overall page layout
- Includes other templates using `{{template "name" .}}`

### 2. `templates/upload-section.html`
- Upload interface with drag-and-drop functionality
- Progress bars and status messages
- Two upload buttons:
  - Upload to `./images/` directory (for img2img)
  - Upload to gallery (for viewing)

### 3. `templates/controls.html`
- Search box for filtering images
- Refresh button
- Download all button

### 4. `templates/gallery.html`
- Gallery container for displaying images
- Loading state and empty state templates

### 5. `templates/modal.html`
- Full-screen image viewer modal
- Close button and overlay

### 6. `templates/scripts.html`
- All JavaScript functionality
- Image loading and display
- Upload handling with drag-and-drop
- API communication
- Search and filtering

## Template Inheritance

The templates use Go's template inheritance pattern:

```
base.html (main layout)
├── upload-section.html (upload interface)
├── controls.html (search/refresh controls)
├── gallery.html (image gallery)
├── modal.html (full-screen viewer)
└── scripts.html (JavaScript)
```

## How It Works

### 1. Template Loading
The Go code loads all templates with:
```go
tmpl, err := template.ParseFiles(
    "templates/base.html",
    "templates/upload-section.html",
    // ... other templates
)
```

### 2. Template Execution
```go
data := map[string]interface{}{
    "Title":     "FLUX Image Gallery",
    "OutputDir": s.OutputDir,
}
tmpl.Execute(w, data)
```

### 3. Template Composition
`base.html` includes other templates:
```html
{{template "upload-section" .}}
{{template "controls" .}}
{{template "gallery" .}}
{{template "modal" .}}
{{template "scripts" .}}
```

## Benefits

### 1. Development
- **Hot Reload**: Edit HTML/CSS/JS without restarting Go server
- **Separation**: Front-end and back-end developers can work independently
- **Version Control**: Easier to track changes to UI components

### 2. Maintenance
- **Readability**: Smaller, focused files instead of one huge Go file
- **Reusability**: Templates can be reused across different pages
- **Testing**: Easier to test individual components

### 3. Performance
- **Compilation**: Templates are parsed once at server start
- **Caching**: Go's template engine caches parsed templates
- **Efficiency**: Only needed templates are loaded into memory

## Editing Templates

### To Modify Styles:
1. Edit `templates/base.html` CSS section
2. No need to recompile Go code
3. Refresh browser to see changes

### To Add New Features:
1. Create new template file (e.g., `templates/new-feature.html`)
2. Add to template loading in `webserver.go`
3. Include in `base.html` if needed
4. Update JavaScript in `templates/scripts.html`

### To Change Layout:
1. Modify `templates/base.html` structure
2. Reorder template inclusions as needed
3. Update CSS for new layout

## Testing

Run the template test:
```bash
go run test_templates.go
```

This will:
1. Check all template files exist
2. Parse all templates
3. Test template execution with dummy data

## Troubleshooting

### Common Issues:

1. **"Template not found"**
   - Check template file paths in `webserver.go`
   - Verify files exist in `templates/` directory

2. **"Failed to parse template"**
   - Check for syntax errors in HTML/Go template syntax
   - Look for unclosed tags or mismatched braces

3. **Template not rendering**
   - Verify template is included in `base.html`
   - Check data passed to template (use `{{.}}` to debug)

4. **JavaScript errors**
   - Check browser console for errors
   - Verify API endpoints match Go handlers

### Debugging:
- Add debug output: `{{printf "%#v" .}}` to see all data
- Check Go server logs for template errors
- Use browser developer tools to inspect network requests

## Best Practices

1. **Keep Templates Small**: Each template should have a single responsibility
2. **Use Meaningful Names**: Template names should describe their purpose
3. **Document Changes**: Update this README when adding new templates
4. **Test Thoroughly**: Test templates after any changes
5. **Backup Originals**: Keep backup of working templates

## Future Enhancements

Potential improvements:
1. **Template Caching**: Cache parsed templates for better performance
2. **Template Minification**: Minify HTML/CSS/JS in production
3. **Internationalization**: Support multiple languages
4. **Theme Support**: Allow different CSS themes
5. **Component Library**: Create reusable UI components