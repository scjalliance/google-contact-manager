package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"google.golang.org/api/admin/directory/v1"

	"github.com/scjalliance/google-api/contacts/v3"
)

const (
	apiPeople   = "people"
	apiContacts = "contacts"
)

var (
	scopes = []string{contacts.ContactsReadonlyScope}
)

type contactFunc func(source ClientSource, account, query string, maxContacts int)

func main() {
	var (
		clientID    string
		secret      string
		keyfile     string
		domain      string
		customer    string
		account     string
		query       string
		api         string
		maxAccounts int
		maxContacts int
	)

	flag.StringVar(&keyfile, "keyfile", "", "OAuth v2 key file")
	flag.StringVar(&clientID, "clientid", "", "OAuth v2 client ID")
	flag.StringVar(&secret, "secret", "", "OAuth v2 client secret")
	flag.StringVar(&domain, "domain", "", "Domain to interrogate via contacts API")
	flag.StringVar(&customer, "customer", "", "Customer to interrogate via contacts API")
	flag.StringVar(&account, "account", "", "Email address of account to interrogate")
	flag.StringVar(&api, "api", apiContacts, fmt.Sprintf("API to use (\"%s\" or \"%s\")", apiContacts, apiPeople))
	flag.StringVar(&query, "q", "", "Restricts returned contacts to those matching query (contacts API only)")
	flag.IntVar(&maxAccounts, "accounts", 500, "Specifies the maximum number of accounts to retrieve")
	flag.IntVar(&maxContacts, "max", 500, "Specifies the maximum number of contacts to retreive per account")

	flag.Parse()

	if keyfile == "" && clientID == "" {
		exitWithUsage(errors.New("No keyfile or client ID provided"))
	}

	var (
		source ClientSource
		err    error
		runner contactFunc
	)

	api = parseAPI(api)

	switch api {
	case apiContacts:
		runner = listContacts
	case apiPeople:
		runner = listPeople
	default:
		exitWithUsage(fmt.Errorf("\"%v\" is not a known API", api))
	}

	if customer != "" || domain != "" {
		scopes = append(scopes, admin.AdminDirectoryUserReadonlyScope)
	}

	if keyfile != "" {
		source, err = NewSourceFromKeyfile(keyfile, scopes...)
	} else {
		if secret == "" {
			exitWithUsage(errors.New("No secret provided"))
		}
		source, err = NewSourceFromToken(context.TODO(), clientID, secret, scopes...)
	}
	if err != nil {
		exit(err)
	}

	// TODO: Handle subcommands

	description := queryDescription(api, customer, domain, account, query, maxAccounts, maxContacts)
	description = fmt.Sprintf("%s\n\n", description)

	switch {
	case customer != "" || domain != "":
		fmt.Print(description)
		listDirectory(source, runner, customer, domain, account, query, maxAccounts, maxContacts)
	case account != "":
		fmt.Print(description)
		runner(source, account, query, maxContacts)
	default:
		exitWithUsage(errors.New("No customer, domain or account specified"))
	}
}
