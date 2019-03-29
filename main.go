package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

var funmap map[string]func(string,int)(string,error)

type Point struct {
	Name string
	Age int
}

var s = &Point{"test",12}

func GetSt() *Point  {
	return s
}

func main()  {

	s1 := GetSt()
	s.Age = 11
	fmt.Printf("s : %v address: %p \n", s, &s)
	fmt.Printf("s1 : %v address: %p \n", s1, &s1)
	fmt.Print(&*s1)

	os.Exit(0)
	funmap = make(map[string]func(string,int)(string,error))
	funmap["test"] = func(str string,i int)(stp string,err error) {
		fmt.Println("i'-a")
		if i == 1 {
			stp = "success"
			return stp,nil
		}else {
			return stp,errors.New("test")
		}
	}
	funmap["to"] = TestFun
	funmap["11"] = tlo

	if result,err := funmap["11"]("s",0);err != nil{
		log.Fatal(err)
	}else {
		fmt.Println(result)
	}

}

var TestFun = func(str string,i int)(stp string,err error) {
	fmt.Println("test",i)
	fmt.Println(str)
	if i == 0{
		stp = str
		return stp,nil
	}else {
		return stp,errors.New("gs")
	}
}

func tlo(str string,i int)(stp string,err error) {
	fmt.Println("test",i)
	fmt.Println(str)
	if i == 0{
		stp = str
		return stp,nil
	}else {
		return stp,errors.New("gs")
	}
}
