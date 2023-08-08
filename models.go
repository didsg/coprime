package coprime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)



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
