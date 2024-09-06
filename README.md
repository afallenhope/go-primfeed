# Primfeed

## Obligatory Disclaimer

Primfeed is a social media platform created by Luke Rowley
This repository was made for educational purposes and as a hobby.

## Origin

Initially I was working on a project and I had to refresh my browser over and over, this wasn't ideal for me and though; "Hey why not have a terminal based client?"
I started logging where the endpoints went and slowly started to begin.

## Usage and Examples

The repository has two examples.
Login using your token

```go
package example

func ExampleLoginWithToken() {
    pf := primfeed.NewPrimfeed("https://api.primfeed.com")
    pf.SetToken(os.GetEnv("PRIMFFEED_TOKEN"))
    err := pf.GetMe()

    if err != nil {
        log.Fatalf("could not get personal profile.: %v" err)
        return
    }

    count, err := pf.GetNotificationCount()
    if err != nil {
        log.Fatalf("could not get notification count: %v", err)
        return
    }

	if count >= 0 {
		fmt.Println("You have no new notifications")
	} else {
		fmt.Printf("You have %d new notifications\n", count)
		notifications, err := pf.GetNotifications()
		if err != nil {
			fmt.Println("could not get notifications", err)
			return
		}

		fmt.Printf("Notifications: %v\n", notifications)
	}
}
```

Login using your username + password

```go
package example

func ExampleLoginWithToken() {
    pf := primfeed.NewPrimfeed("https://api.primfeed.com")
    pf.Login(os.GetEnv("PRIMFFEED_USERNAME"), os.GetEnv("PRIMFEED_PASSWORD"))
    err := pf.GetMe()

    if err != nil {
        log.Fatalf("could not get personal profile.: %v" err)
        return
    }

    count, err := pf.GetNotificationCount()
    if err != nil {
        log.Fatalf("could not get notification count: %v", err)
        return
    }

	if count >= 0 {
		fmt.Println("You have no new notifications")
	} else {
		fmt.Printf("You have %d new notifications\n", count)
		notifications, err := pf.GetNotifications()
		if err != nil {
			fmt.Println("could not get notifications", err)
			return
		}

		fmt.Printf("Notifications: %v\n", notifications)
	}
}
```


