package plot

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/goleveldb/leveldb"
	"os"
)

// DB handles access to the plots leveldb database. It provides abstraction over the database layer so that
// plots may be directly read from it.
type DB struct {
	ldb      *leveldb.DB
	settings Settings
	cache    map[Position]*Plot
}

// OpenDB opens the directory passed as a leveldb database for plots. If the directory does not yet exist, it
// is created.
// If successful, a new DB is returned which may be used to read and write plots.
func OpenDB(dir string, settings Settings) (*DB, error) {
	// Always try to create the directory. If it doesn't work, we've probably created the directory already,
	// and that's fine.
	_ = os.MkdirAll(dir, 0777)

	ldb, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening leveldb database: %w", err)
	}
	return &DB{ldb: ldb, settings: settings, cache: map[Position]*Plot{}}, nil
}

// Plot attempts to read a Plot from the DB at the Position passed.
func (db *DB) Plot(pos Position) (*Plot, error) {
	if p, ok := db.cache[pos]; ok {
		return p, nil
	}
	val, err := db.ldb.Get(pos.Hash(), nil)
	if err != nil {
		return nil, fmt.Errorf("plot: %w", err)
	}
	var p Plot
	if err := json.Unmarshal(val, &p); err != nil {
		return nil, fmt.Errorf("plot: %w", err)
	}
	db.cache[pos] = &p
	return &p, nil
}

// StorePlot attempts to store a Plot at a specific Position in the DB.
func (db *DB) StorePlot(pos Position, p *Plot) error {
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("store plot: %w", err)
	}
	if err := db.ldb.Put(pos.Hash(), b, nil); err != nil {
		return fmt.Errorf("store plot: %w", err)
	}
	db.cache[pos] = p
	return nil
}

// RemovePlot attempts to remove a Plot at a specific Position in the DB.
func (db *DB) RemovePlot(pos Position) error {
	if err := db.ldb.Delete(pos.Hash(), nil); err != nil {
		return fmt.Errorf("remove plot: %w", err)
	}
	delete(db.cache, pos)
	return nil
}

// PlayerPlots attempts to read a list of Positions from the DB for the player.Player passed.
func (db *DB) PlayerPlots(p *player.Player) ([]Position, error) {
	uuid := p.UUID()
	val, err := db.ldb.Get(uuid[:], nil)
	if err != nil {
		return nil, fmt.Errorf("player plots: %w", err)
	}
	var positions []Position
	if err := json.Unmarshal(val, &positions); err != nil {
		return nil, fmt.Errorf("player plots: %w", err)
	}
	return positions, nil
}

// StorePlayerPlots attempts to store the Positions of the plots it owns into the DB.
func (db *DB) StorePlayerPlots(p *player.Player, positions []Position) error {
	val, err := json.Marshal(positions)
	if err != nil {
		return fmt.Errorf("store player plots: %w", err)
	}
	uuid := p.UUID()
	if err := db.ldb.Put(uuid[:], val, nil); err != nil {
		return fmt.Errorf("store player plots: %w", err)
	}
	return nil
}

// Close closes the underlying leveldb database.
func (db *DB) Close() error {
	return db.ldb.Close()
}
