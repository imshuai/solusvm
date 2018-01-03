# solusvm
SolusVM API Client by golang  
---
# Usage  
```golang
package main

import(
	"github.com/imshuai/solusvm"
)

const(
	KEY = "your solusvm api key"
	HASH = "your solusvm api hash"
	HOST = "your solusvm url like https://client.solusvm.com:8787"
)

func main(){
	vm:=solusvm.NewVM(HOST,KEY,HASH)
	err:=vm.Boot()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("virtual machine boot success")
	err=vm.Reboot()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("virtual machine reboot success")
	vi:=&solusvm.VirtualMachineInformation{}
	vi,err=vm.GetStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(vi)
	var jsonstring string
	jsonstring,err=vi.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jsonstring)
	err=vm.Shutdown()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("virtual machine shutdown success")
	err=vi.Update()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(vi)
	jsonstring,err=vi.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jsonstring)
}
```
