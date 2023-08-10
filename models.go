package coprime

type ApiType int

const (
    Prime  ApiType = iota
    Pro
    Advanced
    Sandbox
)

func (w ApiType) String() string {
	return [...]string{"prime", "pro", "advanced", "sandbox"}[w-1]
}

func (w ApiType) EnumIndex() int {
	return int(w)
}


type ServerTime struct {
	ISO   string  `json:"iso"`
	Epoch float64 `json:"epoch,number"`
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
