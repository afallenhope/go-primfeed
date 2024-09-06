package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	primfeed "github.com/afallenhope/primfeed/pkg"
)

// Rather than using 3rd party libs
// we open the file and then parse it ourselves
//
// returns an error if you can't load the file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func main() {
	err := loadEnvFile(".env")

	if err != nil {
		fmt.Println("Error, no .env file. Make sure to rename the .env.example and place your token in it.", err)
		return
	}

	pf := primfeed.NewPrimfeed("https://api.primfeed.com/pf")

	// If you don't have a token, provide username and password
	_, err = pf.Login(os.Getenv("PRIMFEED_USERNAME"), os.Getenv("PRIMFEED_PASSWORD"), nil)

	// Uncomment the following if you want to use inworld code...
	// codeResp, err := pf.GetLoginCode(os.Getenv("PRIMFEED_USERNAME"))
	// if err != nil {
	// 	fmt.Println("could not get login code.", err)
	// 	return
	// }
	//
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Login Code: ")
	// text, err := reader.ReadString('\n')
	// if err != nil {
	// 	fmt.Println("could not get code. %v", err)
	// 	return
	// }
	//
	// _, err = pf.LoginWithCode(codeResp.RequestID, text, "")

	// If you have a token, no need to supply username and password.
	// Uncomment below
	// pf.SetToken(os.Getenv("PRIMFEED_TOKEN"))

	if err != nil {
		fmt.Printf("Could not login: %v", err)
		return
	}

	// fmt.Printf("Setting token: %v", response.Token)

	var names []string

	err = pf.GetMe()
	if err != nil {
		fmt.Printf("Could not get profile: %v", err)
		return
	}

	for _, f := range pf.Me.Followers {
		names = append(names, f.Handle)
	}

	slices.Sort(names)
	fmt.Printf("Followers: %s\n\n", strings.Join(names, "\n"))
	fmt.Printf("Total followers: %d\n\n", len(names))

	err = pf.GetMe()
	if err != nil {
		fmt.Printf("error getting profile: %v", err)
		return
	}

	fmt.Printf("Your Name: %s\n", pf.Me.Profile.User.Name)
	fmt.Printf("You have %d followers, and are following %d\n\n",
		len(pf.Me.Followers), len(pf.Me.Following))

	count, err := pf.GetNotificationCount()
	if err != nil {
		fmt.Printf("error getting notification count: %v", err)
		return
	}

	fmt.Printf("Notification Count: %d\n", count)
	if count <= 0 {
		fmt.Println("You have no new notifications")
	} else {
		fmt.Printf("You have %d new notifications\n", count)
		// notifications, err := pf.GetNotifications()
		// if err != nil {
		// 	fmt.Println("could not get notifications", err)
		// 	return
		// }

		// fmt.Printf("Notifications: %v\n", notifications)
	}

}
