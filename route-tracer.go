package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Redirect struct {
	Path       string
	StatusCode int
}

var maximumRedirects int

type RequestInterceptor struct {
	Transport http.RoundTripper
}

func (ri RequestInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	t := ri.Transport
	resp, err := t.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	fmt.Println("Request URL :", req.URL.String())
	fmt.Println("Response Code received :", resp.StatusCode)
	fmt.Println(strings.Repeat("-", 40))
	return resp, err
}

func RedirectPolicy(req *http.Request, via []*http.Request) error {
	if len(via) >= maximumRedirects {
		return errors.New("Stopping after reaching the maximum redirect limit")
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [-m <Maximum number of redirects>] (http/https)://(website)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n")
		fmt.Fprintf(os.Stderr, "  -m : Number of redirects after which it should stop (Default : 50)\n")
	}
	flag.IntVar(&maximumRedirects, "m", 50, "Maximum number of redirects")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR :: Website is not specified\n")
		flag.Usage()
		os.Exit(0)
	}

	httpClient := &http.Client{
		Transport:     RequestInterceptor{http.DefaultTransport},
		CheckRedirect: RedirectPolicy,
	}

	website := flag.Args()[0]
	if !strings.HasPrefix(website, "http") {
		website = "http://" + website
	}
	fmt.Println(strings.Repeat("-", 40))
	_, err := httpClient.Get(website)
	if err != nil {
		fmt.Println(err.Error())
	}
}
