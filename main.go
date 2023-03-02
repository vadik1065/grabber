package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// выводим ошибку
func shownError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// делает валидное имя
func doValidName(nameString *string) {
	delSymbols := [5]string{"/", "http", ".", "http", ":"}
	for _, delStr := range delSymbols {
		*nameString = strings.ReplaceAll(*nameString, delStr, "")
	}
}

// скачиваем страницу
func downHtm(namePage string) {
	http, err := http.Get(namePage)
	if err == nil {
		doValidName(&namePage)
		body, err := ioutil.ReadAll(http.Body)
		direct := flag.Arg(1)
		formatF := "html"
		err = ioutil.WriteFile(direct+namePage+"."+formatF, body, 0644)

		shownError(err)

	}
	shownError(err)

}

// основная функция
func main() {
	//парсим флаги
	flag.Parse()

	// fmt.Println(flag.Arg(0))
	// fmt.Println(flag.Arg(1))

	// чтение файла
	file, err := os.Open(flag.Arg(0))
	shownError(err)

	fileScaner := bufio.NewScanner(file)
	fileScaner.Split(bufio.ScanLines)

	// пробегаем по всем строкам
	for fileScaner.Scan() {
		puthPage := fileScaner.Text()
		downHtm(puthPage)
	}

	defer file.Close()
}
