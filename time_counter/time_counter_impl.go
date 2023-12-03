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

func (c *t_counter) Init_Counting() { //the way the counter works is by calculating the difference between now and the marker timer which is initalized with this method
	c.marker_time = time.Now().UnixMilli()
}

func (c *t_counter) Return_Current_Counter() (string, bool) {
	var current_time = time.Now().UnixMilli()

	var returntime int64

	if (*c).marker_time != 0 { //calculate difference between the current time and the marker time
		returntime = (*c).time_in_ms - (current_time - (*c).marker_time)
	} else { //if the marker time is 0 return the time left --> this means that the counter wasn't initialized
		returntime = (*c).time_in_ms
	}

	if returntime <= 0 { //negative returntime can happen due to minimal latency in the program --> in that case just return 0
		return "00:00", true
	}
	return convert_time_in_ms_to_string(returntime), false
}

func (c *t_counter) Stop_Counting() { //only if this is executed the time_in_ms var is actually updated
	var current_time = time.Now().UnixMilli()
	if c.marker_time != 0 {
		c.time_in_ms = c.time_in_ms - (current_time - c.marker_time)
	}
	c.marker_time = 0 //set the marker time to 0 to indicate that the timer is not counting
}

func convert_time_in_ms_to_string(time_remaining int64) string { //format the time portable string that is exactly 5 chars long
	var minutes int64 = time_remaining / 60000
	var seconds int64 = time_remaining % 60000
	var milliseconds int64 = seconds % 1000
	seconds = seconds / 1000

	var str_minutes string = "0" + strconv.Itoa(int(minutes))
	var str_seconds string = "0" + strconv.Itoa(int(seconds))
	var str_milliseconds string = "00" + strconv.Itoa(int(milliseconds))

	if str_minutes != "00" {
		return str_minutes[len(str_minutes)-2:] + ":" + str_seconds[len(str_seconds)-2:]
	} else {
		return str_seconds[len(str_seconds)-2:] + ":" + str_milliseconds[len(str_milliseconds)-3:len(str_milliseconds)-1]
	}
}
