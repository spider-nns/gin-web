package custValidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Booking struct {
	CheckIn  time.Time `form:"checkIn" binding:"required,BookableDate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"checkOut" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func BookableDate(fl validator.FieldLevel) bool {
	if  date,ok:=fl.Field().Interface().(time.Time);ok{
		today:=time.Now()
		fmt.Println("date:",date)
		if date.Unix()>today.Unix(){
			fmt.Println("date unix ：",date.Unix())
			return true
		}
	}
	return false
}
