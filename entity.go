package mediawiki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"

	"gitlab.com/tozd/go/errors"
	"gitlab.com/tozd/go/x"
)

type EntityType int

const (
	Item EntityType = iota
	Property
)

var timeRegex = regexp.MustCompile(`^([+-]\d{4,})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})Z$`)

func (t EntityType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t {
	case Item:
		buffer.WriteString("item")
	case Property:
		buffer.WriteString("property")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *EntityType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "item":
		*t = Item
	case "property":
		*t = Property
	default:
		return errors.Errorf("unknown entity type: %s", s)
	}
	return nil
}

type WikiBaseEntityType int

const (
	ItemType WikiBaseEntityType = iota
	PropertyType
	LexemeType
	FormType
	SenseType
)

func (t WikiBaseEntityType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t {
	case ItemType:
		buffer.WriteString("item")
	case PropertyType:
		buffer.WriteString("property")
	case LexemeType:
		buffer.WriteString("lexeme")
	case FormType:
		buffer.WriteString("form")
	case SenseType:
		buffer.WriteString("sense")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *WikiBaseEntityType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "item":
		*t = ItemType
	case "property":
		*t = PropertyType
	case "lexeme":
		*t = LexemeType
	case "form":
		*t = FormType
	case "sense":
		*t = SenseType
	default:
		return errors.Errorf("unknown wikibase entity type: %s", s)
	}
	return nil
}

type StatementType int

const (
	StatementT StatementType = iota
)

func (t StatementType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t { //nolint:gocritic
	case StatementT:
		buffer.WriteString("statement")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *StatementType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "statement":
		*t = StatementT
	default:
		return errors.Errorf("unknown statement type: %s", s)
	}
	return nil
}

type StatementRank int

const (
	Preferred StatementRank = iota
	Normal
	Deprecated
)

func (r StatementRank) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch r {
	case Preferred:
		buffer.WriteString("preferred")
	case Normal:
		buffer.WriteString("normal")
	case Deprecated:
		buffer.WriteString("deprecated")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (r *StatementRank) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "preferred":
		*r = Preferred
	case "normal":
		*r = Normal
	case "deprecated":
		*r = Deprecated
	default:
		return errors.Errorf("unknown statement rank: %s", s)
	}
	return nil
}

type SnakType int

const (
	Value SnakType = iota
	SomeValue
	NoValue
)

func (t SnakType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t {
	case Value:
		buffer.WriteString("value")
	case SomeValue:
		buffer.WriteString("somevalue")
	case NoValue:
		buffer.WriteString("novalue")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *SnakType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "value":
		*t = Value
	case "somevalue":
		*t = SomeValue
	case "novalue":
		*t = NoValue
	default:
		return errors.Errorf("unknown snak type: %s", s)
	}
	return nil
}

type DataType int

const (
	WikiBaseItem DataType = iota
	ExternalID
	String
	Quantity
	Time
	GlobeCoordinate
	CommonsMedia
	MonolingualText
	URL
	GeoShape
	WikiBaseLexeme
	WikiBaseSense
	WikiBaseProperty
	Math
	MusicalNotation
	WikiBaseForm
	TabularData
)

func (t DataType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t {
	case WikiBaseItem:
		buffer.WriteString("wikibase-item")
	case ExternalID:
		buffer.WriteString("external-id")
	case String:
		buffer.WriteString("string")
	case Quantity:
		buffer.WriteString("quantity")
	case Time:
		buffer.WriteString("time")
	case GlobeCoordinate:
		buffer.WriteString("globe-coordinate")
	case CommonsMedia:
		buffer.WriteString("commonsMedia")
	case MonolingualText:
		buffer.WriteString("monolingualtext")
	case URL:
		buffer.WriteString("url")
	case GeoShape:
		buffer.WriteString("geo-shape")
	case WikiBaseLexeme:
		buffer.WriteString("wikibase-lexeme")
	case WikiBaseSense:
		buffer.WriteString("wikibase-sense")
	case WikiBaseProperty:
		buffer.WriteString("wikibase-property")
	case Math:
		buffer.WriteString("math")
	case MusicalNotation:
		buffer.WriteString("musical-notation")
	case WikiBaseForm:
		buffer.WriteString("wikibase-form")
	case TabularData:
		buffer.WriteString("tabular-data")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (t *DataType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "wikibase-item":
		*t = WikiBaseItem
	case "external-id":
		*t = ExternalID
	case "string":
		*t = String
	case "quantity":
		*t = Quantity
	case "time":
		*t = Time
	case "globe-coordinate":
		*t = GlobeCoordinate
	case "commonsMedia":
		*t = CommonsMedia
	case "monolingualtext":
		*t = MonolingualText
	case "url":
		*t = URL
	case "geo-shape":
		*t = GeoShape
	case "wikibase-lexeme":
		*t = WikiBaseLexeme
	case "wikibase-sense":
		*t = WikiBaseSense
	case "wikibase-property":
		*t = WikiBaseProperty
	case "math":
		*t = Math
	case "musical-notation":
		*t = MusicalNotation
	case "wikibase-form":
		*t = WikiBaseForm
	case "tabular-data":
		*t = TabularData
	default:
		return errors.Errorf("unknown data type: %s", s)
	}
	return nil
}

type TimePrecision int

const (
	BillionYears TimePrecision = iota
	HoundredMillionYears
	TenMillionYears
	MillionYears
	HoundredMillenniums
	TenMillenniums
	Millennium
	Century
	Decade
	Year
	Month
	Day
	Hour
	Minute
	Second
)

type CalendarModel int

const (
	Gregorian CalendarModel = iota
	Julian
)

// MarshalJSON implements json.Marshaler interface for CalendarModel.
//
// Go enumeration values are converted to corresponding calendar Wikidata URIs.
// Those might be different (but equivalent) than what it was in the source dump.
func (t CalendarModel) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch t {
	case Gregorian:
		buffer.WriteString("https://www.wikidata.org/wiki/Q1985727")
	case Julian:
		buffer.WriteString("https://www.wikidata.org/wiki/Q1985786")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface for CalendarModel.
//
// It normalizes calendar Wikidata URIs to Go enumeration values.
func (t *CalendarModel) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case "https://www.wikidata.org/wiki/Q1985727":
		*t = Gregorian
	case "http://www.wikidata.org/entity/Q1985727":
		*t = Gregorian
	case "https://www.wikidata.org/wiki/Q1985786":
		*t = Julian
	case "http://www.wikidata.org/entity/Q1985786":
		*t = Julian
	default:
		return errors.Errorf("unknown calendar model: %s", s)
	}
	return nil
}

// ErrorValue represents an error with the value.
//
// When JSON representation contains an error, only error is provided
// as a Go value because any other field might be fail to parse.
type ErrorValue string

type StringValue string

type WikiBaseEntityIDValue struct {
	Type WikiBaseEntityType `json:"entity-type"`
	ID   string             `json:"id"`
}

type GlobeCoordinateValue struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Precision float64 `json:"precision"`
	Globe     string  `json:"globe"`
}

type MonolingualTextValue struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

// Amount is an arbitrary precision number and extends big.Rat.
type Amount struct {
	big.Rat
}

// MarshalJSON implements json.Marshaler interface for Amount.
func (a Amount) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteString(`"`)
	if a.Sign() >= 0 {
		// Sign is required always.
		b.WriteString("+")
	}
	b.WriteString(a.String())
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface for Amount.
func (a *Amount) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	_, ok := a.SetString(s)
	if !ok {
		return errors.Errorf("cannot parse amount into number: %s", s)
	}
	return nil
}

func (a *Amount) String() string {
	l, q := x.RatPrecision(&a.Rat)
	return a.FloatString(l + q)
}

type QuantityValue struct {
	Amount     Amount  `json:"amount"`
	UpperBound *Amount `json:"upperBound,omitempty"`
	LowerBound *Amount `json:"lowerBound,omitempty"`
	Unit       string  `json:"unit"`
}

// TimeValue represents a time value.
//
// While Time is a regular time.Time struct with nanoseconds precision,
// its real precision is available by Precision.
//
// Note that Wikidata uses historical numbering, in which year 0 is undefined
// and 1 BCE is represented by -1, but time.Time uses astronomical numbering,
// in which 1 BCE is represented by 0.
type TimeValue struct {
	Time      time.Time     `json:"time"`
	Precision TimePrecision `json:"precision"`
	Calendar  CalendarModel `json:"calendar"`
}

// MarshalJSON implements json.Marshaler interface for TimeValue.
func (v TimeValue) MarshalJSON() ([]byte, error) {
	type t struct {
		Time      string        `json:"time"`
		Precision TimePrecision `json:"precision"`
		Calendar  CalendarModel `json:"calendarmodel"`
	}
	formatedTime := formatTime(v.Time, v.Precision)
	b, err := json.Marshal(t{
		formatedTime,
		v.Precision,
		v.Calendar,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, err
}

// UnmarshalJSON implements json.Unmarshaler interface for TimeValue.
func (v *TimeValue) UnmarshalJSON(b []byte) error {
	type t struct {
		Time      string        `json:"time"`
		Precision TimePrecision `json:"precision"`
		Calendar  CalendarModel `json:"calendarmodel"`
	}
	var d t
	err := json.Unmarshal(b, &d)
	if err != nil {
		return errors.WithStack(err)
	}
	v.Time, err = parseTime(d.Time)
	if err != nil {
		return errors.WithStack(err)
	}
	v.Precision = d.Precision
	v.Calendar = d.Calendar
	return nil
}

// DataValue provides parsed value as Go value in Value.
//
// Value can be one of ErrorValue, StringValue, WikiBaseEntityIDValue,
// GlobeCoordinateValue, MonolingualTextValue, QuantityValue, and TimeValue.
type DataValue struct {
	Value interface{} `json:"value"`
}

func formatTime(t time.Time, p TimePrecision) string {
	t = t.UTC()
	year := t.Year()
	if year < 1 {
		// Wikidata uses historical numbering, in which year 0 is undefined,
		// but Go uses astronomical numbering, so we subtract 1 here.
		year--
	}
	month := t.Month()
	if p < Month {
		// Wikidata uses 0 when month is unknown or insignificant.
		month = 0
	}
	day := t.Day()
	if p < Day {
		// Wikidata uses 0 when day is unknown or insignificant.
		day = 0
	}
	return fmt.Sprintf("%+05d-%02d-%02dT%02d:%02d:%02dZ", year, month, day, t.Hour(), t.Minute(), t.Second())
}

// MarshalJSON implements json.Marshaler interface for DataValue.
//
// JSON representation of Go values might be different (but equivalent)
// than what it was in the source dump.
func (v DataValue) MarshalJSON() ([]byte, error) {
	switch value := v.Value.(type) {
	case ErrorValue:
		b, err := json.Marshal(struct {
			Error ErrorValue `json:"error"`
		}{value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case StringValue:
		b, err := json.Marshal(struct {
			Type  string      `json:"type"`
			Value StringValue `json:"value"`
		}{"string", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case WikiBaseEntityIDValue:
		b, err := json.Marshal(struct {
			Type  string                `json:"type"`
			Value WikiBaseEntityIDValue `json:"value"`
		}{"wikibase-entityid", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case GlobeCoordinateValue:
		b, err := json.Marshal(struct {
			Type  string               `json:"type"`
			Value GlobeCoordinateValue `json:"value"`
		}{"globecoordinate", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case MonolingualTextValue:
		b, err := json.Marshal(struct {
			Type  string               `json:"type"`
			Value MonolingualTextValue `json:"value"`
		}{"monolingualtext", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case QuantityValue:
		b, err := json.Marshal(struct {
			Type  string        `json:"type"`
			Value QuantityValue `json:"value"`
		}{"quantity", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	case TimeValue:
		b, err := json.Marshal(struct {
			Type  string    `json:"type"`
			Value TimeValue `json:"value"`
		}{"time", value})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, err
	}
	return nil, errors.Errorf(`unknown data value type: %+v`, v.Value)
}

func parseTime(t string) (time.Time, errors.E) {
	match := timeRegex.FindStringSubmatch(t)
	if match == nil {
		return time.Time{}, errors.Errorf(`unable to parse time "%s"`, t)
	}
	year, err := strconv.ParseInt(match[1], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse year "%s"`, t)
	}
	if year < 0 {
		// Wikidata uses historical numbering, in which year 0 is undefined,
		// but Go uses astronomical numbering, so we add 1 here.
		year++
	} else if year == 0 {
		return time.Time{}, errors.New("year cannot be 0")
	}
	month, err := strconv.ParseInt(match[2], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse month "%s"`, t)
	}
	if month == 0 {
		// Wikidata uses 0 when month is unknown or insignificant.
		// Go does not support this, so we set it to 1 here.
		month = 1
	}
	day, err := strconv.ParseInt(match[3], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse day "%s"`, t)
	}
	if day == 0 {
		// Wikidata uses 0 when day is unknown or insignificant.
		// Go does not support this, so we set it to 1 here.
		day = 1
	}
	hour, err := strconv.ParseInt(match[4], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse hour "%s"`, t)
	}
	minute, err := strconv.ParseInt(match[5], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse minute "%s"`, t)
	}
	second, err := strconv.ParseInt(match[6], 10, 0) //nolint:gomnd
	if err != nil {
		return time.Time{}, errors.WithMessagef(err, `unable to parse second "%s"`, t)
	}
	return time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC), nil
}

// UnmarshalJSON implements json.Unmarshaler interface for DataValue.
//
// It normalizes JSON representation to Go values.
func (v *DataValue) UnmarshalJSON(b []byte) error {
	var t struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return errors.WithStack(err)
	}
	if t.Error != "" {
		v.Value = ErrorValue(t.Error)
		return nil
	}
	switch t.Type {
	case "string":
		var t struct {
			Type  string      `json:"type"`
			Value StringValue `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		v.Value = t.Value
	case "wikibase-entityid":
		var t struct {
			Type string `json:"type"`
			// We do not use WikiBaseEntityIDValue because of extra fields.
			Value struct {
				Type WikiBaseEntityType `json:"entity-type"`
				ID   string             `json:"id"`
				// Not available for all entity types. Not recommended to be used. We ignore it.
				NumericID int `json:"numeric-id"`
			} `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		v.Value = WikiBaseEntityIDValue{
			Type: t.Value.Type,
			ID:   t.Value.ID,
		}
	case "globecoordinate":
		var t struct {
			Type string `json:"type"`
			// We do not use GlobeCoordinateValue because of extra fields.
			Value struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				// Altitude is deprecated and no longer used. We ignore it.
				Altitude  float64 `json:"altitude"`
				Precision float64 `json:"precision"`
				Globe     string  `json:"globe"`
			} `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		v.Value = GlobeCoordinateValue{
			Latitude:  t.Value.Latitude,
			Longitude: t.Value.Longitude,
			Precision: t.Value.Precision,
			Globe:     t.Value.Globe,
		}
	case "monolingualtext":
		var t struct {
			Type  string               `json:"type"`
			Value MonolingualTextValue `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		v.Value = t.Value
	case "quantity":
		var t struct {
			Type  string        `json:"type"`
			Value QuantityValue `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		v.Value = t.Value
	case "time":
		var t struct {
			Type string `json:"type"`
			// We do not use TimeValue because of extra fields.
			Value struct {
				Time      string        `json:"time"`
				Precision TimePrecision `json:"precision"`
				Calendar  CalendarModel `json:"calendarmodel"`
				// Defined and declared not used, but sometimes still set. We ignore it.
				Timezone int64 `json:"timezone"`
				// Defined and declared not used, but sometimes still set. We ignore it.
				Before int64 `json:"before"`
				// Defined and declared not used, but sometimes still set. We ignore it.
				After int64 `json:"after"`
			} `json:"value"`
		}
		err := x.UnmarshalWithoutUnknownFields(b, &t)
		if err != nil {
			return err
		}
		parsedTime, err := parseTime(t.Value.Time)
		if err != nil {
			v.Value = ErrorValue(err.Error())
		} else {
			v.Value = TimeValue{
				Time:      parsedTime,
				Precision: t.Value.Precision,
				Calendar:  t.Value.Calendar,
			}
		}
	default:
		return errors.Errorf(`unknown data value type "%s": %s`, t.Type, string(b))
	}
	return nil
}

type LanguageValue struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type SiteLink struct {
	Site   string   `json:"site"`
	Title  string   `json:"title"`
	Badges []string `json:"badges,omitempty"`
	URL    string   `json:"url,omitempty"`
}

type Snak struct {
	Hash      string     `json:"hash,omitempty"`
	SnakType  SnakType   `json:"snaktype"`
	Property  string     `json:"property"`
	DataType  DataType   `json:"datatype"`
	DataValue *DataValue `json:"datavalue,omitempty"`
}

type Reference struct {
	Hash       string            `json:"hash,omitempty"`
	Snaks      map[string][]Snak `json:"snaks,omitempty"`
	SnaksOrder []string          `json:"snaks-order,omitempty"`
}

type Statement struct {
	ID              string            `json:"id"`
	Type            StatementType     `json:"type"`
	MainSnak        Snak              `json:"mainsnak"`
	Rank            StatementRank     `json:"rank"`
	Qualifiers      map[string][]Snak `json:"qualifiers,omitempty"`
	QualifiersOrder []string          `json:"qualifiers-order,omitempty"`
	References      []Reference       `json:"references,omitempty"`
}

// Entity is a Wikidata entities JSON dump entity.
type Entity struct {
	ID           string                     `json:"id"`
	Type         EntityType                 `json:"type"`
	DataType     *DataType                  `json:"datatype,omitempty"`
	Labels       map[string]LanguageValue   `json:"labels,omitempty"`
	Descriptions map[string]LanguageValue   `json:"descriptions,omitempty"`
	Aliases      map[string][]LanguageValue `json:"aliases,omitempty"`
	Claims       map[string][]Statement     `json:"claims,omitempty"`
	SiteLinks    map[string]SiteLink        `json:"sitelinks,omitempty"`
	LastRevID    int64                      `json:"lastrevid"`
}
