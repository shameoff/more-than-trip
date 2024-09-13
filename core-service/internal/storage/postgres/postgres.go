// Здесь описываются запросы в БД

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/shameoff/more-than-trip/core/internal/domain/models"
	// "github.com/shameoff/more-than-trip/core/internal/jaeger"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SavePhoto(ctx context.Context, data models.Photo) error {
	// Сохранение информации о фотографии в базу данных
	query := `
        INSERT INTO photo (coords, description, img_url, place, region_id, trip_id, user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := s.db.ExecContext(ctx, query, data.Coords, data.Description, data.ImgUrl, data.Place, data.RegionId, data.TripId, data.UserId)
	if err != nil {
		return fmt.Errorf("failed to save photo data: %w", err)
	}

	return nil
}

func (s *Storage) GetPhoto(ctx context.Context, photoId uuid.UUID) (models.Photo, error) {
	var photo models.Photo

	query := `
        SELECT id, coords, description, img_url, place, region_id, trip_id, user_id
        FROM photo
        WHERE id = $1
    `
	row := s.db.QueryRowContext(ctx, query, photoId)

	err := row.Scan(&photo.Id, &photo.Coords, &photo.Description, &photo.ImgUrl, &photo.Place, &photo.RegionId, &photo.TripId, &photo.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return photo, fmt.Errorf("photo not found")
		}
		return photo, fmt.Errorf("failed to get photo: %w", err)
	}

	return photo, nil
}

func (s *Storage) GetPhotos(ctx context.Context, filters models.PhotoFiltersDTO) ([]models.Photo, error) {
	var photos []models.Photo

	// Строим запрос в зависимости от фильтров
	query := `SELECT id, coords, description, img_url, place, region_id, trip_id, user_id FROM photo WHERE 1=1`
	args := []interface{}{}

	if filters.RegionId != uuid.Nil {
		query += " AND region_id = $1"
		args = append(args, filters.RegionId)
	}

	if filters.TagKey != "" {
		query += " AND id IN (SELECT photo_id FROM photo_tags WHERE tag_id = $2)"
		args = append(args, filters.TagKey)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Id, &photo.Coords, &photo.Description, &photo.ImgUrl, &photo.Place, &photo.RegionId, &photo.TripId, &photo.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (s *Storage) UpdatePhoto(ctx context.Context, photoId uuid.UUID, data models.Photo) error {
	query := `
        UPDATE photo
        SET coords = $1, description = $2, img_url = $3, place = $4, region_id = $5, trip_id = $6, user_id = $7
        WHERE id = $8
    `
	_, err := s.db.ExecContext(ctx, query, data.Coords, data.Description, data.ImgUrl, data.Place, data.RegionId, data.TripId, data.UserId, photoId)
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

func (s *Storage) DeletePhoto(ctx context.Context, photoId uuid.UUID) error {
	query := `DELETE FROM photo WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, photoId)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}
	return nil
}
func (s *Storage) GetPhotosByTripId(ctx context.Context, tripId uuid.UUID) ([]models.Photo, error) {
	var photos []models.Photo

	query := `
        SELECT id, coords, description, img_url, place, region_id, trip_id, user_id
        FROM photo
        WHERE trip_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, tripId)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos by trip id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Id, &photo.Coords, &photo.Description, &photo.ImgUrl, &photo.Place, &photo.RegionId, &photo.TripId, &photo.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (s *Storage) GetPhotosByUserId(ctx context.Context, userId uuid.UUID) ([]models.Photo, error) {
	var photos []models.Photo

	query := `
        SELECT id, coords, description, img_url, place, region_id, trip_id, user_id
        FROM photo
        WHERE user_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos by user id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Id, &photo.Coords, &photo.Description, &photo.ImgUrl, &photo.Place, &photo.RegionId, &photo.TripId, &photo.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (s *Storage) GetPhotosByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Photo, error) {
	var photos []models.Photo

	query := `
        SELECT id, coords, description, img_url, place, region_id, trip_id, user_id
        FROM photo
        WHERE region_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, regionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos by region id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Id, &photo.Coords, &photo.Description, &photo.ImgUrl, &photo.Place, &photo.RegionId, &photo.TripId, &photo.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (s *Storage) LikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error {
	query := `
        INSERT INTO likes (photo_id, user_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `
	_, err := s.db.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return fmt.Errorf("failed to like photo: %w", err)
	}
	return nil
}

func (s *Storage) DislikePhoto(ctx context.Context, photoId uuid.UUID, userId uuid.UUID) error {
	query := `
        DELETE FROM likes WHERE photo_id = $1 AND user_id = $2
    `
	_, err := s.db.ExecContext(ctx, query, photoId, userId)
	if err != nil {
		return fmt.Errorf("failed to dislike photo: %w", err)
	}
	return nil
}

func (s *Storage) CreateTag(ctx context.Context, tag models.Tag) error {
	query := `
        INSERT INTO tag (name)
        VALUES ($1, $2)
    `
	_, err := s.db.ExecContext(ctx, query, tag.Id, tag.Name)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (s *Storage) DeleteTag(ctx context.Context, tagId string) error {
	query := `DELETE FROM tag WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, tagId)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (s *Storage) GetTags(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag

	query := `SELECT id, name FROM tag`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) GetTagsByPhotoId(ctx context.Context, photoId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag

	query := `
        SELECT t.id, t.name
        FROM tag t
        INNER JOIN photo_tags pt ON t.id = pt.tag_id
        WHERE pt.photo_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, photoId)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by photo id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) CreateRegion(ctx context.Context, region models.Region) error {
	query := `
        INSERT INTO region (id, name, country)
        VALUES ($1, $2, $3)
    `
	_, err := s.db.ExecContext(ctx, query, region.Id, region.Name, region.Country)
	if err != nil {
		return fmt.Errorf("failed to create region: %w", err)
	}
	return nil
}
func (s *Storage) DeleteRegion(ctx context.Context, regionId uuid.UUID) error {
	query := `DELETE FROM region WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, regionId)
	if err != nil {
		return fmt.Errorf("failed to delete region: %w", err)
	}
	return nil
}
func (s *Storage) GetRegionById(ctx context.Context, regionId uuid.UUID) (models.Region, error) {
	var region models.Region

	query := `SELECT id, name, country FROM region WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, regionId).Scan(&region.Id, &region.Name, &region.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			return region, fmt.Errorf("region not found")
		}
		return region, fmt.Errorf("failed to get region: %w", err)
	}

	return region, nil
}

func (s *Storage) GetRegionByKey(ctx context.Context, regionKey string) (models.Region, error) {
	var region models.Region

	query := `SELECT id, name, country FROM region WHERE name = $1`
	err := s.db.QueryRowContext(ctx, query, regionKey).Scan(&region.Id, &region.Name, &region.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			return region, fmt.Errorf("region not found")
		}
		return region, fmt.Errorf("failed to get region by key: %w", err)
	}

	return region, nil
}
func (s *Storage) GetRegions(ctx context.Context) ([]models.Region, error) {
	var regions []models.Region

	query := `SELECT id, name, country FROM region`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var region models.Region
		err := rows.Scan(&region.Id, &region.Name, &region.Country)
		if err != nil {
			return nil, fmt.Errorf("failed to scan region: %w", err)
		}
		regions = append(regions, region)
	}

	return regions, nil
}
func (s *Storage) UpdateRegion(ctx context.Context, regionId uuid.UUID, data models.Region) error {
	query := `
        UPDATE region
        SET name = $1, country = $2
        WHERE id = $3
    `
	_, err := s.db.ExecContext(ctx, query, data.Name, data.Country, regionId)
	if err != nil {
		return fmt.Errorf("failed to update region: %w", err)
	}
	return nil
}

func (s *Storage) CreateTrip(ctx context.Context, trip models.Trip) error {
	query := `
        INSERT INTO trip (name, description, region_id)
        VALUES ($1, $2, $3, $4)
    `
	_, err := s.db.ExecContext(ctx, query, trip.Name, trip.Description, trip.RegionId)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}
	return nil
}

func (s *Storage) DeleteTrip(ctx context.Context, tripId uuid.UUID) error {
	query := `DELETE FROM trip WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, tripId)
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}
	return nil
}
func (s *Storage) GetTripById(ctx context.Context, tripId uuid.UUID) (models.Trip, error) {
	var trip models.Trip

	query := `
        SELECT id, name, description, region_id
        FROM trip
        WHERE id = $1
    `
	row := s.db.QueryRowContext(ctx, query, tripId)

	err := row.Scan(&trip.Id, &trip.Name, &trip.Description, &trip.RegionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return trip, fmt.Errorf("trip not found")
		}
		return trip, fmt.Errorf("failed to get trip: %w", err)
	}

	return trip, nil
}
func (s *Storage) GetTripsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Trip, error) {
	var trips []models.Trip

	query := `
        SELECT t.id, t.name, t.description, t.region_id
        FROM trip t
        INNER JOIN user_trips ut ON t.id = ut.trip_id
        WHERE ut.user_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips by user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(&trip.Id, &trip.Name, &trip.Description, &trip.RegionId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trip: %w", err)
		}
		trips = append(trips, trip)
	}

	return trips, nil
}
func (s *Storage) GetTripsByRegionId(ctx context.Context, regionId uuid.UUID) ([]models.Trip, error) {
	var trips []models.Trip

	query := `
        SELECT id, name, description, region_id
        FROM trip
        WHERE region_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, regionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips by region: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(&trip.Id, &trip.Name, &trip.Description, &trip.RegionId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trip: %w", err)
		}
		trips = append(trips, trip)
	}

	return trips, nil
}
func (s *Storage) GetTripsByTag(ctx context.Context, tagId string) ([]models.Trip, error) {
	var trips []models.Trip

	query := `
        SELECT t.id, t.name, t.description, t.region_id
        FROM trip t
        INNER JOIN trip_tags tt ON t.id = tt.trip_id
        WHERE tt.tag_id = $1
    `
	rows, err := s.db.QueryContext(ctx, query, tagId)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips by tag: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(&trip.Id, &trip.Name, &trip.Description, &trip.RegionId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trip: %w", err)
		}
		trips = append(trips, trip)
	}

	return trips, nil
}

func (s *Storage) UpdateTrip(ctx context.Context, tripId uuid.UUID, data models.Trip) error {
	query := `
        UPDATE trip
        SET name = $1, description = $2, region_id = $3
        WHERE id = $4
    `
	_, err := s.db.ExecContext(ctx, query, data.Name, data.Description, data.RegionId, tripId)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, user models.User) error {
	query := `
        INSERT INTO users (id, username, full_name, birth_date, education, city)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := s.db.ExecContext(ctx, query, user.Id, user.UserName, user.FullName, user.BirthDate, user.Education, user.City)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
func (s *Storage) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *Storage) GetUserById(ctx context.Context, userId uuid.UUID) (models.User, error) {
	var user models.User

	query := `
        SELECT id, username, full_name, birth_date, education, city
        FROM users
        WHERE id = $1
    `
	row := s.db.QueryRowContext(ctx, query, userId)

	err := row.Scan(&user.Id, &user.UserName, &user.FullName, &user.BirthDate, &user.Education, &user.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `SELECT id, username, full_name, birth_date, education, city FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.UserName, &user.FullName, &user.BirthDate, &user.Education, &user.City)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
func (s *Storage) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	query := `
        SELECT id, username, full_name, birth_date, education, city
        FROM users
        WHERE username = $1
    `
	row := s.db.QueryRowContext(ctx, query, username)

	err := row.Scan(&user.Id, &user.UserName, &user.FullName, &user.BirthDate, &user.Education, &user.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
func (s *Storage) UpdateUser(ctx context.Context, userId uuid.UUID, data models.User) error {
	query := `
        UPDATE users
        SET username = $1, full_name = $2, birth_date = $3, education = $4, city = $5
        WHERE id = $6
    `
	_, err := s.db.ExecContext(ctx, query, data.UserName, data.FullName, data.BirthDate, data.Education, data.City, userId)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
