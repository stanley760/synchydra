package middleware

import "github.com/withlin/canal-go/client"

type CanalParam struct {
	ServerAddrs string
	Port        int
	Username    string
	Password    string
	Destination string
	SoTime      int32
	IdleTimeOut int32
}

// NewClient create a new client of canal, and connect it.
func NewClient(param CanalParam) *client.SimpleCanalConnector {
	connector := client.NewSimpleCanalConnector(param.ServerAddrs, param.Port, param.Username, param.Password, param.Destination, param.SoTime, param.IdleTimeOut)
	err := connector.Connect()
	if err != nil {
		panic(err)
	}
	return connector
}
