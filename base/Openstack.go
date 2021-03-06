package base

import (
	"fmt"
	//"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

//openstack基本操作函数，包括拉起、查询等

//创建mysql虚拟机
func CreateMysql(provider *gophercloud.ProviderClient, filename string, flavorID string, imageID string, netID string) (ServerID string) {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	//随机生成虚拟机名称
	name := CreateRandom()
	//读配置文件
	file, err := ReadAll(filename)
	if err != nil {
		log.Error(err)
		return
	}
	ss, err := servers.Create(client, servers.CreateOpts{
		Name:      name,
		FlavorRef: flavorID,
		ImageRef:  imageID,
		//AvailabilityZone:"nova",
		Networks: []servers.Network{
			servers.Network{UUID: netID},
		},
		UserData: file,
		//AdminPass: "root",
	}).Extract()

	if err != nil {
		log.Error(err)
		return
	}
	return ss.ID
}

//创建CDH虚拟机
func CreateCDHVM(provider *gophercloud.ProviderClient, name string, flavorID string, imageID string, netID string) (ServerID string) {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	//随机生成虚拟机名称
	//name := CreateCDHServerName()

	ss, err := servers.Create(client, servers.CreateOpts{
		Name:      name,
		FlavorRef: flavorID,
		ImageRef:  imageID,
		//AvailabilityZone:"nova",
		Networks: []servers.Network{
			servers.Network{UUID: netID},
		},
		//UserData: file,
		//AdminPass: "root",
	}).Extract()

	if err != nil {
		log.Error(err)
		return
	}
	return ss.ID
}

//获取虚拟机相关信息
func GetServerIP(provider *gophercloud.ProviderClient, server_id string) (result *servers.Server) {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Error(err)
		return
	}
	server, _ := servers.Get(client, server_id).Extract()
	return server

}

func ListFlavors(provider *gophercloud.ProviderClient) (result []flavors.Flavor) {
	method := "ListFlavors"
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	opts := flavors.ListOpts{}
	pager := flavors.ListDetail(client, opts)

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)

		for _, f := range flavorList {
			fmt.Println(f)
		}
		return true, err
	})
	return result
}

func GetFlavorID(provider *gophercloud.ProviderClient) (result []flavors.Flavor) {
	flavorname := ""
	method := "GetFlavorID"
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	opts := flavors.ListOpts{}
	pager := flavors.ListDetail(client, opts)

	err = fmt.Errorf("hello error")
	for err != nil {
		fmt.Println("Please enter a flavor name:")
		fmt.Scanln(&flavorname)
		fmt.Println(flavorname)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			flavorID, err := flavors.IDFromName(client, flavorname)
			if err != nil {
				fmt.Printf("%s : %v\n", method, err)
				return true, err
			}
			fmt.Print("The ID of the flavor that you search for is ")
			fmt.Println(flavorID)
			return false, nil
		})
	}
	return result
}

func GetFlavorDetails(provider *gophercloud.ProviderClient) (result []flavors.Flavor) {
	flavorid := ""
	method := "GetFlavorDetails"
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	opts := flavors.ListOpts{}
	pager := flavors.ListDetail(client, opts)

	err = fmt.Errorf("There is No Error Found")
	for err != nil {
		fmt.Println("Please enter a flavor ID:")
		fmt.Scanln(&flavorid)
		fmt.Println(flavorid)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			flavor, err := flavors.Get(client, flavorid).Extract()
			if err != nil {
				fmt.Printf("%s : %v\n", method, err)
				for {
					choice := ""
					fmt.Println("Do you want continue or quit: (c/q)")
					fmt.Scanln(&choice)
					if choice == "q" {
						goto Loop1
					}
					if choice == "c" {
						break
					}
				}
				return true, err
			}
			fmt.Println("The Details of the flavor that you search for is ")
			fmt.Printf("ID: %s\n", flavor.ID)
			fmt.Printf("Name: %s\n", flavor.Name)
			fmt.Printf("vCPUs: %d\n", flavor.VCPUs)
			fmt.Printf("Disk: %dGB\n", flavor.Disk)
			fmt.Printf("RAM: %dMB\n", flavor.RAM)
			fmt.Printf("RxTxFactor: %f\n", flavor.RxTxFactor)
			fmt.Printf("Swap: %dMB\n", flavor.Swap)
		Loop1:
			return false, nil
		})
	}
	return result
}

func ListImages(provider *gophercloud.ProviderClient) (result []images.Image) {
	method := "ListImages"
	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	//	fmt.Println(client)
	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	pager := images.List(client, images.ListOpts{})
	//	fmt.Println(pager)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)

		for _, i := range imageList {
			fmt.Println(i)
		}
		return true, err
	})
	return result
}

func ListNetworks(provider *gophercloud.ProviderClient) (result []networks.Network) {
	method := "ListNetwork"
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	//	fmt.Println(client)
	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	pager := networks.List(client, networks.ListOpts{})
	//	fmt.Println(pager)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)

		for _, n := range networkList {
			fmt.Println(n)
		}
		return true, err
	})
	return result
}

func ListSubNets(provider *gophercloud.ProviderClient) (result []subnets.Subnet) {
	method := "ListNetwork"
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	//	fmt.Println(client)
	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}
	pager := subnets.List(client, subnets.ListOpts{NetworkID: "06d08bad-2b12-4df7-aa9a-80238b067405"})

	// Define an anonymous function to be executed on each page's iteration
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		subnetList, err := subnets.ExtractSubnets(page)

		for _, s := range subnetList {
			fmt.Println(s.ID)
		}
		return true, err
	})
	return result
}

func DeleteServer(provider *gophercloud.ProviderClient) {
	serverID := ""
	fmt.Scanln(&serverID)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Printf("Create : %v", err)
		return
	}

	result := servers.Delete(client, serverID)
	fmt.Println(result)
}

func ListServers(provider *gophercloud.ProviderClient) (result []servers.Server) {
	method := "ListServers"
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Printf("%s : %v", method, err)
		return result
	}

	opts := servers.ListOpts{}
	pager := servers.List(client, opts)

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)

		for _, s := range serverList {
			fmt.Println(s)
		}
		return true, err
	})
	return result
}
