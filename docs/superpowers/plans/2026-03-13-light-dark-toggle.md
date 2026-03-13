# Light/Dark Theme Toggle Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a runtime light/dark theme toggle button to rendered HTML diagrams, persisting user preference via localStorage.

**Architecture:** Extract shared CSS into `base.css`, scope theme variables under `[data-theme]` selectors, embed both themes in every HTML output, and add a fixed toggle button with JS that swaps the `data-theme` attribute on `<html>`.

**Tech Stack:** Go (html/template), CSS custom properties, vanilla JS, jj for version control

---

## File Structure

| Action | File | Responsibility |
|--------|------|---------------|
| Create | `theme/css/base.css` | Shared reset, body, layout, media queries |
| Modify | `theme/css/dark.css` | Dark theme CSS variables only (scoped to `[data-theme="dark"]`) |
| Modify | `theme/css/light.css` | Light theme CSS variables only (scoped to `[data-theme="light"]`) |
| Modify | `theme/theme.go` | Add `LoadAll()` function |
| Modify | `theme/theme_test.go` | Add tests for `LoadAll()`, update existing tests |
| Modify | `render/html.go` | Add `DefaultTheme` to `renderData`, use `LoadAll()` |
| Modify | `render/templates/diagram.html.tmpl` | Add `data-theme`, toggle button, toggle script |
| Modify | `render/templates/index.html.tmpl` | Same toggle pattern |
| Modify | `render/watchserver.go` | Update `handleIndex` to use `LoadAll()` + `DefaultTheme` |
| Modify | `render/html_test.go` | Update assertions, add toggle test |

---

## Chunk 1: CSS Restructuring and Theme Loader

### Task 1: Create `base.css` and restructure theme CSS files

**Files:**
- Create: `theme/css/base.css`
- Modify: `theme/css/dark.css`
- Modify: `theme/css/light.css`

- [ ] **Step 1: Create `theme/css/base.css`**

```css
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
    font-family: var(--font-family);
    background: var(--bg-primary);
    color: var(--text-primary);
    padding: 32px;
    line-height: 1.6;
}

.diagram { max-width: 1200px; margin: 0 auto; }
.diagram > h1 { font-size: 28px; margin-bottom: 8px; color: var(--text-primary); }
.diagram > h2 { font-size: 18px; margin-bottom: 24px; color: var(--text-muted); font-weight: 400; }

.section--half { width: 48%; display: inline-block; vertical-align: top; }

@media (max-width: 768px) {
    .grid { grid-template-columns: 1fr !important; }
    .cards { flex-direction: column; }
    .section--half { width: 100%; }
}
```

- [ ] **Step 2: Reduce `theme/css/dark.css` to scoped variables only**

Replace entire file with:

```css
[data-theme="dark"] {
    --bg-primary: #0d1117;
    --bg-card: #161b22;
    --bg-nested: #0d1117;
    --text-primary: #e6edf3;
    --text-muted: #8b949e;
    --text-detail: #8b949e;
    --border-default: #30363d;
    --border-nested: #21262d;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --font-mono: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
    --radius-section: 12px;
    --radius-card: 8px;
    --radius-badge: 4px;
    --radius-flow-node: 6px;
    --arrow-color: #484f58;
}
```

- [ ] **Step 3: Reduce `theme/css/light.css` to scoped variables only**

Replace entire file with:

```css
[data-theme="light"] {
    --bg-primary: #ffffff;
    --bg-card: #f6f8fa;
    --bg-nested: #ffffff;
    --text-primary: #1f2328;
    --text-muted: #656d76;
    --text-detail: #656d76;
    --border-default: #d0d7de;
    --border-nested: #e1e4e8;
    --font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --font-mono: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
    --radius-section: 12px;
    --radius-card: 8px;
    --radius-badge: 4px;
    --radius-flow-node: 6px;
    --arrow-color: #b0b8c1;
}
```

- [ ] **Step 4: Commit CSS restructuring**

```bash
jj commit -m "refactor: extract base.css and scope theme variables under [data-theme]"
```

### Task 2: Add `LoadAll()` to theme loader (TDD)

**Files:**
- Modify: `theme/theme.go`
- Modify: `theme/theme_test.go`

- [ ] **Step 1: Update existing theme tests for new CSS structure**

The existing tests assert `:root`-style content that has changed (Task 1 already restructured the CSS files). Update `theme/theme_test.go`:

- `TestLoadDark`: Change assertion from `--bg-primary: #0d1117` to `[data-theme="dark"]` (variables are now under this selector)
- `TestLoadLight`: Change assertion from `--bg-primary: #ffffff` to `[data-theme="light"]`

```go
func TestLoadDark(t *testing.T) {
	css, err := theme.Load("dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if css == "" {
		t.Fatal("dark theme CSS is empty")
	}
	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Error("dark theme missing [data-theme] selector")
	}
	if !strings.Contains(css, "--bg-primary: #0d1117") {
		t.Error("dark theme missing --bg-primary")
	}
}

func TestLoadLight(t *testing.T) {
	css, err := theme.Load("light")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(css, `[data-theme="light"]`) {
		t.Error("light theme missing [data-theme] selector")
	}
	if !strings.Contains(css, "--bg-primary: #ffffff") {
		t.Error("light theme missing --bg-primary")
	}
}
```

- [ ] **Step 2: Write failing test for `LoadAll()`**

Add to `theme/theme_test.go`:

```go
func TestLoadAll(t *testing.T) {
	css, err := theme.LoadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if css == "" {
		t.Fatal("combined CSS is empty")
	}
	// Contains base styles
	if !strings.Contains(css, "box-sizing: border-box") {
		t.Error("missing base styles")
	}
	// Contains dark theme variables
	if !strings.Contains(css, `[data-theme="dark"]`) {
		t.Error("missing dark theme selector")
	}
	// Contains light theme variables
	if !strings.Contains(css, `[data-theme="light"]`) {
		t.Error("missing light theme selector")
	}
}
```

- [ ] **Step 3: Run tests to verify `TestLoadAll` fails**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go test ./theme/ -v`
Expected: `TestLoadDark` PASS, `TestLoadLight` PASS, `TestLoadAll` FAIL — `theme.LoadAll` undefined

- [ ] **Step 4: Implement `LoadAll()` in `theme/theme.go`**

Add after the existing `Load()` function:

```go
// LoadAll returns the combined CSS for base styles plus both themes.
func LoadAll() (string, error) {
	base, err := themesFS.ReadFile("css/base.css")
	if err != nil {
		return "", fmt.Errorf("load base CSS: %w", err)
	}
	dark, err := themesFS.ReadFile("css/dark.css")
	if err != nil {
		return "", fmt.Errorf("load dark CSS: %w", err)
	}
	light, err := themesFS.ReadFile("css/light.css")
	if err != nil {
		return "", fmt.Errorf("load light CSS: %w", err)
	}
	return string(base) + "\n" + string(dark) + "\n" + string(light), nil
}
```

- [ ] **Step 5: Run all theme tests**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go test ./theme/ -v`
Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
jj commit -m "feat: add LoadAll() to theme package for dual-theme support"
```

---

## Chunk 2: Renderer and Templates

### Task 3: Update renderer to use `LoadAll()` and pass `DefaultTheme` (TDD)

**Files:**
- Modify: `render/html.go`
- Modify: `render/html_test.go`

- [ ] **Step 1: Update existing tests for new render output**

In `render/html_test.go`, update `TestHTMLRendererMinimal`:

Replace the assertion at line 43:
```go
if !strings.Contains(html, "--bg-primary: #0d1117") {
	t.Error("missing dark theme CSS")
}
```

With:
```go
if !strings.Contains(html, `data-theme="dark"`) {
	t.Error("missing data-theme attribute on html tag")
}
if !strings.Contains(html, `[data-theme="dark"]`) {
	t.Error("missing dark theme CSS selector")
}
if !strings.Contains(html, `[data-theme="light"]`) {
	t.Error("missing light theme CSS selector")
}
if !strings.Contains(html, `id="theme-toggle"`) {
	t.Error("missing theme toggle button")
}
```

Update `TestHTMLRendererLightTheme` — replace the assertion at line 113:
```go
if !strings.Contains(html, "--bg-primary: #ffffff") {
	t.Error("missing light theme CSS")
}
```

With:
```go
if !strings.Contains(html, `data-theme="light"`) {
	t.Error("missing data-theme='light' attribute")
}
if !strings.Contains(html, `[data-theme="dark"]`) {
	t.Error("missing dark theme CSS (both themes should be embedded)")
}
if !strings.Contains(html, `[data-theme="light"]`) {
	t.Error("missing light theme CSS")
}
```

Add new test:
```go
func TestThemeToggleScript(t *testing.T) {
	d := &diagram.Diagram{
		Title: "Toggle Test",
		Theme: "dark",
		Elements: []diagram.Component{
			&diagram.Note{Text: "test", Style: "muted"},
		},
	}

	r, err := render.NewHTMLRenderer()
	if err != nil {
		t.Fatalf("NewHTMLRenderer: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(d, &buf); err != nil {
		t.Fatalf("Render: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "localStorage") {
		t.Error("missing localStorage in toggle script")
	}
	if !strings.Contains(html, "archdiag-theme") {
		t.Error("missing archdiag-theme key in toggle script")
	}
	if !strings.Contains(html, "data-theme") {
		t.Error("missing data-theme in toggle script")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go test ./render/ -run "TestHTMLRendererMinimal|TestHTMLRendererLightTheme|TestThemeToggleScript" -v`
Expected: FAIL — missing `data-theme`, toggle button, toggle script

- [ ] **Step 3: Update `renderData` struct and `renderHTML` in `render/html.go`**

Change the `renderData` struct (line 69-73):
```go
type renderData struct {
	Diagram      *diagram.Diagram
	ThemeCSS     template.CSS
	DefaultTheme string
	WatchScript  template.JS
}
```

Change `renderHTML` (line 85-100):
```go
func (r *HTMLRenderer) renderHTML(d *diagram.Diagram, w io.Writer, watch bool) error {
	css, err := theme.LoadAll()
	if err != nil {
		return fmt.Errorf("load themes: %w", err)
	}

	defaultTheme := d.Theme
	if defaultTheme == "" {
		defaultTheme = "dark"
	}

	data := renderData{
		Diagram:      d,
		ThemeCSS:     template.CSS(css),
		DefaultTheme: defaultTheme,
	}
	if watch {
		data.WatchScript = template.JS(watchScript)
	}

	return r.Tmpl.ExecuteTemplate(w, "diagram.html.tmpl", data)
}
```

- [ ] **Step 4: Update `diagram.html.tmpl`**

Replace entire file `render/templates/diagram.html.tmpl` with:

```html
<!DOCTYPE html>
<html lang="en" data-theme="{{.DefaultTheme}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Diagram.Title}}</title>
<style>
{{.ThemeCSS}}
</style>
</head>
<body>
<button id="theme-toggle" title="Toggle theme" style="
  position: fixed; top: 16px; right: 16px; z-index: 1000;
  background: var(--bg-card); border: 1px solid var(--border-default);
  border-radius: var(--radius-card); padding: 6px 10px;
  cursor: pointer; font-size: 16px; color: var(--text-primary);
  line-height: 1;
" aria-label="Toggle light/dark theme">{{if eq .DefaultTheme "dark"}}&#9728;&#65039;{{else}}&#127769;{{end}}</button>
<div class="diagram">
  <h1>{{.Diagram.Title}}</h1>
  {{- if .Diagram.Subtitle}}
  <h2>{{.Diagram.Subtitle}}</h2>
  {{- end}}
  {{- range .Diagram.Elements}}
  {{renderComponent .}}
  {{- end}}
</div>
<script>
(function() {
  var btn = document.getElementById('theme-toggle');
  var html = document.documentElement;
  var stored = localStorage.getItem('archdiag-theme');
  if (stored === 'dark' || stored === 'light') {
    html.setAttribute('data-theme', stored);
  }
  function updateIcon() {
    btn.textContent = html.getAttribute('data-theme') === 'dark' ? '\u2600\uFE0F' : '\uD83C\uDF19';
  }
  updateIcon();
  btn.addEventListener('click', function() {
    var current = html.getAttribute('data-theme');
    var next = current === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', next);
    localStorage.setItem('archdiag-theme', next);
    updateIcon();
  });
})();
</script>
{{- if .WatchScript}}
<script>{{.WatchScript}}</script>
{{- end}}
</body>
</html>
```

- [ ] **Step 5: Run all render tests**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go test ./render/ -v`
Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
jj commit -m "feat: add theme toggle to diagram renderer and template"
```

### Task 4: Update index template and watch server

**Files:**
- Modify: `render/templates/index.html.tmpl`
- Modify: `render/watchserver.go`

- [ ] **Step 1: Update `render/templates/index.html.tmpl`**

Replace entire file with:

```html
<!DOCTYPE html>
<html lang="en" data-theme="{{.DefaultTheme}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>archdiag — Diagrams</title>
<style>
{{.ThemeCSS}}
</style>
</head>
<body>
<button id="theme-toggle" title="Toggle theme" style="
  position: fixed; top: 16px; right: 16px; z-index: 1000;
  background: var(--bg-card); border: 1px solid var(--border-default);
  border-radius: var(--radius-card); padding: 6px 10px;
  cursor: pointer; font-size: 16px; color: var(--text-primary);
  line-height: 1;
" aria-label="Toggle light/dark theme">{{if eq .DefaultTheme "dark"}}&#9728;&#65039;{{else}}&#127769;{{end}}</button>
<div class="diagram">
  <h1>archdiag</h1>
  <h2>Available diagrams</h2>
  {{- range .Diagrams}}
  <a href="{{.Path}}" style="display: block; padding: 12px 16px; margin-bottom: 8px; border-radius: var(--radius-card); background: var(--bg-card); border: 1px solid var(--border-default); text-decoration: none; color: var(--text-primary);">
    <strong>{{.Title}}</strong>
    {{- if .Subtitle}}
    <span style="display: block; font-size: 12px; color: var(--text-muted);">{{.Subtitle}}</span>
    {{- end}}
  </a>
  {{- end}}
</div>
<script>
(function() {
  var btn = document.getElementById('theme-toggle');
  var html = document.documentElement;
  var stored = localStorage.getItem('archdiag-theme');
  if (stored === 'dark' || stored === 'light') {
    html.setAttribute('data-theme', stored);
  }
  function updateIcon() {
    btn.textContent = html.getAttribute('data-theme') === 'dark' ? '\u2600\uFE0F' : '\uD83C\uDF19';
  }
  updateIcon();
  btn.addEventListener('click', function() {
    var current = html.getAttribute('data-theme');
    var next = current === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', next);
    localStorage.setItem('archdiag-theme', next);
    updateIcon();
  });
})();
</script>
</body>
</html>
```

- [ ] **Step 2: Update `handleIndex` in `render/watchserver.go`**

Change the `indexData` struct and `handleIndex` method (lines 234-251). Replace:

```go
	themeName := "dark"
	if ws.theme != "" {
		themeName = ws.theme
	}
	css, _ := theme.Load(themeName)

	type indexData struct {
		ThemeCSS template.CSS
		Diagrams []entry
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ws.renderer.Tmpl.ExecuteTemplate(w, "index.html.tmpl", indexData{
		ThemeCSS: template.CSS(css),
		Diagrams: entries,
	})
```

With:

```go
	themeName := "dark"
	if ws.theme != "" {
		themeName = ws.theme
	}
	css, _ := theme.LoadAll()

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
```

- [ ] **Step 3: Run full test suite**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go test ./... -v`
Expected: ALL PASS

- [ ] **Step 4: Commit**

```bash
jj commit -m "feat: add theme toggle to index template and watch server"
```

---

## Chunk 3: Manual Verification

### Task 5: End-to-end verification

- [ ] **Step 1: Render a dark-theme diagram and verify toggle**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go run ./cmd/archdiag render examples/01-hello-world.yaml --theme dark -o /tmp/dark-test.html`

Open `/tmp/dark-test.html` in a browser. Verify:
- Page renders in dark theme
- Toggle button visible in top-right corner showing ☀️
- Clicking toggle switches to light theme (button shows 🌙)
- Clicking again switches back to dark

- [ ] **Step 2: Render a light-theme diagram and verify toggle**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go run ./cmd/archdiag render examples/01-hello-world.yaml --theme light -o /tmp/light-test.html`

Open `/tmp/light-test.html`. Verify:
- Page renders in light theme
- Toggle button shows 🌙
- Clicking switches to dark

- [ ] **Step 3: Verify localStorage persistence**

1. Open `/tmp/dark-test.html`
2. Click toggle to switch to light
3. Reload the page
4. Verify it stays on light theme (localStorage persists)

- [ ] **Step 4: Render a complex diagram to verify no visual regressions**

Run: `cd /Users/patrickvollebregt/Projects/archdiag && go run ./cmd/archdiag render examples/10-full-architecture.yaml --theme dark -o /tmp/full-test.html`

Open and verify all components render correctly in both themes.
