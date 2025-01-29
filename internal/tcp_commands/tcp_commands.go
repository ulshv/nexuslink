package tcp_commands

type TCPCommand struct {
	Command string
	Payload []byte
}

type TCPCommandString = string

const (
	// Commands sent from client to server
	CommandClientInit     TCPCommandString = "client/init"
	CommandClientLogin    TCPCommandString = "client/login"
	CommandClientRegister TCPCommandString = "client/register"

	// Commands sent from server to client
	CommandServerLoginPrompt     TCPCommandString = "server/login_prompt"
	CommandServerLoginSuccess    TCPCommandString = "server/login_success"
	CommandServerLoginFailed     TCPCommandString = "server/login_failed"
	CommandServerRegisterSuccess TCPCommandString = "server/register_success"
	CommandServerRegisterFailed  TCPCommandString = "server/register_failed"
)

func NewCommand(command TCPCommandString, payload []byte) TCPCommand {
	return TCPCommand{
		Command: command,
		Payload: payload,
	}
}
