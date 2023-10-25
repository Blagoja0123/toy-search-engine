package internal

import (
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
)

type Result struct {
	NumKeyWords map[string]int
	Path        string
	Accuracy    int
}

func extractKeyWords(s string) []string {
	commonExp := []string{"how", "what", "why", "is", "to", "do"}
	var keyWords []string
	for _, exp := range commonExp {
		s = strings.ReplaceAll(s, exp, "")
	}
	keyWords = append(keyWords, strings.Split(s, " ")...)
	return keyWords
}

func initKWMap(keyWords []string) map[string]int {
	res := make(map[string]int)
	for _, keyWord := range keyWords {
		res[keyWord] = 0
	}
	return res
}

func Search(query string, wg *sync.WaitGroup) []Result {
	var Res []Result
	//wg := sync.WaitGroup{}
	channelRes := make(chan Result)
	keyWords := extractKeyWords(query)
	dirs := getSubDirs("./static")

	for _, dir := range dirs {
		wg.Add(1)
		go func(dir string) {
			files := loadFiles(dir)
			htmlFiles := readFiles(dir, files)
			execSearch(htmlFiles, keyWords, files, channelRes)
			wg.Done()
		}(dir)
	}
	go func() {
		wg.Wait()
		close(channelRes)
	}()
	for res := range channelRes {
		Res = append(Res, res)
	}
	return Res
}

func execSearch(htmlFiles []*tree, keyWords []string, files []string, channelRes chan Result) {
	for i := 0; i < len(htmlFiles); i++ {
		wordMap := initKWMap(keyWords)
		tokens := tokenize(htmlFiles[i].root)
		for _, token := range tokens {
			if slices.Contains(keyWords, token) {
				wordMap[token]++
			}
		}
		channelRes <- Result{
			NumKeyWords: wordMap,
			Path:        files[i],
			Accuracy:    calculateAcc(wordMap),
		}
	}
}

func calculateAcc(mp map[string]int) int {
	score := 0
	for _, val := range mp {
		score += val * 10
	}
	return score
}

func getSubDirs(dir string) []string {
	var res []string
	subDirs, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, currDir := range subDirs {
		res = append(res, dir+"/"+currDir.Name())
	}
	return res
}

func loadFiles(directory string) []string {
	files, err := os.ReadDir(directory)
	var res []string
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if path.Ext(file.Name()) != ".html" {
			continue
		}
		res = append(res, file.Name())
	}
	return res
}

func readFiles(directory string, fileList []string) []*tree {
	var parsedFiles []*tree

	for _, file := range fileList {
		currTree, body := initTree()

		rFile, err := os.ReadFile(directory + "/" + file)
		check(err)
		fileLines := formatHtml(rFile)
		stripFile(body, fileLines, 0)
		parsedFiles = append(parsedFiles, currTree)
	}
	return parsedFiles
}

func initTree() (*tree, *htmlTreeNode) {
	wTree := newHTMLTree()
	root := newHTMLNode("html", []string{""}, nil)
	wTree.setRoot(root)
	head := newHTMLNode("head", []string{""}, root)
	body := newHTMLNode("body", []string{""}, root)
	root.addChild(head)
	root.addChild(body)

	return wTree, body
}

func formatHtml(file []byte) []string {
	fileLines := strings.Split(string(file), "\r\n")
	fileLines = formatFile(fileLines)
	return fileLines[8:]
}
