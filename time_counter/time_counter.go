package time_counter

type Counter interface {
	Init_Counting()
	Stop_Counting()
	Return_Current_Counter() (string, bool)
}
