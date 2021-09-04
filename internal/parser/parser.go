package parser

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"

	"log"
	"os"

	// demo_struct "github.com/hx-w/minidemo-encoder/internal/parser"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const IN_ATTACK = (1 << 0)
const IN_ATTACK2 = (1 << 11)
const IN_JUMP int32 = (1 << 1)
const IN_DUCK int32 = (1 << 2)
const IN_RELOAD int32 = (1 << 13)
const IN_SPEED int32 = (1 << 17)
const IN_ZOOM int32 = (1 << 19)

func init() {
	file, err := os.OpenFile("demoparser.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

var filePath string

func ArgParser() {
	flag.StringVar(&filePath, "filepath", "Unknown", "Demo file path")
	flag.Parse()
}

func convertRoundEndReason(r events.RoundEndReason) string {
	switch reason := r; reason {
	case 1:
		return "TargetBombed"
	case 2:
		return "VIPEscaped"
	case 3:
		return "VIPKilled"
	case 4:
		return "TerroristsEscaped"
	case 5:
		return "CTStoppedEscape"
	case 6:
		return "TerroristsStopped"
	case 7:
		return "BombDefused"
	case 8:
		return "CTWin"
	case 9:
		return "TWin"
	case 10:
		return "Draw"
	case 11:
		return "HostagesRescued"
	case 12:
		return "TargetSaved"
	case 13:
		return "HostagesNotRescued"
	case 14:
		return "TerroristsNotEscaped"
	case 15:
		return "VIPNotEscaped"
	case 16:
		return "GameStart"
	case 17:
		return "TerroristsSurrender"
	case 18:
		return "CTSurrender"
	default:
		return "Unknown"
	}
}

func parsePlayer(player *common.Player) PlayerInfo {
	currentPlayer := PlayerInfo{}
	currentPlayer.PlayerName = player.Name
	currentPlayer.VelocityX = float32(player.Velocity().X)
	currentPlayer.VelocityY = float32(player.Velocity().Y)
	currentPlayer.VelocityZ = float32(player.Velocity().Z)
	currentPlayer.X = float32(player.LastAlivePosition.X)
	currentPlayer.Y = float32(player.LastAlivePosition.Y)
	currentPlayer.Z = float32(player.LastAlivePosition.Z)
	currentPlayer.ViewX = player.ViewDirectionX()
	currentPlayer.ViewY = player.ViewDirectionY()
	currentPlayer.IsAlive = player.IsAlive()
	currentPlayer.PlayerButtons = 0
	if player.IsReloading {
		currentPlayer.PlayerButtons |= IN_RELOAD
	}
	if player.IsDucking() || player.IsDuckingInProgress() || player.IsUnDuckingInProgress() {
		currentPlayer.PlayerButtons |= IN_DUCK
	}
	if player.IsAirborne() {
		currentPlayer.PlayerButtons |= IN_JUMP
	}
	if player.IsWalking() {
		currentPlayer.PlayerButtons |= IN_SPEED
	}
	if player.IsScoped() {
		currentPlayer.PlayerButtons |= IN_ZOOM
	}
	if currentPlayer.IsAlive {
		currentPlayer.ActiveWeapon = player.ActiveWeapon().String()
	}
	//
	// ...
	//
	return currentPlayer
}

// Define cleaning functions
func cleanMapName(mapName string) string {
	lastSlash := strings.LastIndex(mapName, "/")
	if lastSlash == -1 {
		return mapName
	}
	return mapName[lastSlash+1:]
}

func main() {
	// arg info
	ArgParser()
	f, err := os.Open(filePath)
	checkError(err)
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	// flags
	currentFrameIdx := 0
	roundStarted := 0
	roundInEndTime := 0
	roundInFreezetime := 0

	// header
	header, err := p.ParseHeader()
	checkError(err)

	// game object
	currentGame := Game{}
	if p.TickRate() == 0 {
		currentGame.TickRate = 128
	} else {
		currentGame.TickRate = int32(p.TickRate())
	}
	currentGame.PlaybackTicks = int32(header.PlaybackTicks)
	currentGame.ClientName = header.ClientName
	currentGame.ParseRate = 1
	currentGame.MatchName = "demo"
	currentGame.Map = cleanMapName(header.MapName)

	// game round
	currentRound := GameRound{}
	InfoLogger.Printf("Demo is of type %s with tickrate %d \n", currentGame.ClientName, currentGame.TickRate)

	p.RegisterEventHandler(func(e events.RoundStart) {
		gs := p.GameState()
		if roundStarted == 1 {
			currentRound.EndOfficialTick = int32(gs.IngameTick()) - (5 * currentGame.TickRate)
			currentGame.Rounds = append(currentGame.Rounds, currentRound)
		}
		roundStarted = 1
		roundInFreezetime = 1
		roundInEndTime = 0
		currentRound = GameRound{}
		currentRound.RoundNum = int32(len(currentGame.Rounds) + 1)
		currentRound.StartTick = int32(gs.IngameTick())
		currentRound.TScore = int32(gs.TeamTerrorists().Score())
		currentRound.CTScore = int32(gs.TeamCounterTerrorists().Score())
		tTeam := gs.TeamTerrorists().ClanName()
		ctTeam := gs.TeamCounterTerrorists().ClanName()
		currentRound.TTeam = &tTeam
		currentRound.CTTeam = &ctTeam
	})

	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		gs := p.GameState()
		roundInFreezetime = 0
		currentRound.FreezeTimeEnd = int32(gs.IngameTick())
	})

	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		gs := p.GameState()
		if roundInEndTime == 0 {
			currentRound.EndOfficialTick = int32(gs.IngameTick())
			tPlayers := gs.TeamTerrorists().Members()
			aliveT := 0
			ctPlayers := gs.TeamCounterTerrorists().Members()
			aliveCT := 0
			for _, p := range tPlayers {
				if p.IsAlive() && p != nil {
					aliveT = aliveT + 1
				}
			}
			for _, p := range ctPlayers {
				if p.IsAlive() && p != nil {
					aliveCT = aliveCT + 1
				}
			}
			// reasonable ?
			if aliveCT == 0 {
				currentRound.Reason = "TWin"
				currentRound.EndTScore = currentRound.TScore + 1
				currentRound.EndCTScore = currentRound.CTScore
			} else {
				currentRound.Reason = "CTWin"
				currentRound.EndCTScore = currentRound.CTScore + 1
				currentRound.EndTScore = currentRound.TScore
			}
		}
	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		gs := p.GameState()
		if roundStarted == 0 {
			roundStarted = 1
			currentRound.RoundNum = 0
			currentRound.StartTick = 0
			currentRound.TScore = 0
			currentRound.CTScore = 0
			tTeam := gs.TeamTerrorists().ClanName()
			ctTeam := gs.TeamCounterTerrorists().ClanName()
			currentRound.TTeam = &tTeam
			currentRound.CTTeam = &ctTeam
		}
		roundInEndTime = 1
		currentRound.EndTick = int32(gs.IngameTick())
		currentRound.EndOfficialTick = int32(gs.IngameTick())
		currentRound.Reason = convertRoundEndReason(e.Reason)
		switch e.Winner {
		case common.TeamTerrorists:
			currentRound.EndTScore = currentRound.TScore + 1
			currentRound.EndCTScore = currentRound.CTScore
		case common.TeamCounterTerrorists:
			currentRound.EndCTScore = currentRound.CTScore + 1
			currentRound.EndTScore = currentRound.TScore
		}
	})

	p.RegisterEventHandler(func(e events.FrameDone) {
		gs := p.GameState()

		if (roundInFreezetime == 0) && (currentFrameIdx == 0) {
			currentFrame := GameFrame{}
			currentFrame.Tick = int32(gs.IngameTick())
			// Parse T
			currentFrame.T = TeamFrameInfo{}
			currentFrame.T.Side = "T"
			currentFrame.T.Team = gs.TeamTerrorists().ClanName()
			tPlayers := gs.TeamTerrorists().Members()
			for _, player := range tPlayers {
				if player != nil {
					currentFrame.T.Players = append(currentFrame.T.Players, parsePlayer(player))
				}
			}
			// Parse CT
			currentFrame.CT = TeamFrameInfo{}
			currentFrame.CT.Side = "CT"
			currentFrame.CT.Team = gs.TeamCounterTerrorists().ClanName()
			ctPlayers := gs.TeamCounterTerrorists().Members()

			for _, player := range ctPlayers {
				if player != nil {
					currentFrame.CT.Players = append(currentFrame.CT.Players, parsePlayer(player))
				}
			}

			currentRound.Frames = append(currentRound.Frames, currentFrame)
			if currentFrameIdx == (currentGame.ParseRate - 1) {
				currentFrameIdx = 0
			} else {
				currentFrameIdx = currentFrameIdx + 1
			}
		} else {
			if currentFrameIdx == (currentGame.ParseRate - 1) {
				currentFrameIdx = 0
			} else {
				currentFrameIdx = currentFrameIdx + 1
			}
		}
	})
	// Parse to end
	err = p.ParseToEnd()
	checkError(err)

	currentGame.Rounds = append(currentGame.Rounds, currentRound)

	// clean rounds
	if len(currentGame.Rounds) > 0 {
		InfoLogger.Println("Cleaning data")

		// Remove rounds where win reason doesn't exist
		var tempRoundsReason []GameRound
		for i := range currentGame.Rounds {
			currRound := currentGame.Rounds[i]
			if currRound.Reason == "CTWin" || currRound.Reason == "BombDefused" || currRound.Reason == "TargetSaved" || currRound.Reason == "TWin" || currRound.Reason == "TargetBombed" {
				tempRoundsReason = append(tempRoundsReason, currRound)
			}
		}
		currentGame.Rounds = tempRoundsReason

		// Remove rounds with missing end or start tick
		var tempRoundsTicks []GameRound
		for i := range currentGame.Rounds {
			currRound := currentGame.Rounds[i]
			if currRound.StartTick > 0 && currRound.EndTick > 0 {
				tempRoundsTicks = append(tempRoundsTicks, currRound)
			} else {
				if currRound.EndTick > 0 {
					tempRoundsTicks = append(tempRoundsTicks, currRound)
				}
			}
		}
		currentGame.Rounds = tempRoundsTicks

		// Remove rounds that dip in score
		var tempRoundsDip []GameRound
		for i := range currentGame.Rounds {
			if i > 0 && i < len(currentGame.Rounds) {
				prevRound := currentGame.Rounds[i-1]
				currRound := currentGame.Rounds[i]
				if currRound.CTScore+currRound.TScore >= prevRound.CTScore+prevRound.TScore {
					tempRoundsDip = append(tempRoundsDip, currRound)
				}
			} else if i == 0 {
				currRound := currentGame.Rounds[i]
				tempRoundsDip = append(tempRoundsDip, currRound)
			}
		}
		currentGame.Rounds = tempRoundsDip

		// Set first round scores to 0-0
		currentGame.Rounds[0].TScore = 0
		currentGame.Rounds[0].CTScore = 0

		// Remove rounds where score doesn't change
		var tempRounds []GameRound
		for i := range currentGame.Rounds {
			if i < len(currentGame.Rounds)-1 {
				nextRound := currentGame.Rounds[i+1]
				currRound := currentGame.Rounds[i]
				if !(currRound.CTScore+currRound.TScore >= nextRound.CTScore+nextRound.TScore) {
					tempRounds = append(tempRounds, currRound)
				}
			} else {
				currRound := currentGame.Rounds[i]
				tempRounds = append(tempRounds, currRound)
			}

		}
		currentGame.Rounds = tempRounds

		// Find the starting round. Starting round is defined as the first 0-0 round which has following rounds.
		startIdx := 0
		for i, r := range currentGame.Rounds {
			if (i < len(currentGame.Rounds)-3) && (len(currentGame.Rounds) > 3) {
				if (r.TScore+r.CTScore == 0) && (currentGame.Rounds[i+1].TScore+currentGame.Rounds[i+1].CTScore > 0) && (currentGame.Rounds[i+2].TScore+currentGame.Rounds[i+2].CTScore > 0) && (currentGame.Rounds[i+3].TScore+currentGame.Rounds[i+4].CTScore > 0) {
					startIdx = i
				}
			}
		}
		currentGame.Rounds = currentGame.Rounds[startIdx:len(currentGame.Rounds)]

		// Remove rounds with 0-0 scorelines that arent first
		var tempRoundsScores []GameRound
		for i := range currentGame.Rounds {
			currRound := currentGame.Rounds[i]
			if i > 0 {
				if currRound.TScore+currRound.CTScore > 0 {
					tempRoundsScores = append(tempRoundsScores, currRound)
				}
			} else {
				tempRoundsScores = append(tempRoundsScores, currRound)
			}
		}
		currentGame.Rounds = tempRoundsScores

		// Determine scores
		for i := range currentGame.Rounds {
			if i == 15 {
				currentGame.Rounds[i].TScore = currentGame.Rounds[i-1].EndCTScore
				currentGame.Rounds[i].CTScore = currentGame.Rounds[i-1].EndTScore
				if currentGame.Rounds[i].Reason == "CTWin" || currentGame.Rounds[i].Reason == "BombDefused" || currentGame.Rounds[i].Reason == "TargetSaved" {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore + 1
				} else {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore + 1
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore
				}
			} else if i > 0 {
				currentGame.Rounds[i].TScore = currentGame.Rounds[i-1].EndTScore
				currentGame.Rounds[i].CTScore = currentGame.Rounds[i-1].EndCTScore
				if currentGame.Rounds[i].Reason == "CTWin" || currentGame.Rounds[i].Reason == "BombDefused" || currentGame.Rounds[i].Reason == "TargetSaved" {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore + 1
				} else {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore + 1
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore
				}
			} else if i == 0 {
				// Set first round to 0-0, switch other scores
				currentGame.Rounds[i].TScore = 0
				currentGame.Rounds[i].CTScore = 0
				if currentGame.Rounds[i].Reason == "CTWin" || currentGame.Rounds[i].Reason == "BombDefused" || currentGame.Rounds[i].Reason == "TargetSaved" {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore + 1
				} else {
					currentGame.Rounds[i].EndTScore = currentGame.Rounds[i].TScore + 1
					currentGame.Rounds[i].EndCTScore = currentGame.Rounds[i].CTScore
				}
			}
		}

		// Set correct round numbers
		for i := range currentGame.Rounds {
			currentGame.Rounds[i].RoundNum = int32(i + 1)
		}

		InfoLogger.Println("Cleaned data, writing to JSON file")

		// Write the JSON
		// file, _ := json.MarshalIndent(currentGame.Rounds[3], "", " ")
		file, _ := json.MarshalIndent(currentGame, "", " ")
		_ = ioutil.WriteFile("json"+"/"+currentGame.MatchName+".json", file, 0644)
		InfoLogger.Println("Wrote to JSON file to: " + "json" + "/" + currentGame.MatchName + ".json")
	}
}

// Function to handle errors
func checkError(err error) {
	if err != nil {
		ErrorLogger.Println("DEMO STREAM ERROR")
		WarningLogger.Println("Demo stream errors can still write output, check for JSON file")
		ErrorLogger.Println(err.Error())
	}
}
