# Light/Dark Theme Toggle

## Summary

Add a runtime light/dark theme toggle to the rendered HTML output. Currently the theme is chosen at render time and baked in. This change embeds both themes and lets the user switch between them in the browser.

## Decisions

- **Default theme:** Whatever was specified at render time (via `--theme` flag or YAML `theme` field)
- **Toggle placement:** Fixed top-right corner of the viewport
- **Persistence:** `localStorage` key `archdiag-theme` — survives reloads
- **Approach:** Dual CSS with `data-theme` attribute toggle on `<html>`

## Changes

### 1. Theme CSS files — extract base styles, scope variables

Both `dark.css` and `light.css` currently contain identical base styles (reset, body, `.diagram`, `.section--half`, media queries) plus theme-specific CSS variables under `:root`.

**Create `theme/css/base.css`** containing the shared styles:

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

**Reduce `dark.css`** to only the scoped variables:

```css
[data-theme="dark"] {
    --bg-primary: #0d1117;
    --bg-card: #161b22;
    /* ... all dark variables ... */
}
```

**Reduce `light.css`** to only the scoped variables:

```css
[data-theme="light"] {
    --bg-primary: #ffffff;
    --bg-card: #f6f8fa;
    /* ... all light variables ... */
}
```

This prevents drift between duplicate base styles and keeps each file focused. The `[data-theme]` selector on `<html>` has the same specificity as `:root` on `<html>` — CSS variables defined there inherit to all descendants including `body`, so `body { background: var(--bg-primary); }` resolves correctly from whichever theme is active.

### 2. Theme loader (`theme/theme.go`)

Add a `LoadAll()` function that returns base + dark + light CSS concatenated. The renderer calls this instead of `Load()` for the diagram/index templates.

```go
func LoadAll() (string, error) {
    base, err := themesFS.ReadFile("css/base.css")
    if err != nil { return "", err }
    dark, err := themesFS.ReadFile("css/dark.css")
    if err != nil { return "", err }
    light, err := themesFS.ReadFile("css/light.css")
    if err != nil { return "", err }
    return string(base) + "\n" + string(dark) + "\n" + string(light), nil
}
```

Keep existing `Load()` function for any code that needs a single theme (backward compat).

### 3. Render data (`render/html.go`)

Change `renderData` struct:

```go
type renderData struct {
    Diagram      *diagram.Diagram
    ThemeCSS     template.CSS   // now contains base + both themes
    DefaultTheme string         // "dark" or "light"
    WatchScript  template.JS
}
```

In `renderHTML()`:

```go
css, err := theme.LoadAll()
if err != nil { return err }

data := renderData{
    Diagram:      d,
    ThemeCSS:     template.CSS(css),
    DefaultTheme: d.Theme,
}
```

If `d.Theme` is empty, default to `"dark"`.

### 4. Template (`render/templates/diagram.html.tmpl`)

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

The toggle button gets an initial text value from the Go template to avoid a flash of empty content before JS runs.

### 5. Watch server (`render/watchserver.go`)

The `handleIndex` method loads a single theme via `theme.Load(themeName)`. Update it to use `theme.LoadAll()` and pass `DefaultTheme`:

```go
type indexData struct {
    ThemeCSS     template.CSS
    DefaultTheme string
    Diagrams     []entry
}
```

```go
css, _ := theme.LoadAll()
themeName := "dark"
if ws.theme != "" {
    themeName = ws.theme
}
ws.renderer.Tmpl.ExecuteTemplate(w, "index.html.tmpl", indexData{
    ThemeCSS:     template.CSS(css),
    DefaultTheme: themeName,
    Diagrams:     entries,
})
```

### 6. Index template (`render/templates/index.html.tmpl`)

Apply the same pattern as `diagram.html.tmpl`:

- Add `data-theme="{{.DefaultTheme}}"` to `<html>`
- Add toggle button
- Add toggle script

### 7. Toggle button behavior

- Shows ☀️ (sun) when in dark mode (click to switch to light)
- Shows 🌙 (moon) when in light mode (click to switch to dark)
- Uses CSS variables for its own styling so it matches the active theme
- `z-index: 1000` ensures it floats above diagram content

### 8. Testing (`render/html_test.go`)

Existing tests need updating since both themes are now always embedded:

- **`TestHTMLRendererMinimal`** (line 43): Currently asserts `--bg-primary: #0d1117`. Update to also verify:
  - `data-theme="dark"` appears on `<html>` tag
  - Both `[data-theme="dark"]` and `[data-theme="light"]` selectors are present
  - Toggle button `id="theme-toggle"` is present
  - Toggle script is present

- **`TestHTMLRendererLightTheme`** (line 113): Currently asserts `--bg-primary: #ffffff`. Update to:
  - Verify `data-theme="light"` on `<html>` tag
  - Both theme CSS blocks are still present (both always embedded)

- **Add `TestThemeToggleScript`**: Verify the toggle script contains `localStorage`, `archdiag-theme`, and `data-theme`.

- **`TestHTMLRendererWatchScript`**: No changes needed — watch script is independent.

- **`TestHTMLRendererAllPrimitives`**: No changes needed — content assertions are unaffected.

## Size impact

Embedding both themes adds ~20 lines of CSS (~400 bytes) to every rendered HTML file. Negligible for self-contained HTML documents that already inline all styles.

## What does NOT change

- CLI interface — `--theme` flag still works, sets the default
- YAML `theme` field — still works, sets the default
- Watch mode — SSE reload preserves localStorage preference
- All component templates — untouched, they already use CSS variables
