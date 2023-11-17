package main

import (
	"discord_bot/gowiki"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// Bot parameters
var (
	guild = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name:        "wiki-search",
			Description: "Wiki search. Description and name may vary depending on the Language setting",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "wiki-hladaj",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "Wiki hladaj. Popis a název se mohou lišit v závislosti na nastavení Jazyk",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "title",
					Description: "Title. Description and name may vary depending on the Language setting",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Czech: "článek",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Czech: " Clánek. Popis a název se mohou lišit v závislosti na nastavení Jazyk",
					},
					Type:     discordgo.ApplicationCommandOptionString,
					Required: true,
				},
			},
		},
		{
			Name:        "wiki-get-random",
			Description: "GetRandom. Description and name may vary depending on the Language setting",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "wiki-nahodne",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "Nahodne. Popis a název se mohou lišit v závislosti na nastavení Jazyk",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "count",
					Description: "Count. Description and name may vary depending on the Language setting",
					NameLocalizations: map[discordgo.Locale]string{
						discordgo.Czech: "počet",
					},
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.Czech: "Počet. Popis a název se mohou lišit v závislosti na nastavení Jazyk",
					},
					Type:     discordgo.ApplicationCommandOptionInteger,
					Required: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},

		"wiki-search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			responses := map[discordgo.Locale]string{
				discordgo.Czech: "Ahoj! Toto je lokalizovaná zpráva",
			}
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			response := "Hi! This is a localized message with topic: "
			if option, ok := optionMap["title"]; ok {
				response += option.StringValue()
				page, err := gowiki.GetPage(option.StringValue(), -1, false, true)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(page.URL)
				response += page.URL
			}

			if r, ok := responses[i.Locale]; ok {
				response = r
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"wiki-get-random": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			responses := map[discordgo.Locale]string{
				discordgo.Czech: "Ahoj! Toto je lokalizovaná zpráva",
			}
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			response := fmt.Sprintf("Random wiki pages \n")
			if option, ok := optionMap["count"]; ok {
				pages, err := gowiki.GetRandomPages(int(option.IntValue()))
				if err != nil {
					fmt.Println(err)
				}
				for _, page := range pages {
					response += fmt.Sprintf("%s \n", page.URL)
				}
			}

			if r, ok := responses[i.Locale]; ok {
				response = r
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func RemoveCommands(s *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) {

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *guild, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Deletion complete")

}

func AddCommands(s *discordgo.Session) []*discordgo.ApplicationCommand {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *guild, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	log.Println("Adding commands complete")

	return registeredCommands

}
