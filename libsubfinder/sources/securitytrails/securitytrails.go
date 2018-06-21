//
// Written By : @ice3man (Nizamul Rana)
//
// Distributed Under MIT License
// Copyrights (C) 2018 Ice3man
//

// A golang SecurityTrails API client for subdomain discovery.
package securitytrails

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Ice3man543/subfinder/libsubfinder/helper"
)

type securitytrails_object struct {
	Subdomains []string `json:"subdomains"`
}

var securitytrails_data securitytrails_object

// all subdomains found
var subdomains []string

// Query function returns all subdomains found using the service.
func Query(args ...interface{}) interface{} {

	domain := args[0].(string)
	state := args[1].(*helper.State)

	// We have recieved an API Key
	if state.ConfigState.SecurityTrailsKey != "" {

		// Get credentials for performing HTTP Basic Auth
		securitytrailsKey := state.ConfigState.SecurityTrailsKey

		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.securitytrails.com/v1/domain/"+domain+"/subdomains", nil)

		req.Header.Add("APIKEY", securitytrailsKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("\nsecuritytrails: %v\n", err)
			return subdomains
		}

		// Get the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("\nsecuritytrails: %v\n", err)
			return subdomains
		}

		// Decode the json format
		err = json.Unmarshal([]byte(body), &securitytrails_data)
		if err != nil {
			fmt.Printf("\nsecuritytrails: %v\n", err)
			return subdomains
		}

		// Append each subdomain found to subdomains array
		for _, subdomain := range securitytrails_data.Subdomains {
			finalSubdomain := subdomain + "." + domain

			if state.Verbose == true {
				if state.Color == true {
					fmt.Printf("\n[%sSECURITYTRAILS%s] %s", helper.Red, helper.Reset, finalSubdomain)
				} else {
					fmt.Printf("\n[SECURITYTRAILS] %s", finalSubdomain)
				}
			}

			subdomains = append(subdomains, finalSubdomain)
		}

	}

	return subdomains
}
