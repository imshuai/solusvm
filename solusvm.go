package solusvm

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type vminfo struct {
	Hostname  string `xml:"hostname,omitempty"`
	IPaddress string `xml:"ipaddress,omitempty"`
	IPaddr    string `xml:"ipaddr,omitempty"`
	HDD       string `xml:"hdd,omitempty"`
	MEM       string `xml:"mem,omitempty"`
	BW        string `xml:"bw,omitempty"`
	VMStat    string `xml:"stat,omitempty"`
	Status    string `xml:"status,omitempty"`
	StatusMSG string `xml:"statusmsg,omitempty"`
}

func (vi *vminfo) unMarshal(msg []byte) error {
	err := xml.Unmarshal(msg, vi)
	if err != nil {
		return err
	}
	return nil
}

//HardwareInformation define virtual machine's hardware information
type HardwareInformation struct {
	Total       int    `json:"total"`
	Used        int    `json:"used"`
	Free        int    `json:"free"`
	PercentUsed string `json:"percent_used"`
}

//VirtualMachineInformation define virtual machine's information
type VirtualMachineInformation struct {
	vm        *VirtualMachine
	Hostname  string              `json:"hostname"`
	MainIP    string              `json:"main_ip"`
	IPAddress []string            `json:"ipaddress"`
	HDD       HardwareInformation `json:"hdd"`
	BW        HardwareInformation `json:"bandwith"`
	MEM       HardwareInformation `json:"memory"`
	Status    string              `json:"status"`
}

//Update virtual machine's information from solusvm api
func (vi *VirtualMachineInformation) Update() (err error) {
	vi, err = vi.vm.GetStatus()
	if err != nil {
		return err
	}
	return nil
}

//Marshal encode struct to json string
func (vi *VirtualMachineInformation) Marshal() (jsonString string, err error) {
	var byts []byte
	byts, err = json.Marshal(vi)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

//VirtualMachine define virtual machine
type VirtualMachine struct {
	key  string
	hash string
	host string
}

//NewVM create a virtual machine object's pointer
func NewVM(host, key, hash string) *VirtualMachine {
	return &VirtualMachine{
		host: host + "/api/client/command.php",
		key:  key,
		hash: hash,
	}
}

//Boot virtual machine
func (vm *VirtualMachine) Boot() error {
	msg, err := do(vm.host, "boot", vm.key, vm.hash)
	if err != nil {
		return err
	}
	vmi := &vminfo{}
	err = vmi.unMarshal(msg)
	if err != nil {
		return err
	}
	if vmi.Status != "success" {
		return errors.New(vmi.StatusMSG)
	}
	return nil
}

//Reboot virtual machine
func (vm *VirtualMachine) Reboot() error {
	msg, err := do(vm.host, "reboot", vm.key, vm.hash)
	if err != nil {
		return err
	}
	vmi := &vminfo{}
	err = vmi.unMarshal(msg)
	if err != nil {
		return err
	}
	if vmi.Status != "success" {
		return errors.New(vmi.StatusMSG)
	}
	return nil
}

//Shutdown virtual machine
func (vm *VirtualMachine) Shutdown() error {
	msg, err := do(vm.host, "shutdown", vm.key, vm.hash)
	if err != nil {
		return err
	}
	vmi := &vminfo{}
	err = vmi.unMarshal(msg)
	if err != nil {
		return err
	}
	if vmi.Status != "success" {
		return errors.New(vmi.StatusMSG)
	}
	return nil
}

//GetStatus Get virtual machine's information from solusvm api
func (vm *VirtualMachine) GetStatus() (vi *VirtualMachineInformation, err error) {
	var msg []byte
	msg, err = do(vm.host, "info", vm.key, vm.hash, "status", "hdd", "bw", "mem", "ipaddr")
	if err != nil {
		return nil, err
	}
	vmi := &vminfo{}
	err = vmi.unMarshal(msg)
	if err != nil {
		return nil, err
	}
	if vmi.Status != "success" {
		return nil, errors.New(vmi.StatusMSG)
	}
	vi.Hostname = vmi.Hostname
	vi.Status = vmi.VMStat
	vi.MainIP = vmi.IPaddress
	vi.BW = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.BW, ",")
		lm.Total, _ = strconv.Atoi(t[0])
		lm.Used, _ = strconv.Atoi(t[1])
		lm.Free, _ = strconv.Atoi(t[2])
		lm.PercentUsed = t[3]
		return lm
	}()
	vi.HDD = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.HDD, ",")
		lm.Total, _ = strconv.Atoi(t[0])
		lm.Used, _ = strconv.Atoi(t[1])
		lm.Free, _ = strconv.Atoi(t[2])
		lm.PercentUsed = t[3]
		return lm
	}()
	vi.MEM = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.MEM, ",")
		lm.Total, _ = strconv.Atoi(t[0])
		lm.Used, _ = strconv.Atoi(t[1])
		lm.Free, _ = strconv.Atoi(t[2])
		lm.PercentUsed = t[3]
		return lm
	}()
	vi.IPAddress = strings.Split(vmi.HDD, ",")
	return vi, nil
}

func do(u string, action string, flags ...string) ([]byte, error) {
	val := url.Values{}
	val.Add("action", action)
	for _, v := range flags {
		val.Add(v, "true")
	}

	resp, err := http.PostForm(u, val)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	var msg []byte
	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}
