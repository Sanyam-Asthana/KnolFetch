package main

import(
	"fmt"
	"net/http"
	"encoding/xml"
	"io"
	"os"
)

const (
	BOLDGREEN string = "\x1b[1;32m"
	BOLDRED string = "\x1b[1;31m"
	RESET string = "\x1b[0m"

)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title     string   `xml:"title"`
	Authors   []Author `xml:"author"`
	Summary   string   `xml:"summary"`
	ID        string   `xml:"id"`
	Published string   `xml:"published"`
	Links	  []Link   `xml:"link"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
}

type Author struct {
	Name string `xml:"name"`
}

type Category struct {
	Name 	  	  string `xml:"name"`
	MaxResults 	  int    `xml:"maxresults"`
	MaxAbstractLength int    `xml:"maxabstractlength"`
	ShowTitle         bool   `xml:"title,attr"`
	ShowAuthor 	  bool   `xml:"author,attr"`
	ShowAbstract	  bool   `xml:"abstract,attr"`
	ShowLink          bool   `xml:"link,attr"`
}

type Config struct {
	XMLName    xml.Name `xml:"config"`
	Categories []Category `xml:"category"`
}

func main() {

	var config Config

	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Cannot get the home directory!")
	}

	path := homeDir + "/.config/knolfetch/config.xml"
	dat, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Cannot open config file!")
		return
	}

	xml.Unmarshal(dat, &config)

	categories := config.Categories

	for _, category := range categories {

		fmt.Println(BOLDRED + "Results for " + category.Name + ":" + RESET);

		query := fmt.Sprintf("https://export.arxiv.org/api/query?search_query=cat:%s&start=0&max_results=%d&sortBy=submittedDate&sortOrder=descending", category.Name, category.MaxResults)
		resp, err := http.Get(query)

		if err != nil {
			fmt.Println("There was some problem in fetching!")
			return
		}

		text, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("There was some problem in reading the response!")
			return
		}

		var feed Feed
		xml.Unmarshal(text, &feed)

		for _, entry := range feed.Entries {
			fmt.Println("--------------------")
			if category.ShowTitle {
				fmt.Println(BOLDGREEN + "Title: " + RESET + entry.Title)
				fmt.Println()
			}
			
			if category.ShowAuthor {
				fmt.Print(BOLDGREEN + "Authors: " + RESET)
				for _, author := range entry.Authors {
					fmt.Print(author.Name + ", ")
				}
				fmt.Println()
				fmt.Println()
			}
			
			if category.ShowAbstract {
				abstract := entry.Summary
				if len(abstract) > category.MaxAbstractLength {
				    abstract = abstract[:category.MaxAbstractLength]
				}
			
				fmt.Println(BOLDGREEN + "Abstract: " + RESET + abstract)

				fmt.Println()
			}

			if category.ShowLink {
				fmt.Println(BOLDGREEN + "Link to PDF: " + RESET + entry.Links[1].Href)
			}
		}

			resp.Body.Close()
			fmt.Println("--------------------")
		}

}
