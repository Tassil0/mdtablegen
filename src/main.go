package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	output     string
	outputFile *os.File
	err        error
	urlPrefix  string
)

type Book struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Book   string `json:"book"`
	Url    string `json:"url"`
}

type BookTable struct {
	Id       string
	Author   string
	BookLink string
}

func main() {
	output = "./output.md"
	urlPrefix = "https://github.com/POJFM/cetba/tree/main/"
	path, _ := os.Getwd()
	header := []string{"**ID**", "**Autor**", "**Dílo**"}

	if _, err = os.Stat(output); errors.Is(err, os.ErrNotExist) {
		if outputFile, err = os.Create(output); err != nil {
			log.Fatal(err)
		}
		if err = outputFile.Close(); err != nil {
			log.Fatal(err)
		}
	}

	outputFile, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	dirList := getList(path)

	//books := parseRawBooks("./data.json")
	//
	//if err = writeJSON("./data.json", books); err != nil {
	//	log.Fatal(err)
	//}

	books, err := readJSON("./data.json")
	if err != nil {
		log.Fatal(err)
	}

	books = addURL(books, urlPrefix, dirList)

	bookTable := makeBookTable(books)
	//for _, b := range bookTable {
	//	fmt.Println(b.Author)
	//}
	colWidths := []int{6, getLongestRowAuthor(bookTable), getLongestRowBookLink(bookTable)}

	writeHeader(header, colWidths)

	writeRows(bookTable, colWidths)

	//for _, e := range dirList {
	//	_, err := outputFile.WriteString(urlPrefix + e + "\n")
	//	//log.Println("written: ", n)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	//fmt.Println(e)
	//}

	_ = outputFile.Close()
}

// fill end of string so the table lines up
func colFill(row string, colWidth int) string {
	if len(row) != colWidth {
		for o := true; o; o = len(row) < colWidth {
			row += " "
		}
	}
	return row
}

func writeCol(row string, colWidth int) {
	if _, err = outputFile.WriteString(" " + colFill(row, colWidth) + " "); err != nil {
		log.Fatal(err)
	}
}

func writeRows(rows []BookTable, colWidths []int) {
	for _, row := range rows {
		writeRow(row, colWidths)
	}
}

func writeRow(row BookTable, colWidths []int) {
	for i, colWidth := range colWidths {
		writeDivider()
		switch i {
		case 0:
			writeCol(row.Id, colWidth)
		case 1:
			writeCol(row.Author, colWidth)
		case 2:
			writeCol(row.BookLink, colWidth)
		}
	}
	writeDivider()
	writeNewline()
}

func writeHeader(header []string, colWidths []int) {
	for i, row := range header {
		writeDivider()
		writeCol(row, colWidths[i])
	}
	writeDivider()
	writeNewline()
	writeHr(colWidths)
}

func writeDivider() {
	if _, err = outputFile.WriteString("|"); err != nil {
		log.Fatal(err)
	}
}

func writeNewline() {
	if _, err = outputFile.WriteString("\n"); err != nil {
		log.Fatal(err)
	}
}

func writeHr(colWidths []int) {
	for _, colWidth := range colWidths {
		writeDivider()
		var col string
		for o := true; o; o = len(col) < colWidth {
			col += "-"
		}
		if _, err = outputFile.WriteString(" " + col + " "); err != nil {
			log.Fatal(err)
		}
	}
	writeDivider()
	writeNewline()
}

func makeBookTable(books []Book) (bookTable []BookTable) {
	var appendBook BookTable
	for _, book := range books {
		appendBook.Id = fmt.Sprintf("%02d", book.Id)
		appendBook.Author = book.Author
		if book.Url != "" {
			appendBook.BookLink = fmt.Sprintf("[%s](%s)", book.Book, book.Url)
		} else {
			appendBook.BookLink = book.Book
		}
		bookTable = append(bookTable, appendBook)
	}
	return
}

func getLongestRowAuthor(rows []BookTable) (n int) {
	index := 0
	for i, row := range rows {
		if len(rows[index].Author) < len(row.Author) {
			index = i
		}
	}
	return len(rows[index].Author)
}

func getLongestRowBookLink(rows []BookTable) (n int) {
	index := 0
	for i, row := range rows {
		if len(rows[index].BookLink) < len(row.BookLink) {
			index = i
		}
	}
	return len(rows[index].BookLink)
}

func getList(path string) []string {
	var dirs []string
	dirList, _ := os.ReadDir(path)
	for _, dir := range dirList {
		if dir.IsDir() && (dir.Name() == ".git" || dir.Name() == ".idea") {
			continue
		}
		if dir.IsDir() {
			dirs = append(dirs, dir.Name())
		}
	}
	return dirs
}

func addURL(books []Book, prefix string, dirs []string) []Book {
	for i, book := range books {
		for _, dir := range dirs {
			if dirId, err := strconv.Atoi(dir[0:2]); err == nil {
				if dirId == book.Id {
					books[i].Url = prefix + dir
				}
			} else {
				log.Fatal(err)
			}
		}
	}
	return books
}

func writeJSON(path string, books []Book) error {
	booksJson, _ := json.MarshalIndent(books, "", " ")
	if err := os.WriteFile(path, booksJson, 0644); err != nil {
		return err
	}
	return nil
}

func readJSON(path string) (books []Book, err error) {
	if data, err := os.ReadFile(path); err == nil {
		if err = json.Unmarshal(data, &books); err != nil {
			return []Book{}, err
		}
	} else {
		return []Book{}, err
	}
	return books, nil
}

func parseRawBooks(path string) (books []Book) {
	if _, err = os.Stat(output); errors.Is(err, os.ErrNotExist) {
		log.Fatal("data.json with raw books data not found")
		return []Book{}
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal("could not open data.json")
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		data := strings.Split(line, "\t")
		// trim trailing whitespace
		for i, d := range data {
			if i != 2 {
				data[i] = trimTrailingWhitespace(d)
			}
		}
		id, err := strconv.Atoi(data[0])
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, Book{
			Id:     id,
			Author: data[1],
			Book:   data[2],
			Url:    "",
		})
	}
	_ = f.Close()
	return
}

func trimTrailingWhitespace(s string) string {
	return string(s[0 : len(s)-1])
}
