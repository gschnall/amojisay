package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	amoji_map "github.com/gschnall/amojisay/amoji_map"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/urfave/cli"
)

var (
	columns *int32
)

var app = cli.App{
	Action: func(c *cli.Context) error {
		return nil
	},
}

func getAmojiFromMap(amojiName string) interface{} {
	return amoji_map.Amojis[amojiName]
}

func listAllAmojisInJSONFile() {
	content, err := ioutil.ReadFile("./amojis.json")
	if err != nil {
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	amojis := []string{}
	for key := range payload {
		amojis = append(amojis, key)
	}

	sort.Strings(amojis)
	for i := 0; i < len(amojis); i++ {
		fmt.Print(amojis[i] + " ")
	}
}

func getAmojiNameFromTemplate(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func searchAndPrintSimiliarAmojis(search string, payload map[string]interface{}) {
	if len(search) < 2 {
		fmt.Print("amoji " + search + " not found.\namojisay -l |> list all available amojis\n")
	} else {
		amojis := []string{}
		for key := range payload {
			similarity := strutil.Similarity(key, search, metrics.NewLevenshtein())

			if similarity >= .75 || (len(search) > 2 && strings.HasPrefix(key, search) || len(search) > 3 && strings.Contains(key, search)) {
				amojis = append(amojis, key)
			}
		}

		if len(amojis) > 0 {
			sort.Strings(amojis)
			fmt.Print("¯\\_(ツ)_/¯ couldn't find amoji " + search + "\n\nMaybe you meant to use one of these:\n")
			for i := 0; i < len(amojis); i++ {
				fmt.Print(amojis[i] + " ")
			}
		} else {
			fmt.Print("(╥﹏╥) amoji " + search + " not found.\namojisay -l |> list all available amojis\n")
		}
	}
}

func setupCliApp() {
	app.Name = "amojisay"
	app.Usage = "One line ascii emoji cli"
	app.Author = "Gabe Schnall"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{}
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "a",
			Usage:    "amoji name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "s",
			Usage:    "formatted string with amojis. Example: \"-s %{tada} Adverture Time with %{finn} & %{jake}\"",
			Required: true,
		},
		&cli.BoolFlag{
			Name: "p", Usage: "prepend text",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "l",
			Usage:    "list all amojis",
			Required: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		inputs := os.Args

		message := inputs[len(inputs)-1]
		// handle pipes
		// r := bufio.NewReader(os.Stdin)
		// buf := make([]byte, 0, 4*1024)
		// n, _ := r.Read(buf[:cap(buf)])
		// buf = buf[:n]
		// pipedString := string(buf)
		// if pipedString != "" {
		// 	message = pipedString
		// }

		if c.Bool("l") {
			listAllAmojisInJSONFile()
		} else if c.String("s") != "" {
			var re = regexp.MustCompile(`%{[^{}]*}`)
			matches := re.FindAllStringSubmatch(c.String("s"), -1)

			amojiString := c.String("s")
			for _, v := range matches {
				amojiVar := v[0]
				amojiName := strings.ToLower(amojiVar[2 : len(amojiVar)-1])
				amoji := amoji_map.Amojis[amojiName]
				if amoji == nil {
					searchAndPrintSimiliarAmojis(amojiName, amoji_map.Amojis)
					// fmt.Printf("amoji [ " + amojiName + " ] not found.\namojisay -l |> list all available amojis\n")
					os.Exit(0)
				} else {
					amojiString = strings.Replace(amojiString, amojiVar, fmt.Sprintf("%v", amoji), -1)
				}
			}
			fmt.Print(amojiString + "\n")
		} else if c.String("a") != "" {
			amojiName := strings.ToLower(c.String("a"))
			amoji := getAmojiFromMap(amojiName)

			if amoji == nil {
				searchAndPrintSimiliarAmojis(amojiName, amoji_map.Amojis)
				// fmt.Printf("amoji [ " + amojiName + " ] not found.\namojisay -l |> list all available amojis\n")
				os.Exit(0)
			} else if c.Bool("p") {
				fmt.Printf("%s %s\n", message, amoji)
			} else {
				fmt.Printf("%s %s\n", amoji, message)
			}
		}
		return nil
	}
}

func main() {
	setupCliApp()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
