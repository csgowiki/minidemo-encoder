package parser

import (
	// ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

const (
	IN_ATTACK    int32 = (1 << 0)
	IN_JUMP            = (1 << 1)
	IN_DUCK            = (1 << 2)
	IN_FORWARD         = (1 << 3)
	IN_BACK            = (1 << 4)
	IN_USE             = (1 << 5)
	IN_CANCEL          = (1 << 6)
	IN_LEFT            = (1 << 7)
	IN_RIGHT           = (1 << 8)
	IN_MOVELEFT        = (1 << 9)
	IN_MOVERIGHT       = (1 << 10)
	IN_ATTACK2         = (1 << 11)
	IN_RUN             = (1 << 12)
	IN_RELOAD          = (1 << 13)
	IN_ALT1            = (1 << 14)
	IN_ALT2            = (1 << 15)
	IN_SCORE           = (1 << 16) /**< Used by client.dll for when scoreboard is held down */
	IN_SPEED           = (1 << 17) /**< Player is holding the speed key */
	IN_WALK            = (1 << 18) /**< Player holding walk key */
	IN_ZOOM            = (1 << 19) /**< Zoom key for HUD zoom */
	IN_WEAPON1         = (1 << 20) /**< weapon defines these bits */
	IN_WEAPON2         = (1 << 21) /**< weapon defines these bits */
	IN_BULLRUSH        = (1 << 22)
	IN_GRENADE1        = (1 << 23) /**< grenade 1 */
	IN_GRENADE2        = (1 << 24) /**< grenade 2 */
	IN_ATTACK3         = (1 << 25)
)

func ButtonConvert(player *common.Player, addonButton int32) int32 {
	var button int32 = addonButton
	if player.Flags().DuckingKeyPressed() {
		button |= IN_DUCK
	}
	if player.IsWalking() {
		button |= IN_SPEED
	}
	if player.IsReloading {
		button |= IN_RELOAD
	}
	return button
}
