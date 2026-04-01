package database

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
}

func New() *Database {
	return &Database{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
	}
}
