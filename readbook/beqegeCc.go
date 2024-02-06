package readbook

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zhshch2002/goreq"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetBooksBeqegeCc : https://www.beqege.cc/ //start:1 12300-...-12399-123100
func GetBooksBeqegeCc(websiteOfBooks, cookie, userAgent, path, bookname string, from, num int) error {
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
		rd, err := Read(web, cookie, userAgent, map[string]string{})
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

// parsebeqegeCc: parse https://www.beqege.cc/ //start:1 12300-...-12399-123100
func parsebeqegeCc(httpBody interface{}, path, bookName string, ele int) {
	doc := &goquery.Document{}
	var err error
	switch httpBody.(type) {
	case io.ReadCloser:
		hBody := httpBody.(io.ReadCloser)
		defer hBody.Close()
		doc, err = goquery.NewDocumentFromReader(hBody)
		if err != nil {
			log.Panicf("parsebeqegeCc NewDocumentFromReader err:%v", err)
		}
	case *goreq.Response:
		hBody := httpBody.(*goreq.Response)
		doc, err = hBody.HTML()
		if err != nil {
			log.Panicf("parsebeqegeCc HTML() err:%v", err)
		}
	default:
		log.Panicf("parsebeqegeCc 未兼容类型")
	}
	body := doc.Find("body")
	divBookName := body.Find("div.divBookName") //.是class
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
	content := body.Find("div#content") //#是id
	content.Find("p").Each(func(i int, selection *goquery.Selection) {
		f.WriteString(selection.Text() + "\n")
	},
	)
	defer f.Close()
}
