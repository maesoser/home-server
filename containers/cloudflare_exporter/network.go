package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type networkRespDataStruct struct {
	NetworkViewer NetworkViewer `json:"networkViewer"`
}
type NetworkDimensions struct {
	AttackID             string `json:"attackId"`
	AttackMitigationType string `json:"attackMitigationType"`
	AttackProtocol       string `json:"attackProtocol"`
	AttackType           string `json:"attackType"`
	ColoCountry          string `json:"coloCountry"`
	DestinationPort      int    `json:"destinationPort"`
}
type Sum struct {
	Bits    int `json:"bits"`
	Packets int `json:"packets"`
}
type AttackHistory struct {
	NetworkDimensions NetworkDimensions `json:"networkDimensions"`
	Sum               Sum               `json:"sum"`
}
type Accounts struct {
	AttackHistory []AttackHistory `json:"attackHistory"`
}
type NetworkViewer struct {
	Accounts []Accounts `json:"accounts"`
}

func buildNetworkGraphQLQuery(startDate string, endDate string, accountID string) *graphql.Request {

	query := graphql.NewRequest(`
	{
		networkViewer:viewer {
		  accounts(filter: { accountTag: $accountTag }) {
			attackHistory: ipFlows1mGroups(
			  limit: 10000
			  filter: {datetimeMinute_geq: $startDate, datetimeMinute_leq: $endDate}
			  orderBy: [sum_packets_DESC]
			) {
			  sum {
				bits
				packets
			  }
			  networkDimensions:dimensions {
				attackId
        		coloCountry
        		destinationPort
        		attackType
        		attackMitigationType
        		attackProtocol
			  }
			}
		  }
		}
	  }
	  `)

	// set any variables
	query.Var("accountTag", accountID)
	query.Var("startDate", startDate)
	query.Var("endDate", endDate)

	return query
}

// Get cloudflare metrics from GraphQL using the provided api-email and api-key parameters and returns a marshalled JSON struct and an error if something went wrong during the fetch
func getCloudflareNetworkMetrics(query *graphql.Request, apiEmail string, apiKey string) (respData networkRespDataStruct, err error) {
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
