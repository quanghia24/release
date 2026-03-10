package main

import (
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"

	pb "ride-sharing/shared/proto/driver"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleDriverWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WS: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("Missing userID in query params")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Printf("Missing packageSlug in query params")
		return
	}

	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func(){
		driver, err := driverService.Client.UnregisterDriver(r.Context(), &pb.RegisterDriverRequest{
			DriverID: userID,
		})
		if err != nil {
			log.Printf("Error calling UnregisterDriver: %v\n", err)
		}
		driverService.Close()

		log.Println("Driver unregister:", driver)
	}()


	req := &pb.RegisterDriverRequest{
		DriverID:    userID,
		PackageSlug: packageSlug,
	}

	driverRegister, err := driverService.Client.RegisterDriver(r.Context(), req)
	if err != nil {
		log.Printf("Error calling RegisterDriver: %v\n", err)
		http.Error(w, "Failed to register driver", http.StatusInternalServerError)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverRegister.Driver,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}

func handleRiderWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WS: %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("Missing userID in query params")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}
