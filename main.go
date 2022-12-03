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
	var base_url = "https://" + os.Getenv("PV_HOST")
	nodeName := getNodeName(auth, base_url)
	var cur_qemu_url = base_url + "api2/json/nodes/" + nodeName + "/qemu/"
	// getProxUrl(auth, cur_qemu_url)
	vmID := getVMID(auth, cur_qemu_url, "avorion")
	if vmID != "0" {
		fmt.Println("Attempting to start VM")
		response, err := startVM(auth, cur_qemu_url, vmID)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(response)
	}
}

func getNodeName(auth string, base_url string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", base_url+"api2/json/nodes/", nil)
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

func getProxUrl(auth, cur_qemu_url string) {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", cur_qemu_url, nil)
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

func startVM(auth string, cur_qemu_url string, vmID string) (string, error) {
	check_if_running := getVMStatus(auth, cur_qemu_url, vmID)
	if check_if_running == "stopped" {
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", cur_qemu_url+vmID+"/status/start", nil)
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
	if check_if_running == "paused" {
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		req, err := http.NewRequest("POST", cur_qemu_url+vmID+"/status/resume", nil)
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
		return "Resumed VM " + vmID, nil
	}

	return "VM current state is: " + check_if_running, nil
}

func getVMStatus(auth, cur_qemu_url, vmID string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", cur_qemu_url+vmID+"/status/current", nil)
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
func getVMID(auth, cur_qemu_url, vmname string) string {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", cur_qemu_url, nil)
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
