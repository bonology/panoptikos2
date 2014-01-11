// This program must be executed with the standard Go tools (not with the App
// Engine SDK Go tools) because it depends on Java programs.
package main

import (
	"flag"
	"fmt"
	"github.com/ChristianSiegert/panoptikos/assetcompiler/asset"
	"github.com/ChristianSiegert/panoptikos/assetcompiler/base"
	"github.com/ChristianSiegert/panoptikos/assetcompiler/sanitizer"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

// Command-line flags
var (
	jsCompilationLevel = flag.String("js-compilation-level", asset.JS_COMPILATION_LEVEL_SIMPLE_OPTIMIZATIONS, "Either WHITESPACE_ONLY, SIMPLE_OPTIMIZATIONS or ADVANCED_OPTIMIZATIONS. See https://developers.google.com/closure/compiler/docs/compilation_levels. Advanced optimizations can break your code.")
	verbose            = flag.Bool("verbose", false, "Whether additional information should be displayed after compiling.")
)

var cssCompilerArguments = []string{
	// Ignore non-standard CSS functions and unrecognized CSS properties that
	// we use or else Closure Stylesheets won’t compile our CSS.
	"--allowed-non-standard-function", "color-stop",
	"--allowed-non-standard-function", "progid:DXImageTransform.Microsoft.gradient",
	"--allowed-unrecognized-property", "tap-highlight-color",
}

func main() {
	// token is a Unix timestamp in base 62
	token, err := base.Convert(uint64(time.Now().Unix()), base.DefaultCharacters)

	if err != nil {
		log.Printf("Failed to create token: %s", err)
		return
	}

	// Read index.html
	sourceFilename := "./app/webroot/dev-partials/index.html"
	indexHtml, err := ioutil.ReadFile(sourceFilename)

	if err != nil {
		log.Printf("Failed to read file '%s'.", sourceFilename)
		return
	}

	// CSS
	stylesheetRegExp := regexp.MustCompile("<link href=\"([^\"]+)\" rel=\"stylesheet\" type=\"text/css\">")
	matches := stylesheetRegExp.FindAllStringSubmatch(string(indexHtml), -1)

	cssFilenames := make([]string, 0, len(matches))

	for _, match := range matches {
		url := match[1]

		// If URL begins with “http://”, “https://” or “//”, skip it.
		if found, err := regexp.MatchString("^(https?:)?//", url); found {
			continue
		} else if err != nil {
			log.Printf("Failed while searching CSS links: %s", err)
		}

		cssFilenames = append(cssFilenames, url)
		indexHtml = []byte(strings.Replace(string(indexHtml), match[0], "", 1))
	}

	// JavaScript
	jsRegExp := regexp.MustCompile("<script src=\"([^\"]+)\"></script>")
	matches = jsRegExp.FindAllStringSubmatch(string(indexHtml), -1)

	jsFilenames := make([]string, 0, len(matches))

	for _, match := range matches {
		url := match[1]

		// If URL begins with “http://”, “https://” or “//”, skip it.
		if found, err := regexp.MatchString("^(https?:)?//", url); found {
			continue
		} else if err != nil {
			log.Printf("Failed while searching JS links: %s", err)
		}

		jsFilenames = append(jsFilenames, url)
		indexHtml = []byte(strings.Replace(string(indexHtml), match[0], "", 1))
	}

	// Compile
	cssDestinationBaseName, jsDestinationBaseName, err := compileCssJs(token, cssFilenames, jsFilenames)

	if err != nil {
		log.Printf("assetcompiler: Compiling CSS/JS failed: %s", err)
		return
	}

	cssLink := fmt.Sprintf("<link href=\"%s\" rel=\"stylesheet\" type=\"text/css\">", cssDestinationBaseName)
	jsLink := fmt.Sprintf("<script src=\"%s\"></script>", jsDestinationBaseName)

	indexHtml = []byte(strings.Replace(string(indexHtml), "<!-- COMPILED_CSS_HERE -->", cssLink, 1))
	indexHtml = []byte(strings.Replace(string(indexHtml), "<!-- COMPILED_JS_HERE -->", jsLink, 1))
	indexHtml = sanitizer.RemoveHtmlComments(indexHtml)
	indexHtml = sanitizer.RemoveWhitespace(indexHtml)

	destinationFilename := "./app/webroot/compiled-partials/index-" + token + ".html"
	ioutil.WriteFile(destinationFilename, indexHtml, 0666)
}

// compileCssJs compiles CSS and/or JavaScript. Progress and error messages are
// logged.
func compileCssJs(uniqueKey string, cssSourceFilenames, jsSourceFilenames []string) (cssDestinationBaseName, jsDestinationBaseName string, err error) {
	cssResultChan := make(chan string)
	cssProgressChan := make(chan string)
	cssErrorChan := make(chan error)

	jsResultChan := make(chan string)
	jsProgressChan := make(chan string)
	jsErrorChan := make(chan error)

	defer close(cssResultChan)
	defer close(cssProgressChan)
	defer close(cssErrorChan)

	defer close(jsResultChan)
	defer close(jsProgressChan)
	defer close(jsErrorChan)

	cssDestinationBaseName = uniqueKey + ".css"
	jsDestinationBaseName = uniqueKey + ".js"

	cssDestinationFilename := "/app/webroot/compiled-css/" + cssDestinationBaseName
	jsDestinationFilename := "/app/webroot/compiled-js/" + jsDestinationBaseName

	go asset.CompileCss(cssSourceFilenames, cssDestinationFilename, cssCompilerArguments, cssResultChan, cssProgressChan, cssErrorChan)
	go asset.CompileJs(jsSourceFilenames, jsDestinationFilename, *jsCompilationLevel, *verbose, jsResultChan, jsProgressChan, jsErrorChan)

	for isCompilingCss, isCompilingJs := true, true; isCompilingCss || isCompilingJs; {
		select {
		case _ = <-cssResultChan:
			isCompilingCss = false
		case _ = <-jsResultChan:
			isCompilingJs = false

		case cssProgress := <-cssProgressChan:
			log.Println(cssProgress)
		case jsProgress := <-jsProgressChan:
			log.Println(jsProgress)

		case cssError := <-cssErrorChan:
			log.Println("Compiling CSS failed: ", cssError)
		case jsError := <-jsErrorChan:
			log.Println("Compiling JavaScript failed: ", jsError)
		}
	}

	return cssDestinationBaseName, jsDestinationBaseName, nil
}
