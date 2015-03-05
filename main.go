package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var s = flag.String("s", "hikari", "speaker")

func say() int {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return 1
	}

	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}

	b, err := ioutil.ReadFile(filepath.Join(home, ".voicetext-apikey"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	apikey := strings.TrimSpace(string(b))

	params := url.Values{}
	params.Set("text", strings.Join(flag.Args(), " "))
	req, err := http.NewRequest("POST", "https://api.voicetext.jp/v1/tts", strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	req.SetBasicAuth(apikey, "")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}

	f, err := ioutil.TempFile(os.TempDir(), "say")
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		os.Exit(1)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}

	err = play(f.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(say())
}
