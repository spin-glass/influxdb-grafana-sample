package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func stringToTime(DateTime string) time.Time {
	date, _ := time.Parse("01/02/2006 03:04:00 PM", DateTime)
	return date
}

func main() {
	bucket_name := "bike"

	client := influxdb2.NewClient("http://localhost:8086", "token")
	writeAPI := client.WriteAPIBlocking("org", bucket_name)

	f, err := os.Open("../data/Eco-Totem_Broadway_Bicycle_Count.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	header, _ := r.Read()
	fmt.Println(header)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// DateTime, Day, Date, Time, Total, Westbound, Eastbound := record[0], record[1], record[2], record[3], record[4], record[5], record[6]
		DateTime, Day := record[0], record[1]
		Westbound, _ := strconv.Atoi(record[5])
		Eastbound, _ := strconv.Atoi(record[6])

		datetime := stringToTime(DateTime)

		println(DateTime, Westbound, Eastbound)

		west := influxdb2.NewPointWithMeasurement("bike").
			AddTag("WeekDay", Day).
			AddTag("Direction", "west").
			AddField("Amount", Westbound).
			SetTime(datetime)

		east := influxdb2.NewPointWithMeasurement("bike").
			AddTag("WeekDay", Day).
			AddTag("Direction", "east").
			AddField("Amount", Eastbound).
			SetTime(datetime)

		writeAPI.WritePoint(context.Background(), west)
		writeAPI.WritePoint(context.Background(), east)
	}

	queryAPI := client.QueryAPI("org")

	result, err := queryAPI.Query(context.Background(), `from(bucket:bike)|> range(start: -1y) |> filter(fn: (r) => r._measurement == "bike")`)

	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}

		client.Close()
	}
}
