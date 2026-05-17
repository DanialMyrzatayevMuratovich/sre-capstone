<template>
  <router-link :to="`/movies/${movie._id || movie.id}`" class="movie-card">
    <div class="poster-wrap">
      <img
        v-if="movie.posterUrl && !imgError"
        :src="movie.posterUrl"
        :alt="movie.title"
        class="poster-img"
        @error="imgError = true"
      />
      <div v-else class="poster-placeholder">
        <span class="placeholder-icon">🎬</span>
        <span class="placeholder-title">{{ movie.title }}</span>
      </div>

      <div class="age-badge" v-if="movie.ageRestriction > 0">{{ movie.ageRestriction }}+</div>
      <div class="rating-badge">⭐ {{ movie.imdbRating?.toFixed(1) }}</div>

      <div class="hover-overlay">
        <div class="overlay-genres">
          <span v-for="g in (movie.genres || []).slice(0, 2)" :key="g" class="overlay-genre">{{ g }}</span>
        </div>
        <p class="overlay-desc">{{ movie.description?.slice(0, 100) }}...</p>
        <div class="overlay-footer">
          <span class="overlay-duration">🕐 {{ movie.duration }} мин</span>
          <span class="overlay-cta">Смотреть →</span>
        </div>
      </div>
    </div>

    <div class="card-info">
      <h3 class="card-title">{{ movie.titleRu || movie.title }}</h3>
      <div class="card-meta">
        <span class="card-year">{{ new Date(movie.releaseDate).getFullYear() }}</span>
        <span class="card-rating-text">IMDb {{ movie.imdbRating?.toFixed(1) }}</span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  movie: { type: Object, required: true }
})

const imgError = ref(false)
</script>

<style scoped>
.movie-card {
  display: block;
  text-decoration: none;
  color: var(--text);
  cursor: pointer;
  transition: transform 0.25s ease;
}

.movie-card:hover { transform: translateY(-6px); }

.poster-wrap {
  position: relative;
  aspect-ratio: 2/3;
  border-radius: 12px;
  overflow: hidden;
  background: var(--dark-light);
  margin-bottom: 10px;
}

.poster-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.movie-card:hover .poster-img { transform: scale(1.05); }

.poster-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: linear-gradient(135deg, var(--dark-light), var(--dark-lighter));
  padding: 20px;
  text-align: center;
}

.placeholder-icon { font-size: 48px; opacity: 0.4; }
.placeholder-title { font-size: 13px; font-weight: 600; color: var(--text-gray); line-height: 1.3; }

.age-badge {
  position: absolute;
  top: 10px;
  left: 10px;
  background: rgba(0,0,0,0.7);
  border: 1px solid rgba(255,255,255,0.25);
  color: white;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
}

.rating-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(0,0,0,0.75);
  backdrop-filter: blur(4px);
  color: #fbbf24;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 700;
}

.hover-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.95) 0%, rgba(0,0,0,0.55) 55%, transparent 100%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 14px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.movie-card:hover .hover-overlay { opacity: 1; }

.overlay-genres { display: flex; gap: 6px; margin-bottom: 8px; flex-wrap: wrap; }
.overlay-genre {
  background: rgba(229,9,20,0.85);
  color: white;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.overlay-desc {
  font-size: 11px;
  color: rgba(255,255,255,0.8);
  line-height: 1.5;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.overlay-footer { display: flex; justify-content: space-between; align-items: center; }
.overlay-duration { font-size: 12px; color: rgba(255,255,255,0.65); }
.overlay-cta { font-size: 13px; font-weight: 700; color: var(--primary); }

.card-info { padding: 0 2px; }
.card-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-meta { display: flex; justify-content: space-between; font-size: 12px; color: var(--text-gray); }
.card-rating-text { color: #fbbf24; font-weight: 600; }
</style>
