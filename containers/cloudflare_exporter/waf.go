package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type FwEvent struct {
	Count      int          `json:"count"`
	Dimensions FwDimensions `json:"dimensions"`
}

type FwDimensions struct {
	Action  string `json:"action"`
	ASName  string `json:"clientASNDescription"`
	Country string `json:"clientCountryName"`
	RuleID  string `json:"ruleId"`
}

func buildWAFGraphQLQuery(startDate string, endDate string, zoneID string) *graphql.Request {

	query := graphql.NewRequest(`
	{ 
		viewer {
		zones( filter: { zoneTag: $zoneTag } ) {
		  fwEvents: firewallEventsAdaptiveGroups(
			  limit: 5000, 
			  filter: {datetimeMinute_geq: $startDate, datetimeMinute_leq: $endDate}
		  ) {
			count
			dimensions {
			  action
			  clientCountryName
			  clientASNDescription
			  ruleId
			}
		  }  
		}
		}
	}
  `)

	// set any variables
	query.Var("zoneTag", zoneID)
	query.Var("startDate", startDate)
	query.Var("endDate", endDate)

	return query
}

// Get cloudflare metrics from GraphQL using the provided api-email and api-key parameters and returns a marshalled JSON struct and an error if something went wrong during the fetch
func getCloudflareWAFMetrics(query *graphql.Request, apiEmail string, apiKey string) (respData httpRespDataStruct, err error) {
	client := graphql.NewClient("https://api.cloudflare.com/client/v4/graphql")

	req := query

	// set header fields -> token and email!
	req.Header.Set("x-auth-key", apiKey)
	req.Header.Set("x-auth-email", apiEmail)
	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response

	if err := client.Run(ctx, req, &respData); err != nil {
		return respData, err
	}
	return respData, nil
}
