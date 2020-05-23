package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func SliceUniq(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

type Blog struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Path       string `json:"posts"`
	Author     string `json:"author"`
	AuthorURL  string `json:"author_url"`
	Year       string `json:"year"`
	Domain     string `json:"domain"`
	Tags       []string
	Posts      []Post
	PrevPage   int
	ActualPage int
	NextPage   int
}

func (u *Blog) Load(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("unable to read configuration file: %v", err))
	}
	err = json.Unmarshal(bytes, &u)
	if err != nil {
		panic(fmt.Sprintf("unable to unmarshal configuration file: %v", err))
	}
}

func (blog *Blog) Compile() {
	f, err := os.Open(blog.Path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		post := Post{}
		post.Blog = blog
		if file.IsDir() {
			log.Printf("[INFO] Compiling %s", file.Name())
			err := post.Compile(blog.Path + "/" + file.Name())
			if err != nil {
				log.Printf("[FAIL] %v", err)
			} else {
				if post.Draft == false {
					for _, tag := range post.Tags {
						blog.Tags = append(blog.Tags, tag)
					}
					blog.Posts = append(blog.Posts, post)
				}
			}
		}
	}
	blog.Tags = SliceUniq(blog.Tags)
	log.Printf("[INFO] %d/%d posts succesfully generated.\n", len(blog.Posts), len(files))
}

func (h *Blog) ServeMain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageNum := 0
	pageNum, err := strconv.Atoi(vars["page"])
	if err != nil {
		pageNum = 1
	}
	indexTemplate, err := Asset("assets/main.html")
	if err != nil {
		log.Printf("[ERR] %s\n", err)
	}
	maxNPosts := 5
	numPages := 1 + (len(h.Posts) / maxNPosts)
	h.ActualPage = pageNum
	h.PrevPage = pageNum
	if h.ActualPage != 1 {
		h.PrevPage = h.ActualPage - 1
	}
	if h.ActualPage != numPages {
		h.NextPage = h.ActualPage + 1
	}
	tmpl, err := template.New("index").Parse(string(indexTemplate))
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, h)
	if err != nil {
		log.Printf("[ERR] %s\n", err)
	}
}

func (h *Blog) NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not found\n")
}