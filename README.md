lets make a game

2D, top-down, pixelized magic action RPG

contributions very welcome (see roadmap)

# [AERPG](https://github.com/aerth/rpg)

**demo**


[![Build Status](https://travis-ci.org/aerth/rpg.svg?branch=master)](https://travis-ci.org/aerth/rpg)

![screenshot](https://raw.githubusercontent.com/aerth/rpg/master/doc/screenshot.png)

## INCLUDED TOOLS

  * aerpg - play the demo rpg: explore map, kill skeletons, pick up loot, gain xp, collect magic items

  * mapmaker - create/edit a map:  read source code for keymap  

  * mapgen - generate a map: run `mapgen ______` to use specific seed such as `mapgen mycoolseed`

## FETCHING DEPENDENCIES

  * install Go: [go](https://golang.org)

  * install C dependencies: `apt-get install xorg-dev libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libopenal-dev libasound2-dev`

  * fetch this source code and dependencies: `go get -v -d -u github.com/aerth/rpg/cmd/...`

## INSTALLING AERPG GAME

### Requirements

If you're using Windows and having trouble building Pixel, please check [this
guide](https://github.com/faiface/pixel/wiki/Building-Pixel-on-Windows) on the
[wiki](https://github.com/faiface/pixel/wiki).

[PixelGL](https://godoc.org/github.com/faiface/pixel/pixelgl) backend uses OpenGL to render
graphics. Because of that, OpenGL development libraries are needed for compilation. The dependencies
are same as for [GLFW](https://github.com/go-gl/glfw).

The OpenGL version used is **OpenGL 3.3**.

- On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required
  headers and libraries.
- On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
- On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel
  libXinerama-devel mesa-libGL-devel libXi-devel` packages.
- See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.

**The combination of Go 1.8, macOS and latest XCode seems to be problematic** as mentioned in issue
[#7](https://github.com/faiface/pixel/issues/7). This issue is probably not related to Pixel.
**Upgrading to Go 1.8.1 fixes the issue.**

### Compiling

```go get -v -d github.com/aerth/rpg```
```GOBIN=$PWD go install github.com/aerth/rpg/cmd/...```

## keymap:

  * Pause, Inventory, Char Stats: `i`
  
  * Movement: `arrows`, `asdw`, `hjkl`, `hold right mouse`

  * Zoom: `mouse wheel`

  * Identify tile: `left click`

  * Pick up loot: `left click`

  * Attack (manastorm): `space` `middle click`

  * Attack (magic bullet): `B` `left click (point/shoot)`

  * Quit: ctrl+Q

## CHEATS

  * Toggle show enemy paths: `=`

  * Toggle fly mode: `caps lock`

  * Mana potion: `1`

  * Health potion: `2`

  * XP potion: `3`

  * Speed up time: `LSHIFT`

  * Slow motion: `TAB`

  * Random Loot: `ctrl+L` (random location), `ctrl+K` (under mouse)

  * Spawn fresh mob: `M` (watch FPS go down)

## ROADMAP

  * [ ] Regions (separated by doors)

  * [ ] Doors/Portals connect regions

  * [ ] Spawn tiles (instead of random tile) (should be Rectangle)

  * [ ] Map editor improvements (toolbar, pallet, fix offsets)
 
  * [ ] Text boxes (space to speed past conversations, pgup pgdown scroll)
  
  * [ ] Text Input (cheat codes, debug, chat, user input)

  * [x] Pick up loot

  * [ ] Drop item

  * [ ] proper Inventory and Wearing

  * [ ] Optimization

  * [ ] Replace spritesheets, allow texturepacks, skins

  * [ ] "Stage 1" map and missions, villages with markets, npcs, enemies, and a generated dungeon with bad guys and a boss

  * [ ] D2 style multiplayer co-operative and chat (no p2p)


### questions / support / donations

donations support the author and will make more frequent updates

BTC: ![https://blockchain.info/address/1ANjiTNvdEM6Me3yc4EBFSkDb4db4XW6pr](https://blockchain.info/qr?data=1ANjiTNvdEM6Me3yc4EBFSkDb4db4XW6pr&size=200)

PayPal Me: https://www.paypal.me/aerth

questions and comments can be directed to the email address published at https://github.com/aerth


### credits

  * font 'Admtas' by [adem taş](http://www.dafont.com/profile.php?user=980017)
  * font '[Terminus TTF Font](http://files.ax86.net/terminus-ttf/)' (SIL Open Font License, version 1.1)
  * [sprite generator](http://gaurav.munjal.us/Universal-LPC-Spritesheet-Character-Generator/)
  * [main character sprite](http://mmorpgmakerxb.com/p/characters-sprites-generator)
  * big thanks to the [pixel](https://github.com/faiface/pixel) library
  * additional credits in `assets/sprites/credits.txt` file

