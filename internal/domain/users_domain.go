package domain

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID            string   `json:"id"`
	CreatedAt     string   `json:"created_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	Salutation    string   `json:"salutation,omitempty"`
	Title         string   `json:"title,omitempty"`
	BirthDate     string   `json:"birth_date"`
	BirthCity     string   `json:"birth_city"`
	BirthCountry  string   `json:"birth_country"`
	BirthName     string   `json:"birth_name,omitempty"`
	Nationalities []string `json:"nationalities"`
	PostalAddress *Address `json:"postal_address,omitempty"`
	Address       Address  `json:"address"`
	Status        string   `json:"status,omitempty"`
}

type Address struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	Postcode     string `json:"postcode"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country"`
}

var (
	// Allowed values for salutation and title.
	validSalutations = map[string]struct{}{
		"":                          {},
		"SALUTATION_MALE":           {},
		"SALUTATION_FEMALE":         {},
		"SALUTATION_FEMALE_MARRIED": {},
		"SALUTATION_DIVERSE":        {},
	}
	validTitles = map[string]struct{}{
		"":         {},
		"DR":       {},
		"PROF":     {},
		"PROF_DR":  {},
		"DIPL_ING": {},
		"MAGISTER": {},
	}
	// Regex for postcode validation.
	postcodeRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\s\-]{0,8}[a-zA-Z0-9]?$`)
	// ISO 3166-1 alpha-2 country codes.
	validCountries = map[string]struct{}{
		"AD": {}, "AE": {}, "AF": {}, "AG": {}, "AI": {}, "AL": {}, "AM": {}, "AO": {}, "AQ": {}, "AR": {},
		"AS": {}, "AT": {}, "AU": {}, "AW": {}, "AX": {}, "AZ": {}, "BA": {}, "BB": {}, "BD": {}, "BE": {},
		"BF": {}, "BG": {}, "BH": {}, "BI": {}, "BJ": {}, "BL": {}, "BM": {}, "BN": {}, "BO": {}, "BQ": {},
		"BR": {}, "BS": {}, "BT": {}, "BV": {}, "BW": {}, "BY": {}, "BZ": {}, "CA": {}, "CC": {}, "CD": {},
		"CF": {}, "CG": {}, "CH": {}, "CI": {}, "CK": {}, "CL": {}, "CM": {}, "CN": {}, "CO": {}, "CR": {},
		"CU": {}, "CV": {}, "CW": {}, "CX": {}, "CY": {}, "CZ": {}, "DE": {}, "DJ": {}, "DK": {}, "DM": {},
		"DO": {}, "DZ": {}, "EC": {}, "EE": {}, "EG": {}, "EH": {}, "ER": {}, "ES": {}, "ET": {}, "FI": {},
		"FJ": {}, "FK": {}, "FM": {}, "FO": {}, "FR": {}, "GA": {}, "GB": {}, "GD": {}, "GE": {}, "GF": {},
		"GG": {}, "GH": {}, "GI": {}, "GL": {}, "GM": {}, "GN": {}, "GP": {}, "GQ": {}, "GR": {}, "GS": {},
		"GT": {}, "GU": {}, "GW": {}, "GY": {}, "HK": {}, "HM": {}, "HN": {}, "HR": {}, "HT": {}, "HU": {},
		"ID": {}, "IE": {}, "IL": {}, "IM": {}, "IN": {}, "IO": {}, "IQ": {}, "IR": {}, "IS": {}, "IT": {},
		"JE": {}, "JM": {}, "JO": {}, "JP": {}, "KE": {}, "KG": {}, "KH": {}, "KI": {}, "KM": {}, "KN": {},
		"KP": {}, "KR": {}, "KW": {}, "KY": {}, "KZ": {}, "LA": {}, "LB": {}, "LC": {}, "LI": {}, "LK": {},
		"LR": {}, "LS": {}, "LT": {}, "LU": {}, "LV": {}, "LY": {}, "MA": {}, "MC": {}, "MD": {}, "ME": {},
		"MF": {}, "MG": {}, "MH": {}, "MK": {}, "ML": {}, "MM": {}, "MN": {}, "MO": {}, "MP": {}, "MQ": {},
		"MR": {}, "MS": {}, "MT": {}, "MU": {}, "MV": {}, "MW": {}, "MX": {}, "MY": {}, "MZ": {}, "NA": {},
		"NC": {}, "NE": {}, "NF": {}, "NG": {}, "NI": {}, "NL": {}, "NO": {}, "NP": {}, "NR": {}, "NU": {},
		"NZ": {}, "OM": {}, "PA": {}, "PE": {}, "PF": {}, "PG": {}, "PH": {}, "PK": {}, "PL": {}, "PM": {},
		"PN": {}, "PR": {}, "PS": {}, "PT": {}, "PW": {}, "PY": {}, "QA": {}, "RE": {}, "RO": {}, "RS": {},
		"RU": {}, "RW": {}, "SA": {}, "SB": {}, "SC": {}, "SD": {}, "SE": {}, "SG": {}, "SH": {}, "SI": {},
		"SJ": {}, "SK": {}, "SL": {}, "SM": {}, "SN": {}, "SO": {}, "SR": {}, "SS": {}, "ST": {}, "SV": {},
		"SX": {}, "SY": {}, "SZ": {}, "TC": {}, "TD": {}, "TF": {}, "TG": {}, "TH": {}, "TJ": {}, "TK": {},
		"TL": {}, "TM": {}, "TN": {}, "TO": {}, "TR": {}, "TT": {}, "TV": {}, "TW": {}, "TZ": {}, "UA": {},
		"UG": {}, "UM": {}, "US": {}, "UY": {}, "UZ": {}, "VA": {}, "VC": {}, "VE": {}, "VG": {}, "VI": {},
		"VN": {}, "VU": {}, "WF": {}, "WS": {}, "XK": {}, "YE": {}, "YT": {}, "ZA": {}, "ZM": {}, "ZW": {},
	}
)

// Validate checks if the user object adheres to the spec.
func (u *User) Validate() error {
	if len(u.FirstName) < 2 || len(u.FirstName) > 100 {
		return errors.New("first_name must be between 2 and 100 characters")
	}
	if len(u.LastName) < 2 || len(u.LastName) > 100 {
		return errors.New("last_name must be between 2 and 100 characters")
	}
	if _, valid := validSalutations[u.Salutation]; !valid {
		return errors.New("invalid salutation")
	}
	if _, valid := validTitles[u.Title]; !valid {
		return errors.New("invalid title")
	}
	if _, err := time.Parse("2006-01-02", u.BirthDate); err != nil {
		return errors.New("birth_date must be in YYYY-MM-DD format")
	}
	if len(u.BirthCity) < 1 || len(u.BirthCity) > 85 {
		return errors.New("birth_city must be between 1 and 85 characters")
	}
	if _, valid := validCountries[u.BirthCountry]; !valid {
		return errors.New("invalid birth_country code")
	}
	if len(u.BirthName) > 100 {
		return errors.New("birth_name must be at most 100 characters")
	}
	if len(u.Nationalities) == 0 {
		return errors.New("at least one nationality is required")
	}
	for _, nationality := range u.Nationalities {
		if _, valid := validCountries[nationality]; !valid {
			return errors.New("invalid nationality code")
		}
	}
	if err := u.Address.Validate(); err != nil {
		return err
	}
	if u.PostalAddress != nil {
		if err := u.PostalAddress.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks if the address object adheres to the spec.
func (a *Address) Validate() error {
	if len(a.AddressLine1) == 0 || len(a.AddressLine1) > 100 {
		return errors.New("address_line1 must be between 1 and 100 characters")
	}
	if len(a.AddressLine2) > 100 {
		return errors.New("address_line2 must be at most 100 characters")
	}
	if !postcodeRegex.MatchString(a.Postcode) {
		return errors.New("postcode must match the required pattern")
	}
	if len(a.City) < 1 || len(a.City) > 85 {
		return errors.New("city must be between 1 and 85 characters")
	}
	if len(a.State) > 50 {
		return errors.New("state must be at most 50 characters")
	}
	if _, valid := validCountries[a.Country]; !valid {
		return errors.New("invalid country code")
	}
	return nil
}
