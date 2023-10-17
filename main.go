package main

//Import Packages
import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	Token        string `json:"Token"`
	SQLPassword  string `json:"SQLPassword"`
	RapidAPIKey  string `json:"X-RapidAPI-Key"`
	RapidAPIHost string `json:"X-RapidAPI-Host"`

	AlexNames    []string `json:"AlexNames"`
	AlexIDs      []string `json:"AlexIDs"`
	PatrickNames []string `json:"PatrickNames"`
	PatrickIDs   []string `json:"PatrickIDs"`
	KrikorNames  []string `json:"KrikorNames"`
	KrikorIDs    []string `json:"KrikorIDs"`
	FelixNames   []string `json:"FelixNames"`
	FelixIDs     []string `json:"FelixIDs"`
	AlecNames    []string `json:"AlecNames"`
	AlecIDs      []string `json:"AlecIDs"`
	GabrielNames []string `json:"GabrielNames"`
	GabrielIDs   []string `json:"GabrielIDs"`
	KennethNames []string `json:"KennethNames"`
	KennethIDs   []string `json:"KennethIDs"`
	QuincyNames  []string `json:"QuincyNames"`
	QuincyIDs    []string `json:"QuincyIDs"`

	AlexNamesPlayin    []string `json:"AlexNamesPlayin"`
	AlexIDsPlayin      []string `json:"AlexIDsPlayin"`
	PatrickNamesPlayin []string `json:"PatrickNamesPlayin"`
	PatrickIDsPlayin   []string `json:"PatrickIDsPlayin"`
	KrikorNamesPlayin  []string `json:"KrikorNamesPlayin"`
	KrikorIDsPlayin    []string `json:"KrikorIDsPlayin"`
	FelixNamesPlayin   []string `json:"FelixNamesPlayin"`
	FelixIDsPlayin     []string `json:"FelixIDsPlayin"`
	AlecNamesPlayin    []string `json:"AlecNamesPlayin"`
	AlecIDsPlayin      []string `json:"AlecIDsPlayin"`
	GabrielNamesPlayin []string `json:"GabrielNamesPlayin"`
	GabrielIDsPlayin   []string `json:"GabrielIDsPlayin"`
	KennethNamesPlayin []string `json:"KennethNamesPlayin"`
	KennethIDsPlayin   []string `json:"KennethIDsPlayin"`
	QuincyNamesPlayin  []string `json:"QuincyNamesPlayin"`
	QuincyIDsPlayin    []string `json:"QuincyIDsPlayin"`

	PlayinList   []string `json:"PlayinList"`
	PlayinListID []string `json:"PlayinListID"`
}

type Outcomes_Info struct {
	Team string  `json:"name"`
	Odds float64 `json:"price"`
}

type Markets_Info struct {
	Outcomes []Outcomes_Info `json:"outcomes"`
}

type Bookmarkers_Info struct {
	Title   string         `json:"title"`
	Markets []Markets_Info `json:"markets"`
}

type Match_Info struct {
	Home_team   string             `json:"home_team"`
	Away_team   string             `json:"away_team"`
	Date        string             `json:"commence_time"`
	League      string             `json:"sport_title"`
	Bookmarkers []Bookmarkers_Info `json:"bookmakers"`
}

const chatString = "1159737801647079475"
const spreadsheetID = "1qC-yEWpy304or_rcPRA06QLqCy5WbO1ZqWOsDmO4NRs"

var db *sql.DB

func main() {
	//Load Config File
	//go loggingTool()

	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	dsn := "root:" + config.SQLPassword + "@tcp(127.0.0.1:3306)/aris_bot"

	//Connect to database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Create Discordgo session
	sess, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
	}

	//Process Discord request
	sess.AddHandler(messageCreate)

	//Set intents of reading from discord server
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	//Open/Close Session
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	//Web Scrap
	fmt.Println("Bot online!")

	//Put concurrent functions here

	go sendMessages(sess)

	//go timedMessages(sess)

	//Make sure bot doesnt end after starting, only exits if Ctrl + C on cmd
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	//If session ID matches dont do anything
	if m.Author.ID == s.State.User.ID {
		return
	}
	userInput := m.Content
	//Parse command
	if strings.HasPrefix(userInput, "!") {
		parts := strings.Split(userInput, " ")
		if len(parts) > 1 && parts[0] == "!worlds" {
			switch strings.ToLower(parts[1]) {
			case "alex":
				updateWorlds(s, config.AlexNames, config.AlexIDs, "alex", true)
			case "patrick":
				updateWorlds(s, config.PatrickNames, config.PatrickIDs, "patrick", true)
			case "krikor":
				updateWorlds(s, config.KrikorNames, config.KrikorIDs, "krikor", true)
			case "felix":
				updateWorlds(s, config.FelixNames, config.FelixIDs, "felix", true)
			case "alec":
				updateWorlds(s, config.AlecNames, config.AlecIDs, "alec", true)
			case "gabriel":
				updateWorlds(s, config.GabrielNames, config.GabrielIDs, "gabriel", true)
			case "kenneth":
				updateWorlds(s, config.KennethNames, config.KennethIDs, "kenneth", true)
			case "quincy":
				updateWorlds(s, config.QuincyNames, config.QuincyIDs, "quincy", true)
			case "update":
				updateWorlds(s, config.AlexNames, config.AlexIDs, "alex", true)
				updateWorlds(s, config.PatrickNames, config.PatrickIDs, "patrick", true)
				updateWorlds(s, config.KrikorNames, config.KrikorIDs, "krikor", true)
				updateWorlds(s, config.FelixNames, config.FelixIDs, "felix", true)
				updateWorlds(s, config.AlecNames, config.AlecIDs, "alec", true)
				updateWorlds(s, config.GabrielNames, config.GabrielIDs, "gabriel", true)
				updateWorlds(s, config.KennethNames, config.KennethIDs, "kenneth", true)
				updateWorlds(s, config.QuincyNames, config.QuincyIDs, "quincy", true)
				updatePlayinTeams(s, config.PlayinList, config.PlayinListID)
				s.ChannelMessageSend(chatString, "Done Updating")
			default:
				getWorldsID(s, parts[1], true)
			}
		}
		if len(parts) > 1 && parts[0] == "!playin" {
			// switch strings.ToLower(parts[1]) {
			// case "alex":
			// 	updateWorlds(s, config.AlexNamesPlayin, config.AlexIDsPlayin, "alex", config.PlayinList, false)
			// case "patrick":
			// 	updateWorlds(s, config.PatrickNamesPlayin, config.PatrickIDsPlayin, "patrick", config.PlayinList, false)
			// case "krikor":
			// 	updateWorlds(s, config.KrikorNamesPlayin, config.KrikorIDsPlayin, "krikor", config.PlayinList, false)
			// case "felix":
			// 	updateWorlds(s, config.FelixNamesPlayin, config.FelixIDsPlayin, "felix", config.PlayinList, false)
			// case "alec":
			// 	updateWorlds(s, config.AlecNamesPlayin, config.AlecIDsPlayin, "alec", config.PlayinList, false)
			// case "gabriel":
			// 	updateWorlds(s, config.GabrielNamesPlayin, config.GabrielIDsPlayin, "gabriel", config.PlayinList, false)
			// case "kenneth":
			// 	updateWorlds(s, config.KennethNamesPlayin, config.KennethIDsPlayin, "kenneth", config.PlayinList, false)
			// case "quincy":
			// 	updateWorlds(s, config.QuincyNamesPlayin, config.QuincyIDsPlayin, "quincy", config.PlayinList, false)
			// case "update":
			// 	updateWorlds(s, config.AlexNamesPlayin, config.AlexIDsPlayin, "alex", config.PlayinList, false)
			// 	updateWorlds(s, config.PatrickNamesPlayin, config.PatrickIDsPlayin, "patrick", config.PlayinList, false)
			// 	updateWorlds(s, config.KrikorNamesPlayin, config.KrikorIDsPlayin, "krikor", config.PlayinList, false)
			// 	updateWorlds(s, config.FelixNamesPlayin, config.FelixIDsPlayin, "felix", config.PlayinList, false)
			// 	updateWorlds(s, config.AlecNamesPlayin, config.AlecIDsPlayin, "alec", config.PlayinList, false)
			// 	updateWorlds(s, config.GabrielNamesPlayin, config.GabrielIDsPlayin, "gabriel", config.PlayinList, false)
			// 	updateWorlds(s, config.KennethNamesPlayin, config.KennethIDsPlayin, "kenneth", config.PlayinList, false)
			// 	updateWorlds(s, config.QuincyNamesPlayin, config.QuincyIDsPlayin, "quincy", config.PlayinList, false)
			// default:
			getWorldsID(s, parts[1], false)
			// }

			// } else {
			// switch strings.ToLower(userInput) {
			// case "!ping":
			// 	s.ChannelMessageSend(m.ChannelID, "pong") //Check if someone types hello in chat
			// case "!profile": //Command for Profile info (Currency balance)
			// 	{
			// 		balance := getBalance(m.Author.ID, s)
			// 		s.ChannelMessageSend(chatString, "User: "+m.Author.Username+"\nBalance: "+strconv.Itoa(balance))
			// 	}
			// // case "!odds":
			// // 	GetOdds(s)
			// case "!register":
			// 	registerUser(m.Author.ID, m.Author.Username, s)
			// case "!update":
			// 	updateBalance(m.Author.ID, s)
			// case "!scores":
			// 	getPastGames(s)
			// case "!games":
			// 	getTodaysGames(s)
			// default: //!bet, ....
			// 	{
			// 		betNBAGame(m.Author.ID, userInput, s)
			// 		return
			// 	}
			// }
		}
	}
}

func getWorldsID(s *discordgo.Session, playerName string, mainEvent bool) {
	id := "NOT FOUND"
	listOfMatchURL := []string{}
	url := "https://gol.gg/players/list/season-ALL/split-ALL/tournament-World%20Championship%202023/"
	if !mainEvent {
		url = "https://gol.gg/players/list/season-ALL/split-ALL/tournament-Worlds%20Play-In%202023/"
	}
	c := colly.NewCollector(colly.AllowedDomains("www.gol.gg", "gol.gg"))

	//Before HTTP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnHTML("table.table_list tbody tr", func(h *colly.HTMLElement) {
		h.DOM.Find("a").Each(func(_ int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			name := s.Text()
			if exists {
				listOfMatchURL = append(listOfMatchURL, link)
				if strings.EqualFold(name, playerName) {
					id = strings.TrimSpace((strings.Split(link, "/"))[2])
					return
				}
			}
		})
	})

	//Error occured during request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error")
	})

	c.Visit(url)
	qualFlag := false
	qualNames := []string{"adam", "crownie", "gori", "huhi", "labrov", "licorice", "river", "sheo", "stixxay", "nuc"}
	for _, s := range qualNames {
		if playerName == s {
			qualFlag = true
		}
	}

	if qualFlag {
		c.Visit("https://gol.gg/players/list/season-ALL/split-ALL/tournament-Worlds%20Qualifying%20Series%202023/")
	}

	if id == "NOT FOUND" {
		s.ChannelMessageSend(chatString, "Not found")

	} else {
		points := getWorldsStats(s, id, playerName, true, mainEvent, false)
		msg := fmt.Sprintf("# Player: %s Total Points: %.2f", playerName, points)
		s.ChannelMessageSend(chatString, msg)
	}
}

func updatePlayinTeams(s *discordgo.Session, playinList []string, playinListID []string) {
	for index, value := range playinListID {
		points := getWorldsStats(s, value, playinList[index], false, true, false)
		points += getWorldsStats(s, value, playinList[index], false, false, true)
		msg := fmt.Sprintf("### Player: %s Points: %.2f", playinList[index], points)
		insertInSheet(spreadsheetID, points, "", index, true)
		s.ChannelMessageSend(chatString, msg)
	}
}

func updateWorlds(s *discordgo.Session, listUsers []string, list []string, user string, mainEvent bool) {
	for index, value := range list {
		points := getWorldsStats(s, value, listUsers[index], false, mainEvent, false)
		msg := fmt.Sprintf("### Player: %s Points: %.2f", listUsers[index], points)
		insertInSheet(spreadsheetID, points, user, index+3, mainEvent)
		s.ChannelMessageSend(chatString, msg)
	}
}

func getWorldsStats(s *discordgo.Session, id string, playerName string, infoFlag bool, mainEvent bool, qualFlagCheck bool) float64 {
	totalPoints := 0.0
	list := []string{}
	listOfMatchURL := []string{}
	skip := false
	resultIndex := 0
	time := ""
	result := ""

	url := "https://gol.gg/players/player-matchlist/" + id + "/season-ALL/split-ALL/tournament-World%20Championship%202023/"
	if !mainEvent {
		url = "https://gol.gg/players/player-matchlist/" + id + "/season-ALL/split-ALL/tournament-Worlds%20Play-In%202023/"
	}
	c := colly.NewCollector(colly.AllowedDomains("www.gol.gg", "gol.gg"))

	//Before HTTP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnHTML("table.table_list tbody tr", func(h *colly.HTMLElement) {
		h.DOM.Find("td.text-center").Each(func(_ int, s *goquery.Selection) {
			list = append(list, s.Text())
			link, exists := s.Find("a").Attr("href")
			if exists {
				if !skip { //Skips 2nd link
					listOfMatchURL = append(listOfMatchURL, link)
					skip = true
				} else {
					skip = false
				}
			}
		})
	})

	//Error occured during request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error")
	})

	c.Visit(url)
	qualFlag := false
	qualNames := []string{"adam", "crownie", "gori", "huhi", "labrov", "licorice", "river", "sheo", "stixxay", "nuc"}
	for _, str := range qualNames {
		if playerName == str {
			qualFlag = true
		}
	}

	if qualFlag && !qualFlagCheck {
		c.Visit("https://gol.gg/players/player-matchlist/" + id + "/season-ALL/split-ALL/tournament-Worlds%20Qualifying%20Series%202023/")
	}

	for _, value := range listOfMatchURL {
		points, msg := checkGameStats(value, playerName)
		result = list[resultIndex]
		if result == "Victory" {
			msg += " *Win bonus +2* "
			points += 2
			resultIndex += 2
			time = strings.TrimSpace(list[resultIndex])
			if convertTime(time) {
				msg += " *Time bonus +2* "
				points += 2
			}
		} else {
			resultIndex += 2
			time = strings.TrimSpace(list[resultIndex])
		}
		resultIndex += 4
		msg += "\n**Result:** " + result + " **Duration:** " + time
		msg += "\n**POINTS EARNED: **" + strconv.FormatFloat(points, 'f', 2, 64)
		if infoFlag {
			s.ChannelMessageSend(chatString, msg)
		}
		totalPoints += points
	}
	return totalPoints
}

func checkGameStats(linkString string, player string) (float64, string) {
	kda, total, msg := "0/0/0", 0.0, ""
	url := "https://gol.gg/" + linkString[2:]
	c := colly.NewCollector(colly.AllowedDomains("www.gol.gg", "gol.gg"))
	list := []string{}
	listTeamInfo := []string{}
	side := 0
	listFB := []bool{}
	killBonus := 0.0
	assistBonus := 0.0
	listFormatString := []string{}
	listFormatFloat := []float64{}
	//Before HTTP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
		listFormatString = append(listFormatString, player, url)
	})

	c.OnHTML("div.row table.playersInfosLine tbody", func(h *colly.HTMLElement) {
		h.DOM.Find("td:not([colspan])").Each(func(_ int, s *goquery.Selection) {
			//Find what side player is on
			text := s.Find("a").Text()
			if len(text) > 1 {
				list = append(list, text)
			}
			styleAttr, exists := s.Attr("style")
			if exists && strings.Contains(styleAttr, "text-align") {
				text := s.Text()
				list = append(list, text)
			}
		})
	})

	c.OnHTML("div.col-2", func(h *colly.HTMLElement) {
		h.DOM.Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			if len(text) > 1 {
				listTeamInfo = append(listTeamInfo, strings.TrimSpace(text))
			}
		})
	})

	c.OnHTML("div.col-2", func(h *colly.HTMLElement) {
		h.DOM.Find("img").Each(func(_ int, s *goquery.Selection) {
			altAttr, exists := s.Attr("alt")
			if exists && altAttr == "First Blood" {
				listFB = append(listFB, true) // Blue = 1, Red = 10
			} else {
				listFB = append(listFB, false)
			}
		})
	})

	//Error occured during request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error")
	})

	c.Visit(url)

	for index, value := range list {
		if strings.EqualFold(strings.ToLower(value), strings.ToLower(player)) {
			if index <= 14 {
				side = 0
			} else {
				side = 7
			}
			kda = list[index+1]
			minions := strings.TrimSpace(list[index+2])
			fbFlag := (listFB[1] && side == 0) || (!listFB[1] && side == 7)
			listFormatString = append(listFormatString, kda, minions, listTeamInfo[side+1], listTeamInfo[side+2], listTeamInfo[side+3], strconv.FormatBool(fbFlag))
			kills, _ := strconv.ParseFloat((strings.Split(kda, "/"))[0], 64)
			if kills >= 10.0 {
				killBonus = 2.0
			}
			kills = kills * 3
			deaths, _ := strconv.ParseFloat((strings.Split(kda, "/"))[1], 64)
			assists, _ := strconv.ParseFloat((strings.Split(kda, "/"))[2], 64)
			if assists >= 10.0 {
				assistBonus = 2.0
			}
			assists = assists * 2
			minionsFloat, _ := strconv.ParseFloat(minions, 64)
			minionsFloat = minionsFloat * 0.02
			towers, _ := strconv.ParseFloat(listTeamInfo[side+1], 64)
			dragons, _ := strconv.ParseFloat(listTeamInfo[side+2], 64)
			dragons = dragons * 2
			barons, _ := strconv.ParseFloat(listTeamInfo[side+3], 64)
			barons = barons * 3
			fb := 0.0
			if fbFlag {
				fb = 2.0
			}
			listFormatFloat = append(listFormatFloat, kills, deaths, assists, minionsFloat, towers, dragons, barons)

			total += kills + assists + minionsFloat + towers + dragons + barons + fb + killBonus + assistBonus
			total -= deaths
			msg = formatStatMsg(listFormatString, listFormatFloat)
			break
		}
	}
	return total, msg
}

func insertInSheet(spreadsheetID string, points float64, user string, index int, mainEvent bool) {
	type Credentials struct {
		Token string `json:"token"`
	}
	credentialsJSON, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Println("Error reading credentials file:", err)
		return
	}

	// Create a Credentials struct to hold the data
	var credentials Credentials

	// Unmarshal the JSON into the Credentials struct
	err = json.Unmarshal(credentialsJSON, &credentials)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Load your credentials
	creds, err := google.JWTConfigFromJSON([]byte(credentialsJSON), sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Create a new JWT client
	client := creds.Client(context.Background())

	// Initialize the Sheets API client
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	rangeToWrite := "Main!"
	if !mainEvent {
		index += 10
	}
	switch user {
	case "alex":
		rangeToWrite += "C" + strconv.Itoa(index)
	case "patrick":
		rangeToWrite += "E" + strconv.Itoa(index)
	case "krikor":
		rangeToWrite += "G" + strconv.Itoa(index)
	case "felix":
		rangeToWrite += "I" + strconv.Itoa(index)
	case "alec":
		rangeToWrite += "K" + strconv.Itoa(index)
	case "gabriel":
		rangeToWrite += "M" + strconv.Itoa(index)
	case "kenneth":
		rangeToWrite += "O" + strconv.Itoa(index)
	case "quincy":
		rangeToWrite += "Q" + strconv.Itoa(index)
	case "":
		rangeToWrite += playinTeamsUpdate(index)
	}

	//Alex C3-C9 Patrick E3-E9 Krikor G3-G9 Felix I3-I9
	//Alec K3-K9 Gabriel M3-M9 Kenneth O3-O9 Quincy Q3-Q9

	//Alex C13-C17 Patrick E13-E17 Krikor G13-G17 Felix I13-I17
	//Alec K13-K17 Gabriel M13-M17 Kenneth O13-O17 Quincy Q13-Q17
	value := points

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{value}},
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, rangeToWrite, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Fatalf("Unable to update value: %v", err)
	}

	log.Println("Value updated successfully.")
}

func playinTeamsUpdate(index int) string {
	switch index {
	case 0:
		return "I" + strconv.Itoa(13)
	case 1:
		return "I" + strconv.Itoa(14)
	case 2:
		return "E" + strconv.Itoa(15)
	case 3:
		return "E" + strconv.Itoa(16)
	case 4:
		return "I" + strconv.Itoa(17)
	case 5:
		return "O" + strconv.Itoa(13)
	case 6:
		return "Q" + strconv.Itoa(14)
	case 7:
		return "Q" + strconv.Itoa(15)
	case 8:
		return "M" + strconv.Itoa(16)
	case 9:
		return "K" + strconv.Itoa(17)
	}
	return ""
}

func formatStatMsg(listFormatString []string, listFormatFloat []float64) string {
	for _, v := range listFormatFloat {
		listFormatString = append(listFormatString, fmt.Sprintf("%f", v))
	}
	msg := "** " + strings.ToUpper(listFormatString[0]) + " STATS **<" + listFormatString[1] + ">\n" +
		">>> __**PLAYER STATS:**__\n" +
		"KDA: " + listFormatString[2] + "\n" +
		"MINIONS: " + listFormatString[3] + "\n\n" +
		"__**TEAM STATS:**__\n" +
		"TOWERS: " + listFormatString[4] + " DRAGONS: " + listFormatString[5] +
		" BARONS: " + listFormatString[6] + " FB: " + listFormatString[7] + "\n\n" +
		"__**POINT DISTRIBUTION:**__\n" +
		"K: " + listFormatString[8][:len(listFormatString[8])-4] + " D: -" + listFormatString[9][:len(listFormatString[9])-4] +
		" A: " + listFormatString[10][:len(listFormatString[10])-4] + " M: " + listFormatString[11][:len(listFormatString[11])-4] + "\n" +
		"T: " + listFormatString[12][:len(listFormatString[12])-4] + " DG: " + listFormatString[13][:len(listFormatString[13])-4] +
		" B: " + listFormatString[14][:len(listFormatString[14])-4] + "\n"
	if listFormatFloat[0]/3 >= 10 {
		msg += " *Kill Bonus +2* "
	}
	if listFormatFloat[2]/2 >= 10 {
		msg += " *Assist Bonus +2* "
	}
	if listFormatString[7] == "true" {
		msg += " *FB Bonus +2* "
	}
	return msg
}

func registerUser(id string, username string, s *discordgo.Session) {
	currDate := time.Now().Format("2006-01-02")
	//Insert User Entry
	query := "INSERT INTO profiles VALUES(?, ?, ?)"
	_, err := db.Exec(query, id, username, currDate) //Look into context timeout
	if err != nil {
		//Case where duplicate entry
		if strings.Contains(err.Error(), "Duplicate entry") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "Already registered user! Hmph~!")
		} else {
			panic(err.Error())
		}
	}

	//Insert default balance
	uuid := strings.Replace(uuid.New().String(), "-", "", -1)
	query = "INSERT INTO balance VALUES(?, 1000, ?, ?);"
	_, err = db.Exec(query, uuid, currDate, id) //Look into context timeout
	if err != nil {
		//Case where duplicate entry
		if strings.Contains(err.Error(), "Duplicate entry") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "Already registered balance! Hmph~!")
		} else {
			panic(err.Error())
		}
	} else {
		s.ChannelMessageSend(chatString, "Registered user "+username)
	}
}

func updateBalance(id string, s *discordgo.Session) {
	randNum := rand.Intn(100)
	currDate := time.Now().Format("2006-01-02")
	if currDate != getLoginDate(id, s) {
		s.ChannelMessageSend(chatString, "Adding "+strconv.Itoa(randNum)+" credits\n")
		balance := getBalance(id, s) + randNum
		query := "UPDATE balance SET balance = ?, dailyLoginDate = ? WHERE ProfileID = ?;"
		_, err := db.Exec(query, balance, currDate, id) //Look into context timeout
		if err != nil {
			//Case where no matching rows
			if strings.Contains(err.Error(), "no rows") {
				// Handle the permission-denied error
				s.ChannelMessageSend(chatString, "No matching rows! Hmph~!")
			} else {
				panic(err.Error())
			}
		} else {
			s.ChannelMessageSend(chatString, "Updated Balance to "+strconv.Itoa(balance)+"\n")
		}
	} else {
		s.ChannelMessageSend(chatString, "Already logged!")
	}
}

func getBalance(id string, s *discordgo.Session) int {
	var number int
	err := db.QueryRow("SELECT balance from balance where ProfileID = ?", id).Scan(&number) //Gets and stores number
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "No matching rows! Hmph~!")
		} else {
			panic(err.Error())
		}
	}
	return number
}

func getBalanceId(id string, s *discordgo.Session) string {
	var idString string
	err := db.QueryRow("SELECT ID from balance where ProfileID = ?", id).Scan(&idString) //Gets and stores number
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "No matching rows! Hmph~!")
		} else {
			panic(err.Error())
		}
	}
	return idString
}

func getLoginDate(id string, s *discordgo.Session) string {
	var date string
	err := db.QueryRow("SELECT dailyLoginDate from balance where ProfileID = ?", id).Scan(&date) //Gets and stores number
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "No matching rows! Hmph~!")
		} else {
			panic(err.Error())
		}
	} else {
		s.ChannelMessageSend(chatString, "Last login date: "+date)
	}
	return date
}

func GetOdds(s *discordgo.Session) {
	url := "https://odds.p.rapidapi.com/v4/sports/basketball_nba/odds?regions=us&oddsFormat=decimal&markets=h2h%2Cspreads&dateFormat=iso"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "")
	req.Header.Add("X-RapidAPI-Host", "")

	//req.Header.Add("X-RapidAPI-Key", config.RapidAPIKey)
	//req.Header.Add("X-RapidAPI-Host", config.RapidAPIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var inform []Match_Info
	json.Unmarshal(body, &inform)

	fmt.Printf("Home: %s, Away: %s\n", inform[0].Home_team, inform[0].Away_team)
	fmt.Printf("Length of Bookmarkers: %v\n", len(inform[0].Bookmarkers))
	var totalOne float64
	var totalTwo float64
	for i := 0; i < len(inform[0].Bookmarkers); i++ {
		totalOne = totalOne + inform[0].Bookmarkers[i].Markets[0].Outcomes[0].Odds
		totalTwo = totalTwo + inform[0].Bookmarkers[i].Markets[0].Outcomes[1].Odds
	}

	averageOne := totalOne / float64(len(inform[0].Bookmarkers))
	teamOne := inform[0].Bookmarkers[0].Markets[0].Outcomes[0].Team

	averageTwo := totalTwo / float64(len(inform[0].Bookmarkers))
	teamTwo := inform[0].Bookmarkers[0].Markets[0].Outcomes[1].Team
	fmt.Printf("Avg Odds for %s: %f\n", teamOne, averageOne)
	fmt.Printf("Avg Odds for %s: %f\n", teamTwo, averageTwo)
}

func sendMessages(s *discordgo.Session) {
	//8am caculations
	sum := 1
	num := 1
	for sum < 100 {
		fmt.Println(num)
		time.Sleep(2 * time.Second)
		//s.ChannelMessageSend(chatString, "Rise and shine!")
		num++
	}
}

func scrapSetup(past time.Duration) (string, *colly.Collector) {
	todaysDate := time.Now()
	yesterday := todaysDate.Add(past * time.Hour)
	month := fmt.Sprintf("%02d", int(yesterday.Month()))
	day := strconv.Itoa(yesterday.Day())
	year := strconv.Itoa(yesterday.Year())
	scrapURL := "https://www.cbssports.com/nba/schedule/" + year + month + day
	c := colly.NewCollector(colly.AllowedDomains("www.cbssports.com", "cbssports.com"))
	return scrapURL, c
}

func getTodaysGames(s *discordgo.Session) {
	scrapURL, c := scrapSetup(0)
	list := []string{}
	postLine := ""

	//Before HTTP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnHTML("table.TableBase-table tbody tr.TableBase-bodyTr td.TableBase-bodyTd span.CellLogoNameLockup", func(h *colly.HTMLElement) {
		selection := h.DOM
		list = append(list, selection.Find("a").Text())
	})

	//Error occured during request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error")
	})

	c.Visit(scrapURL)

	// //Recieve response
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Page visited: ", r.Request.URL)
	// })

	// //After OnResponse if content is HTML
	// c.OnHTML("li.product", func(e *colly.HTMLElement) {
	// 	// printing all URLs associated with the a links in the page
	// 	fmt.Println("%v", e.Attr("href"))
	// })

	//After all OnHMTL
	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println(listLosers)
	// 	fmt.Println(listWinners)
	// })
	postLine += "DATE: " + time.Now().Format("2006-01-02") + "\n"
	if len(list) == 0 {
		postLine += "No Games Today"
	} else {
		for index, value := range list {
			if index%2 == 0 {
				postLine += "AWAY: " + value + "  VS  "
			} else {
				postLine += " HOME: " + value + "\n"
			}
		}
	}
	s.ChannelMessageSend(chatString, postLine)
}

func getPastGames(s *discordgo.Session) {
	scrapURL, c := scrapSetup(-24)
	// listAway := []string{}
	// listHome := []string{}
	list := []string{}
	postLine := ""
	//Before HTTP request
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnHTML("table.TableBase-table tbody tr.TableBase-bodyTr td.TableBase-bodyTd div.CellGame", func(h *colly.HTMLElement) {
		selection := h.DOM
		list = append(list, selection.Find("a").Text())
	})

	//Error occured during request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error")
	})

	c.Visit(scrapURL)

	// //Recieve response
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Page visited: ", r.Request.URL)
	// })

	// //After OnResponse if content is HTML
	// c.OnHTML("li.product", func(e *colly.HTMLElement) {
	// 	// printing all URLs associated with the a links in the page
	// 	fmt.Println("%v", e.Attr("href"))
	// })

	//After all OnHMTL
	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println(listLosers)
	// 	fmt.Println(listWinners)
	// })
	postLine += "DATE: " + time.Now().Format("2006-01-02") + "\n"
	if len(list) == 0 {
		postLine += "No Games Yesterday"
	} else {
		for _, value := range list {
			postLine += "RESULTS: " + value + "\n"
		}
	}
	s.ChannelMessageSend(chatString, postLine)
}

func betNBAGame(id string, input string, s *discordgo.Session) {
	currBalance := getBalance(id, s)
	amount := ""
	team := ""
	if strings.HasPrefix(input, "!bet ") {
		betStr := strings.TrimPrefix(input, "!bet ")
		betList := strings.Split(betStr, " ")
		if len(betList) != 2 {
			fmt.Println("Error, not enough parameters")
		} else {
			for _, value := range betList {
				_, err := strconv.Atoi(value)
				if err != nil {
					team = matchNBATeam(value)
				} else {
					amount = value
				}
			}
		}
		if amount == "" || team == "" {
			fmt.Println("Error, wrong parameters")
			return
		}
	}
	numAmount, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Println("Error, not a number!")
		return
	}
	//Case if amount bet is over actual amount (Negative)
	if (currBalance - numAmount) < 0 {
		fmt.Println("Error, too high of a bet")
		return
	}
	query := "UPDATE balance SET balance = balance - ? where ProfileID = ?;"
	_, err = db.Exec(query, amount, id) //Look into context timeout
	if err != nil {
		//Case where no matching rows
		if strings.Contains(err.Error(), "no rows") {
			// Handle the permission-denied error
			s.ChannelMessageSend(chatString, "No matching rows! Hmph~!")
		} else {
			panic(err.Error())
		}
	} else {
		//Update transaction log
		uuid := strings.Replace(uuid.New().String(), "-", "", -1)
		balanceId := getBalanceId(id, s)
		typeOfBet := "ML"
		comment := "Testing"
		currDate := time.Now().Format("2006-01-02 15:04:05")
		query = "INSERT INTO log_book VALUES(?, ?, ?, ?, ?, ?);"
		_, err = db.Exec(query, uuid, numAmount, typeOfBet, comment, currDate, balanceId) //Look into context timeout
		if err != nil {
			//Case where duplicate entry
			if strings.Contains(err.Error(), "Duplicate entry") {
				// Handle the permission-denied error
				s.ChannelMessageSend(chatString, "Already registered balance transaction! Hmph~!")
			} else {
				panic(err.Error())
			}
		} else {
			fmt.Println("Bet " + amount + " on " + team + " with XX odds")
			s.ChannelMessageSend(chatString, "Updated Balance to "+strconv.Itoa(currBalance-numAmount)+"\n")
		}
	}
}

func matchNBATeam(input string) string {
	input = strings.ToLower(input)
	switch input {
	case "atlanta hawks", "hawks", "atl":
		return "Atlanta Hawks"
	case "boston celtics", "celtics", "bos", "boston":
		return "Boston Celtics"
	case "brooklyn nets", "nets", "bkn":
		return "Brooklyn Nets"
	case "charlotte hornets", "hornets", "cha":
		return "Charlotte Hornets"
	case "chicago bulls", "bulls", "chi":
		return "Chicago Bulls"
	case "cleveland cavaliers", "cavaliers", "cle", "clev":
		return "Cleveland Cavaliers"
	case "dallas mavericks", "mavericks", "dal", "dallas":
		return "Dallas Mavericks"
	case "denver nuggets", "nuggets", "den":
		return "Denver Nuggets"
	case "detroit pistons", "pistons", "det":
		return "Detroit Pistons"
	case "golden state warriors", "warriors", "gs", "gsw":
		return "Golden State Warriors"
	case "houston rockets", "rockets", "hou":
		return "Houston Rockets"
	case "indiana pacers", "pacers", "ind":
		return "Indiana Pacers"
	case "los angeles clippers", "clippers", "lac", "clip":
		return "Los Angeles Clippers"
	case "los angeles lakers", "lakers", "lal":
		return "Los Angeles Lakers"
	case "memphis grizzlies", "grizzlies", "mem", "grizz":
		return "Memphis Grizzlies"
	case "miami heat", "heat", "mia", "miami":
		return "Miami Heat"
	case "milwaukee bucks", "bucks", "mil":
		return "Milwaukee Bucks"
	case "minnesota timberwolves", "timberwolves", "min", "wolves":
		return "Minnesota Timberwolves"
	case "new orleans pelicans", "pelicans", "nop":
		return "New Orleans Pelicans"
	case "new york knicks", "knicks", "nyk", "ny":
		return "New York Knicks"
	case "oklahoma city thunder", "thunder", "okc":
		return "Oklahoma City Thunder"
	case "orlando magic", "magic", "orl":
		return "Orlando Magic"
	case "philadelphia 76ers", "76ers", "phi":
		return "Philadelphia 76ers"
	case "phoenix suns", "suns", "phx":
		return "Phoenix Suns"
	case "portland trail blazers", "trail blazers", "blazers", "por":
		return "Portland Trail Blazers"
	case "sacramento kings", "kings", "sac":
		return "Sacramento Kings"
	case "san antonio spurs", "spurs", "sas":
		return "San Antonio Spurs"
	case "toronto raptors", "raptors", "tor":
		return "Toronto Raptors"
	case "utah jazz", "jazz", "uta", "utah":
		return "Utah Jazz"
	case "washington wizards", "wizards", "was":
		return "Washington Wizards"
	default:
		fmt.Println("Error, can't find team")
		return ""
	}
}

func convertTime(s string) bool {
	// Split the time string into minutes and seconds
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		fmt.Println("Invalid time format")
		return false
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("Error converting minutes:", err)
		return false
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("Error converting seconds:", err)
		return false
	}

	// Convert to total seconds
	totalSeconds := (minutes * 60) + seconds
	if totalSeconds < 1800 {
		return true
	}
	return false
}

func elementExists(list []string, target string) bool {
	for _, item := range list {
		println(item + " " + target)
		if strings.EqualFold(item, target) {
			return true
		}
	}
	return false
}
