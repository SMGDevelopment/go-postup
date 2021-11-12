# Go-PostUp
[![Go Reference](https://pkg.go.dev/badge/github.com/cheddartv/go-postup.svg)](https://pkg.go.dev/github.com/cheddartv/go-postup)
[![Go Report Card](https://goreportcard.com/badge/github.com/cheddartv/go-postup)](https://goreportcard.com/report/github.com/cheddartv/go-postup)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client for the [Upland PostUp](https://uplandsoftware.com/postup/) [API](https://apidocs.postup.com/docs).

## How to use
The PostUp API uses basic authentication, this means that you only need to provide a PostUp username and password to the `postup.NewPostUp` method. It is safe to call the client from multiple Goroutines.

### Example
```Go
package main

import (
    "context"
    "fmt"

    "github.com/cheddartv/go-postup"
)

func main() {
    var (
        ctx = context.Background()

        username = "..."
        password = "..."

        client = postup.NewPostUp(username, password, nil)
    )

    lists, err := client.GetLists(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println(lists)
}
```

## Implemented
### Importing Data
- [ ] Member imports through an import template
- [ ] Return a single import status
- [ ] Return import stats

### Import Templates
- [ ] Create an import template
- [ ] Update an import template
- [ ] Return import template configurations

### Recipients
- [x] Create a recipient
- [x] Update existing recipient data
- [x] Return recipient data by email address
- [x] Return recipient data by external id 
- [x] Return recipient data by recipient id 

### Recipient Privacy
- [ ] Return recipient privacy data
- [ ] Delete all recipient privacy data
- [x] Delete recipient from database
- [ ] Delete recipient location data
- [ ] Delete recipient contact data
- [ ] Delete recipient demographics data
- [ ] Delete recipient origin data

### Lists
- [ ] Create a new list
- [ ] Update existing list
- [x] Return all lists
- [x] Return list by brand id
- [ ] Return list counts

### List Subscriptions
- [x] Subscribe recipient to a list
- [x] Unsubscribe recipient from a list
- [ ] Return list subscriptions for a recipient
- [ ] Return member subscriptions for a list
- [ ] Check if recipient is subscribed to a list

### Triggered Mailings
- [ ] Create a send template
- [ ] Update a send template
- [ ] Send a triggered mailing
- [ ] Return send template details

### Scheduled Mailings
- [ ] Create a campaign
- [ ] Create scheduled mailing
- [ ] Create A/B split mailing
- [ ] Send test message from mailing draft
- [ ] Schedule or update mailing draft
- [ ] Check mailing status
- [ ] Return brand information

### Reports
- [ ] Return all campaigns
- [ ] Return campaign reports
- [ ] Return report for a specific mailing
- [ ] Return detailed click report for a mailing
- [ ] Return mailings under a specific campaign
- [ ] Recipient level reporting
- [ ] Return recipients by engagement

### Content
- [ ] Create content folder in CMS
- [ ] Create new content in CMS
- [ ] Update existing content in CMS
- [ ] Return content from a specific directory

### Custom Fields
- [ ] Create a custom field
- [ ] Update a custom field
- [ ] Return custom field data

### Site Monitoring
- [ ] Return site information

### Data Fields
- [ ] List property data fields
- [ ] Recipient endpoint data fields
- [ ] List subscription endpoint data fields
- [ ] Import endpoint data fields
- [ ] Import template endpoint data fields
- [ ] Send template endpoint data fields
- [ ] Template mailing endpoint data fields
- [ ] Mailing endpoint data fields
- [ ] Mailing report data fields
- [ ] Link statistics data fields
- [ ] Campaign statistics data fields
- [ ] Content data fields
- [ ] Content upload data fields
- [ ] Custom field data fields
