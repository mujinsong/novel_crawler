package readbook

type BookInfo struct {
	Website      string `json:"website"`
	Cookie       string `json:"cookie"`
	UserAgent    string `json:"user_agent"`
	SavePath     string `json:"save_path"`
	SaveName     string `json:"save_name"`
	StartChapter int    `json:"start_chapter"`
	ChapterNum   int    `json:"chapter_num"`
	Switch       bool   `json:"switch"`
}
