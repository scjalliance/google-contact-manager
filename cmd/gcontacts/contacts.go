package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/scjalliance/google-api/contacts/v3"
)

func listContacts(source ClientSource, account, query string, maxContacts int) {
	if account == "" {
		exitWithUsage(errors.New("No account specified"))
	}

	fmt.Printf("%v:", account)

	client, err := source.Client(context.TODO(), account)
	if err != nil {
		exit(err)
	}
	svc, err := contacts.New(client)
	if err != nil {
		exit(err)
	}

	call := svc.Contacts.Feed(account, "full")
	call.MaxResults(maxContacts)
	if query != "" {
		call.Query(query)
	}
	response, err := call.Do()
	if err != nil {
		exit(err)
	}
	//fmt.Printf("Parsed response: %#v", response)

	fmt.Printf(" %v contact(s)\n", len(response.Feed.Contacts))

	for i, contact := range response.Feed.Contacts {
		var (
			title          = contact.Title
			emailAddresses = fmtContactEmailAddresses(contact.EmailAddresses)
			phoneNumbers   = fmtContactPhoneNumbers(contact.PhoneNumbers)
		)
		fmt.Printf("  %3d: \"%v\"%v%v\n", i, title, emailAddresses, phoneNumbers)
	}
}

func fmtContactEmailAddresses(addresses []*contacts.EmailAddress) (formatted string) {
	if len(addresses) == 0 {
		return
	}
	formatted = " "
	for i, address := range addresses {
		if i > 0 {
			formatted += ", "
		}
		formatted += "<" + address.Value + ">"
	}
	return
}

func fmtContactPhoneNumbers(numbers []*contacts.PhoneNumber) (formatted string) {
	if len(numbers) == 0 {
		return
	}
	formatted = " [ "
	for i, number := range numbers {
		if i > 0 {
			formatted += ", "
		}
		formatted += number.Value
	}
	formatted += " ]"
	return
}
