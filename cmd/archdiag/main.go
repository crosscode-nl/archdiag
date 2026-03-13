package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/crosscode-nl/archdiag/parse"
	"github.com/crosscode-nl/archdiag/render"
)

func main() {
	root := &cobra.Command{
		Use:   "archdiag",
		Short: "Generate HTML architecture diagrams from YAML",
	}

	// render command
	var outputDir string
	var lightTheme, darkTheme bool

	renderCmd := &cobra.Command{
		Use:   "render <path>",
		Short: "Render YAML diagram(s) to HTML",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			theme := ""
			if lightTheme {
				theme = "light"
			} else if darkTheme {
				theme = "dark"
			}
			return runRender(args[0], outputDir, theme)
		},
	}
	renderCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	renderCmd.Flags().BoolVar(&lightTheme, "light", false, "Use light theme")
	renderCmd.Flags().BoolVar(&darkTheme, "dark", false, "Use dark theme")

	root.AddCommand(renderCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate YAML diagram(s) without rendering",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0])
		},
	}
	root.AddCommand(validateCmd)

	var watchPort int
	var watchOpen bool
	var watchLight, watchDark bool

	watchCmd := &cobra.Command{
		Use:   "watch <path>",
		Short: "Watch YAML and serve with live reload",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			theme := ""
			if watchLight {
				theme = "light"
			} else if watchDark {
				theme = "dark"
			}
			return runWatch(args[0], watchPort, watchOpen, theme)
		},
	}
	watchCmd.Flags().IntVarP(&watchPort, "port", "p", 3210, "HTTP server port")
	watchCmd.Flags().BoolVar(&watchOpen, "open", false, "Open browser on start")
	watchCmd.Flags().BoolVar(&watchLight, "light", false, "Use light theme")
	watchCmd.Flags().BoolVar(&watchDark, "dark", false, "Use dark theme")
	root.AddCommand(watchCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRender(path, outputDir, themeOverride string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	var files []string
	if info.IsDir() {
		matches, err := filepath.Glob(filepath.Join(path, "*.yaml"))
		if err != nil {
			return err
		}
		files = matches
	} else {
		files = []string{path}
	}

	if len(files) == 0 {
		return fmt.Errorf("no YAML files found in %s", path)
	}

	renderer, err := render.NewHTMLRenderer()
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := renderFile(f, outputDir, themeOverride, renderer); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}
	}

	return nil
}

func renderFile(yamlPath, outputDir, themeOverride string, renderer *render.HTMLRenderer) error {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	d, err := parse.Parse(data)
	if err != nil {
		return err
	}

	if errs := parse.Validate(d); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "%s: %v\n", yamlPath, e)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(errs))
	}

	if themeOverride != "" {
		d.Theme = themeOverride
	}

	// Determine output path
	outPath := strings.TrimSuffix(yamlPath, filepath.Ext(yamlPath)) + ".html"
	if outputDir != "" {
		outPath = filepath.Join(outputDir, filepath.Base(outPath))
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := renderer.Render(d, outFile); err != nil {
		outFile.Close()
		os.Remove(outPath)
		return err
	}

	fmt.Printf("rendered %s\n", outPath)
	return nil
}

func runWatch(path string, port int, openBrowser bool, themeOverride string) error {
	renderer, err := render.NewHTMLRenderer()
	if err != nil {
		return err
	}

	ws := render.NewWatchServer(renderer, themeOverride, port)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return ws.WatchDir(path, openBrowser)
	}
	return ws.WatchFile(path, openBrowser)
}

func runValidate(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	var files []string
	if info.IsDir() {
		matches, err := filepath.Glob(filepath.Join(path, "*.yaml"))
		if err != nil {
			return err
		}
		files = matches
	} else {
		files = []string{path}
	}

	hasErrors := false
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			hasErrors = true
			continue
		}

		d, err := parse.Parse(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			hasErrors = true
			continue
		}

		if errs := parse.Validate(d); len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "%s: %v\n", f, e)
			}
			hasErrors = true
		} else {
			fmt.Printf("%s: valid\n", f)
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}
	return nil
}
