package ob

import (
	//"bytes"
	"encoding/json"
	"strconv"

	//"fmt"
	//"io/ioutil"
	"net/http"
	//"time"

	"github.com/labstack/echo"
	//"github.com/pkg/sftp"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"github.com/tuzkibug/auto-echo/base"
)

//对象方法：创建server虚拟机名称
func (dd *CDHCluster) CreateServerVM(provider *gophercloud.ProviderClient, no int, id chan string, name chan string) {
	server_name := base.CreateCDHServerName() + strconv.Itoa(no)
	server_id := base.CreateCDHServer(provider, server_name, dd.ServerFlavorID, dd.ServerImageID, dd.NetworkID)
	name <- server_name
	id <- server_id
}

//对象方法：创建agent虚拟机
func (dd *CDHCluster) CreateAgentVM(provider *gophercloud.ProviderClient, no int, id chan string, name chan string) {
	a_name := base.CreateCDHAgentName() + strconv.Itoa(no)
	agent_id := base.CreateCDHAgent(provider, a_name, dd.AgentFlavorID, dd.AgentImageID, dd.NetworkID)
	name <- a_name
	id <- agent_id
}

func BuildCDHCluster(c echo.Context) (err error) {
	d := new(CDHCluster)
	if err = c.Bind(d); err != nil {
		return
	}
	//openstack用户认证
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://" + d.OpenstackIP + ":5000/v3",
		Username:         d.Username,
		Password:         d.Password,
		DomainName:       d.DomainName,
		TenantID:         d.TenantID,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Error(err)
		return
	}

	//定义数组，存放server主机名，主机ID
	var names1 []string
	//var ids1 []string
	//定义channel提供并发
	namechs1 := make([]chan string, d.SeverNum)
	idschs1 := make([]chan string, d.SeverNum)
	for i := 0; i < d.SeverNum; i++ {
		namechs1[i] = make(chan string)
		idschs1[i] = make(chan string)
		go d.CreateServerVM(provider, i, idschs1[i], namechs1[i])
	}

	for _, namech1 := range namechs1 {
		names1 = append(names1, <-namech1)
	}

	//定义数组，存放agent主机名，主机ID
	var names2 []string
	//var ids2 []string
	//定义channel提供并发
	namechs2 := make([]chan string, d.AgentNum)
	idschs2 := make([]chan string, d.AgentNum)
	for j := 0; j < d.AgentNum; j++ {
		namechs2[j] = make(chan string)
		idschs2[j] = make(chan string)
		go d.CreateAgentVM(provider, j, idschs2[j], namechs2[j])
	}

	for _, namech2 := range namechs2 {
		names2 = append(names2, <-namech2)
	}

	serverName, _ := json.Marshal(names1)
	agentName, _ := json.Marshal(names2)

	return c.String(http.StatusOK, string(serverName)+string(agentName))
}
