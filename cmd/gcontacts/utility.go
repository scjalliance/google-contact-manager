package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(2)
}

func exitWithUsage(err error) {
	fmt.Fprintf(os.Stderr, "%v. ", err)
	flag.Usage()
	os.Exit(2)
}

func parseAPI(api string) string {
	api = strings.ToLower(api)

	switch api {
	case "contact":
		api = "contacts"
	case "person":
		api = "people"
	}

	return api
}

func queryDescription(api, customer, domain, account, query string, maxAccounts, maxContacts int) (description string) {
	var (
		isUsingDirectory = customer != "" || domain != ""
		isUsingSearch    = query != "" && api == apiContacts
		queryDesc        string
		scopeDesc        string
		apiDesc          string
	)

	// API description
	switch api {
	case apiContacts:
		//apiDesc = fmt.Sprintf("Google Contacts API v%v", contacts.Version)
		apiDesc = "Google Contacts API"
	case apiPeople:
		apiDesc = "Google People API"
	}

	if isUsingDirectory {
		apiDesc = fmt.Sprintf("%s and Google Admin Directory API", apiDesc)
	}

	// Scope description
	if isUsingSearch {
		scopeDesc = "contacts"
	} else {
		scopeDesc = "all contacts"
	}
	if isUsingDirectory {
		if customer != "" {
			scopeDesc = fmt.Sprintf("%s for customer \"%s\"", scopeDesc, customer)
		}
		if domain != "" {
			scopeDesc = fmt.Sprintf("%s within domain \"%s\"", scopeDesc, domain)
		}
	} else {
		scopeDesc = fmt.Sprintf("%s for account \"%s\"", scopeDesc, account)
	}

	// Query description
	if isUsingSearch {
		queryDesc = fmt.Sprintf("matching \"%s\"", query)
	}

	// Delegation description
	/*
	   if isUsingDirectory && account != "" {
	     scopeDesc = fmt.Sprintf("%s using account \"%s\" for delegated access", scopeDesc, account)
	   }
	*/

	// Final assembly
	switch {
	case isUsingSearch:
		description = fmt.Sprintf("Showing %s %s using %s.", scopeDesc, queryDesc, apiDesc)
	default:
		description = fmt.Sprintf("Showing %s using %s.", scopeDesc, apiDesc)
	}
	return
}
