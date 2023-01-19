# Go github exercise 

This repository contains an intelligence service for public Github events.
The service periodically watches public Github events as they happen, through the relevant github events API, and process them in order to allow the functionality that's described below.
The service also contains an server for API requests

## Clone the project

```
$ git clone https://github.com/Yoni7/GoGithubEvents
$ cd GoGithubEvents
```
https://github.com/Yoni7/GoGithubEvents is the canonical Git repository.

## Service description

The service includes of 2 main parts:
1. `github_events.go`
Periodically watches public github event via an API call to `https://api.github.com/events`. 
    The required data (event types, actors, repo urls, unique emails) are stored in MongoDB.
    The Entry point to this section is `GetPublicEvents`

2. `server.go`
Webserver the listens on port `8080` and expose 4 routes to retireve the information about the github global events from MongoDB via GET method:
    - `/github/event`: All the event types names and how many time have we seen them
    - `/github/actors`: Unique name of the last 50 actors.
    - `/github/repos`: The last 20 repository URLs and there stars counter
    - `/github/emails`: All of the unique email addresses found in github global events

The MongoDB is hosted at MongoDB Atlas

## Docker
This project can also be used with a docker image. Check out `Dockerfile`.

## Notes
The env varialbes fallback in the code is only for local debugging purposes (in a real production environment they must come has ENV)
