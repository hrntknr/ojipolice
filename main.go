package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hrntknr/ojipolice/analyzer"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

var token = os.Getenv("TOKEN")

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	err := newBot()
	if err != nil {
		log.Fatal(err)
	}
}

func newBot() error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	dg.AddHandler(messageCreated)

	err = dg.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	return nil
}

func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	results := analyzer.CheckOjiLevel(m.Content)

	for _, result := range results {
		switch result.Level {
		case analyzer.Alert:
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("%s おじさんを検出しました！あなたは完全におじさんです！！\n「%s」", m.Author.Mention(), result.Sentence),
			)
		case analyzer.Warn:
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("%s ちょっとおじさんを検出しました！おじさん化に気をつけてください\n「%s」", m.Author.Mention(), result.Sentence),
			)
		}
	}
}
