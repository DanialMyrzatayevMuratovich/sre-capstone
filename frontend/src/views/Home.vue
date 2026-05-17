<template>
  <div class="home">
    <!-- Hero Banner -->
    <div class="hero-banner" v-if="featuredMovie">
      <div
        class="hero-bg"
        :style="featuredMovie.posterUrl ? `background-image: url(${featuredMovie.posterUrl})` : ''"
      ></div>
      <div class="hero-overlay"></div>
      <div class="hero-content container">
        <div class="hero-badge">⭐ Рекомендуем</div>
        <h1 class="hero-title">{{ featuredMovie.titleRu || featuredMovie.title }}</h1>
        <div class="hero-meta">
          <span class="hero-rating">⭐ {{ featuredMovie.imdbRating }}</span>
          <span class="hero-dot">•</span>
          <span>{{ featuredMovie.duration }} мин</span>
          <span class="hero-dot">•</span>
          <span v-for="g in (featuredMovie.genres || []).slice(0, 2)" :key="g" class="hero-genre">{{ g }}</span>
        </div>
        <p class="hero-desc">{{ featuredMovie.description?.slice(0, 160) }}...</p>
        <router-link :to="`/movies/${featuredMovie._id || featuredMovie.id}`" class="btn btn-primary hero-btn">
          🎬 Смотреть детали
        </router-link>
      </div>
    </div>

    <div class="container">
      <!-- Search -->
      <div class="search-wrap">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск фильмов..."
            class="search-input"
            @input="searchMovies"
          />
          <button v-if="searchQuery" class="search-clear" @click="clearSearch">✕</button>
        </div>
      </div>

      <!-- Genre Filters -->
      <div class="genre-filters">
        <button
          v-for="g in genres"
          :key="g.value"
          class="genre-pill"
          :class="{ active: selectedGenre === g.value }"
          @click="selectGenre(g.value)"
        >
          {{ g.icon }} {{ g.label }}
        </button>
      </div>

      <!-- Sort -->
      <div class="sort-row">
        <span class="results-count" v-if="!loading">{{ totalMovies }} фильмов</span>
        <select v-model="sortBy" class="sort-select" @change="fetchMovies">
          <option value="imdbRating_desc">⭐ По рейтингу</option>
          <option value="releaseDate_desc">🆕 Новинки</option>
          <option value="title_asc">🔤 По названию</option>
        </select>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading-grid">
        <div v-for="i in 12" :key="i" class="skeleton-card"></div>
      </div>

      <!-- Movies Grid -->
      <div v-else-if="movies.length > 0" class="movies-grid">
        <MovieCard
          v-for="movie in movies"
          :key="movie._id || movie.id"
          :movie="movie"
        />
      </div>

      <!-- No Results -->
      <div v-else class="no-results">
        <div class="no-results-icon">🎬</div>
        <h3>Фильмы не найдены</h3>
        <p>Попробуйте изменить критерии поиска</p>
        <button class="btn btn-primary" @click="clearSearch">Сбросить фильтры</button>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button class="page-arrow" :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">‹</button>
        <button
          v-for="page in visiblePages"
          :key="page"
          class="page-btn"
          :class="{ active: page === currentPage }"
          @click="goToPage(page)"
        >{{ page }}</button>
        <button class="page-arrow" :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">›</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import MovieCard from '../components/MovieCard.vue'
import api from '../services/api'

const movies = ref([])
const loading = ref(true)
const searchQuery = ref('')
const selectedGenre = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const totalMovies = ref(0)
const sortBy = ref('imdbRating_desc')
const featuredMovie = ref(null)
const limit = 12

const genres = [
  { value: '', label: 'Все', icon: '🎬' },
  { value: 'Action', label: 'Экшн', icon: '💥' },
  { value: 'Sci-Fi', label: 'Фантастика', icon: '🚀' },
  { value: 'Drama', label: 'Драма', icon: '🎭' },
  { value: 'Horror', label: 'Ужасы', icon: '👻' },
  { value: 'Animation', label: 'Анимация', icon: '🎨' },
  { value: 'Comedy', label: 'Комедия', icon: '😂' },
  { value: 'Thriller', label: 'Триллер', icon: '🔪' },
  { value: 'Adventure', label: 'Приключения', icon: '🗺️' },
  { value: 'Family', label: 'Семейный', icon: '👨‍👩‍👧' },
]

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - 2)
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  if (end - start < maxVisible - 1) start = Math.max(1, end - maxVisible + 1)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const fetchMovies = async () => {
  loading.value = true
  try {
    const [sortField, sortOrder] = sortBy.value.split('_')
    const params = {
      page: currentPage.value,
      limit,
      isActive: 'true',
      sortBy: sortField,
      sortOrder,
    }
    if (searchQuery.value) params.search = searchQuery.value
    if (selectedGenre.value) params.genre = selectedGenre.value

    const response = await api.getMovies(params)
    movies.value = response.data.data
    totalPages.value = response.data.pagination?.totalPages || 1
    totalMovies.value = response.data.pagination?.total || movies.value.length

    if (!featuredMovie.value && movies.value.length > 0) {
      featuredMovie.value = movies.value.find(m => m.imdbRating >= 8) || movies.value[0]
    }
  } catch (error) {
    console.error('Failed to fetch movies:', error)
    movies.value = []
  } finally {
    loading.value = false
  }
}

const searchMovies = () => {
  currentPage.value = 1
  fetchMovies()
}

const selectGenre = (genre) => {
  selectedGenre.value = genre
  currentPage.value = 1
  fetchMovies()
}

const clearSearch = () => {
  searchQuery.value = ''
  selectedGenre.value = ''
  currentPage.value = 1
  fetchMovies()
}

const goToPage = (page) => {
  currentPage.value = page
  fetchMovies()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => fetchMovies())
</script>

<style scoped>
.home {
  min-height: calc(100vh - 70px);
  padding-bottom: 60px;
}

/* ── Hero Banner ── */
.hero-banner {
  position: relative;
  height: 500px;
  overflow: hidden;
  margin-bottom: 48px;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center top;
  filter: blur(2px) brightness(0.4);
  transform: scale(1.05);
  transition: transform 8s ease;
}

.hero-banner:hover .hero-bg { transform: scale(1.08); }

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to right, rgba(0,0,0,0.85) 40%, transparent 100%),
              linear-gradient(to top, rgba(0,0,0,0.6) 0%, transparent 60%);
}

.hero-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: 600px;
  padding-top: 40px;
}

.hero-badge {
  display: inline-block;
  background: var(--primary);
  color: white;
  padding: 4px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 1px;
  text-transform: uppercase;
  margin-bottom: 16px;
  width: fit-content;
}

.hero-title {
  font-size: 42px;
  font-weight: 800;
  line-height: 1.15;
  margin-bottom: 12px;
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
}

.hero-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  font-size: 14px;
  color: rgba(255,255,255,0.8);
}

.hero-rating { color: #fbbf24; font-weight: 700; }
.hero-dot { opacity: 0.4; }
.hero-genre {
  background: rgba(255,255,255,0.15);
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
}

.hero-desc {
  font-size: 15px;
  color: rgba(255,255,255,0.75);
  line-height: 1.6;
  margin-bottom: 24px;
}

.hero-btn {
  width: fit-content;
  padding: 12px 28px;
  font-size: 16px;
  font-weight: 700;
}

/* ── Search ── */
.search-wrap {
  margin-bottom: 24px;
}

.search-box {
  position: relative;
  max-width: 560px;
  margin: 0 auto;
}

.search-icon {
  position: absolute;
  left: 18px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 18px;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 14px 48px;
  background: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  border-radius: 30px;
  color: var(--text);
  font-size: 16px;
  transition: all 0.25s;
}

.search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(229,9,20,0.12);
  outline: none;
}

.search-clear {
  position: absolute;
  right: 18px;
  top: 50%;
  transform: translateY(-50%);
  background: var(--dark-lighter);
  border: none;
  color: var(--text-gray);
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ── Genre Pills ── */
.genre-filters {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.genre-pill {
  padding: 8px 18px;
  background: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  border-radius: 20px;
  color: var(--text-gray);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.genre-pill:hover { border-color: var(--primary); color: var(--text); }
.genre-pill.active {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

/* ── Sort row ── */
.sort-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.results-count {
  font-size: 14px;
  color: var(--text-gray);
}

.sort-select {
  background: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  border-radius: 10px;
  color: var(--text);
  padding: 8px 14px;
  font-size: 14px;
  cursor: pointer;
}

/* ── Skeleton loader ── */
.loading-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.skeleton-card {
  aspect-ratio: 2/3;
  background: linear-gradient(90deg, var(--dark-light) 25%, var(--dark-lighter) 50%, var(--dark-light) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 12px;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ── Movies Grid ── */
.movies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

/* ── No Results ── */
.no-results {
  text-align: center;
  padding: 80px 20px;
}

.no-results-icon { font-size: 72px; margin-bottom: 20px; opacity: 0.5; }
.no-results h3 { font-size: 24px; margin-bottom: 10px; }
.no-results p { color: var(--text-gray); margin-bottom: 24px; }

/* ── Pagination ── */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 40px;
}

.page-arrow {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  color: var(--text);
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.page-arrow:disabled { opacity: 0.3; cursor: not-allowed; }
.page-arrow:not(:disabled):hover { border-color: var(--primary); }

.page-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--dark-light);
  border: 2px solid var(--dark-lighter);
  color: var(--text);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover { border-color: var(--primary); }
.page-btn.active { background: var(--primary); border-color: var(--primary); color: white; }

@media (max-width: 768px) {
  .hero-banner { height: 380px; }
  .hero-title { font-size: 28px; }
  .hero-desc { display: none; }
  .movies-grid { grid-template-columns: repeat(2, 1fr); gap: 12px; }
  .loading-grid { grid-template-columns: repeat(2, 1fr); gap: 12px; }
}
</style>
