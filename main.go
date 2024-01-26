package main

import (
	"encoding/json"
	"fmt"
	"log"
	"novel_crawler/readbook"
	"os"
)

func main() {
	data, err := os.ReadFile("book_info_list.json")
	if err != nil {
		log.Println("read config file err:", err)
		return
	}
	configs := make([]readbook.BookInfo, 0)
	err = json.Unmarshal(data, &configs)
	if err != nil {
		log.Println("covert json err:", err)
		return
	}
	for _, config := range configs {
		if !config.Switch {
			continue
		}
		cookie := config.Cookie
		userAgent := config.UserAgent
		from := config.StartChapter
		num := config.ChapterNum
		readbook.GetBooks(config.Website,
			cookie, userAgent, config.SavePath,
			fmt.Sprintf("%s%v-%v", config.SaveName, from, num), from, num)
	}

}
