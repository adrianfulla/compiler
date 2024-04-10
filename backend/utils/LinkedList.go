package utils
import(
	"fmt"
)

type LinkedNode struct {
    Value interface{}
    Prev  *LinkedNode
    Next  *LinkedNode
}

type DoublyLinkedList struct {
    Head *LinkedNode
    Tail *LinkedNode
	Size int
}

func (l *DoublyLinkedList) Prepend(value interface{}) {
    newNode := &LinkedNode{Value: value}
    if l.Head == nil {
        l.Head = newNode
        l.Tail = newNode
    } else {
        l.Head.Prev = newNode
        newNode.Next = l.Head
        l.Head = newNode
    }
	l.Size++
}

func (l *DoublyLinkedList) Append(value interface{}) {
    newNode := &LinkedNode{Value: value}
    if l.Tail == nil {
        l.Head = newNode
        l.Tail = newNode
    } else {
        l.Tail.Next = newNode
        newNode.Prev = l.Tail
        l.Tail = newNode
    }
	l.Size++
}

func (l *DoublyLinkedList) DeleteWithValue(value interface{}) {
    current := l.Head
    for current != nil {
        if current.Value == value {
            if current.Prev != nil {
                current.Prev.Next = current.Next
            } else {
                l.Head = current.Next
            }
            if current.Next != nil {
                current.Next.Prev = current.Prev
            } else {
                l.Tail = current.Prev
            }
			l.Size++
            return
        }
        current = current.Next
    }
}

func (l *DoublyLinkedList) PrintForward() {
    current := l.Head
    for current != nil {
        fmt.Println(current.Value)
        current = current.Next
    }
}

func (l *DoublyLinkedList) PrintReverse() {
    current := l.Tail
    for current != nil {
        fmt.Println(current.Value)
        current = current.Prev
    }
}


