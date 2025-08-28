package entity

const (
	GenderInvalid Gender = "invalid"
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
)

type Gender string

func (g Gender) String() string {
	return string(g)
}
