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

1. Find the existing save state for your current game. For Windows, this is usually stored at ``%USERPROFILE%\AppData\LocalLow\Midjiwan\Polytopia\00000000-0000-0000-0000-000000000000\Data\Singleplayer``.
2. Before modifying the save state, you should backup the save state in case you made any changes that break the save state.
3. Open the editor. Open file and click on the save state that you want to modify.
4. Make changes to the save state in the editor.
5. Once you are done making changes, click "File", "Save Map State..." and click "Yes" to overwrite save state file.
6. Copy the new save state file to the original save state location in step 1 and overwrite the existing file.

Before overwriting the existing save file, you must quit your current Polytopia game and go to the main menu or exit the game. If you overwrite the file while the game is still in progress, the game will overwrite the file when you leave the game and none of your new changes will apply.
