package time_counter

import (
	"strconv"
	"time"
)

type Time_Counter struct {
	time_in_ms  int64
	marker_time int64
}

func New(starting_time int64) *Time_Counter {
	var c *Time_Counter = new(Time_Counter)
	(*c).time_in_ms = starting_time
	(*c).marker_time = 0

	return c
}

func (c *Time_Counter) Init_Counting() {
	(*c).marker_time = time.Now().UnixMilli()
}

func (c *Time_Counter) Return_Current_Counter() (string, int64) {
	var current_time = time.Now().UnixMilli()

	var returntime int64

	if (*c).marker_time != 0 {
		// fmt.Println("retruntime:", returntime)
		// fmt.Println("current_time:", current_time, "marker:", (*c).marker_time)
		returntime = (*c).time_in_ms - (current_time - (*c).marker_time)
		// fmt.Println("retruntime:", returntime)
	} else {
		returntime = (*c).time_in_ms
	}
	return convert_time_in_ms_to_string(returntime), returntime
}

func (c *Time_Counter) Stop_Counting() {
	var current_time = time.Now().UnixMilli()
	if (*c).marker_time != 0 {
		(*c).time_in_ms = (*c).time_in_ms - (current_time - (*c).marker_time)
	}
	(*c).marker_time = 0
}

func convert_time_in_ms_to_string(time_remaining int64) string {
	var minutes int64 = time_remaining / 60000
	var seconds int64 = time_remaining % 60000
	var milliseconds int64 = seconds % 1000
	seconds = seconds / 1000

	//fmt.Println(minutes, ":", seconds, ":", milliseconds)

	var str_minutes string = "0" + strconv.Itoa(int(minutes))
	var str_seconds string = "0" + strconv.Itoa(int(seconds))
	var str_milliseconds string = "00" + strconv.Itoa(int(milliseconds))

	if str_minutes != "00" {
		return str_minutes[len(str_minutes)-2:] + ":" + str_seconds[len(str_seconds)-2:]
	} else {
		return str_seconds[len(str_seconds)-2:] + ":" + str_milliseconds[len(str_milliseconds)-3:len(str_milliseconds)-1]
	}
}
