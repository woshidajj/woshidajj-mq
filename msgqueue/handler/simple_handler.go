package handler

import "fmt"

type SimpleHandler struct {
}

func (h *SimpleHandler) HandleMsg(m interface{}) error {

	v, ok := m.(string)

	if ok {
		fmt.Println("----------start-----------")
		fmt.Printf("MSG - %s \n", v)
		fmt.Println("----------end-----------")
	} else {
		fmt.Println("----------start-----------")
		fmt.Println("INVALID M")
		fmt.Println("----------end-----------")
	}

	return nil

}
