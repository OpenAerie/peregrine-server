package main

import (
	"fmt"
	"log"
	"os"

	"github.com/topritchett/game-server/proxmox"
	// "github.com/topritchett/game-server/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	// Get the values from the .env file
	var auth = "PVEAPIToken=" + os.Getenv("PV_TOKEN") + "=" + os.Getenv("PV_SECRET")
	var baseUrl = "https://" + os.Getenv("PV_HOST") + "/"
	nodeName := proxmox.GetNodeName(auth, baseUrl)
	var qemuUrl = baseUrl + "api2/json/nodes/" + nodeName + "/qemu/"

	proxmox.GetProxUrl(auth, qemuUrl)
	vmName := "avorion"
	vmID := proxmox.GetVMID(auth, qemuUrl, vmName)
	if vmID != "0" {
		fmt.Println("Attempting to start VM")
		response, err := proxmox.StartVM(auth, qemuUrl, vmID)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(response)
	}
	proxmox.PauseAllVMs(auth, qemuUrl)

	// mux := http.NewServeMux()
	// server.New(mux)

	// http.ListenAndServe(":8080", mux)
}
