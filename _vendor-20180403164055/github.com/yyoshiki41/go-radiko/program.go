package radiko

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"time"

	"github.com/yyoshiki41/go-radiko/internal/util"
)

// Stations is a slice of Station.
type Stations []Station

// Station is a struct.
type Station struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Scd   Scd    `xml:"scd,omitempty"`
	Progs Progs  `xml:"progs,omitempty"`
}

// Scd is a struct.
type Scd struct {
	Progs Progs `xml:"progs"`
}

// Progs is a slice of Prog.
type Progs struct {
	Date  string `xml:"date"`
	Progs []Prog `xml:"prog"`
}

// Prog is a struct.
type Prog struct {
	Ft       string `xml:"ft,attr"`
	To       string `xml:"to,attr"`
	Ftl      string `xml:"ftl,attr"`
	Tol      string `xml:"tol,attr"`
	Dur      string `xml:"dur,attr"`
	Title    string `xml:"title"`
	SubTitle string `xml:"sub_title"`
	Desc     string `xml:"desc"`
	Pfm      string `xml:"pfm"`
	Info     string `xml:"info"`
	URL      string `xml:"url"`
}

// GetStations returns the program's meta-info.
func (c *Client) GetStations(ctx context.Context, date time.Time) (Stations, error) {
	apiEndpoint := path.Join(apiV3,
		"program/date", util.ProgramsDate(date),
		fmt.Sprintf("%s.xml", c.AreaID()))

	req, err := c.newRequest(ctx, "GET", apiEndpoint, &Params{})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d stationsData
	if err = decodeStationsData(resp.Body, &d); err != nil {
		return nil, err
	}
	return d.stations(), nil
}

// GetNowPrograms returns the program's meta-info which are currently on the air.
func (c *Client) GetNowPrograms(ctx context.Context) (Stations, error) {
	apiEndpoint := apiPath(apiV2, "program/now")

	req, err := c.newRequest(ctx, "GET", apiEndpoint, &Params{
		query: map[string]string{
			"area_id": c.AreaID(),
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d stationsData
	if err = decodeStationsData(resp.Body, &d); err != nil {
		return nil, err
	}
	return d.stations(), nil
}

// GetProgramByStartTime returns a specified program.
// This API wraps GetStations.
func (c *Client) GetProgramByStartTime(ctx context.Context, stationID string, start time.Time) (*Prog, error) {
	if stationID == "" {
		return nil, errors.New("StationID is empty")
	}

	stations, err := c.GetStations(ctx, start)
	if err != nil {
		return nil, err
	}

	ft := util.Datetime(start)
	var prog *Prog
	for _, s := range stations {
		if s.ID == stationID {
			for _, p := range s.Progs.Progs {
				if p.Ft == ft {
					prog = &p
					break
				}
			}
		}
	}
	if prog == nil {
		return nil, errors.New("program is not found")
	}
	return prog, nil
}

// GetWeeklyPrograms returns the weekly programs.
func (c *Client) GetWeeklyPrograms(ctx context.Context, stationID string) (Stations, error) {
	apiEndpoint := path.Join(apiV3,
		"program/station/weekly",
		fmt.Sprintf("%s.xml", stationID))

	req, err := c.newRequest(ctx, "GET", apiEndpoint, &Params{})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d stationsData
	if err = decodeStationsData(resp.Body, &d); err != nil {
		return nil, err
	}
	return d.stations(), nil
}

// stationsData includes a response struct for client's users.
type stationsData struct {
	XMLName     xml.Name `xml:"radiko"`
	XMLStations struct {
		XMLName  xml.Name `xml:"stations"`
		Stations Stations `xml:"station"`
	} `xml:"stations"`
}

// stations returns Stations which is a response struct for client's users.
func (d *stationsData) stations() Stations {
	return d.XMLStations.Stations
}

// decodeStationsData parses the XML-encoded data and stores the result.
func decodeStationsData(input io.Reader, stations *stationsData) error {
	b, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal(b, stations); err != nil {
		return err
	}
	return nil
}
