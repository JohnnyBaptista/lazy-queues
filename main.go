package main

import (
	"lazy-queues/client"
	"lazy-queues/monitoring"
) 


func main() {
	client.Execute()	
	monitoring.FetchMetricDescriptors()
}


