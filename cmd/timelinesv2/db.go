package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

var errImageWriteError error = errors.New("Error writing image to blob store")

type (
	db struct {
		db *sql.DB
	}

	userWithHash struct {
		User
		hash string
	}

	privateDetails struct {
		UserID uint64
		IsPrivate bool
	}
)

func newDB(cfg *databaseConfig) *db {
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.user, cfg.password, cfg.host, cfg.port, cfg.name, "disable")
	database, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}
	log.WithFields(log.Fields{
		"User": cfg.user,
		"Host": cfg.host,
		"Port": cfg.port,
		"Name": cfg.name,
	}).Info("Connected to database")
	return &db{db: database}
}

func (db *db) readGroups() ([]*Group, error) {
	var groups []*Group
	sql := `SELECT * FROM groups`

	rows, err := db.db.Query(sql)
	if err != nil {
		log.Error("Error reading timeline groups from db")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.Title, &group.CreatedAt, &group.UpdatedAt, &group.Private, &group.UserId, &group.Uuid)
		if err != nil {
			log.Error("Error scanning timeline group")
			return nil, err
		}

		if group.Timelines, err = db.readTimelinesWithGroupID(group.Id); err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	return groups, nil
}
func (db *db) readUsers() ([]*User, error) {
	var users []*User
	sql := `SELECT id, email, created_at, updated_at, is_admin FROM users`

	rows, err := db.db.Query(sql)
	if err != nil {
		log.Error("Error reading users from db")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin)
		if err != nil {
			log.Error("Error scanning user")
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
func (db *db) readTimelines() ([]*Timeline, error) {
	var timelines []*Timeline
	sql := `SELECT * FROM timelines`

	rows, err := db.db.Query(sql)
	if err != nil {
		log.Error("Error reading timelines from db")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timeline Timeline
		err := rows.Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &timeline.CreatedAt, &timeline.UpdatedAt)
		if err != nil {
			log.Error("Error scanning timeline")
			return nil, err
		}

		if timeline.Tags, err = db.readTagsWithTimelineID(timeline.Id); err != nil {
			return nil, err
		}

		if timeline.Events, err = db.readEvents(timeline.Id); err != nil {
			return nil, err
		}

		timelines = append(timelines, &timeline)
	}

	return timelines, err
}
func (db *db) readGroup(id string) (*Group, error) {
	var group Group
	var sql string
	if uuidRegex.MatchString(id) {
		sql = `SELECT * from groups WHERE uuid = $1;`
	} else {
		sql = `SELECT * from groups WHERE id = $1;`
	}
	err := db.db.QueryRow(sql, id).Scan(&group.Id, &group.Title, &group.CreatedAt, &group.UpdatedAt, &group.Private, &group.UserId, &group.Uuid, &group.Views)
	if err != nil {
		log.Error("Error reading timeline group from db")
		return nil, err
	}

	if group.Timelines, err = db.readTimelinesWithGroupID(group.Id); err != nil {
		return nil, err
	}
	return &group, nil
}
func (db *db) readTimeline(id string) (*Timeline, error) {
	var timeline Timeline
	sql := `SELECT id, group_id, title, created_at, updated_at FROM timelines WHERE id = $1;`

	err := db.db.QueryRow(sql, id).Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &timeline.CreatedAt, &timeline.UpdatedAt)
	if err != nil {
		log.Error("Error reading timeline from db")
		return nil, err
	}

	if timeline.Tags, err = db.readTagsWithTimelineID(timeline.Id); err != nil {
		return nil, err
	}

	if timeline.Events, err = db.readEvents(timeline.Id); err != nil {
		return nil, err
	}

	return &timeline, nil
}
func (db *db) readEvent(id string) (*Event, error) {
	var event Event
	var imageURL sql.NullString
	sql := `SELECT id, timeline_id, title, timestamp, description, content, created_at, updated_at, image_url
			FROM events
			WHERE id = $1
			`
	err := db.db.QueryRow(sql, id).Scan(&event.Id, &event.TimelineId, &event.Title, &event.Timestamp, &event.Description, &event.Content, &event.CreatedAt, &event.UpdatedAt, &imageURL)
	if err != nil {
		log.Error("Error reading event from db")
		return nil, err
	}
	if imageURL.Valid {
		event.ImageUrl = imageURL.String
	}

	return &event, nil
}
func (db *db) readTimelineEvents(id string) ([]*Event, error) {
	var events []*Event
	var imageURL sql.NullString

	sql := `SELECT id, timeline_id, title, timestamp, description, content, created_at, updated_at, image_url
			FROM events
			WHERE timeline_id = $1
			`
	rows, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error reading events with timeline_id from db")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event Event
		rows.Scan(&event.Id, &event.TimelineId, &event.Title, &event.Timestamp, &event.Description, &event.Content, &event.CreatedAt, &event.UpdatedAt, &imageURL)
		if imageURL.Valid {
			event.ImageUrl = imageURL.String
		}
		events = append(events, &event)
	}

	return events, nil
}
func (db *db) readTimelinesWithGroupID(gid uint64) ([]*Timeline, error) {
	var timelines []*Timeline
	sql := `SELECT id, group_id, title, created_at, updated_at from timelines WHERE group_id = $1;`
	rows, err := db.db.Query(sql, gid)
	if err != nil {
		log.Error("Error reading timelines with group_id from db")
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var timeline Timeline
		err := rows.Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &timeline.CreatedAt, &timeline.UpdatedAt)
		if err != nil {
			log.Error("Error scanning timeline")
			return nil, err
		}
		if timeline.Tags, err = db.readTagsWithTimelineID(timeline.Id); err != nil {
			return nil, err
		}

		if timeline.Events, err = db.readEvents(timeline.Id); err != nil {
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
func (db *db) readEvents(tid uint64) ([]*Event, error) {
	var events []*Event
	var imageURL sql.NullString

	sql := `SELECT id, timeline_id, title, timestamp, description, content, created_at, updated_at, image_url
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
		var event Event
		rows.Scan(&event.Id, &event.TimelineId, &event.Title, &event.Timestamp, &event.Description, &event.Content,
			&event.CreatedAt, &event.UpdatedAt, &imageURL)
		if imageURL.Valid {
			event.ImageUrl = imageURL.String
		}
		events = append(events, &event)
	}
	return events, nil
}
func (db *db) readUserByEmail(email string) (*userWithHash, error) {
	var user User
	var hash string

	sql := `SELECT id, email, hash, created_at, updated_at, is_admin from users WHERE email = $1;`

	err := db.db.QueryRow(sql, email).Scan(&user.Id, &user.Email, &hash, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin)
	if err != nil {
		log.Error("Error reading user from db")
		return nil, err
	}

	return &userWithHash{
		User: user,
		hash: hash,
	}, nil
}
func (db *db) readUserByID(id uint64) (*User, error) {
	var user User
	sql := `SELECT id, email, created_at, updated_at, is_admin from users WHERE id = $1;`
	err := db.db.QueryRow(sql, id).Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin)
	if err != nil {
		log.Error("Error reading user from db")
		return nil, err
	}

	return &user, nil
}
func (db *db) readUserTimelineGroups(id uint64) ([]*Group, error) {
	var groups []*Group
	sql := `SELECT * FROM groups WHERE user_id = $1`
	rows, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error reading user timeline groups from db")
		return nil, err
	}
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.Id, &group.Title, &group.CreatedAt, &group.UpdatedAt, &group.Private, &group.UserId, &group.Uuid)
		if err != nil {
			log.Error("Error scanning group from db")
			return nil, err
		}
		if group.Timelines, err = db.readTimelinesWithGroupID(group.Id); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}
func (db *db) readImageUrlFromEvent(id string) (string, error) {
	var imageURL sql.NullString
	sql := `SELECT image_url FROM events WHERE id = $1`
	err := db.db.QueryRow(sql, id).Scan(&imageURL)
	if err != nil || !imageURL.Valid {
		log.Error("Error reading event image url from db.")
		return "", err
	}

	return imageURL.String, nil
}

func (db *db) insertGroup(title string, userID uint64, isPrivate bool) (*Group, error) {
	var group Group
	uid := uuid.New()
	sql := `INSERT INTO groups(title, private, user_id, uuid) VALUES($1, $2, $3, $4) RETURNING id, title, created_at, updated_at, private, user_id, uuid;`
	err := db.db.QueryRow(sql, title, isPrivate, userID, uid.String()).Scan(&group.Id, &group.Title, &group.CreatedAt, &group.UpdatedAt, &group.Private, &group.UserId, &group.Uuid)
	if err != nil {
		log.Error("Error writing timeline group to DB")
		return nil, err
	}

	return &group, nil
}
func (db *db) insertTimeline(gid uint64, title string) (*Timeline, error) {
	var timeline Timeline
	sql := `INSERT INTO timelines(group_id, title)
			SELECT *
			FROM (SELECT $1::integer, $2) x
			WHERE (SELECT COUNT(*) FROM timelines WHERE group_id = $1::integer) < 2
			RETURNING id, group_id, title, created_at, updated_at;`

	err := db.db.QueryRow(sql, gid, title).Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &timeline.CreatedAt, &timeline.UpdatedAt)
	if err != nil {
		log.Error("Error writing timeline to DB")
		return nil, err
	}
	return &timeline, nil
}
func (db *db) insertTimelineEvent(tid, title, description, content string, timestamp time.Time) (*Event, error) {
	var event Event
	var imageURL sql.NullString
	sql := `INSERT INTO events(timeline_id, title, timestamp, description, content) 
			VALUES($1, $2, $3, $4, $5) 
			RETURNING id, timeline_id, title, timestamp, description, content, created_at, updated_at, image_url;
			`
	err := db.db.QueryRow(sql, tid, title, timestamp, description, content).Scan(&event.Id, &event.TimelineId, &event.Title, &event.Timestamp, &event.Description, &event.Content, &event.CreatedAt, &event.UpdatedAt, &imageURL)
	if err != nil {
		log.Error("Error writing timeline event to DB")
		return nil, err
	}
	if imageURL.Valid {
		event.ImageUrl = imageURL.String
	}
	return &event, nil
}
func (db *db) insertTag(tag string, tid uint64) (uint64, error) {
	var id uint64
	sql := `INSERT INTO tags(tag, timeline_id) VALUES($1, $2) RETURNING id;`

	err := db.db.QueryRow(sql, tag, tid).Scan(&id)
	if err != nil {
		log.Error("Error writing tag to DB")
		return 0, err
	}
	return id, nil
}
func (db *db) insertUser(username, hash string) (*User, error) {
	var user User
	sql := `INSERT INTO users(email, hash) VALUES($1, $2) RETURNING id, email, created_at, updated_at, is_admin;`
	err := db.db.QueryRow(sql, username, hash).Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin)
	if err != nil {
		log.Error("Error writing user to db")
		return nil, err
	}

	return &user, nil
}

func (db *db) updateEvent(id, title, description, content string, timestamp time.Time) (*Event, error) {
	var event Event
	var imageURL sql.NullString
	sql := `UPDATE events
			SET title = $1, timestamp = $2, description = $3, content = $4, updated_at = $5
			WHERE id = $6
			RETURNING id, timeline_id, title, timestamp, description, content, created_at, updated_at, image_url;`

	now := time.Now()

	err := db.db.QueryRow(sql, title, timestamp, description, content, now, id).
		Scan(&event.Id, &event.TimelineId, &event.Title, &event.Timestamp, &event.Description, &event.Content, &event.CreatedAt, &event.UpdatedAt, &imageURL)

	if err != nil {
		log.Error("Error updating timeline event db")
		return nil, err
	}
	if imageURL.Valid {
		event.ImageUrl = imageURL.String
	}

	return &event, nil
}
func (db *db) updateGroup(id, title string, isPrivate bool) (*Group, error) {
	var group Group
	sql := `UPDATE groups 
		SET title = $1, updated_at = $2, private = $3 
		WHERE id = $4 
		RETURNING id, title, created_at, updated_at, private, user_id, uuid
	`
	err := db.db.QueryRow(sql, title, time.Now(), isPrivate, id).
		Scan(&group.Id, &group.Title, &group.CreatedAt, &group.UpdatedAt, &group.Private, &group.UserId, &group.Uuid)
	if err != nil {
		log.Error("Error updating timeline group")
		return nil, err
	}

	if group.Timelines, err = db.readTimelinesWithGroupID(group.Id); err != nil {
		return nil, err
	}

	return &group, nil
}
func (db *db) updateTimeline(id, title string, tags []string) (*Timeline, error) {
	var timeline Timeline
	timelineSql := `UPDATE timelines 
		SET title = $1, updated_at = $2 WHERE id = $3 
		RETURNING id, group_id, title, created_at, updated_at
	`
	err := db.db.QueryRow(timelineSql, title, time.Now(), id).
		Scan(&timeline.Id, &timeline.GroupId, &timeline.Title, &timeline.CreatedAt, &timeline.UpdatedAt)
	if err != nil {
		log.Error("Error updating timeline")
		return nil, err
	}

	if tags != nil {
		_, err := db.db.Query("DELETE FROM tags WHERE timeline_id = $1", id)
		if err != nil {
			log.Error("Error deleting tags during update timeline")
			return &timeline, err
		}
		for _, tag := range tags {
			db.insertTag(tag, timeline.Id)
		}
		timeline.Tags = tags
	} else {
		if timeline.Tags, err = db.readTagsWithTimelineID(timeline.Id); err != nil {
			log.Error("Error reading tags during update timeline")
			return &timeline, err
		}
	}

	if timeline.Events, err = db.readTimelineEvents(id); err != nil {
		log.Error("Error reading events during update timeline")
		return &timeline, err
	}

	return &timeline, nil
}
func (db *db) updateTimelineEventWithImageURL(id, url string) (string, error) {
	var imageURL sql.NullString
	sql := `UPDATE events
			SET image_url = $1, updated_at = $2
			WHERE id = $3
			RETURNING image_url;`
	now := time.Now()
	err := db.db.QueryRow(sql, url, now, id).Scan(&imageURL)
	if err != nil {
		log.Error("Error updating timeline event with image URL")
		return "", err
	}
	if imageURL.Valid {
		return imageURL.String, nil
	} else {
		return "", errImageWriteError
	}
}

func (db *db) deleteTimeline(id string) error {
	sql := `DELETE FROM timelines WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting timeline from DB")
		return err
	}
	return nil
}
func (db *db) deleteEvent(id string) error {
	sql := `DELETE FROM events WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting event from DB")
		return err
	}
	return nil
}
func (db *db) deleteGroup(id string) error {
	sql := `DELETE FROM groups WHERE id = $1;`
	_, err := db.db.Query(sql, id)
	if err != nil {
		log.Error("Error deleting group from DB")
		return err
	}
	return nil
}
func (db *db) deleteImageUrlFromEvent(id string) error {
	sql := `UPDATE events SET image_url = null WHERE id = $1`
	db.db.QueryRow(sql, id)
	return nil
}
// func (db *db) readEventPrivateDetails(id string) (privateDetails, error) {
// 	var p privateDetails
// 	sql := `
// 		SELECT groups.user_id, groups.private 
// 		FROM events 
// 		FULL JOIN timelines ON events.timeline_id = timelines.id 
// 		FULL JOIN groups ON timelines.group_id = groups.id 
// 		WHERE events.id = $1;`

// 	err := db.db.QueryRow(sql, id).Scan(&p.UserID, &p.IsPrivate)
// 	if err != nil {
// 		log.Error("Error selecting timeline events private details")
// 	}

// 	return p, err
// }
// func (db *db) readTimelinePrivateDetails(id string) (privateDetails, error) {
// 	var p privateDetails
// 	sql := `
// 		SELECT groups.user_id, groups.private
// 		FROM timelines
// 		FULL JOIN groups on timelines.group_id = groups.id
// 		WHERE timelines.id = $1`

// 	err := db.db.QueryRow(sql, id).Scan(&p.UserID, &p.IsPrivate)
// 	if err != nil {
// 		log.Error("Error selecting timelines private details")
// 	}

// 	return p, err
// }

func (db *db) isEventPrivate(id string) (bool, error) {
	var b bool
	sql := `
		SELECT groups.private 
		FROM events 
		FULL JOIN timelines ON events.timeline_id = timelines.id 
		FULL JOIN groups ON timelines.group_id = groups.id 
		WHERE events.id = $1;`

	err := db.db.QueryRow(sql, id).Scan(&b)
	if err != nil {
		log.Error("Error selecting timeline events private details")
	}

	return b, err
}

func (db *db) isTimelinePrivate(id string) (bool, error) {
	var b bool
	sql := `
		SELECT groups.private
		FROM timelines
		FULL JOIN groups on timelines.group_id = groups.id
		WHERE timelines.id = $1`

	err := db.db.QueryRow(sql, id).Scan(&b)
	if err != nil {
		log.Error("Error selecting timelines private details")
	}

	return b, err
}

func (db *db) incrementGroupViews(id string) error {
	ctx := context.Background()
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE groups SET views = views + 1 WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *db) insertGroupView(id string) error {
	sql := `INSERT INTO group_views (id) VALUES ($1);`
	_, err := db.db.Query(sql, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *db) getTrendingGroups() error {
	//TODO
	sql := `select id, count(*) as views
	from group_views
	where time between date_trunc('day', time) and date_trunc('day', time) + '24 hours'
	group by id
	order by views desc`
	return nil
}