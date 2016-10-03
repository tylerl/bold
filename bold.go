// Make matching lines bold.
// Sort of like grep(1) but all lines get printed, with terminal coloring on the
// lines that match.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	offsetFg = 30
	offsetBg = 40
	escape   = "\x1b"
	reset    = "\x1b[0m"
)

var colors = [...]string{
	"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: bold [-fg <color>] [-bg <color>] <regex>\n\n"+
			"Valid colors include:\n"+
			"(dark|light)-{black,red,green,yellow,blue,magenta,cyan,white}\n"+
			"Use \"none\" to explicitly select no color.\n\n"+
			"Regex can be specified as `-re <regex>` or simply as the last argument.\n"+
			"Ex: \n    ps | bold bash\n    ps | bold -re bash -fg black -bg light-white\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	fg := flag.String("fg", "light-red", "Highlight color")
	bg := flag.String("bg", "", "Background color")
	match := flag.String("re", "", "Regex to match")
	flag.Usage = usage
	flag.Parse()

	if *match == "" {
		s := strings.Join(flag.Args(), " ")
		match = &s
	}
	if *match == "" {
		usage()
		os.Exit(1)
	}
	re := regexp.MustCompile(*match)
	ctrans, err := colorfunc(*fg, *bg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	filterStdio(re, ctrans)
}

func colorfunc(fg, bg string) (func(string) string, error) {
	fgs, err := decodecolor(fg, true, true)
	if err != nil {
		return nil, err
	}
	bgs, err := decodecolor(bg, false, false)
	if err != nil {
		return nil, err
	}
	start := fgs + bgs
	end := ""
	if start != "" {
		end = reset
	}
	rtn := func(s string) string {
		return fmt.Sprintf("%s%s%s", start, s, end)
	}
	return rtn, nil
}

func filterStdio(reg *regexp.Regexp, t func(string) string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if reg.MatchString(line) {
			fmt.Println(t(line))
		} else {
			fmt.Println(line)
		}
	}
}

func decodecolor(s string, bold, fg bool) (string, error) {
	if s == "" || s == "none" {
		return "", nil
	}
	if s == "gray" {
		s = "light-black"
	}
	if s == "bold" {
		s = "light-white"
	}
	if s == "black" {
		s = "dark-black"
	}
	if strings.Contains(s, "-") {
		var ss = strings.SplitN(s, "-", 2)
		p := ss[0]
		s = ss[1]
		if p == "dark" {
			bold = false
		} else if p == "light" || p == "bold" || p == "bright" {
			bold = true
		} else {
			return "", fmt.Errorf("unrecognized prefix %s", p)
		}
	}
	for i, v := range colors {
		if v == s {
			n := i
			if fg {
				n += offsetFg
			} else {
				n += offsetBg
			}
			bolds := ""
			if bold {
				bolds = "1;"
			}
			code := fmt.Sprintf("%s[%s%dm", escape, bolds, n)
			return code, nil
		}
	}
	return "", fmt.Errorf("unrecognized color %s", s)
}
