package entity

type Programs struct {
	Stations Stations `xml:"stations"`
}

type Stations struct {
	Stations []Station `xml:"station"`
}

type Station struct {
	Name string `xml:"name"`
	ID   string `xml:"id,attr"`
}
