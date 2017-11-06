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
	".JPG":     ".JPG",
	".JPEG":    ".JPEG",
	".PNG":     ".PNG",
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
	".PSD":     ".PSD",
}

var ignoreExtensionMap = map[string]int{}

var crc32q = crc32.MakeTable(0xD5828281)

func main() {
	c := make(chan chanDirInfo)
	go mapDirectories("d:/Images", c)
	dirinfo := <-c
	fmt.Printf("Main directory contains %d\n\n", dirinfo.cnt)

	/* 	for _, v := range dirinfo.dirs {
		fmt.Printf("%-60s %-40s %d\n", v.path, v.name, v.files)
	} */

	for k, v := range ignoreExtensionMap {
		fmt.Printf("%-40s %-10d\n", k, v)
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
		f := path + "/" + file.Name()
		if file.IsDir() {
			c := make(chan chanDirInfo)
			go mapDirectories(f, c)
			dirinfo := <-c
			dirmap = append(dirmap, dirInfo{file.Name(), path, dirinfo.cnt})
			dirmap = append(dirmap, dirinfo.dirs...)
		} else {
			filename := strings.ToUpper(filepath.Ext(file.Name()))
			if _, ok := extesionMap[filename]; ok {
				filecount++
				fcontent, err := ioutil.ReadFile(f)

				if err != nil {
					fmt.Println("Error reading file", err)
					continue
				}
				fmt.Printf("%08x\n", crc32.Checksum(fcontent, crc32q))
			} else {
				//collect all not checked file extension
				ignoreExtensionMap[filename]++
			}
		}
	}
	c <- chanDirInfo{dirmap, filecount}
}
