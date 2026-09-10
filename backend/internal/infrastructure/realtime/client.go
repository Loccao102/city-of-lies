package realtime

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coder/websocket"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	sessionID string
	send      chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(hub *Hub, conn *websocket.Conn, sessionID string, queueSize int) *Client {
	if queueSize <= 0 {
		queueSize = 256
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		hub:       hub,
		conn:      conn,
		sessionID: sessionID,
		send:      make(chan []byte, queueSize),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.cancel()
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, msgBytes, err := c.conn.Read(c.ctx)
		if err != nil {
			return
		}

		var resume ClientResumeMessage
		if err := json.Unmarshal(msgBytes, &resume); err == nil && resume.Type == "client.resume" {
			// Handle client resume / replay if requested
			c.hub.HandleResume(c, resume.LastSequence)
		}
	}
}

func (c *Client) WritePump(writeTimeout, pingInterval time.Duration) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				c.conn.Close(websocket.StatusNormalClosure, "")
				return
			}
			writeCtx, writeCancel := context.WithTimeout(c.ctx, writeTimeout)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			writeCancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			pingCtx, pingCancel := context.WithTimeout(c.ctx, writeTimeout)
			err := c.conn.Ping(pingCtx)
			pingCancel()
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) Send(data []byte) bool {
	select {
	case c.send <- data:
		return true
	default:
		// Queue full: drop and disconnect slow client
		c.cancel()
		return false
	}
}
