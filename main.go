package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
	tld "github.com/jpillora/go-tld"
)

// Globals
const regexStr = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`

var founds []string
var match_baseurl bool
var wg sync.WaitGroup

//Main Function
func main() {

	// Local Vars
	var dl string
	flag.StringVar(&dl, "dl", "", "Input list of js file urls.")
	flag.BoolVar(&match_baseurl, "m", false, "Only find URLs which have same domain as main domain.")
	var n int
	flag.IntVar(&n, "n", 1, "Number of GoRoutine to use at once.")
	var fname string
	flag.StringVar(&fname, "o", "None", "Output file to write.")
	count := 1

	// Usage Error
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			err_msg := f.Name + ", is not set use -h for help."
			color.Red(err_msg)
			os.Exit(0)
		}
	})

	color.HiGreen("\n==== Current Configuration ====\n")
	fmt.Println("Target file: ", dl)
	fmt.Println("Concurrency: ", n)
	fmt.Println("Matching base domain? ", match_baseurl)
	fmt.Println("Output file: ", fname)
	color.HiGreen("===============================")
	fmt.Println()
	color.HiGreen("[+] Opening the input file")

	// File operation
	file, err := os.Open(string(dl))
	if err != nil {
		log.Fatal("[-] ", err)
		os.Exit(0)
	}
	color.HiGreen("[+] File Open Completed")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal("[-]", err)
	}
	color.HiGreen("[+] File Read Completed")
	for scanner.Scan() {
		if count > n {
			wg.Wait()
			count = 0
		}
		wg.Add(1)
		func() {
			go getFile(scanner.Text())
		}()
		count++
	}
	wg.Wait()
	if fname == "None" {
		printFounds()
	} else {
		writeToFile(fname)
	}
	fmt.Println()
	color.HiGreen("[+] Completed")
}

func getFile(fileurl string) {
	defer wg.Done()
	u, _ := tld.Parse(fileurl)
	baseDomain := u.Domain
	res, err := http.Get(fileurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	matchAndAdd(string(body), baseDomain)
}

func matchAndAdd(content string, baseDomain string) []string {

	regExp, err := regexp.Compile(regexStr)
	if err != nil {
		log.Fatal(err)
	}

	links := regExp.FindAllString(content, -1)
	linksLength := len(links)
	if linksLength > 1 {
		for i := 0; i < linksLength; i++ {
			if match_baseurl == true {
				if strings.Contains(links[i], baseDomain) {
					founds = append(founds, links[i])
				}
			} else {
				founds = append(founds, links[i])
			}
		}

	}
	return founds
}

func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func printFounds() {
	founds = Unique(founds)
	fmt.Println()
	color.HiGreen("[+] %d Results Found: \n", len(founds))
	for i := range founds {
		fmt.Println(founds[i])
	}
}

func writeToFile(fname string) {
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range founds {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}
