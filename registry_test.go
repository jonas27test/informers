package main

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/net/publicsuffix"
)

// go test -run TestRegistry
func TestRegistry(t *testing.T) {
	before("TestRegistry")

	domains := []string{
		"amazon.co.uk",
		"books.amazon.co.uk",
		"www.books.amazon.co.uk",
		"amazon.com",
		"",
		"example0.debian.net",
		"example1.debian.org",
		"",
		"golang.dev",
		"golang.net",
		"play.golang.org",
		"gophers.in.space.museum",
		"",
		"0emm.com",
		"a.0emm.com",
		"b.c.d.0emm.com",
		"",
		"there.is.no.such-tld",
		"",
		// Examples from the PublicSuffix function's documentation.
		"foo.org",
		"foo.co.uk",
		"foo.dyndns.org",
		"foo.blogspot.co.uk",
		"cromulent",
	}

	for _, domain := range domains {
		if domain == "" {
			fmt.Println(">")
			continue
		}
		eTLD, icann := publicsuffix.PublicSuffix(domain)
		eff, err := publicsuffix.EffectiveTLDPlusOne(domain)
		log.Println(eff)
		log.Println(eTLD)
		log.Println(icann)
		if err != nil {
			log.Println(err)
		}
		return
	}

}
