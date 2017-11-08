package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

var crc32FileMap = map[int]map[uint64][]string{}

var jobtimer struct {
	start  string
	finish string
}

func main() {
	jobtimer.start = time.Now().String()
	c := make(chan chanDirInfo)
	go mapDirectories("d:/Images", c)
	<-c

	f, err := os.Create("d:/Images/dupe-result.txt")
	if err != nil {
		fmt.Println("Error creating result file", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	jobtimer.finish = time.Now().String()

	_, err2 := w.WriteString(fmt.Sprintf("Job start-stop times: %v \r\n", jobtimer))
	if err2 != nil {
		fmt.Println("Error writing to buffer", err2)
	}

	_, err3 := w.WriteString("Found but ignored extensions:\r\n")
	if err3 != nil {
		fmt.Println("Error writing to buffer", err3)
	}

	for k, v := range ignoreExtensionMap {
		_, err := w.WriteString(fmt.Sprintf("%-40s %-10d\r\n", k, v))
		if err != nil {
			fmt.Println("Error writing to buffer", err)
		}
	}

	for _, h := range crc32FileMap {
		for _, v := range h {
			if len(v) > 1 {
				_, err := w.WriteString(fmt.Sprintf("\r\n%s\r\n", v[0]))
				if err != nil {
					fmt.Println("Error writing to buffer", err)
				}
				for _, fl := range v[1:] {
					_, err := w.WriteString(fmt.Sprintf("\t%s\r\n", fl))
					if err != nil {
						fmt.Println("Error writing to buffer", err)
					}
				}
			}
		}
	}

	w.Flush()

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
			//fmt.Printf("\r checked files in %s\t\t\t\t\t", f)
			dirmap = append(dirmap, dirInfo{file.Name(), path, dirinfo.cnt})
			dirmap = append(dirmap, dirinfo.dirs...)
		} else {
			ext := strings.ToUpper(filepath.Ext(f))
			if _, ok := extesionMap[ext]; ok {
				if file.Size() < 10000000 {
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
	/* body, err := ioutil.ReadFile(file) */

	b := make([]byte, 150000)

	f, err := os.Open(file)

	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	_, err2 := f.Read(b)

	if err2 != nil {
		fmt.Printf("Error reading body to limited byte array in %s with error %s\n", file, err2)
		return
	}

	//make a hash based on crc32
	//crc := crc32.Checksum(b, crc32q)

	crc := fingerprint(b)

	f.Close()

	//add to hashmap
	addToMap(crc, file)
}

func addToMap(crc uint64, loc string) {
	first := int(crc % 25)

	if _, ok := crc32FileMap[first]; !ok {
		crc32FileMap[first] = map[uint64][]string{}
	}

	if _, ok := crc32FileMap[first][crc]; !ok {
		crc32FileMap[first][crc] = []string{}
	}

	crc32FileMap[first][crc] = append(crc32FileMap[first][crc], loc)
}

func fingerprint(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}
