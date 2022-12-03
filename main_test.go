package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_getNodeName(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	type args struct {
		auth     string
		base_url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"tpdev-pv01", args{"PVEAPIToken=" + os.Getenv("PV_TOKEN") + "=" + os.Getenv("PV_SECRET"), "https://" + os.Getenv("PV_HOST") + "/"}, "tpdev-pv01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNodeName(tt.args.auth, tt.args.base_url); got != tt.want {
				t.Errorf("nodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProxUrl(t *testing.T) {
	type args struct {
		auth    string
		qemuUrl string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getProxUrl(tt.args.auth, tt.args.qemuUrl)
		})
	}
}

func Test_startVM(t *testing.T) {
	type args struct {
		auth    string
		qemuUrl string
		vmID    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := startVM(tt.args.auth, tt.args.qemuUrl, tt.args.vmID)
			if (err != nil) != tt.wantErr {
				t.Errorf("startVM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("startVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVMStatus(t *testing.T) {
	type args struct {
		auth    string
		qemuUrl string
		vmID    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VMStatus(tt.args.auth, tt.args.qemuUrl, tt.args.vmID); got != tt.want {
				t.Errorf("VMStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVMID(t *testing.T) {
	type args struct {
		auth    string
		qemuUrl string
		vmname  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVMID(tt.args.auth, tt.args.qemuUrl, tt.args.vmname); got != tt.want {
				t.Errorf("getVMID() = %v, want %v", got, tt.want)
			}
		})
	}
}
