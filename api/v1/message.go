package v1

import (
	"app/lib/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConnectWebsocket(c *gin.Context) {
	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	token := c.Query("token")
	ws.WebsocketServer.RegisterConn(token, unsafeConn)
	// util.NewWebsocketConnection(token, unsafeConn)
	c.Status(http.StatusNoContent)
}

func DisconnectWebsocket(c *gin.Context) {
	token := c.Query("token")
	ws.WebsocketServer.UnRegisterConn(token)
	c.Status(http.StatusNoContent)
}
