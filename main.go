package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token         string
	hostname      string
	adminUsername string
	botUsername   string
)

func init() {
	flag.StringVar(&botUsername, "b", "", "Discord bot username")
	flag.StringVar(&adminUsername, "u", "", "Admin username")
	flag.StringVar(&hostname, "h", "", "Light controller hostname")
	flag.StringVar(&token, "t", "", "Discord bot token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println(botUsername + " is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

// Need to eventually abstract these details out.
func lights(name, state string) error {
	url := fmt.Sprintf("http://%s/api/1/group?name=%s&state=%s", hostname, name, state)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

// Format of valid inbound messages look like:
//   @gir lights hallway on
//   @gir lights hallway off
//   @gir lights hallway blue
//   @gir lights hallway reading
//   @gir lights hallway concentrate
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Ignore all messages that aren't from admin.
	if m.Author.Username != adminUsername {
		_, err := s.ChannelMessageSend(m.ChannelID, "You aren't an admin user.")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	normalizedMessage := m.ContentWithMentionsReplaced()
	if strings.HasPrefix(normalizedMessage, "<@&") {
		_, err := s.ChannelMessageSend(m.ChannelID, "Roles are not supported.")
		if err != nil {
			fmt.Println(err)
		}
	}
	normalizedMessage = strings.Replace(normalizedMessage, "@"+botUsername+" ", "", 1)
	tokens := strings.Split(normalizedMessage, " ")
	command := tokens[0]
	if command == "lights" {
		if len(tokens) != 3 {
			_, err := s.ChannelMessageSend(m.ChannelID, "The command _lights_ needs more arguments.")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		group := tokens[1]
		state := tokens[2]
		err := lights(group, state)
		if err != nil {
			fmt.Println(err)
		}
		feedbackMessage := fmt.Sprintf("%s lights are %s.", group, state)
		_, err = s.ChannelMessageSend(m.ChannelID, feedbackMessage)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry the command _%s_ is not supported.", command))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
