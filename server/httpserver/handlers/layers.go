package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"generator/collections/layers"
	"generator/internal/logger"
)

// ErrLayers is an internal error type for layers controller.
var ErrLayers = errs.Class("layers controller error")

// Layers is a mvc controller that handles all layers related views.
type Layers struct {
	log logger.Logger

	layers *layers.Service
}

// NewLayers is a constructor for layers controller.
func NewLayers(log logger.Logger, layers *layers.Service) *Layers {
	layersController := &Layers{
		log:    log,
		layers: layers,
	}

	return layersController
}

// Create is an endpoint that creates layer.
func (controller *Layers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	type request struct {
		Name     string `json:"name"`
		Priority int    `json:"priority"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrLayers.Wrap(err))
		return
	}

	err = controller.layers.Create(ctx, projectID, req.Name, req.Priority)
	if err != nil {
		controller.log.Error("could not create layer", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// List is an endpoint that returns all project layers.
func (controller *Layers) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	projectLayers, err := controller.layers.List(ctx, projectID)
	if err != nil {
		if layers.ErrNoLayerFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not list project layers", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(projectLayers); err != nil {
		controller.log.Error("failed to write json response", ErrLayers.Wrap(err))
		return
	}
}

// ListLayerImageNames is an endpoint that returns all names of layer images.
func (controller *Layers) ListLayerImageNames(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	layerImagesName, err := controller.layers.ListLayerImageNames(ctx, projectID, layerID)
	if err != nil {
		if layers.ErrNoLayerFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not list names of layer images", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(layerImagesName); err != nil {
		controller.log.Error("failed to write json response", ErrLayers.Wrap(err))
		return
	}
}

// UpdatePriority is an endpoint that updates layer priority.
func (controller *Layers) UpdatePriority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	type request struct {
		Priority int `json:"priority"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrLayers.Wrap(err))
		return
	}

	err = controller.layers.UpdatePriority(ctx, layerID, req.Priority)
	if err != nil {
		if layers.ErrNoLayerFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		if layers.ErrNoLayer.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not update layer priority", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// UpdateWeights is an endpoint that sets the weights of the layer images and of the layer being absent.
func (controller *Layers) UpdateWeights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	var req struct {
		NoneWeight int            `json:"noneWeight"`
		Weights    map[string]int `json:"weights"`
	}
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	err = controller.layers.UpdateWeights(ctx, layerID, req.NoneWeight, req.Weights)
	if err != nil {
		switch {
		case layers.ErrInvalidWeight.Has(err):
			controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		case layers.ErrNoLayer.Has(err), layers.ErrNoLayerFile.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
		default:
			controller.log.Error("could not update layer weights", ErrLayers.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		}
		return
	}
}

// Delete is an endpoint that deletes layer.
func (controller *Layers) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	err = controller.layers.Delete(ctx, layerID)
	if err != nil {
		if layers.ErrNoLayerFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		if layers.ErrNoLayer.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not delete layer", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// CreateLayerImage is an endpoint that creates layer image.
func (controller *Layers) CreateLayerImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
		Name  string `json:"name"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrLayers.Wrap(err))
		return
	}

	layerImg, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.New("invalid name"))
		return
	}

	err = controller.layers.CreateLayerImage(ctx, layerID, req.Name, bytes.NewBuffer(layerImg))
	if err != nil {
		controller.log.Error("could not update name of layer", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// UpdateImageName is an endpoint that updates name of layer image.
func (controller *Layers) UpdateImageName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	oldName, ok := vars["name"]
	if !ok {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.New("invalid name"))
		return
	}

	type request struct {
		NewName string `json:"newName"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrLayers.Wrap(err))
		return
	}

	err = controller.layers.UpdateNameOfLayerImage(ctx, layerID, req.NewName, oldName)
	if err != nil {
		if layers.ErrNoLayer.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not update name of layer", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// DeleteLayerImage is an endpoint that deletes layer image.
func (controller *Layers) DeleteLayerImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	layerID, err := uuid.Parse(vars["layerId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.Wrap(err))
		return
	}

	name, ok := vars["name"]
	if !ok {
		controller.serveError(w, http.StatusBadRequest, ErrLayers.New("invalid name"))
		return
	}

	err = controller.layers.DeleteLayerImage(ctx, layerID, name)
	if err != nil {
		if layers.ErrNoLayer.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrLayers.Wrap(err))
			return
		}
		controller.log.Error("could not update name of layer", ErrLayers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrLayers.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Layers) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrLayers.Wrap(err))
	}
}
