# PolytopiaMapEditor

This project is an open source map editor for The Battle of Polytopia. The game generates random maps for each game, but there is no way to edit the maps. This editor solves this problem by giving users the ability to modify the save state containing the map and customize it to their liking.

<div style="display:inline-block;">
<img src="https://raw.githubusercontent.com/samuelyuan/PolytopiaMapEditor/master/screenshots/mapeditor.png" alt="earth" width="600" height="500" />
</div>

### Features

* Modify any tile, unit, improvement in the editor
* Enable cheats to reveal all tiles and unlock all tech
* Save all changes to the save state

### Usage

To use this editor, you must have provide a save state (.state) from the game. This editor should only be used for singleplayer games.

You should backup the save state before modifying it. Once you have modified the save state, make sure you quit your current Polytopia game and go to the main menu before overwriting the existing save file. If you overwrite the file while the game is still in progress, the game will overwrite the file when you leave and none of your new changes will apply.
