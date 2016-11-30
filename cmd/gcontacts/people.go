package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/people/v1"
)

const (
	fields     = "connections(emailAddresses,names,phoneNumbers,resourceName)"
	maskFields = "person.emailAddresses,person.names,person.phoneNumbers"
)

func listPeople(source ClientSource, account, query string, maxContacts int) {
	if account == "" {
		exitWithUsage(errors.New("No account specified"))
	}

	fmt.Printf("%v:", account)

	client, err := source.Client(context.TODO(), account)
	if err != nil {
		exit(err)
	}
	svc, err := people.New(client)
	if err != nil {
		exit(err)
	}

	call := svc.People.Connections.List("people/me")
	call.Fields(fields)
	call.RequestMaskIncludeField(maskFields)
	call.PageSize(int64(maxContacts))
	response, err := call.Do()
	if err != nil {
		exit(err)
	}

	// FIXME: Handle paginated responses?

	fmt.Printf(" %v contact(s)\n", len(response.Connections))

	//log.Printf("Response: %v", response.ServerResponse)

	for i, person := range response.Connections {
		var (
			name           = fmtPeopleNames(person.Names)
			emailAddresses = fmtPeopleEmailAddresses(person.EmailAddresses)
			phoneNumbers   = fmtPeoplePhoneNumbers(person.PhoneNumbers)
		)
		//fmt.Printf("Contact [%v:%s]: \"%v\"%v%v\n", i, person.ResourceName, name, emailAddresses, phoneNumbers)
		fmt.Printf("  %3d: \"%v\"%v%v\n", i, name, emailAddresses, phoneNumbers)
	}
}

func fmtPeopleNames(names []*people.Name) (formatted string) {
	if len(names) == 0 {
		return ""
	}
	return names[0].DisplayName
}

func fmtPeopleEmailAddresses(addresses []*people.EmailAddress) (formatted string) {
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

func fmtPeoplePhoneNumbers(numbers []*people.PhoneNumber) (formatted string) {
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
