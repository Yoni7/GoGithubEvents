package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sort"
	"regexp"
)


func GetPublicEvents() {
	body, err := getDataFromGithubApi("https://api.github.com/events?per_page=200")
	if err != nil {
		fmt.Println(err)
		return
	}

	parseGlobalEvents(body)
}

func getDataFromGithubApi(url string) ([]byte, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	var b []byte
	if err != nil {
		fmt.Println(err)
		return b, err
	}

	req.Header.Add("Authorization", "Bearer " + GetEnv("GITHUB_TOKEN", "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"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return b, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return b, err
	}

	return body, nil
}


func parseGlobalEvents(body []byte) {
	var eventsObjects []map[string]interface{}
	json.Unmarshal([]byte(string(body)), &eventsObjects)

	for _, eventObj := range eventsObjects {
		eventType := eventObj["type"].(string)
		go UpdateEventType(eventType)

		actorObj := eventObj["actor"]
		actorObjData := actorObj.(map[string]interface {})
		actorLogin := actorObjData["login"].(string)
		go UpdateEventActor(actorLogin)

		repoObj := eventObj["repo"]
		repoObjData := repoObj.(map[string]interface {})
		repoUrl := repoObjData["url"].(string)
		repoName := repoObjData["name"].(string)
		go UpdateEventRepo(repoName, repoUrl)

		lookForEmailsInEvent(eventObj)
	}
}

func lookForEmailsInEvent(eventObj map[string]interface {}) {
	for key, value := range eventObj {
		if key == "email" {
			emailAddr := value.(string)
			if isEmailValid(emailAddr) {
				go UpdateEventEmail(emailAddr)
			}
       	} else if v, ok := value.(map[string]any); ok {
            lookForEmailsInEvent(v)
        } else if v, ok := value.([]any); ok {
			for _, valueObj := range v {
				if v, ok := valueObj.(map[string]any); ok {
					lookForEmailsInEvent(v)
				}
			}
        }
	}
}

func isEmailValid(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(e)
}


func GetEventTypes(w http.ResponseWriter, r *http.Request) {
	events := GetEventsDocs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func GetActors(w http.ResponseWriter, r *http.Request) {
	actors := GetActorsDocs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors.Actors)
}

func GetEmails(w http.ResponseWriter, r *http.Request) {
	emails := GetEmailsDocs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}


	
type repoStars struct {
    Name 	string
    Url  	string
	Stars	int
}
var STARS_WG sync.WaitGroup
func GetRepoUrls(w http.ResponseWriter, r *http.Request) {
	repos := GetRepoDocs()

	repoItems := []repoStars{}
	for _, repoObj := range repos.Repos {
		STARS_WG.Add(1)
		repoStarObj := repoStars{Name: repoObj.Name, Url: repoObj.Url, Stars: 0}
		resultStars := make(chan int)
		go getRepoStarsFromGithub(repoObj.Name, resultStars)
		repoStarObj.Stars = <-resultStars
		repoItems = append(repoItems, repoStarObj)
	}

	STARS_WG.Wait()
	fmt.Println("finished getting stars for all repos")

	// sort from high to low (stars)
	sort.Slice(repoItems[:], func(i, j int) bool {
		return repoItems[i].Stars > repoItems[j].Stars
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repoItems)
}


func getRepoStarsFromGithub(repoName string, resultStars chan int){
	defer STARS_WG.Done()

	if repoName == "" {
		fmt.Println("Error: missing repo name")
		return
	}

	var url = "https://api.github.com/repos/"+repoName

	body, err := getDataFromGithubApi(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	var repoInfo map[string]interface{}
	json.Unmarshal([]byte(string(body)), &repoInfo)

	if repoInfo["stargazers_count"] == nil {
		resultStars <- 0
		return
	}

	stars := int(repoInfo["stargazers_count"].(float64))
	fmt.Printf("getRepoStarsFromGithub: stars %v for %v\n", stars, url)
	resultStars <- stars
}
