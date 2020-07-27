package main

import (
	"context"

   	"github.com/machinebox/graphql"
)

type workersRespDataStruct struct {
	WorkersViewer WorkersViewer `json:"workersViewer"`
}
type WorkersViewer struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Workers []Worker `json:"workers"`
}

type Worker struct {
	Info  		Info  		`json:"info"`
	Quantiles 	Quantiles 	`json:"quantiles"`
	Sum 		WorkersSum  		`json:"sum"`
}

type Info struct {
	Name string `json:"scriptName"`
}

type Quantiles struct {
	CpuTimeP50 float64 `json:"cpuTimeP50"`
	CpuTimeP75 float64 `json:"cpuTimeP75"`
	CpuTimeP99 float64 `json:"cpuTimeP99"`
	CpuTimeP999 float64 `json:"cpuTimeP999"`
}

type WorkersSum struct {
	Errors int `json:"errors"`
	Requests int `json:"requests"`
	SubRequests int `json:"subrequests"`

}

func buildWorkersGraphQLQuery(startDate string, endDate string, accountID string) *graphql.Request {
	query := graphql.NewRequest(`
		{
		workersViewer:viewer {
			accounts(filter: {accountTag: $accountTag}) {
			workers:workersInvocationsAdaptive(
				limit: 10000
				filter: {datetimeHour_geq: $startDate, datetimeHour_leq: $endDate }
			) {
				sum {
					subrequests
					requests
					errors
				}
				quantiles {
					cpuTimeP50
					cpuTimeP75
					cpuTimeP99
					cpuTimeP999
				}
				info:dimensions {
					scriptName
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
func getCloudflareWorkerMetrics(query *graphql.Request, apiEmail string, apiKey string) (respData workersRespDataStruct, err error) {
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
