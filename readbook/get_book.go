package readbook

import (
	"fmt"
	"github.com/zhshch2002/goreq"
	"io"
	"log"
	"net/http"
	"strings"
)

func ReadReq(website, cookie, userAgent string, header map[string]string) (*goreq.Response, error) {
	c := http.Cookie{
		Name:  "cookie",
		Value: cookie,
	}
	rsp := goreq.Get(website).AddHeaders(header).AddCookie(&c).AddHeader("User-Agent", userAgent).Do()
	if rsp.Err != nil {
		log.Println("ReadReq err:", rsp.Err)
		return nil, rsp.Err
	}
	return rsp, nil
}

// Deprecated: Read :读取web
func Read(website, cookie, userAgent string, header map[string]string) (io.ReadCloser, error) {
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
	for s, s2 := range header {
		req.Header.Add(s, s2)
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
	if strings.Contains(websiteOfBooks, "www.beqege.cc") {
		err := GetBooksBeqegeCc(websiteOfBooks, cookie, userAgent, path, bookname, from, num)
		if err != nil {
			log.Println(err)
			return err
		}
	} else if strings.Contains(websiteOfBooks, "www.xbiquge.bz") {
		err := GetBooksXbiqugeBz(websiteOfBooks, cookie, userAgent, path, bookname, from, num)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {

	}
	return nil
}
