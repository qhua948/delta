package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Mark struct {
	Reference
	Point
	Span

	GroundRelationship float64
}

type MarkList []Mark

func (m MarkList) Len() int           { return len(m) }
func (m MarkList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MarkList) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m MarkList) encode() [][]string {
	data := [][]string{{
		"Mark Code",
		"Network Code",
		"Mark Name",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Ground Relationship",
		"Start Time",
		"End Time",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Name),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
			strconv.FormatFloat(v.GroundRelationship, 'g', -1, 64),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *MarkList) decode(data [][]string) error {
	var marks []Mark
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 10 {
				return fmt.Errorf("incorrect number of installed mark fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[3], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[4], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}

			var ground float64
			if ground, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[8]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[9]); err != nil {
				return err
			}
			marks = append(marks, Mark{
				Reference: Reference{
					Code:    strings.TrimSpace(d[0]),
					Network: strings.TrimSpace(d[1]),
					Name:    strings.TrimSpace(d[2]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[6]),
				},
				GroundRelationship: ground,
			})
		}

		*m = MarkList(marks)
	}
	return nil
}

func LoadMarks(path string) ([]Mark, error) {
	var m []Mark

	if err := LoadList(path, (*MarkList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MarkList(m))

	return m, nil
}
