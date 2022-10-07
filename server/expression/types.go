package expression

type Type uint

const (
	TYPE_NIL       Type = 0
	TYPE_STRING    Type = 1
	TYPE_NUMBER    Type = 2
	TYPE_ATTRIBUTE Type = 3
	TYPE_DURATION  Type = 4
)
