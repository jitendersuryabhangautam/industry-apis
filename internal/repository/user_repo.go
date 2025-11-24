package repository

import (
	"context"
	"fmt"
	"industry-api/internal/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {

	query := `
		INSERT INTO users ( name, email, password_hash,Phone, role)
		VALUES ($1,$2,$3,$4,$5)
		Returning id, created_at
	`
	return r.db.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.Phone, user.Role).Scan(&user.ID, &user.CreatedAt)

}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, name, email, password_hash, phone, role, created_at
		FROM users
		WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *UserRepository) GetUserList(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error) {
	// Build dynamic query
	baseQuery := `
		SELECT id, name, email, phone, role, is_active, created_at, updated_at 
		FROM users 
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"

	var args []interface{}
	var conditions []string
	argPos := 1

	// Add role condition if provided
	if role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argPos))
		args = append(args, role)
		argPos++
	}

	// Add is_active condition if provided
	if isActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *isActive)
		argPos++
	}

	// Add search condition if provided (search in name or email)
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(email) LIKE $%d)", argPos, argPos))
		args = append(args, searchTerm)
		argPos++
	}

	// Add WHERE conditions if any
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Add ordering
	baseQuery += " ORDER BY created_at DESC"

	// Execute count query to get total records
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	var users []models.User

	// Handle pagination: if limit is 0, return all records without pagination
	if limit == 0 {
		// Return all results (no pagination)
		rows, err := r.db.Query(ctx, baseQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to query users: %w", err)
		}
		defer rows.Close()

		users, err = r.scanUsers(rows)
		if err != nil {
			return nil, err
		}
	} else {
		// Apply pagination with LIMIT and OFFSET
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)

		// Calculate offset
		offset := (page - 1) * limit
		if offset < 0 {
			offset = 0
		}

		paginationArgs := append(args, limit, offset)

		rows, err := r.db.Query(ctx, baseQuery, paginationArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to query users: %w", err)
		}
		defer rows.Close()

		users, err = r.scanUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserListResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: r.calculateTotalPages(total, limit),
		Users:      users,
	}, nil
}

// Helper function to calculate total pages
func (r *UserRepository) calculateTotalPages(total, limit int) int {
	if limit == 0 || total == 0 {
		return 1
	}
	pages := total / limit
	if total%limit > 0 {
		pages++
	}
	return pages
}

// Enhanced helper function to scan users from rows
func (r *UserRepository) scanUsers(rows pgx.Rows) ([]models.User, error) {
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}
