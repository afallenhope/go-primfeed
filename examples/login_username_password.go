package examples

import (
	"log"
	"os"

	primfeed "github.com/afallenhope/primfeed/pkg"
)

func LoginWithUsernameAndPassword() {
	loadEnvFile(".env")
	pf := primfeed.NewPrimfeed("api.primfeed.com")
	pf.Login(os.Getenv("PRIMFEED_USERNAME"), os.Getenv("PRIMFEED_PASSWORD"), nil)

	err := pf.GetMe()
	if err != nil {
		log.Fatalf("could not get profile. %v", err)
		return
	}

	log.Printf("Your name: %s\n", pf.Me.Profile.User.Name)
	log.Printf("Followers %d\n", len(pf.Me.Followers))
	log.Printf("Following %d\n", len(pf.Me.Following))
}
