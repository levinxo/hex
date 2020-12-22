package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"html/template"
	ttemplate "text/template"
	"time"
)

type Article struct {
	Title string
	Author string
	Date string
	DateUnix string
	DateSEO string
	Modified string
	ModifiedUnix string
	ModifiedSEO string
	Days string
	OutdatedNotice string
	disqus_identifier string
	DisqusIdentifier string
	Slug string
	Category string
	Tags string
	TagSlice []string
	Content string
	ContentHTML template.HTML
}

type Config struct {
	SiteName string `yaml:"site_name"`
	SiteLang string `yaml:"site_lang"`
	Author string `yaml:"author"`
	TimeLocation string `yaml:"time_location"`
	ProjectDir string `yaml:"project_dir"`
	ContentDir string `yaml:"content_dir"`
	OutputDir string `yaml:"output_dir"`
	DevSiteUrl string `yaml:"dev_site_url"`
	PublishSiteUrl string `yaml:"publish_site_url"`
	ArticleUrl string `yaml:"article_url"`
	ThemeName string `yaml:"theme_name"`
	GoogleAnalytics string `yaml:"google_analytics"`
	DisqusSiteName string `yaml:"disqus_site_name"`
	SwiftTypeInstallCode string `yaml:"swift_type_install_code"`
}

type Meta struct {
	Description string
	PageTitle string
	Author string
	Keywords string
	OgType string
	OgUrl string
	OgSiteName string
	TwCard string
}

type SiteMap struct {
	Loc string
	Lastmod string
	Changefreq string
	Priority string
}

type BlogMaker struct {
	SiteName string
	Author string
	SiteLang string
	PublishSiteUrl string
	ProjectDir string
	ContentDir string
	OutputDir string
	GoogleAnalytics string
	DisqusSiteName string
	SwiftypeInstallCode string
	ArticleUrl string
	ThemeName string
	TimeLocation *time.Location
	SiteMap []SiteMap
	ArticleList []Article
	ReverseArticleList []Article
	Categories map [string] []Article
	CurrentArticle Article
	Meta Meta
	Config Config
}

func NewBlogMaker() *BlogMaker {
	b := &BlogMaker{}
	b.init()

	var fileName,fileContent1,fileContent2 string
	//var nil error
        var ifErr bool

	files, err := ioutil.ReadDir(b.ContentDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files{

		fileName = b.ContentDir+"/"+file.Name()
		fmt.Println(fileName)

                fileContent1, fileContent2, ifErr = b.parseFile(fileName)

                if ifErr == false {
                    continue
                }

		//生成文章头部数据结构
		article := b.parseArticleHeader(fileContent1)

		//文章body fullfill
		_ = setField(&article, "Content", fileContent2)

		article.ContentHTML = template.HTML(fileContent2)

		b.ArticleList = append(b.ArticleList, article)
	}

	for i:=len(b.ArticleList)-1; i>=0; i-- {
		b.ReverseArticleList = append(b.ReverseArticleList, b.ArticleList[i])
		//生成分类
		category := b.ArticleList[i].Category
		if _, ok := b.Categories[category]; !ok {
			b.Categories[category] = []Article{}
		}
		b.Categories[category] = append(b.Categories[category], b.ArticleList[i])
	}

	//渲染模板
	b.render()

	//生成sitemap
	b.sitemap()

	b.moveStatic()
	return b
}

func (b *BlogMaker) parseFile(filePath string) (string, string, bool){
	var fileByte,fileByte1,fileByte2 []byte
	var fileContent1,fileContent2 string
	var err error

	fileByte, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile(`\n\s*\n`)
	loc := reg.FindIndex(fileByte)

        if len(loc) == 0 {
            return fileContent1, fileContent2, false
        }

	fileByte1 = fileByte[0:loc[0]]
	fileByte2 = blackfriday.Run(fileByte[loc[1]:])

	fileContent1 = string(fileByte1)
	fileContent2 = string(fileByte2)

	return fileContent1, fileContent2, true
}

func (b *BlogMaker) parseArticleHeader(content string) (article Article ){
	stringSlice := strings.Split(content, "\n")

	for _, stringSlice := range stringSlice{

		//这里注意每篇文章有多次循环

		idx := strings.Index(stringSlice, ":")

		rowSlice1 := strings.TrimSpace(string([]rune(stringSlice)[0:idx]))
		rowSlice2 := strings.TrimSpace(string([]rune(stringSlice)[idx+1:]))

		if 0 == strings.Compare("Date", rowSlice1) {
			t, _ := time.ParseInLocation("2006-01-02 15:04", rowSlice2, b.TimeLocation)
			_ = setField(&article, "DateUnix", t.String())
			_ = setField(&article, "DateSEO", t.In(b.TimeLocation).Format("2006-01-02T15:04:05")+"+08:00")
			_ = setField(&article, "Days", strconv.Itoa(int(time.Now().Sub(t).Hours()) / 24))
		}

		if 0 == strings.Compare("Modified", rowSlice1) {
			t, _ := time.ParseInLocation("2006-01-02 15:04", rowSlice2, b.TimeLocation)
			_ = setField(&article, "ModifiedUnix", t.String())
			_ = setField(&article, "ModifiedSEO", t.In(b.TimeLocation).Format("2006-01-02T15:04:05")+"+08:00")
			_ = setField(&article, "Days", strconv.Itoa(int(time.Now().Sub(t).Hours()) / 24))
		}

		//信息过时的文章
		if days, _ := strconv.Atoi(article.Days); days > 100 {
			//_ = setField(&article, "OutdatedNotice", article.Days)
			_ = setField(&article, "OutdatedNotice", "")
		} else {
			_ = setField(&article, "OutdatedNotice", "")
		}

		if 0 == strings.Compare("Tags", rowSlice1){
			article.TagSlice = strings.Split(rowSlice2, ",")
		}

		if 0 == strings.Compare("disqus_identifier", rowSlice1){
			_ = setField(&article, "DisqusIdentifier", rowSlice2)
		}

		_ = setField(&article, rowSlice1, rowSlice2)

	}
	return article
}

func (b *BlogMaker) init(){
	fileByte, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	err = yaml.Unmarshal(fileByte, &config)
	if err != nil {
		log.Fatal(err)
	}

	b.Config = config
	b.SiteName = config.SiteName
	b.Author = config.Author
	b.SiteLang = config.SiteLang
	b.TimeLocation, _ = time.LoadLocation(config.TimeLocation)
	b.ProjectDir = config.ProjectDir
	b.ContentDir = config.ContentDir
	b.OutputDir = config.OutputDir
	b.PublishSiteUrl = config.PublishSiteUrl
	b.ArticleUrl = config.ArticleUrl
	b.ThemeName = config.ThemeName
	b.GoogleAnalytics = config.GoogleAnalytics
	b.DisqusSiteName = config.DisqusSiteName
	b.SwiftypeInstallCode = config.SwiftTypeInstallCode

	b.Categories = make(map [string] []Article)
	_ = os.MkdirAll(b.OutputDir, 0755)
}


func (b *BlogMaker) render(){
	var meta Meta
	tplPathPrefix := "theme/"+b.ThemeName+"/template"

	//生成首页
	meta.PageTitle = b.SiteName
	meta.Author = b.Author
	meta.OgType = "article"
	meta.OgUrl = "."
	meta.OgSiteName = b.SiteName
	meta.TwCard = "summary"
	b.Meta = meta

	fIndex, _ := os.Create(b.OutputDir+"/index.html")
	defer fIndex.Close()
	tIndex, _ := template.ParseFiles(
		tplPathPrefix+"/index.html",
		tplPathPrefix+"/component/header.html",
		tplPathPrefix+"/component/navi.html",
		tplPathPrefix+"/component/footer.html",
		)
	err := tIndex.Execute(fIndex, b)
	if err != nil {
		log.Fatal(err)
	}

	//生成归档页
	meta = Meta{}
	meta.Description = "Full archives of "+b.SiteName
	meta.PageTitle = "All Posts · "+b.SiteName
	meta.Author = b.Author
	meta.OgType = "article"
	meta.OgUrl = "./archives.html"
	meta.OgSiteName = b.SiteName
	meta.TwCard = "summary"
	b.Meta = meta

	fArchive, _ := os.Create(b.OutputDir+"/archives.html")
	defer fArchive.Close()
	tArchive, _ := template.ParseFiles(
		tplPathPrefix+"/archives.html",
		tplPathPrefix+"/component/header.html",
		tplPathPrefix+"/component/navi.html",
		tplPathPrefix+"/component/footer.html",
		)
	_ = tArchive.Execute(fArchive, b)

	//生成分类页
	meta = Meta{}
	meta.Description = "All categories of "+b.SiteName
	meta.PageTitle = "All Categories · "+b.SiteName
	meta.Author = b.Author
	meta.OgType = "article"
	meta.OgUrl = "./categories.html"
	meta.OgSiteName = b.SiteName
	meta.TwCard = "summary"
	b.Meta = meta

	fCategories, _ := os.Create(b.OutputDir+"/categories.html")
	defer fCategories.Close()
	tCategories, _ := template.ParseFiles(
		tplPathPrefix+"/categories.html",
		tplPathPrefix+"/component/header.html",
		tplPathPrefix+"/component/navi.html",
		tplPathPrefix+"/component/footer.html",
	)
	_ = tCategories.Execute(fCategories, b)

	//生成单独的文章页
	for _, article := range b.ArticleList {

		meta = Meta{}
		meta.Description = ""
		meta.PageTitle = article.Title + " · " + b.SiteName
		meta.Author = article.Author
		meta.Keywords = article.Tags + "," + article.Category
		meta.OgType = "article"
		meta.OgUrl = "../../"+b.ArticleUrl+"/"+article.Slug+"/"
		meta.OgSiteName = b.SiteName
		meta.TwCard = "summary"
		b.Meta = meta

		_ = os.MkdirAll(b.OutputDir+"/"+b.ArticleUrl+"/"+article.Slug, 0755)
		fArticle, _ := os.Create(b.OutputDir+"/"+b.ArticleUrl+"/"+article.Slug+"/index.html")

		tArticle, _ := template.ParseFiles(
			tplPathPrefix+"/article.html",
			tplPathPrefix+"/component/header.html",
			tplPathPrefix+"/component/navi.html",
			tplPathPrefix+"/component/footer.html",
			)

		b.CurrentArticle = article
		_ = tArticle.Execute(fArticle, b)

		fArticle.Close()
	}
}

func (b *BlogMaker) sitemap(){

	//生成sitemap数据结构
	for _, article := range b.ReverseArticleList {
		var sMap SiteMap

		sMap.Loc = b.PublishSiteUrl+"/"+b.ArticleUrl+"/"+article.Slug+"/"
		sMap.Lastmod = article.DateSEO
		if article.ModifiedSEO != "" {
			sMap.Lastmod = article.ModifiedSEO
		}
		sMap.Changefreq = "monthly"
		sMap.Priority = "0.5"

		b.SiteMap = append(b.SiteMap, sMap)
	}

	//渲染
	f, _ := os.Create(b.OutputDir+"/sitemap.xml")
	defer f.Close()
	t, _ := ttemplate.ParseFiles("sitemap.xml")
	_ = t.Execute(f, b)
}

func (b *BlogMaker) moveStatic(){

	files, _ := ioutil.ReadDir(b.ProjectDir+"/theme/"+b.ThemeName+"/css")
	_ = os.MkdirAll(b.OutputDir+"/theme/css", 0755)
	for _, file := range files{
		src, _ := os.Open(b.ProjectDir+"/theme/"+b.ThemeName+"/css/"+file.Name())
		dst, _ := os.OpenFile(b.OutputDir+"/theme/css/"+file.Name(), os.O_WRONLY|os.O_CREATE, 0644)
		_, _ = io.Copy(dst, src)
		src.Close()
		dst.Close()
	}

}
