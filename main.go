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

func getAmojiJSONFile() map[string]interface{} {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile("./amojis.json")
	if err != nil {
		log.Fatal("Error could not ReadFile ./amojis.json): ", err)
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

func getAmojiFromJSONFile(amojiName string) interface{} {
	payload := getAmojiJSONFile()
	return payload[amojiName]
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
			amojiJSONFile := getAmojiJSONFile()
			var re = regexp.MustCompile(`%{[^{}]*}`)
			matches := re.FindAllStringSubmatch(c.String("s"), -1)

			amojiString := c.String("s")
			for _, v := range matches {
				amojiVar := v[0]
				amojiName := amojiVar[2 : len(amojiVar)-1]
				amoji := amojiJSONFile[amojiName]
				if amoji == nil {
					fmt.Printf("amoji [ " + amojiName + " ] not found.\namoji -l |> list all available amojis\n")
					os.Exit(0)
				} else {
					amojiString = strings.Replace(amojiString, amojiVar, fmt.Sprintf("%v", amoji), -1)
				}
			}
			fmt.Print(amojiString + "\n")
		} else if c.String("a") != "" {
			amojiName := c.String("a")
			amoji := getAmojiFromJSONFile(amojiName)

			if amoji == nil {
				fmt.Printf("amoji [ " + amojiName + " ] not found.\namoji -l |> list all available amojis\n")
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
