package main

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

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
	log.WithFields(log.Fields{
		"User": cfg.user,
		"Host": cfg.host,
		"Port": cfg.port,
		"Name": cfg.name,
	}).Info("Connected to database")
	return db
}

func (db *db) insertTimelineGroup(title string) (*models.TimelineGroup, error) {
	var createdAt, updatedAt time.Time
	var id uint64
	sql := `INSERT INTO groups(title) VALUES($1) RETURNING id, created_at, updated_at;`

	err := db.db.QueryRow(sql, title).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error writing timeline group to DB")
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
		log.Error("Error writing timeline to DB")
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
		log.Error("Error writing timeline event to DB")
		return nil, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)

	return &models.TimelineEvent{
		Id:          tid,
		EventId:     id,
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
		log.Error("Error writing tag to DB")
		return 0, err
	}

	created, _ := convertTime(createdAt)
	updated, _ := convertTime(updatedAt)
	fmt.Println(created, updated)

	return id, nil
}

func (db *db) updateTimelineEvent(eventID uint64, title, description, content string, t *timestamp.Timestamp) (*models.TimelineEvent, error) {
	event := &models.TimelineEvent{}
	sql := `UPDATE events
			SET title = $1, timestamp = $2, description = $3, content = $4, updated_at = $5
			WHERE id = $6
			RETURNING id, timeline_id, title, timestamp, description, content, created_at, updated_at;`

	now := time.Now()
	timestamp, _ := convertTimestamp(t)
	var ts, createdAt, updatedAt time.Time
	err := db.db.QueryRow(sql, title, timestamp, description, content, now, eventID).
		Scan(&event.EventId, &event.Id, &event.Title, &ts, &event.Description, &event.Content, &createdAt, &updatedAt)

	if err != nil {
		log.Error("Error updating timeline event db")
		return nil, err
	}

	if event.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if event.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}
	if event.Timestamp, err = convertTime(ts); err != nil {
		return nil, err
	}

	return event, nil
}
func (db *db) updateTimeline(id uint64, title string, tags []string) (*models.Timeline, error) {
	timeline := &models.Timeline{}

	timelineSql := `UPDATE timelines 
		SET title = $1, updated_at = $2 WHERE id = $3 
		RETURNING id, group_id, title, created_at, updated_at
	`

	var createdAt, updatedAt time.Time

	err := db.db.QueryRow(timelineSql, title, time.Now(), id).
		Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error updating timeline")
		return nil, err
	}

	if timeline.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if timeline.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	if tags != nil {
		_, err := db.db.Query("DELETE FROM tags WHERE timeline_id = $1", id)
		if err != nil {
			log.Error("Error deleting tags during update timeline")
			return timeline, err
		}

		for _, tag := range tags {
			db.insertTag(tag, id)
		}

		timeline.Tags = tags
	} else {
		if timeline.Tags, err = db.readTagsWithTimelineID(id); err != nil {
			log.Error("Error reading tags during update timeline")
			return timeline, err
		}
	}

	if timeline.Events, err = db.readTimelineEvents(id); err != nil {
		log.Error("Error reading events during update timeline")
		return timeline, err
	}

	return timeline, nil
}
func (db *db) updateTimelineGroup(id uint64, title string) (*models.TimelineGroup, error) {
	group := &models.TimelineGroup{}

	sql := `UPDATE groups 
		SET title = $1, updated_at = $2 WHERE id = $3 
		RETURNING id, title, created_at, updated_at
	`

	var createdAt, updatedAt time.Time

	err := db.db.QueryRow(sql, title, time.Now(), id).
		Scan(&group.Id, &group.Title, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error updating timeline")
		return nil, err
	}

	if group.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if group.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	if group.Timelines, err = db.readTimelinesWithGroupID(id); err != nil {
		return nil, err
	}

	return group, nil
}

func (db *db) readTimelineGroup(id uint64) (*models.TimelineGroup, error) {
	var tg models.TimelineGroup
	var createdAt, updatedAt time.Time
	sql := `SELECT * from groups WHERE id = $1;`
	err := db.db.QueryRow(sql, id).Scan(&tg.Id, &tg.Title, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error reading timeline group from db")
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
	sql := `SELECT id, group_id, title, created_at, updated_at from timelines WHERE group_id = $1;`
	rows, err := db.db.Query(sql, gid)
	if err != nil {
		log.Error("Error reading timelines with group_id from db")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var timeline models.Timeline
		var createdAt, updatedAt time.Time

		err := rows.Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &createdAt, &updatedAt)
		if err != nil {
			log.Error("Error scanning timeline")
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

		if timeline.Events, err = db.readTimelineEvents(timeline.GetId()); err != nil {
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
		log.Error("Error reading tags with timeline_id from db")
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			log.Error("Error scanning tag from db")
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
func (db *db) readTimeline(id uint64) (*models.Timeline, error) {
	var timeline models.Timeline
	var createdAt, updatedAt time.Time

	sql := `SELECT id, group_id, title, created_at, updated_at FROM timelines WHERE id = $1;`

	err := db.db.QueryRow(sql, id).Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error reading timeline from db")
		return nil, err
	}

	if timeline.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if timeline.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}

	if timeline.Tags, err = db.readTagsWithTimelineID(id); err != nil {
		return nil, err
	}

	if timeline.Events, err = db.readTimelineEvents(id); err != nil {
		return nil, err
	}

	return &timeline, nil
}
func (db *db) readTimelineEvents(tid uint64) ([]*models.TimelineEvent, error) {
	var events []*models.TimelineEvent

	sql := `SELECT id, timeline_id, title, timestamp, description, content, created_at, updated_at
			FROM events
			WHERE timeline_id = $1
			`
	rows, err := db.db.Query(sql, tid)
	if err != nil {
		log.Error("Error reading events with timeline_id from db")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.TimelineEvent
		var timestamp, createdAt, updatedAt time.Time
		rows.Scan(&event.EventId, &event.Id, &event.Title, &timestamp, &event.Description, &event.Content, &createdAt, &updatedAt)
		if event.CreatedAt, err = convertTime(createdAt); err != nil {
			return nil, err
		}
		if event.UpdatedAt, err = convertTime(updatedAt); err != nil {
			return nil, err
		}
		if event.Timestamp, err = convertTime(timestamp); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}
func (db *db) readTimelineEvent(id uint64) (*models.TimelineEvent, error) {
	var event *models.TimelineEvent
	var timestamp, createdAt, updatedAt time.Time
	sql := `SELECT id, timeline_id, title, timestamp, description, content, created_at, updated_at
			FROM events
			WHERE id = $1
			`
	err := db.db.QueryRow(sql, id).Scan(&event.EventId, &event.Id, &event.Title, &timestamp, &event.Description, &event.Content, &createdAt, &updatedAt)
	if err != nil {
		log.Error("Error reading event from db")
		return nil, err
	}
	if event.CreatedAt, err = convertTime(createdAt); err != nil {
		return nil, err
	}
	if event.UpdatedAt, err = convertTime(updatedAt); err != nil {
		return nil, err
	}
	if event.Timestamp, err = convertTime(timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func (db *db) deleteTimelineGroup(id uint64) error {
	sql := `DELETE FROM groups WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting group from DB")
		return err
	}
	return nil
}
func (db *db) deleteTimeline(id uint64) error {
	sql := `DELETE FROM timelines WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting timeline from DB")
		return err
	}
	return nil
}
func (db *db) deleteTimelineEvent(id uint64) error {
	sql := `DELETE FROM events WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting event from DB")
		return err
	}
	return nil
}
func (db *db) deleteTag(id uint64) error {
	sql := `DELETE FROM tags WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting tags from DB")
		return err
	}
	return nil
}

func convertTimestamp(t *timestamp.Timestamp) (time.Time, error) {
	ti, err := ptypes.Timestamp(t)
	if err != nil {
		log.Error("Error converting *timestamp.Timestamp to time.Time")
		return time.Time{}, err
	}

	return ti, nil
}
func convertTime(t time.Time) (*timestamp.Timestamp, error) {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		log.Error("Error converting time.Time to *timestamp.Timestamp")
		return nil, err
	}
	return ts, nil
}
