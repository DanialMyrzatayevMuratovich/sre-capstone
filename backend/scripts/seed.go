package scripts

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Глобальные переменные для хранения ID созданных объектов
var (
	userIDs     []primitive.ObjectID
	cinemaIDs   []primitive.ObjectID
	hallIDs     []primitive.ObjectID
	movieIDs    []primitive.ObjectID
	showtimeIDs []primitive.ObjectID
)

func SeedDatabase() {
	ctx := context.Background()

	log.Println("🌱 Starting database seeding...")

	// Очистить существующие данные (опционально)
	clearCollections(ctx)

	// 1. Создать пользователей
	seedUsers(ctx)

	// 2. Создать кинотеатры
	seedCinemas(ctx)

	// 3. Создать залы
	seedHalls(ctx)

	// 4. Создать фильмы
	seedMovies(ctx)

	// 5. Создать сеансы
	seedShowtimes(ctx)

	// 6. Создать тестовые брони
	seedBookings(ctx)

	// 7. Создать тестовые транзакции
	seedTransactions(ctx)

	log.Println("✅ Database seeding completed successfully!")
}

// Очистка коллекций
func clearCollections(ctx context.Context) {
	collections := []string{"users", "cinemas", "halls", "movies", "showtimes", "bookings", "transactions"}

	for _, collName := range collections {
		coll := config.GetCollection(collName)
		if err := coll.Drop(ctx); err != nil {
			log.Printf("⚠️ Warning: Could not drop collection %s: %v", collName, err)
		}
	}

	log.Println("🗑️  Collections cleared")
}

// 1. USERS
func seedUsers(ctx context.Context) {
	collection := config.GetCollection("users")

	// Hash пароль функция
	hashPassword := func(password string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hash)
	}

	users := []models.User{
		{
			Email:    "admin@cinema.kz",
			Password: hashPassword("admin123"),
			FullName: "Admin User",
			Phone:    "+77001234567",
			Role:     "admin",
			Wallet: models.Wallet{
				Balance:  0,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "manager@cinema.kz",
			Password: hashPassword("manager123"),
			FullName: "Cinema Manager",
			Phone:    "+77001234568",
			Role:     "cinema_manager",
			Wallet: models.Wallet{
				Balance:  0,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user1@example.com",
			Password: hashPassword("user123"),
			FullName: "Айдар Қасымов",
			Phone:    "+77771234567",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  5000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user2@example.com",
			Password: hashPassword("user123"),
			FullName: "Асель Нұрғалиева",
			Phone:    "+77771234568",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  10000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user3@example.com",
			Password: hashPassword("user123"),
			FullName: "Ерлан Сыдықов",
			Phone:    "+77771234569",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  3000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for i := range users {
		result, err := collection.InsertOne(ctx, users[i])
		if err != nil {
			log.Printf("❌ Error inserting user: %v", err)
			continue
		}
		userIDs = append(userIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d users", len(userIDs))
}

// 2. CINEMAS
func seedCinemas(ctx context.Context) {
	collection := config.GetCollection("cinemas")

	cinemas := []models.Cinema{
		{
			Name:    "Chaplin Cinemas Mega Alma-Ata",
			City:    "Almaty",
			Address: "ул. Розыбакиева, 247А, ТРЦ Mega",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.8512, 43.2061}, // [longitude, latitude]
			},
			Facilities:   []string{"IMAX", "4DX", "VIP", "Parking", "Food Court"},
			HallIDs:      []primitive.ObjectID{}, // заполним позже
			Rating:       4.6,
			TotalReviews: 1523,
			Images:       []string{"/uploads/cinemas/chaplin_mega.jpg"},
			CreatedAt:    time.Now(),
		},
		{
			Name:    "Kinopark Sary-Arka",
			City:    "Almaty",
			Address: "пр. Сары-Арка, 10, ТРЦ Sary-Arka",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.9286, 43.2425},
			},
			Facilities:   []string{"3D", "VIP", "Parking"},
			HallIDs:      []primitive.ObjectID{},
			Rating:       4.3,
			TotalReviews: 892,
			Images:       []string{"/uploads/cinemas/kinopark_saryarka.jpg"},
			CreatedAt:    time.Now(),
		},
		{
			Name:    "Arman Cinema Dostyk Plaza",
			City:    "Almaty",
			Address: "пр. Достык, 111, Достык Плаза",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.9539, 43.2324},
			},
			Facilities:   []string{"3D", "IMAX", "VIP", "Dolby Atmos"},
			HallIDs:      []primitive.ObjectID{},
			Rating:       4.7,
			TotalReviews: 2105,
			Images:       []string{"/uploads/cinemas/arman_dostyk.jpg"},
			CreatedAt:    time.Now(),
		},
	}

	for i := range cinemas {
		result, err := collection.InsertOne(ctx, cinemas[i])
		if err != nil {
			log.Printf("❌ Error inserting cinema: %v", err)
			continue
		}
		cinemaIDs = append(cinemaIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d cinemas", len(cinemaIDs))
}

// 3. HALLS
func seedHalls(ctx context.Context) {
	collection := config.GetCollection("halls")

	// Генератор мест для зала
	generateSeats := func(rows int, seatsPerRow int, hallType string) []models.Seat {
		seats := []models.Seat{}
		rowLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		basePrice := 2000.0
		switch hallType {
		case "VIP":
			basePrice = 4000.0
		case "IMAX":
			basePrice = 3500.0
		}

		for r := 0; r < rows; r++ {
			for s := 1; s <= seatsPerRow; s++ {
				seatType := "regular"
				price := basePrice

				// VIP места в середине
				if hallType == "VIP" || (r >= rows/3 && r <= 2*rows/3 && s >= seatsPerRow/4 && s <= 3*seatsPerRow/4) {
					seatType = "vip"
					price = basePrice * 1.5
				}

				// Couple seats в последних рядах
				if r >= rows-2 && s%2 == 1 && s < seatsPerRow {
					seatType = "couple"
					price = basePrice * 1.3
				}

				seats = append(seats, models.Seat{
					Row:    string(rowLetters[r]),
					Number: s,
					Type:   seatType,
					Price:  price,
				})
			}
		}
		return seats
	}

	// Залы для каждого кинотеатра
	halls := []models.Hall{}

	// Chaplin Mega (3 зала)
	if len(cinemaIDs) > 0 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "IMAX Hall 1",
				HallNumber: 1,
				Capacity:   240,
				Type:       "IMAX",
				Seats:      generateSeats(12, 20, "IMAX"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "Standard Hall 2",
				HallNumber: 2,
				Capacity:   150,
				Type:       "Standard",
				Seats:      generateSeats(10, 15, "Standard"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "VIP Hall 3",
				HallNumber: 3,
				Capacity:   50,
				Type:       "VIP",
				Seats:      generateSeats(5, 10, "VIP"),
			},
		)
	}

	// Kinopark Sary-Arka (2 зала)
	if len(cinemaIDs) > 1 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[1],
				Name:       "Standard Hall 1",
				HallNumber: 1,
				Capacity:   180,
				Type:       "Standard",
				Seats:      generateSeats(12, 15, "Standard"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[1],
				Name:       "3D Hall 2",
				HallNumber: 2,
				Capacity:   120,
				Type:       "Standard",
				Seats:      generateSeats(10, 12, "Standard"),
			},
		)
	}

	// Arman Dostyk (3 зала)
	if len(cinemaIDs) > 2 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "IMAX Hall 1",
				HallNumber: 1,
				Capacity:   280,
				Type:       "IMAX",
				Seats:      generateSeats(14, 20, "IMAX"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "VIP Hall 2",
				HallNumber: 2,
				Capacity:   60,
				Type:       "VIP",
				Seats:      generateSeats(6, 10, "VIP"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "Standard Hall 3",
				HallNumber: 3,
				Capacity:   200,
				Type:       "Standard",
				Seats:      generateSeats(12, 17, "Standard"),
			},
		)
	}

	// Вставить все залы
	for i := range halls {
		result, err := collection.InsertOne(ctx, halls[i])
		if err != nil {
			log.Printf("❌ Error inserting hall: %v", err)
			continue
		}
		hallID := result.InsertedID.(primitive.ObjectID)
		hallIDs = append(hallIDs, hallID)

		// Обновить cinema с hallID
		cinemaCol := config.GetCollection("cinemas")
		cinemaCol.UpdateOne(
			ctx,
			primitive.M{"_id": halls[i].CinemaID},
			primitive.M{"$push": primitive.M{"hallIds": hallID}},
		)
	}

	log.Printf("✅ Created %d halls", len(hallIDs))
}

// 4. MOVIES
func seedMovies(ctx context.Context) {
	collection := config.GetCollection("movies")

	movies := []models.Movie{
		// --- Новинки 2026 ---
		{
			Title: "Dune: Part Three", TitleKz: "Құм: 3-бөлім", TitleRu: "Дюна: Часть третья",
			Description:    "The epic conclusion to Denis Villeneuve's Dune trilogy. Paul Atreides leads the Fremen revolution against the Imperium.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/d5NXSklXo0qyIYkgV94XAgMIckC.jpg",
			Director:       "Denis Villeneuve",
			Cast:           []string{"Timothée Chalamet", "Zendaya", "Austin Butler", "Florence Pugh"},
			Genres:         []string{"Sci-Fi", "Adventure", "Drama"},
			Duration:       165, ReleaseDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.9,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian", "English"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "The Batman 2", TitleKz: "Бэтмен 2", TitleRu: "Бэтмен 2",
			Description:    "Batman continues his war on crime in Gotham City while facing the Joker and uncovering a city-wide conspiracy.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/74xTEgt7R36Fpooo50r9T25onhq.jpg",
			Director:       "Matt Reeves",
			Cast:           []string{"Robert Pattinson", "Zoë Kravitz", "Barry Keoghan", "Colin Farrell"},
			Genres:         []string{"Action", "Crime", "Thriller"},
			Duration:       155, ReleaseDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.5,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		// --- Новинки 2025 ---
		{
			Title: "Avatar: Fire and Ash", TitleKz: "Аватар: От және Күл", TitleRu: "Аватар: Огонь и Пепел",
			Description:    "Jake Sully and Neytiri face a new threat as a volcanic clan of Na'vi known as the Ash People emerges on Pandora.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/3ttduWMCgBCYjzDEC6bDqrtbne0.jpg",
			Director:       "James Cameron",
			Cast:           []string{"Sam Worthington", "Zoe Saldaña", "Sigourney Weaver", "Kate Winslet"},
			Genres:         []string{"Sci-Fi", "Adventure", "Fantasy"},
			Duration:       190, ReleaseDate: time.Date(2025, 12, 19, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.7,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian", "English"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Mission: Impossible – The Final Reckoning", TitleKz: "Мүмкін емес миссия 8", TitleRu: "Миссия невыполнима 8",
			Description:    "Ethan Hunt and his IMF team face their most impossible mission — stopping an AI threat that could end civilization.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/z53D72EAOxGRqdr7KXXWp9dJiDe.jpg",
			Director:       "Christopher McQuarrie",
			Cast:           []string{"Tom Cruise", "Hayley Atwell", "Ving Rhames", "Simon Pegg"},
			Genres:         []string{"Action", "Thriller", "Adventure"},
			Duration:       163, ReleaseDate: time.Date(2025, 5, 23, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.3,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Sinners", TitleKz: "Күнәкарлар", TitleRu: "Грешники",
			Description:    "Twin brothers trying to leave their troubled lives behind return to their hometown, only to discover that a greater evil lurks.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/lqpoeLoyCfvSbGLeSZ2mdzpyQaD.jpg",
			Director:       "Ryan Coogler",
			Cast:           []string{"Michael B. Jordan", "Hailee Steinfeld", "Jack O'Connell", "Wunmi Mosaku"},
			Genres:         []string{"Horror", "Drama", "Thriller"},
			Duration:       137, ReleaseDate: time.Date(2025, 4, 18, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 8.2,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Captain America: Brave New World", TitleKz: "Капитан Америка: Жаңа дүние", TitleRu: "Капитан Америка: Дивный новый мир",
			Description:    "Sam Wilson's first outing as Captain America leads him into a conflict involving the Red Hulk and a global conspiracy.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/pzIddUEMWhWzfvLI3TwxUG2wGoi.jpg",
			Director:       "Julius Onah",
			Cast:           []string{"Anthony Mackie", "Harrison Ford", "Danny Ramirez", "Shira Haas"},
			Genres:         []string{"Action", "Adventure", "Sci-Fi"},
			Duration:       118, ReleaseDate: time.Date(2025, 2, 14, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 6.8,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Thunderbolts*", TitleKz: "Найзағай отряды", TitleRu: "Громовержцы",
			Description:    "A team of antiheroes and reformed villains must work together when they find themselves in an impossible situation threatening the world.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/hqcexYHbiTBfDIdDWxrxPtVndBX.jpg",
			Director:       "Jake Schreier",
			Cast:           []string{"Florence Pugh", "Sebastian Stan", "David Harbour", "Wyatt Russell"},
			Genres:         []string{"Action", "Adventure", "Comedy"},
			Duration:       127, ReleaseDate: time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 7.4,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "A Minecraft Movie", TitleKz: "Майнкрафт фильмі", TitleRu: "Фильм: Майнкрафт",
			Description:    "Four misfits find themselves struggling in the Overworld, a bizarre place that thrives on imagination. To get back home, they must master the world.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/yFHHfHcUgGAxziP1C3lLt0q2T4s.jpg",
			Director:       "Jared Hess",
			Cast:           []string{"Jason Momoa", "Jack Black", "Jennifer Coolidge", "Emma Myers"},
			Genres:         []string{"Animation", "Adventure", "Comedy", "Family"},
			Duration:       101, ReleaseDate: time.Date(2025, 4, 4, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 6.5,
			Language: []string{"English", "Russian", "Kazakh"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 6, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		// --- Хиты 2024 ---
		{
			Title: "Dune: Part Two", TitleKz: "Құм: 2-бөлім", TitleRu: "Дюна: Часть вторая",
			Description:    "Paul Atreides unites with Chani and the Fremen while on a warpath of revenge against the conspirators who destroyed his family.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/cdqLnri3NEGcmfnqwk2TSIYtddg.jpg",
			Director:       "Denis Villeneuve",
			Cast:           []string{"Timothée Chalamet", "Zendaya", "Rebecca Ferguson", "Josh Brolin"},
			Genres:         []string{"Sci-Fi", "Adventure", "Drama"},
			Duration:       166, ReleaseDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.5,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian", "English"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Deadpool & Wolverine", TitleKz: "Дэдпул және Росомаха", TitleRu: "Дэдпул и Росомаха",
			Description:    "Deadpool is recruited by the Time Variance Authority and must work with the reluctant Wolverine to save the universe.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/kKy5USqxtYbsZNugPFStBHTzeIr.jpg",
			Director:       "Shawn Levy",
			Cast:           []string{"Ryan Reynolds", "Hugh Jackman", "Emma Corrin", "Matthew Macfadyen"},
			Genres:         []string{"Action", "Comedy", "Sci-Fi"},
			Duration:       128, ReleaseDate: time.Date(2024, 7, 26, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.7,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Inside Out 2", TitleKz: "Ішкі әлем 2", TitleRu: "Головоломка 2",
			Description:    "Riley enters her teenage years and new emotions — Anxiety, Envy, Ennui, and Embarrassment — arrive at Headquarters, threatening to take over.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
			Director:       "Kelsey Mann",
			Cast:           []string{"Amy Poehler", "Maya Hawke", "Phyllis Smith", "Lewis Black"},
			Genres:         []string{"Animation", "Adventure", "Comedy", "Family"},
			Duration:       96, ReleaseDate: time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 7.9,
			Language: []string{"English", "Russian", "Kazakh"}, Subtitles: []string{"Kazakh", "Russian", "English"},
			IsActive: true, AgeRestriction: 6, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Alien: Romulus", TitleKz: "Бөтен: Ромул", TitleRu: "Чужой: Ромул",
			Description:    "A group of young space colonizers come face to face with the most terrifying life form in the universe on an abandoned space station.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/b33nnKl1GSFbao4l3fZDDqsMx0F.jpg",
			Director:       "Fede Álvarez",
			Cast:           []string{"Cailee Spaeny", "David Jonsson", "Archie Renaux", "Isabela Merced"},
			Genres:         []string{"Horror", "Sci-Fi", "Thriller"},
			Duration:       119, ReleaseDate: time.Date(2024, 8, 16, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.3,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Gladiator II", TitleKz: "Гладиатор II", TitleRu: "Гладиатор II",
			Description:    "Years after the death of Maximus, Lucius — raised in North Africa — is enslaved and brought to Rome, where he must fight in the Colosseum.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/2cxhvwyEwRlysAmRH4iodkvo0z5.jpg",
			Director:       "Ridley Scott",
			Cast:           []string{"Paul Mescal", "Denzel Washington", "Pedro Pascal", "Connie Nielsen"},
			Genres:         []string{"Action", "Adventure", "Drama"},
			Duration:       148, ReleaseDate: time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.3,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 16, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Wicked", TitleKz: "Зұлым", TitleRu: "Злая",
			Description:    "The untold story of the witches of Oz — the unlikely friendship between Elphaba and Glinda before Dorothy's arrival changed everything.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/xDGbZ0JJ3mYaGKy4Nzd9Kph6M9L.jpg",
			Director:       "Jon M. Chu",
			Cast:           []string{"Cynthia Erivo", "Ariana Grande", "Michelle Yeoh", "Jeff Goldblum"},
			Genres:         []string{"Musical", "Fantasy", "Romance"},
			Duration:       160, ReleaseDate: time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 7.8,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 6, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Nosferatu", TitleKz: "Носферату", TitleRu: "Носферату",
			Description:    "A gothic tale of obsession between a haunted young woman and the terrifying ancient vampire who has followed her across centuries.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/tRYN0qW4RZjdv5bizOlCAobZSpY.jpg",
			Director:       "Robert Eggers",
			Cast:           []string{"Bill Skarsgård", "Lily-Rose Depp", "Nicholas Hoult", "Willem Dafoe"},
			Genres:         []string{"Horror", "Fantasy", "Drama"},
			Duration:       132, ReleaseDate: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.5,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "The Wild Robot", TitleKz: "Жабайы робот", TitleRu: "Дикий робот",
			Description:    "After a shipwreck, robot ROZZUM 7134 is stranded on an island and must learn to survive by adopting an orphaned gosling.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/eG9lz41mJqsI4J6ubMtVqD26q2J.jpg",
			Director:       "Chris Sanders",
			Cast:           []string{"Lupita Nyong'o", "Pedro Pascal", "Kit Connor", "Bill Nighy"},
			Genres:         []string{"Animation", "Adventure", "Drama", "Family"},
			Duration:       102, ReleaseDate: time.Date(2024, 9, 27, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 8.3,
			Language: []string{"English", "Russian", "Kazakh"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 6, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Kingdom of the Planet of the Apes", TitleKz: "Маймылдар планетасы: Патшалық", TitleRu: "Планета обезьян: Новое царство",
			Description:    "Many generations after Caesar's rule, an ape king builds an empire while a young ape embarks on a journey that challenges his understanding of the past.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/gKkl37BQuKTanygYQG1pyYgLVgf.jpg",
			Director:       "Wes Ball",
			Cast:           []string{"Owen Teague", "Freya Allan", "Kevin Durand", "William H. Macy"},
			Genres:         []string{"Action", "Adventure", "Sci-Fi"},
			Duration:       145, ReleaseDate: time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 6.9,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Furiosa: A Mad Max Saga", TitleKz: "Фуриоса: Мад Макс сагасы", TitleRu: "Фуриоса: История Безумного Макса",
			Description:    "The origin story of renegade warrior Furiosa before she teamed up with Mad Max in the post-apocalyptic world.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/iADOJ8Zymht2JPMoy3R7xceZprc.jpg",
			Director:       "George Miller",
			Cast:           []string{"Anya Taylor-Joy", "Chris Hemsworth", "Tom Burke", "Alyla Browne"},
			Genres:         []string{"Action", "Adventure", "Sci-Fi"},
			Duration:       148, ReleaseDate: time.Date(2024, 5, 24, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.8,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 16, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Beetlejuice Beetlejuice", TitleKz: "Битлджус Битлджус", TitleRu: "Битлджус Битлджус",
			Description:    "Lydia Deetz returns to Winter River, unknowingly unleashing the mischievous demon Beetlejuice once again.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/g8TbOXrNMuqq7AaKqdvqS2oG4ob.jpg",
			Director:       "Tim Burton",
			Cast:           []string{"Michael Keaton", "Winona Ryder", "Jenna Ortega", "Monica Bellucci"},
			Genres:         []string{"Comedy", "Fantasy", "Horror"},
			Duration:       104, ReleaseDate: time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 7.0,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Joker: Folie à Deux", TitleKz: "Джокер: Фолі а Де", TitleRu: "Джокер: Безумие на двоих",
			Description:    "Arthur Fleck is institutionalized at Arkham, where he encounters Harley Quinn and embarks on a musical journey.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/gSOVog7ydsaF1YpgAqBqnKYFGY.jpg",
			Director:       "Todd Phillips",
			Cast:           []string{"Joaquin Phoenix", "Lady Gaga", "Brendan Gleeson", "Catherine Keener"},
			Genres:         []string{"Drama", "Crime", "Musical"},
			Duration:       138, ReleaseDate: time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 5.5,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 16, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Moana 2", TitleKz: "Моана 2", TitleRu: "Моана 2",
			Description:    "Moana ventures into the far seas with a new crew on an unexpected voyage, following an ancient calling from the wayfinders of old.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/g6r4nsGZchNSZJcRMeMVF0sIImD.jpg",
			Director:       "David Derrick Jr.",
			Cast:           []string{"Auli'i Cravalho", "Dwayne Johnson", "Alan Tudyk", "Temuera Morrison"},
			Genres:         []string{"Animation", "Adventure", "Comedy", "Family"},
			Duration:       100, ReleaseDate: time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 7.0,
			Language: []string{"English", "Russian", "Kazakh"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 0, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "A Quiet Place: Day One", TitleKz: "Тыныш орын: Бірінші күн", TitleRu: "Тихое место: День первый",
			Description:    "Experience the day the world went quiet — the terrifying story of survival in New York City on the day of the alien invasion.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/kpKwsC9DlBifITw8M04Q0WZnDvd.jpg",
			Director:       "Michael Sarnoski",
			Cast:           []string{"Lupita Nyong'o", "Joseph Quinn", "Djimon Hounsou", "Alex Wolff"},
			Genres:         []string{"Horror", "Sci-Fi", "Thriller"},
			Duration:       99, ReleaseDate: time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 7.0,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 13, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Twisters", TitleKz: "Торнадо", TitleRu: "Смерч",
			Description:    "Storm chasers pursue deadly tornadoes across the Oklahoma plains in this high-octane thrill ride.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/pjnD08FlMAIXsfOLKQbvmO0f0MD.jpg",
			Director:       "Lee Isaac Chung",
			Cast:           []string{"Daisy Edgar-Jones", "Glen Powell", "Anthony Ramos", "Maura Tierney"},
			Genres:         []string{"Action", "Adventure", "Thriller"},
			Duration:       122, ReleaseDate: time.Date(2024, 7, 19, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 7.2,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "The Substance", TitleKz: "Субстанция", TitleRu: "Субстанция",
			Description:    "A fading celebrity uses a black-market drug to generate a younger, better version of herself — with disturbing consequences.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/lqoMzCcZYEFK729d6qzt349fB4o.jpg",
			Director:       "Coralie Fargeat",
			Cast:           []string{"Demi Moore", "Margaret Qualley", "Dennis Quaid"},
			Genres:         []string{"Horror", "Sci-Fi", "Drama"},
			Duration:       140, ReleaseDate: time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 7.5,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		// --- Классика в повторном прокате ---
		{
			Title: "Oppenheimer", TitleKz: "Оппенгеймер", TitleRu: "Оппенгеймер",
			Description:    "The story of J. Robert Oppenheimer's role in the development of the atomic bomb during World War II.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
			Director:       "Christopher Nolan",
			Cast:           []string{"Cillian Murphy", "Emily Blunt", "Matt Damon", "Robert Downey Jr."},
			Genres:         []string{"Drama", "History", "Thriller"},
			Duration:       180, ReleaseDate: time.Date(2023, 7, 21, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 8.9,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 16, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Barbie", TitleKz: "Барби", TitleRu: "Барби",
			Description:    "Barbie and Ken go on an extraordinary journey of self-discovery in the real world after Barbie starts having existential thoughts.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/iuFNMS8U5cb6xfzi51Dbkovj7vM.jpg",
			Director:       "Greta Gerwig",
			Cast:           []string{"Margot Robbie", "Ryan Gosling", "America Ferrera", "Kate McKinnon"},
			Genres:         []string{"Comedy", "Adventure", "Fantasy"},
			Duration:       114, ReleaseDate: time.Date(2023, 7, 21, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 6.9,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Poor Things", TitleKz: "Жарлы жандар", TitleRu: "Бедные-несчастные",
			Description:    "A young woman brought back to life by an eccentric scientist embarks on a journey of self-discovery across Europe.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/kCGlIMHnOm8JPXq3rXM6c5wMxcT.jpg",
			Director:       "Yorgos Lanthimos",
			Cast:           []string{"Emma Stone", "Mark Ruffalo", "Willem Dafoe", "Ramy Youssef"},
			Genres:         []string{"Comedy", "Drama", "Romance", "Sci-Fi"},
			Duration:       141, ReleaseDate: time.Date(2023, 12, 8, 0, 0, 0, 0, time.UTC),
			Rating: "R", IMDBRating: 8.0,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 18, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Wonka", TitleKz: "Вонка", TitleRu: "Вонка",
			Description:    "The story of a young Willy Wonka and how his incredible imagination and irresistible chocolate led him to the creation of his famous factory.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/3ySgD2xwasTHOK6R9bNZiEwKgYo.jpg",
			Director:       "Paul King",
			Cast:           []string{"Timothée Chalamet", "Calah Lane", "Keegan-Michael Key", "Olivia Colman"},
			Genres:         []string{"Adventure", "Comedy", "Family", "Fantasy"},
			Duration:       116, ReleaseDate: time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 7.0,
			Language: []string{"English", "Russian", "Kazakh"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 6, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Conclave", TitleKz: "Конклав", TitleRu: "Конклав",
			Description:    "A cardinal tasked with running the secret election for a new pope uncovers a shocking secret that could shake the foundation of the Catholic Church.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/oil3EZwKFp3CWxZnfGfGglesvm9.jpg",
			Director:       "Edward Berger",
			Cast:           []string{"Ralph Fiennes", "Stanley Tucci", "John Lithgow", "Isabella Rossellini"},
			Genres:         []string{"Drama", "Mystery", "Thriller"},
			Duration:       120, ReleaseDate: time.Date(2024, 10, 25, 0, 0, 0, 0, time.UTC),
			Rating: "PG", IMDBRating: 7.6,
			Language: []string{"English"}, Subtitles: []string{"Kazakh", "Russian"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
		{
			Title: "Interstellar (Re-release)", TitleKz: "Жұлдызаралық (Қайта шығарылым)", TitleRu: "Интерстеллар (Переиздание)",
			Description:    "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival. Christopher Nolan's sci-fi masterpiece in 4K IMAX.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/yQvGrMoipbRoddT0ZR8tPoR7NfX.jpg",
			Director:       "Christopher Nolan",
			Cast:           []string{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain", "Michael Caine"},
			Genres:         []string{"Sci-Fi", "Drama", "Adventure"},
			Duration:       169, ReleaseDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			Rating: "PG-13", IMDBRating: 8.7,
			Language: []string{"English", "Russian"}, Subtitles: []string{"Kazakh", "Russian", "English"},
			IsActive: true, AgeRestriction: 12, Reviews: []models.Review{}, CreatedAt: time.Now(),
		},
	}

	// Добавить тестовые отзывы
	if len(userIDs) >= 3 {
		movies[0].Reviews = []models.Review{
			{UserID: userIDs[2], Rating: 9.5, Comment: "Шынайы ғылыми-фантастикалық шедевр!", CreatedAt: time.Now().Add(-48 * time.Hour)},
			{UserID: userIDs[3], Rating: 9.0, Comment: "Лучший sci-fi фильм последних лет!", CreatedAt: time.Now().Add(-24 * time.Hour)},
		}
		movies[8].Reviews = []models.Review{ // Dune Part Two
			{UserID: userIDs[4], Rating: 9.5, Comment: "Зрелищно и эпично!", CreatedAt: time.Now().Add(-72 * time.Hour)},
		}
		movies[10].Reviews = []models.Review{ // Inside Out 2
			{UserID: userIDs[2], Rating: 8.5, Comment: "Балаларға керемет мультфильм!", CreatedAt: time.Now().Add(-36 * time.Hour)},
		}
	}

	for i := range movies {
		result, err := collection.InsertOne(ctx, movies[i])
		if err != nil {
			log.Printf("❌ Error inserting movie: %v", err)
			continue
		}
		movieIDs = append(movieIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d movies", len(movieIDs))
}

// 5. SHOWTIMES
func seedShowtimes(ctx context.Context) {
	collection := config.GetCollection("showtimes")

	now := time.Now()
	showtimes := []models.Showtime{}

	sessionTimes := []int{9, 12, 15, 18, 21}

	// Залы по кинотеатрам: [chaplin: 0,1,2], [kinopark: 3,4], [arman: 5,6,7]
	cinemaHalls := [][]int{{0, 1, 2}, {3, 4}, {5, 6, 7}}

	// Форматы и цены для разных залов
	hallFormats := map[int]struct{ format string; price float64; lang string }{
		0: {"IMAX", 3500, "Russian"},
		1: {"2D", 2000, "Russian"},
		2: {"VIP", 4000, "Kazakh"},
		3: {"3D", 2500, "Russian"},
		4: {"2D", 2000, "Kazakh"},
		5: {"IMAX", 3500, "Russian"},
		6: {"VIP", 4500, "Russian"},
		7: {"2D", 2000, "Russian"},
	}

	for day := 0; day < 14; day++ { // 14 дней вперёд
		currentDate := now.AddDate(0, 0, day)

		for movieIdx, movieID := range movieIDs {
			// Каждый фильм идёт в 1-2 кинотеатрах
			cinemasForMovie := []int{movieIdx % 3, (movieIdx + 1) % 3}

			for _, cinemaIdx := range cinemasForMovie {
				if cinemaIdx >= len(cinemaIDs) {
					continue
				}
				cinemaID := cinemaIDs[cinemaIdx]
				hallsForCinema := cinemaHalls[cinemaIdx]

				// 1-2 зала в этом кинотеатре
				maxHalls := 2
				if len(hallsForCinema) < maxHalls {
					maxHalls = len(hallsForCinema)
				}

				for hallSlot := 0; hallSlot < maxHalls; hallSlot++ {
					hallIdx := hallsForCinema[hallSlot]
					if hallIdx >= len(hallIDs) {
						continue
					}
					hallID := hallIDs[hallIdx]
					fmtInfo := hallFormats[hallSlot+(cinemaIdx*3)]

					// 2-3 сеанса в день для каждого зала
					showtimesPerDay := 2 + (movieIdx % 2)
					for t := 0; t < showtimesPerDay; t++ {
						hour := sessionTimes[t+(hallSlot*2)%len(sessionTimes)]
						startTime := time.Date(
							currentDate.Year(), currentDate.Month(), currentDate.Day(),
							hour, 0, 0, 0, time.Local,
						)
						if startTime.Before(now) {
							continue
						}

						duration := 100 + (movieIdx % 10 * 10)
						endTime := startTime.Add(time.Duration(duration) * time.Minute)

						capacity := 150
						booked := movieIdx * 3 % 50

						showtime := models.Showtime{
							MovieID:        movieID,
							CinemaID:       cinemaID,
							HallID:         hallID,
							StartTime:      startTime,
							EndTime:        endTime,
							BasePrice:      fmtInfo.price,
							Format:         fmtInfo.format,
							Language:       fmtInfo.lang,
							Subtitles:      "Kazakh",
							AvailableSeats: capacity - booked,
							BookedSeats:    []models.BookedSeat{},
							CreatedAt:      time.Now(),
						}
						showtimes = append(showtimes, showtime)
					}
				}
			}
		}
	}

	for i := range showtimes {
		result, err := collection.InsertOne(ctx, showtimes[i])
		if err != nil {
			log.Printf("❌ Error inserting showtime: %v", err)
			continue
		}
		showtimeIDs = append(showtimeIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d showtimes", len(showtimeIDs))
}

// 6. BOOKINGS
func seedBookings(ctx context.Context) {
	collection := config.GetCollection("bookings")

	// Создать несколько тестовых броней
	if len(userIDs) < 3 || len(showtimeIDs) < 2 {
		log.Println("⚠️ Not enough data to create bookings")
		return
	}

	bookings := []models.Booking{
		{
			BookingNumber: "BK-20260201-001234",
			UserID:        userIDs[2],
			ShowtimeID:    showtimeIDs[0],
			Seats: []models.BookingSeat{
				{Row: "E", Number: 10, Price: 2000},
				{Row: "E", Number: 11, Price: 2000},
			},
			TotalAmount: 4000,
			Status:      "confirmed",
			Payment: models.Payment{
				Method:        "wallet",
				TransactionID: "TXN-" + time.Now().Format("20060102150405"),
				PaidAt:        time.Now().Add(-2 * time.Hour),
				Status:        "completed",
			},
			QRCode:    "QR-BK-20260201-001234",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			BookingNumber: "BK-20260201-001235",
			UserID:        userIDs[3],
			ShowtimeID:    showtimeIDs[1],
			Seats: []models.BookingSeat{
				{Row: "D", Number: 5, Price: 3000},
				{Row: "D", Number: 6, Price: 3000},
				{Row: "D", Number: 7, Price: 3000},
			},
			TotalAmount: 9000,
			Status:      "pending",
			Payment: models.Payment{
				Method: "card",
				Status: "pending",
			},
			QRCode:    "QR-BK-20260201-001235",
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now().Add(-5 * time.Minute),
			UpdatedAt: time.Now().Add(-5 * time.Minute),
		},
	}

	for i := range bookings {
		result, err := collection.InsertOne(ctx, bookings[i])
		if err != nil {
			log.Printf("❌ Error inserting booking: %v", err)
			continue
		}
		_ = result.InsertedID.(primitive.ObjectID)
	}

	log.Printf("✅ Created %d bookings", len(bookings))
}

// 7. TRANSACTIONS
func seedTransactions(ctx context.Context) {
	collection := config.GetCollection("transactions")

	if len(userIDs) < 3 {
		log.Println("⚠️ Not enough users to create transactions")
		return
	}

	transactions := []models.Transaction{
		{
			UserID:      userIDs[2],
			Type:        "wallet_topup",
			Amount:      5000,
			Status:      "completed",
			Description: "Пополнение кошелька через карту",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			UserID:      userIDs[2],
			Type:        "booking",
			Amount:      -4000,
			Status:      "completed",
			Description: "Оплата билетов на 'Dune: Part Three'",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:      userIDs[3],
			Type:        "wallet_topup",
			Amount:      10000,
			Status:      "completed",
			Description: "Пополнение кошелька",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
	}

	for i := range transactions {
		_, err := collection.InsertOne(ctx, transactions[i])
		if err != nil {
			log.Printf("❌ Error inserting transaction: %v", err)
			continue
		}
	}

	log.Printf("✅ Created %d transactions", len(transactions))
}
