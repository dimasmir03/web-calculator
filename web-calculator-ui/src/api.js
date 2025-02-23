import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost/api/v1', // Укажите URL вашего бэкенда
})

export default {
  calculate(expression) {
    return api.post('/calculate', { expression })
  },
}
