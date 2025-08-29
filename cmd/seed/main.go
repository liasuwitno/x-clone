package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/v5/stdlib" // driver "pgx"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func hashPassword(pw string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b)
}

func main() {
	// load .env
	_ = godotenv.Load()

	dsn := mustEnv("DATABASE_URL")

	// open db with pgx driver
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("open db:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// seed users
	users := []struct {
		username, email, password string
	}{
		{"lia", "lia@example.com", "secret123"},
		{"joko", "joko@example.com", "secret123"},
		{"nina", "nina@example.com", "secret123"},
	}

	for _, u := range users {
		_, err := db.ExecContext(ctx, `
			INSERT INTO users (username, email, password, bio)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (username) DO NOTHING
		`, u.username, u.email, hashPassword(u.password), "Hello, I'm "+u.username)
		if err != nil {
			log.Fatal("insert user:", err)
		}
	}

	// ambil user IDs
	var liaID, jokoID, ninaID int64
	_ = db.QueryRowContext(ctx, `SELECT id FROM users WHERE username='lia'`).Scan(&liaID)
	_ = db.QueryRowContext(ctx, `SELECT id FROM users WHERE username='joko'`).Scan(&jokoID)
	_ = db.QueryRowContext(ctx, `SELECT id FROM users WHERE username='nina'`).Scan(&ninaID)

	// follows
	_, _ = db.ExecContext(ctx, `INSERT INTO follows (following_user_id, followed_user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, liaID, jokoID)
	_, _ = db.ExecContext(ctx, `INSERT INTO follows (following_user_id, followed_user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, liaID, ninaID)
	_, _ = db.ExecContext(ctx, `INSERT INTO follows (following_user_id, followed_user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, jokoID, liaID)

	// tweets
	var t1, t2, t3 int64
	_ = db.QueryRowContext(ctx, `
		INSERT INTO tweets (title, body, user_id, status)
		VALUES ('First Post','Halo dunia!', $1, 'published')
		RETURNING id
	`, liaID).Scan(&t1)

	_ = db.QueryRowContext(ctx, `
		INSERT INTO tweets (title, body, user_id, status)
		VALUES ('Go Tips','Gunakan context untuk cancelable ops.', $1, 'published')
		RETURNING id
	`, jokoID).Scan(&t2)

	_ = db.QueryRowContext(ctx, `
		INSERT INTO tweets (title, body, user_id, status)
		VALUES ('React Hooks','SWR untuk data fetching elegan.', $1, 'published')
		RETURNING id
	`, ninaID).Scan(&t3)

	// likes
	_, _ = db.ExecContext(ctx, `INSERT INTO likes (user_id, tweet_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, liaID, t2)
	_, _ = db.ExecContext(ctx, `INSERT INTO likes (user_id, tweet_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, jokoID, t1)
	_, _ = db.ExecContext(ctx, `INSERT INTO likes (user_id, tweet_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, ninaID, t1)

	// edit history
	_, _ = db.ExecContext(ctx, `INSERT INTO edit_history (tweet_id, previous_body) VALUES ($1,$2)`, t1, "Halo~")

	log.Println("✅ Seed completed.")
}
