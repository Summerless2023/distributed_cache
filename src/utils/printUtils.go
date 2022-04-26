package utils

import (
	"container/list"
	"fmt"
)

func PrintList(myList *list.List) {
	for i := myList.Front(); i != nil; {
		fmt.Println(i.Value)
		i = i.Next()
	}
}
