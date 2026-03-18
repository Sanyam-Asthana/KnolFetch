package main

import(
	"fmt"
	"net/http"
	"encoding/xml"
	"io"
	"os"
	"flag"
)

const (
	BOLDGREEN string = "\x1b[1;32m"
	BOLDRED string = "\x1b[1;31m"
	RESET string = "\x1b[0m"
)

type Results struct {
	XMLName       xml.Name       `xml:"results"`
	CategoryFeeds []CategoryFeed `xml:"categoryfeed"`
}

type CategoryFeed struct {
	Name             string `xml:"name,attr"`
	MaxAbstractLength int   `xml:"maxabstractlength,attr"`
	HideTitle        bool   `xml:"hidetitle,attr"`
	HideAuthor       bool   `xml:"hideauthor,attr"`
	HideAbstract     bool   `xml:"hideabstract,attr"`
	HideLink         bool   `xml:"hidelink,attr"`
	Feed Feed `xml:"feed"`
}

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title     string   `xml:"title"`
	Authors   []Author `xml:"author"`
	Summary   string   `xml:"summary"`
	ID        string   `xml:"id"`
	Published string   `xml:"published"`
	Links     []Link   `xml:"link"`
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
	Name             string `xml:"name"`
	MaxResults       int    `xml:"maxresults"`
	MaxAbstractLength int   `xml:"maxabstractlength"`
	HideTitle        bool   `xml:"hidetitle,attr"`
	HideAuthor       bool   `xml:"hideauthor,attr"`
	HideAbstract     bool   `xml:"hideabstract,attr"`
	HideLink         bool   `xml:"hidelink,attr"`
}

type Config struct {
	XMLName    xml.Name   `xml:"config"`
	Categories []Category `xml:"category"`
}

func main() {

	fetchBoolPtr := flag.Bool("fetch", false, "fetch the latest research")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get the home directory!")
	}

	if *fetchBoolPtr {

		var config Config

		path := homeDir + "/.config/knolfetch/config.xml"
		dat, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Cannot open config file!")
			return
		}

		xml.Unmarshal(dat, &config)
		categories := config.Categories

		cacheDir, err := os.UserCacheDir()
		if err != nil {
			fmt.Println("Cannot get cache directory!")
			return
		}
		cacheDir = cacheDir + "/knolfetch/"
		os.MkdirAll(cacheDir, 0755)

		cachePath := cacheDir + "cache.xml"
		f, err := os.Create(cachePath)
		if err != nil {
			fmt.Println("Error while creating cache file:", err)
			return
		}
		defer f.Close()

		f.WriteString("<results>\n")

		for _, category := range categories {

			query := fmt.Sprintf(
				"https://export.arxiv.org/api/query?search_query=cat:%s&start=0&max_results=%d&sortBy=submittedDate&sortOrder=descending",
				category.Name, category.MaxResults,
			)
			resp, err := http.Get(query)
			if err != nil {
				fmt.Println("There was some problem in fetching!")
				return
			}
			defer resp.Body.Close()

			text, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("There was some problem in reading the response!")
				return
			}

			f.WriteString(fmt.Sprintf(
				"<categoryfeed name=\"%s\" maxabstractlength=\"%d\" hidetitle=\"%t\" hideauthor=\"%t\" hideabstract=\"%t\" hidelink=\"%t\">\n",
				category.Name,
				category.MaxAbstractLength,
				category.HideTitle,
				category.HideAuthor,
				category.HideAbstract,
				category.HideLink,
			))
			f.Write(text)
			f.WriteString("</categoryfeed>\n")
		}

		f.WriteString("</results>\n")
		fmt.Println("Fetched the latest research!")

	} else {

		cacheDir, err := os.UserCacheDir()
		if err != nil {
			fmt.Println("Cannot get cache directory!")
			return
		}
		cachePath := cacheDir + "/knolfetch/cache.xml"

		text, err := os.ReadFile(cachePath)
		if err != nil {
			fmt.Println("Error fetching cache!")
			return
		}

		var results Results
		xml.Unmarshal(text, &results)

		for _, cf := range results.CategoryFeeds {
			fmt.Println(BOLDRED + "Results for " + cf.Name + ":" + RESET)
			for _, entry := range cf.Feed.Entries {
				fmt.Println("--------------------")

				if !cf.HideTitle {
					fmt.Println(BOLDGREEN + "Title: " + RESET + entry.Title)
					fmt.Println()
				}

				if !cf.HideAuthor {
					fmt.Print(BOLDGREEN + "Authors: " + RESET)
					for _, author := range entry.Authors {
						fmt.Print(author.Name + ", ")
					}
					fmt.Println()
					fmt.Println()
				}

				if !cf.HideAbstract {
					abstract := entry.Summary
					if cf.MaxAbstractLength > 0 && len(abstract) > cf.MaxAbstractLength {
						abstract = abstract[:cf.MaxAbstractLength]
					}
					fmt.Println(BOLDGREEN + "Abstract: " + RESET + abstract)
					fmt.Println()
				}

				if !cf.HideLink {
					for _, link := range entry.Links {
						if link.Type == "application/pdf" {
							fmt.Println(BOLDGREEN + "Link to PDF: " + RESET + link.Href)
							break
						}
					}
				}
			}
			fmt.Println("--------------------")
		}
	}
}
