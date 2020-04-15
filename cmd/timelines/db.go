package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zhughes3/kruonis/cmd/timelines/models"
)

type db struct {
	db *sql.DB
}

func NewDB(cfg *dbConfig) *sql.DB {
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.user, cfg.password, cfg.host, cfg.port, cfg.name, "disable")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database.")
	return db
}

func (db *db) insertTimelineGroup(title string) (*models.TimelineGroup, error) {
	var createdAt, updatedAt time.Time
	var id uint64
	sql := `INSERT INTO groups(title) VALUES($1) RETURNING id, created_at, updated_at;`

	err := db.db.QueryRow(sql, title).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error writing timeline group to DB")
		return nil, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)

	return &models.TimelineGroup{
		Id:        id,
		Title:     title,
		Timelines: nil,
		CreatedAt: created,
		UpdatedAt: updated,
	}, nil
}

func (db *db) insertTimeline(gid uint64, title string) (*models.Timeline, error) {
	var timeline models.Timeline
	var createdAt, updatedAt time.Time
	sql := `INSERT INTO timelines(group_id, title) VALUES($1, $2) RETURNING id, group_id, title, created_at, updated_at;`

	err := db.db.QueryRow(sql, gid, title).Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error writing timeline to DB")
		return nil, err
	}

	if timeline.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err

	}
	if timeline.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}
	return &timeline, nil
}

func (db *db) insertTimelineEvent(tid uint64, title, description, content string, timestamp *timestamp.Timestamp) (*models.TimelineEvent, error) {
	var createdAt, updatedAt time.Time
	var id uint64
	sql := `INSERT INTO events(timeline_id, title, timestamp, description, content) 
			VALUES($1, $2, $3, $4, $5) 
			RETURNING id, created_at, updated_at;
			`
	t, _ := convertTimestamp(timestamp)
	err := db.db.QueryRow(sql, tid, title, t, description, content).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error writing timeline event to DB")
		return nil, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)

	return &models.TimelineEvent{
		Id:          id,
		TimelineId:  tid,
		Title:       title,
		Timestamp:   timestamp,
		Description: description,
		Content:     content,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}, nil
}

func (db *db) insertTag(tag string, tid uint64) (uint64, error) {
	var createdAt, updatedAt time.Time
	var id uint64
	sql := `INSERT INTO tags(tag, timeline_id) VALUES($1, $2) RETURNING id;`

	err := db.db.QueryRow(sql, tag, tid).Scan(&id)
	if err != nil {
		fmt.Println("Error writing tag to DB")
		return 0, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)
	fmt.Println(created, updated)

	return id, nil
}

func (db *db) readTimelineGroup(id uint64) (*models.TimelineGroup, error) {
	var tg models.TimelineGroup
	var createdAt, updatedAt time.Time
	sql := `SELECT * from groups WHERE id = $1;`
	err := db.db.QueryRow(sql, id).Scan(&tg.Id, &tg.Title, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error reading timeline group from db")
		return nil, err
	}

	if tg.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}

	if tg.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	if tg.Timelines, err = db.readTimelinesWithGroupID(tg.GetId()); err != nil {
		return nil, err
	}

	return &tg, nil

}

func (db *db) readTimelinesWithGroupID(gid uint64) ([]*models.Timeline, error) {
	var timelines []*models.Timeline
	sql := `SELECT id, title, created_at, updated_at from timelines WHERE group_id = $1;`
	rows, err := db.db.Query(sql, gid)
	if err != nil {
		fmt.Println("Error reading timelines with group_id from db")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var timeline models.Timeline
		var createdAt, updatedAt time.Time

		err := rows.Scan(&timeline.Id, &timeline.Title, createdAt, updatedAt)
		if err != nil {
			fmt.Println("Error scanning timeline")
			return nil, err
		}

		if timeline.CreatedAt, err = convertTime(createdAt); err != nil {
			return nil, err
		}
		if timeline.UpdatedAt, err = convertTime(updatedAt); err != nil {
			return nil, err
		}

		if timeline.Tags, err = db.readTagsWithTimelineID(timeline.GetId()); err != nil {
			return nil, err
		}

		timelines = append(timelines, &timeline)
	}

	return timelines, nil
}

func (db *db) readTagsWithTimelineID(tid uint64) ([]string, error) {
	var tags []string
	sql := `SELECT tag FROM tags WHERE timeline_id = $1;`
	rows, err := db.db.Query(sql, tid)
	if err != nil {
		fmt.Println("Error reading tags with timeline_id from db")
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			fmt.Println("Error scanning tag from db")
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (db *db) readTimeline(id uint64) (*models.Timeline, error) {
	var timeline models.Timeline
	var createdAt, updatedAt time.Time

	sql := `SELECT id, title, created_at, updated_at FROM timelines WHERE id = $1;`

	err := db.db.QueryRow(sql, id).Scan(&timeline.Id, &timeline.Title, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Error reading timeline from db")
		return nil, err
	}

	if timeline.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if timeline.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	return &timeline, nil
}

func (db *db) readTimelineEvent(id uint64) (*models.TimelineEvent, error) {
	var timelineEvent models.TimelineEvent
	var timestamp, createdAt, updatedAt time.Time

	sql := `SELECT id, title, timestamp, description, content, created_at, updated_at
			FROM events
			WHERE timeline_id = $1
			`

	err := db.db.QueryRow(sql, id).Scan(&timelineEvent.Id, &timelineEvent.Title, timestamp, &timelineEvent.Description,
		&timelineEvent.Content, createdAt, updatedAt)
	if err != nil {
		fmt.Println("Error reading timeline event from db")
		return nil, err
	}

	if timelineEvent.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if timelineEvent.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}
	if timelineEvent.Timestamp, err = convertTime(timestamp); err != nil {
		return nil, err
	}

	return &timelineEvent, nil
}

func (db *db) deleteTimelineGroup(id uint64) error {
	sql := `DELETE FROM groups WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		fmt.Println("Error deleting group from DB")
		return err
	}
	return nil
}
func (db *db) deleteTimeline(id uint64) error {
	sql := `DELETE FROM timelines WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		fmt.Println("Error deleting timeline from DB")
		return err
	}
	return nil
}
func (db *db) deleteTimelineEvent(id uint64) error {
	sql := `DELETE FROM events WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		fmt.Println("Error deleting event from DB")
		return err
	}
	return nil
}
func (db *db) deleteTag(id uint64) error {
	sql := `DELETE FROM tags WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		fmt.Println("Error deleting tags from DB")
		return err
	}
	return nil
}

func convertTimestamp(t *timestamp.Timestamp) (time.Time, error) {
	ti, err := ptypes.Timestamp(t)
	if err != nil {
		fmt.Println("Error converting *timestamp.Timestamp to time.Time")
		return time.Time{}, err
	}

	return ti, nil
}
func convertTime(t time.Time) (*timestamp.Timestamp, error) {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		fmt.Println("Error converting time.Time to timestamp.Timestamp")
		return nil, err
	}
	return ts, nil
}
