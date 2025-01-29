package chat

type ChatServer struct {
	ID    int
	rooms []ChatRoom
}

type ChatRoom struct {
	ID int
}
