package dbhelpers

import (
	"math"

	"github.com/go-pg/pg/v10"
	"github.com/jivison/gowon-indexer/lib/customerrors"
	"github.com/jivison/gowon-indexer/lib/db"
)

func SelectArtistsWhereInMany(artists []string, itemsPerChunk float64) ([]db.Artist, error) {
	var chunks [][]interface{}
	var allArtists []db.Artist

	if len(artists) == 0 {
		return allArtists, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(artists))/itemsPerChunk))+1)

	for index, artist := range artists {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], artist)
	}

	for _, chunk := range chunks {
		var selectedArtists []db.Artist

		err := db.Db.Model((*db.Artist)(nil)).
			Where(
				"artist.name IN (?)",
				pg.In(
					chunk,
				),
			).Select(&selectedArtists)

		if err != nil {
			return allArtists, customerrors.DatabaseUnknownError()
		}

		allArtists = append(allArtists, selectedArtists...)
	}

	return allArtists, nil
}

func SelectAlbumsWhereInMany(albums []interface{}, itemsPerChunk float64) ([]db.Album, error) {
	var chunks [][]interface{}
	var allAlbums []db.Album

	if len(albums) == 0 {
		return allAlbums, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(albums))/itemsPerChunk))+1)

	for index, album := range albums {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], album)
	}

	for _, chunk := range chunks {
		var selectedAlbums []db.Album

		err := db.Db.Model((*db.Album)(nil)).Relation("Artist").
			Where(
				"(artist.name, album.name) IN (?)",
				pg.InMulti(
					chunk...,
				),
			).Select(&selectedAlbums)

		if err != nil {
			return allAlbums, customerrors.DatabaseUnknownError()
		}

		allAlbums = append(allAlbums, selectedAlbums...)
	}

	return allAlbums, nil
}

func SelectTracksWhereInMany(tracks []interface{}, itemsPerChunk float64) ([]db.Track, error) {
	var chunks [][]interface{}
	var allTracks []db.Track

	if len(tracks) == 0 {
		return allTracks, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(tracks))/itemsPerChunk))+1)

	for index, track := range tracks {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], track)
	}

	for _, chunk := range chunks {
		var selectedTracks []db.Track

		err := db.Db.Model((*db.Track)(nil)).
			Relation("Artist").
			Relation("Album").
			Where(
				"(artist.name, track.name, album.name) IN (?)",
				pg.InMulti(
					chunk...,
				),
			).Select(&selectedTracks)

		if err != nil {
			return allTracks, customerrors.DatabaseUnknownError()
		}

		allTracks = append(allTracks, selectedTracks...)
	}

	return allTracks, nil
}

func SelectArtistCountsWhereInMany(artistIDs []interface{}, userID int64, itemsPerChunk float64) ([]db.ArtistCount, error) {
	var chunks [][]interface{}
	var allArtistCounts []db.ArtistCount

	if len(artistIDs) == 0 {
		return allArtistCounts, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(artistIDs))/itemsPerChunk))+1)

	for index, artistID := range artistIDs {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], artistID)
	}

	for _, chunk := range chunks {
		var selectedArtistCounts []db.ArtistCount

		err := db.Db.Model((*db.ArtistCount)(nil)).
			Relation("Artist").
			Where(
				"artist_id IN (?)",
				pg.In(chunk),
			).
			Where("user_id = ?", userID).
			Select(&selectedArtistCounts)

		if err != nil {
			return allArtistCounts, customerrors.DatabaseUnknownError()
		}

		allArtistCounts = append(allArtistCounts, selectedArtistCounts...)
	}

	return allArtistCounts, nil
}

func SelectAlbumCountsWhereInMany(albumIDs []interface{}, userID int64, itemsPerChunk float64) ([]db.AlbumCount, error) {
	var chunks [][]interface{}
	var allAlbumCounts []db.AlbumCount

	if len(albumIDs) == 0 {
		return allAlbumCounts, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(albumIDs))/itemsPerChunk))+1)

	for index, albumID := range albumIDs {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], albumID)
	}

	for _, chunk := range chunks {
		var selectedAlbumCounts []db.AlbumCount

		err := db.Db.Model((*db.AlbumCount)(nil)).
			Relation("Album").
			Relation("Album.Artist").
			Where(
				"album_id IN (?)",
				pg.In(chunk),
			).
			Where("user_id = ?", userID).
			Select(&selectedAlbumCounts)

		if err != nil {
			return allAlbumCounts, customerrors.DatabaseUnknownError()
		}

		allAlbumCounts = append(allAlbumCounts, selectedAlbumCounts...)
	}

	return allAlbumCounts, nil
}

func SelectTrackCountsWhereInMany(trackIDs []interface{}, userID int64, itemsPerChunk float64) ([]db.TrackCount, error) {
	var chunks [][]interface{}
	var allTrackCounts []db.TrackCount

	if len(trackIDs) == 0 {
		return allTrackCounts, nil
	}

	chunks = make([][]interface{}, int(math.Floor(float64(len(trackIDs))/itemsPerChunk))+1)

	for index, trackID := range trackIDs {
		chunkIndex := int(math.Floor(float64(index+1) / itemsPerChunk))

		if chunks[chunkIndex] == nil {
			chunks[chunkIndex] = make([]interface{}, 0)
		}

		chunks[chunkIndex] = append(chunks[chunkIndex], trackID)
	}

	for _, chunk := range chunks {
		var selectedTrackCounts []db.TrackCount

		err := db.Db.Model((*db.TrackCount)(nil)).
			Relation("Track").
			Relation("Track.Artist").
			Relation("Track.Album").
			Where(
				"album_id IN (?)",
				pg.In(chunk),
			).
			Where("user_id = ?", userID).
			Select(&selectedTrackCounts)

		if err != nil {
			return allTrackCounts, customerrors.DatabaseUnknownError()
		}

		allTrackCounts = append(allTrackCounts, selectedTrackCounts...)
	}

	return allTrackCounts, nil
}
