package solusvm

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
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
	VMStat    string `xml:"vmstat,omitempty"`
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

type conversionInt int64

func (c conversionInt) MarshalText() ([]byte, error) {
	return []byte(unitConversion(int64(c))), nil
}

//HardwareInformation define virtual machine's hardware information
type HardwareInformation struct {
	Total       int64 `json:"total"`
	Used        int64 `json:"used"`
	Free        int64 `json:"free"`
	PercentUsed int64 `json:"percent_used"`
}
type hardwareInformationConversion struct {
	Total       conversionInt `json:"total"`
	Used        conversionInt `json:"used"`
	Free        conversionInt `json:"free"`
	PercentUsed string        `json:"percent_used"`
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

//ConversionMarshalIndent is like ConversionMarshal but applies Indent to format the output.
func (vi *VirtualMachineInformation) ConversionMarshalIndent(prefix, indent string) (jsonString string, err error) {
	viConversion := &struct {
		Hostname  string                        `json:"hostname"`
		MainIP    string                        `json:"main_ip"`
		IPAddress []string                      `json:"ipaddress"`
		HDD       hardwareInformationConversion `json:"hdd"`
		BW        hardwareInformationConversion `json:"bandwith"`
		MEM       hardwareInformationConversion `json:"memory"`
		Status    string                        `json:"status"`
	}{
		Hostname:  vi.Hostname,
		MainIP:    vi.MainIP,
		IPAddress: vi.IPAddress,
		HDD: hardwareInformationConversion{
			Total:       conversionInt(vi.HDD.Total),
			Used:        conversionInt(vi.HDD.Used),
			Free:        conversionInt(vi.HDD.Free),
			PercentUsed: strconv.FormatInt(vi.HDD.PercentUsed, 10) + "%",
		},
		BW: hardwareInformationConversion{
			Total:       conversionInt(vi.BW.Total),
			Used:        conversionInt(vi.BW.Used),
			Free:        conversionInt(vi.BW.Free),
			PercentUsed: strconv.FormatInt(vi.BW.PercentUsed, 10) + "%",
		},
		MEM: hardwareInformationConversion{
			Total:       conversionInt(vi.MEM.Total),
			Used:        conversionInt(vi.MEM.Used),
			Free:        conversionInt(vi.MEM.Free),
			PercentUsed: strconv.FormatInt(vi.MEM.PercentUsed, 10) + "%",
		},
		Status: vi.Status,
	}
	byts, err := json.MarshalIndent(viConversion, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

//ConversionMarshal encode struct to json string with unit conversion
func (vi *VirtualMachineInformation) ConversionMarshal() (jsonString string, err error) {
	return vi.ConversionMarshalIndent("", "")
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
	msg, err := do(vm, "boot")
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
	msg, err := do(vm, "reboot")
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
	msg, err := do(vm, "shutdown")
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
	msg, err = do(vm, "info", "status", "hdd", "bw", "mem", "ipaddr")
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
	vi = new(VirtualMachineInformation)
	vi.vm = vm
	vi.Hostname = vmi.Hostname
	vi.Status = vmi.VMStat
	vi.MainIP = vmi.IPaddress
	vi.BW = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.BW, ",")
		lm.Total, _ = strconv.ParseInt(t[0], 10, 0)
		lm.Used, _ = strconv.ParseInt(t[1], 10, 0)
		lm.Free, _ = strconv.ParseInt(t[2], 10, 0)
		lm.PercentUsed, _ = strconv.ParseInt(t[3], 10, 0)
		return lm
	}()
	vi.HDD = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.HDD, ",")
		lm.Total, _ = strconv.ParseInt(t[0], 10, 0)
		lm.Used, _ = strconv.ParseInt(t[1], 10, 0)
		lm.Free, _ = strconv.ParseInt(t[2], 10, 0)
		lm.PercentUsed, _ = strconv.ParseInt(t[3], 10, 0)
		return lm
	}()
	vi.MEM = func() HardwareInformation {
		lm := HardwareInformation{}
		t := strings.Split(vmi.MEM, ",")
		lm.Total, _ = strconv.ParseInt(t[0], 10, 0)
		lm.Used, _ = strconv.ParseInt(t[1], 10, 0)
		lm.Free, _ = strconv.ParseInt(t[2], 10, 0)
		lm.PercentUsed, _ = strconv.ParseInt(t[3], 10, 0)
		return lm
	}()
	vi.IPAddress = strings.Split(vmi.IPaddr, ",")
	return vi, nil
}

func do(vm *VirtualMachine, action string, flags ...string) ([]byte, error) {
	val := url.Values{}
	val.Add("action", action)
	val.Add("key", vm.key)
	val.Add("hash", vm.hash)
	for _, v := range flags {
		val.Add(v, "true")
	}
	resp, err := http.PostForm(vm.host, val)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	var msg []byte
	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return []byte("<vminfo>" + string(msg) + "</vminfo>"), nil
}

func unitConversion(b int64) string {
	if b == 0 {
		return "0.0B"
	}
	var unit string
	units := map[int]string{
		0: "B",
		1: "KB",
		2: "MB",
		3: "GB",
		4: "TB",
		5: "PB",
		6: "EB",
		7: "ZB",
		8: "YB",
	}
	pow := int(math.Floor(math.Log(float64(b)) / math.Log(1024)))
	if v, ok := units[pow]; ok {
		unit = v
	} else {
		unit = "GB"
	}
	return fmt.Sprintf("%.2f"+unit, float64(b)/math.Pow(1024, float64(pow)))
}
