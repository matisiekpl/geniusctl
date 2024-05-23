package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/term"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Cyan = "\033[36m"
var White = "\033[97m"
var DarkGray = "\u001b[38;5;238m"

func printManyTimes(contents string, n int, color string) {
	out := color
	for i := 0; i < n; i++ {
		out += contents
	}
	fmt.Print(out)
}

func printLyrics(header string, lyrics string) {
	offset := 2
	lines := strings.Split(lyrics, "\n")
	width, _, _ := term.GetSize(0)

	printManyTimes("─", offset, DarkGray)
	fmt.Print("┬")
	printManyTimes("─", width-offset-1, DarkGray)
	fmt.Println()
	printManyTimes(" ", offset, DarkGray)
	fmt.Print("│")
	fmt.Println(White + " " + header)
	printManyTimes("─", offset, DarkGray)
	fmt.Print("┼")
	printManyTimes("─", width-offset-1, DarkGray)

	for _, line := range lines {
		printManyTimes(" ", offset, DarkGray)
		fmt.Print("│")
		if strings.Contains(line, "[") {
			fmt.Println(" " + Cyan + line)
		} else {
			fmt.Println(" " + White + line)
		}
	}

	printManyTimes("─", offset, DarkGray)
	fmt.Print("┴")
	printManyTimes("─", width-offset-1, DarkGray)
	fmt.Print(Reset)
}

func printSong(song Song) {
	resp, _ := http.Get(song.Url)
	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	doc, _ := htmlquery.Parse(strings.NewReader(strings.ReplaceAll(string(html), "<br/>", "\n")))
	list := htmlquery.Find(doc, "//div[@data-lyrics-container]")
	lyrics := ""
	for _, el := range list {
		lyrics += htmlquery.InnerText(el)
	}
	printLyrics(song.Title+" - "+song.ArtistNames, lyrics)
}

type Song struct {
	Id          int    `json:"id"`
	Path        string `json:"path"`
	Url         string `json:"url"`
	ArtistNames string `json:"artist_names"`
	Title       string `json:"title"`
}

type Hit struct {
	Index  string `json:"index"`
	Type   string `json:"type"`
	Result Song   `json:"result"`
}

type Section struct {
	Type string `json:"type"`
	Hits []Hit  `json:"hits"`
}

type SearchResponseResponse struct {
	Sections []Section `json:"sections"`
}

type SearchResponse struct {
	Response SearchResponseResponse `json:"response"`
}

func search(query string) Song {
	resp, _ := http.Get("https://genius.com/api/search/multi?per_page=5&q=" + url.QueryEscape(query))
	defer resp.Body.Close()
	responseData, _ := ioutil.ReadAll(resp.Body)
	var result SearchResponse
	json.Unmarshal(responseData, &result)
	if len(result.Response.Sections) < 2 {
		fmt.Println(Red + "Song not found")
		os.Exit(1)
	}
	if len(result.Response.Sections[1].Hits) < 1 {
		fmt.Println(Red + "Song not found")
		os.Exit(1)
	}
	return result.Response.Sections[1].Hits[0].Result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(Red + "Song title not given")
		os.Exit(1)
	}
	printSong(search(strings.Join(os.Args[1:], " ")))
}
