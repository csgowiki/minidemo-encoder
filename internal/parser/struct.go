package parser

// Game is the overall struct that holds everything
type Game struct {
	MatchName     string      `json:"MatchId"`
	ClientName    string      `json:"ClientName"`
	Map           string      `json:"MapName"`
	TickRate      int32       `json:"TickRate"`
	PlaybackTicks int32       `json:"PlaybackTicks"`
	ParseRate     int         `json:"ParseRate"`
	Rounds        []GameRound `json:"GameRounds"`
}

// GameRound information and all of the associated events
type GameRound struct {
	RoundNum        int32       `json:"RoundNum"`
	StartTick       int32       `json:"StartTick"`
	FreezeTimeEnd   int32       `json:"FreezeTimeEnd"`
	EndTick         int32       `json:"EndTick"`
	EndOfficialTick int32       `json:"EndOfficialTick"`
	Reason          string      `json:"RoundEndReason"`
	TScore          int32       `json:"TScore"`
	CTScore         int32       `json:"CTScore"`
	EndTScore       int32       `json:"EndTScore"`
	EndCTScore      int32       `json:"EndCTScore"`
	CTTeam          *string     `json:"CTTeam"`
	TTeam           *string     `json:"TTeam"`
	Frames          []GameFrame `json:"Frames"`
}

// GameFrame (game state at time t)
type GameFrame struct {
	Tick int32         `json:"Tick"`
	T    TeamFrameInfo `json:"T"`
	CT   TeamFrameInfo `json:"CT"`
}

// TeamFrameInfo at time t
type TeamFrameInfo struct {
	Side    string       `json:"Side"`
	Team    string       `json:"TeamName"`
	Players []PlayerInfo `json:"Players"`
}

// PlayerInfo at time t
type PlayerInfo struct {
	PlayerName    string  `json:"Name"`
	X             float32 `json:"X"`
	Y             float32 `json:"Y"`
	Z             float32 `json:"Z"`
	ViewX         float32 `json:"ViewX"`
	ViewY         float32 `json:"ViewY"`
	VelocityX     float32 `json:"VelocityX"`
	VelocityY     float32 `json:"VelocityY"`
	VelocityZ     float32 `json:"VelocityZ"`
	PlayerButtons int32   `json:"Buttons"`
	ActiveWeapon  string  `json:"ActiveWeapon"`
	IsAlive       bool    `json:"IsAlive"`
}
