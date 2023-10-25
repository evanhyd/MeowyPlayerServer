package server

var state = makeServerState()

func GetInstance() *ServerState {
	return &state
}

type ServerState struct {
	serverAnalytics
	serverManager
}

func makeServerState() ServerState {
	return ServerState{makeServerAnalytics(), makeServerManager()}
}
