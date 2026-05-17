<template>
  <div class="profile">
    <div class="container">
      <h1 class="page-title">Мой профиль</h1>

      <!-- User Info Card -->
      <div class="profile-card">
        <div class="profile-header">
          <div class="user-avatar">
            <span class="avatar-icon">👤</span>
          </div>
          <div class="user-info">
            <h2 class="user-name">{{ authStore.user?.fullName }}</h2>
            <p class="user-email">{{ authStore.user?.email }}</p>
            <p class="user-role">{{ getRoleText(authStore.user?.role) }}</p>
          </div>
        </div>

        <div class="wallet-info">
          <div class="wallet-card">
            <div class="wallet-icon">💳</div>
            <div class="wallet-details">
              <div class="wallet-label">Баланс кошелька</div>
              <div class="wallet-balance">
                {{ formatCurrency(authStore.user?.wallet?.balance || 0) }}
              </div>
            </div>
            <router-link to="/topup" class="btn btn-primary topup-btn">
              + Пополнить
            </router-link>
          </div>
        </div>
      </div>

      <!-- Success Alert -->
      <div v-if="showSuccessMessage" class="alert alert-success">
        ✅ Бронирование успешно создано!
      </div>

      <!-- Bookings Section -->
      <div class="bookings-section">
        <div class="section-header">
          <h2 class="section-title">Мои бронирования</h2>
          
          <!-- Filter Tabs -->
          <div class="filter-tabs">
            <button
              v-for="status in statusFilters"
              :key="status.value"
              class="tab-btn"
              :class="{ active: filterStatus === status.value }"
              @click="filterStatus = status.value; fetchBookings()"
            >
              {{ status.label }}
            </button>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loadingBookings" class="loading">
          <div class="spinner"></div>
        </div>

        <!-- Bookings List -->
        <div v-else-if="bookings.length > 0" class="bookings-list">
          <BookingCard
            v-for="booking in bookings"
            :key="booking._id || booking.id"
            :booking="booking"
            :movie-title="getMovieTitle(booking)"
            @confirm="confirmBooking"
            @cancel="cancelBooking"
          />
        </div>

        <!-- No Bookings -->
        <div v-else class="no-bookings">
          <span class="no-bookings-icon">🎫</span>
          <h3>Нет бронирований</h3>
          <p>Вы еще не забронировали ни одного билета</p>
          <router-link to="/" class="btn btn-primary">
            Посмотреть фильмы
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../store/auth'
import BookingCard from '../components/BookingCard.vue'
import api from '../services/api'
import { formatCurrency } from '../utils/formatters'

const route = useRoute()
const authStore = useAuthStore()

const bookings = ref([])
const loadingBookings = ref(true)
const filterStatus = ref('')
const showSuccessMessage = ref(false)

const statusFilters = [
  { label: 'Все', value: '' },
  { label: 'Подтверждено', value: 'confirmed' },
  { label: 'Ожидает оплаты', value: 'pending' },
  { label: 'Отменено', value: 'cancelled' }
]

const fetchBookings = async () => {
  loadingBookings.value = true
  try {
    const params = {}
    if (filterStatus.value) {
      params.status = filterStatus.value
    }

    const response = await api.getMyBookings(params)
    bookings.value = response.data.data
  } catch (error) {
    console.error('Failed to fetch bookings:', error)
    bookings.value = []
  } finally {
    loadingBookings.value = false
  }
}

const confirmBooking = async (bookingId) => {
  try {
    await api.confirmBooking(bookingId)
    
    // Обновить список броней
    await fetchBookings()
    
    // Обновить профиль
    await authStore.fetchProfile()
    
    alert('✅ Бронирование подтверждено!')
  } catch (error) {
    alert(error.response?.data?.error || 'Не удалось подтвердить бронирование')
  }
}

const cancelBooking = async (bookingId) => {
  if (!confirm('Вы уверены, что хотите отменить это бронирование?')) {
    return
  }

  try {
    await api.cancelBooking(bookingId)
    
    // Обновить список броней
    await fetchBookings()
    
    // Обновить профиль (баланс вернулся)
    await authStore.fetchProfile()
    
    alert('✅ Бронирование отменено. Средства возвращены на кошелек.')
  } catch (error) {
    alert(error.response?.data?.error || 'Не удалось отменить бронирование')
  }
}

const getRoleText = (role) => {
  const roles = {
    user: 'Пользователь',
    admin: 'Администратор',
    cinema_manager: 'Менеджер кинотеатра'
  }
  return roles[role] || role
}

const getMovieTitle = (booking) => {
  return booking.movieTitle || 'Фильм'
}

onMounted(async () => {
  // Показать сообщение об успешном создании брони
  if (route.query.bookingCreated) {
    showSuccessMessage.value = true
    setTimeout(() => {
      showSuccessMessage.value = false
    }, 5000)
  }

  // Обновить профиль
  await authStore.fetchProfile()
  
  // Загрузить брони
  await fetchBookings()
})
</script>

<style scoped>
.profile {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.page-title {
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 30px;
}

.profile-card {
  background: linear-gradient(135deg, var(--dark-light), var(--dark-lighter));
  padding: 30px;
  border-radius: 16px;
  margin-bottom: 30px;
  border: 2px solid var(--dark-lighter);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 4px;
}

.user-email {
  font-size: 16px;
  color: var(--text-gray);
  margin-bottom: 4px;
}

.user-role {
  font-size: 14px;
  color: var(--secondary);
  font-weight: 600;
}

.wallet-info {
  margin-top: 24px;
}

.wallet-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background-color: var(--dark);
  border-radius: 12px;
  border: 2px solid var(--primary);
}

.wallet-icon {
  font-size: 48px;
}

.wallet-label {
  font-size: 14px;
  color: var(--text-gray);
  margin-bottom: 6px;
}

.wallet-balance {
  font-size: 32px;
  font-weight: bold;
  color: var(--primary);
}

.topup-btn {
  margin-left: auto;
  white-space: nowrap;
  flex-shrink: 0;
}

.bookings-section {
  margin-top: 40px;
}

.section-header {
  margin-bottom: 24px;
}

.section-title {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 16px;
}

.filter-tabs {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 10px 20px;
  background-color: var(--dark-light);
  color: var(--text);
  border: 2px solid var(--dark-lighter);
  border-radius: 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.tab-btn:hover {
  border-color: var(--primary);
}

.tab-btn.active {
  background-color: var(--primary);
  border-color: var(--primary);
  color: white;
}

.bookings-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-bookings {
  text-align: center;
  padding: 80px 20px;
  background-color: var(--dark-light);
  border-radius: 16px;
}

.no-bookings-icon {
  font-size: 80px;
  display: block;
  margin-bottom: 20px;
  opacity: 0.5;
}

.no-bookings h3 {
  font-size: 24px;
  margin-bottom: 12px;
}

.no-bookings p {
  color: var(--text-gray);
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .wallet-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>