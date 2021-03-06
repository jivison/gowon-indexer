package indexing

import (
	"github.com/jivison/gowon-indexer/lib/db"
	apihelpers "github.com/jivison/gowon-indexer/lib/helpers/api"
	helpers "github.com/jivison/gowon-indexer/lib/helpers/generic"
	"github.com/jivison/gowon-indexer/lib/services/lastfm"
)

func (i Indexing) GenerateCountsFromScrobbles(scrobbles []lastfm.RecentTrack, user db.User) ([]db.ArtistCount, []db.AlbumCount, []db.TrackCount, []db.Play, error) {
	artistsMap, albumsMap, tracksMap, err := i.ConvertScrobbles(scrobbles)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	var dbScrobbles []db.Play

	artistCounts := make(map[string]int)
	albumCounts := make(map[string]map[string]int)
	trackCounts := make(map[string]map[string]map[string]int)

	for _, scrobble := range scrobbles {
		artist := artistsMap[scrobble.Artist.Text]
		artistCounts[artist.Name] += 1

		if _, ok := trackCounts[artist.Name]; !ok {
			trackCounts[scrobble.Artist.Text] = make(map[string]map[string]int)
		}
		if _, ok := trackCounts[artist.Name][scrobble.Album.Text]; !ok {
			trackCounts[artist.Name][scrobble.Album.Text] = make(map[string]int)
		}

		track := tracksMap[scrobble.Artist.Text][scrobble.Album.Text][scrobble.Name]
		trackCounts[artist.Name][scrobble.Album.Text][scrobble.Name] += 1

		if scrobble.Album.Text != "" {
			if _, ok := albumCounts[artist.Name]; !ok {
				albumCounts[artist.Name] = make(map[string]int)
			}

			album := albumsMap[scrobble.Artist.Text][scrobble.Album.Text]
			albumCounts[artist.Name][album.Name] += 1
		}

		timestamp, _ := apihelpers.ParseUnix(scrobble.Timestamp.UTS)
		dbScrobbles = append(dbScrobbles, db.Play{
			UserID: user.ID,
			User:   &user,

			TrackID: track.ID,
			Track:   &track,

			ScrobbledAt: timestamp,
		})
	}

	var dbArtistCounts []db.ArtistCount
	var dbAlbumCounts []db.AlbumCount
	var dbTrackCounts []db.TrackCount

	for artist, count := range artistCounts {
		dbArtist := artistsMap[artist]
		dbArtistCounts = append(dbArtistCounts, db.ArtistCount{Artist: &dbArtist, ArtistID: dbArtist.ID, User: &user, UserID: user.ID, Playcount: int32(count)})
	}

	for artist, artistAlbums := range albumCounts {
		for album, count := range artistAlbums {
			dbAlbum := albumsMap[artist][album]
			dbAlbumCounts = append(dbAlbumCounts, db.AlbumCount{Album: &dbAlbum, AlbumID: dbAlbum.ID, User: &user, UserID: user.ID, Playcount: int32(count)})
		}
	}

	for artist, artistAlbums := range trackCounts {
		for album, albumTracks := range artistAlbums {
			for track, count := range albumTracks {
				dbTrack := tracksMap[artist][album][track]
				dbTrackCounts = append(dbTrackCounts, db.TrackCount{Track: &dbTrack, TrackID: dbTrack.ID, User: &user, UserID: user.ID, Playcount: int32(count)})
			}
		}
	}

	return dbArtistCounts, dbAlbumCounts, dbTrackCounts, dbScrobbles, nil
}

func (i Indexing) ConvertScrobbles(scrobbles []lastfm.RecentTrack) (ArtistsMap, AlbumsMap, TracksMap, error) {
	artistsList, albumsList, tracksList := i.GenerateUniqueLists(scrobbles)

	artistsMap, err := i.ConvertArtists(artistsList)

	if err != nil {
		return nil, nil, nil, err
	}

	albumsMap, err := i.ConvertAlbums(albumsList, &artistsMap)

	if err != nil {
		return nil, nil, nil, err
	}

	tracksMap, err := i.ConvertTracks(tracksList, &artistsMap, &albumsMap)

	if err != nil {
		return nil, nil, nil, err
	}

	return artistsMap, albumsMap, tracksMap, nil
}

func (I Indexing) GenerateUniqueLists(scrobbles []lastfm.RecentTrack) ([]ArtistToConvert, []AlbumToConvert, []TrackToConvert) {
	artists := make(map[string]int)
	albums := make(map[string]map[string]int)
	tracks := make(map[string]map[string]map[string]int)

	for _, scrobble := range scrobbles {
		if _, ok := albums[scrobble.Artist.Text]; !ok {
			albums[scrobble.Artist.Text] = make(map[string]int)
		}
		if _, ok := tracks[scrobble.Artist.Text]; !ok {
			tracks[scrobble.Artist.Text] = make(map[string]map[string]int)
		}
		if _, ok := tracks[scrobble.Artist.Text][scrobble.Album.Text]; !ok {
			tracks[scrobble.Artist.Text][scrobble.Album.Text] = make(map[string]int)
		}

		artists[scrobble.Artist.Text] = 1
		albums[scrobble.Artist.Text][scrobble.Album.Text] = 1
		tracks[scrobble.Artist.Text][scrobble.Album.Text][scrobble.Name] = 1
	}

	var artistsList []ArtistToConvert
	var albumsList []AlbumToConvert
	var tracksList []TrackToConvert

	for artist := range artists {
		artistsList = append(artistsList, artist)
	}

	for artist, artistAlbums := range albums {
		for album := range artistAlbums {
			albumsList = append(albumsList, AlbumToConvert{
				AlbumName:  album,
				ArtistName: artist,
			})
		}
	}

	for artist, artistAlbums := range tracks {
		for album, albumTracks := range artistAlbums {
			for track := range albumTracks {
				var albumName *string

				if album != "" {
					copiedAlbumName := helpers.DeepCopy(album)
					albumName = &copiedAlbumName
				}

				tracksList = append(tracksList, TrackToConvert{
					ArtistName: artist,
					TrackName:  track,
					AlbumName:  albumName,
				})
			}
		}
	}

	return artistsList, albumsList, tracksList
}
