package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/admin/directory/v1"
)

// Documentation:
//
// https://developers.google.com/admin-sdk/directory/
// https://developers.google.com/admin-sdk/directory/v1/quickstart/go
// https://godoc.org/google.golang.org/api/admin/directory/v1

func listDirectory(source ClientSource, runner contactFunc, customer, domain, account, query string, maxAccounts int, maxContacts int) {
	if customer == "" && domain == "" {
		exit(errors.New("A customer or domain must be provided in order to list a directory"))
	}

	client, err := source.Client(context.TODO(), account)
	if err != nil {
		exit(fmt.Errorf("Unable to create authenticate HTTP client: %v", err))
	}
	svc, err := admin.New(client)
	if err != nil {
		exit(fmt.Errorf("Unable to retrieve directory client: %v", err))
	}

	call := svc.Users.List()
	if customer != "" {
		call.Customer(customer)
	}
	if customer != "" {
		call.Domain(domain)
	}
	call.MaxResults(int64(maxAccounts))
	response, err := call.Do()
	if err != nil {
		exit(fmt.Errorf("Unable to retrieve directory listing: %v", err))
	}

	if len(response.Users) == 0 {
		fmt.Print("No users found.\n")
	} else {
		for _, account := range response.Users {
			//fmt.Printf("%s <%s>:\n", account.Name.FullName, account.PrimaryEmail)
			if account.PrimaryEmail != "" {
				runner(source, account.PrimaryEmail, query, maxContacts)
			}
		}
	}
}
