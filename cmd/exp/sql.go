package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func (cfg *PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode)
}

func main() {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DbName:   "lenslocked",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	defer db.Close()

	//============================================================================

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT UNIQUE NOT NULL
	//   );
	//   	CREATE TABLE IF NOT EXISTS orders (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	amount INT,
	// 	description TEXT
	//   );
	// `)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Tables created.")

	// //Inserting data
	// row := db.QueryRow(`
	// 	INSERT INTO users(name, email)
	// 	VALUES($1, $2) RETURNING id;
	// `, "Leo Messi", "Leoo @gamil.com")
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User: ", id, " created")

	// //Querying Single Row
	// id := 5
	// row := db.QueryRow(`
	// 	SELECT name, email
	// 	FROM users WHERE id = $1
	// `, id)
	// var name, email string
	// err = row.Scan(&name, &email)
	// if err == sql.ErrNoRows {
	// 	fmt.Println(sql.ErrNoRows)
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User %d, name: %s, email: %s\n", id, name, email)

	//============================================================================

	// userID := 1 // Pick an ID that exists in your DB
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err := db.Exec(`
	// 	INSERT INTO orders(user_id, amount, description)
	// 	VALUES($1, $2, $3)`, userID, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("Created fake orders.")

	//============================================================================

	// type Order struct {
	// 	Id, UserId, Amount int
	// 	Description        string
	// }

	// var orders []Order
	// userId := 1
	// rows, err := db.Query(`
	// 	SELECT id, amount, description
	// 	FROM orders WHERE user_id = $1
	// `, userId)
	// defer rows.Close()

	// for rows.Next() {
	// 	var order Order
	// 	order.UserId = userId
	// 	err = rows.Scan(&order.Id, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(orders)

	//============================================================================

	// us := models.UserService{DB: db}
	// u, err := us.Create("mm", "sdsd")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(u)

}
