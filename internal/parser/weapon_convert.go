package parser

import ilog "github.com/hx-w/minidemo-encoder/internal/logger"

type CSWeaponID int32

const (
	CSWeapon_NONE CSWeaponID = iota
	CSWeapon_P228
	CSWeapon_GLOCK
	CSWeapon_SCOUT
	CSWeapon_HEGRENADE
	CSWeapon_XM1014
	CSWeapon_C4
	CSWeapon_MAC10
	CSWeapon_AUG
	CSWeapon_SMOKEGRENADE
	CSWeapon_ELITE
	CSWeapon_FIVESEVEN
	CSWeapon_UMP45
	CSWeapon_SG550
	CSWeapon_GALIL
	CSWeapon_FAMAS
	CSWeapon_USP
	CSWeapon_AWP
	CSWeapon_MP5NAVY
	CSWeapon_M249
	CSWeapon_M3
	CSWeapon_M4A1
	CSWeapon_TMP
	CSWeapon_G3SG1
	CSWeapon_FLASHBANG
	CSWeapon_DEAGLE
	CSWeapon_SG552
	CSWeapon_AK47
	CSWeapon_KNIFE
	CSWeapon_P90
	CSWeapon_SHIELD
	CSWeapon_KEVLAR
	CSWeapon_ASSAULTSUIT
	CSWeapon_NIGHTVISION //Anything below is CS:GO ONLY
	CSWeapon_GALILAR
	CSWeapon_BIZON
	CSWeapon_MAG7
	CSWeapon_NEGEV
	CSWeapon_SAWEDOFF
	CSWeapon_TEC9
	CSWeapon_TASER
	CSWeapon_HKP2000
	CSWeapon_MP7
	CSWeapon_MP9
	CSWeapon_NOVA
	CSWeapon_P250
	CSWeapon_SCAR17
	CSWeapon_SCAR20
	CSWeapon_SG556
	CSWeapon_SSG08
	CSWeapon_KNIFE_GG
	CSWeapon_MOLOTOV
	CSWeapon_DECOY
	CSWeapon_INCGRENADE
	CSWeapon_DEFUSER
	CSWeapon_HEAVYASSAULTSUIT
	//The rest are actual item definition indexes for CS:GO
	CSWeapon_CUTTERS               = 56
	CSWeapon_HEALTHSHOT            = 57
	CSWeapon_KNIFE_T               = 59
	CSWeapon_M4A1_SILENCER         = 60
	CSWeapon_USP_SILENCER          = 61
	CSWeapon_CZ75A                 = 63
	CSWeapon_REVOLVER              = 64
	CSWeapon_TAGGRENADE            = 68
	CSWeapon_FISTS                 = 69
	CSWeapon_BREACHCHARGE          = 70
	CSWeapon_TABLET                = 72
	CSWeapon_MELEE                 = 74
	CSWeapon_AXE                   = 75
	CSWeapon_HAMMER                = 76
	CSWeapon_SPANNER               = 78
	CSWeapon_KNIFE_GHOST           = 80
	CSWeapon_FIREBOMB              = 81
	CSWeapon_DIVERSION             = 82
	CSWeapon_FRAGGRENADE           = 83
	CSWeapon_SNOWBALL              = 84
	CSWeapon_BUMPMINE              = 85
	CSWeapon_MAX_WEAPONS_NO_KNIFES // Max without the knife item defs, useful when treating all knives as a regular knife.
	CSWeapon_BAYONET               = 500
	CSWeapon_KNIFE_CLASSIC         = 503
	CSWeapon_KNIFE_FLIP            = 505
	CSWeapon_KNIFE_GUT             = 506
	CSWeapon_KNIFE_KARAMBIT        = 507
	CSWeapon_KNIFE_M9_BAYONET      = 508
	CSWeapon_KNIFE_TATICAL         = 509
	CSWeapon_KNIFE_FALCHION        = 512
	CSWeapon_KNIFE_SURVIVAL_BOWIE  = 514
	CSWeapon_KNIFE_BUTTERFLY       = 515
	CSWeapon_KNIFE_PUSH            = 516
	CSWeapon_KNIFE_CORD            = 517
	CSWeapon_KNIFE_CANIS           = 518
	CSWeapon_KNIFE_URSUS           = 519
	CSWeapon_KNIFE_GYPSY_JACKKNIFE = 520
	CSWeapon_KNIFE_OUTDOOR         = 521
	CSWeapon_KNIFE_STILETTO        = 522
	CSWeapon_KNIFE_WIDOWMAKER      = 523
	CSWeapon_KNIFE_SKELETON        = 525
	CSWeapon_MAX_WEAPONS           //THIS MUST BE LAST, EASY WAY TO CREATE LOOPS. When looping, do CS_IsValidWeaponID(i), to check.
)

var WeaponMap map[string]CSWeaponID

func init() {
	WeaponMap = map[string]CSWeaponID{
		// Melee
		"Knife":           CSWeapon_KNIFE,
		"Butterfly Knife": CSWeapon_KNIFE_BUTTERFLY,
		"Skeleton Knife":  CSWeapon_KNIFE_SKELETON,
		"Karambit":        CSWeapon_KNIFE_KARAMBIT,
		"Paracord Knife":  CSWeapon_KNIFE_CORD,
		"Bayonet":         CSWeapon_BAYONET,
		"Classic Knife":   CSWeapon_KNIFE_CLASSIC,
		"M9 Baynonet":     CSWeapon_KNIFE_M9_BAYONET,
		"Bowie Knife":     CSWeapon_KNIFE_SURVIVAL_BOWIE,
		"Falchion Knife":  CSWeapon_KNIFE_FALCHION,
		"Flip Knife":      CSWeapon_KNIFE_FLIP,
		"Gut Knife":       CSWeapon_KNIFE_GUT,
		"Huntsman Knife":  CSWeapon_KNIFE_TATICAL,
		"Navaja Knife":    CSWeapon_KNIFE_GYPSY_JACKKNIFE,
		"Shadow Daggers":  CSWeapon_KNIFE_PUSH,
		"Stiletto Knife":  CSWeapon_KNIFE_STILETTO,
		"Survival Knife":  CSWeapon_KNIFE_CANIS,
		"Talon Knife":     CSWeapon_KNIFE_WIDOWMAKER,
		"Ursus Knife":     CSWeapon_KNIFE_URSUS,
		// Pistols
		"Desert Eagle":  CSWeapon_DEAGLE,
		"R8 Revolver":   CSWeapon_REVOLVER,
		"Glock-18":      CSWeapon_GLOCK,
		"USP-S":         CSWeapon_USP,
		"Five-SeveN":    CSWeapon_FIVESEVEN,
		"Tec-9":         CSWeapon_TEC9,
		"p250":          CSWeapon_P250,
		"CZ75-Auto":     CSWeapon_CZ75A,
		"Dual Berettas": CSWeapon_ELITE,
		"P2000":         CSWeapon_HKP2000,
		// Shotguns
		"XM1014":    CSWeapon_XM1014,
		"Nova":      CSWeapon_NOVA,
		"MAG-7":     CSWeapon_MAG7,
		"Sawed-Off": CSWeapon_SAWEDOFF,
		// Submachine guns
		"MAC-10":   CSWeapon_MAC10,
		"MP5-SD":   CSWeapon_MP5NAVY,
		"MP7":      CSWeapon_MP7,
		"MP9":      CSWeapon_MP9,
		"P90":      CSWeapon_P90,
		"PP-Bizon": CSWeapon_BIZON,
		"UMP-45":   CSWeapon_UMP45,
		// Rifles
		"AK-47":    CSWeapon_AK47,
		"AUG":      CSWeapon_AUG,
		"AWP":      CSWeapon_AWP,
		"FAMAS":    CSWeapon_FAMAS,
		"G3SG1":    CSWeapon_G3SG1,
		"Galil AR": CSWeapon_GALILAR,
		"M4A1":     CSWeapon_M4A1_SILENCER,
		"M4A4":     CSWeapon_M4A1,
		"SCAR-20":  CSWeapon_SCAR20,
		"SG 553":   CSWeapon_SG556,
		"SSG 08":   CSWeapon_SSG08,
		// Machine guns
		"M249":  CSWeapon_M249,
		"Negev": CSWeapon_NEGEV,
		// Grenades
		"Decoy Grenade":      CSWeapon_DECOY,
		"Flashbang":          CSWeapon_FLASHBANG,
		"HE Grenade":         CSWeapon_HEGRENADE,
		"Incendiary Grenade": CSWeapon_INCGRENADE,
		"Molotov":            CSWeapon_MOLOTOV,
		"Smoke Grenade":      CSWeapon_SMOKEGRENADE,
		"Snowball":           CSWeapon_SNOWBALL,
		"Zeus x27":           CSWeapon_TASER,
		"C4":                 CSWeapon_C4,
	}
}

func WeaponStr2ID(weaponName string) CSWeaponID {
	if WeaponID, ok := WeaponMap[weaponName]; ok {
		return WeaponID
	} else {
		ilog.WarningLogger.Printf("[WeaponConvert] <%s> missing", weaponName)
		return CSWeapon_NONE
	}
}
