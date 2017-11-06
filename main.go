package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type dirInfo struct {
	name  string
	path  string
	files int32
}

type chanDirInfo struct {
	dirs []dirInfo
	cnt  int32
}

func main() {
	c := make(chan chanDirInfo)
	go mapDirectories("d:/Images", c)
	dirinfo := <-c
	fmt.Printf("Main directory contains %d\n\n", dirinfo.cnt)

	for _, v := range dirinfo.dirs {
		fmt.Printf("%-60s %-40s %d\n", v.path, v.name, v.files)
	}
}

func mapDirectories(path string, c chan chanDirInfo) {
	dirmap := []dirInfo{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var filecount int32

	for _, file := range files {
		if file.IsDir() {
			c := make(chan chanDirInfo)
			go mapDirectories(path+"/"+file.Name(), c)
			dirinfo := <-c
			dirmap = append(dirmap, dirInfo{file.Name(), path, dirinfo.cnt})
			dirmap = append(dirmap, dirinfo.dirs...)
		} else {
			filename := strings.ToUpper(filepath.Ext(file.Name()))
			if filename == ".JPG" || filename == ".JPEG" || filename == ".PNG" || filename == ".GIF" || filename == ".RAW" || filename == ".RW2" {
				filecount++
			}
		}
	}
	c <- chanDirInfo{dirmap, filecount}
}
