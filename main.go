package main

func main() {
	// Mongo will store all the github events data
	ConnectToMongoDB()

	// Periodically watch public Github events
	go GetPublicEventsPeriodically()

	// starting server for API requests
	RunServer()
}
