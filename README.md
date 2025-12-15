# Spotify playlist sorter

This script sorts your Spotify playlist.

Sorting of the playlist follows these rules:

1. Duplicate tracks will be removed.
2. Tracks are sorted into groups of albums, that are ordered in the manner of the album.
3. Albums are grouped by artist and release date. (Artist is determined from the album rather than the track)
4. Artist groups are sorted depending on the first tracks appearance date on the playlist.

# Getting started

Create a Spotify application to get your `CLIENT_ID` and `CLIENT_SECRET` from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications).

Next you need to choose you variant of the application.

1. Terminal App (lazysps)
2. CLI App (sps)

Install you prefered variation from [releases](https://github.com/j0h-dev/my-spotify-playlist-sorter-go/releases/latest).

Just run the application and follow the instructions.

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
