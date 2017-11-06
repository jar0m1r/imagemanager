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

var crc32FileMap = map[uint32][]string{}

func main() {
	c := make(chan chanDirInfo)
	go mapDirectories("c:/Users/jarom/Pictures", c)
	<-c

	for _, v := range crc32FileMap {
		if len(v) > 1 {
			fmt.Println(v)
		}
	}

	/* 	for k, v := range ignoreExtensionMap {
		fmt.Printf("%-40s %-10d\n", k, v)
	} */
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
			fmt.Printf("\r checked files in %s\t\t\t\t\t", f)
			dirmap = append(dirmap, dirInfo{file.Name(), path, dirinfo.cnt})
			dirmap = append(dirmap, dirinfo.dirs...)
		} else {
			fileext := strings.ToUpper(filepath.Ext(file.Name()))
			if _, ok := extesionMap[fileext]; ok {
				filecount++
				body, err := ioutil.ReadFile(f)

				if err != nil {
					fmt.Println("Error reading file", err)
					continue
				}
				checksum := crc32.Checksum(body, crc32q)
				crc32FileMap[checksum] = append(crc32FileMap[checksum], f)
			} else {
				//collect all not checked file extension
				ignoreExtensionMap[fileext]++
			}
		}
	}
	c <- chanDirInfo{dirmap, filecount}
}
