package main

//Import Packages
import (
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

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/google/uuid"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token        string `json:"Token"`
	SQLPassword  string `json:"SQLPassword"`
	RapidAPIKey  string `json:"X-RapidAPI-Key"`
	RapidAPIHost string `json:"X-RapidAPI-Host"`
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

const chatString = "795763169191788585"

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
	//If session ID matches dont do anything
	if m.Author.ID == s.State.User.ID {
		return
	}
	userInput := m.Content
	//Parse command
	if strings.HasPrefix(userInput, "!") {
		switch userInput {
		case "!ping":
			s.ChannelMessageSend(m.ChannelID, "pong") //Check if someone types hello in chat
		case "!profile": //Command for Profile info (Currency balance)
			{
				balance := getBalance(m.Author.ID, s)
				s.ChannelMessageSend(chatString, "User: "+m.Author.Username+"\nBalance: "+strconv.Itoa(balance))
			}
		// case "!odds":
		// 	GetOdds(s)
		case "!register":
			registerUser(m.Author.ID, m.Author.Username, s)
		case "!update":
			updateBalance(m.Author.ID, s)
		case "!scores":
			getPastGames(s)
		case "!games":
			getTodaysGames(s)
		default: //!bet, ....
			{
				betNBAGame(userInput)
				return
			}
		}
	}

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

func betNBAGame(input string) {
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
		}
	}

	//Subtract from balance
	//Update transaction log
	//At 5am EST update for NBA (LCS will be different)
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
	case "miami heat", "heat", "mia":
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
