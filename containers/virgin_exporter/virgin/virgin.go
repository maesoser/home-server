package virgin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// A Router defined the router object
type Router struct {
	HTTPClient http.Client
	Token      string
	Auth       string
	Nonce      string
	Address    string
	User       string
	Pass       string
	Model      string
	Credential string
	Verbose    bool
	LastScrape time.Time
}

// LoginResponse is the response given by the router to a login attempt
type LoginResponse struct {
	UniqueID         string `json:"unique"`
	FamilyID         string `json:"family"`
	ModelName        string `json:"modelname"`
	Name             string `json:"name"`
	Tech             bool   `json:"tech"`
	Moca             int    `json:"moca"`
	Wifi             int    `json:"wifi"`
	ConnType         string `json:"conType"`
	GwWAN            string `json:"gwWan"`
	DefPasswdChanged string `json:"DefPasswdChanged"`
}

// NewRouter configures the http client and generates the cookie
func NewRouter(address, user, pass string) *Router {
	router := &Router{
		User:    user,
		Pass:    pass,
		Address: address,
	}
	now := time.Now()
	router.Auth = base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	router.Nonce = fmt.Sprintf("_n=%d&_=%d", rand.Intn(99999-10000)+10000, now.UnixNano()/1000000)
	router.HTTPClient = http.Client{Timeout: time.Second * 10}
	return router
}

// Login retrieves the Token needed to perform requests to the router
func (r *Router) Login() error {
	reqURL := "http://" + r.Address + "/login?arg=" + r.Auth + "&" + r.Nonce
	//log.Println(reqURL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	response, err := r.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	msg, err := DecodeLogin(body)
	if err != nil {
		return err
	}
	//log.Println(msg)
	r.Credential = string(body)
	r.Model = msg.ModelName
	return nil
}

// DecodeLogin decodes and insert into an structure the response to a login request
func DecodeLogin(data []byte) (LoginResponse, error) {
	var msg LoginResponse
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return msg, err
	}
	err = json.Unmarshal([]byte(decodedData), &msg)
	if err != nil {
		return msg, err
	}
	return msg, err
}

//Get makes a get request adding the authentication needed.
func (r *Router) Get(url string) (string, error) {
	req, err := http.NewRequest("GET", "http://"+r.Address+url, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "credential", Value: r.Credential})
	// req.Header.Set("Referer", url)
	response, err := r.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// SNMPGetOID return the result executing that specific OID
func (r *Router) SNMPGetOID(oid string) (map[string]interface{}, error) {
	dataSt := make(map[string]interface{})
	url := "/snmpGet?oids=" + ";" + oid + ";&" + r.Nonce
	data, err := r.Get(url)
	if err != nil {
		return dataSt, err
	}
	err = json.Unmarshal([]byte(data), &dataSt)
	if err != nil {
		return dataSt, err
	}
	return dataSt, nil
}

// SNMPGetOIDs return the result executing that specific OID
func (r *Router) SNMPGetOIDs(oids []string) (map[string]interface{}, error) {
	dataSt := make(map[string]interface{})
	oidList := strings.Join(oids, ",;")
	url := "/snmpGet?oids=" + ";" + oidList + ";&" + r.Nonce
	data, err := r.Get(url)
	if err != nil {
		return dataSt, err
	}
	err = json.Unmarshal([]byte(data), &dataSt)
	if err != nil {
		return dataSt, err
	}
	return dataSt, nil
}

// SNMPGetOIDs return the result executing that specific OID
func (r *Router) SNMPWalkOID(oid string) (map[string]interface{}, error) {
	original := make(map[string]interface{})
	output := make(map[string]interface{})
	url := "/walk?oids=" + ";" + oid + ";&" + r.Nonce
	data, err := r.Get(url)
	if err != nil {
		return original, err
	}
	err = json.Unmarshal([]byte(data), &original)
	if err != nil {
		return original, err
	}
	for k, v := range original {
		newKey := strings.Replace(k, oid, "", -1)
		output[newKey] = v
	}
	return output, nil
}

func TranslateSNMP(data map[string]interface{}, dict []string) (map[int]map[string]string, error) {
	output := map[int]map[string]string{}
	for key, val := range data {
		if val == "Finish" || val == "" {
			continue
		}
		OID := strings.Split(key, ".")
		index, err := strconv.Atoi(OID[2])
		if err != nil {
			return output, err
		}
		chann, err := strconv.Atoi(OID[3])
		if err != nil {
			return output, err
		}
		name := dict[index-1]
		if output[chann] == nil {
			output[chann] = make(map[string]string)
		}
		output[chann][name] = fmt.Sprintf("%v", val)
	}
	return output, nil
}

func SNMPTimeTranslate(bstring string) time.Time {
	bstring = bstring[1:] // To remove the $ sign
	year, _ := strconv.ParseInt("0x"+bstring[0:4], 0, 64)
	month, _ := strconv.ParseInt("0x"+bstring[4:6], 0, 64)
	day, _ := strconv.ParseInt("0x"+bstring[6:8], 0, 64)
	hour, _ := strconv.ParseInt("0x"+bstring[8:10], 0, 64)
	min, _ := strconv.ParseInt("0x"+bstring[10:12], 0, 64)
	sec, _ := strconv.ParseInt("0x"+bstring[12:14], 0, 64)
	dts := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	dt, err := time.Parse("2006-01-02 15:04:05", dts)
	if err != nil {
		fmt.Println(bstring)
		fmt.Print(err)
	}
	return dt
}

func (r *Router) GetClientList(w http.ResponseWriter) error {
	output := map[string]map[string]string{}
	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4115.1.20.1.1.3.42.1")
	if err != nil {
		return err
	}
	for k, c := range data {
		if c == "Finish" || c == "" {
			continue
		}
		OID := strings.Split(k, ".")
		index, err := strconv.Atoi(OID[1])
		if err != nil {
			return err
		}
		ipLen, err := strconv.Atoi(OID[4])
		if err != nil {
			return err
		}
		key := WiFiClientEntries[index-1]
		ipAddr := strings.Join(OID[len(OID)-ipLen:], ".")
		if output[ipAddr] == nil {
			output[ipAddr] = make(map[string]string)
		}
		output[ipAddr][key] = fmt.Sprintf("%v", c)
	}
	for _, channelData := range output {
		for key, val := range channelData {
			if _, err := strconv.Atoi(val); err == nil {
				fmt.Fprintf(w, "virgin_%s{hostname=\"%s\", addr=\"%s\", vendor=\"%s\"} %s\n",
					key,
					channelData["wifi_client_HostName"],
					channelData["wifi_client_IPAddrTextual"],
					channelData["wifi_client_MACMfg"], val)
			}
		}
	}
	return nil
}

// Get Downstream returns the values of docsIfDownstreamChannelTable
func (r *Router) GetDownstream(w http.ResponseWriter) error {
	data, err := r.SNMPWalkOID("1.3.6.1.2.1.10.127.1.1.1")
	if err != nil {
		return err
	}
	output, err := TranslateSNMP(data, DownChannelEntries)
	if err != nil {
		return err
	}
	for _, channelData := range output {
		for key, val := range channelData {
			fmt.Fprintf(w, "virgin_%s{channelID=\"%s\"} %s\n", key, channelData["down_chann_ID"], val)
		}
	}
	return nil
}

// GetUpstream returns the values of docsIfUpstreamChannelTable
func (r *Router) GetUpstream(w http.ResponseWriter) error {
	data, err := r.SNMPWalkOID("1.3.6.1.2.1.10.127.1.1.2")
	if err != nil {
		return err
	}
	output, err := TranslateSNMP(data, UpChannelEntries)
	if err != nil {
		return err
	}
	for _, channelData := range output {
		for key, val := range channelData {
			fmt.Fprintf(w, "virgin_%s{channelID=\"%s\"} %s\n", key, channelData["up_channel_ID"], val)
		}
	}
	return nil
}

func (r *Router) GetSignalQualityTable(w http.ResponseWriter) error {

	data, err := r.SNMPWalkOID("1.3.6.1.2.1.10.127.1.1.4")
	if err != nil {
		return err
	}
	for k, c := range data {
		if c == "Finish" || c == "" {
			continue
		}
		OID := strings.Split(k, ".")
		index, err := strconv.Atoi(OID[2])
		if err != nil {
			return err
		}
		chann, err := strconv.Atoi(OID[3])
		if err != nil {
			return err
		}
		key := SignalQualityEntry[index-1]
		if _, err := strconv.Atoi(fmt.Sprintf("%v", c)); err == nil {
			fmt.Fprintf(w, "virgin_%s{channel=\"%d\"} %s\n", key, chann, c)
		}
	}
	return nil
}

func (r *Router) GetCurrentDateTime(w http.ResponseWriter) error {
	data, err := r.SNMPGetOID("1.3.6.1.4.1.4115.1.20.1.1.5.15.0")
	if err != nil {
		return err
	}
	for _, val := range data {
		r.LastScrape = SNMPTimeTranslate(fmt.Sprintf("%v", val))
	}
	return nil
}

func (r *Router) GetDevEventTable(w http.ResponseWriter) error {
	events := make(map[string]int)
	data, err := r.SNMPWalkOID("1.3.6.1.2.1.69.1.5.8")
	if err != nil {
		return err
	}
	output, err := TranslateSNMP(data, DevEventEntry)
	if err != nil {
		return err
	}
	for _, edata := range output {
		eventTime := SNMPTimeTranslate(fmt.Sprintf("%v", edata["docsDevEvLastTime"]))
		if eventTime.After(r.LastScrape) {
			event := strings.Split(edata["docsDevEvText"], "-")[0]
			event = strings.Split(event, ":")[0]
			event = strings.Split(event, ";")[0]
			events[event] += 1
			fmt.Printf("[%v] %v\n", eventTime, edata["docsDevEvText"])
		}
	}
	for key, val := range events {
		fmt.Fprintf(w, "virgin_event{type=\"%s\"} %d\n", key, val)
	}
	return nil
}

func (r *Router) GetFwEvents(w http.ResponseWriter) error {
	events := make(map[string]int)
	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4115.1.20.1.1.5.19.1.1")
	if err != nil {
		return err
	}
	output := map[int]map[string]string{}
	for key, val := range data {
		if val == "Finish" || val == "" {
			continue
		}
		OID := strings.Split(key, ".")
		entryType, err := strconv.Atoi(OID[2])
		if err != nil {
			return err
		}
		eventID, err := strconv.Atoi(OID[3])
		if err != nil {
			return err
		}
		name := FwEventEntry[entryType-1]
		if output[eventID] == nil {
			output[eventID] = make(map[string]string)
		}
		output[eventID][name] = fmt.Sprintf("%v", val)
	}
	for _, val := range output {
		eventTime := SNMPTimeTranslate(fmt.Sprintf("%v", val["FwEvTime"]))
		if eventTime.After(r.LastScrape) {
			event := strings.Split(val["FwEvText"], " ")
			private := 0
			if event[0] == "[PRIV" {
				private = 1
			}
			ip_type := event[1]
			src_port := strings.Split(event[6], ",")[1]
			src_addr := strings.Split(strings.Split(event[6], ",")[0], ":")[1]
			dst_port := strings.Split(event[7], ",")[1]
			key := fmt.Sprintf("%d:%s:%s:%s:%s", private, src_addr, src_port, ip_type, dst_port)
			events[key] += 1
			fmt.Printf("[%v] FW Hit: [P:SRCADDR:SRCPORT:TYPE:DSTPORT] %v\n", eventTime, key)
		}
	}
	for key, val := range events {
		keys := strings.Split(key, ":")
		names, _ := net.LookupAddr(keys[1])
		hostname := "unknown"
		if len(names) >= 1 {
			hostname = names[0]
		}
		fmt.Fprintf(w, "virgin_fw_event{private=\"%s\",type=\"%s\",srcport=\"%s\",srcaddr=\"%s\",srchost=\"%s\",dstport=\"%s\"} %d\n",
			keys[0], keys[3], keys[2], keys[1], hostname, keys[4], val)
	}
	return nil
}

// 	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4491.2.1.21.1.3")
// 	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4491.2.1.21.1.2")
//  data, err := r.SNMPWalkOID("1.3.6.1.4.1.4491.2.1.20.1.1")
// return r.SNMPWalkOID("1.3.6.1.4.1.4115.1.3.4.1.9.2")
// 	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4491.2.1.20.1.2")
// 	data, err := r.SNMPWalkOID("1.3.6.1.4.1.4491.2.1.20.1.24")

// Logout return the result executing that specific OID
func (r *Router) Logout() (string, error) {
	url := "/logout?" + r.Nonce
	return r.Get(url)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := router.Login()
	if err != nil {
		log.Printf("Error logging: %v", err)
		return
	}
	err = router.GetCurrentDateTime(w)
	if err != nil {
		log.Printf("Error retrieving datetime events:\n %v\n", err)
	}
	err = router.GetDownstream(w)
	if err != nil {
		log.Printf("Error retrieving upstream stats:\n %v\n", err)
	}
	err = router.GetUpstream(w)
	if err != nil {
		log.Printf("Error retrieving downstream stats:\n %v\n", err)
	}
	err = router.GetClientList(w)
	if err != nil {
		log.Printf("Error retrieving client list:\n %v\n", err)
	}
	err = router.GetSignalQualityTable(w)
	if err != nil {
		log.Printf("Error retrieving signal quality:\n %v\n", err)
	}
	err = router.GetDevEventTable(w)
	if err != nil {
		log.Printf("Error retrieving log events:\n %v\n", err)
	}
	err = router.GetFwEvents(w)
	if err != nil {
		log.Printf("Error retrieving fw events:\n %v\n", err)
	}
	router.Logout()

}
