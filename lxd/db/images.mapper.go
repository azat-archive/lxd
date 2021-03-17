//go:build linux && cgo && !agent
// +build linux,cgo,!agent

package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"database/sql"
	"fmt"
	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
	"github.com/pkg/errors"
)

var _ = api.ServerEnvironment{}

var imageObjects = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProject = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE project = ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndCached = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE project = ? AND images.cached = ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndPublic = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE project = ? AND images.public = ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndFingerprint = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE project = ? AND images.fingerprint LIKE ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndFingerprintAndPublic = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE project = ? AND images.fingerprint LIKE ? AND images.public = ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByFingerprint = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE images.fingerprint LIKE ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByCached = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE images.cached = ? ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByAutoUpdate = cluster.RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images JOIN projects ON images.project_id = projects.id
  WHERE images.auto_update = ? ORDER BY projects.id, images.fingerprint
`)

// GetImages returns all available images.
func (c *ClusterTx) GetImages(filter ImageFilter) ([]Image, error) {
	// Result slice.
	objects := make([]Image, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Fingerprint != "" {
		criteria["Fingerprint"] = filter.Fingerprint
	}
	if filter.Public != false {
		criteria["Public"] = filter.Public
	}
	if filter.Cached != false {
		criteria["Cached"] = filter.Cached
	}
	if filter.AutoUpdate != false {
		criteria["AutoUpdate"] = filter.AutoUpdate
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Fingerprint"] != nil && criteria["Public"] != nil {
		stmt = c.stmt(imageObjectsByProjectAndFingerprintAndPublic)
		args = []interface{}{
			filter.Project,
			filter.Fingerprint,
			filter.Public,
		}
	} else if criteria["Project"] != nil && criteria["Public"] != nil {
		stmt = c.stmt(imageObjectsByProjectAndPublic)
		args = []interface{}{
			filter.Project,
			filter.Public,
		}
	} else if criteria["Project"] != nil && criteria["Fingerprint"] != nil {
		stmt = c.stmt(imageObjectsByProjectAndFingerprint)
		args = []interface{}{
			filter.Project,
			filter.Fingerprint,
		}
	} else if criteria["Project"] != nil && criteria["Cached"] != nil {
		stmt = c.stmt(imageObjectsByProjectAndCached)
		args = []interface{}{
			filter.Project,
			filter.Cached,
		}
	} else if criteria["Project"] != nil {
		stmt = c.stmt(imageObjectsByProject)
		args = []interface{}{
			filter.Project,
		}
	} else if criteria["Fingerprint"] != nil {
		stmt = c.stmt(imageObjectsByFingerprint)
		args = []interface{}{
			filter.Fingerprint,
		}
	} else if criteria["Cached"] != nil {
		stmt = c.stmt(imageObjectsByCached)
		args = []interface{}{
			filter.Cached,
		}
	} else if criteria["AutoUpdate"] != nil {
		stmt = c.stmt(imageObjectsByAutoUpdate)
		args = []interface{}{
			filter.AutoUpdate,
		}
	} else {
		stmt = c.stmt(imageObjects)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Image{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].Project,
			&objects[i].Fingerprint,
			&objects[i].Type,
			&objects[i].Filename,
			&objects[i].Size,
			&objects[i].Public,
			&objects[i].Architecture,
			&objects[i].CreationDate,
			&objects[i].ExpiryDate,
			&objects[i].UploadDate,
			&objects[i].Cached,
			&objects[i].LastUseDate,
			&objects[i].AutoUpdate,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch images")
	}

	return objects, nil
}

// GetImage returns the image with the given key.
func (c *ClusterTx) GetImage(project string, fingerprint string) (*Image, error) {
	filter := ImageFilter{}
	filter.Project = project
	filter.Fingerprint = fingerprint

	objects, err := c.GetImages(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch Image")
	}

	switch len(objects) {
	case 0:
		return nil, ErrNoSuchObject
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one image matches")
	}
}
