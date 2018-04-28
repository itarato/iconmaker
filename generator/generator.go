package generator

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func Generate(title string, url string, icon_file io.Reader) (string, error) {
	re := regexp.MustCompile("[^a-z0-9]+")
	package_name := re.ReplaceAllString(strings.ToLower(title), "_")

	log.Printf("Package \"%s\" referring %s is requested.", title, url)

	file, err := os.Create("/tmp/" + package_name + ".zip")
	if err != nil {
		return "", err
	}

	w := zip.NewWriter(file)
	defer w.Close()

	manifest_source := `{
    "name": "` + title + `",
    "version": "1.0",
    "description": "Shortcut to ` + title + `",
    "manifest_version": 2,
    "background": {
      "scripts": ["background.js"],
      "persistent": false
    },
    "browser_action": {
        "default_title": "Go to ` + title + `",
        "default_icon": "images/sample_128.png"
    },
    "permissions": [
      "activeTab"
    ],
    "icons": {
      "128": "images/sample_128.png"
    }
	}`
	background_js_source := `chrome.browserAction.onClicked.addListener(function(tab) {
		var action_url = "` + url + `";
		chrome.tabs.update(tab.id, {url: action_url});
	});`

	icon_bytes := bytes.NewBufferString("")
	io.Copy(icon_bytes, icon_file)

	var files = []struct {
		Name, Body string
	}{
		{package_name + "/manifest.json", manifest_source},
		{package_name + "/background.js", background_js_source},
		{package_name + "/images/sample_128.png", icon_bytes.String()},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return "", err
		}

		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return "", err
		}
	}

	return package_name, nil
}
