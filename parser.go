package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"strconv"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

type UtilityRecord struct {
	player_name      string
	steamid          uint64
	utType           string
	throw_pitch      float32
	throw_yaw        float32
	throw_posX       float32
	throw_posY       float32
	throw_posZ       float32
  	velocity_x       float32
 	velocity_y       float32
  	velocity_z       float32
	end_posX         float32
	end_posY         float32
	end_posZ         float32
	round            int
	valid            bool
	start_time       time.Duration
	air_time         float32
	match_throw_time float32
	teamname         string
	is_walk          bool
	is_duck          bool
	is_jump          bool
	entity_posX		 float32
	entity_posY		 float32
	entity_posZ		 float32
}

var utrecord_collector map[int64]UtilityRecord
var type_map map[string]string
var filePath, mapName, resPath string
var tickRate float64

func ArgParser() {
	flag.StringVar(&filePath, "filepath", "Unknown", "Demo file path")
	flag.StringVar(&resPath, "topath", "Unknown", "Demo file path")
	flag.Parse()
}

func f2str(input_num float32) string {
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}

func JsonFomat(ut UtilityRecord, round int) string {
	var json_list = []string {
		f2str(ut.throw_pitch),
		f2str(ut.throw_yaw),
		f2str(ut.air_time),
		f2str(ut.end_posX),
		f2str(ut.end_posY),
		f2str(ut.end_posZ),
		strconv.FormatBool(ut.is_duck),
		strconv.FormatBool(ut.is_jump),
		strconv.FormatBool(ut.is_walk),
		strconv.Itoa(round),
		f2str(ut.match_throw_time),
		ut.player_name,
		strconv.FormatUint(ut.steamid, 10),
		ut.teamname,
		f2str(ut.throw_posX),
		f2str(ut.throw_posY),
		f2str(ut.throw_posZ),
		type_map[ut.utType],
		f2str(ut.velocity_x),
		f2str(ut.velocity_y),
		f2str(ut.velocity_z),
		f2str(ut.entity_posX),
		f2str(ut.entity_posY),
		f2str(ut.entity_posZ),
	}
	str, err := json.Marshal(json_list)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	const he_flash_time float32 = 1.63

	// arg info
	ArgParser()
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	round := 0
	var round_start_time time.Duration

	header, err := p.ParseHeader()
	if err != nil {
		panic(err)
	}
	mapName = header.MapName
	tickRate = p.TickRate()
	count := 0

	var infoPath = resPath
	var infoFile *os.File
	var infoError error

	if !checkFileIsExist(infoPath) {
		infoFile, infoError = os.Create(infoPath)
	}
	infoFile, infoError = os.OpenFile(infoPath, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if infoError != nil {
		panic(infoError)
	}
	defer infoFile.Close()

	// init
	infoFile.WriteString("[[]")

	type_map = make(map[string]string)
	type_map["Smoke Grenade"] = "smokegrenade"
	type_map["HE Grenade"] = "hegrenade"
	type_map["Flashbang"] = "flashbang"
	type_map["Incendiary Grenade"] = "incgrenade"
	type_map["Molotov"] = "molotov"

	utrecord_collector = make(map[int64]UtilityRecord)

	p.RegisterEventHandler(func(e events.MatchStartedChanged) {
		round = 1
		round_start_time = p.CurrentTime()
	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		round++
		round_start_time = p.CurrentTime()
	})

	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		round_start_time = p.CurrentTime()
	})

	p.RegisterEventHandler(func(e events.GrenadeProjectileThrow) {
		uId := int64(e.Projectile.WeaponInstance.UniqueID())
		_, ok := utrecord_collector[uId]
		if !ok {
			utrecord_collector[int64(e.Projectile.WeaponInstance.UniqueID())] = UtilityRecord{
				player_name:      string(e.Projectile.Thrower.Name),
				steamid:          uint64(e.Projectile.Thrower.SteamID64),
				throw_yaw:        float32(e.Projectile.Thrower.ViewDirectionX()),
				throw_pitch:      float32(e.Projectile.Thrower.ViewDirectionY()),
				throw_posX:       float32(e.Projectile.Thrower.LastAlivePosition.X),
				throw_posY:       float32(e.Projectile.Thrower.LastAlivePosition.Y),
				throw_posZ:       float32(e.Projectile.Thrower.LastAlivePosition.Z),
       			velocity_x:       float32(e.Projectile.Velocity().X),
        		velocity_y:       float32(e.Projectile.Velocity().Y),
        		velocity_z:       float32(e.Projectile.Velocity().Z),
				utType:           string(e.Projectile.WeaponInstance.String()),
				round:            int(round),
				valid:            false,
				start_time:       p.CurrentTime(),
				match_throw_time: float32((p.CurrentTime() - round_start_time).Seconds()),
				teamname:         string(e.Projectile.Thrower.TeamState.ClanName()),
				is_walk:          e.Projectile.Thrower.IsWalking(),
				is_jump:          e.Projectile.Thrower.IsAirborne(),
				is_duck:          e.Projectile.Thrower.Flags().DuckingKeyPressed(),
				entity_posX:	  float32(e.Projectile.Entity.Position().X),
				entity_posY:	  float32(e.Projectile.Entity.Position().Y),
				entity_posZ:	  float32(e.Projectile.Entity.Position().Z),
			}
		}
	})

	// SMOKE DETONATE
	p.RegisterEventHandler(func(e events.SmokeStart) {
		uId := int64(e.Grenade.UniqueID())
		utrecord, ok := utrecord_collector[uId]
		if ok && !utrecord.valid {
			utrecord.valid = true
			utrecord.end_posX = float32(e.Position.X)
			utrecord.end_posY = float32(e.Position.Y)
			utrecord.end_posZ = float32(e.Position.Z)
			end_time := p.CurrentTime()
			utrecord.air_time = float32((end_time - utrecord.start_time).Seconds())
			count++

			json_str := JsonFomat(utrecord, round)
			io.WriteString(infoFile, ","+json_str)

			utrecord_collector[uId] = utrecord
		}
	})

	// MOLOTOV & INC GRENADE DETONATE
	p.RegisterEventHandler(func(e events.GrenadeProjectileDestroy) {
		if e.Projectile.WeaponInstance.Type.String() != string("Incendiary Grenade") && e.Projectile.WeaponInstance.Type.String() != string("Molotov") {
			return
		}
		uId := int64(e.Projectile.WeaponInstance.UniqueID())

		utrecord, ok := utrecord_collector[uId]
		if ok && !utrecord.valid {
			utrecord.valid = true
			utrecord.end_posX = float32(e.Projectile.Position().X)
			utrecord.end_posY = float32(e.Projectile.Position().Y)
			utrecord.end_posZ = float32(e.Projectile.Position().Z)

			end_time := p.CurrentTime()
			utrecord.air_time = float32((end_time - utrecord.start_time).Seconds())
			count++

			json_str := JsonFomat(utrecord, round)
			io.WriteString(infoFile, ","+json_str)
			utrecord_collector[uId] = utrecord
		}
	})

	// FLASH DETONATE
	p.RegisterEventHandler(func(e events.FlashExplode) {
		uId := int64(e.Grenade.UniqueID())
		utrecord, ok := utrecord_collector[uId]
		if ok && !utrecord.valid {
			utrecord.valid = true
			utrecord.end_posX = float32(e.Position.X)
			utrecord.end_posY = float32(e.Position.Y)
			utrecord.end_posZ = float32(e.Position.Z)
			utrecord.air_time = he_flash_time
			count++

			json_str := JsonFomat(utrecord, round)
			io.WriteString(infoFile, ","+json_str)
			utrecord_collector[uId] = utrecord
		}
	})

	// HE GRENADE DETONATE
	p.RegisterEventHandler(func(e events.HeExplode) {
		uId := int64(e.Grenade.UniqueID())
		utrecord, ok := utrecord_collector[uId]
		if ok && !utrecord.valid {
			utrecord.valid = true
			utrecord.end_posX = float32(e.Position.X)
			utrecord.end_posY = float32(e.Position.Y)
			utrecord.end_posZ = float32(e.Position.Z)
			utrecord.air_time = he_flash_time
			count++
			json_str := JsonFomat(utrecord, round)
			io.WriteString(infoFile, ","+json_str)

			utrecord_collector[uId] = utrecord
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	// if err != nil {
	// 	panic(err)
	// }
	n, infoError := io.WriteString(infoFile, "]")
	if infoError != nil {
		panic(infoError)
	}
	fmt.Printf("writen %d bytes | %d\n", n, count)
}
