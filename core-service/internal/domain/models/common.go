package models

import (
	"github.com/google/uuid"
)

type Photo struct {
	Coords      string
	Description string
	Id          uuid.UUID
	ImgUrl      string
	Place       string
	RegionId    uuid.UUID
	TripId      uuid.UUID
	UserId      uuid.UUID
}

type Region struct {
	Id        uuid.UUID
	Name      string
	Country   string
	ObjectKey string
}

type Trip struct {
	Description string
	Id          uuid.UUID
	Name        string
	UserId      uuid.UUID
	RegionId    uuid.UUID
	Place       string
}

type Tag struct {
	Id   uuid.UUID
	Name string
}

type User struct {
	AvatarUrl string
	FullName  string
	Id        uuid.UUID
	Trips     []uuid.UUID
	UserName  string
	BirthDate string
	Education string
	City      string
}

type UserShort struct {
	AvatarUrl  string
	TripsCount int
	UserName   string
}

type Like struct {
	Id      uuid.UUID
	PhotoId uuid.UUID
	UserId  uuid.UUID
}

type PhotoFiltersDTO struct {
	RegionId uuid.UUID
	TripId   uuid.UUID
	TagKey   string
}
