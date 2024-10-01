package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// urlは最終的に引数から撮れるようにする
	url := "https://natalie.mu/comic/anime/season/2024-autumn"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ab := doc.Find("div.NA_article_body")
	ab.Find("h2").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		// 放映時間はh2の次の次の要素に含まれる
		table := s.Next().Next()
		// 最速放映時間が含まれない場合もあるので containsで走査
		z := table.Find("tbody tr:contains('最速')").Find("td p")
		log.Println("-START-")
		log.Println(title)
		log.Println(z.Text())
		log.Println("-END-")
	})

}
