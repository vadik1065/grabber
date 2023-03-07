package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

// makeValidName - делает валидное имя
func makeValidName(nameString string) string {
	delSymbols := [4]string{"www.", "https://", "http://", ".html"}
	for _, delStr := range delSymbols {
		nameString = strings.ReplaceAll(nameString, delStr, "")
	}
	nameString = strings.ReplaceAll(nameString, "/", "-")
	nameString = strings.TrimSuffix(nameString, "-")
	return nameString
}

//  downloadHtml - скачиваем страницу
func downloadHtml(namePage string, directory string) {
	fmt.Printf("start %s \n", namePage)
	http, err := http.Get(namePage)
	defer wg.Done()

	if err != nil {
		fmt.Println(err)
		return
	}
	defer http.Body.Close()

	namePage = makeValidName(namePage)
	body, err := ioutil.ReadAll(http.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fileFormat := "html"
	nameComponents := []string{namePage, ".", fileFormat}
	nameFile := strings.Join(nameComponents, "")
	fullNameFile := path.Join(directory, nameFile)

	err = ioutil.WriteFile(fullNameFile, body, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("end %s \n", namePage)
}

var wg sync.WaitGroup

// main - основная функция
func main() {

	//парсим флаги
	var directOutput = flag.String("directOutput", "", "sets the directory where to save files")
	var fileInput = flag.String("fileInput", "sites.txt", "path to the file from where to get html page")
	flag.Parse()

	// чтение файлa
	file, err := os.Open(*fileInput)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileScaner := bufio.NewScanner(file)
	fileScaner.Split(bufio.ScanLines)

	// пробегаем по всем строкам файлы и скачиваем
	for fileScaner.Scan() && err == nil {
		wg.Add(1)
		puthPage := fileScaner.Text()
		go downloadHtml(puthPage, *directOutput)
	}

	wg.Wait()
}
