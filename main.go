package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	appName = filepath.Base(os.Args[0])
	in      = flag.String("in", "", "input from file instead of STDIN")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", appName)
	fmt.Fprintf(os.Stderr, "Upload document to hastebin.com.")
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func _main() error {
	flag.Usage = usage
	flag.Parse()

	reader := os.Stdin
	if *in != "" {
		f, err := os.Open(*in)
		if err != nil {
			return err
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			return fmt.Errorf("error during stat: %v", err)
		}
		if stat.Size() > 512*1024 {
			return fmt.Errorf("file is larger than 512KiB")
		}

		reader = f
	}

	url := "https://hastebin.com/documents"
	resp, err := http.Post(url, "text/plain", reader)
	if err != nil {
		return fmt.Errorf("unable to post document: %v", err)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Server responded with error %s", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read reply contents: %v", err)
	}

	fmt.Println(string(data))

	return nil
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}
