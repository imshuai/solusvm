package solusvm

import (
	"io/ioutil"
	"net/http"
)

import (
	"net/url"
)

type vminfo struct {
	IPaddress string `xml:"ipaddress"`
	IPaddr    string `xml:"ipaddr"`
	HDD       string `xml:"hdd"`
	MEM       string `xml:"mem"`
	BW        string `xml:"bw"`
	VMStat    string `xml:"stat"`
	Status    string `xml:"status"`
	StatusMSG string `xml:"statusmsg"`
}

type limitation struct {
	Total       int    `json:"total"`
	Used        int    `json:"used"`
	Free        int    `json:"free"`
	PercentUsed string `json:"percent_used"`
}

//VirtualMachine comment to write here
type VirtualMachine struct {
	key  string
	hash string
	host string
}

//VirtualMachineInformation comment to write here
type VirtualMachineInformation struct {
	Hostname  string     `json:"hostname"`
	MainIP    string     `json:"main_ip"`
	IPAddress []string   `json:"ipaddress"`
	HDD       limitation `json:"hdd"`
	BW        limitation `json:"bandwith"`
	MEM       limitation `json:"memory"`
	Status    string     `json:"status"`
}

func NewVM(host,key,hash string)*VirtualMachine{
	return &VirtualMachine{
		host:host,
		key:key,
		hash:hash,
	}
}

func (vm *VirtualMachine)Boot()error{
	msg,err:=do(vm.host,"boot",vm.key,vm.hash)
	if err!=nil{
		return err
	}
	return nil
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
