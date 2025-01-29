package tcp_commands

// type TCPCommand struct {
// 	Command string
// 	Payload []byte
// }

type TCPCommand = string

const (
	// Commands sent from client to server
	CommandClientInit     TCPCommand = "client/init"
	CommandClientLogin    TCPCommand = "client/login"
	CommandClientRegister TCPCommand = "client/register"

	// Commands sent from server to client
	CommandServerLoginPrompt     TCPCommand = "server/login_prompt"
	CommandServerLoginSuccess    TCPCommand = "server/login_success"
	CommandServerLoginFailed     TCPCommand = "server/login_failed"
	CommandServerRegisterSuccess TCPCommand = "server/register_success"
	CommandServerRegisterFailed  TCPCommand = "server/register_failed"
)

// func NewCommand(command TCPCommandString, payload []byte) TCPCommand {
// 	return TCPCommand{
// 		Command: command,
// 		Payload: payload,
// 	}
// }
