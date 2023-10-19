package internal

import (
	"log"
	"os"
	"path"
	"strings"
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

func Search(query string) []Result {
	var Res []Result
	keyWords := extractKeyWords(query)
	files := loadFiles("./static")
	htmlFiles := readFiles("./static", files)
	wordMap := initKWMap(keyWords)
	for i := 0; i < len(htmlFiles); i++ {
		tokens := tokenize(htmlFiles[i].root, 0)
		for _, token := range tokens {
			if token == keyWords[0] || token == keyWords[1] {
				wordMap[token]++
			}
		}
		Res = append(Res, Result{
			NumKeyWords: wordMap,
			Path:        files[i],
			Accuracy:    calculateAcc(wordMap),
		})
	}
	return Res
}

func calculateAcc(mp map[string]int) int {
	score := 0
	for _, val := range mp {
		score += val * 10
	}
	return score
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
	root := newHTMLNode("html", "", nil)
	wTree.setRoot(root)
	head := newHTMLNode("head", "", root)
	body := newHTMLNode("body", "", root)
	root.addChild(head)
	root.addChild(body)

	return wTree, body
}

func formatHtml(file []byte) []string {
	fileLines := strings.Split(string(file), "\r\n")
	fileLines = formatFile(fileLines)
	return fileLines[8:]
}
