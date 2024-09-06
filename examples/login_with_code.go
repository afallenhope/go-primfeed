package examples

import (
	"bufio"
	"fmt"
	"log"
	"os"

	primfeed "github.com/afallenhope/primfeed/pkg"
)

func LoginWithCode() {

	loadEnvFile(".env")
	pf := primfeed.NewPrimfeed("https://api.primfeed.com")

	codeResp, err := pf.GetLoginCode(os.Getenv("PRIMFEED_USERNAME"))
	if err != nil {
		log.Fatalf("could not get login code. %v", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Login Code: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("could not get code. %v", err)
		return
	}

	_, err = pf.LoginWithCode(codeResp.RequestID, text, "")
	if err != nil {
		log.Fatalf("could not login. %v", err)
		return
	}

	err = pf.GetMe()
	if err != nil {
		log.Fatalf("could not get profile. %v", err)
		return
	}

	log.Printf("Your name: %s\n", pf.Me.Profile.User.Name)
	log.Printf("Followers %d\n", len(pf.Me.Followers))
	log.Printf("Following %d\n", len(pf.Me.Following))
}

func main() {
	LoginWithCode()
}
