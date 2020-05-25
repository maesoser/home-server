package blog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Name    string    `json:"name"`
	Title   string    `json:"title"`
	Date    time.Time `json:"datetime"`
	Path    string    `json:"path"`
	MDFile  string    `json:"mdfile"`
	Tags    []string  `json:"tags"`
	Draft   bool      `json:"draft"`
	Content string    `json:"content"`
	Blog    *Blog
}

func (p *Post) ReadMarkdown() (string, string, error) {
	bytes, err := ioutil.ReadFile(p.Path + "/" + p.MDFile)
	if err != nil {
		return "nil", "nil", err
	}
	text := string(bytes)
	var re = regexp.MustCompile(`---\ntitle:((.?|\n)*)---\n\n`)
	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("malformed header o post")
	}
	return matches[0], text[len(matches[0]):len(text)], nil

}

func (p *Post) ParseMetadata(data string) error {
	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(line, "title:") {
			p.Title = strings.TrimPrefix(line, "title: ")
		}
		if strings.HasPrefix(line, "draft:") {
			if strings.Contains(line, "False") {
				p.Draft = false
			} else if strings.Contains(line, "false") {
				p.Draft = false
			} else {
				p.Draft = true
			}
		}
		if strings.HasPrefix(line, "tags:") {
			p.Tags = strings.Split(strings.TrimPrefix(line, "tags: "), ",")
		}
	}
	return nil
}

func (p *Post) SummarizeContent(size int) {
	var re = regexp.MustCompile(`(?m)!*\[(.*?)\]\(.*?\)`)
	p.Content = re.ReplaceAllString(p.Content, "$1")
	p.Content = re.ReplaceAllString(p.Content, "$1")
	if len(p.Content) > size {
		p.Content = p.Content[0:size] + "..."
	}
}

func (p *Post) ParseMarkdown(data string) error {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes | parser.Attributes
	parser := parser.NewWithExtensions(extensions)
	p.Content = string(markdown.ToHTML([]byte(data), parser, renderer))

	postTemplate, err := Asset("assets/post.html")
	if err != nil {
		return err
	}
	f, err := os.Create(p.Path + "/index.html")
	if err != nil {
		return err
	}
	wr := bufio.NewWriter(f)
	tmpl, err := template.New("post").Parse(string(postTemplate))
	if err != nil {
		return err
	}
	err = tmpl.Execute(wr, p)
	if err != nil {
		return err
	}
	wr.Flush()
	f.Close()
	p.Content = data
	p.SummarizeContent(512)
	return nil
}

func (p *Post) findMarkdown() error {
	files, err := ioutil.ReadDir(p.Path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".md" {
			p.Date = f.ModTime()
			p.MDFile = f.Name()
			return nil
		}
	}
	return fmt.Errorf("Markdown file not found")
}

// Compile fills the metadata info and then generates HTML file
func (p *Post) Compile(path string) error {
	p.Path = path
	items := strings.Split(path, "/")
	p.Name = items[len(items)-1]
	err := p.findMarkdown()
	if err != nil {
		return err
	}
	metadata, data, err := p.ReadMarkdown()
	if err != nil {
		return err
	}
	err = p.ParseMetadata(metadata)
	if err != nil {
		return err
	}
	err = p.ParseMarkdown(data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) TagExist(tag string) bool{
    for _, t := range p.Tags {
		if tag == t {
			return true
		}
	}
	return false
}
