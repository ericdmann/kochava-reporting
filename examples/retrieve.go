package main

import (
	"fmt"
	"time"

	"github.com/ericdmann/kochava-reporting"
)

const (
	apiKey  = "<>"
	appGUID = "<>"
)

func main() {
	kClient, err := kochava.NewClient(apiKey, appGUID)
	if err != nil {
		fmt.Println("Error creating Kochava reporting client: ", err)
	}

	SendAndCheckDetailRequest(kClient)
	SendAndCheckSummaryRequest(kClient)

}

func SendAndCheckSummaryRequest(kClient kochava.Kochava) {
	fmt.Println("\nSubmitting summary report.")
	//  Requests a report be delivered by email, for the last 24 hours, in UTC
	reportRequest, err := kClient.NewSummaryRequest("click",
		"1",
		[]string{"network", "country"},
		"json",
		[]string{"you@your.com"},
		time.Now().Add(-time.Hour*24),
		time.Now(),
		"UTC")

	if err != nil {
		fmt.Println("Error creating report request: ", err)
		return
	}

	token, err := kClient.SubmitRequest(reportRequest)
	if err != nil {
		fmt.Println(err)
		return
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
			break
		}
	}
}

func SendAndCheckDetailRequest(kClient kochava.Kochava) {
	fmt.Println("\nSubmitting detail report.")

	//  Requests a report be delivered by email, for the last 24 hours, in UTC
	reportRequest, err := kClient.NewDetailRequest("click",
		"json",
		[]string{"emann@kochava.com"},
		time.Now().Add(-time.Hour*24),
		time.Now(),
		"UTC")

	if err != nil {
		fmt.Println("Error creating report request: ", err)
		return
	}

	token, err := kClient.SubmitRequest(reportRequest)
	if err != nil {
		fmt.Println(err)
		return
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
			break
		}
	}
}
