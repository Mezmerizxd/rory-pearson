package socket

import (
	"encoding/json"
	"errors"
	"net"
	"rory-pearson/pkg/log"
	"sync"
)

var (
	ErrorConnectionNotFound = errors.New("connection not found")

	ErrorFailedToAcceptConnection     = errors.New("failed to accept connection")
	ErrorFailedToDecodeConnectRequest = errors.New("failed to decode connect request")

	ErrorUserAlreadyConnected = errors.New("user already connected")
	ErrorUserNameEmpty        = errors.New("user name is empty")
)

type Config struct {
	Listener net.Listener
}

type Socket struct {
	Mutex       sync.Mutex
	Log         log.Log
	Listener    net.Listener
	Connections []Connection
	Endpoints   []Endpoint
}

func New(c Config) *Socket {
	l := log.New(log.Config{
		ID:            "socket",
		FileOutput:    false,
		ConsoleOutput: true,
	})

	return &Socket{
		Listener: c.Listener,
		Log:      l,
	}
}

func (s *Socket) Start() error {
	go func() {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				s.Log.Error().Err(err).Msg("Failed to accept connection")
			}

			go s.Routine(conn)
		}
	}()

	return nil
}

func (s *Socket) Stop() {
	s.Log.Info().Msg("Shutting down")

	err := s.Listener.Close()
	if err != nil {
		s.Log.Error().Err(err).Msg("Failed to close listener")
	}

	s.Log.Close()
}

func (s *Socket) Routine(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to close connection")
			return
		}
	}(conn)

	buffer := make([]byte, 1024)

	for {
		// Read from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			// If the connection is closed by the client, remove it from the list of connections
			if err.Error() == "EOF" {
				for i, c := range s.Connections {
					if c.Connection == conn {
						s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
						s.Log.Info().Msgf("Connection closed: %s : %s", c.User.Name, conn.RemoteAddr().String())
						break
					}
				}
				return
			}

			s.Log.Error().Err(err).Msgf("Error reading connection from: %s", conn.RemoteAddr().String())
			return
		}

		s.Log.Info().Msgf("Received message from %s: %s", conn.RemoteAddr().String(), buffer[:n])

		// Decode the JSON
		var request Request
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			s.Log.Error().Err(err).Msg("Failed to decode JSON")
			return
		}

		// Check if it is a connect request
		if request.Namespace == CONNECT {
			// Handle the connect request
			s.Log.Info().Msgf("Connection request from %s", conn.RemoteAddr().String())

			connection, err := s.HandleConnection(conn, request)
			if err != nil {
				s.Emit(conn, Response{
					Status:    400,
					Namespace: CONNECT,
					Message:   err.Error(),
				})
				s.Log.Error().Err(err).Msg("Error handling connection request")
				return
			}

			jsonData, err := json.Marshal(ConnectResponseData{
				User: *connection.User,
			})
			if err != nil {
				s.Emit(conn, Response{
					Status:    400,
					Namespace: CONNECT,
					Message:   err.Error(),
				})
				s.Log.Error().Err(err).Msg("Error marshalling connection")
				return
			}

			s.Emit(conn, Response{
				Status:    200,
				Namespace: CONNECT,
				Message:   "success",
				Data:      string(jsonData),
			})

			s.Log.Info().Msgf("Connection established: %s : %s", connection.User.Name, conn.RemoteAddr().String())
		} else {
			// Handle the request
			for _, e := range s.Endpoints {
				if e.Namespace == request.Namespace {
					c, err := s.GetConnection(conn)
					if err != nil {
						s.Emit(conn, Response{
							Status:    400,
							Namespace: e.Namespace,
							Message:   err.Error(),
						})
						s.Log.Error().Err(err).Msg("Error getting connection")
						return
					}

					err = e.Function(*c, request.Data)
					if err != nil {
						s.Log.Error().Err(err).Msgf("Error handling request for namespace: %s", e.Namespace)
						return
					}
					break
				}
			}
		}

		buffer = make([]byte, 1024)
	}
}

func (s *Socket) HandleConnection(conn net.Conn, request Request) (*Connection, error) {
	var data ConnectRequestData
	err := json.Unmarshal([]byte(request.Data), &data)
	if err != nil {
		return nil, ErrorFailedToDecodeConnectRequest
	}

	// Check the name is not empty
	if data.Name == "" {
		return nil, ErrorUserNameEmpty
	}

	// Check if the user is already connected
	for _, c := range s.Connections {
		if c.User.Name == data.Name {
			return nil, ErrorUserAlreadyConnected
		}
	}

	connection := Connection{
		Connection: conn,
		User: &User{
			Name: data.Name,
		},
	}

	s.Connections = append(s.Connections, connection)

	return &connection, nil
}

func (s *Socket) AddEndpoint(endpoint Endpoint) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Endpoints = append(s.Endpoints, endpoint)

	s.Log.Info().Msgf("Added endpoint: %s", endpoint.Namespace)
}

func (s *Socket) Emit(conn net.Conn, response Response) {
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(response)
	if err != nil {
		s.Log.Error().Err(err).Msg("Error encoding response")
	}
	s.Log.Info().Msgf("Response sent: %s", response.Message)
}

func (s *Socket) EmitAll(response Response) {
	for _, conn := range s.Connections {
		s.Log.Info().Msgf("Emitting to: %s", conn.User.Name)
		s.Emit(conn.Connection, response)
	}
}

func (s *Socket) GetConnection(c net.Conn) (*Connection, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, conn := range s.Connections {
		if conn.Connection == c {
			return &conn, nil
		}
	}

	return nil, ErrorConnectionNotFound
}
