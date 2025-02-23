<!-- eslint-disable vue/multi-word-component-names -->
<script setup>
import { ref } from 'vue'
import api from '@/api'

const expression = ref('')
const result = ref(null)
const error = ref('')
const isLoading = ref(false)

const calculate = async () => {
  error.value = ''
  result.value = null
  isLoading.value = true

  if (!expression.value) {
    error.value = 'Пожалуйста, введите выражение.'
    isLoading.value = false
    return
  }

  try {
    const response = await api.calculate(expression.value)
    result.value = response.data.result
  } catch (err) {
    error.value = 'Ошибка при вычислении выражения.'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="calculator">
    <input v-model="expression" placeholder="Введите выражение" />
    <button @click="calculate" :disabled="isLoading">
      {{ isLoading ? 'Вычисление...' : 'Вычислить' }}
    </button>
    <div v-if="result !== null">
      <h3>Результат: {{ result }}</h3>
    </div>
    <div v-if="error">
      <p style="color: red">{{ error }}</p>
    </div>
  </div>
</template>

<style scoped>
.calculator {
  margin: 20px;
}

input {
  padding: 10px;
  font-size: 16px;
  margin-right: 10px;
}

button {
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
}
</style>
