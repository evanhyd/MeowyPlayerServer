package server

var state = MakeServer()

func GetInstance() *Server {
	return &state
}

type Server struct {
	serverAnalytics
	serverManager
}

func MakeServer() Server {
	return Server{makeServerAnalytics(), makeServerManager()}
}
