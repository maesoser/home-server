package nmap

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type Vulnerability struct {
	Name        string  `json:"name"`
	Description string  `json:"descr"`
	Url         string  `json:"url"`
	Type        string  `json:"type"`
	Severiy     string  `json:"severity"`
	Score       float64 `json:"score"`
}

type Service struct {
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	RawVersion      string          `json:"raw_version"`
	Endpoints       []string        `json:"endpoints"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type NMAPScan struct {
	OutputPath string
	Targets    []string  `json:"targets"`
	UUID       string    `json:"uuid"`
	StartDate  time.Time `json:"start"`
	EndDate    time.Time `json:"end"`
	Command    string    `json:"nmap_cmd"`
	RaWOutput  string    `json:"nmap_raw"`
	Done       bool
	Services   []Service `json:"services"`
}

func RandomString(length int) string {
	var charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewScan(targets []string) *NMAPScan {
	nmapScan := &NMAPScan{
		Targets:    targets,
		UUID:       RandomString(12),
		StartDate:  time.Now(),
		Done:       false,
		Command:    "",
		OutputPath: "/scans/",
	}
	return nmapScan
}

func (n *NMAPScan) ScanFolder() string {
	return fmt.Sprintf("%s%s_%s",
		n.OutputPath,
		n.StartDate.Format("20060102T150405"),
		n.UUID,
	)
}

func (n *NMAPScan) XmlFileName(target string) string {
	return n.ScanFolder() + "/" + target + ".xml"
}

func (n *NMAPScan) Run(args string) {
	n.Command = "nmap -sV -oX <output-file> -oN - -v1 --script=vulners/vulners.nse " + args

	if _, err := os.Stat(n.ScanFolder()); os.IsNotExist(err) {
		log.Printf("Creating results folder\n")
		err := os.Mkdir(n.ScanFolder(), 0755)
		if err != nil {
			log.Printf("Error creating folder: %s", err)
		}
	}
	for _, target := range n.Targets {
		Command := fmt.Sprintf("nmap -sV -oX %s -oN - -v1 --script=vulners/vulners.nse %s %s", n.XmlFileName(target), args, target)
		log.Printf("[%s] %s\n", target, Command)
		Cmd := exec.Command("/bin/sh", "-c", Command)
		out, err := Cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing command: %s\nError out: %s\n", err, out)
		} else {
			//n.RaWOutput = fmt.Sprintf("%s\n", out)
			n.EndDate = time.Now()
			n.Done = true
			err := ioutil.WriteFile(n.ScanFolder()+"/"+target+".txt", out, 0644)
			if err != nil {
				log.Printf("Error saving stdout to file: %s\n", err)
			}
			log.Printf("[%s] Finished\n", target)
		}
	}
}

/*
<port protocol="tcp" portid="3306">
	<state state="open" reason="syn-ack" reason_ttl="64"/>
	<service name="mysql" product="MySQL" version="5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1" method="probed" conf="10">
		<cpe>cpe:/a:mysql:mysql:5.5.5-10.1.43-mariadb-0ubuntu0.18.04.1</cpe>
	</service>
	<script id="vulners" output="&#10;  MySQL 5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1: &#10;    &#9;NODEJS:602&#9;0.0&#9;https://vulners.com/nodejs/NODEJS:602">
		<table key="MySQL 5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1">
			<table>
				<elem key="cvss">0.0</elem>
				<elem key="is_exploit">false</elem>
				<elem key="type">nodejs</elem>
				<elem key="id">NODEJS:602</elem>
			</table>
		</table>
	</script>
</port>


type Service struct {
	Name            string          `json:"name"`
	Endpoints       []string        `json:"endpoints"`
		[IP]:Port
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
			Name        string  `json:"name"`
			Description string  `json:"descr"`
			Url         string  `json:"url"`
			Severiy     string  `json:"severity"`
			Score       float32 `json:"score"`
}
*/
func (n *NMAPScan) Parse() error {
	for _, target := range n.Targets {
		xmlFile, err := os.Open(n.XmlFileName(target))
		if err != nil {
			return err
		}
		defer xmlFile.Close()
		byteValue, _ := ioutil.ReadAll(xmlFile)
		var xmlRaw NmapRaw
		err = xml.Unmarshal(byteValue, &xmlRaw)
		if err != nil {
			return err
		}
		for _, host := range xmlRaw.Host {
			if host.Status.State == "up" {
				log.Println(host.Address.Addr)
				for _, port := range host.Ports.Port {
					for _, script := range port.Script {
						for _, cve := range script.Table.Table {
							var vuln Vulnerability
							for _, elem := range cve {
								if elem.key == "id" {
									vuln.Name = elem.text
								}
								if elem.key == "cvss" {
									vuln.Score, _ = strconv.ParseFloat(elem.text, 64)
								}
								if elem.key == "type" {
									vuln.Type == elem.text
								}
							}
							vuln = completeVuln(vuln)
						}
						//n.AddVuln(port.Service.Product, port.Service.Version, port.Service.Cpe[0], vuln)
					}
				}
			}
		}
	}
	return nil
}

/*
func (n *NMAPScan) AddVuln(Product, Ver, Cpe string, vuln Vulnerability) error {
	for _, service := range n.Services {
		if service.Name == serviceName {

		}
	}
}
*/

func getField(body []byte, name string) interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing CVE description: %s\n", err)
	}
	return data[name]
}
func completeVuln(vuln Vulnerability) Vulnerability {
	if vuln.Type == "cve" {
		url := fmt.Sprintf("https://vulners.com/api/v3/search/id/?id=%v", vuln.Name)
		spaceClient := http.Client{
			Timeout: time.Second * 5,
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("Error getting CVE description: %s\n", err)
		}
		req.Header.Set("User-Agent", "glanScan")

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Printf("Error getting CVE description: %s\n", err)
		}
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Printf("Error getting CVE description: %s\n", err)
		}
		data := getField( ([]byte) getField( ([]byte) getField(body, "data"), "documents"), vuln)
		vuln.Description = getField(([]byte) data, "description")
		vuln.Severiy = getField( ([]byte) getField( ([]byte) getField( ([]byte) data, "cvss3"), "cvssV3"), "baseSeverity")
	} else {
		vuln.Description = "Unknown"
		vuln.Severiy = "Unknown"
	}
	return vuln
}
