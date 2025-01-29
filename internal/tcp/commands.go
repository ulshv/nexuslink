package tcp

type TCPCommand = string

const (
	// Commands sent from client to server
	CommandClientInit     TCPCommand = "client/init"
	CommandClientLogin    TCPCommand = "client/login"
	CommandClientRegister TCPCommand = "client/register"

	// Commands sent from server to client
	CommandServerInit            TCPCommand = "server/init"
	CommandServerLoginPrompt     TCPCommand = "server/login_prompt"
	CommandServerLoginSuccess    TCPCommand = "server/login_success"
	CommandServerLoginFailed     TCPCommand = "server/login_failed"
	CommandServerRegisterSuccess TCPCommand = "server/register_success"
	CommandServerRegisterFailed  TCPCommand = "server/register_failed"
)
