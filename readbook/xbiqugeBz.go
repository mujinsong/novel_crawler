package readbook

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkumza/numcn"
	"github.com/zhshch2002/goreq"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// parseXbiqugeBz: parse www.xbiquge.bz/ satrt:00 12300-...-12399-12400
func parseXbiqugeBz(httpBody interface{}, path, bookName string, ele int) {
	hBody, ok := httpBody.(*goreq.Response)
	if !ok {
		log.Panicf("parseXbiqugeBz-不是*goreq.Response类型")
	}
	doc, err := hBody.HTML()
	if err != nil {
		panic(err)
	}
	body := doc.Find("body")
	contentRead := body.Find("div.content_read")
	divbookname := contentRead.Find("div.bookname")
	title := divbookname.Find("h1").Text()
	fmt.Println(title)
	path = strings.Trim(path, " ")
	title = strings.Trim(title, " ")
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.txt", path, bookName), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("\n\n\n")
	f.WriteString(strconv.Itoa(ele) + title + "\n")
	content := body.Find("div#content") //#是id
	s := content.Text()
	ss := strings.Split(s, " ")
	for i, sr := range ss {
		if i < 2 {
			continue
		}
		if i == 2 {
			sr = strings.TrimPrefix(sr, "！")
		}
		_, err := f.WriteString(sr + "\n")
		if err != nil {
			log.Println("parseXbiqugeBz写入失败")
			return
		}
	}
	return
}

func getTable(web, cookie, userAgent string) (map[int]string, error) {
	mp := make(map[string]string)
	mp["Accept-Encoding"] = "gzip, deflate, br"
	rd, err := ReadReq(web, cookie, userAgent, map[string]string{})
	if err != nil {
		log.Println(web)
		return nil, err
	}
	res := make(map[int]string)
	doc, err := rd.HTML()
	if err != nil {
		log.Println("NewDocument err")
		return nil, err
	}
	body := doc.Find("body")
	list1 := body.Find("div#list")
	listdl := list1.Find("dl")
	listdl.Find("dd").Each(func(i int, selection *goquery.Selection) {
		a := selection.Find("a")
		title := a.Text()
		//fmt.Println(title)
		if strings.Contains(title, "章") && strings.Contains(title, "第") {
			tableNumstr := ""
			fmt.Sscanf(title, "第%s章 %s", &tableNumstr)
			tableNumstr = strings.Replace(tableNumstr, "章", "", 1)
			tableNumstr = strings.Trim(tableNumstr, " ")
			//fmt.Println(tableNumstr)
			tableNum := 0
			if num, err := strconv.Atoi(tableNumstr); err != nil {
				temp, err := numcn.DecodeToInt64(tableNumstr)
				if err != nil {
					log.Panicf("数字转换 err:%v", err)
				}
				tableNum = int(temp)
			} else {
				tableNum = num
			}
			if tableNum == 0 {
				fmt.Println(tableNumstr)
				log.Panicf("获取章节有问题:%s", title)
			}
			val, ext := a.Attr("href")
			if !ext {
				log.Panicf("获取章节有问题:%v", a.Nodes)
			}
			//fmt.Println(title, val)
			res[tableNum] = val
		}
	},
	)
	return res, nil
}

func GetBooksXbiqugeBz(websiteOfBooks, cookie, userAgent, path, bookname string, from, num int) error {
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
	tables, err := getTable(websiteOfBooks, cookie, userAgent)
	if err != nil {
		log.Println("getTable err:", err)
		return err
	}
	//fmt.Println(tables)
	for i := 0; i < num; i++ {
		url := tables[from+i]
		if url == "" {
			log.Println("url is nil err")
			return errors.New("url is nil err")
		}
		web := websiteOfBooks + url
		fmt.Println(web)
		rd, err := ReadReq(web, cookie, userAgent, map[string]string{})
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
