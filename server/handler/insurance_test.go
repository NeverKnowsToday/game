package handler

import (
	"fmt"
	"testing"
)

func TestInsurance(t *testing.T) {

	v := []string{"1","1","1","1","2","2","2","3","3"}


	for i := 0; i < len(v); {
		fmt.Println( "i = ", i)
        if i == len(v)-1 {
        	break
		}
		for j := i+1; j < len(v) ; j++ {
			fmt.Println( "j = ", j)

			if v[i] == v[j] {
				v[j]= "0"
			}else {
				i = j
				break
			}

		}

	}


	for _, value:= range v {
		fmt.Println( value)
	}




	a:=20%10
	fmt.Println("a= ", a)



}


