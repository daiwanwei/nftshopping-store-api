package utils

import (
	"encoding/json"
	"errors"
)

type Role string

const (
	TenantRole   Role = "Tenant"
	MerchantRole Role = "Merchant"
)

var AllOfRole = []Role{TenantRole, MerchantRole}

func (o *Role) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	role := Role(s)
	switch role {
	case TenantRole, MerchantRole:
		*o = role
		return nil
	default:
		return errors.New("Invalid role type")
	}
}

type Orientation string

const (
	South Orientation = "South"
	North Orientation = "North"
)

func (o *Orientation) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	orientation := Orientation(s)
	switch orientation {
	case South, North:
		*o = orientation
		return nil
	default:
		return errors.New("Invalid orientation type")
	}
}

type Furniture string

const (
	SingleBed Furniture = "SingleBed"
	DoubleBed Furniture = "DoubleBed"
	Closet    Furniture = "Closet"
	Sofa      Furniture = "Sofa"
	Chair     Furniture = "Chair"
	Shelf     Furniture = "Shelf"
	Desk      Furniture = "Desk"
)

func (f *Furniture) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	furniture := Furniture(s)
	switch furniture {
	case SingleBed, DoubleBed, Closet, Sofa, Desk, Chair, Shelf:
		*f = furniture
		return nil
	default:
		return errors.New("Invalid furniture type")
	}
}

type Group string

const (
	Couple       Group = "Couple"
	Family       Group = "Family"
	OfficeWorker Group = "OfficeWorker"
	Citizen      Group = "Citizen"
	Foreigner    Group = "Foreigner"
)

func (g *Group) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	group := Group(s)
	switch group {
	case Couple, Family, OfficeWorker, Citizen, Foreigner:
		*g = group
		return nil
	default:
		return errors.New("Invalid group type")
	}
}

type Amenity string

const (
	Washer          Amenity = "Washer"
	AirConditioning Amenity = "AirConditioning"
	Television      Amenity = "Television"
	Cable           Amenity = "Cable"
	Network         Amenity = "Network"
	Heater          Amenity = "Heater"
	Gas             Amenity = "Gas"
	Refrigerator    Amenity = "Refrigerator"
	Microwave       Amenity = "Microwave"
	GasCookTop      Amenity = "GasCookTop"
	WaterDispenser  Amenity = "WaterDispenser"
)

func (f *Amenity) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	amenity := Amenity(s)
	switch amenity {
	case Washer, AirConditioning, Television, Cable, Network, Heater,
		Gas, Refrigerator, Microwave, GasCookTop, WaterDispenser:
		*f = amenity
		return nil
	default:
		return errors.New("Invalid amenity type")
	}
}

type FireFighting string

const (
	FireExtinguisher  FireFighting = "FireExtinguisher"
	EscapeSling       FireFighting = "EscapeSling"
	EmergencyLighting FireFighting = "EmergencyLighting"
	Monitor           FireFighting = "Monitor"
)

func (f *FireFighting) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	fireFighting := FireFighting(s)
	switch fireFighting {
	case FireExtinguisher, EscapeSling, EmergencyLighting, Monitor:
		*f = fireFighting
		return nil
	default:
		return errors.New("Invalid fireFighting type")
	}
}

type People string

const (
	UnLimit People = "UnLimit"
	One     People = "One"
	Two     People = "Two"
)

func (f *People) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	people := People(s)
	switch people {
	case UnLimit, One, Two:
		*f = people
		return nil
	default:
		return errors.New("Invalid people type")
	}
}

type Restrict string

const (
	Pet       Restrict = "Pet"
	Shrine    Restrict = "Shrine"
	Household Restrict = "Household"
)

func (f *Restrict) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	restrict := Restrict(s)
	switch restrict {
	case Pet, Shrine, Household:
		*f = restrict
		return nil
	default:
		return errors.New("Invalid restrict type")
	}
}

type LifeFunction string

const (
	ConvenienceStore  LifeFunction = "ConvenienceStore"
	TraditionalMarket LifeFunction = "TraditionalMarket"
	Hypermarket       LifeFunction = "Hypermarket"
	DepartmentStore   LifeFunction = "DepartmentStore"
	Park              LifeFunction = "Park"
	Hospital          LifeFunction = "Hospital"
	NightMarket       LifeFunction = "NightMarket"
	School            LifeFunction = "School"
)

func (f *LifeFunction) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	lifeFunction := LifeFunction(s)
	switch lifeFunction {
	case ConvenienceStore, TraditionalMarket, Hypermarket, DepartmentStore,
		Park, Hospital, NightMarket, School:
		*f = lifeFunction
		return nil
	default:
		return errors.New("Invalid lifeFunction type")
	}
}
