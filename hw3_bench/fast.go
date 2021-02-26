package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	seenBrowsers := make([]string, 0, 114)
	uniqueBrowsers := 0
	foundUsers := make([]byte, 0, 4188)

	user := &User{}
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		// fmt.Printf("%v %v\n", err, line)
		err = user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		foundUsers = append(foundUsers, '[')
		foundUsers = append(foundUsers, strconv.Itoa(i)...)
		foundUsers = append(foundUsers, "] "...)
		foundUsers = append(foundUsers, user.Name...)
		foundUsers = append(foundUsers, " <"...)
        for _, ch := range user.Email {
            if ch == '@' {
                foundUsers = append(foundUsers, " [at] "...)
            } else {
                foundUsers = append(foundUsers, string(ch)...)
            }
        }
		foundUsers = append(foundUsers, ">\n"...)
	}

	fmt.Fprintln(out, "found users:\n"+string(foundUsers))
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
