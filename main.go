package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	var auth = "PVEAPIToken=" + os.Getenv("PV_TOKEN") + "=" + os.Getenv("PV_SECRET")
	var baseUrl = "https://" + os.Getenv("PV_HOST")
	nodeName := getNodeName(auth, baseUrl)
	var qemuUrl = baseUrl + "api2/json/nodes/" + nodeName + "/qemu/"
	// getProxUrl(auth, qemuUrl)
	vmID := getVMID(auth, qemuUrl, "avorion")
	if vmID != "0" {
		fmt.Println("Attempting to start VM")
		response, err := startVM(auth, qemuUrl, vmID)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(response)
	}
}

func getNodeName(auth string, baseUrl string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", baseUrl+"api2/json/nodes/", nil)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var nodename map[string]interface{}
	err = json.Unmarshal([]byte(body), &nodename)
	if err != nil {
		fmt.Println("Error", err)
	}
	nodeName := nodename["data"].([]interface{})[0].(map[string]interface{})["node"]
	return nodeName.(string)
}

func getProxUrl(auth, qemuUrl string) {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", qemuUrl, nil)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, body, "", "  "); err != nil {
		panic(err)
	}
	fmt.Println(dst)
}

func startVM(auth string, qemuUrl string, vmID string) (string, error) {
	checkStatus := VMStatus(auth, qemuUrl, vmID)
	if checkStatus == "stopped" {
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", qemuUrl+vmID+"/status/start", nil)
		req.Header.Set("Authorization", auth)
		req.Header.Set("Accept", "application/json")
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		response, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, body, "", "  "); err != nil {
			panic(err)
		}
		return "Started VM " + vmID, nil
	}
	if checkStatus == "paused" {
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", qemuUrl+vmID+"/status/resume", nil)
		req.Header.Set("Authorization", auth)
		req.Header.Set("Accept", "application/json")
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		response, err := client.Do(req)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		return "Resumed VM " + vmID, nil
	}
	if checkStatus == "running" {
		return "VM " + vmID + " is already running", nil
	}
	return "VM current state is: " + checkStatus, nil
}

func VMStatus(auth, qemuUrl, vmID string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", qemuUrl+vmID+"/status/current", nil)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var vmstatus map[string]interface{}
	err = json.Unmarshal([]byte(body), &vmstatus)
	if err != nil {
		fmt.Println("Error", err)
	}
	vmStatus := vmstatus["data"].(map[string]interface{})["qmpstatus"]
	return vmStatus.(string)
}

// Iterate over the json and get the "vmid" of the vm based on the "name"
func getVMID(auth, qemuUrl, vmname string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", qemuUrl, nil)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	var vmlist map[string]interface{}
	err = json.Unmarshal([]byte(body), &vmlist)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, vm := range vmlist["data"].([]interface{}) {

		if vm.(map[string]interface{})["name"] == vmname {
			return fmt.Sprint(vm.(map[string]interface{})["vmid"].(float64))
		}
	}
	return "0"
}
