package readbook

import (
	"log"
	"strings"
)

// .是class #是id

// 为兼容老代码，body可以为io.ReadCloser，但后续都要传*goreq.Response
func Parse(website, path, bookname string, body interface{}, ele int) {
	if body == nil {
		log.Println("body is nil")
		return
	}
	if strings.Contains(website, "www.beqege.cc") {
		parsebeqegeCc(body, path, bookname, ele)
	} else if strings.Contains(website, "www.xbiquge.bz") {
		parseXbiqugeBz(body, path, bookname, ele)
	} else {
		log.Printf("未知的源网站，请支持\n")
	}
}
