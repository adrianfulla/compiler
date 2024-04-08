package utils

// import "fmt"

// Stack representa una pila donde los elementos se añaden o se quitan desde el mismo lado (FILO).
type Stack struct {
	Elements []interface{} `json:"elements"` // Utiliza interface{} para permitir almacenar cualquier tipo de dato
}

// NewStack crea y devuelve una nueva pila vacía.
func NewStack() *Stack {
	return &Stack{
		Elements: make([]interface{}, 0),
	}
}

// Push añade un elemento a la parte superior de la pila.
func (s *Stack) Push(element interface{}) {
	s.Elements = append(s.Elements, element)
}

// Pop elimina y devuelve el elemento superior de la pila.
// Retorna nil si la pila está vacía.
func (s *Stack) Pop() interface{} {
	if len(s.Elements) == 0 {
		return nil
	}

	// Obtener el último elemento
	top := s.Elements[len(s.Elements)-1]
	s.Elements = s.Elements[:len(s.Elements)-1]
	return top
}

// Peek devuelve el elemento superior de la pila sin eliminarlo.
// Retorna nil si la pila está vacía.
func (s *Stack) Peek() interface{} {
	if len(s.Elements) == 0 {
		return nil
	}

	return s.Elements[len(s.Elements)-1]
}

// IsEmpty devuelve true si la pila no tiene elementos.
func (s *Stack) IsEmpty() bool {
	return len(s.Elements) == 0
}

// Size devuelve el número de elementos en la pila.
func (s *Stack) Size() int {
	return len(s.Elements)
}

func (s *Stack) ElemInStack(element interface{}) bool {
	for i := 0; i < len(s.Elements); i++ {
		if s.Elements[i] == element {
			return true
		}
	}
	return false
}
