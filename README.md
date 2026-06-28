<h1 align="center">
  PKGFetch
</h1>

# Setup
## Setting Up Configuration
PKGFetch uses a toml file for it's config file and you can get the default version here [Default Config From Github Repo](https://github.com/PirateIsTaken/PKGFetch/blob/main/DefaultFiles/pkgfetch/pkgf.toml)
\
1. Download the config file into this folder: `~/.config/pkgfetch/pkgf.toml` (create the folders if you don't have them).

## Installing
Because there is no install script, you will have to do this manually too... sorry.
\
1. Download pkgf binary from github releases [here](https://github.com/PirateIsTaken/PKGFetch/releases)
2. Now, move it to `~/.local/bin` which is where pkgf will put all the symlnks of installed packages too (Ones that feature is added)
3. Now, make sure you have executable perms for `pkgf`. You can do that by running this command `chmod +x pkgf` replace *pkgf* with where your *pkgf* binary is located
4. Make sure the dir you moved it to is in your `$PATH` so you can access it from anywhere. If you don't know how to add a dir to `$PATH`, just google on how to add a directory to `$PATH` for your shell (To know which shell you're using, type `echo $SHELL` into your terminal).

# Usage
## Searching Packages
You can search packages from github repos.\
NOTE: Searching a package doesn't tell you if it has a release or not. (This is due to Github API Rate Limit)
\
Command: `pkgf search <package_name>`

## Installing Packages
Just like searching, while installing, you will see all the available packages from github.
\
\
1. you have to run install command: `pkgf install <package_name>`\
2. If there are releases available, It will show you all the available releases from that repo. You have to choose by typing a number between 1 to 10.
3. If the release has a suported format, you will get to choose what file you want to download and install \
3a. Note that the suported formats are only (rpm, appimage). Trying to install others will just return and do nothin or print out a message saying it's not supported.
4. After selecting the file, You will be asked for confirmation with information about the file.
5. Then the package will be installed accordingly.
6. At the end it will ask if you want to delete the cache file. And unless, you are testing pkgf, I suggest deleting it.
\
\
**That's It**
\
\
NOTE: If you don't see the app after the installation finished, Run this command which will make DEs like KDE and Gnome reload their Desktop Entries \
`update-desktop-database ~/.local/share/applications/` \
OR, you could just logout and log back in.
