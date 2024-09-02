package seatsmanage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestFindAvailableSeat(t *testing.T) {
	seats := []Seats{}
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 1, SeatNumber: 1},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 2, SeatNumber: 2},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 3, SeatNumber: 3},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 4, SeatNumber: 4},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 5, SeatNumber: 5},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 6, SeatNumber: 6},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 7, SeatNumber: 7},
	// 	{SeatStatus: false, Section: "A", Price: 50.0, ID: 8, SeatNumber: 8},
	// 	{SeatStatus: true, Section: "A", Price: 50.0, ID: 9, SeatNumber: 9},
	// }
	for i := 0; i < 1000000; i++ {
		rand.Seed(time.Now().UnixNano())
		seats = append(seats, Seats{SeatStatus: rand.Intn(2) == 1, Section: "A", Price: 50.0, ID: int32(i), SeatNumber: int32(i)})
	}
	start := time.Now()
	data, bo := FindAvailableSeat(seats, 2)
	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
	fmt.Println(data)
	fmt.Println(bo)
}
