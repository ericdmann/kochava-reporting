package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ericdmann/kochava-reporting"
	"github.com/ericdmann/satori-go"
)

const (
	KochavaAPIKey  = "<>"
	KochavaAppGUID = "<>"
	RTMEndpoint    = "<>"
	RTMAppKey      = "<>"
	RTMRoleName    = "<>"
	RTMRoleSecret  = "<>"
)

func main() {
	kClient, err := kochava.NewClient(KochavaAPIKey, KochavaAppGUID)
	if err != nil {
		panic("Error creating Kochava reporting client: " + err.Error())
		return
	}

	rtmClient, err := rtm.NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)
	if err != nil {
		panic("Error creating RPM client: " + err.Error())
		return
	}

	reportURL := SendAndDeliverDetailReport(kClient, time.Now().Add(-time.Hour*3), time.Now().Add(-time.Hour*2))

	if reportURL != "" {
		err = SubmitToSatori(rtmClient, reportURL)
		if err != nil {
			panic(err)
		}
	}
}

func SubmitToSatori(rtmClient rtm.RTMClient, reportURL string) error {

	resp, err := http.Get(reportURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var reportResult []map[string]interface{}
	err = json.Unmarshal(responseData, &reportResult)
	if err != nil {
		panic(err)
	}

	for _, k := range reportResult {
		str, _ := json.Marshal(k)
		rtmClient.Publish("from_kochava", string(str))
	}
	return err
}

func SendAndDeliverDetailReport(kClient kochava.Kochava, startTime time.Time, endTime time.Time) string {

	fmt.Println("\nSubmitting detail report.")

	reportRequest, err := kClient.NewDetailRequest("click",
		"json",
		[]string{},
		startTime,
		endTime,
		"UTC")

	if err != nil {
		fmt.Println("Error creating report request: ", err)
		return ""
	}

	token, err := kClient.SubmitRequest(reportRequest)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("Report submitted: ", token)

	for range time.Tick(time.Second * 10) {
		repResponse, err := kClient.CheckRequest(token)
		if err != nil {
			fmt.Print("Error checking status: ", err)
			continue
		}

		fmt.Println("["+time.Now().Local().String()+"] Status update: ", repResponse.Status, "Progress: ", repResponse.Progress)
		if repResponse.Status != "queued" && repResponse.Status != "running" {
			fmt.Println("Report link: ", repResponse.ReportLink)
			return repResponse.ReportLink
		}
	}
	return ""
}
