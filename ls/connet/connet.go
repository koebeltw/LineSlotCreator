package connet

import (
	"github.com/koebeltw/LineSlotCreator/ls/server/connectHandler"
	"github.com/koebeltw/LineSlotCreator/ls/server/eventHandler"
	"github.com/koebeltw/LineSlotCreator/ls/loadConfig"
	_ "github.com/koebeltw/LineSlotCreator/ls/loadData"
	"github.com/koebeltw/Common/util"
	"github.com/koebeltw/Common/type"
	"github.com/koebeltw/LineSlotCreator/ls/sockets"
	"github.com/koebeltw/LineSlotCreator/ls/webSocketHandler"
	"github.com/koebeltw/Common/tcp"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	//"time"
	//"flag"
)

func init() {
	//ok := flag.Bool("ok", false, "is ok")
	//id := flag.Int("id", 0, "id")
	//port := flag.String("port", ":8080", "http listen port")
	//var name string
	//flag.StringVar(&name, "name", "123", "name")
	//
	//flag.Parse() // 這一句至關重要！！必須先解析才能拿到參數值
	//
	//fmt.Println("ok:", *ok)
	//fmt.Println("id:", *id)
	//fmt.Println("port:", *port)
	//fmt.Println("name:", name)

	sockets.Server = tcp.CreateBaseServer()

	coder := tcp.NewMsgHead()
	serverEventHandler := tcp.GetEventHandler(eventHandler.EventHandler{}, nil)
	serverConnectHandler := &connectHandler.ConnectHandler{}
	sockets.Server.SetEventHandler(serverEventHandler)
	sockets.Server.SetCoder(coder)
	sockets.Server.SetUserHandler(serverConnectHandler)
	sockets.Server.SetServerHandler(serverConnectHandler)
	sockets.Server.SetSessionMgr(tcp.NewSessionsMgr(1))
	sockets.Server.SetServerAddr(Type.Addr{
		ID:   0,
		Name: "LS",
		IP:   util.GetIPs()[0],
		Port: int(util.ConvInt32(loadConfig.Config.LSPort)),
	})

	sockets.Server.SetClientsAddr([]Type.Addr{{
		ID:   0,
		Name: "LS",
		IP:   util.GetIPs()[0],
		Port: int(util.ConvInt32(loadConfig.Config.LSPort)),
	}})

	go sockets.Server.Start()

	log.Println("Socket start at " + util.GetIPs()[0] + ":" +  loadConfig.Config.LSPort)
}

func init() {
	go func() {
		var conn *websocket.Conn
		upgrader := &websocket.Upgrader{
			//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		http.Handle("/", http.FileServer(http.Dir(".//build")))
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			var err error
			//if conn != nil {
			//	fmt.Println(conn.RemoteAddr().String())
			//	WebSocketUsers.Delete(conn.RemoteAddr().String())
			//	log.Println("pre disconnect !!")
			//	conn.Close()
			//}

			conn, err = upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err, conn.RemoteAddr().String())
				//WebSocketUsers.Delete(conn.RemoteAddr().String())
				log.Println("upgrade:", err)
				return
			}
			//if v, _:= WebSocketUsers.LoadOrStore(conn.RemoteAddr().String(), conn); v != nil {
			//	WebSocketUsers.Delete(conn.RemoteAddr().String())
			//	log.Println("pre disconnect !!")
			//	conn.Close()
			//}
			fmt.Println(conn.RemoteAddr().String())
			//WebSocketUsers.Store(conn.RemoteAddr().String(), conn)

			go func(conn *websocket.Conn) {
				defer func() {
					log.Println(conn.RemoteAddr().String(), " disconnect !!")
					//WebSocketUsers.Delete(conn.RemoteAddr().String())
					conn.Close()
				}()

				isLogin := false
				loginCount := 0
				for {
					mtype, msg, err := conn.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						break
					}
					fmt.Println("Msg:", string(msg))

					var ok bool
					if isLogin {
						msg = webSocketHandler.GetMsg(msg)
					} else {
						msg, ok = webSocketHandler.CanLogin(msg)
						if ok {
							isLogin = true
						} else {
							loginCount++

							if loginCount >= 3 {
								conn.Close()
							}
						}
					}

					if len(msg) == 0 {
						continue
					}

					log.Println(msg)
					err = conn.WriteMessage(mtype, msg)
					if err != nil {
						log.Println("write:", err)
						break
					}
				}
			}(conn)
		})

		go func() {
			log.Println("WebSocket start at " + util.GetIPs()[0] + ":7688")
			log.Fatal(http.ListenAndServe(util.GetIPs()[0]+":7688", nil))
		}()

		go func() {
			log.Println("WebSocket start at 127.0.0.1:7688")
			log.Fatal(http.ListenAndServe(":7688", nil))
		}()
	}()
}
