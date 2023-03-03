package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// выводим ошибку
func shownError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// делает валидное имя
func makeValidName(nameString *string) {
	delSymbols := [4]string{"www.", "https://", "http://", ".html"}
	for _, delStr := range delSymbols {
		*nameString = strings.ReplaceAll(*nameString, delStr, "")
	}
	*nameString = strings.ReplaceAll(*nameString, "/", "-")
	*nameString = strings.TrimSuffix(*nameString, "-")
}

// делает валидную дикерторию
func makeValidDir(directory *string) {
	if len(*directory) != 0 {
		*directory += "/"
	}
}

// основная функция
func main() {

	var wg sync.WaitGroup

	// скачиваем страницу
	downHtm := func(namePage string, direct string) {
		fmt.Println("start " + namePage)
		http, err := http.Get(namePage)
		if err == nil {
			makeValidName(&namePage)
			body, err := ioutil.ReadAll(http.Body)
			formatF := "html"
			err = ioutil.WriteFile(direct+namePage+"."+formatF, body, 0644)
			shownError(err)
		}
		shownError(err)
		fmt.Println("end " + namePage)
		defer wg.Done()
	}

	//парсим флаги
	var directOutput = flag.String("directOutput", "", "sets the directory where to save files")
	var fileInput = flag.String("fileInput", "sites.txt", "path to the file from where to get html page")
	flag.Parse()

	makeValidDir(directOutput)

	// чтение файлa
	file, err := os.Open(*fileInput)
	fileScaner := bufio.NewScanner(file)
	fileScaner.Split(bufio.ScanLines)
	shownError(err)

	// пробегаем по всем строкам
	for fileScaner.Scan() && err == nil {
		wg.Add(1)
		puthPage := fileScaner.Text()
		go downHtm(puthPage, *directOutput)
	}

	wg.Wait()

	defer file.Close()
}
