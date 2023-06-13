package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Read enviroment var of .env
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	service := os.Getenv("SERVICE")
	//Exec command hostname -f to obtain the hostname of server
	hostname, _ := exec.Command("hostname", "-f").Output()

	//Obtain the public ip address of server
	web := randomWebCheckIpAddress()
	ipAddress, _ := exec.Command("curl", web).Output()

	//If ipAddress is 0 it means that it did not receive the IP
	if len(ipAddress) == 0 {
		fmt.Printf("[PROBLEM] at obtain IP from %s\n\n", web)
		for len(ipAddress) == 0 {
			web := randomWebCheckIpAddress()
			newIpAddress, _ := exec.Command("curl", web).Output()
			ipAddress = newIpAddress
			if len(ipAddress) > 0 {
				break
			}
		}
	}
	var updateDomain string

	if (strings.ToUpper(service)) == "GOOGLEDOMAIN" {
		updateDomain = fmt.Sprintf("https://%s:%s@domains.google.com/nic/update?hostname=%s&myip=%s", username, password, []byte(strings.TrimSpace(string(hostname))), []byte(strings.TrimSpace(string(ipAddress))))
		fmt.Println(updateDomain)
	} else if strings.ToUpper(service) == "CLOUDFLARE" {
		updateDomain = fmt.Sprintf("https://%s:%s@domains.google.com/nic/update?hostname=%s&myip=%s", username, password, []byte(strings.TrimSpace(string(hostname))), []byte(strings.TrimSpace(string(ipAddress))))
		fmt.Println(updateDomain)
	} else {
		fmt.Printf("[PROBLEM] at .env SERVICE: %s\n\n", service)
	}
	updatingDomain := exec.Command("curl", updateDomain)
	errDomain := updatingDomain.Run()

	if errDomain != nil {
		log.Fatal("[PROBLEM] at update domain")
	}

	//to do, get public ip address from other service
	//create log file for check if web is down or not work?
	//update the dynamic ip address to google domain service with curl or other?
	//test in production?
}

// Get random web to check the public ip address from server
func randomWebCheckIpAddress() string {
	rand.NewSource(time.Now().UnixNano())
	web := []string{"ifconfig.co", "ifconfig.me", "icanhazip.com"}
	return web[rand.Intn(len(web))]
}
