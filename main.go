package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/junioryono/Riot-API-Golang/apiclient"
	"github.com/ryantbvt/riot-agent/internal/framework"
)

func main() {
	// Load configs
	conf := framework.LoadEnv()

	// Initialize Riot API client (handles rate limiting via Riot-API-Golang)
	riotClient := apiclient.New(conf.RiotToken)

	// Initialize discord bot
	discordServer, err := discordgo.New("Bot " + conf.DiscToken)
	if err != nil {
		log.Fatal("Error starting Discord bot", err)
	}

	// Add handlers
	handler := framework.NewHandler(riotClient)
	discordServer.AddHandler(handler.MessageHandler)

	// for scale, but not needed
	// discordServer.Identify.Intents = discordgo.IntentGuildMessages

	// Open Discord server
	if err := discordServer.Open(); err != nil {
		log.Fatal("Error opening Discord connection", err)
	}

	defer discordServer.Close()

	log.Println("Bot is now running")

	// Wait for ctrl + C or kill sig
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("Shutting down bot")
}
