package main

import (
	"net/http"
	"refresh/lib/database"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
)

// func pingAddress(address string) {
// 	pinger, err := ping.NewPinger(address)
// 	if err != nil {
// 		panic(err)
// 	}
// 	pinger.Count = 3
// 	pinger.Timeout = time.Second * 3
// 	pinger.SetPrivileged(true)

// 	pinger.OnFinish = func(stats *ping.Statistics) {
// 		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
// 		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
// 			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
// 		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
// 			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
// 		ping := database.Ping{Loss: stats.PacketLoss, Response: int64(stats.AvgRtt), Time: time.Now().UTC()}
// 		database.UpdateDatabase(address, ping)
// 	}

// 	err = pinger.Run()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func pingAsyncContinuous() {
// 	for _, address := range database.Websites {
// 		go pingAddress(address)
// 	}
// 	for range time.Tick(time.Minute * 60) {
// 		for _, address := range database.Websites {
// 			go pingAddress(address)
// 		}
// 	}
// }

// func pingAsync() {
// 	var wg sync.WaitGroup
// 	for _, address := range database.Websites {
// 		wg.Add(1)
// 		go func(address string) {
// 			defer wg.Done()
// 			pingAddress(address)
// 		}(address)
// 	}
// 	wg.Wait()
// }

// func pingSync() {
// 	for _, address := range database.Websites {
// 		pingAddress(address)
// 	}
// }

func fetchHttpAddress(collection *mongo.Collection, address string) {
	// Measure time to perform a HTTP GET request on the address
	// HTTP alternative to pingAddress
	httpAddress := "http://" + address
	startTime := time.Now()
	_, err := http.Get(httpAddress)
	if err != nil {
		panic(err)
	}
	duration := time.Since(startTime)
	print(duration)

	ping := database.Ping{Loss: 0, Response: int64(duration), Time: time.Now().UTC()}

	database.UpdateDatabase(collection, address, ping)
}

// func httpAsyncContinuous() {
// 	for _, address := range database.Websites {
// 		go fetchHttpAddress(address)
// 	}
// 	for range time.Tick(time.Minute * 60) {
// 		for _, address := range database.Websites {
// 			go fetchHttpAddress(address)
// 		}
// 	}
// }

// func httpAsync() {
// 	var wg sync.WaitGroup
// 	for _, address := range database.Websites {
// 		wg.Add(1)
// 		go func(address string) {
// 			defer wg.Done()
// 			fetchHttpAddress(address)
// 		}(address)
// 	}
// 	wg.Wait()
// }

func httpSync() {
	collection := database.ConnectToDatabase()
	for _, address := range database.Websites {
		fetchHttpAddress(collection, address)
	}
}

func LambdaHandler() error {
	httpSync()
	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
