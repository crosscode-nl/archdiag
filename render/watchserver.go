package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/crosscode-nl/archdiag/parse"
	"github.com/crosscode-nl/archdiag/theme"
)

// WatchServer serves rendered diagrams with live reload via SSE.
type WatchServer struct {
	renderer *HTMLRenderer
	theme    string
	port     int

	mu      sync.RWMutex
	html    map[string][]byte // rendered HTML keyed by diagram name
	clients map[chan struct{}]struct{}
}

// NewWatchServer creates a watch server.
func NewWatchServer(renderer *HTMLRenderer, theme string, port int) *WatchServer {
	return &WatchServer{
		renderer: renderer,
		theme:    theme,
		port:     port,
		html:     make(map[string][]byte),
		clients:  make(map[chan struct{}]struct{}),
	}
}

// WatchFile watches a single YAML file and serves it with live reload.
func (ws *WatchServer) WatchFile(yamlPath string, openBrowser bool) error {
	key := filepath.Base(yamlPath)
	if err := ws.renderFile(key, yamlPath); err != nil {
		return fmt.Errorf("initial render: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(yamlPath); err != nil {
		return err
	}

	// Start HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.servePage(key, w)
	})
	mux.HandleFunc("/events", ws.handleSSE)

	addr := fmt.Sprintf("localhost:%d", ws.port)
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Printf("serving at http://%s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	if openBrowser {
		openURL(fmt.Sprintf("http://%s", addr))
	}

	// Watch loop
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) {
				start := time.Now()
				if err := ws.renderFile(key, yamlPath); err != nil {
					log.Printf("render error: %v", err)
					continue
				}
				log.Printf("re-rendered %s (%s)", filepath.Base(yamlPath), time.Since(start).Round(time.Millisecond))
				ws.notifyClients()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

// WatchDir watches a directory of YAML files and serves an index with links.
func (ws *WatchServer) WatchDir(dirPath string, openBrowser bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(dirPath); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "" {
			ws.handleIndex(dirPath, w, r)
			return
		}
		// Serve individual diagram
		name := strings.TrimPrefix(r.URL.Path, "/")
		name = strings.TrimSuffix(name, ".html")
		key := name + ".yaml"
		yamlFile := filepath.Join(dirPath, key)
		if err := ws.renderFile(key, yamlFile); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		ws.servePage(key, w)
	})
	mux.HandleFunc("/events", ws.handleSSE)

	addr := fmt.Sprintf("localhost:%d", ws.port)
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Printf("serving directory at http://%s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	if openBrowser {
		openURL(fmt.Sprintf("http://%s", addr))
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) && strings.HasSuffix(event.Name, ".yaml") {
				log.Printf("changed: %s", filepath.Base(event.Name))
				ws.notifyClients()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

func (ws *WatchServer) renderFile(key, yamlPath string) error {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	d, err := parse.Parse(data)
	if err != nil {
		return err
	}

	if ws.theme != "" {
		d.Theme = ws.theme
	}

	var buf bytes.Buffer
	if err := ws.renderer.RenderWithWatch(d, &buf); err != nil {
		return err
	}

	ws.mu.Lock()
	ws.html[key] = buf.Bytes()
	ws.mu.Unlock()

	return nil
}

func (ws *WatchServer) servePage(key string, w http.ResponseWriter) {
	ws.mu.RLock()
	html := ws.html[key]
	ws.mu.RUnlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func (ws *WatchServer) handleIndex(dirPath string, w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob(filepath.Join(dirPath, "*.yaml"))

	type entry struct {
		Path     string
		Title    string
		Subtitle string
	}

	var entries []entry
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		d, err := parse.Parse(data)
		if err != nil {
			continue
		}
		name := strings.TrimSuffix(filepath.Base(f), ".yaml")
		entries = append(entries, entry{
			Path:     "/" + name + ".html",
			Title:    d.Title,
			Subtitle: d.Subtitle,
		})
	}

	// Render index template
	themeName := "dark"
	if ws.theme != "" {
		themeName = ws.theme
	}
	css, err := theme.LoadAll()
	if err != nil {
		http.Error(w, "failed to load theme CSS", http.StatusInternalServerError)
		return
	}

	type indexData struct {
		ThemeCSS     template.CSS
		DefaultTheme string
		Diagrams     []entry
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ws.renderer.Tmpl.ExecuteTemplate(w, "index.html.tmpl", indexData{
		ThemeCSS:     template.CSS(css),
		DefaultTheme: themeName,
		Diagrams:     entries,
	})
}

func (ws *WatchServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 1)

	ws.mu.Lock()
	ws.clients[ch] = struct{}{}
	ws.mu.Unlock()

	defer func() {
		ws.mu.Lock()
		delete(ws.clients, ch)
		ws.mu.Unlock()
	}()

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (ws *WatchServer) notifyClients() {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	for ch := range ws.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func openURL(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		return
	}
	exec.Command(cmd, url).Start()
}
