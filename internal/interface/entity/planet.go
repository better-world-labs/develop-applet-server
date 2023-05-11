package entity

import "time"

type (
	PlanetRole         uint
	PlanetMemberStatus uint
)

const (
	PlanetRoleMember PlanetRole = 0
	PlanetRoleAdmin  PlanetRole = 1
	PlanetRoleRoot   PlanetRole = 2
)

const (
	PlanetMemberStatusOK    PlanetMemberStatus = 0
	PlanetMemberStatusBlack PlanetMemberStatus = 1
)

const MoyuPlanetId = 1

type (
	Planet struct {
		Id         int64     `json:"id"`
		Name       string    `json:"name"`
		Icon       string    `json:"icon"`
		FrontCover string    `json:"frontCover"`
		CreatedAt  time.Time `json:"createdAt"`
	}

	PlanetMember struct {
		Id        int64              `json:"id"`
		PlanetId  int64              `json:"planetId"`
		UserId    int64              `json:"userId"`
		Role      PlanetRole         `json:"role"`
		Status    PlanetMemberStatus `json:"status"`
		CreatedAt time.Time          `json:"createdAt"`
	}
)
