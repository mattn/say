package main

import (
	"encoding/json"
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
	"time"

	"github.com/mattn/go-soundplayer"
)

var (
	speaker       = flag.String("speaker", "hikari", "show/haruka/hikari/takeru/santa/bear")
	emotion       = flag.String("emotion", "", "happiness/anger/sadness")
	emotion_level = flag.Int("emotion_level", 1, "1/2")
	pitch         = flag.Int("pitch", 100, "50% - 200%")
	speed         = flag.Int("speed", 100, "50% - 400%")
	volume        = flag.Int("volume", 100, "50% - 200%")
	keep          = flag.Bool("keep", false, "remain wav file")
)

type response struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

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
	params.Set("speaker", *speaker)
	if *emotion != "" {
		params.Set("emotion", *emotion)
		params.Set("emotion_level", fmt.Sprint(*emotion_level))
	}
	params.Set("pitch", fmt.Sprint(*pitch))
	params.Set("speed", fmt.Sprint(*speed))
	params.Set("volume", fmt.Sprint(*volume))
	req, err := http.NewRequest("POST", "https://api.voicetext.jp/v1/tts", strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	req.SetBasicAuth(apikey, "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		var resp response
		if err = json.NewDecoder(res.Body).Decode(&resp); err == nil {
			fmt.Fprintln(os.Stderr, "say", resp.Error.Message)
		} else {
			fmt.Fprintln(os.Stderr, "say", "something wrong")
		}
		return 1
	}

	dir, err := ioutil.TempDir(os.TempDir(), "say")
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	defer os.RemoveAll(dir)

	f, err := os.Create(filepath.Join(dir, "say.wav"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}
	f.Close()

	err = soundplayer.Play(f.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, "say", err)
		return 1
	}

	if *keep {
		now := time.Now().Format("say20060102030405.wav")
		copyFile(f.Name(), filepath.Base(now))
	}

	return 0
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	err = out.Close()
	if err != nil {
		return err
	}
	return err
}

func main() {
	os.Exit(say())
}
