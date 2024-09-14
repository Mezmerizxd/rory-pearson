package socket

import "net"

type EndpointNamespace string

const (
	CONNECT EndpointNamespace = "connect"
	CHAT    EndpointNamespace = "chat"
)

type Request struct {
	Namespace EndpointNamespace
	Data      string
}

type Response struct {
	Status    int `json:"status"`
	Namespace EndpointNamespace
	Message   string `json:"message"`
	Data      string `json:"data"`
}

type Endpoint struct {
	Url 		 string
	Description  string
	Namespace EndpointNamespace
	Function  func(Connection, string) error
}

type Connection struct {
	Connection net.Conn
	User       *User
}

type User struct {
	Name string `json:"name"`
}

type ConnectRequestData struct {
	Name string `json:"name"`
}

type ConnectResponseData struct {
	User User `json:"user"`
}
