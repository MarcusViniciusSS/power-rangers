package entities

import (
    "time"
)

type PedidoRetirada struct {
    Id int
    CreateAt time.Time 
    UpdatedAt time.Time
    DeletedAt *time.Time
    Product string
    Acquired uint64
    Available uint64
    Safety int64
    Price float64
    Overbooking bool
}

