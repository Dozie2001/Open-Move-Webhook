package main

import (
	"github.com/Dozie2001/Open-Move-Webhook/api"
	"github.com/Dozie2001/Open-Move-Webhook/configs"
	"fmt"
	"os"
	"strconv"
)

func main() {
	configs.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to conver PORT to integer: %v", err))
	}

	srv := api.NewServer(uint16(port), api.BuildRoutesHandler())
	srv.Listen()
}