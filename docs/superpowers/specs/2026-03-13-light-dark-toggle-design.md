# Light/Dark Theme Toggle

## Summary

Add a runtime light/dark theme toggle to the rendered HTML output. Currently the theme is chosen at render time and baked in. This change embeds both themes and lets the user switch between them in the browser.

## Decisions

- **Default theme:** Whatever was specified at render time (via `--theme` flag or YAML `theme` field)
- **Toggle placement:** Fixed top-right corner of the viewport
- **Persistence:** `localStorage` key `archdiag-theme` — survives reloads
- **Approach:** Dual CSS with `data-theme` attribute toggle on `<html>`

## Changes

### 1. Theme CSS files (`theme/css/dark.css`, `theme/css/light.css`)

Both files currently define variables under `:root`. Change to scope under `[data-theme]`:

- `dark.css`: `:root {` → `[data-theme="dark"] {`
- `light.css`: `:root {` → `[data-theme="light"] {`

Shared base styles (reset, body, `.diagram`, media queries) are identical in both files. Extract them to remain under `:root` or `html` so they apply regardless of theme. Keep one copy — in `dark.css` is fine since it's the primary theme file, or create a `base.css`. Decision: keep shared styles in both files for now (they're small and identical), avoid adding a third file. The `[data-theme]` selector scopes only the CSS variables block; everything else stays as-is.

### 2. Theme loader (`theme/theme.go`)

Add a `LoadBoth()` function that returns dark CSS and light CSS strings. Alternatively, the renderer can just call `Load("dark")` and `Load("light")` separately — no new function needed. Keep it simple: two `Load()` calls in the renderer.

### 3. Render data (`render/html.go`)

Change `renderData` struct:

```go
type renderData struct {
    Diagram      *diagram.Diagram
    DarkCSS      template.CSS
    LightCSS     template.CSS
    DefaultTheme string          // "dark" or "light"
    WatchScript  template.JS
}
```

In `renderHTML()`, load both themes:

```go
darkCSS, err := theme.Load("dark")
lightCSS, err := theme.Load("light")
```

Set `DefaultTheme` from `d.Theme`.

### 4. Template (`render/templates/diagram.html.tmpl`)

```html
<!DOCTYPE html>
<html lang="en" data-theme="{{.DefaultTheme}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Diagram.Title}}</title>
<style>
{{.DarkCSS}}
{{.LightCSS}}
</style>
</head>
<body>
<button id="theme-toggle" title="Toggle theme" style="
  position: fixed; top: 16px; right: 16px; z-index: 1000;
  background: var(--bg-card); border: 1px solid var(--border-default);
  border-radius: var(--radius-card); padding: 6px 10px;
  cursor: pointer; font-size: 16px; color: var(--text-primary);
  line-height: 1;
" aria-label="Toggle light/dark theme">
</button>
<div class="diagram">
  <!-- existing content -->
</div>
<script>
(function() {
  var btn = document.getElementById('theme-toggle');
  var html = document.documentElement;
  var defaultTheme = '{{.DefaultTheme}}';
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

### 5. Index template (`render/templates/index.html.tmpl`)

Apply the same pattern if the index page uses theme styling. Needs investigation — if it embeds theme CSS, add the toggle there too.

## Toggle Button Behavior

- Shows ☀️ (sun) when in dark mode (click to switch to light)
- Shows 🌙 (moon) when in light mode (click to switch to dark)
- Uses CSS variables for its own styling so it matches the active theme
- `z-index: 1000` ensures it floats above diagram content

## What Does NOT Change

- CLI interface — `--theme` flag still works, sets the default
- YAML `theme` field — still works, sets the default
- Watch mode — still works, SSE reload preserves localStorage preference
- All component templates — untouched, they already use CSS variables

## Testing

- Render a diagram with `--theme dark`, verify toggle switches to light and back
- Render with `--theme light`, verify toggle defaults to light
- Switch theme, reload page, verify localStorage persists the choice
- Verify toggle button is visible and doesn't overlap diagram content
- Existing `html_test.go` tests need updating for new `renderData` fields
