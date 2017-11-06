package main

import (
	"fmt"
	"hash/crc32"
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

var extesionMap = map[string]string{
	".JPG":  ".JPG",
	".JPEG": ".JPEG",
	/* 	".PNG":     ".PNG",
	   	".GIF":     ".GIF",
	   	".RAW":     ".RAW",
	   	".RW2":     ".RW2",
	   	".TIF":     ".TIF",
	   	".AI":      ".AI",
	   	".MOV":     ".MOV",
	   	".RAR":     ".RAR",
	   	".MPG":     ".MPG",
	   	".AFPHOTO": ".AFPHOTO",
	   	".AVI":     ".AVI",
	   	".M4V":     ".M4V",
	   	".WMV":     ".WMV",
	   	".MP4":     ".MP4",
	   	".MTS":     ".MTS",
	   	".CDR":     ".CDR",
	   	".BMP":     ".BMP",
	   	".NEF":     ".NEF",
	   	".ZIP":     ".ZIP",
	   	".3GP":     ".3GP",
	   	".CR2":     ".CR2",
	   	".PSD":     ".PSD", */
}

var ignoreExtensionMap = map[string]int{}

var crc32q = crc32.MakeTable(0xD5828281)

var crc32FileMap = map[rune]map[string][]string{}

func main() {
	c := make(chan chanDirInfo)
	go mapDirectories("c:/Users/jarom/Pictures", c)
	<-c

	for _, h := range crc32FileMap {
		for _, v := range h {
			if len(v) > 1 {
				fmt.Println(v)
			}
		}
	}

	/* 	for k, v := range ignoreExtensionMap {
		fmt.Printf("%-40s %-10d\n", k, v)
	} */
}

func mapDirectories(path string, c chan chanDirInfo) {
	dirmap := []dirInfo{}
	var filecount int32

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f := path + "/" + file.Name()
		if file.IsDir() {
			c := make(chan chanDirInfo)
			go mapDirectories(f, c)
			dirinfo := <-c
			fmt.Printf("\r checked files in %s\t\t\t\t\t", f)
			dirmap = append(dirmap, dirInfo{file.Name(), path, dirinfo.cnt})
			dirmap = append(dirmap, dirinfo.dirs...)
		} else {
			ext := strings.ToUpper(filepath.Ext(f))
			if _, ok := extesionMap[ext]; ok {
				if file.Size() < 2000000 {
					filecount++
					handleFile(f)
				}
			} else {
				//collect all not checked file extension
				ignoreExtensionMap[ext]++
			}
		}
	}
	c <- chanDirInfo{dirmap, filecount}
}

func handleFile(file string) {

	//Read complete file
	body, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	//make a hash based on crc32
	crc := string(crc32.Checksum(body, crc32q))

	//add to hashmap
	addToMap(crc, file)
}

func addToMap(crc string, loc string) {
	first := rune(string(crc)[0])

	if _, ok := crc32FileMap[first]; !ok {
		crc32FileMap[first] = map[string][]string{}
	}

	if _, ok := crc32FileMap[first][crc]; !ok {
		crc32FileMap[first][crc] = []string{}
	}

	crc32FileMap[first][crc] = append(crc32FileMap[first][crc], loc)
}
