package readbook

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Parse(website, path, bookname string, body io.ReadCloser, ele int) {
	if body == nil {
		log.Println("body is nil")
		return
	}
	defer body.Close()
	if strings.Contains(website, "www.beqege.cc") {
		parsebeqegeCc(body, path, bookname, ele)
	} else {
		log.Printf("未知的源网站，请支持\n")
	}
}

// parsebeqegeCc: parse https://www.beqege.cc/
func parsebeqegeCc(httpBody io.ReadCloser, path, bookName string, ele int) {
	doc, err := goquery.NewDocumentFromReader(httpBody)
	if err != nil {
		panic(err)
	}
	body := doc.Find("body")
	divBookName := body.Find("div.divBookName")
	title := divBookName.Find("h1").Text()
	fmt.Println(title)
	path = strings.Trim(path, " ")
	title = strings.Trim(title, " ")
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.txt", path, bookName), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("\n\n\n")
	f.WriteString(strconv.Itoa(ele) + title + "\n")
	content := body.Find("div#content")
	content.Find("p").Each(func(i int, selection *goquery.Selection) {
		f.WriteString(selection.Text() + "\n")
	},
	)
	defer f.Close()
}
func Read(website, cookie, userAgent string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", website, nil)
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	} else {
		log.Println("cookie nil")
	}
	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}
	//req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err := client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	log.Println("body:", resp.Body)
	return resp.Body, nil
}
func GetBooks(websiteOfBooks, cookie, userAgent, path, bookname string, from, num int) error {
	f, err := os.Create(path + bookname + ".txt")
	if err != nil {
		log.Println(err)
		return err
	}
	err = f.Close()
	if err != nil {
		log.Println("close file err:", err)
		return err
	}
	for i := 0; i < num; i++ {
		web := websiteOfBooks + fmt.Sprintf("%v.html", from+i)
		rd, err := Read(web, cookie, userAgent)
		if err != nil {
			log.Println(web)
			log.Println(from, " ", i, "err:", err)
			return err
		}
		Parse(websiteOfBooks, path, bookname, rd, from+i)
		time.Sleep(3 * time.Second)
	}
	return nil
}
