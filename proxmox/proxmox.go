package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/topritchett/game-server/config"
)

var pvToken string = config.GoDotEnvVariable("PVE_TOKEN")
var pvHost string = config.GoDotEnvVariable("PVE_HOST")
var pvSecret string = config.GoDotEnvVariable("PVE_SECRET")
var Auth = "PVEAPIToken=" + pvToken + "=" + pvSecret
var BaseUrl = "https://" + pvHost + "/"

var NodeName = GetNodeName(Auth, BaseUrl)
var QemuUrl = BaseUrl + "api2/json/nodes/" + NodeName + "/qemu/"

func HttpReq(reqType, Auth, BaseUrl string) ([]byte, error) {
	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(reqType, BaseUrl, nil)
	req.Header.Set("Authorization", Auth)
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

	return body, err
}

func GetNodeName(Auth, BaseUrl string) string {
	body, err := HttpReq("GET", Auth, BaseUrl+"api2/json/nodes/")
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

func GetProxUrl(Auth, QemuUrl string) string {
	body, err := HttpReq("GET", Auth, QemuUrl)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, body, "", "  "); err != nil {
		panic(err)
	}
	return dst.String()
}

func StartVM(Auth string, QemuUrl string, vmID string) (string, error) {
	checkStatus := VMStatus(Auth, QemuUrl, vmID)
	if checkStatus == "stopped" {
		HttpReq("POST", Auth, QemuUrl+vmID+"/status/start")
		return "Started VM " + vmID, nil
	}
	if checkStatus == "paused" {
		HttpReq("POST", Auth, QemuUrl+vmID+"/status/resume")
		return "Resumed VM " + vmID, nil
	}
	if checkStatus == "running" {
		return "VM " + vmID + " is already running", nil
	}
	return "VM current state is: " + checkStatus, nil
}

func VMStatus(Auth, QemuUrl, vmID string) string {
	body, err := HttpReq("GET", Auth, QemuUrl+vmID+"/status/current")
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
func GetVMID(Auth, QemuUrl, vmname string) string {
	body, err := HttpReq("GET", Auth, QemuUrl)
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

func PauseVM(Auth, QemuUrl, vmID string) (string, error) {
	HttpReq("POST", Auth, QemuUrl+vmID+"/status/suspend")
	return "VM " + vmID + " successfully paused\n", nil
}

func PauseAllVMs(Auth, QemuUrl string) {
	body, err := HttpReq("GET", Auth, QemuUrl)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var vms map[string]interface{}
	err = json.Unmarshal([]byte(body), &vms)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, vm := range vms["data"].([]interface{}) {
		vmID := fmt.Sprint(vm.(map[string]interface{})["vmid"].(float64))
		vmStatus := VMStatus(Auth, QemuUrl, vmID)
		if vmStatus == "running" {
			fmt.Println("Attempting to pause VM", vmID)
			response, err := PauseVM(Auth, QemuUrl, vmID)
			if err != nil {
				fmt.Println("Error", err)
			}
			fmt.Println(response)
		}
	}
}
