package radiko

import (
	"net/http"

	"golang.org/x/net/html"
)

const (
	areaURL = "http://radiko.jp/area"
)

// AreaID returns areaID.
func AreaID() (string, error) {
	resp, err := http.Get(areaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	return processSpanNode(doc), nil
}

func processSpanNode(n *html.Node) string {
	var areaID string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" && len(n.Attr) > 0 {
			areaID = n.Attr[0].Val
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	return areaID
}
