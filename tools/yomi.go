package main

import (
	"encoding/json"
	"github.com/bluele/mecab-golang"
	"io/ioutil"
	"log"
	"strings"
)

type Precure struct {
	Title     string     `json:"title"`
	Subtitles []Subtitle `json:"subtitles"`
}

type Subtitle struct {
	Episode  string   `json:"episode"`
	Raw      string   `json:"raw"`
	Yomi     string   `json:"yomi"`
	Keitaiso []string `json:"keitaiso"`
}

func main() {
	bytes, err := ioutil.ReadFile("../precure-subtitle.json")
	if err != nil {
		log.Fatal(err)
	}

	var precures []Precure
	if err := json.Unmarshal(bytes, &precures); err != nil {
		log.Fatal(err)
	}

	m, err := mecab.New("-Oyomi -d /usr/share/mecab/dic/mecab-ipadic-neologd")
	if err != nil {
		log.Fatal(err)
	}
	defer m.Destroy()

	for pi, p := range precures {
		for si, s := range p.Subtitles {
			tg, err := m.NewTagger()
			if err != nil {
				log.Fatal(err)
			}
			defer tg.Destroy()

			lt, err := m.NewLattice(s.Raw)
			if err != nil {
				log.Fatal(err)
			}
			defer lt.Destroy()

			precures[pi].Subtitles[si].Yomi = strings.TrimRight(tg.Parse(lt), "\n")
		}
	}

	bytes2, err := json.Marshal(precures)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("../precure-subtitle.json", bytes2, 0644)
}
