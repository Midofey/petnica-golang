package main

import "fmt"

func capitaliseWords(s string) (string, int) {
	temp := ""
	brojac := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			if i == 0 || s[i-1] == ' ' {
				c -= 32
				brojac++
			}
		}
		temp += string(c)
	}
	return temp, brojac
}

func main() {
	var unos string = ""
	fmt.Print("Unesite string: ")
	fmt.Scan(&unos)
	capitalised, _ := capitaliseWords("hello world 2")
	fmt.Printf("%s", capitalised)
}
