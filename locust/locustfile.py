"""
CinemaHub Load Test — Locust
Run: locust -f locustfile.py --host http://localhost:8080
Web UI: http://localhost:8089
"""

from locust import HttpUser, task, between
import random

MOVIE_IDS = []  # populated from /api/movies on first request


class CinemaHubUser(HttpUser):
    """Simulates a mix of anonymous browsing and authenticated booking actions."""

    wait_time = between(1, 3)
    token = None

    def on_start(self):
        self._fetch_movie_ids()
        self._login()

    # ------------------------------------------------------------------
    # Setup helpers
    # ------------------------------------------------------------------

    def _fetch_movie_ids(self):
        global MOVIE_IDS
        if MOVIE_IDS:
            return
        with self.client.get("/api/movies", name="[setup] GET /api/movies",
                             catch_response=True) as resp:
            if resp.status_code == 200:
                data = resp.json()
                movies = data if isinstance(data, list) else data.get("movies", data.get("data", []))
                MOVIE_IDS = [m.get("_id") or m.get("id") for m in movies
                             if m.get("_id") or m.get("id")]
                resp.success()

    def _login(self):
        credentials = {"email": "loadtest@cinema.com", "password": "LoadTest123!"}

        with self.client.post("/api/auth/login", json=credentials,
                              name="[setup] POST /api/auth/login",
                              catch_response=True) as resp:
            if resp.status_code == 200:
                self.token = resp.json().get("token")
                resp.success()
                return
            resp.success()  # ignore 401 — will register below

        # Register a new test account if login failed
        reg_data = {"name": "Load Tester", "email": "loadtest@cinema.com",
                    "password": "LoadTest123!"}
        with self.client.post("/api/auth/register", json=reg_data,
                              name="[setup] POST /api/auth/register",
                              catch_response=True) as r:
            r.success()  # ignore duplicate-email errors

        with self.client.post("/api/auth/login", json=credentials,
                              name="[setup] POST /api/auth/login",
                              catch_response=True) as resp2:
            if resp2.status_code == 200:
                self.token = resp2.json().get("token")
            resp2.success()

    def _auth(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    # ------------------------------------------------------------------
    # Tasks — weight = relative frequency
    # ------------------------------------------------------------------

    @task(10)
    def browse_movies(self):
        self.client.get("/api/movies", name="GET /api/movies")

    @task(7)
    def view_movie_detail(self):
        if MOVIE_IDS:
            mid = random.choice(MOVIE_IDS)
            self.client.get(f"/api/movies/{mid}", name="GET /api/movies/:id")

    @task(5)
    def browse_cinemas(self):
        self.client.get("/api/cinemas", name="GET /api/cinemas")

    @task(5)
    def get_showtimes(self):
        self.client.get("/api/showtimes", name="GET /api/showtimes")

    @task(3)
    def health_check(self):
        self.client.get("/api/health", name="GET /api/health")

    @task(2)
    def view_profile(self):
        if self.token:
            self.client.get("/api/profile", headers=self._auth(),
                            name="GET /api/profile")

    @task(1)
    def view_my_bookings(self):
        if self.token:
            self.client.get("/api/bookings/my", headers=self._auth(),
                            name="GET /api/bookings/my")
