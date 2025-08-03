package states

type User struct {
	Name string
	ID   string
}

type Messages struct {
	RoomName string
	Message  string
	User     User
}

type AppState any

type RoomListState struct {
	Rooms        []string
	SelectedRoom int
}

type ChatState struct {
	RoomName string
	Messages []Messages
}

type AppScreen interface {
	Render() string
	OnPressKey(key string) int
}
