// Copyright (C) 2019-2024 Chrystian Huot <chrystian@huot.qc.ca>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Access     *Access
	AuthCount  int
	Controller *Controller
	Conn       *websocket.Conn
	Send       chan *Message
	Systems    []System
	GroupsData []Group
	GroupsMap  GroupsMap
	TagsData   []Tag
	TagsMap    TagsMap
	Livefeed   *Livefeed
	SystemsMap SystemsMap
	request    *http.Request
}

func (client *Client) Init(controller *Controller, request *http.Request, conn *websocket.Conn) error {
	const (
		pongWait   = 60 * time.Second
		pingPeriod = pongWait / 10 * 9
		writeWait  = 10 * time.Second
	)

	if conn == nil {
		return errors.New("client.init: no websocket connection")
	}

	if controller.Clients.Count() >= int(controller.Options.MaxClients) {
		conn.Close()
		return nil
	}

	client.Access = NewAccess()
	client.Controller = controller
	client.Conn = conn
	client.Livefeed = NewLivefeed()
	client.Send = make(chan *Message, 8192)
	client.request = request

	go func() {
		defer func() {
			controller.Unregister <- client

			if len(client.Access.Ident) > 0 {
				controller.Logs.LogEvent(LogLevelInfo, fmt.Sprintf("listener disconnected from ip %s with ident %s", client.GetRemoteAddr(), client.Access.Ident))

			} else {
				controller.Logs.LogEvent(LogLevelInfo, fmt.Sprintf("listener disconnected from ip %s", client.GetRemoteAddr()))
			}

			client.Conn.Close()
		}()

		if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return
		}

		client.Conn.SetPongHandler(func(string) error {
			if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
				return err
			}

			return nil
		})

		for {
			_, b, err := client.Conn.ReadMessage()
			if err != nil {
				return
			}

			message := &Message{}
			if err = message.FromJson(b); err != nil {
				log.Println(fmt.Errorf("client.message.fromjson: %v", err))
				continue
			}

			if err = client.Controller.ProcessMessage(client, message); err != nil {
				log.Println(fmt.Errorf("client.processmessage: %v", err))
				continue
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(pingPeriod)

		timer := time.AfterFunc(pongWait, func() {
			client.Conn.Close()
		})

		defer func() {
			client.Send = nil

			ticker.Stop()

			if timer != nil {
				timer.Stop()
			}

			client.Conn.Close()
		}()

		for {
			select {
			case message, ok := <-client.Send:
				if !ok {
					return
				}

				if message.Command == MessageCommandConfig {
					if timer != nil {
						timer.Stop()
						timer = nil

						controller.Register <- client

						if len(client.Access.Ident) > 0 {
							controller.Logs.LogEvent(LogLevelInfo, fmt.Sprintf("new listener from ip %s with ident %s", client.GetRemoteAddr(), client.Access.Ident))

						} else {
							controller.Logs.LogEvent(LogLevelInfo, fmt.Sprintf("new listener from ip %s", client.GetRemoteAddr()))
						}
					}
				}

				if b, err := message.ToJson(); err != nil {
					log.Println(fmt.Errorf("client.message.tojson: %v", err))

				} else {
					if err = client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
						return
					}

					if err = client.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
						return
					}
				}

			case <-ticker.C:
				if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
					return
				}

				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	return nil
}

func (client *Client) GetRemoteAddr() string {
	return GetRemoteAddr(client.request)
}

func (client *Client) SendConfig(groups *Groups, options *Options, systems *Systems, tags *Tags) {
	client.SystemsMap = systems.GetScopedSystems(client, groups, tags, options.SortTalkgroups)
	client.GroupsData = groups.GetGroupsData(&client.SystemsMap)
	client.GroupsMap = groups.GetGroupsMap(&client.SystemsMap)
	client.TagsData = tags.GetTagsData(&client.SystemsMap)
	client.TagsMap = tags.GetTagsMap(&client.SystemsMap)

	var payload = map[string]any{
		"alerts":             Alerts,
		"branding":           options.Branding,
		"dimmerDelay":        options.DimmerDelay,
		"email":              options.Email,
		"groups":             client.GroupsMap,
		"groupsData":         client.GroupsData,
		"keypadBeeps":        GetKeypadBeeps(options),
		"playbackGoesLive":   options.PlaybackGoesLive,
		"showListenersCount": options.ShowListenersCount,
		"systems":            client.SystemsMap,
		"tags":               client.TagsMap,
		"tagsData":           client.TagsData,
		"time12hFormat":      options.Time12hFormat,
	}

	client.Send <- &Message{Command: MessageCommandConfig, Payload: payload}
}

func (client *Client) SendListenersCount(count int) {
	client.Send <- &Message{
		Command: MessagecommandListenersCount,
		Payload: count,
	}
}

type Clients struct {
	Map   map[*Client]bool
	mutex sync.Mutex
}

func NewClients() *Clients {
	return &Clients{
		Map:   map[*Client]bool{},
		mutex: sync.Mutex{},
	}
}

func (clients *Clients) AccessCount(client *Client) uint {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	count := uint(0)

	for c := range clients.Map {
		if c.Access == client.Access {
			count++
		}
	}

	return count
}

func (clients *Clients) Add(client *Client) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	clients.Map[client] = true
}

func (clients *Clients) Count() int {
	return len(clients.Map)
}

func (clients *Clients) EmitCall(call *Call, restricted bool) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	for c := range clients.Map {
		if (!restricted || c.Access.HasAccess(call)) && c.Livefeed.IsEnabled(call) {
			c.Send <- &Message{Command: MessageCommandCall, Payload: call}
		}
	}
}

func (clients *Clients) EmitConfig(controller *Controller) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	count := len(clients.Map)
	restricted := controller.Accesses.IsRestricted()
	showListenersCount := controller.Options.ShowListenersCount

	for c := range clients.Map {
		if restricted {
			c.Send <- &Message{Command: MessageCommandPin}

		} else {
			c.SendConfig(controller.Groups, controller.Options, controller.Systems, controller.Tags)
		}

		if showListenersCount {
			c.SendListenersCount(count)
		}
	}
}

func (clients *Clients) EmitListenersCount() {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	count := len(clients.Map)

	for c := range clients.Map {
		c.SendListenersCount(count)
	}
}

func (clients *Clients) Remove(client *Client) {
	clients.mutex.Lock()
	defer clients.mutex.Unlock()

	delete(clients.Map, client)
}
