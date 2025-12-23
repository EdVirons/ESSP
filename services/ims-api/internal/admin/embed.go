package admin

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed dashboard/dist/*
var dashboardFS embed.FS

// ServeDashboard serves the embedded dashboard static files
// It handles SPA routing by serving index.html for client-side routes
func ServeDashboard(r chi.Router) {
	// Get the dist subdirectory
	distFS, err := fs.Sub(dashboardFS, "dashboard/dist")
	if err != nil {
		// Fallback: serve a simple message if dashboard is not built
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>ESSP Admin Dashboard</title></head>
<body>
<h1>ESSP Admin Dashboard</h1>
<p>Dashboard static files not found. Please build the dashboard first:</p>
<pre>cd dashboard && npm install && npm run build</pre>
<p>Then rebuild the ims-api service.</p>
</body>
</html>`))
		})
		return
	}

	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		path := chi.URLParam(req, "*")
		if path == "" {
			path = "index.html"
		}

		// Try to serve the exact file
		if !strings.HasSuffix(path, "/") {
			// Check if file exists
			content, err := fs.ReadFile(distFS, path)
			if err == nil {
				// Determine content type
				contentType := "application/octet-stream"
				if strings.HasSuffix(path, ".html") {
					contentType = "text/html; charset=utf-8"
				} else if strings.HasSuffix(path, ".css") {
					contentType = "text/css; charset=utf-8"
				} else if strings.HasSuffix(path, ".js") {
					contentType = "application/javascript; charset=utf-8"
				} else if strings.HasSuffix(path, ".svg") {
					contentType = "image/svg+xml"
				} else if strings.HasSuffix(path, ".png") {
					contentType = "image/png"
				} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
					contentType = "image/jpeg"
				} else if strings.HasSuffix(path, ".json") {
					contentType = "application/json; charset=utf-8"
				} else if strings.HasSuffix(path, ".woff2") {
					contentType = "font/woff2"
				} else if strings.HasSuffix(path, ".woff") {
					contentType = "font/woff"
				}
				w.Header().Set("Content-Type", contentType)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(content)
				return
			}
		}

		// For SPA routes or if file doesn't exist, serve index.html
		indexContent, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			http.Error(w, "Dashboard not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexContent)
	})
}

// DashboardHandler returns an http.Handler for serving the dashboard
// This is an alternative interface if needed
func DashboardHandler() http.Handler {
	distFS, err := fs.Sub(dashboardFS, "dashboard/dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Dashboard not found", http.StatusNotFound)
		})
	}

	return &spaHandler{
		staticHandler: http.FileServer(http.FS(distFS)),
		staticFS:      distFS,
	}
}

// spaHandler handles SPA routing for the dashboard
type spaHandler struct {
	staticHandler http.Handler
	staticFS      fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" {
		path = "/"
	}

	// Check if requesting a static file
	if path != "/" && !strings.HasSuffix(path, "/") {
		if _, err := fs.Stat(h.staticFS, strings.TrimPrefix(path, "/")); err == nil {
			// Serve the static file
			r.URL.Path = path
			h.staticHandler.ServeHTTP(w, r)
			return
		}
	}

	// Serve index.html for SPA routes
	indexContent, err := fs.ReadFile(h.staticFS, "index.html")
	if err != nil {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexContent)
}
