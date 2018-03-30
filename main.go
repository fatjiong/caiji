package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

type book struct {
	author string
	title string
	urlList []string
}

type chapter struct {
	title string
	content string
}

//字符串转换方法
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//获取书籍信息(书名，作者，章节链接)
func bookInfo(url string) *book{
	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Println(err)
	}

	//获取书名
	title := doc.Find("h1").Text()
	title = ConvertToString(title, "gbk", "utf-8")

	//获取作者
	author := doc.Find("p").First().Text()
	author = ConvertToString(author, "gbk", "utf-8")
	authors := strings.Split(author, "：")
	author = authors[1]

	//获取章节列表
	urlList := []string{}
	doc.Find("div#list a").Each(func(i int, selection *goquery.Selection) {
		if i>8 {
			href,err :=selection.Attr("href")
			if err!=true {
				fmt.Println(err)
			}
			urlList = append(urlList,href)
		}
	})

	return &book{
		author,
			title,
			urlList,
	}
}

//获取章节信息
func chapterInfo(url string) *chapter{
	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Println(err)
	}

	title := doc.Find("title").Text()
	title = ConvertToString(title, "gbk", "utf-8")
	titles := strings.Split(title, "_")
	title = titles[0]

	content:=	doc.Find("div#content").Text()
	content = ConvertToString(content, "gbk", "utf-8")
	content = strings.Replace(content,"聽","　",-1)

	return &chapter{
		title,
		content,
	}
}


func main() {
	//默认链接
	baseUrl := "http://www.81xsw.com"
	//书籍主页
	bookUrl := baseUrl + "/0_355/"
	//获取书籍信息
	bookInfo := bookInfo(bookUrl)

	fmt.Println(bookInfo.title+"，本书是由:"+bookInfo.author+"所著")

	timeStart := time.Now()
	//循环章节链接，获取内容
	for _,item := range bookInfo.urlList {
		chapterUrl:=baseUrl+item
		chapter := chapterInfo(chapterUrl)
		fmt.Println("章节:"+chapter.title+"的链接地址是:"+item)
		// 剩余的任务就是组装写入数据库
	}

	timeEnd := time.Now()
	wasteTime := time.Time.Sub(timeEnd,timeStart)

	fmt.Println(wasteTime)
}
