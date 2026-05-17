<template>
  <div class="seat-selection">
    <div class="container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
      </div>

      <div v-else-if="showtime && hall">
        <!-- Header -->
        <div class="booking-header">
          <button class="btn btn-secondary" @click="goBack">
            ← Назад
          </button>
          <h1 class="page-title">Выбор мест</h1>
        </div>

        <!-- Movie Info -->
        <div class="movie-info-card">
          <div class="movie-info-content">
            <h2 class="movie-title">{{ movieTitle }}</h2>
            <div class="session-info">
              <span class="info-item">📅 {{ formatDate(showtime.startTime) }}</span>
              <span class="info-item">🏢 {{ cinemaName }}</span>
              <span class="info-item">🎬 {{ showtime.format }}</span>
              <span class="info-item">🗣️ {{ showtime.language }}</span>
            </div>
          </div>
        </div>

        <!-- Seat Map -->
        <div class="seat-selection-area">
          <SeatMap
            :seats="hall.seats"
            :booked-seats="showtime.bookedSeats"
            :selected-seats="selectedSeats"
            @seat-selected="toggleSeat"
          />
        </div>

        <!-- Booking Summary -->
        <div v-if="selectedSeats.length > 0" class="booking-summary">
          <div class="summary-content">
            <div class="selected-seats-info">
              <h3>Выбранные места:</h3>
              <div class="seats-list">
                <span
                  v-for="seat in selectedSeats"
                  :key="`${seat.row}-${seat.number}`"
                  class="seat-badge"
                >
                  {{ seat.row }}-{{ seat.number }}
                  <button class="remove-seat" @click="removeSeat(seat)">×</button>
                </span>
              </div>
            </div>

            <div class="total-price">
              <span class="price-label">Итого:</span>
              <span class="price-value">{{ formatCurrency(totalPrice) }}</span>
            </div>
          </div>

          <!-- Booking Timer -->
          <div class="booking-timer" :class="{ urgent: timerUrgent }">
            <span class="timer-icon">⏱</span>
            <span class="timer-label">Время на бронирование:</span>
            <span class="timer-value">{{ timerFormatted }}</span>
          </div>

          <div class="payment-method">
            <h3>Способ оплаты:</h3>
            <div class="payment-options">
              <label class="payment-option">
                <input v-model="paymentMethod" type="radio" value="wallet" name="payment" />
                <span class="option-content">
                  <span class="option-icon">💼</span>
                  <span class="option-text">
                    <strong>Кошелёк</strong>
                    <small>Баланс: {{ formatCurrency(userBalance) }}</small>
                  </span>
                </span>
              </label>

              <label class="payment-option">
                <input v-model="paymentMethod" type="radio" value="card" name="payment" />
                <span class="option-content">
                  <span class="option-icon">💳</span>
                  <span class="option-text">
                    <strong>Банковская карта</strong>
                    <small>Visa / Mastercard</small>
                  </span>
                </span>
              </label>

              <label class="payment-option">
                <input v-model="paymentMethod" type="radio" value="kaspi" name="payment" />
                <span class="option-content">
                  <span class="option-icon">📱</span>
                  <span class="option-text">
                    <strong>Kaspi QR</strong>
                    <small>Быстрая оплата</small>
                  </span>
                </span>
              </label>

              <label class="payment-option">
                <input v-model="paymentMethod" type="radio" value="cash" name="payment" />
                <span class="option-content">
                  <span class="option-icon">💵</span>
                  <span class="option-text">
                    <strong>Наличные</strong>
                    <small>Оплата в кассе</small>
                  </span>
                </span>
              </label>
            </div>

            <!-- Визуальная карта -->
            <div v-if="paymentMethod === 'card'" class="card-form">
              <div class="credit-card" :class="{ flipped: cardFlipped }" @click="cardFlipped = !cardFlipped">
                <div class="card-front">
                  <div class="card-chip">▬</div>
                  <div class="card-number-display">{{ cardNumberDisplay }}</div>
                  <div class="card-bottom">
                    <div>
                      <div class="card-label">Держатель</div>
                      <div class="card-name">{{ cardHolder || 'FULL NAME' }}</div>
                    </div>
                    <div>
                      <div class="card-label">Срок</div>
                      <div class="card-expiry">{{ cardExpiry || 'MM/YY' }}</div>
                    </div>
                    <div class="card-type">{{ cardType }}</div>
                  </div>
                </div>
                <div class="card-back">
                  <div class="card-stripe"></div>
                  <div class="card-cvv-wrap">
                    <span class="card-label">CVV</span>
                    <span class="card-cvv-val">{{ cardCvv || '•••' }}</span>
                  </div>
                </div>
              </div>

              <div class="card-inputs">
                <input
                  :value="cardNumber"
                  @input="cardNumber = formatCardNumber($event.target.value)"
                  placeholder="Номер карты"
                  class="card-input"
                  maxlength="19"
                />
                <input
                  v-model="cardHolder"
                  placeholder="Имя держателя"
                  class="card-input"
                  @input="cardHolder = $event.target.value.toUpperCase()"
                />
                <div class="card-input-row">
                  <input
                    :value="cardExpiry"
                    @input="cardExpiry = formatExpiry($event.target.value)"
                    placeholder="MM/YY"
                    class="card-input"
                    maxlength="5"
                  />
                  <input
                    v-model="cardCvv"
                    placeholder="CVV"
                    class="card-input"
                    maxlength="3"
                    @focus="cardFlipped = true"
                    @blur="cardFlipped = false"
                  />
                </div>
              </div>
            </div>

            <!-- Kaspi QR -->
            <div v-if="paymentMethod === 'kaspi'" class="kaspi-block">
              <div class="kaspi-header">
                <span class="kaspi-logo">🟡 Kaspi Pay</span>
                <span class="kaspi-amount">{{ formatCurrency(totalPrice) }}</span>
              </div>
              <div class="kaspi-qr-wrap">
                <div class="kaspi-qr-mock">
                  <div class="qr-corner tl"></div>
                  <div class="qr-corner tr"></div>
                  <div class="qr-corner bl"></div>
                  <div class="qr-corner br"></div>
                  <div class="qr-center">📱</div>
                </div>
              </div>
              <p class="kaspi-hint">Откройте Kaspi.kz → Платить → Сканируйте QR</p>
            </div>
          </div>

          <div v-if="error" class="alert alert-error">
            {{ error }}
          </div>

          <button
            class="btn btn-primary btn-large"
            :disabled="bookingInProgress"
            @click="createBooking"
          >
            <span v-if="bookingInProgress">Бронирование...</span>
            <span v-else>Забронировать {{ formatCurrency(totalPrice) }}</span>
          </button>
        </div>

        <div v-else class="no-seats-selected">
          <p>Выберите места для бронирования</p>
        </div>
      </div>

      <div v-else class="error-state">
        <h2>Ошибка загрузки</h2>
        <p>Не удалось загрузить информацию о сеансе</p>
        <button class="btn btn-primary" @click="goBack">
          Вернуться назад
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import SeatMap from '../components/SeatMap.vue'
import api from '../services/api'
import { formatDate, formatCurrency } from '../utils/formatters'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const bookingInProgress = ref(false)
const showtime = ref(null)
const hall = ref(null)
const movieTitle = ref('')
const cinemaName = ref('')
const selectedSeats = ref([])
const paymentMethod = ref('wallet')
const error = ref('')

// Таймер брони
const timerSeconds = ref(15 * 60)
const timerInterval = ref(null)
const timerFormatted = computed(() => {
  const m = Math.floor(timerSeconds.value / 60).toString().padStart(2, '0')
  const s = (timerSeconds.value % 60).toString().padStart(2, '0')
  return `${m}:${s}`
})
const timerUrgent = computed(() => timerSeconds.value < 120)

// Визуальная карта
const cardNumber = ref('')
const cardHolder = ref('')
const cardExpiry = ref('')
const cardCvv = ref('')
const cardFlipped = ref(false)

const formatCardNumber = (val) => {
  const digits = val.replace(/\D/g, '').slice(0, 16)
  return digits.replace(/(.{4})/g, '$1 ').trim()
}

const formatExpiry = (val) => {
  const digits = val.replace(/\D/g, '').slice(0, 4)
  if (digits.length >= 3) return digits.slice(0, 2) + '/' + digits.slice(2)
  return digits
}

const cardNumberDisplay = computed(() => {
  const v = cardNumber.value.replace(/\s/g, '')
  if (!v) return '•••• •••• •••• ••••'
  return formatCardNumber(v)
})

const cardType = computed(() => {
  const n = cardNumber.value.replace(/\s/g, '')
  if (n.startsWith('4')) return 'VISA'
  if (n.startsWith('5')) return 'MC'
  if (n.startsWith('6')) return 'Kaspi'
  return '••••'
})

const userBalance = computed(() => {
  return authStore.user?.wallet?.balance || 0
})

const totalPrice = computed(() => {
  return selectedSeats.value.reduce((sum, seat) => sum + seat.price, 0)
})

const fetchShowtimeDetails = async () => {
  loading.value = true
  try {
    const response = await api.getShowtimeById(route.params.showtimeId)
    const showtimeData = response.data.data

    if (!showtimeData) {
      throw new Error('Showtime not found')
    }

    showtime.value = showtimeData
    movieTitle.value = showtimeData.movieDetails?.title || 'Фильм'
    cinemaName.value = showtimeData.cinemaDetails?.name || 'Кинотеатр'

    // Использовать реальные места из зала если есть, иначе генерировать
    if (showtimeData.hallDetails?.seats?.length) {
      hall.value = showtimeData.hallDetails
    } else {
      hall.value = generateHallSeats(showtimeData)
    }
  } catch (err) {
    console.error('Failed to fetch showtime:', err)
    error.value = 'Не удалось загрузить информацию о сеансе'
  } finally {
    loading.value = false
  }
}

const generateHallSeats = (showtime) => {
  // Генерация мест зала (10 рядов по 15 мест)
  const rows = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J']
  const seatsPerRow = 15
  const seats = []

  rows.forEach((row, rowIndex) => {
    for (let number = 1; number <= seatsPerRow; number++) {
      let type = 'regular'
      let price = showtime.basePrice

      // VIP места (середина зала, средние ряды)
      if (rowIndex >= 3 && rowIndex <= 6 && number >= 5 && number <= 11) {
        type = 'vip'
        price = showtime.basePrice * 1.5
      }

      // Couple seats (последние 2 ряда, парные места)
      if (rowIndex >= 8 && number % 2 === 1 && number < seatsPerRow) {
        type = 'couple'
        price = showtime.basePrice * 1.3
      }

      seats.push({ row, number, type, price })
    }
  })

  return {
    _id: showtime.hallId,
    seats
  }
}

const toggleSeat = (seat) => {
  const index = selectedSeats.value.findIndex(
    s => s.row === seat.row && s.number === seat.number
  )

  if (index > -1) {
    selectedSeats.value.splice(index, 1)
  } else {
    if (selectedSeats.value.length >= 10) {
      error.value = 'Максимум 10 мест за одно бронирование'
      return
    }
    selectedSeats.value.push(seat)
    error.value = ''
  }
}

const removeSeat = (seat) => {
  const index = selectedSeats.value.findIndex(
    s => s.row === seat.row && s.number === seat.number
  )
  if (index > -1) {
    selectedSeats.value.splice(index, 1)
  }
}

const createBooking = async () => {
  if (selectedSeats.value.length === 0) {
    error.value = 'Выберите хотя бы одно место'
    return
  }

  if (paymentMethod.value === 'wallet' && totalPrice.value > userBalance.value) {
    error.value = `Недостаточно средств на кошельке. Необходимо: ${formatCurrency(totalPrice.value)}, доступно: ${formatCurrency(userBalance.value)}`
    return
  }

  bookingInProgress.value = true
  error.value = ''

  try {
    const bookingData = {
      showtimeId: route.params.showtimeId,
      seats: selectedSeats.value.map(s => ({
        row: s.row,
        number: s.number
      })),
      paymentMethod: paymentMethod.value
    }

    const response = await api.createBooking(bookingData)

    // Успешно создано
    const booking = response.data.data

    // Обновить профиль (баланс изменился)
    await authStore.fetchProfile()

    // Перейти в профиль
    router.push({
      name: 'Profile',
      query: { bookingCreated: booking._id || booking.id }
    })
  } catch (err) {
    error.value = err.response?.data?.error || 'Не удалось создать бронь'
    console.error('Booking failed:', err)
  } finally {
    bookingInProgress.value = false
  }
}

const goBack = () => {
  router.go(-1)
}

const startTimer = () => {
  timerInterval.value = setInterval(() => {
    if (timerSeconds.value > 0) {
      timerSeconds.value--
    } else {
      clearInterval(timerInterval.value)
      error.value = 'Время бронирования истекло. Пожалуйста, начните заново.'
    }
  }, 1000)
}

onMounted(() => {
  fetchShowtimeDetails()
  startTimer()
})

onUnmounted(() => {
  if (timerInterval.value) clearInterval(timerInterval.value)
})
</script>

<style scoped>
.seat-selection {
  padding: 40px 0;
  min-height: calc(100vh - 70px);
}

.booking-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.page-title {
  font-size: 32px;
  font-weight: bold;
}

.movie-info-card {
  background: linear-gradient(135deg, var(--dark-light), var(--dark-lighter));
  padding: 24px;
  border-radius: 16px;
  margin-bottom: 40px;
  border: 2px solid var(--dark-lighter);
}

.movie-title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 12px;
}

.session-info {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  font-size: 14px;
  color: var(--text-gray);
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.seat-selection-area {
  margin-bottom: 40px;
}

.booking-summary {
  background-color: var(--dark-light);
  padding: 30px;
  border-radius: 16px;
  border: 2px solid var(--primary);
}

.summary-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 20px;
}

.selected-seats-info h3 {
  font-size: 18px;
  margin-bottom: 12px;
}

.seats-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.seat-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: var(--primary);
  color: white;
  border-radius: 20px;
  font-weight: 600;
  font-size: 14px;
}

.remove-seat {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.remove-seat:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.total-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.price-label {
  font-size: 14px;
  color: var(--text-gray);
  margin-bottom: 4px;
}

.price-value {
  font-size: 36px;
  font-weight: bold;
  color: var(--primary);
}

/* ── Timer ── */
.booking-timer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: rgba(70, 211, 105, 0.1);
  border: 1px solid rgba(70, 211, 105, 0.3);
  border-radius: 12px;
  margin-bottom: 20px;
  font-size: 15px;
}
.booking-timer.urgent {
  background: rgba(229, 9, 20, 0.1);
  border-color: rgba(229, 9, 20, 0.4);
  animation: pulse 1s infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
.timer-icon { font-size: 20px; }
.timer-label { color: var(--text-gray); flex: 1; }
.timer-value { font-size: 22px; font-weight: 800; font-variant-numeric: tabular-nums; color: var(--success); }
.booking-timer.urgent .timer-value { color: var(--danger); }

/* ── Credit Card Visual ── */
.card-form { margin-top: 24px; }

.credit-card {
  width: 340px;
  height: 200px;
  margin: 0 auto 24px;
  position: relative;
  perspective: 1000px;
  cursor: pointer;
}

.card-front, .card-back {
  position: absolute;
  inset: 0;
  border-radius: 18px;
  padding: 24px;
  backface-visibility: hidden;
  transition: transform 0.6s ease;
}

.card-front {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 40%, #0f3460 100%);
  border: 1px solid rgba(255,255,255,0.1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-back {
  background: linear-gradient(135deg, #0f3460 0%, #16213e 100%);
  transform: rotateY(180deg);
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 20px;
}

.credit-card.flipped .card-front { transform: rotateY(-180deg); }
.credit-card.flipped .card-back { transform: rotateY(0deg); }

.card-chip {
  font-size: 28px;
  letter-spacing: -4px;
  color: #fbbf24;
}

.card-number-display {
  font-size: 20px;
  letter-spacing: 3px;
  font-family: 'Courier New', monospace;
  color: white;
  text-align: center;
}

.card-bottom {
  display: flex;
  align-items: flex-end;
  gap: 20px;
}

.card-label {
  font-size: 9px;
  color: rgba(255,255,255,0.5);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 2px;
}

.card-name, .card-expiry {
  font-size: 14px;
  color: white;
  font-weight: 600;
  letter-spacing: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120px;
}

.card-type {
  margin-left: auto;
  font-size: 18px;
  font-weight: 900;
  color: #fbbf24;
}

.card-stripe {
  height: 44px;
  background: rgba(0,0,0,0.7);
  margin: 0 -24px;
}

.card-cvv-wrap {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  padding: 0 12px;
}

.card-cvv-val {
  background: white;
  color: #111;
  padding: 6px 16px;
  border-radius: 6px;
  font-family: monospace;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 4px;
}

.card-inputs { display: flex; flex-direction: column; gap: 12px; }
.card-input {
  width: 100%;
  padding: 12px 16px;
  background: var(--dark);
  border: 2px solid var(--dark-lighter);
  border-radius: 10px;
  color: var(--text);
  font-size: 15px;
  transition: border-color 0.2s;
}
.card-input:focus { border-color: var(--primary); outline: none; }
.card-input-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

/* ── Kaspi QR ── */
.kaspi-block {
  margin-top: 20px;
  padding: 24px;
  background: #1a1a1a;
  border-radius: 16px;
  border: 2px solid #fbbf24;
  text-align: center;
}

.kaspi-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  font-size: 16px;
  font-weight: 700;
}

.kaspi-logo { color: #fbbf24; }
.kaspi-amount { color: var(--primary); font-size: 20px; }

.kaspi-qr-wrap { display: flex; justify-content: center; margin-bottom: 16px; }

.kaspi-qr-mock {
  width: 160px;
  height: 160px;
  background: white;
  border-radius: 12px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
}

.qr-corner {
  position: absolute;
  width: 28px;
  height: 28px;
  border: 4px solid #111;
}
.qr-corner.tl { top: 8px; left: 8px; border-right: none; border-bottom: none; }
.qr-corner.tr { top: 8px; right: 8px; border-left: none; border-bottom: none; }
.qr-corner.bl { bottom: 8px; left: 8px; border-right: none; border-top: none; }
.qr-corner.br { bottom: 8px; right: 8px; border-left: none; border-top: none; }
.qr-center { font-size: 48px; }

.kaspi-hint { font-size: 13px; color: var(--text-gray); }

.payment-method {
  margin-bottom: 24px;
}

.payment-method h3 {
  font-size: 18px;
  margin-bottom: 16px;
}

.payment-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.payment-option {
  position: relative;
  cursor: pointer;
}

.payment-option input[type="radio"] {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background-color: var(--dark-lighter);
  border: 2px solid var(--dark-lighter);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.payment-option input[type="radio"]:checked + .option-content {
  border-color: var(--primary);
  background-color: rgba(229, 9, 20, 0.1);
}

.option-icon {
  font-size: 32px;
}

.option-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.option-text strong {
  font-size: 16px;
}

.option-text small {
  font-size: 13px;
  color: var(--text-gray);
}

.btn-large {
  width: 100%;
  padding: 16px 32px;
  font-size: 18px;
  font-weight: bold;
}

.no-seats-selected {
  text-align: center;
  padding: 60px 20px;
  background-color: var(--dark-light);
  border-radius: 16px;
  color: var(--text-gray);
}

.error-state {
  text-align: center;
  padding: 80px 20px;
}

.error-state h2 {
  font-size: 32px;
  margin-bottom: 16px;
}

.error-state p {
  color: var(--text-gray);
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .summary-content {
    flex-direction: column;
    align-items: stretch;
  }

  .total-price {
    align-items: flex-start;
  }

  .payment-options {
    grid-template-columns: 1fr;
  }
}
</style>
