package commands

import (
	"discord_bot/ISSNow"
	"discord_bot/gowiki"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "iss-now",
			Description: "ISS current location. (International Space Station",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "iss-lokator",
			},
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Czech: "Poloha ISS. Medzinarodna vesmirna stanica.",
			},
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

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"iss-now": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := godotenv.Load(".env")

			if err != nil {
				log.Fatalf("Error loading .env file")
			}
			apikey := os.Getenv("API_KEY")
			data := ISSNow.GetLocation(apikey)

			response := fmt.Sprintf("latitude: %s\n longitude: %s\n ", data.Latitude, data.Longitude)
			fmt.Println(data.Location)
			if len(data.Location) == 0 {
				response += fmt.Sprintf("ISS is above ocean or something, not being able to retrieve location name")
			} else {
				response += fmt.Sprintf("%s \n", data.Location)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"wiki-search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			responses := map[discordgo.Locale]string{
				discordgo.Czech:     "Vysledok \n",
				discordgo.EnglishUS: "Result \n",
			}
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			response := "If you see this, it's error"

			option, ok := optionMap["článek"]
			if !ok {
				option, ok = optionMap["title"]
			}
			if ok {
				if r, ok := responses[i.Locale]; ok {
					response = r
				}
				page, err := gowiki.GetPage(option.StringValue(), -1, false, true)
				if err != nil {
					log.Fatal(err)
				}
				if len(page.URL) < 1 {
					response += "Not found"
				}
				response += page.URL
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
				discordgo.Czech:     "Random wiki stranky\n",
				discordgo.EnglishUS: "Random wiki pages\n",
			}
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			response := "If you see this, it's error"

			option, ok := optionMap["počet"]
			if !ok {
				option, ok = optionMap["count"]
			}

			if ok {
				if r, ok := responses[i.Locale]; ok {
					response = r
				}
				pages, err := gowiki.GetRandomPages(int(option.IntValue()))
				if err != nil {
					fmt.Println(err)
				}
				for _, page := range pages {
					response += fmt.Sprintf("%s \n", page.URL)
				}
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
