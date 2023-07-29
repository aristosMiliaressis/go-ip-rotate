package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/aristosMiliaressis/go-ip-rotate/pkg/iprotate"
	"github.com/projectdiscovery/gologger"
)

type Options struct {
	Url         *url.URL
	Profile     string
	ApiToDelete string
}

func main() {
	opts := parseArguments()

	endpoint, err := iprotate.CreateApi(opts.Profile, opts.Url)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
		os.Exit(1)
	}

	fmt.Printf("ApiId: %s\n", endpoint.ApiId)
	fmt.Printf("DeploymentId: %s\n", endpoint.DeploymentId)
	fmt.Printf("ProxyUrl: %s\n", endpoint.ProxyUrl)
}

func parseArguments() *Options {
	opts := &Options{}
	var baseUrl string
	var err error

	flag.StringVar(&baseUrl, "url", "", "Base Url for apigateway proxy.")
	flag.StringVar(&opts.Profile, "profile", "default", "AWS profile to use.")
	flag.StringVar(&opts.ApiToDelete, "delete", "", "Delete api Id.")
	flag.Parse()

	opts.Url, err = url.Parse(baseUrl)
	if err != nil || baseUrl == "" {
		gologger.Fatal().Msgf("Invalid Url Provided: %s", baseUrl)
		os.Exit(1)
	}

	return opts
}
