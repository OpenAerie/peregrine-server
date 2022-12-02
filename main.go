package main

import (
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
	var cur_vm_url = "api2/json/nodes/" + nodeName + "/qemu/"
	// getProxUrl(auth, base_url, cur_vm_url)
	// searchVMID(auth, base_url, cur_vm_url)
	fmt.Println(getVMID(auth, base_url, cur_vm_url, "avorion"))
	// client := &http.Client{}
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// req, err := http.NewRequest("GET", base_url+"api2/json/nodes/tpdev-pv01/qemu/", nil)
	// req.Header.Set("Authorization", auth)
	// req.Header.Set("Accept", "application/json")
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }
	// response, err := client.Do(req)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }
	// defer response.Body.Close()

	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }
	// // This Looks for Avorion server and attemps to start it.
	// var stuff map[string]interface{}
	// err = json.Unmarshal([]byte(body), &stuff)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// if stuff["data"].(map[string]interface{})["name"] == "avorion" {
	// 	vmid := stuff["data"].(map[string]interface{})["vmid"]
	// 	fmt.Println(vmid)
	// }
	// if stuff["data"].(map[string]interface{})["status"] != "running" {
	// 	fmt.Println("VM is stopped")
	// 	fmt.Println("Attempting to start VM")
	// 	startVM(auth, base_url)
	// }
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
	// dst := &bytes.Buffer{}
	// if err := json.Indent(dst, body, "", "  "); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Node Name:", dst)
	var nodename map[string]interface{}
	err = json.Unmarshal([]byte(body), &nodename)
	if err != nil {
		fmt.Println("Error", err)
	}
	nodeName := nodename["data"].([]interface{})[0].(map[string]interface{})["node"]
	return nodeName.(string)
}

// func getProxUrl(auth string, base_url string, ext_url string) {
// 	client := &http.Client{}
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	req, err := http.NewRequest("GET", base_url+ext_url, nil)
// 	req.Header.Set("Authorization", auth)
// 	req.Header.Set("Accept", "application/json")
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	response, err := client.Do(req)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	defer response.Body.Close()

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	dst := &bytes.Buffer{}
// 	if err := json.Indent(dst, body, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(dst)
// }

// func startVM(auth string, base_url string) {
// 	client := &http.Client{}
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	req, err := http.NewRequest("POST", base_url+"api2/json/nodes/tpdev-pv01/qemu/109/status/start", nil)
// 	req.Header.Set("Authorization", auth)
// 	req.Header.Set("Accept", "application/json")
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	response, err := client.Do(req)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	defer response.Body.Close()

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 		os.Exit(1)
// 	}
// 	dst := &bytes.Buffer{}
// 	if err := json.Indent(dst, body, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(dst)
// }

// Iterate over the json and get the "vmid" of the vm based on the "name"
func getVMID(auth string, base_url string, ext_url string, vmname string) float64 {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest("GET", base_url+ext_url, nil)
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
	// dst := &bytes.Buffer{}
	// if err := json.Indent(dst, body, "", "  "); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Node Name:", dst)
	var vmlist map[string]interface{}
	err = json.Unmarshal([]byte(body), &vmlist)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, vm := range vmlist["data"].([]interface{}) {

		if vm.(map[string]interface{})["name"] == vmname {
			return vm.(map[string]interface{})["vmid"].(float64)
		}
	}
	return 0
}
