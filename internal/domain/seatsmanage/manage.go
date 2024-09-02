package seatsmanage

import (
	"sort"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
)

type Seats struct {
	SeatStatus bool
	Section    string
	Price      float64
	ID         int32
	SeatNumber int32
	Types      string
}

func FindAvailableSeat(seats []Seats, amount int) ([]Seats, bool) {
	// 排序座位按SeatNumber
	sort.Slice(seats, func(i, j int) bool {
		return seats[i].SeatNumber < seats[j].SeatNumber
	})

	// 查找連續座位
	for i := 0; i <= len(seats)-amount; i++ {
		if !seats[i].SeatStatus {
			continue
		}
		available := true
		for j := 0; j < amount; j++ {
			if !seats[i+j].SeatStatus || seats[i+j].SeatNumber != seats[i].SeatNumber+int32(j) {
				available = false
				break
			}
		}
		if available {
			return seats[i : i+amount], true
		}
	}

	// 如果找不到連續的，找任意可用座位
	result := []Seats{}
	for _, s := range seats {
		if s.SeatStatus {
			result = append(result, s)
			if len(result) == amount {
				return result, true
			}
		}
	}

	return result, len(result) == amount
}

type BySeatNumber []PostgresDB.Seat

func (a BySeatNumber) Len() int           { return len(a) }
func (a BySeatNumber) Less(i, j int) bool { return a[i].SeatNumber < a[j].SeatNumber }
func (a BySeatNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func SortSeatsBySectionAndSeatNumber(seats []PostgresDB.Seat) []PostgresDB.Seat {
	// Group seats by section
	sections := make(map[string][]PostgresDB.Seat)
	sectionOrder := []string{}

	for _, seat := range seats {
		if _, found := sections[seat.Section]; !found {
			sectionOrder = append(sectionOrder, seat.Section)
		}
		sections[seat.Section] = append(sections[seat.Section], seat)
	}

	// Sort each section by SeatNumber
	sortedSeats := []PostgresDB.Seat{}
	for _, section := range sectionOrder {
		sectionSeats := sections[section]
		sort.Sort(BySeatNumber(sectionSeats))
		sortedSeats = append(sortedSeats, sectionSeats...)
	}

	return sortedSeats
}

type Section struct {
	Section string
	Price   float64
	Count   int32
}

func ArrageSectionNumbersPrices(sectionName []string, prices []float64, count map[string]int32) []Section {
	section := []Section{}

	for i := range sectionName {
		section = append(section, Section{
			Section: sectionName[i],
			Price:   prices[i],
			Count:   count[sectionName[i]],
		})
	}

	return section
}
