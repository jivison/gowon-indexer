package controllers

import (
	"github.com/jivison/gowon-indexer/lib/graph/model"
	"github.com/jivison/gowon-indexer/lib/presenters"
	"github.com/jivison/gowon-indexer/lib/services/analysis"
	"github.com/jivison/gowon-indexer/lib/services/indexing"
)

// WhoKnowsArtist returns a list of who has plays of an artist
func WhoKnowsArtist(artistInput model.ArtistInput, settings *model.WhoKnowsSettings) (*model.WhoKnowsArtistResponse, error) {
	indexingService := indexing.CreateService()
	analysisService := analysis.CreateService()

	artist, err := indexingService.GetArtist(artistInput, false)

	if err != nil {
		return nil, err
	}

	whoKnows, err := analysisService.WhoKnowsArtist(artist, settings)

	return presenters.PresentWhoKnowsArtistResponse(artist, whoKnows), nil
}

// WhoKnowsAlbum returns a list of who has plays of an album
func WhoKnowsAlbum(albumInput model.AlbumInput, settings *model.WhoKnowsSettings) (*model.WhoKnowsAlbumResponse, error) {
	indexingService := indexing.CreateService()
	analysisService := analysis.CreateService()

	album, err := indexingService.GetAlbum(albumInput, false)

	if err != nil {
		return nil, err
	}

	whoKnows, err := analysisService.WhoKnowsAlbum(album, settings)

	return presenters.PresentWhoKnowsAlbumResponse(album, whoKnows), nil
}
