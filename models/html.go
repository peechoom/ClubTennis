package models

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func IsCleanHTML(data string) bool {
	resultChan := make(chan bool)
	var wg sync.WaitGroup

	var unsafes = []string{
		"<script>",
		"</script>",
		"onload",
		"onerror",
		"onclick",
		"onmouseover",
		"onfocus",
		"onblur",
		"onsubmit",
		"onchange",
		"onkeypress",
		"onkeydown",
		"onkeyup",
		"onmousedown",
		"onmouseup",
		"style",
		"<iframe>",
		"<object>",
		"<embed>",
		"<applet>",
		"<meta http-equiv=\"refresh\">",
		"<form>",
		"<input>",
		"<button>",
		"<select>",
		"<textarea>",
		"<!--",
	}

	for _, s := range unsafes {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			resultChan <- strings.Contains(data, s)
		}(s)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var ret bool = false
	for b := range resultChan {
		if b {
			ret = b
		}
	}
	return !ret
}

// strips image objects from an html string and replaces them with <img src={link}/> tags.
func StripImages(htmlData *string, path string) ([]*Image, error) {
	path += "/"
	re := regexp.MustCompile(`data:image/([^\s";]+);base64,([^\s"]+)`)

	doc, err := html.Parse(strings.NewReader(*htmlData))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, err
	}

	// Replace Base64 data
	var images []*Image
	findAndReplace(doc, path, &images, re)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, fmt.Errorf("error rendering HTML: %v", err)
	}
	*htmlData = buf.String()

	return images, nil
}

func findAndReplace(n *html.Node, path string, images *[]*Image, re *regexp.Regexp) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for i, attr := range n.Attr {
			if attr.Key == "src" && strings.HasPrefix(attr.Val, "data:image") && strings.Contains(attr.Val, "base64") {
				match := re.FindSubmatch([]byte(attr.Val))
				var data string
				var ex string
				if len(match) > 2 {
					ex = string(match[1])
					data = string(match[2])
				} else {
					break
				}

				img := NewImage(data, ex)
				*images = append(*images, img)
				n.Attr[i].Val = path + img.FileName
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findAndReplace(c, path, images, re)
	}
}
