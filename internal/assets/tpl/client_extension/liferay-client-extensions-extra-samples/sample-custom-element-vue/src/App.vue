<script setup>
import {ref} from 'vue'
import LiferayService from "./services/LiferayService.js";

const user = ref(null)
const error = ref(false)
const errorMessage = ref('')
const loading = ref(true)

LiferayService.get('/o/headless-admin-user/v1.0/my-user-account')
.then((response) => {
  user.value = response
  loading.value = false
})
.catch((err) => {
  error.value = true
  errorMessage.value = err.message
})
.finally(() => {
  loading.value = false
})
</script>

<template>
  <span
      v-if="loading"
      aria-hidden="true"
      class="loading-animation-squares loading-animation-primary loading-animation-md"
  ></span>

  <div
      v-if="error"
      class="alert alert-danger">
    <strong class="lead">Error:</strong>{errorMessage}
  </div>
  <h1
      v-else
      style="color: var(--primary)">
    Hello from Vue.js{{ ', ' + user.name + '!' || '!' }}
  </h1>
</template>
