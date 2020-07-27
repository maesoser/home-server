package main

import (
	"context"

	"github.com/machinebox/graphql"
)

type httpRespDataStruct struct {
	Viewer Viewer `json:"viewer"`
}
type Viewer struct {
	Zones []Zones `json:"zones"`
}
type Zones struct {
	Caching  []Caching  `json:"caching"`
	Requests []Requests `json:"requests"`
	FwEvents []FwEvent  `json:"fwEvents"`
}

type Requests struct {
	RequestsData RequestsData `json:"requestsData"`
}

type RequestsData struct {
	Bytes             int `json:"bytes"`
	CachedBytes       int `json:"cachedBytes"`
	EncryptedBytes    int `json:"encryptedBytes"`
	Requests          int `json:"requests"`
	CachedRequests    int `json:"cachedRequests"`
	EncryptedRequests int `json:"encryptedRequests"`

	ResponseStatusMap    []ResponseStatusMap    `json:"responseStatusMap"`
	ClientSSLMap         []ClientSSLMap         `json:"clientSSLMap"`
	ClientHTTPVersionMap []ClientHTTPVersionMap `json:"clientHTTPVersionMap"`
	ContentTypeMap       []ContentTypeMap       `json:"contentTypeMap"`
	CountryMap           []CountryMap           `json:"countryMap"`
}

type Caching struct {
	Dimensions           Dimensions           `json:"dimensions"`
	SumEdgeResponseBytes SumEdgeResponseBytes `json:"sumEdgeResponseBytes"`
}

type Dimensions struct {
	CacheStatus     string `json:"cacheStatus"`
	HTTPMethod      string `json:"clientRequestHTTPMethodName"`
	CountryName     string `json:"clientCountryName"`
	ContentTypeName string `json:"edgeResponseContentTypeName"`
}

type SumEdgeResponseBytes struct {
	EdgeResponseBytes int `json:"edgeResponseBytes"`
}

type ResponseStatusMap struct {
	EdgeResponseStatus int `json:"edgeResponseStatus"`
	Requests           int `json:"requests"`
}

type ClientHTTPVersionMap struct {
	ClientHTTPProtocol string `json:"clientHTTPProtocol"`
	Requests           int    `json:"requests"`
}

type ContentTypeMap struct {
	ContentTypeName string `json:"edgeResponseContentTypeName"`
	Requests        int    `json:"requests"`
	Bytes           int    `json:"bytes"`
}

type ClientSSLMap struct {
	ClientSSLProtocol string `json:"clientSSLProtocol"`
	Requests          int    `json:"requests"`
}

type CountryMap struct {
	CountryName string `json:"clientCountryName"`
	Requests    int    `json:"requests"`
	Bytes       int    `json:"bytes"`
	Threats     int    `json:"threats"`
}

func buildHttpGraphQLQuery(startDate string, endDate string, zoneID string) *graphql.Request {
	query := graphql.NewRequest(`
	{
		viewer {
			zones(filter: { zoneTag: $zoneTag }) {
				caching:httpRequestsCacheGroups(
					limit: 10000
					filter: {datetimeMinute_geq: $startDate, datetimeMinute_leq: $endDate}
				) {
					dimensions {
						cacheStatus
						clientCountryName
						clientRequestHTTPMethodName
						edgeResponseContentTypeName
					}
					SumEdgeResponseBytes:sum {	
						edgeResponseBytes	
					}
				}
				requests: httpRequests1mGroups(
					limit: 10000, 
					filter: {datetimeMinute_geq: $startDate, datetimeMinute_leq: $endDate}
				) {
					requestsData:sum {
						bytes
						cachedBytes
						requests
						cachedRequests
						encryptedBytes
						encryptedRequests
						clientSSLMap{
							requests
							clientSSLProtocol
						}
						responseStatusMap{
							edgeResponseStatus
							requests
						}
						clientHTTPVersionMap{
							requests
							clientHTTPProtocol
						}
						contentTypeMap{
							requests
							bytes
							edgeResponseContentTypeName
						}
						countryMap{
							requests
							threats
							clientCountryName
							bytes
						}
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
func getCloudflareHTTPMetrics(query *graphql.Request, apiEmail string, apiKey string) (respData httpRespDataStruct, err error) {
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
