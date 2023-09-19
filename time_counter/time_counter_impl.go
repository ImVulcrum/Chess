package time_counter

import (
	"strconv"
	"time"
)

type t_counter struct {
	time_in_ms  int64
	marker_time int64
}

func New(starting_time int64) *t_counter {
	var c *t_counter = new(t_counter)
	c.time_in_ms = starting_time
	c.marker_time = 0

	return c
}

func (c *t_counter) Init_Counting() {
	c.marker_time = time.Now().UnixMilli()
}

func (c *t_counter) Return_Current_Counter() (string, bool) {
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

	if returntime <= 0 {
		return "00:00", true
	}
	return convert_time_in_ms_to_string(returntime), false
}

func (c *t_counter) Stop_Counting() {
	var current_time = time.Now().UnixMilli()
	if c.marker_time != 0 {
		c.time_in_ms = c.time_in_ms - (current_time - c.marker_time)
	}
	c.marker_time = 0
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
