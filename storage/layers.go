package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections/layers"
)

// ensures that layersDB implements layers.DB.
var _ layers.DB = (*layersDB)(nil)

// ErrLayers indicates that there was an error in the database.
var ErrLayers = errs.Class("layers repository error")

// layersDB provide access to layers DB.
//
// architecture: Database
type layersDB struct {
	conn *sql.DB
}

// Create creates layer in the database.
func (layersDB *layersDB) Create(ctx context.Context, layer layers.Layer) error {
	query := `INSERT INTO layers(id, project_id, name, priority)
 	          VALUES ($1, $2, $3, $4)`

	_, err := layersDB.conn.ExecContext(ctx, query, layer.ID, layer.ProjectID, layer.Name, layer.Priority)

	return ErrLayers.Wrap(err)
}

// List returns all project layers ordered by priority from the database.
func (layersDB *layersDB) List(ctx context.Context, projectID uuid.UUID) (_ []layers.Layer, err error) {
	var allLayers []layers.Layer
	query := `SELECT id, project_id, name, priority, none_weight
 	          FROM layers
 	          WHERE project_id=$1
 	          ORDER BY priority ASC`

	rows, err := layersDB.conn.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var layer layers.Layer
		err = rows.Scan(&layer.ID, &layer.ProjectID, &layer.Name, &layer.Priority, &layer.NoneWeight)
		if err != nil {
			return nil, ErrLayers.Wrap(err)
		}

		allLayers = append(allLayers, layer)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrLayers.Wrap(err)
	}

	return allLayers, nil
}

// Get returns layer by id from the database.
func (layersDB *layersDB) Get(ctx context.Context, id uuid.UUID) (layers.Layer, error) {
	var layer layers.Layer
	query := `SELECT id, project_id, name, priority, none_weight
 	          FROM layers
 	          WHERE id=$1`

	err := layersDB.conn.QueryRowContext(ctx, query, id).Scan(&layer.ID,
		&layer.ProjectID, &layer.Name, &layer.Priority, &layer.NoneWeight)
	if errs.Is(err, sql.ErrNoRows) {
		return layer, layers.ErrNoLayer.New("")
	}

	return layer, ErrLayers.Wrap(err)
}

// UpdatePriorities updates priorities of layers in the database.
func (layersDB *layersDB) UpdatePriorities(ctx context.Context, layersToUpdate []layers.Layer) (err error) {
	query := `UPDATE layers 
	          SET priority=$1 
	          WHERE id = $2`

	tx, err := layersDB.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	for _, layer := range layersToUpdate {
		result, err := tx.ExecContext(ctx, query, layer.Priority, layer.ID)
		if err != nil {
			return errs.Combine(ErrLayers.Wrap(err), tx.Rollback())
		}

		rowNum, err := result.RowsAffected()
		if rowNum == 0 {
			return errs.Combine(layers.ErrNoLayer.New(""), tx.Rollback())
		}
		if err != nil {
			return errs.Combine(ErrLayers.Wrap(err), tx.Rollback())
		}
	}

	return ErrLayers.Wrap(tx.Commit())
}

// Delete deletes layer from the database.
func (layersDB *layersDB) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := layersDB.conn.ExecContext(ctx, "DELETE FROM layers WHERE id=$1", id)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return layers.ErrNoLayer.New("")
	}

	return ErrLayers.Wrap(err)
}

// ListWeights returns the set weights of the images of every project layer, by layer id and file name.
func (layersDB *layersDB) ListWeights(ctx context.Context, projectID uuid.UUID) (_ map[uuid.UUID]map[string]int, err error) {
	query := `SELECT w.layer_id, w.name, w.weight
	          FROM layer_weights w JOIN layers l ON l.id = w.layer_id
	          WHERE l.project_id = $1`

	rows, err := layersDB.conn.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	weights := make(map[uuid.UUID]map[string]int)
	for rows.Next() {
		var (
			layerID uuid.UUID
			name    string
			weight  int
		)
		if err = rows.Scan(&layerID, &name, &weight); err != nil {
			return nil, ErrLayers.Wrap(err)
		}
		if weights[layerID] == nil {
			weights[layerID] = make(map[string]int)
		}
		weights[layerID][name] = weight
	}

	return weights, ErrLayers.Wrap(rows.Err())
}

// UpdateWeights replaces the weight of the layer being absent and the weights of its images in one transaction.
func (layersDB *layersDB) UpdateWeights(ctx context.Context, layerID uuid.UUID, noneWeight int, weights map[string]int) (err error) {
	tx, err := layersDB.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrLayers.Wrap(err)
	}
	defer func() {
		if err != nil {
			err = errs.Combine(err, tx.Rollback())
		}
	}()

	result, err := tx.ExecContext(ctx, "UPDATE layers SET none_weight = $1 WHERE id = $2", noneWeight, layerID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}
	if rowNum, err := result.RowsAffected(); err != nil || rowNum == 0 {
		return errs.Combine(layers.ErrNoLayer.New(""), err)
	}

	if _, err = tx.ExecContext(ctx, "DELETE FROM layer_weights WHERE layer_id = $1", layerID); err != nil {
		return ErrLayers.Wrap(err)
	}
	for name, weight := range weights {
		_, err = tx.ExecContext(ctx, "INSERT INTO layer_weights(layer_id, name, weight) VALUES ($1, $2, $3)", layerID, name, weight)
		if err != nil {
			return ErrLayers.Wrap(err)
		}
	}

	return ErrLayers.Wrap(tx.Commit())
}

// RenameWeight moves the weight of a renamed image to its new name.
func (layersDB *layersDB) RenameWeight(ctx context.Context, layerID uuid.UUID, oldName, newName string) error {
	_, err := layersDB.conn.ExecContext(ctx, "UPDATE layer_weights SET name = $1 WHERE layer_id = $2 AND name = $3",
		newName, layerID, oldName)
	return ErrLayers.Wrap(err)
}

// DeleteWeight deletes the weight of a deleted image.
func (layersDB *layersDB) DeleteWeight(ctx context.Context, layerID uuid.UUID, name string) error {
	_, err := layersDB.conn.ExecContext(ctx, "DELETE FROM layer_weights WHERE layer_id = $1 AND name = $2", layerID, name)
	return ErrLayers.Wrap(err)
}
