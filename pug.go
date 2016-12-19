package main

import (
	"fmt"
	"os"
	"sync"
	"io/ioutil"
	"strings"
	"github.com/fatih/color"
)

var wg sync.WaitGroup
var search_string string = ""

var cyan = color.New(color.FgCyan)
var boldCyan = cyan.Add(color.Bold)

func walk_r(dir string){
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range list{
		if strings.HasPrefix(v.Name(), ".")==false && strings.HasSuffix(v.Name(), ".log")==false{ //ignore hidden ones
			path := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), v.Name())
			if v.IsDir() {
				wg.Add(1)
				go walk_r(path)
			}else{
				file, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Println(err)
				}
				str := string(file)
				line := 0
				temp := strings.Split(str, "\n")
				matched_lines := []string{}
				for _, item := range temp {
					if strings.Contains(item, search_string){
						t1 := fmt.Sprintf("%d \t %s", line, item)
						t2 := len(t1)
						if len(t1) > 255{
							t2 = 255
						}
						matched_lines = append(matched_lines, t1[:t2])
					}
					line++
				}
				if len(matched_lines) > 0{
					
					boldCyan.Println(path)
					for _, str := range matched_lines{
						color.Yellow(str)
					}
				}
				
			}
		}
	}
	wg.Done()
}

func main() {
	args := os.Args[1:]
	
	dir := args[0]
	search_string = args[1]
	
	wg.Add(1)
	go walk_r(dir)
	wg.Wait()
}
