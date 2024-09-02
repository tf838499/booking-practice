package RedisCache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
)

type SeatKeyParams struct {
	EventID    int32
	Section    string
	SeatNumber []int32
}

func (q *RedisRepository) SetSeatKey(ctx context.Context, arg SeatKeyParams) (bool, error) {

	SetKey := strconv.FormatInt(int64(arg.EventID), 10) + ":" + arg.Section + ":Key"
	// data := make(map[int32]interface{})
	flag := false
	SetData := []int32{}
	for i := range arg.SeatNumber {
		field := strconv.FormatInt(int64(arg.SeatNumber[i]), 10)
		data := q.Client.HSetNX(ctx, SetKey, field, 1)
		if data.Err() != nil {
			flag = true
		}
		if !data.Val() {
			flag = true
		} else {
			SetData = append(SetData, arg.SeatNumber[i])
		}
	}
	if flag {
		for j := range SetData {
			field := strconv.FormatInt(int64(SetData[j]), 10)
			err := q.Client.HDel(ctx, SetKey, field).Err()
			if err != nil {
				return false, err
			}
		}
		return false, errors.New("seat get fail")
	}
	// q.Client.HGetAll()
	return true, nil
}
func (q *RedisRepository) DelSeatKey(ctx context.Context, arg SeatKeyParams) error {

	SetKey := strconv.FormatInt(int64(arg.EventID), 10) + ":" + arg.Section + ":Key"

	for i := range arg.SeatNumber {
		field := strconv.FormatInt(int64(arg.SeatNumber[i]), 10)
		err := q.Client.HDel(ctx, SetKey, field).Err()
		if err != nil {
			return errors.New("seat key delete fail")
		}
	}

	return nil
}

func (q *RedisRepository) SetSeatValueAndSectionAndIncrSeatNumber(ctx context.Context, sections []string, prices []float64, seat []PostgresDB.Seat) error {
	SetKey := strconv.FormatInt(int64(seat[0].EventID), 10) + ":"
	for ind := range sections {
		SetSetcionKey := SetKey + "sections"
		err := q.Client.HSet(ctx, SetSetcionKey, sections[ind], prices[ind]).Err()
		if err != nil {
			return err
		}
	}
	for i := range seat {
		SetKey := strconv.FormatInt(int64(seat[i].EventID), 10) + ":" + seat[i].Section + ":" + "number"
		err := q.Client.SAdd(ctx, SetKey, seat[i].SeatNumber).Err()
		if err != nil {
			return err
		}
		SectionTotalKey := "{" + strconv.FormatInt(int64(seat[i].EventID), 10) + "}:" + seat[i].Section + ":total"
		err = q.Client.Incr(ctx, SectionTotalKey).Err()
		if err != nil {
			return err
		}
	}
	// IncrEventSeatNumberByEventID
	return nil
}
func (q *RedisRepository) SetSeatAndIncrSeatNumber(ctx context.Context, sections []string, seat []PostgresDB.Seat) error {

	for i := range seat {
		SetKey := strconv.FormatInt(int64(seat[i].EventID), 10) + ":" + seat[i].Section + ":" + "number"
		err := q.Client.SAdd(ctx, SetKey, seat[i].SeatNumber).Err()
		if err != nil {
			return err
		}
		SectionTotalKey := "{" + strconv.FormatInt(int64(seat[i].EventID), 10) + "}:" + seat[i].Section + ":total"
		err = q.Client.Incr(ctx, SectionTotalKey).Err()
		if err != nil {
			return err
		}
	}
	// IncrEventSeatNumberByEventID
	return nil
}

type GetSeatParams struct {
	EventId int
	Section []string
	Price   []float64
}

func (q *RedisRepository) DelSeatValue(ctx context.Context, arg []PostgresDB.Seat) error {
	for i := range arg {
		SetKey := strconv.FormatInt(int64(arg[i].EventID), 10) + ":" + arg[i].Section + ":" + "number"
		err := q.Client.SRem(ctx, SetKey, arg[i].SeatNumber).Err()
		if err != nil {
			return err
		}
		SectionTotalKey := "{" + strconv.FormatInt(int64(arg[i].EventID), 10) + "}:" + arg[i].Section + ":total"
		err = q.Client.Decr(ctx, SectionTotalKey).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func (q *RedisRepository) GetSeatValue(ctx context.Context, arg GetSeatParams) ([]PostgresDB.Seat, error) {
	Seatdata := []PostgresDB.Seat{}
	for i := range arg.Section {
		SetKey := strconv.FormatInt(int64(arg.EventId), 10) + ":" + arg.Section[i] + ":" + "number"

		item := q.Client.SMembers(ctx, SetKey)

		err := item.Err()
		if err != nil {
			return Seatdata, err
		}

		for _, seatnumber := range item.Val() {
			Number, _ := strconv.ParseInt(seatnumber, 10, 32)
			data := PostgresDB.Seat{
				Section:    arg.Section[i],
				SeatStatus: true,
				SeatNumber: int32(Number),
				Price:      arg.Price[i],
				EventID:    int32(arg.EventId),
			}
			Seatdata = append(Seatdata, data)
		}

	}

	return Seatdata, nil
}
func (q *RedisRepository) GetSectionsAndPricesByEventID(ctx context.Context, eventid int32) ([]PostgresDB.SelectSectionsAndPricesByEventIDRow, error) {
	SetKey := strconv.FormatInt(int64(eventid), 10) + ":" + "sections"
	item := q.Client.HGetAll(ctx, SetKey)
	data := []PostgresDB.SelectSectionsAndPricesByEventIDRow{}
	if item.Err() != nil {
		return data, item.Err()
	}

	for section, priceString := range item.Val() {
		pirces, _ := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
		data = append(data, PostgresDB.SelectSectionsAndPricesByEventIDRow{
			Section: section,
			Price:   pirces,
		})
	}
	return data, nil
}
func (q *RedisRepository) GetSectionsPricesBySectionsName(ctx context.Context, eventid int32, section string) (float64, error) {
	SetKey := strconv.FormatInt(int64(eventid), 10) + ":" + "sections"
	item := q.Client.HGet(ctx, SetKey, section)
	if item.Err() != nil {
		return 0, item.Err()
	}
	priceString := item.Val()
	pirces, _ := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
	return pirces, nil
}

type GetSectionSeatCount struct {
	EventId int
	Section []string
}

func (q *RedisRepository) GetSectionSeatCount(ctx context.Context, arg GetSectionSeatCount) (map[string]int32, error) {
	sectionAmount := map[string]int32{}
	var keys []string
	for i := range arg.Section {
		SectionTotalKey := "{" + strconv.FormatInt(int64(arg.EventId), 10) + "}:" + arg.Section[i] + ":total"
		// SectionTotalKey := strconv.FormatInt(int64(arg.EventId), 10) + ":" + arg.Section[i] + ":total"
		keys = append(keys, SectionTotalKey)
	}

	luaScript := `
    local result = {}
    for i, key in ipairs(KEYS) do
        local amount = redis.call("GET", key)
        if amount then
            table.insert(result, key)
            table.insert(result, tonumber(amount))
        else
            table.insert(result, key)
            table.insert(result, 0)
        end
    end
    return result
    `

	result, err := q.Client.Eval(ctx, luaScript, keys).Result()
	if err != nil {
		return nil, err
	}

	// 解析Lua腳本返回的列表
	if res, ok := result.([]interface{}); ok {
		for i := 0; i < len(res); i += 2 {
			key := res[i].(string)
			value := res[i+1].(int64)
			keyParts := strings.Split(key, ":")
			if len(keyParts) >= 3 {
				section := keyParts[1]
				sectionAmount[section] = int32(value)
			}
		}
	}

	return sectionAmount, nil
}
func (q *RedisRepository) IncrEventSeatNumberByEventID(ctx context.Context, eventid int32) {
	SetKey := strconv.FormatInt(int64(eventid), 10) + ":" + "total"
	item := q.Client.Incr(ctx, SetKey)
	// q.Client.Decr()
	fmt.Print(item.Val())
	// data := []PostgresDB.SelectSectionsAndPricesByEventIDRow{}
	// if item.Err() != nil {
	// 	return data, item.Err()
	// }

	// for section, priceString := range item.Val() {
	// 	pirces, _ := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
	// 	data = append(data, PostgresDB.SelectSectionsAndPricesByEventIDRow{
	// 		Section: section,
	// 		Price:   pirces,
	// 	})
	// }
	// return data, nil
}

type DelSeatParams struct {
	EventId int
	Section []string
}

func (q *RedisRepository) FlushAll(ctx context.Context, arg DelSeatParams) {
	for i := range arg.Section {
		SetKey := strconv.FormatInt(int64(arg.EventId), 10) + ":" + arg.Section[i] + ":" + "number"
		err := q.Client.Del(ctx, SetKey).Err()
		if err != nil {
			fmt.Println(err)
		}
		SetKeyNx := strconv.FormatInt(int64(arg.EventId), 10) + ":" + arg.Section[i] + ":Key"
		err = q.Client.Del(ctx, SetKeyNx).Err()
		if err != nil {
			fmt.Println(err)
		}
		SectionTotalKey := "{" + strconv.FormatInt(int64(arg.EventId), 10) + "}:" + arg.Section[i] + ":total"
		err = q.Client.Del(ctx, SectionTotalKey).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
	err := q.Client.Del(ctx, "OrderTradeNo").Err()
	if err != nil {
		fmt.Println(err)
	}
	return

}

func (q *RedisRepository) GetOrderTradeNoList(ctx context.Context) ([]string, error) {
	SetKey := "OrderTradeNo"
	item := q.Client.SMembers(ctx, SetKey)
	odertadeNo, err := item.Result()
	odertadeNo = item.Val()

	if err != nil {
		return []string{}, err
	}

	if len(odertadeNo) > 0 {
		if odertadeNo[0] == "" {
			return odertadeNo[1:], nil
		}
	}
	return odertadeNo, nil
}
func (q *RedisRepository) GetOrderTradeNoExist(ctx context.Context, orderTradeNo string) (bool, error) {
	SetKey := "OrderTradeNo"
	item := q.Client.SIsMember(ctx, SetKey, orderTradeNo)
	exist, err := item.Result()
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (q *RedisRepository) SetOrderTradeNo(ctx context.Context, orderTradeNo string) error {
	SetKey := "OrderTradeNo"
	err := q.Client.SAdd(ctx, SetKey, orderTradeNo).Err()
	if err != nil {
		return err
	}
	return nil
}
func (q *RedisRepository) DelOrderTradeNo(ctx context.Context, orderTradeNo []string) error {
	// SetKey := "OrderTradeNo"
	// for i := range orderTradeNo {
	// 	err := q.Client.SRem(ctx, SetKey, orderTradeNo[i]).Err()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// err := q.Client.SRem(ctx, SetKey, orderTradeNo).Err()
	// if err != nil {
	// 	return err
	// }
	return nil
}
