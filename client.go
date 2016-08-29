package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan *mesage

	// room is the room this client is chatting in.
	room *room

	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *mesage
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("Monday 2 Jan, 3:4")
			msg.Name = c.userData["name"].(string)
			if avatarUrl, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarUrl.(string)
			}
		} else {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
