package cmd

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ytakky2014/ani-sche-natalie/anime"
	"github.com/ytakky2014/ani-sche-natalie/google"
)

// AnimeSchedule は指定されたurlからパースしたAnimeのスケジュールをGoogleカレンダーへ登録する
func AnimeSchedule() {
	// urlは最終的に引数から撮れるようにする
	// calenderIDも引数から取る
	url := os.Args[1]
	calendarID := os.Args[2]
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
	Animes := []anime.Anime{}
	ab.Find("h2").EachWithBreak(func(i int, s *goquery.Selection) bool {
		an := anime.Anime{}
		title := s.Text()
		an.Title = title
		// 放映時間はh2の次の次の要素に含まれる
		table := s.Next().Next()
		// 最速放映時間が含まれない場合もあるので containsで走査
		z := table.Find("tbody tr:contains('最速放送')").Find("td p")
		if z.Text() == "" {
			log.Printf("%s は放送日未確定のためSKIP", an.Title)
			return true
		}
		broadcast, statTime, endTime, err := sanitize(z.Text())
		if err != nil {
			log.Printf("%s は放送日の処理でエラーのためSKIP err: %s", an.Title, err)
			return true
		}
		an.Broadcast = broadcast
		an.StartTime = statTime
		an.EndTime = endTime
		Animes = append(Animes, an)
		return true
	})

	log.Println("-----------------------------------")

	for _, anime := range Animes {
		log.Printf("%s は放送日: %s ~ %s , 放送局 %s", anime.Title, anime.StartTime, anime.EndTime, anime.Broadcast)
		out := google.CreateCalendar(anime, calendarID)
		log.Println(out)
	}

}

// sanitize はscheduleから放送局とスケジュールを取得する
// schedule は"2024年10月3日（木）21:00～（AT-X）" のように与えられる
func sanitize(s string) (broadcast, startTime, endTime string, err error) {
	// 最後に含まれる（を起点に放送局と放送日時を分割
	in := strings.LastIndex(s, "（")
	broadcast = s[in:]
	// 放送局に含まれる不要な（を削除
	broadcast = strings.TrimLeft(broadcast, "（")
	broadcast = strings.TrimRight(broadcast, "）")

	start, end, err := converTime(s[:in])
	if err != nil {
		return broadcast, "", "", err
	}

	return broadcast, start, end, nil
}

// convertTime は放送日時から開始時刻と終了時刻を2006-01-02T15:04:05-07:00形式で返す
func converTime(schedule string) (startTime, endTime string, err error) {
	if schedule == "" {
		return "", "", nil
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	parseLayout := "2006年1月2日"
	// 不要な～を削除
	schedule = strings.TrimRight(schedule, "～")
	// 不要な台を削除
	schedule = strings.ReplaceAll(schedule, "台", "")
	// goでは25時のような時刻は直接parse出来ないので、考慮する
	// 年月日のみ取得
	ymd := strings.SplitAfter(schedule, "日")
	ymdt, err := time.ParseInLocation(parseLayout, ymd[0], loc)
	if err != nil {
		return "", "", err
	}

	// 時間のみ取得
	hm := strings.Split(ymd[len(ymd)-1], "）")
	h := strings.Split(hm[len(hm)-1], ":")

	hi, err := strconv.ParseInt(h[0], 10, 64)
	if err != nil {
		return "", "", err
	}
	// 24時を過ぎている場合、24時間引いて1日を足す
	if hi >= 24 {
		hi = hi - 24
		ymdt = ymdt.AddDate(0, 0, 1)
	}

	m, err := strconv.ParseInt(h[1], 10, 64)
	if err != nil {
		return "", "", err
	}
	ymdt = ymdt.Add(time.Hour * time.Duration(hi))
	ymdt = ymdt.Add(time.Minute * time.Duration(m))

	end := ymdt.Add(30 * time.Minute)

	outputlayout := "2006-01-02T15:04:05-07:00"
	return ymdt.Format(outputlayout), end.Format(outputlayout), nil
}
