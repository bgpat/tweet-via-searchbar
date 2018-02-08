package opensearch

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bgpat/tweet-via-searchbar/middleware"
	"github.com/bgpat/twtr"
)

var (
	baseURL = os.Getenv("BASE_URL")
)

type OpenSearch struct {
	XMLName        xml.Name `xml:"OpenSearchDescription"`
	XMLNS          string   `xml:"xmlns,attr"`
	ShortName      string   `xml:"ShortName"`
	LongName       string   `xml:"LongName"`
	Description    string   `xml:"Description"`
	Image          Image    `xml:"Image"`
	Site           string   `xml:"site"`
	InputEncoding  string   `xml:"InputEncoding"`
	OutputEncoding string   `xml:"OutputEncoding"`
	URL            URL      `xml:"Url"`
}

type Image struct {
	XMLName xml.Name `xml:"Image"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Source  string   `xml:",chardata"`
}

type URL struct {
	XMLName  xml.Name `xml:"Url"`
	Type     string   `xml:"type,attr"`
	Method   string   `xml:"method,attr"`
	Template string   `xml:"template,attr"`
	Params   []Param  `xml:"Params"`
}

type Param struct {
	XMLName xml.Name `xml:"Param"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func NewOpenSearch(user *twtr.User, client *middleware.Client) *OpenSearch {
	return &OpenSearch{
		XMLNS:       "http://a9.com/-/spec/opensearch/1.1/",
		ShortName:   fmt.Sprintf("@%s(%s)", user.ScreenName, user.IDStr),
		LongName:    fmt.Sprintf("@%sでツイート", user.ScreenName),
		Description: fmt.Sprintf("検索窓ツイート (@%s)", user.ScreenName),
		Image: Image{
			Width:  16,
			Height: 16,
			Source: user.ProfileImageURLHttps,
		},
		Site:           baseURL + "/",
		InputEncoding:  "UTF-8",
		OutputEncoding: "UTF-8",
		URL: URL{
			Type:     "text/html",
			Method:   "POST",
			Template: baseURL + "/search",
			Params: []Param{
				Param{
					Name:  "q",
					Value: "{searchTerms}",
				},
				Param{
					Name:  "token",
					Value: client.AccessToken.Token,
				},
				Param{
					Name:  "secret",
					Value: client.AccessToken.Secret,
				},
				Param{
					Name:  "redirect",
					Value: client.Config.Redirect,
				},
			},
		},
	}
}

func (o *OpenSearch) ToString() (string, error) {
	buf, err := xml.MarshalIndent(o, "", "    ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(buf), nil
}

func (o *OpenSearch) Render(w http.ResponseWriter) error {
	str, err := o.ToString()
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, str)
	return err
}

func (o *OpenSearch) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = []string{"application/opensearchdescription+xml; charset=utf-8"}
}
