package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	var urlFlag string
	var cookie string
	var customHeader string
	var proxyFlag string
	var urlOnly bool
	var pathOnly bool
	var silentMode bool
	var client *http.Client

	flag.StringVar(&urlFlag, "t", "", "Specify the target URL (e.g., domain.com or https://domain.com)")
	flag.StringVar(&cookie, "c", "", "Specify cookies (e.g., user_token=g3p21ip21h; )")
	flag.StringVar(&customHeader, "r", "", "Specify headers (e.g., Myheader: test )")
	flag.StringVar(&proxyFlag, "p", "", "Specify the proxy URL (e.g., 127.0.0.1:8080)")
	flag.BoolVar(&urlOnly, "uO", false, "Extract only URLs")
	flag.BoolVar(&pathOnly, "pO", false, "Extract only paths")
	flag.BoolVar(&silentMode, "s", false, "Avoid printing banner and other messages")

	helpFlag := flag.Bool("h", false, "Display help")
	flag.Parse()

	if *helpFlag {
		fmt.Println("gourlex is a tool for extracting URLs from a webpage.")
		fmt.Println("\nUsage:")
		fmt.Printf("  %s [arguments]\n", os.Args[0])
		fmt.Println("\nThe arguments are:")
		fmt.Println("  -t string    Specify the target URL (e.g., domain.com or https://domain.com)")
		fmt.Println("  -c string    Specify cookies (e.g., user_token=g3p21ip21h; )")
		fmt.Println("  -r string    Specify headers (e.g., Myheader: test )")
		fmt.Println("  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)")
		fmt.Println("  -s           Silent Mode, avoid printing banner and other messages")
		fmt.Println("  -uO          Extract only full URLs")
		fmt.Println("  -pO          Extract only URL paths")
		fmt.Println("  -h           Display help")

		fmt.Println("\nExample:")
		fmt.Println("  gourlex -t domain.com ")

		return
	}

	if !silentMode {
		printBanner()
	}

	if urlFlag == "" {
		fmt.Println("Please provide a target.\n Example usage: gourlex -t domain.com\n Use -h for help.")
		return
	}

	validUrl, err := validateUrl(urlFlag)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	if urlOnly && pathOnly {
		fmt.Println("You can't use both -uO and -pO flags together")
		return
	}

	if len(proxyFlag) > 0 {
		if isValidProxy(proxyFlag) {
			if !silentMode {
				fmt.Printf("Using proxy: %s\n\n", proxyFlag)
			}
			client, err = createHTTPClientWProxy(proxyFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("Invalid proxy:", proxyFlag)
			fmt.Println("Please insert a valid proxy in the ip:port format")
			return
		}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", validUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	if customHeader != "" {
		headerParts := strings.SplitN(customHeader, ":", 2)
		if len(headerParts) == 2 {
			req.Header.Add(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
		} else {
			fmt.Printf("Invalid header format: %s\n", customHeader)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if !silentMode {
		fmt.Printf("\033[1;97mExtracting URLs from: %s\033[0m\n\n", urlFlag)
	}

	urls, paths, err := extractURLsAndPathsFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	printURLsAndPaths(urls, paths, urlOnly, pathOnly, silentMode)
}

func extractURLsAndPathsFromResponse(resp *http.Response) ([]string, []string, error) {
	tokenizer := html.NewTokenizer(resp.Body)
	var urls []string
	var paths []string

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return urls, paths, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			var hrefValue string
			var srcValue string
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					hrefValue = strings.TrimSpace(attr.Val)
				} else if attr.Key == "src" {
					srcValue = strings.TrimSpace(attr.Val)
				}
			}

			if hrefValue != "" && hrefValue != "#" {
				if strings.HasPrefix(hrefValue, "http://") || strings.HasPrefix(hrefValue, "https://") {
					urls = append(urls, hrefValue)
				} else if u, err := url.Parse(hrefValue); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
					urls = append(urls, u.String())
				} else {
					paths = append(paths, hrefValue)
				}
			}

			if srcValue != "" && hrefValue != "#" {
				if strings.HasPrefix(srcValue, "http://") || strings.HasPrefix(srcValue, "https://") {
					urls = append(urls, srcValue)
				} else if u, err := url.Parse(srcValue); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
					urls = append(urls, u.String())
				} else {
					paths = append(paths, srcValue)
				}
			}
		}
	}
}

func printURLsAndPaths(urls []string, paths []string, urlOnly, pathOnly, silentMode bool) {

	if !pathOnly {
		if !silentMode {
			str := "Extracted URLs:\n\n"
			coloredUrls := colorize(str, "\033[1;32m")
			fmt.Print(coloredUrls)
		}

		for _, u := range urls {
			fmt.Println(u)
		}
	}

	fmt.Println()

	if !urlOnly {
		if !silentMode {
			str2 := "Extracted Paths:\n\n"
			coloredPaths := colorize(str2, "\033[1;32m")
			fmt.Print(coloredPaths)
		}
		for _, p := range paths {
			fmt.Println(p)
		}
	}

}

func validateUrl(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %v", err)
	}

	if u.Scheme == "" {
		inputURL = "https://" + inputURL
		u, _ = url.Parse(inputURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("Invalid URL scheme")
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	_, err = net.LookupHost(host)
	if err != nil {
		return "", err
	}

	if port != "" {
		inputURL = fmt.Sprintf("%s://%s:%s%s", u.Scheme, host, port, u.RequestURI())
	} else {
		inputURL = fmt.Sprintf("%s://%s%s", u.Scheme, host, u.RequestURI())
	}

	return inputURL, nil
}

func createHTTPClientWProxy(proxy string) (*http.Client, error) {
	parts := strings.Split(proxy, ":")
	proxyIP := parts[0]
	proxyPortStr := parts[1]
	proxyPort, err := strconv.Atoi(proxyPortStr)
	if err != nil {
		return nil, fmt.Errorf("error converting proxy port to integer: %v", err)
	}

	client := &http.Client{}
	if proxyIP != "" && proxyPort != 0 {
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", proxyIP, proxyPort))
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy URL: %v", err)
		}
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client, nil
}

func isValidProxy(input string) bool {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return false
	}

	ip := parts[0]
	portStr := parts[1]

	if net.ParseIP(ip) == nil {
		return false
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return false
	}
	return true
}

func colorize(text string, colorCode string) string {
	resetColor := "\033[0m"
	return colorCode + text + resetColor
}

func printBanner() {
	purple := "\033[1;35m"

	fmt.Println(colorize("                         _           ", purple))
	fmt.Println(colorize("   __ _  ___  _   _ _ __| | _____  __", purple))
	fmt.Println(colorize("  / _` |/ _ \\| | | | '__| |/ _ \\ \\/ /", purple))
	fmt.Println(colorize(" | (_| | (_) | |_| | |  | |  __/>  < ", purple))
	fmt.Println(colorize("  \\__, |\\___/ \\__,_|_|  |_|\\___/_/\\_\\", purple))
	fmt.Println(colorize("  |___/                                  ", purple))

	fmt.Println("")
	fmt.Print(colorize("Gourlex - WebPage Urls Extractor Tool\n\n", purple))
}
