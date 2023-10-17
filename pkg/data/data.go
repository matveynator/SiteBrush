package Data

type RawData  struct {
	TagId                       string `xml:"TagID"`
	DiscoveryUnixTime           int64 
	ReaderIP                    string
	Antenna                     uint
	ProxyIP											string
}

type RawXMLData  struct {
	TagId                       string `xml:"TagID"`
	DiscoveryTime           		string
	Antenna                     uint
	ReaderIP										string
}


type Lap struct {
	Id                          int64
	SportsmanId                 string
	TagId                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryMaximalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	AverageResultsCount					uint
	RaceId                      uint
	RacePosition                uint
	RaceTotalTime               int64
	RaceFinished                bool

	LapNumber                   uint
	LapTime                     int64
	LapPosition                 uint
	TimeBehindTheLeader         int64
	LapIsLatest                 bool

	BestLapTime                 int64
	BestLapNumber               uint
	FastestLapInThisRace				bool
	FasterOrSlowerThanPreviousLapTime int64
	LapIsStrange                bool
}


